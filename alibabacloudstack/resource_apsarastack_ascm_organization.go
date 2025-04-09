package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmOrganization() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1",
			},
			"person_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmOrganizationCreate, resourceAlibabacloudStackAscmOrganizationRead, resourceAlibabacloudStackAscmOrganizationUpdate, resourceAlibabacloudStackAscmOrganizationDelete)
	return resource
}

func resourceAlibabacloudStackAscmOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	var requestInfo *ecs.Client
	name := d.Get("name").(string)
	check, err := ascmService.DescribeAscmOrganization(name)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_organization", "ORG alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	parentid := d.Get("parent_id").(string)

	if len(check.Data) == 0 {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateOrganization", "/ascm/auth/organization/add")
		request.QueryParams["parentId"] = parentid
		request.QueryParams["name"] = name
		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf("response of raw CreateOrganization is : %s", bresponse)

		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_organization", "CreateOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("CreateOrganization", bresponse, requestInfo, request)

		if bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_organization", "CreateOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		// TODO: 父组织未找到这里不会报错，因为HttpStatus依旧为200
		addDebug("CreateOrganization", bresponse, requestInfo, bresponse.GetHttpContentString())
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmOrganization(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})

	d.SetId(fmt.Sprint(check.Data[0].ID))

	return nil
}

func resourceAlibabacloudStackAscmOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	name := d.Get("name").(string)
	attributeUpdate := false
	check, err := ascmService.DescribeAscmOrganization(name)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsOrganizationExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data[0].Name = name
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data[0].Name = name
	}
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateOrganization", "/ascm/auth/organization/update")
	request.QueryParams["name"] = name
	request.QueryParams["id"] = d.Id()

	if attributeUpdate {
		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf(" response of raw UpdateOrganization : %s", bresponse)

		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_instance", "ConsoleInstanceCreate", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), bresponse, request)
	}

	d.SetId(fmt.Sprint(check.Data[0].ID))

	return nil
}

func resourceAlibabacloudStackAscmOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmOrganization(d.Get("name").(string))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if len(object.Data) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("org_id", object.Data[0].UUID)
	d.Set("name", object.Data[0].Name)
	d.Set("parent_id", strconv.Itoa(object.Data[0].ParentID))

	return nil
}

func resourceAlibabacloudStackAscmOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmOrganization(d.Get("name").(string))
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsOrganizationExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	addDebug("IsOrganizationExist", check, requestInfo, map[string]string{"id": d.Id()})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		if len(check.Data) != 0 {
			request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveOrganization", "/ascm/auth/organization/delete")
			request.QueryParams["id"] = d.Id()

			bresponse, err := client.ProcessCommonRequest(request)
			if err != nil {
				if bresponse == nil {
					return resource.RetryableError(err)
				}
				errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_organization", "RemoveOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			check, err = ascmService.DescribeAscmOrganization(d.Id())

			if err != nil {
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})
	return nil
}
