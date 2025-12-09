package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmRamRole() *schema.Resource {
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
			"role_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"role_range": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"roleRange.orgAndSubOrgs", "roleRange.allOrganizations", "roleRange.userGroup"}, false),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmRamRoleCreate, resourceAlibabacloudStackAscmRamRoleRead, resourceAlibabacloudStackAscmRamRoleUpdate, resourceAlibabacloudStackAscmRamRoleDelete)
	return resource
}

func resourceAlibabacloudStackAscmRamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	name := d.Get("role_name").(string)
	description := d.Get("description").(string)
	rolerange := d.Get("role_range").(string)
	organizationvisibility := d.Get("organization_visibility").(string)
	check, err := ascmService.DescribeAscmRamRole(name)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ram_role", "check role failed", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(check.Data) > 0 {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ram_role", "role alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateRole", "/ascm/auth/role/createRole")
	mergeMaps(request.QueryParams, map[string]string{
		"roleName":               name,
		"description":            description,
		"roleRange":              rolerange,
		"roleType":               "ROLETYPE_RAM",
		"organizationVisibility": organizationvisibility,
	})

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_role", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	resp := CreateAscmRolesResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, "Process Common Request Failed")
	}
	d.SetId(name + COLON_SEPARATED + fmt.Sprint(resp.Data.ID))
	return nil
}

func resourceAlibabacloudStackAscmRamRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmRamRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if strings.Contains(object.Data[0].OrganizationVisibility, "organizationVisibility.") {
		object.Data[0].OrganizationVisibility = strings.TrimPrefix(object.Data[0].OrganizationVisibility, "organizationVisibility.")
	}
	d.Set("role_name", did[0])
	d.Set("organization_visibility", object.Data[0].OrganizationVisibility)
	d.Set("role_id", object.Data[0].ID)
	d.Set("description", object.Data[0].Description)
	d.Set("role_range", object.Data[0].RoleRange)
	return nil
}

func resourceAlibabacloudStackAscmRamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	noUpdateAllowedFields := []string{"role_name", "description", "organization_visibility", "role_range"}

	return noUpdatesAllowedCheck(d, noUpdateAllowedFields)
}

func resourceAlibabacloudStackAscmRamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	_, err := ascmService.DescribeAscmRamRole(d.Id())
	if errmsgs.NotFoundError(err) {
		return nil
	}
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsRamRoleExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRole", "/ascm/auth/role/removeRole")
		request.QueryParams["roleName"] = did[0]

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_role", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "RemoveRole", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
