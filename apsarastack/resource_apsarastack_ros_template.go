package apsarastack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApsaraStackRosTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackRosTemplateCreate,
		Read:   resourceApsaraStackRosTemplateRead,
		Update: resourceApsaraStackRosTemplateUpdate,
		Delete: resourceApsaraStackRosTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"template_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackRosTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	action := "CreateTemplate"
	request := make(map[string]interface{})
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("template_body"); ok {
		request["TemplateBody"] = v
	}

	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "ROS"
	request["product"] = "ROS"
	request["OrganizationId"] = client.Department
	//wait := incrementalWait(3*time.Second, 3*time.Second)
	//err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	//if err != nil {
	//	if NeedRetry(err) {
	//		wait()
	//		return resource.RetryableError(err)
	//	}
	//	return resource.NonRetryableError(err)
	//}
	//addDebug(action, response, request)
	//return nil
	//})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ros_template", action, ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return resourceApsaraStackRosTemplateRead(d, meta)
}
func resourceApsaraStackRosTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_ros_template rosService.DescribeRosTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	//d.Set("description", object["Description"])
	d.Set("template_body", object["TemplateBody"])
	//d.Set("template_name", object["TemplateName"])

	//listTagResourcesObject, err := rosService.ListTagResources(d.Id(), "template")
	//if err != nil {
	//	return WrapError(err)
	//}
	//d.Set("tags", tagsToMap(listTagResourcesObject))
	return nil
}
func resourceApsaraStackRosTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	//rosService := RosService{client}
	var response map[string]interface{}
	d.Partial(true)

	//if d.HasChange("tags") {
	//	if err := rosService.SetResourceTags(d, "template"); err != nil {
	//		return WrapError(err)
	//	}
	//	d.SetPartial("tags")
	//}
	update := false
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("template_body") {
		update = true
		request["TemplateBody"] = d.Get("template_body")
	}
	if !d.IsNewResource() && d.HasChange("template_name") {
		update = true
		request["TemplateName"] = d.Get("template_name")
	}
	if update {
		if _, ok := d.GetOk("template_url"); ok {
			request["TemplateURL"] = d.Get("template_url")
		}
		action := "UpdateTemplate"
		request["RegionId"] = client.RegionId
		request["Product"] = "ROS"
		request["product"] = "ROS"
		request["OrganizationId"] = client.Department
		conn, err := client.NewRosClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		//d.SetPartial("description")
		//d.SetPartial("template_body")
		//d.SetPartial("template_name")
	}
	d.Partial(false)
	return resourceApsaraStackRosTemplateRead(d, meta)
}
func resourceApsaraStackRosTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	action := "DeleteTemplate"
	var response map[string]interface{}
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "ROS"
	request["product"] = "ROS"
	request["OrganizationId"] = client.Department
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ChangeSetNotFound", "StackNotFound", "TemplateNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	return nil
}
