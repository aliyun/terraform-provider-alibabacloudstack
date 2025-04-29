package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackCRNamespace() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackCRNamespaceCreate, resourceAlibabacloudStackCRNamespaceRead, resourceAlibabacloudStackCRNamespaceUpdate, resourceAlibabacloudStackCRNamespaceDelete)
	return resource
}

func resourceAlibabacloudStackCRNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	resp := crResponse{}
	namespaceName := d.Get("name").(string)
	request := client.NewCommonRequest("PUT", "cr", "2016-06-07", "CreateNamespace", "/namespace")
	body := map[string]interface{}{
		"namespace": map[string]interface{}{
			"namespace":     namespaceName,
			"haApsaraStack": "false",
			"arch":          "x86_64",
		},
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	request.SetContentType(requests.Json)
	request.SetContent(jsonData)
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	log.Printf("response for create %v", bresponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error Unmarshal to JSON: %v", err))
	}
	log.Printf("unmarshalled response for create %v", resp)
	addDebug(request.GetActionName(), bresponse, request)
	create := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)
	if create == false || visibility == "PUBLIC" {
		request := client.NewCommonRequest("POST", "cr", "2016-06-07", "UpdateNamespace", fmt.Sprintf("/namespace/%s", namespaceName))
		body = map[string]interface{}{
			"namespace": map[string]interface{}{
				"AutoCreate":        fmt.Sprintf("%t", create),
				"DefaultVisibility": visibility,
			},
		}
		request.QueryParams["Namespace"] = namespaceName
		jsonData, err := json.Marshal(body)
		if err != nil {
			return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
		}
		request.SetContentType(requests.Json)
		request.SetContent(jsonData)
		uresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if uresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(uresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		err = json.Unmarshal(uresponse.GetHttpContentBytes(), &resp)
		log.Printf("response for update %v", &resp)
		if err != nil {
			return errmsgs.WrapError(fmt.Errorf("Error Unmarshal to JSON: %v", err))
		}
		addDebug(request.GetActionName(), uresponse, request)
	}

	d.SetId(namespaceName)

	return nil
}

func resourceAlibabacloudStackCRNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	create := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)
	if d.HasChanges("auto_create", "default_visibility") {
		request := client.NewCommonRequest("POST", "cr", "2016-06-07", "UpdateNamespace", fmt.Sprintf("/namespace/%s", d.Id()))
		body := map[string]interface{}{
			"namespace": map[string]interface{}{
				"AutoCreate":        fmt.Sprintf("%t", create),
				"DefaultVisibility": visibility,
			},
		}
		request.QueryParams["Namespace"] = d.Id()
		jsonData, err := json.Marshal(body)
		if err != nil {
			return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
		}
		request.SetContentType(requests.Json)
		request.SetContent(jsonData)
		uresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if uresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(uresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), uresponse, request)
	}

	return nil
}

func resourceAlibabacloudStackCRNamespaceRead(d *schema.ResourceData, meta interface{}) error {
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
	request := client.NewCommonRequest("DELETE", "cr", "2016-06-07", "DeleteNamespace", fmt.Sprintf("/namespace/%s", d.Id()))
	request.QueryParams["Namespace"] = d.Id()
	uresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if uresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(uresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(uresponse.GetHttpContentBytes(), &resp)
	log.Printf("response for delete %v", &resp)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error Unmarshal to JSON: %v", err))
	}

	addDebug(request.GetActionName(), uresponse, request)
	return nil
}
