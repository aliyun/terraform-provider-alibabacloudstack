package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCrEeNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCrEeNamespaceCreate,
		Read:   resourceAlibabacloudStackCrEeNamespaceRead,
		Update: resourceAlibabacloudStackCrEeNamespaceUpdate,
		Delete: resourceAlibabacloudStackCrEeNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"auto_create": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_visibility": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{RepoTypePublic, RepoTypePrivate}, false),
			},
		},
	}
}

func resourceAlibabacloudStackCrEeNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("name").(string)
	autoCreate := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)

	request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "CreateNamespace", "")
	mergeMaps(request.QueryParams, map[string]string{
		"InstanceId":      instanceId,
		"NamespaceName":   namespace,
		"AutoCreateRepo":  fmt.Sprintf("%t", autoCreate),
		"DefaultRepoType": visibility,
	})

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("create ee namespace failed, %s", response["asapiErrorMessage"].(string))
	}
	d.SetId(crService.GenResourceId(instanceId, namespace))

	return resourceAlibabacloudStackCrEeNamespaceRead(d, meta)
}

func resourceAlibabacloudStackCrEeNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	response, err := crService.DescribeCrEeNamespace(d.Id())

	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("read ee namespace failed, %s", response["asapiErrorMessage"].(string))
	}

	d.Set("instance_id", response["InstanceId"].(string))
	d.Set("name", response["NamespaceName"].(string))
	d.Set("auto_create", response["AutoCreateRepo"].(bool))
	d.Set("default_visibility", response["DefaultRepoType"].(string))

	return nil
}

func resourceAlibabacloudStackCrEeNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("name").(string)
	if d.HasChanges("auto_create", "default_visibility") {
		autoCreate := d.Get("auto_create").(bool)
		visibility := d.Get("default_visibility").(string)

		request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "UpdateNamespace", "")
		mergeMaps(request.QueryParams, map[string]string{
			"InstanceId":      instanceId,
			"NamespaceName":   namespace,
			"AutoCreateRepo":  fmt.Sprintf("%t", autoCreate),
			"DefaultRepoType": visibility,
		})

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		response := make(map[string]interface{})
		addDebug(request.GetActionName(), response, request, request.QueryParams)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if !response["asapiSuccess"].(bool) {
			return fmt.Errorf("update ee namespace failed, %s", response["asapiErrorMessage"].(string))
		}
	}

	return resourceAlibabacloudStackCrEeNamespaceRead(d, meta)
}

func resourceAlibabacloudStackCrEeNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	strRet := crService.ParseResourceId(d.Id())
	instanceId := strRet[0]
	namespaceName := strRet[1]

	request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "DeleteNamespace", "")
	mergeMaps(request.QueryParams, map[string]string{
		"InstanceId":    instanceId,
		"NamespaceName": namespaceName,
	})

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), response, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("delete ee namespace failed, %s", response["asapiErrorMessage"].(string))
	}

	return nil
}
