package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
				Type:       schema.TypeString,
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
	name := d.Get("name").(string)
	var organizationId string
	if _, ok := d.GetOk("organization_id"); ok {
		organizationId = d.Get("organization_id").(string)
	} else {
		organizationId = client.Department
	}
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateResourceGroup", "/ascm/auth/resource_group/create_resource_group")
	request.QueryParams["ProductName"] = "ascm"
	request.QueryParams["resource_group_name"] = name
	request.QueryParams["organization_id"] = organizationId
	request.QueryParams["OrganizationId"] = organizationId
	request.QueryParams["Department"] = organizationId
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
	response := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	resourceGroupId, err := jsonpath.Get("$.data.id", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group", "CreateResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, "resource group id not found in response")
	}
	resourceId := fmt.Sprintf("%s:%s", organizationId, fmt.Sprint(resourceGroupId))
	d.SetId(resourceId)
	return nil
}

func resourceAlibabacloudStackAscmResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if d.IsNewResource() {
		return nil
	}
	if d.HasChange("name") {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateResourceGroup", "/ascm/auth/resource_group/update_resource_group")
		request.QueryParams["resourceGroupName"] = d.Get("name").(string)
		if len(did) < 2 {
			return fmt.Errorf("invalid resource group id format, expected 2 parts separated by colon, got: %s", d.Id())
		}
		request.QueryParams["id"] = did[1]
		request.QueryParams["OrganizationId"] = did[0]
		request.QueryParams["Department"] = did[0]
		request.QueryParams["ResourceGroup"] = did[1]
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

	return nil
}

func resourceAlibabacloudStackAscmResourceGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmResourceGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object.ResourceGroupName)
	d.Set("rg_id", strconv.Itoa(object.ID))
	d.Set("organization_id", strconv.Itoa(object.OrganizationID))

	return nil
}

func resourceAlibabacloudStackAscmResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	did := strings.Split(d.Id(), COLON_SEPARATED)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveResourceGroup", "/ascm/auth/resource_group/delete_resource_group")
	request.QueryParams["OrganizationId"] = did[0]
	request.QueryParams["Department"] = did[0]
	request.QueryParams["ResourceGroup"] = did[1]
	request.QueryParams["resource_group_id"] = did[1]
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("RemoveResourceGroup", bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group", "RemoveResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}
