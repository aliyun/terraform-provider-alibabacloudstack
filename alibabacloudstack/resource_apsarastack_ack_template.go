package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAckTemplate() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_with_hist_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_hash_code_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ali_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceAlibabacloudStackAckTemplateCreate,
		Read:   resourceAlibabacloudStackAckTemplateRead,
		Update: resourceAlibabacloudStackAckTemplateUpdate,
		Delete: resourceAlibabacloudStackAckTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return resource
}

func resourceAlibabacloudStackAckTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare the request body
	reqBody := map[string]interface{}{
		"template":      d.Get("template").(string),
		"name":          d.Get("name").(string),
		"description":   d.Get("description").(string),
		"template_type": d.Get("template_type").(string),
	}

	// Call the API to create the template
	resp, err := client.DoTeaRequest("POST", "cs", "2015-12-15", "CreateTemplate", "/templates", nil, nil, reqBody)
	if err != nil {
		return err
	}

	// Extract template ID from response
	templateId, ok := resp["template_id"].(string)
	if !ok || templateId == "" {
		return fmt.Errorf("failed to retrieve template_id from CreateTemplate response")
	}

	// Set the resource ID temporarily
	d.SetId(templateId)

	// No status check is required according to the documentation, so we directly set the final ID
	return resourceAlibabacloudStackAckTemplateRead(d, meta)
}

func resourceAlibabacloudStackAckTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ackTemplateService := CsService{client}

	object, err := ackTemplateService.DescribeAckTemplate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ack_template ackTemplateService.DescribeAckTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("template", object["template"])
	d.Set("name", object["name"])
	d.Set("description", object["description"])
	d.Set("template_type", object["template_type"])
	d.Set("template_id", object["id"])
	d.Set("template_with_hist_id", object["template_with_hist_id"])
	d.Set("template_hash_code_version", object["template_hash_code_version"])
	d.Set("created", object["created"])
	d.Set("acl", object["acl"])
	d.Set("version", object["version"])
	d.Set("tags", object["tags"])
	d.Set("ali_uid", object["ali_uid"])
	d.Set("updated", object["updated"])

	return nil
}

func resourceAlibabacloudStackAckTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ackTemplateService := CsService{client}
	if d.IsNewResource() {
		return nil
	}
	object, err := ackTemplateService.DescribeAckTemplate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ack_template ackTemplateService.DescribeAckTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	if d.HasChanges("template", "description") {
		templateId := d.Id()
		pathPattern := fmt.Sprintf("/templates/%s", templateId)

		body := make(map[string]interface{})
		body["template"] = object["template"]
		body["template_with_hist_id"] = object["template_with_hist_id"]
		body["template_hash_code_version"] = object["template_hash_code_version"]
		body["acl"] = object["acl"]
		body["name"] = object["name"]
		body["template_type"] = object["template_type"]
		body["id"] = templateId
		body["description"] = object["description"]
		if d.HasChange("template") {
			body["template"] = d.Get("template")
		}
		if d.HasChange("description") {
			body["description"] = d.Get("description")
		}

		_, err := client.DoTeaRequest("PUT", "cs", "2015-12-15", "UpdateTemplate", pathPattern, nil, nil, body)
		if err != nil {
			return fmt.Errorf("updating Ack Template failed: %v", err)
		}
	}

	return resourceAlibabacloudStackAckTemplateRead(d, meta)
}

func resourceAlibabacloudStackAckTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	templateId := d.Id()

	// Define the path pattern and other required parameters for the delete API
	pathPattern := fmt.Sprintf("/templates/%s", templateId)
	method := "DELETE"
	popCode := "cs"
	version := "2015-12-15"
	apiName := "DeleteTemplate"

	// Call the delete API
	_, err := client.DoTeaRequest(method, popCode, version, apiName, pathPattern, nil, nil, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, templateId, apiName, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Since there is no response data to consume, we directly return
	return nil
}
