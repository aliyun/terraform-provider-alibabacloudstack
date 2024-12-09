package alibabacloudstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSlbServerGroupCreate,
		Read:   resourceAlibabacloudStackSlbServerGroupRead,
		Update: resourceAlibabacloudStackSlbServerGroupUpdate,
		Delete: resourceAlibabacloudStackSlbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "tf-server-group",
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'vserver_group_name' instead.",
				ConflictsWith: []string{"vserver_group_name"},
			},
			"vserver_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "tf-server-group",
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							MinItems: 1,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(ECS),
							ValidateFunc: validation.StringInSlice([]string{"eni", "ecs"}, false),
						},
					},
				},
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlibabacloudStackSlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := slb.CreateCreateVServerGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	if v, ok := connectivity.GetResourceDataOk(d, "vserver_group_name", "name"); ok {
		request.VServerGroupName = v.(string)
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateVServerGroup(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*slb.CreateVServerGroupResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_server_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*slb.CreateVServerGroupResponse)
	d.SetId(response.VServerGroupId)
	d.Set("load_balancer_id", d.Get("load_balancer_id").(string))
	return resourceAlibabacloudStackSlbServerGroupUpdate(d, meta)
}

func resourceAlibabacloudStackSlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbServerGroup(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.VServerGroupName, "vserver_group_name", "name")
	d.Set("load_balancer_id", object.LoadBalancerId)

	servers := make([]map[string]interface{}, 0)
	portAndWeight := make(map[string][]string)
	for _, server := range object.BackendServers.BackendServer {
		key := fmt.Sprintf("%d%s%d%s%s", server.Port, COLON_SEPARATED, server.Weight, COLON_SEPARATED, server.Type)
		if v, ok := portAndWeight[key]; !ok {
			portAndWeight[key] = []string{server.ServerId}
		} else {
			v = append(v, server.ServerId)
			portAndWeight[key] = v
		}
	}
	for key, value := range portAndWeight {
		k := strings.Split(key, COLON_SEPARATED)
		p, e := strconv.Atoi(k[0])
		if e != nil {
			return errmsgs.WrapError(e)
		}
		w, e := strconv.Atoi(k[1])
		if e != nil {
			return errmsgs.WrapError(e)
		}
		t := k[2]
		s := map[string]interface{}{
			"server_ids": value,
			"port":       p,
			"weight":     w,
			"type":       t,
		}
		servers = append(servers, s)
	}

	if err := d.Set("servers", servers); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackSlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	d.Partial(true)
	var removeserverSet, addServerSet, updateServerSet *schema.Set
	serverUpdate := false
	step := 20
	if d.HasChange("servers") {
		o, n := d.GetChange("servers")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()
		oldIdPort := getIdPortSetFromServers(remove)
		newIdPort := getIdPortSetFromServers(add)
		updateServerSet = oldIdPort.Intersection(newIdPort)
		removeserverSet = oldIdPort.Difference(newIdPort)
		addServerSet = newIdPort.Difference(oldIdPort)
		if removeserverSet.Len() > 0 {
			rmservers := make([]interface{}, 0)
			for _, rmserver := range remove {
				rms := rmserver.(map[string]interface{})
				if v, ok := rms["server_ids"]; ok {
					server_ids := v.([]interface{})
					for _, id := range server_ids {
						idPort := fmt.Sprintf("%s:%d", id, rms["port"])
						if removeserverSet.Contains(idPort) {
							rmsm := map[string]interface{}{
								"server_id": id,
								"port":      rms["port"],
								"type":      rms["type"],
								"weight":    rms["weight"],
							}
							rmservers = append(rmservers, rmsm)
						}
					}
				}
			}
			request := slb.CreateRemoveVServerGroupBackendServersRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.VServerGroupId = d.Id()
			segs := len(rmservers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(rmservers) {
					end = len(rmservers)
				}
				request.BackendServers = expandBackendServersWithPortToString(rmservers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.RemoveVServerGroupBackendServers(request)
				})
				if err != nil {
					errmsg := ""
					if response, ok := raw.(*slb.RemoveVServerGroupBackendServersResponse); ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}
		if addServerSet.Len() > 0 {
			addservers := make([]interface{}, 0)
			for _, addserver := range add {
				adds := addserver.(map[string]interface{})
				if v, ok := adds["server_ids"]; ok {
					server_ids := v.([]interface{})
					for _, id := range server_ids {
						idPort := fmt.Sprintf("%s:%d", id, adds["port"])
						if addServerSet.Contains(idPort) {
							addsm := map[string]interface{}{
								"server_id": id,
								"port":      adds["port"],
								"type":      adds["type"],
								"weight":    adds["weight"],
							}
							addservers = append(addservers, addsm)
						}
					}
				}
			}
			request := slb.CreateAddVServerGroupBackendServersRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.VServerGroupId = d.Id()
			segs := len(addservers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(addservers) {
					end = len(addservers)
				}
				request.BackendServers = expandBackendServersWithPortToString(addservers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.AddVServerGroupBackendServers(request)
				})
				if err != nil {
					errmsg := ""
					if response, ok := raw.(*slb.AddVServerGroupBackendServersResponse); ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}
	}
	name := connectivity.GetResourceData(d, "vserver_group_name", "name").(string)
	nameUpdate := false
	if d.HasChanges("name","vserver_group_name") {
		nameUpdate = true
	}
	if d.HasChange("servers") {
		serverUpdate = true
	}
	if serverUpdate || nameUpdate {
		request := slb.CreateSetVServerGroupAttributeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.VServerGroupId = d.Id()
		if nameUpdate{
			request.VServerGroupName = name
		}
		if serverUpdate {
			servers := make([]interface{}, 0)
			for _, server := range d.Get("servers").(*schema.Set).List() {
				s := server.(map[string]interface{})
				if v, ok := s["server_ids"]; ok {
					server_ids := v.([]interface{})
					for _, id := range server_ids {
						idPort := fmt.Sprintf("%s:%d", id, s["port"])
						if updateServerSet.Contains(idPort) {
							sm := map[string]interface{}{
								"server_id": id,
								"port":      s["port"],
								"type":      s["type"],
								"weight":    s["weight"],
							}
							servers = append(servers, sm)
						}
					}
				}
			}
			segs := len(servers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(servers) {
					end = len(servers)
				}
				request.BackendServers = expandBackendServersWithPortToString(servers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.SetVServerGroupAttribute(request)
				})
				if err != nil {
					errmsg := ""
					if response, ok := raw.(*slb.SetVServerGroupAttributeResponse); ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		} else {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetVServerGroupAttribute(request)
			})
			if err != nil {
				errmsg := ""
				if response, ok := raw.(*slb.SetVServerGroupAttributeResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
	}
	d.Partial(false)

	return resourceAlibabacloudStackSlbServerGroupRead(d, meta)
}

func resourceAlibabacloudStackSlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return errmsgs.WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return errmsgs.WrapError(fmt.Errorf("Current VServerGroup's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the group.", lbId))
		}
	}

	request := slb.CreateDeleteVServerGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VServerGroupId = d.Id()
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteVServerGroup(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"RspoolVipExist"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if response, ok := raw.(*slb.DeleteVServerGroupResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"The specified VServerGroupId does not exist", "InvalidParameter"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return errmsgs.WrapError(slbService.WaitForSlbServerGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
