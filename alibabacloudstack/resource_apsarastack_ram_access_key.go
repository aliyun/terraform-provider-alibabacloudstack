package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudstackRamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudstackAscmAccessKeyCreate,
		Read:   resourceAlibabacloudstackAscmAccessKeyRead,
		Update: resourceAlibabacloudstackAscmAccessKeyUpdate,
		Delete: resourceAlibabacloudstackAscmAccessKeyDelete,

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
				Default:      Active,
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
}

func resourceAlibabacloudstackAscmAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	//  userName:=d.Get("user_name")

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2015-05-01"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "RamCreateAccessKey"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "RamCreateAccessKey",
		"Version":         "2015-05-01",
		"ProductName":     "ascm",
		//"": userName,
	}
	var response = AccessKeyInCreateAccessKey{}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_access_key", "", raw)
	}

	addDebug("RamCreateAccessKey", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_access_key", "", AlibabacloudStackSdkGoERROR)
	}
	addDebug("RamCreateAccessKey", raw, requestInfo, bresponse.GetHttpContentString())
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	d.SetId(fmt.Sprint(response.AccessKeyId))

	return resourceAlibabacloudstackAscmAccessKeyUpdate(d, meta)
}

func resourceAlibabacloudstackAscmAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudstackAscmAccessKeyRead(d, meta)
}

func resourceAlibabacloudstackAscmAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmKeypolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", object) //add read data from the struct
	return nil
}

func resourceAlibabacloudstackAscmAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmKeypolicy(d.Id())

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsResourceGroupExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsKeyExist", check, requestInfo, map[string]string{"resourceGroupName": d.Id()})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			
			"Product":         "ascm",
			"Action":          "RamDeleteAccessKey",
			"Version":         "2015-05-01",
			"ProductName":     "ascm",
			"id":              d.Id(),
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "RamDeleteAccessKey"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		_, err = ascmService.DescribeAscmKeypolicy(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ram_access_key", "RamDeleteAccessKey", AlibabacloudStackSdkGoERROR)
	}
	return nil

}

func (s *AscmService) DescribeAscmKeypolicy(id string) (response *AccessKeyInCreateAccessKey, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "RamListAccessKeys",
		"Version":         "2015-05-01",
		"id":              id,
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2015-05-01"
	request.ServiceCode = "ascm"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "RamListAccessKeys"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	request.Domain = s.client.Domain
	var value = &AccessKeyInCreateAccessKey{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorQuotaNotFound"}) {
			return value, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return value, WrapErrorf(err, DefaultErrorMsg, "GetQuota", AlibabacloudStackSdkGoERROR)

	}
	addDebug("RamListAccessKeys", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), value)
	if err != nil {
		return value, WrapError(err)
	}
	if value.AccessKeyId == "200" {
		return value, WrapError(err)
	}

	return value, nil
}

type AccessKeyInCreateAccessKey struct {
	AccessKeyId     string `json:"AccessKeyId" xml:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret" xml:"AccessKeySecret"`
	Status          string `json:"Status" xml:"Status"`
	CreateDate      string `json:"CreateDate" xml:"CreateDate"`
}
