package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksUser() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"role_project_owner", "role_project_admin", "role_project_dev", "role_project_pe", "role_project_deploy", "role_project_guest", "role_project_security"}, false),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksUserCreate, resourceAlibabacloudStackDataWorksUserRead, resourceAlibabacloudStackDataWorksUserUpdate, resourceAlibabacloudStackDataWorksUserDelete)
	return resource
}

func resourceAlibabacloudStackDataWorksUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateProjectMember"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("project_id"); ok {
		request["ProjectId"] = v.(string)
	}

	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v.(string)
	}

	if v, ok := d.GetOk("role_code"); ok {
		request["RoleCode"] = v.(string)
	}

	request["ClientToken"] = fmt.Sprint(uuid.NewRandom())

	response, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_data_works_folder", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RequestId"], ":", request["ProjectId"], ":", request["UserId"]))

	return nil
}

func resourceAlibabacloudStackDataWorksUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksUser(d.Id())
	log.Printf(fmt.Sprint(object))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("user_id", parts[2])
	d.Set("project_id", parts[1])

	return nil
}

func resourceAlibabacloudStackDataWorksUserUpdate(d *schema.ResourceData, meta interface{}) error {
	noUpdateAllowedFields := []string{"project_id", "user_id", "role_code"}
	return noUpdatesAllowedCheck(d, noUpdateAllowedFields)
}

func resourceAlibabacloudStackDataWorksUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteProjectMember"
	request := map[string]interface{}{
		"ProjectId": parts[1],
		"UserId":    parts[2],
	}

	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
