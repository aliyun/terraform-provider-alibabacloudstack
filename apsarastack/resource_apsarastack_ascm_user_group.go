package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceApsaraStackAscmUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmUserGroupCreate,
		Read:   resourceApsaraStackAscmUserGroupRead,
		Update: resourceApsaraStackAscmUserGroupUpdate,
		Delete: resourceApsaraStackAscmUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//"role_ids": {
			//	Type:     schema.TypeList,
			//	Computed: true,
			//	Elem:     &schema.Schema{Type: schema.TypeInt},
			//},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceApsaraStackAscmUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client

	groupName := d.Get("group_name").(string)

	organizationId := d.Get("organization_id").(string)
	if organizationId == "" {
		organizationId = client.Department
	}

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ascm",
		"Action":          "CreateUserGroup",
		"Version":         "2019-05-10",
		"ProductName":     "ascm",
		"groupName":       groupName,
		"OrganizationId":  organizationId,
		//"roleIdList":      ,
	}

	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "CreateUserGroup"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw CreateUserGroup is : %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_group", "CreateUserGroup", raw)
	}

	addDebug("CreateUserGroup", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)

	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_group", "CreateUserGroup", ApsaraStackSdkGoERROR)
	}
	addDebug("CreateUserGroup", raw, requestInfo, bresponse.GetHttpContentString())

	d.SetId(groupName)

	return resourceApsaraStackAscmUserGroupUpdate(d, meta)
}

func resourceApsaraStackAscmUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceApsaraStackAscmUserGroupRead(d, meta)
}

func resourceApsaraStackAscmUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if len(object.Data) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("user_group_id", strconv.Itoa(object.Data[0].Id))
	d.Set("group_name", object.Data[0].GroupName)
	d.Set("organization_id", strconv.Itoa(object.Data[0].OrganizationId))

	var roleIds []string
	for _, role := range object.Data[0].Roles {
		roleIds = append(roleIds, strconv.Itoa(role.Id))
	}
	d.Set("role_ids", roleIds)

	return nil
}

func resourceApsaraStackAscmUserGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	check, err := ascmService.DescribeAscmUserGroup(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsUserGroupExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsUserGroupExist", check, requestInfo, map[string]string{"groupName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "DeleteUserGroup",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"userGroupId":     strconv.Itoa(check.Data[0].Id),
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "DeleteUserGroup"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = ascmService.DescribeAscmUserGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteUserGroup", raw, request)
		return nil
	})
	return nil
}
