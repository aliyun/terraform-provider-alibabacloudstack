package alibabacloudstack

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOnsGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOnsGroupId,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"read_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackOnsGroupCreate, resourceAlibabacloudStackOnsGroupRead, resourceAlibabacloudStackOnsGroupUpdate, resourceAlibabacloudStackOnsGroupDelete)
	return resource
}

func resourceAlibabacloudStackOnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ons.Client

	instanceId := d.Get("instance_id").(string)
	groupId := d.Get("group_id").(string)
	remark := d.Get("remark").(string)

	request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleGroupCreate", "")
	mergeMaps(request.QueryParams, map[string]string{
		"ProductName":  "Ons-inner",
		"PreventCache": "",
		"GroupId":      groupId,
		"Remark":       remark,
		"OnsRegionId":  client.RegionId,
		"InstanceId":   instanceId,
	})
	grp_resp := OGroup{}

	raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_group", "ConsoleGroupCreate", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ConsoleGroupCreate", raw, requestInfo, request)

	if !bresponse.IsSuccess() {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_group", "ConsoleGroupCreate", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &grp_resp)
	if grp_resp.Success != true {
		return errmsgs.WrapErrorf(errors.New(grp_resp.Message), errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_group", "ConsoleGroupCreate", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if err != nil {
		return errmsgs.WrapError(err)
	}

	log.Printf("groupid and instanceid %s %s", groupId, instanceId)
	d.SetId(groupId + COLON_SEPARATED + instanceId)

	return nil
}

func resourceAlibabacloudStackOnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}

	object, err := onsService.DescribeOnsGroup(d.Id())
	if err != nil {
		// Handle exceptions
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", object.Data[0].NamespaceID)
	d.Set("group_id", object.Data[0].GroupID)
	d.Set("remark", object.Data[0].Remark)

	return nil
}

func resourceAlibabacloudStackOnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlibabacloudStackOnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	var requestInfo *ons.Client
	check, err := onsService.DescribeOnsGroup(d.Id())
	parts, err := ParseResourceId(d.Id(), 2)

	if err != nil {
		errmsg := ""
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, parts[0], "IsGroupExist", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("IsGroupExist", check, requestInfo, map[string]string{"GroupId": parts[0]})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleGroupDelete", "")
		mergeMaps(request.QueryParams, map[string]string{
			"ProductName":  "Ons-inner",
			"PreventCache": "",
			"GroupId":      parts[0],
			"OnsRegionId":  client.RegionId,
			"InstanceId":   parts[1],
		})

		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_group", "ConsoleGroupDelete", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		check, err = onsService.DescribeOnsGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return nil
}
