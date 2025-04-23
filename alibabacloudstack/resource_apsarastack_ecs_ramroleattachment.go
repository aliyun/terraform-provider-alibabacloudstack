package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackRamRoleAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, 
		resourceAlibabacloudStackInstanceRoleAttachmentCreate,
		resourceAlibabacloudStackInstanceRoleAttachmentRead,
		nil,
		resourceAlibabacloudStackInstanceRoleAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackInstanceRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var instanceId string
	var instanceIds []string
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIds = expandStringList(v.(*schema.Set).List())
		for i, k := range instanceIds {
			if i != 0 {
				instanceId = fmt.Sprintf("%s\",\"%s", instanceId, k)
			} else {
				instanceId = k
			}
		}
	}
	request := ecs.CreateAttachInstanceRamRoleRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", instanceId)
	request.RamRoleName = d.Get("role_name").(string)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachInstanceRamRole(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
				return resource.RetryableError(errmsgs.WrapError(errmsgs.Error("Please trying again.")))
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.AttachInstanceRamRoleResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ram_role_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetId(d.Get("role_name").(string) + COLON_SEPARATED + instanceId)
		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ram_role_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func resourceAlibabacloudStackInstanceRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	roleName := parts[0]
	client := meta.(*connectivity.AlibabacloudStackClient)
	ramService := RamService{client}
	object, err := ramService.DescribeRamRoleAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	instRoleSets := object.InstanceRamRoleSets.InstanceRamRoleSet
	var instIds []string
	for _, item := range instRoleSets {
		if item.RamRoleName == roleName {
			instIds = append(instIds, item.InstanceId)
		}
	}
	d.Set("role_name", object.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName)
	d.Set("instance_ids", instIds)
	return nil
}

func resourceAlibabacloudStackInstanceRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ramService := RamService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	roleName := parts[0]
	instanceIds := parts[1]

	request := ecs.CreateDetachInstanceRamRoleRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RamRoleName = roleName
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", instanceIds)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachInstanceRamRole(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.DetachInstanceRamRoleResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return errmsgs.WrapError(ramService.WaitForRamRoleAttachment(d.Id(), Deleted, DefaultTimeout))
}
