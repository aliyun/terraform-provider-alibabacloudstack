package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackOosTemplate() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_delete_executions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_trigger": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"template_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOosTemplateCreate, resourceAlibabacloudStackOosTemplateRead, resourceAlibabacloudStackOosTemplateUpdate, resourceAlibabacloudStackOosTemplateDelete)
	return resource
}

func resourceAlibabacloudStackOosTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateTemplate"
	request := make(map[string]interface{})
	request["Content"] = d.Get("content")
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request["Tags"] = respJson
	}
	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("version_name"); ok {
		request["VersionName"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	response, err := client.DoTeaRequest("POST", "Oos", "2019-06-01", action, "", nil, nil, request)
	if err != nil {
		errmsg := ""
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oos_template", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	responseTemplate := response["Template"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseTemplate["TemplateName"]))
	return nil
}

func resourceAlibabacloudStackOosTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosTemplate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_oos_template oosService.DescribeOosTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("template_name", d.Id())
	d.Set("created_by", object["CreatedBy"])
	d.Set("created_date", object["CreatedDate"])
	d.Set("description", object["Description"])
	d.Set("has_trigger", object["HasTrigger"])
	d.Set("share_type", object["ShareType"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v))
	}
	d.Set("template_format", object["TemplateFormat"])
	d.Set("template_id", object["TemplateId"])
	d.Set("template_type", object["TemplateType"])
	d.Set("template_version", object["TemplateVersion"])
	d.Set("updated_by", object["UpdatedBy"])
	d.Set("updated_date", object["UpdatedDate"])
	return nil
}

func resourceAlibabacloudStackOosTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := map[string]interface{}{
		"TemplateName": d.Id(),
	}
	if d.HasChange("content") {
		update = true
	}
	request["Content"] = d.Get("content")
	if d.HasChange("tags") {
		update = true
		respJson, err := convertMaptoJsonString(d.Get("tags").(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oos_template", "UpdateTemplate", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		request["Tags"] = respJson
	}
	if d.HasChange("version_name") {
		update = true
		request["VersionName"] = d.Get("version_name")
	}
	if update {
		action := "UpdateTemplate"
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
		_, err := client.DoTeaRequest("POST", "Oos", "2019-06-01", action, "", nil, nil, request)
		if err != nil {
			errmsg := ""
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	return nil
}

func resourceAlibabacloudStackOosTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteTemplate"
	request := map[string]interface{}{
		"TemplateName": d.Id(),
	}

	if v, ok := d.GetOkExists("auto_delete_executions"); ok {
		request["AutoDeleteExecutions"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	_, err := client.DoTeaRequest("POST", "Oos", "2019-06-01", action, "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExists.Template"}) {
			return nil
		}
		return err
	}
	return nil
}