package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
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

	user_name := d.Get("user_name").(string)
	action := "CreateOdpsUser"
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	request.SetDomain(client.Config.Endpoints[connectivity.ASAPICode])
	mergeMaps(request.QueryParams, map[string]string{
		"UserName":       user_name,
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
	users, err := maxcomputeService.DescribeMaxcomputeUsers()
	if err != nil || len(users) == 0 {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	var userid float64
	for _, user := range users {
		user_map := user.(map[string]interface{})
		if user_map["userName"].(string) == user_name {
			userid = user_map["id"].(float64)
			break
		}
	}
	id := fmt.Sprintf("%v", userid)
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
	d.Set("account_id", int(object["id"].(float64)))
	d.Set("user_id", object["userId"].(string))
	d.Set("user_pk", object["aasPk"].(string))
	d.Set("user_name", object["userName"].(string))
	d.Set("user_type", object["userType"].(string))
	d.Set("description", object["description"].(string))
	return nil
}

func resourceAlibabacloudStackMaxcomputeUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if !d.IsNewResource() && d.HasChanges("user_name", "description") {
		action := "UpdateOdpsUser"
		user_name := d.Get("user_name").(string)
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
		request.SetDomain(client.Config.Endpoints[connectivity.ASAPICode])
		mergeMaps(request.QueryParams, map[string]string{
			"UserName":       user_name,
			"Id":             d.Id(),
			"UserId":         d.Get("user_id").(string),
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
