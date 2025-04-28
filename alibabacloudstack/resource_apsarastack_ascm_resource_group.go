package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmResourceGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"organization_id": {
				Type:     schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field 'organization_id' has been deprecated. Use the organization to which the current user belongs",
			},
			"rg_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmResourceGroupCreate, resourceAlibabacloudStackAscmResourceGroupRead, resourceAlibabacloudStackAscmResourceGroupUpdate, resourceAlibabacloudStackAscmResourceGroupDelete)
	return resource
}

func resourceAlibabacloudStackAscmResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	name := d.Get("name").(string)
	check, err := ascmService.DescribeAscmResourceGroup(name)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_resource_group", "RG alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	var organizationId string
	if _, ok:= d.GetOk("organization_id"); ok {
		organizationId = d.Get("organization_id").(string)
	} else {
		organizationId = client.Department
	}

	if len(check.Data) == 0 {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateResourceGroup", "/ascm/auth/resource_group/create_resource_group")
		request.QueryParams["ProductName"] = "ascm"
		request.QueryParams["resource_group_name"] = name
		request.QueryParams["organization_id"] = organizationId
		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug("CreateResourceGroup", bresponse, request, request.QueryParams)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group", "CreateResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		if bresponse.GetHttpStatus() != 200 {
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group", "CreateResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("CreateResourceGroup", bresponse, requestInfo, bresponse.GetHttpContentString())
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmResourceGroup(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	d.SetId(check.Data[0].ResourceGroupName + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))

	return nil
}

func resourceAlibabacloudStackAscmResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	name := d.Get("name").(string)
	attributeUpdate := false
	check, err := ascmService.DescribeAscmResourceGroup(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsResourceGroupExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data[0].ResourceGroupName = name
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data[0].ResourceGroupName = name
	}

	if attributeUpdate {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateResourceGroup", "/ascm/auth/resource_group/update_resource_group")
		request.QueryParams["resourceGroupName"] = name
		request.QueryParams["id"] = did[1]
		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"

		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf(" response of raw UpdateResourceGroup : %s", bresponse)

		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group", "UpdateResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), bresponse, request)
	}
	d.SetId(name + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))

	return nil
}

func resourceAlibabacloudStackAscmResourceGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmResourceGroup(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
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

	d.Set("name", did[0])
	d.Set("rg_id", did[1])
	d.Set("organization_id", strconv.Itoa(object.Data[0].OrganizationID))

	return nil
}

func resourceAlibabacloudStackAscmResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmResourceGroup(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsResourceGroupExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsResourceGroupExist", check, requestInfo, map[string]string{"resourceGroupName": did[0]})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveResourceGroup", "/ascm/auth/resource_group/delete_resource_group")
		request.QueryParams["resourceGroupName"] = did[0]
		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group", "RemoveResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		log.Printf(" response of raw RemoveResourceGroup : %s", bresponse)
		_, err = ascmService.DescribeAscmResourceGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}