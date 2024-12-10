package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkInterfaceCreate,
		Read:   resourceNetworkInterfaceRead,
		Update: resourceNetworkInterfaceUpdate,
		Delete: resourceNetworkInterfaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:true,
				Computed:true,
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'network_interface_name' instead.",
				ConflictsWith: []string{"network_interface_name"},
			},
			"network_interface_name": {
				Type:          schema.TypeString,
				Optional:true,
				Computed:true,
				ConflictsWith: []string{"name"},
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				MinItems: 1,
			},

			"private_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'private_ip' is deprecated and will be removed in a future release. Please use new field 'primary_ip_address' instead.",
				ConflictsWith: []string{"primary_ip_address"},
			},
			"primary_ip_address": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"private_ip"},
			},
			"private_ips": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				MaxItems:      10,
				ConflictsWith: []string{"private_ips_count"},
			},
			"private_ips_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.IntBetween(0, 10),
				ConflictsWith: []string{"private_ips"},
			},
			"mac": {
				Type:          schema.TypeString,
				Computed:      true,
				Deprecated:    "Field 'mac' is deprecated and will be removed in a future release. Please use new field 'mac_address' instead.",
			},
			"mac_address": {
				Type:          schema.TypeString,
				Computed:      true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	request := ecs.CreateCreateNetworkInterfaceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VSwitchId = d.Get("vswitch_id").(string)
	groups := d.Get("security_groups").(*schema.Set).List()

	request.SecurityGroupId = groups[0].(string)

	if primaryIpAddress, ok := connectivity.GetResourceDataOk(d, "primary_ip_address", "private_ip"); ok {
		request.PrimaryIpAddress = primaryIpAddress.(string)
	}
	if name, ok := connectivity.GetResourceDataOk(d, "network_interface_name", "name"); ok {
		request.NetworkInterfaceName = name.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateNetworkInterface(request)
	})
	bresponse, ok := raw.(*ecs.CreateNetworkInterfaceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_network_interface", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(bresponse.NetworkInterfaceId)

	if err := ecsService.WaitForNetworkInterface(d.Id(), Available, 600); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceNetworkInterfaceUpdate(d, meta)
}

func resourceNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeNetworkInterface(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.NetworkInterfaceName, "network_interface_name", "name")
	d.Set("description", object.Description)
	d.Set("vswitch_id", object.VSwitchId)
	connectivity.SetResourceData(d, object.PrivateIpAddress, "primary_ip_address", "private_ip")
	d.Set("security_groups", object.SecurityGroupIds.SecurityGroupId)
	privateIps := make([]string, 0, len(object.PrivateIpSets.PrivateIpSet))
	for i := range object.PrivateIpSets.PrivateIpSet {
		if !object.PrivateIpSets.PrivateIpSet[i].Primary {
			privateIps = append(privateIps, object.PrivateIpSets.PrivateIpSet[i].PrivateIpAddress)
		}
	}
	d.Set("private_ips", privateIps)
	d.Set("private_ips_count", len(privateIps))
	connectivity.SetResourceData(d, object.MacAddress, "mac_address", "mac")

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceEni)
	if err != nil && !errmsgs.NotFoundError(err) {
		return errmsgs.WrapError(err)
	}

	if len(tags) > 0 {
		d.Set("tags", ecsService.tagsToMap(tags))
	}

	return nil
}

func resourceNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	d.Partial(true)

	attributeUpdate := false
	request := ecs.CreateModifyNetworkInterfaceAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.NetworkInterfaceId = d.Id()
	if !d.IsNewResource() && d.HasChange("description") {
		request.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if !d.IsNewResource() && (d.HasChange("network_interface_name") || d.HasChange("name")) {
		request.NetworkInterfaceName = connectivity.GetResourceData(d, "network_interface_name", "name").(string)
		attributeUpdate = true
	}

	if d.HasChange("security_groups") {
		securityGroups := expandStringList(d.Get("security_groups").(*schema.Set).List())
		if len(securityGroups) > 1 || !d.IsNewResource() {
			request.SecurityGroupId = &securityGroups
			attributeUpdate = true
		}
	}

	if attributeUpdate {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyNetworkInterfaceAttribute(request)
		})
		bresponse, ok := raw.(*ecs.ModifyNetworkInterfaceAttributeResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("private_ips") {
		oldIps, newIps := d.GetChange("private_ips")
		oldIpsSet := oldIps.(*schema.Set)
		newIpsSet := newIps.(*schema.Set)

		unAssignIps := oldIpsSet.Difference(newIpsSet)
		if unAssignIps.Len() > 0 {
			unAssignIpList := expandStringList(unAssignIps.List())
			unAssignPrivateIpAddressesRequest := ecs.CreateUnassignPrivateIpAddressesRequest()
			client.InitRpcRequest(*unAssignPrivateIpAddressesRequest.RpcRequest)
			unAssignPrivateIpAddressesRequest.NetworkInterfaceId = d.Id()
			unAssignPrivateIpAddressesRequest.PrivateIpAddress = &unAssignIpList
			err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.UnassignPrivateIpAddresses(unAssignPrivateIpAddressesRequest)
				})
				bresponse, ok := raw.(*ecs.UnassignPrivateIpAddressesResponse)
				if err != nil {
					if errmsgs.IsExpectedErrors(err, errmsgs.NetworkInterfaceInvalidOperations) {
						return resource.RetryableError(err)
					}
					errmsg := ""
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
					}
					return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), unAssignPrivateIpAddressesRequest.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure, errmsg))
				}
				addDebug(unAssignPrivateIpAddressesRequest.GetActionName(), raw, unAssignPrivateIpAddressesRequest.RpcRequest, unAssignPrivateIpAddressesRequest)
				return nil
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), unAssignPrivateIpAddressesRequest.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure)
			}
		}

		assignIps := newIpsSet.Difference(oldIpsSet)
		if assignIps.Len() > 0 {
			assignIpList := expandStringList(assignIps.List())
			assignPrivateIpAddressesRequest := ecs.CreateAssignPrivateIpAddressesRequest()
			client.InitRpcRequest(*assignPrivateIpAddressesRequest.RpcRequest)
			assignPrivateIpAddressesRequest.NetworkInterfaceId = d.Id()
			assignPrivateIpAddressesRequest.PrivateIpAddress = &assignIpList
			err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.AssignPrivateIpAddresses(assignPrivateIpAddressesRequest)
				})
				bresponse, ok := raw.(*ecs.AssignPrivateIpAddressesResponse)
				if err != nil {
					errmsg := ""
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
					}
					return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure, errmsg))
				}
				addDebug(assignPrivateIpAddressesRequest.GetActionName(), raw, assignPrivateIpAddressesRequest.RpcRequest, assignPrivateIpAddressesRequest)
				return nil
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure)
			}
		}

		if err := ecsService.WaitForPrivateIpsListChanged(d.Id(), expandStringList(newIpsSet.List())); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChange("private_ips_count") {
		privateIpList := expandStringList(d.Get("private_ips").(*schema.Set).List())
		oldIpsCount, newIpsCount := d.GetChange("private_ips_count")
		if oldIpsCount != nil && newIpsCount != nil && newIpsCount != len(privateIpList) {
			diff := newIpsCount.(int) - oldIpsCount.(int)
			if diff > 0 {
				assignPrivateIpAddressesRequest := ecs.CreateAssignPrivateIpAddressesRequest()
				client.InitRpcRequest(*assignPrivateIpAddressesRequest.RpcRequest)
				assignPrivateIpAddressesRequest.NetworkInterfaceId = d.Id()
				assignPrivateIpAddressesRequest.SecondaryPrivateIpAddressCount = requests.NewInteger(diff)
				err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
					raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.AssignPrivateIpAddresses(assignPrivateIpAddressesRequest)
					})
					bresponse, ok := raw.(*ecs.AssignPrivateIpAddressesResponse)
					if err != nil {
						if errmsgs.IsExpectedErrors(err, errmsgs.NetworkInterfaceInvalidOperations) {
							errmsg := ""
							if ok {
								errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
							}
							return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure, errmsg))
						}
						errmsg := ""
						if ok {
							errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
						}
						return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure, errmsg))
					}
					addDebug(assignPrivateIpAddressesRequest.GetActionName(), raw, assignPrivateIpAddressesRequest.RpcRequest, assignPrivateIpAddressesRequest)
					return nil
				})
				if err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure)
				}
			}

			if diff < 0 {
				diff *= -1
				unAssignIps := privateIpList[:diff]
				unAssignPrivateIpAddressesRequest := ecs.CreateUnassignPrivateIpAddressesRequest()
				client.InitRpcRequest(*unAssignPrivateIpAddressesRequest.RpcRequest)
				unAssignPrivateIpAddressesRequest.NetworkInterfaceId = d.Id()
				unAssignPrivateIpAddressesRequest.PrivateIpAddress = &unAssignIps
				err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
					raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.UnassignPrivateIpAddresses(unAssignPrivateIpAddressesRequest)
					})
					if err != nil {
						if errmsgs.IsExpectedErrors(err, errmsgs.NetworkInterfaceInvalidOperations) {
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(unAssignPrivateIpAddressesRequest.GetActionName(), raw, unAssignPrivateIpAddressesRequest.RpcRequest, unAssignPrivateIpAddressesRequest)
					return nil
				})
				if err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), unAssignPrivateIpAddressesRequest.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure)
				}
			}

			err := ecsService.WaitForPrivateIpsCountChanged(d.Id(), newIpsCount.(int))
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}

	if err := setTags(client, TagResourceEni, d); err != nil {
		return errmsgs.WrapError(err)
	}

	d.Partial(false)

	return resourceNetworkInterfaceRead(d, meta)
}

func resourceNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteNetworkInterfaceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.NetworkInterfaceId = d.Id()

	err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteNetworkInterface(request)
		})
		bresponse, ok := raw.(*ecs.DeleteNetworkInterfaceResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure)
	}
	return errmsgs.WrapError(ecsService.WaitForNetworkInterface(d.Id(), Deleted, DefaultTimeoutMedium))
}
