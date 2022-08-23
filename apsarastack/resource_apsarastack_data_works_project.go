package apsarastack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackDataWorksProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDataWorksProjectCreate,
		Read:   resourceApsaraStackDataWorksProjectRead,
		Update: resourceApsaraStackDataWorksProjectUpdate,
		Delete: resourceApsaraStackDataWorksProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PROJECT",
			},
		},
	}
}

func resourceApsaraStackDataWorksProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	action := "CreateProject"
	request := make(map[string]interface{})
	conn, err := client.NewDataworksPrivateClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("project_name"); ok {
		request["ProjectName"] = v.(string)
		request["ProjectIdentifier"] = v.(string)
		request["ProjectDesc"] = v.(string)
	}

	if v, ok := d.GetOk("task_auth_type"); ok {
		request["TaskAuthType"] = v.(string)
	}

	request["RegionId"] = "default"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-17"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_data_works_project", action, ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RequestId"], ":", response["Data"]))

	return resourceApsaraStackDataWorksProjectRead(d, meta)
}
func resourceApsaraStackDataWorksProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksProject(d.Id())
	log.Printf(fmt.Sprint(object))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_data_works_folder dataworksPublicService.DescribeDataWorksProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("project_id", parts[1])

	return nil
}
func resourceApsaraStackDataWorksProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	// 没有对应 API
	return resourceApsaraStackDataWorksProjectRead(d, meta)
}
func resourceApsaraStackDataWorksProjectDelete(d *schema.ResourceData, meta interface{}) error {
	// 没有对应 API
	return nil
}
