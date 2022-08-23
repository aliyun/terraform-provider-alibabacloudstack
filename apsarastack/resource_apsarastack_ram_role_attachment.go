package apsarastack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackRamRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackInstanceRoleAttachmentCreate,
		Read:   resourceApsaraStackInstanceRoleAttachmentRead,
		Delete: resourceApsaraStackInstanceRoleAttachmentDelete,

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
}

func resourceApsaraStackInstanceRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
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
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"Product":         "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup,
	}
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", instanceId)
	request.RamRoleName = d.Get("role_name").(string)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachInstanceRamRole(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
				return resource.RetryableError(WrapError(Error("Please trying again.")))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, "ram_role_attachment", request.GetActionName(), ApsaraStackSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetId(d.Get("role_name").(string) + COLON_SEPARATED + instanceId)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_role_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return resourceApsaraStackInstanceRoleAttachmentRead(d, meta)
}

func resourceApsaraStackInstanceRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	roleName := parts[0]
	client := meta.(*connectivity.ApsaraStackClient)
	ramService := RamService{client}
	object, err := ramService.DescribeRamRoleAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceApsaraStackInstanceRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ramService := RamService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	roleName := parts[0]
	instanceIds := parts[1]

	request := ecs.CreateDetachInstanceRamRoleRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RamRoleName = roleName
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", instanceIds)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachInstanceRamRole(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(ramService.WaitForRamRoleAttachment(d.Id(), Deleted, DefaultTimeout))
}
