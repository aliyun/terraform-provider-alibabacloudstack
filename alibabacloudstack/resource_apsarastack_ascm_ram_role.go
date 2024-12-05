package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceAlibabacloudStackAscmRamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmRamRoleCreate,
		Read:   resourceAlibabacloudStackAscmRamRoleRead,
		Update: resourceAlibabacloudStackAscmRamRoleUpdate,
		Delete: resourceAlibabacloudStackAscmRamRoleDelete,
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
}

func resourceAlibabacloudStackAscmRamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	ascmService := AscmService{client}
	name := d.Get("role_name").(string)
	description := d.Get("description").(string)
	rolerange := d.Get("role_range").(string)
	organizationvisibility := d.Get("organization_visibility").(string)
	check, err := ascmService.DescribeAscmRamRole(name)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ram_role", "role alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(check.Data) == 0 {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateRole", "")
		mergeMaps(request.QueryParams, map[string]string{
			"roleName":            name,
			"description":         description,
			"roleRange":           rolerange,
			"roleType":            "ROLETYPE_RAM",
			"organizationVisibility": organizationvisibility,
		})

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" rsponse of CreateRole : %s", raw)

		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*responses.CommonResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_role", "CreateRole", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("CreateRole", raw, requestInfo, request)

		bresponse, ok := raw.(*responses.CommonResponse)
		if !ok {
			return fmt.Errorf("failed to cast response to CommonResponse")
		}
		if bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*responses.CommonResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_role", "CreateRole", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("CreateRole", raw, requestInfo, bresponse.GetHttpContentString())
	}
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmRamRole(name)
		if err != nil {
			if errmsgs.IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DescribeAscmRamRole", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.SetId(name + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))
	return resourceAlibabacloudStackAscmRamRoleUpdate(d, meta)
}

func resourceAlibabacloudStackAscmRamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudStackAscmRamRoleRead(d, meta)
}

func resourceAlibabacloudStackAscmRamRoleRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
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
	return nil
}

func resourceAlibabacloudStackAscmRamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmRamRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsRamRoleExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsRamRoleExist", check, requestInfo, map[string]string{"roleName": did[0]})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRole", "")
		request.QueryParams["roleName"] = did[0]

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*responses.CommonResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_role", "RemoveRole", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		_, err = ascmService.DescribeAscmRamRole(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "RemoveRole", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
