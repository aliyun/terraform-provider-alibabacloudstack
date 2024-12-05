package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCrEENamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCrEENamespaceCreate,
		Read:   resourceAlibabacloudStackCrEENamespaceRead,
		Update: resourceAlibabacloudStackCrEENamespaceUpdate,
		Delete: resourceAlibabacloudStackCrEENamespaceDelete,
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

func resourceAlibabacloudStackCrEENamespaceCreate(d *schema.ResourceData, meta interface{}) error {
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

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("create ee namespace failed, %s", response["asapiErrorMessage"].(string))
	}
	d.SetId(crService.GenResourceId(instanceId, namespace))

	return resourceAlibabacloudStackCrEENamespaceRead(d, meta)
}

func resourceAlibabacloudStackCrEENamespaceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	strRet := crService.ParseResourceId(d.Id())
	instanceId := strRet[0]
	namespaceName := strRet[1]

	request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "ListNamespace", "")
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["NamespaceName"] = namespaceName

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("read ee namespace failed, %s", response["asapiErrorMessage"].(string))
	}
	namespaceList := response["Namespaces"].([]interface{})
	if len(namespaceList) == 0 {
		return errmsgs.WrapError(fmt.Errorf("namespace %s not found", namespaceName))
	}
	item := namespaceList[0].(map[string]interface{})
	d.Set("instance_id", item["InstanceId"].(string))
	d.Set("name", item["NamespaceName"].(string))
	d.Set("auto_create", item["AutoCreateRepo"].(bool))
	d.Set("default_visibility", item["DefaultRepoType"].(string))

	return nil
}

func resourceAlibabacloudStackCrEENamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
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

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		response := make(map[string]interface{})
		addDebug(request.GetActionName(), raw, request, request.QueryParams)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if !response["asapiSuccess"].(bool) {
			return fmt.Errorf("update ee namespace failed, %s", response["asapiErrorMessage"].(string))
		}
	}

	return resourceAlibabacloudStackCrEENamespaceRead(d, meta)
}

func resourceAlibabacloudStackCrEENamespaceDelete(d *schema.ResourceData, meta interface{}) error {
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

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("delete ee namespace failed, %s", response["asapiErrorMessage"].(string))
	}

	return nil
}
