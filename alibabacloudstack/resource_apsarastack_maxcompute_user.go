package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackMaxcomputeUser() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"user_pk": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_type": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"organization_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 255),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackMaxcomputeUserCreate, resourceAlibabacloudStackMaxcomputeUserRead, resourceAlibabacloudStackMaxcomputeUserUpdate, resourceAlibabacloudStackMaxcomputeUserDelete)
	return resource
}

func resourceAlibabacloudStackMaxcomputeUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	organization_id := client.Department
	if v, ok := d.GetOk("organization_id"); ok {
		organization_id = fmt.Sprint(v.(int))
	}
	user_name := d.Get("user_name").(string)
	action := "CreateOdpsUser"
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	mergeMaps(request.QueryParams, map[string]string{
		"Region":         client.RegionId,
		"Action":         action,
		"AccessKeyId":    client.AccessKey,
		"UserName":       user_name,
		"OrganizationId": organization_id,
		"Description":    d.Get("description").(string),
	})
	response := make(map[string]interface{})
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(action, bresponse, request, request.QueryParams)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	maxcomputeService := MaxcomputeService{client}
	users, err := maxcomputeService.DescribeMaxcomputeUsers(organization_id)
	if err != nil || len(users) == 0 {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	var userid json.Number
	for _, user := range users {
		user_map := user.(map[string]interface{})
		if user_map["userName"].(string) == user_name {
			userid = user_map["id"].(json.Number)
			break
		}
	}
	id := fmt.Sprintf("%s:%v", organization_id, userid)
	d.SetId(id)
	return
}

func resourceAlibabacloudStackMaxcomputeUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	object, err := maxcomputeService.DescribeMaxcomputeUser(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project_user maxcomputeService.DescribeMaxcomputeUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("account_id", object["id"].(json.Number))
	d.Set("user_id", object["userId"].(string))
	d.Set("user_pk", object["aasPk"].(string))
	d.Set("user_name", object["userName"].(string))
	d.Set("user_type", object["userType"].(string))
	d.Set("organization_id", object["organizationId"].(json.Number))
	d.Set("organization_name", object["organizationName"].(string))
	d.Set("description", object["description"].(string))
	return nil
}

func resourceAlibabacloudStackMaxcomputeUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if !d.IsNewResource() && d.HasChanges("user_name", "description") {
		params := strings.Split(d.Id(), ":")
		action := "UpdateOdpsUser"
		user_name := d.Get("user_name").(string)
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
		mergeMaps(request.QueryParams, map[string]string{
			"Region":         client.RegionId,
			"Action":         action,
			"AccessKeyId":    client.AccessKey,
			"UserName":       user_name,
			"Id":             params[1],
			"UserId":         d.Get("user_id").(string),
			"OrganizationId": params[0],
			"Description":    d.Get("description").(string),
		})

		// response, err = client.DoTeaRequest("POST", "ascm", "2019-05-10", action, "", nil, request, request)
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(action, bresponse, request, request.QueryParams)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func resourceAlibabacloudStackMaxcomputeUserDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
