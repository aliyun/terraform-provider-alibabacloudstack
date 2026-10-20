package alibabacloudstack

import (
	"fmt"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmRole() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_visibility": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"privileges": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmRoleCreate, resourceAlibabacloudStackAscmRoleRead, resourceAlibabacloudStackAscmRoleUpdate, resourceAlibabacloudStackAscmRoleDelete)
	return resource
}

func resourceAlibabacloudStackAscmRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	name := d.Get("role_name").(string)
	description := d.Get("description").(string)
	roleRange := d.Get("role_range").(string)
	var privs []string
	if v, ok := d.GetOk("privileges"); ok {
		privs = expandStringList(v.(*schema.Set).List())
	}
	check, err := ascmService.DescribeAscmCustomRole(name)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_custom_role", "role alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	organizationVisibility := d.Get("organization_visibility").(string)
	if len(check.Data) == 0 {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateRole", "/ascm/auth/role/createRole")
		requestBody := map[string]interface{}{
			"roleName":               name,
			"description":            description,
			"roleRange":              roleRange,
			"organizationVisibility": organizationVisibility,
			"privileges":             privs,
		}
		bodyJSON, _ := json.Marshal(requestBody)
		request.Content = bodyJSON
		request.SetContentType("application/json")
		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf("response of raw CreateRole: %s", bresponse.GetHttpContentString())
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_custom_role", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmCustomRole(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	d.SetId(name + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))
	return nil
}

func resourceAlibabacloudStackAscmRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlibabacloudStackAscmRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmCustomRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	log.Printf("Privileges for did[0]:%v", object.Data[0].Privileges)
	d.Set("role_name", did[0])
	d.Set("organization_visibility", object.Data[0].OrganizationVisibility)
	d.Set("role_id", object.Data[0].ID)
	d.Set("description", object.Data[0].Description)
	d.Set("role_range", object.Data[0].RoleRange)
	return nil
}

func resourceAlibabacloudStackAscmRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	check, err := ascmService.DescribeAscmCustomRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsRoleExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsCustomRoleExist", check, nil, map[string]string{"roleName": did[0]})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRole", "/ascm/auth/role/removeRole")
		request.QueryParams["roleName"] = did[0]
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_custom_role", "RemoveRole", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		_, err = ascmService.DescribeAscmCustomRole(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}
