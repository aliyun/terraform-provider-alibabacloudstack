package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
			"primary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aliyunid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmOrganizationCreate, resourceAlibabacloudStackAscmOrganizationRead, resourceAlibabacloudStackAscmOrganizationUpdate, resourceAlibabacloudStackAscmOrganizationDelete)
	return resource
}

func resourceAlibabacloudStackAscmOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	parentid := d.Get("parent_id").(string)
	var requestInfo *ecs.Client
	name := d.Get("name").(string)
	object, err := ascmService.DescribeAscmOrganizationByName(parentid, name)
	if err != nil && !errmsgs.NotFoundError(err) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_organization", "ORG alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	var resourceId string
	if object == nil {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateOrganization", "/ascm/auth/organization/add")
		request.QueryParams["parentId"] = parentid
		request.QueryParams["name"] = name
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug("CreateOrganization", bresponse, requestInfo, request)

		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_organization", "CreateOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		if bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_organization", "CreateOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		// TODO: Parent organization not found will not report an error here, because HttpStatus is still 200
		addDebug("CreateOrganization", bresponse, requestInfo, bresponse.GetHttpContentString())
		response := make(map[string]interface{})
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		id, err := jsonpath.Get("$.data.id", response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		resourceId = fmt.Sprint(id)
	} else {
		resourceId = fmt.Sprint(object.ID)
	}
	d.SetId(resourceId)
	return nil
}

func resourceAlibabacloudStackAscmOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChange("name") && !d.IsNewResource() {

		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateOrganization", "/ascm/auth/organization/update")
		request.QueryParams["id"] = d.Id()
		request.QueryParams["name"] = d.Get("name").(string)

		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf(" response of raw UpdateOrganization : %s", bresponse)

		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_instance", "UpdateOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), bresponse, request)
	}
	return nil
}

func resourceAlibabacloudStackAscmOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmOrganization(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("org_id", object.UUID)
	d.Set("name", object.Name)
	d.Set("parent_id", strconv.Itoa(object.ParentID))
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetPrivateCloudAccountByOrganizationId", "/ascm/auth/user/getPrivateCloudAccountByOrganizationId")
	request.QueryParams["organizationId"] = d.Id()
	request.QueryParams["OrganizationId"] = d.Id()
	bresponse, err := client.ProcessCommonRequest(request)
	log.Printf(" response of raw UpdateOrganization : %s", bresponse)

	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_organizations", "GetPrivateCloudAccountByOrganizationId", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	} else {
		var resp OrganizationIdResponse
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("aliyunid", resp.Data.AliyunId)
		d.Set("primary_key", resp.Data.PrimaryKey)
	}
	return nil
}

func resourceAlibabacloudStackAscmOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmOrganization(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsOrganizationExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		if object != nil {
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
			return nil
		}
		return nil
	})
	return err
}
