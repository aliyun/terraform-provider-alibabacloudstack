package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackCRNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCRNamespaceCreate,
		Read:   resourceAlibabacloudStackCRNamespaceRead,
		Update: resourceAlibabacloudStackCRNamespaceUpdate,
		Delete: resourceAlibabacloudStackCRNamespaceDelete,
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

func resourceAlibabacloudStackCRNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	resp := crResponse{}
	namespaceName := d.Get("name").(string)
	request := client.NewCommonRequest("POST", "cr", "2016-06-07", "CreateNamespace", "")
	request.SetContentType("application/json")
	request.SetContent([]byte("{}")) // 必须指定，否则SDK会将类型修改为www-form，最终导致cr有一定的随机概率失败
	request.QueryParams["NamespaceName"] = namespaceName
	request.QueryParams["Arch"] = "x86_64"
	request.QueryParams["HaApsaraStack"] = "false"
	request.QueryParams["SignatureVersion"] = "2.1"
	request.QueryParams["Language"] = "zh"
	request.QueryParams["x-acs-body"] = fmt.Sprintf("{\"%s\":{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%d\",\"%s\":\"%d\"}}",
		"namespace", "NamespaceName", namespaceName, "namespace", namespaceName, "Language", "zh", "haApsaraStack", "false", "arch", "x86_64", "RegionId", "cn-wulan-env48-d01", "Department", 37, "ResourceGroup", 124)
	raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	log.Printf("response for create %v", bresponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	log.Printf("unmarshalled response for create %v", resp)
	if err != nil {
	}
	create := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)
	if create == false || visibility == "PUBLIC" {
		request.ApiName = "UpdateNamespace"
		request.Headers["x-acs-instanceId"] = namespaceName
		request.Headers["x-acs-content-type"] = "application/json;charset=UTF-8"
		request.Headers["Content-type"] = "application/json;charset=UTF-8"
		request.Headers["x-acs-body"] = fmt.Sprintf("{\"%s\":{\"%s\":%t,\"%s\":\"%s\"}}",
			"Namespace", "AutoCreate", create, "DefaultVisibility", visibility)
		request.QueryParams["Namespace"] = namespaceName
		request.QueryParams["AutoCreate"] = fmt.Sprintf("%t", create)
		request.QueryParams["DefaultVisibility"] = visibility
		raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.ProcessCommonRequest(request)
		})
		uresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(uresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		err = json.Unmarshal(uresponse.GetHttpContentBytes(), &resp)
		log.Printf("response for update %v", &resp)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(uresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	addDebug(request.GetActionName(), raw, request)

	d.SetId(namespaceName)

	return resourceAlibabacloudStackCRNamespaceUpdate(d, meta)
}

func resourceAlibabacloudStackCRNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	create := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)
	if d.HasChanges("auto_create", "default_visibility") {
		request := client.NewCommonRequest("POST", "cr", "2016-06-07", "UpdateNamespace", "")
		request.SetContentType("application/json")
		request.SetContent([]byte("{}")) // 必须指定，否则SDK会将类型修改为www-form，最终导致cr有一定的随机概率失败
		request.Headers["x-acs-instanceId"] = d.Id()
		request.Headers["x-acs-content-type"] = "application/json;charset=UTF-8"
		request.Headers["Content-type"] = "application/json;charset=UTF-8"
		request.Headers["x-acs-body"] = fmt.Sprintf("{\"%s\":{\"%s\":%t,\"%s\":\"%s\"}}",
			"Namespace", "AutoCreate", create, "DefaultVisibility", visibility)
		request.QueryParams["Namespace"] = d.Id()
		request.QueryParams["AutoCreate"] = fmt.Sprintf("%t", create)
		request.QueryParams["DefaultVisibility"] = visibility
		raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.ProcessCommonRequest(request)
		})
		uresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(uresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request)
	}

	return resourceAlibabacloudStackCRNamespaceRead(d, meta)
}

func resourceAlibabacloudStackCRNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := CrService{client}

	object, err := crService.DescribeCrNamespace(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object.Data.Namespace.Namespace)
	d.Set("auto_create", object.Data.Namespace.AutoCreate)
	d.Set("default_visibility", object.Data.Namespace.DefaultVisibility)

	return nil
}

func resourceAlibabacloudStackCRNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	resp := crResponse{}
	request := client.NewCommonRequest("POST", "cr", "2016-06-07", "DeleteNamespace", "")
	request.Headers["x-acs-instanceId"] = d.Id()
	request.Headers["x-acs-content-type"] = "application/json;charset=UTF-8"
	request.Headers["Content-type"] = "application/json;charset=UTF-8"
	request.QueryParams["Namespace"] = d.Id()
	raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	uresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(uresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(uresponse.GetHttpContentBytes(), &resp)
	log.Printf("response for delete %v", &resp)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(uresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request)
	return nil
}
