package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudstackAscmAccessKey() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Active",
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"pgp_key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudstackAscmAccessKeyCreate, resourceAlibabacloudstackAscmAccessKeyRead, nil, resourceAlibabacloudstackAscmAccessKeyDelete)
	return resource
}

func resourceAlibabacloudstackAscmAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "Ascm", "2015-05-01", "RamCreateAccessKey", "")
	request.SetDomain(client.Config.Endpoints[connectivity.ASAPICode])
	var response = AccessKeyInCreateAccessKey{}
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_access_key", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	d.SetId(fmt.Sprint(response.AccessKeyId))

	return nil
}

func resourceAlibabacloudstackAscmAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmKeypolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("status", object) //add read data from the struct
	return nil
}

func resourceAlibabacloudstackAscmAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	check, err := ascmService.DescribeAscmKeypolicy(d.Id())

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsResourceGroupExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsKeyExist", check, nil, map[string]string{"resourceGroupName": d.Id()})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RamDeleteAccessKey", "/ascm/auth/ramCompatible/deleteAccessKey")
		request.SetDomain(client.Config.Endpoints[connectivity.ASAPICode])
		request.QueryParams["id"] = d.Id()

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_access_key", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		_, err = ascmService.DescribeAscmKeypolicy(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ram_access_key", "RamDeleteAccessKey", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func (s *AscmService) DescribeAscmKeypolicy(id string) (response *AccessKeyInCreateAccessKey, err error) {
	request := s.client.NewCommonRequest("POST", "Ascm", "2015-05-01", "RamListAccessKeys", "")
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])
	request.QueryParams["id"] = id

	var value = &AccessKeyInCreateAccessKey{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return value, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return value, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_access_key", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), value)
	if err != nil {
		return value, errmsgs.WrapError(err)
	}
	if value.AccessKeyId == "200" {
		return value, errmsgs.WrapError(err)
	}

	return value, nil
}

type AccessKeyInCreateAccessKey struct {
	AccessKeyId     string `json:"AccessKeyId" xml:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret" xml:"AccessKeySecret"`
	Status          string `json:"Status" xml:"Status"`
	CreateDate      string `json:"CreateDate" xml:"CreateDate"`
}
