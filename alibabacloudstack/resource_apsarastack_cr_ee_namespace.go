package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
	// request := cr_ee.CreateCreateNamespaceRequest()
	request := requests.NewCommonRequest()
	request.RegionId = crService.client.RegionId
	request.Product = "cr-ee"
	request.Method = "POST"
	request.Domain = client.Domain
	request.Version = "2018-12-01"
	request.ApiName = "CreateNamespace"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		"Product":         "cr-ee",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "CreateNamespace",
		"Version":         "2018-12-01",
		"InstanceId":      instanceId,
		"NamespaceName":   namespace,
		"AutoCreateRepo":  fmt.Sprintf("%t", autoCreate),
		"DefaultRepoType": visibility,
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	// raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	// 	return ecsClient.DescribeInstanceAutoRenewAttribute(request)
	// })
	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(raw.(*responses.CommonResponse).GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
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
	request := requests.NewCommonRequest()
	request.RegionId = client.RegionId
	request.Product = "cr-ee"
	request.Method = "POST"
	request.Domain = client.Domain
	request.Version = "2018-12-01"
	request.ApiName = "ListNamespace"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		"Product":         "cr-ee",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "ListNamespace",
		"Version":         "2018-12-01",
		"InstanceId":      instanceId,
		"NamespaceName":   namespaceName,
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(raw.(*responses.CommonResponse).GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("read ee namespace failed, %s", response["asapiErrorMessage"].(string))
	}
	namespaceList := response["Namespaces"].([]interface{})
	if len(namespaceList) == 0 {
		return WrapError(fmt.Errorf("namespace %s not found", namespaceName))
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
		request := requests.NewCommonRequest()
		request.RegionId = client.RegionId
		request.Product = "cr-ee"
		request.Method = "POST"
		request.Domain = client.Domain
		request.Version = "2018-12-01"
		request.ApiName = "UpdateNamespace"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			
			"Product":         "cr-ee",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Action":          "UpdateNamespace",
			"Version":         "2018-12-01",
			"InstanceId":      instanceId,
			"NamespaceName":   namespace,
			"AutoCreateRepo":  fmt.Sprintf("%t", autoCreate),
			"DefaultRepoType": visibility,
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		response := make(map[string]interface{})
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		err = json.Unmarshal(raw.(*responses.CommonResponse).GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
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
	request := requests.NewCommonRequest()
	request.RegionId = client.RegionId
	request.Product = "cr-ee"
	request.Method = "POST"
	request.Domain = client.Domain
	request.Version = "2018-12-01"
	request.ApiName = "DeleteNamespace"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		"Product":         "cr-ee",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "DeleteNamespace",
		"Version":         "2018-12-01",
		"InstanceId":      instanceId,
		"NamespaceName":   namespaceName,
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(raw.(*responses.CommonResponse).GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("delete ee namespace failed, %s", response["asapiErrorMessage"].(string))
	}

	return nil
}
