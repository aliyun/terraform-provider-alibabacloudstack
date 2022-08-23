package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
	//"encoding/json"
)

func resourceApsaraStackAscmUserGroupUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmUserGroupUserCreate,
		Read:   resourceApsaraStackAscmUserGroupUserRead,
		Delete: resourceApsaraStackAscmUserGroupUserDelete,
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"login_names": {
				Type: schema.TypeSet,
				//Computed: true,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			//"login_name": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	ForceNew: true,
			//},
		},
	}
}

func resourceApsaraStackAscmUserGroupUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client

	userGroupId := d.Get("user_group_id").(string)
	var loginNamesList []string

	if v, ok := d.GetOk("login_names"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())

		for _, loginName := range loginNames {
			loginNamesList = append(loginNamesList, loginName)
		}
	}

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}

	request.Headers["x-ascm-product-name"] = "ascm"
	request.Headers["x-ascm-product-version"] = "2019-05-10"

	QueryParams := map[string]interface{}{
		"userGroupId":   userGroupId,
		"LoginNameList": loginNamesList,
	}

	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	requeststring, err := json.Marshal(QueryParams)

	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers["Content-Type"] = requests.Json
	request.SetContent(requeststring)
	request.PathPattern = "/roa/ascm/auth/user/addUsersToUserGroup"
	request.ApiName = "AddUsersToUserGroup"
	request.RegionId = client.RegionId
	request.Headers["RegionId"] = client.RegionId
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw AddUsersToUserGroup is : %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_usergroup_user", "AddUsersToUserGroup", raw)
	}

	addDebug("AddUsersToUserGroup", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)

	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_usergroup_user", "AddUsersToUserGroup", ApsaraStackSdkGoERROR)
	}
	addDebug("AddUsersToUserGroup", raw, requestInfo, bresponse.GetHttpContentString())

	d.SetId(userGroupId)

	return resourceApsaraStackAscmUserGroupUserRead(d, meta)
}

func resourceApsaraStackAscmUserGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUsergroupUser(d.Id())
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

	var loginNames []string
	for _, data := range object.Data {
		loginNames = append(loginNames, data.LoginName)
	}

	d.Set("login_names", loginNames)

	return nil
}

func resourceApsaraStackAscmUserGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	var login_names []string
	userGroupId := d.Get("user_group_id").(string)
	if v, ok := d.GetOk("login_names"); ok {
		login_names = expandStringList(v.(*schema.Set).List())
	}

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}

	request.Headers["x-ascm-product-name"] = "ascm"
	request.Headers["x-ascm-product-version"] = "2019-05-10"

	QueryParams := map[string]interface{}{
		"userGroupId":   userGroupId,
		"LoginNameList": login_names,
	}

	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	requeststring, err := json.Marshal(QueryParams)

	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers["Content-Type"] = requests.Json
	request.SetContent(requeststring)
	request.PathPattern = "/roa/ascm/auth/user/RemoveUsersFromUserGroup"
	request.ApiName = "RemoveUsersFromUserGroup"
	request.RegionId = client.RegionId
	request.Headers["RegionId"] = client.RegionId
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw AddUsersToUserGroup is : %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_usergroup_user", "AddUsersToUserGroup", raw)
	}

	return nil
}
