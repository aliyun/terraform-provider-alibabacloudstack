package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
	"time"
)

func resourceApsaraStackAscmUserRoleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmUserRoleBindingCreate,
		Read:   resourceApsaraStackAscmUserRoleBindingRead,
		Update: resourceApsaraStackAscmUserRoleBindingUpdate,
		Delete: resourceApsaraStackAscmUserRoleBindingDelete,
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceApsaraStackAscmUserRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	lname := d.Get("login_name").(string)
	flag := false
	var roleids []int
	if v, ok := d.GetOk("role_ids"); ok {
		roleids = expandIntList(v.(*schema.Set).List())
	}
	log.Printf("roleids is %v", roleids)
	flag = true
	if flag {
		for i := range roleids {
			request := requests.NewCommonRequest()
			if client.Config.Insecure {
				request.SetHTTPSInsecure(client.Config.Insecure)
			}
			request.QueryParams = map[string]string{
				"RegionId":        client.RegionId,
				"AccessKeySecret": client.SecretKey,
				"Product":         "Ascm",
				"Action":          "AddRoleToUser",
				"Version":         "2019-05-10",
				"ProductName":     "ascm",
				"LoginName":       lname,
				"RoleId":          fmt.Sprint(roleids[i]),
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
			request.ApiName = "AddRoleToUser"
			request.RegionId = client.RegionId
			request.Headers = map[string]string{"RegionId": client.RegionId}
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			log.Printf("response of raw AddRoleToUser Role(%d) is : %s", roleids[i], raw)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_role_binding", "AddRoleToUser", raw)
			}

			addDebug("AddRoleToUser", raw, requestInfo, request)

			bresponse, _ := raw.(*responses.CommonResponse)
			if bresponse.GetHttpStatus() != 200 {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_user_role_binding", "AddRoleToUser", ApsaraStackSdkGoERROR)
			}
			addDebug("AddRoleToUser", raw, requestInfo, bresponse.GetHttpContentString())
			log.Printf("response of queryparams AddRoleToUser is : %s", request.QueryParams)
		}

	}

	d.SetId(lname)

	return resourceApsaraStackAscmUserRoleBindingRead(d, meta)
}

func resourceApsaraStackAscmUserRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserRoleBinding(d.Id())
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
	d.Set("login_name", object.Data[0].LoginName)

	return nil
}

func resourceApsaraStackAscmUserRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmUserRoleBindingCreate(d, meta)

}

func resourceApsaraStackAscmUserRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	var roleid int
	flag := false
	var roleids []int

	if v, ok := d.GetOk("role_ids"); ok {
		roleids = expandIntList(v.(*schema.Set).List())
		for i := range roleids {
			if len(roleids) > 1 {
				roleid = roleids[i]
				flag = true
			} else {
				roleid = roleids[0]
				flag = true
			}
		}
	}
	log.Printf("roleid is %v", roleid)
	log.Printf("roleids is %v", roleids)
	check, err := ascmService.DescribeAscmUserRoleBinding(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBindingExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"loginName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		if flag {
			request := requests.NewCommonRequest()
			if client.Config.Insecure {
				request.SetHTTPSInsecure(client.Config.Insecure)
			}
			request.QueryParams = map[string]string{
				"RegionId":        client.RegionId,
				"AccessKeySecret": client.SecretKey,
				"Product":         "ascm",
				"Action":          "RemoveRoleFromUser",
				"Version":         "2019-05-10",
				"ProductName":     "ascm",
				"LoginName":       d.Id(),
				"RoleId":          fmt.Sprint(roleid),
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
			request.ApiName = "RemoveRoleFromUser"
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.RegionId = client.RegionId

			_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
				return csClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return resource.RetryableError(err)
			}
			check, err = ascmService.DescribeAscmUserRoleBinding(d.Id())

			if err != nil {
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "RemoveRoleFromUser", ApsaraStackSdkGoERROR)
	}
	return nil
}
