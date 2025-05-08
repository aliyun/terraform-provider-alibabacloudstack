package alibabacloudstack

import (
	"encoding/json"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEssAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 20,
				MinItems: 1,
			},

			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackEssAttachmentCreate,
		resourceAlibabacloudStackEssAttachmentRead,
		resourceAlibabacloudStackEssAttachmentUpdate,
		resourceAlibabacloudStackEssAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackEssAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("scaling_group_id").(string))
	return nil
}

func resourceAlibabacloudStackEssAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	d.Partial(true)

	if d.HasChange("instance_ids") {
		object, err := essService.DescribeEssScalingGroup(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if object.LifecycleState == string(Inactive) {
			return errmsgs.WrapError(errmsgs.Error("Scaling group current status is %s, please active it before attaching or removing ECS instances.", object.LifecycleState))
		} else {
			if err := essService.WaitForEssScalingGroup(object.ScalingGroupId, Active, DefaultTimeout); err != nil {
				return errmsgs.WrapError(err)
			}
		}
		o, n := d.GetChange("instance_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := convertArrayInterfaceToArrayString(ns.Difference(os).List())

		if len(add) > 0 {
			request := ess.CreateAttachInstancesRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.ScalingGroupId = d.Id()
			request.InstanceId = &add
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.AttachInstances(request)
				})
				response, ok := raw.(*ess.AttachInstancesResponse)
				if err != nil {
					errmsg := ""
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
					raw_data := make(map[string]interface{})
					err := json.Unmarshal(response.BaseResponse.GetHttpContentBytes(), &raw_data)
					errCode := raw_data["errorCode"].(string)
					if errCode == "IncorrectCapacity.MaxSize" {
						instances, err := essService.DescribeEssAttachment(d.Id(), make([]string, 0))
						if !errmsgs.NotFoundError(err) {
							return resource.NonRetryableError(err)
						}
						var autoAdded, attached []string
						if len(instances) > 0 {
							for _, inst := range instances {
								if inst.CreationType == "Attached" {
									attached = append(attached, inst.InstanceId)
								} else {
									autoAdded = append(autoAdded, inst.InstanceId)
								}
							}
						}
						if len(add) > object.MaxSize {
							return resource.NonRetryableError(errmsgs.WrapError(errmsgs.Error("To attach %d instances, the total capacity will be greater than the scaling group max size %d. "+
								"Please enlarge scaling group max size.", len(add), object.MaxSize)))
						}

						if len(autoAdded) > 0 {
							if d.Get("force").(bool) {
								if err := essService.EssRemoveInstances(d.Id(), autoAdded); err != nil {
									return resource.NonRetryableError(errmsgs.WrapError(err))
								}
								time.Sleep(5)
								return resource.RetryableError(errmsgs.WrapError(err))
							} else {
								return resource.NonRetryableError(errmsgs.WrapError(errmsgs.Error("To attach the instances, the total capacity will be greater than the scaling group max size %d."+
									"Please enlarge scaling group max size or set 'force' to true to remove autocreated instances: %#v.", object.MaxSize, autoAdded)))
							}
						}

						if len(attached) > 0 {
							return resource.NonRetryableError(errmsgs.WrapError(errmsgs.Error("To attach the instances, the total capacity will be greater than the scaling group max size %d. "+
								"Please enlarge scaling group max size or remove already attached instances: %#v.", object.MaxSize, attached)))
						}
					}
					if errCode == "ScalingActivityInProgress" {
						time.Sleep(5)
						return resource.RetryableError(errmsgs.WrapError(err))
					}
					return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errmsgs.WrapError(err)
			}

			err = resource.Retry(3*time.Minute, func() *resource.RetryError {
				instances, err := essService.DescribeEssAttachment(d.Id(), add)
				if err != nil {
					return resource.NonRetryableError(errmsgs.WrapError(err))
				}
				if len(instances) < 0 {
					return resource.RetryableError(errmsgs.WrapError(errmsgs.Error("There are no ECS instances have been attached.")))
				}

				for _, inst := range instances {
					if inst.LifecycleState != string(InService) {
						return resource.RetryableError(errmsgs.WrapError(errmsgs.Error("There are still ECS instances are not %s.", string(InService))))
					}
				}
				return nil
			})
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
		if len(remove) > 0 {
			if err := essService.EssRemoveInstances(d.Id(), convertArrayInterfaceToArrayString(remove)); err != nil {
				return errmsgs.WrapError(err)
			}
		}

		//d.SetPartial("instance_ids")
	}

	d.Partial(false)

	return nil
}

func resourceAlibabacloudStackEssAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	object, err := essService.DescribeEssAttachment(d.Id(), make([]string, 0))

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	var instanceIds []string
	for _, inst := range object {
		instanceIds = append(instanceIds, inst.InstanceId)
	}

	d.Set("scaling_group_id", object[0].ScalingGroupId)
	d.Set("instance_ids", instanceIds)

	return nil
}

func resourceAlibabacloudStackEssAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	removed := convertArrayInterfaceToArrayString(d.Get("instance_ids").(*schema.Set).List())

	if len(removed) < 1 {
		return nil
	}
	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if err := essService.WaitForEssScalingGroup(object.ScalingGroupId, Active, DefaultTimeout); err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := ess.CreateRemoveInstancesRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.ScalingGroupId = d.Id()

		if len(removed) > 0 {
			request.InstanceId = &removed
		} else {
			return nil
		}
		raw, err := essService.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.RemoveInstances(request)
		})
		response, ok := raw.(*ess.RemoveInstancesResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			raw_data := make(map[string]interface{})
			err := json.Unmarshal(response.BaseResponse.GetHttpContentBytes(), &raw_data)
			errmsg := ""
			if err != nil && ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			errCode := raw_data["errorCode"].(string)
			if errCode == "IncorrectCapacity.MinSize" {
				instances, err := essService.DescribeEssAttachment(d.Id(), removed)
				if len(instances) > 0 {
					if object.MinSize == 0 {
						return resource.RetryableError(errmsgs.WrapError(err))
					}
					return resource.NonRetryableError(errmsgs.WrapError(errmsgs.Error("To remove %d instances, the total capacity will be lesser than the scaling group min size %d. "+
						"Please shorten scaling group min size and try again.", len(removed), object.MinSize)))
				}
			}
			if errCode == "ScalingActivityInProgress" || errCode == "IncorrectScalingGroupStatus" {
				time.Sleep(5)
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			if errCode == "InvalidScalingGroupId.NotFound" {
				return nil
			}
			if err != nil {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
		}
		time.Sleep(3 * time.Second)
		instances, err := essService.DescribeEssAttachment(d.Id(), removed)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		if len(instances) > 0 {
			removed = make([]string, 0)
			for _, inst := range instances {
				removed = append(removed, inst.InstanceId)
			}
			return resource.RetryableError(errmsgs.WrapError(errmsgs.Error("There are still ECS instances in the scaling group.")))
		}

		return nil
	}); err != nil {
		return errmsgs.WrapError(err)
	}

	return errmsgs.WrapError(essService.WaitForEssAttachment(d.Id(), Deleted, DefaultTimeout))
}

func convertArrayInterfaceToArrayString(elm []interface{}) (arr []string) {
	if len(elm) < 1 {
		return
	}
	for _, e := range elm {
		arr = append(arr, e.(string))
	}
	return
}
