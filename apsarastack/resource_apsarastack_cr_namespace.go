package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackCRNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackCRNamespaceCreate,
		Read:   resourceApsaraStackCRNamespaceRead,
		Update: resourceApsaraStackCRNamespaceUpdate,
		Delete: resourceApsaraStackCRNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false),
			},
		},
	}
}

func resourceApsaraStackCRNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var uresponse, bresponse *responses.CommonResponse
	resp := crResponse{}
	namespaceName := d.Get("name").(string)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "cr"
	request.Domain = client.Domain
	request.Version = "2016-06-07"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "CreateNamespace"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "cr-ee",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"Action":          "CreateNamespace",
		"Version":         "2016-06-07",
		"X-acs-body": fmt.Sprintf("{\"%s\":{\"%s\":\"%s\"}}",
			"Namespace", "Namespace", namespaceName),
	}
	raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	bresponse, _ = raw.(*responses.CommonResponse)
	log.Printf("response for create %v", bresponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	log.Printf("umarshalled response for create %v", resp)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	create := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)
	if create == false || visibility == "PUBLIC" {
		request.ApiName = "UpdateNamespace"
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "cr",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"RegionId":        client.RegionId,
			"Action":          "UpdateNamespace",
			"Version":         "2016-06-07",
			"X-acs-body": fmt.Sprintf("{\"%s\":{\"%s\":%t,\"%s\":\"%s\"}}",
				"Namespace", "AutoCreate", create, "DefaultVisibility", visibility),
			"Namespace": namespaceName,
		}
		raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
			return crClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		uresponse, _ = raw.(*responses.CommonResponse)
		err = json.Unmarshal(uresponse.GetHttpContentBytes(), &resp)
		log.Printf("response for update %v", &resp)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

	}
	addDebug(request.GetActionName(), raw, request)

	d.SetId(namespaceName)

	return resourceApsaraStackCRNamespaceUpdate(d, meta)
}

func resourceApsaraStackCRNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	//
	create := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)
	if d.HasChange("auto_create") || d.HasChange("default_visibility") {
		request := requests.NewCommonRequest()
		request.Method = "POST"
		request.Product = "cr"
		request.Domain = client.Domain
		request.Version = "2016-06-07"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "UpdateNamespace"
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "cr",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"RegionId":        client.RegionId,
			"Action":          "UpdateNamespace",
			"Version":         "2016-06-07",
			"X-acs-body": fmt.Sprintf("{\"%s\":{\"%s\":%t,\"%s\":\"%s\"}}",
				"Namespace", "AutoCreate", create, "DefaultVisibility", visibility),
			"Namespace": d.Id(),
		}
		raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
			return crClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request)
	}

	return resourceApsaraStackCRNamespaceRead(d, meta)
}

func resourceApsaraStackCRNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	crService := CrService{client}

	object, err := crService.DescribeCrNamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Data.Namespace.Namespace)
	d.Set("auto_create", object.Data.Namespace.AutoCreate)
	d.Set("default_visibility", object.Data.Namespace.DefaultVisibility)

	return nil
}

func resourceApsaraStackCRNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	resp := crResponse{}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "cr"
	request.Domain = client.Domain
	request.Version = "2016-06-07"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DeleteNamespace"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "cr",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"Action":          "DeleteNamespace",
		"Version":         "2016-06-07",
		"Namespace":       d.Id(),
	}
	raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	uresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(uresponse.GetHttpContentBytes(), &resp)
	log.Printf("response for delete %v", &resp)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request)
	return nil
}
