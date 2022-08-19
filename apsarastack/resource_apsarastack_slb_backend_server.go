package apsarastack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackSlbBackendServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackSlbBackendServersCreate,
		Read:   resourceApsaraStackSlbBackendServersRead,
		Update: resourceApsaraStackSlbBackendServersUpdate,
		Delete: resourceApsaraStackSlbBackendServersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backend_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(0, 100),
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

func resourceApsaraStackSlbBackendServersCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := slb.CreateAddBackendServersRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	if v, ok := d.GetOk("backend_servers"); ok {
		request.BackendServers = expandBackendServersInfoToString(v.(*schema.Set).List())
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.AddBackendServers(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_slb_backend_servers", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.AddBackendServersResponse)
	d.SetId(response.LoadBalancerId)

	return resourceApsaraStackSlbBackendServersRead(d, meta)
}

func resourceApsaraStackSlbBackendServersRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}
	resource_id := d.Id()
	object, err := slbService.DescribeSlb(resource_id)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("load_balancer_id", object.LoadBalancerId)

	servers := make([]map[string]interface{}, 0)

	for _, server := range object.BackendServers.BackendServer {
		s := map[string]interface{}{
			"server_id": server.ServerId,
			"weight":    server.Weight,
		}
		servers = append(servers, s)
	}

	if err := d.Set("backend_servers", servers); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceApsaraStackSlbBackendServersUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	d.Partial(true)
	step := 20
	var removeSet, addSet, updateSet *schema.Set

	if d.HasChange("backend_servers") {
		o, n := d.GetChange("backend_servers")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		oldIds := getIdSetFromServers(remove)
		newIds := getIdSetFromServers(add)
		updateSet = oldIds.Intersection(newIds)
		addSet = newIds.Difference(oldIds)
		removeSet = oldIds.Difference(newIds)

		if removeSet.Len() > 0 {
			rmservers := make([]interface{}, 0)
			for _, rmserver := range remove {
				rms := rmserver.(map[string]interface{})
				if removeSet.Contains(rms["server_id"]) {
					rmsm := map[string]interface{}{
						"server_id": rms["server_id"],
						"weight":    rms["weight"],
					}
					rmservers = append(rmservers, rmsm)
				}
			}
			request := slb.CreateRemoveBackendServersRequest()
			request.RegionId = client.RegionId
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			request.LoadBalancerId = d.Id()

			segs := len(rmservers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(rmservers) {
					end = len(rmservers)
				}
				request.BackendServers = expandBackendServersInfoToString(rmservers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.RemoveBackendServers(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				//d.SetPartial("backend_servers")
			}

		}

		if addSet.Len() > 0 {
			addservers := make([]interface{}, 0)
			for _, addserver := range add {
				adds := addserver.(map[string]interface{})
				if addSet.Contains(adds["server_id"]) {
					addsm := map[string]interface{}{
						"server_id": adds["server_id"],
						"weight":    adds["weight"],
					}
					addservers = append(addservers, addsm)
				}
			}
			request := slb.CreateAddBackendServersRequest()
			request.RegionId = client.RegionId
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			request.LoadBalancerId = d.Id()

			segs := len(addservers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(addservers) {
					end = len(addservers)
				}
				request.BackendServers = expandBackendServersInfoToString(addservers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.AddBackendServers(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				//d.SetPartial("backend_servers")
			}
		}

		servers := make([]interface{}, 0)
		for _, server := range d.Get("backend_servers").(*schema.Set).List() {
			s := server.(map[string]interface{})
			if updateSet.Contains(s["server_id"]) {
				sm := map[string]interface{}{
					"server_id": s["server_id"],
					"weight":    s["weight"],
				}
				servers = append(servers, sm)
			}
		}

		if len(servers) > 0 {

			segs := len(servers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(servers) {
					end = len(servers)
				}
				request := slb.CreateSetBackendServersRequest()
				request.RegionId = client.RegionId
				if strings.ToLower(client.Config.Protocol) == "https" {
					request.Scheme = "https"
				} else {
					request.Scheme = "http"
				}
				request.Headers = map[string]string{"RegionId": client.RegionId}
				request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
				request.LoadBalancerId = d.Id()
				request.BackendServers = expandBackendServersInfoToString(servers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.SetBackendServers(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				//d.SetPartial("backend_servers")
			}
		}
	}
	d.Partial(false)

	return resourceApsaraStackSlbBackendServersRead(d, meta)
}

func resourceApsaraStackSlbBackendServersDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	instanceSet := d.Get("backend_servers").(*schema.Set)
	step := 20
	if len(instanceSet.List()) > 0 {

		slbService := SlbService{client}
		if d.Get("delete_protection_validation").(bool) {
			lbInstance, err := slbService.DescribeSlb(d.Id())
			if err != nil {
				if NotFoundError(err) {
					return nil
				}
				return WrapError(err)
			}
			if lbInstance.DeleteProtection == "on" {
				return WrapError(fmt.Errorf("Current backend servers' SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", d.Id()))
			}
		}

		servers := make([]interface{}, 0)
		for _, rmserver := range instanceSet.List() {
			rms := rmserver.(map[string]interface{})
			rmsm := map[string]interface{}{
				"server_id": rms["server_id"],
				"weight":    rms["weight"],
			}
			servers = append(servers, rmsm)
		}

		request := slb.CreateRemoveBackendServersRequest()
		request.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.LoadBalancerId = d.Id()

		segs := len(servers)/step + 1
		for i := 0; i < segs; i++ {
			start := i * step
			end := (i + 1) * step
			if end >= len(servers) {
				end = len(servers)
			}

			request.BackendServers = expandBackendServersWithoutTypeToString(servers[start:end])
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.RemoveBackendServers(request)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"RspoolVipExist", "ObtainIpFail"}) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
			}
		}
	}

	return nil
}
