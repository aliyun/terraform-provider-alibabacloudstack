package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksProject() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksProjectCreate,
		resourceAlibabacloudStackDataWorksProjectRead, nil, resourceAlibabacloudStackDataWorksProjectDelete)
	return resource
}

func resourceAlibabacloudStackDataWorksProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateProject"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("project_name"); ok {
		request["ProjectName"] = v.(string)
		request["ProjectIdentifier"] = v.(string)
		request["ProjectDesc"] = v.(string)
	}

	if v, ok := d.GetOk("task_auth_type"); ok {
		request["TaskAuthType"] = v.(string)
	}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	response, err := client.DoTeaRequest("POST", "dataworks-public", "2019-01-17", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["RequestId"], ":", response["Data"]))

	return nil
}

func resourceAlibabacloudStackDataWorksProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksProject(d.Id())
	log.Printf(fmt.Sprint(object))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		errmsg := ""
		if object != nil {
			errmsg = errmsgs.GetAsapiErrorMessage(object)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_data_works_project", "DescribeDataWorksProject", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("project_id", parts[1])

	return nil
}

func resourceAlibabacloudStackDataWorksProjectDelete(d *schema.ResourceData, meta interface{}) error {
	// 没有对应 API
	return nil
}