package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlbMasterSlaveServerGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'master_slave_server_group_name' instead.",
				ConflictsWith: []string{"master_slave_server_group_name"},
			},

			"master_slave_server_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},

			"servers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"server_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Master", "Slave"}, false),
						},
					},
				},
				MaxItems: 2,
				MinItems: 2,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackSlbMasterSlaveServerGroupCreate,
		resourceAlibabacloudStackSlbMasterSlaveServerGroupRead,
		resourceAlibabacloudStackSlbMasterSlaveServerGroupUpdate,
		resourceAlibabacloudStackSlbMasterSlaveServerGroupDelete)
	return resource
}

func resourceAlibabacloudStackSlbMasterSlaveServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := slb.CreateCreateMasterSlaveServerGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	if v, ok := connectivity.GetResourceDataOk(d, "master_slave_server_group_name", "name"); ok {
		request.MasterSlaveServerGroupName = v.(string)
	}
	if v, ok := d.GetOk("servers"); ok {
		request.MasterSlaveBackendServers = expandMasterSlaveBackendServersToString(v.(*schema.Set).List())
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateMasterSlaveServerGroup(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*slb.CreateMasterSlaveServerGroupResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_master_slave_server_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateMasterSlaveServerGroupResponse)
	d.SetId(response.MasterSlaveServerGroupId)

	return nil
}

func resourceAlibabacloudStackSlbMasterSlaveServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbMasterSlaveServerGroup(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.MasterSlaveServerGroupName, "master_slave_server_group_name", "name")
	if object.LoadBalancerId != "" {
		// 在专有云的实际环境中可能不会返回相关值导致每次apply都会去销毁资源后重建
		d.Set("load_balancer_id", object.LoadBalancerId)
	}

	servers := make([]map[string]interface{}, 0)

	for _, server := range object.MasterSlaveBackendServers.MasterSlaveBackendServer {
		s := map[string]interface{}{
			"server_id":   server.ServerId,
			"port":        server.Port,
			"weight":      server.Weight,
			"server_type": server.ServerType,
		}
		servers = append(servers, s)
	}

	if err := d.Set("servers", servers); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackSlbMasterSlaveServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	noUpdateAllowedFields := []string{"delete_protection_validation"}
	return noUpdatesAllowedCheck(d, noUpdateAllowedFields)
}

func resourceAlibabacloudStackSlbMasterSlaveServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
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
			return errmsgs.WrapError(fmt.Errorf("Current master-slave server group's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", lbId))
		}
	}

	request := slb.CreateDeleteMasterSlaveServerGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.MasterSlaveServerGroupId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteMasterSlaveServerGroup(request)
		})
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*slb.DeleteMasterSlaveServerGroupResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"RspoolVipExist"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"The specified MasterSlaveGroupId does not exist", "InvalidParameter"}) {
			return nil
		}
		return err
	}

	return errmsgs.WrapError(slbService.WaitForSlbMasterSlaveServerGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
