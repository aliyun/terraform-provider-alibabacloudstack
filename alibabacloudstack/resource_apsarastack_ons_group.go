package alibabacloudstack

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackOnsGroupCreate,
		Read:   resourceAlibabacloudStackOnsGroupRead,
		Update: resourceAlibabacloudStackOnsGroupUpdate,
		Delete: resourceAlibabacloudStackOnsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
}

func resourceAlibabacloudStackOnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client

	instanceId := d.Get("instance_id").(string)
	groupId := d.Get("group_id").(string)
	remark := d.Get("remark").(string)
	request := requests.NewCommonRequest()

	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "Ons-inner",
		"Action":          "ConsoleGroupCreate",
		"Version":         "2018-02-05",
		"ProductName":     "Ons-inner",
		"PreventCache":    "",
		"GroupId":         groupId,
		"Remark":          remark,
		"OnsRegionId":     client.RegionId,
		"InstanceId":      instanceId,
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ConsoleGroupCreate"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	grp_resp := OGroup{}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ons_group", "ConsoleGroupCreate", raw)
	}
	addDebug("ConsoleGroupCreate", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.IsSuccess() != true {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ons_group", "ConsoleGroupCreate", AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &grp_resp)
	if grp_resp.Success != true {
		return WrapErrorf(errors.New(grp_resp.Message), DefaultErrorMsg, "alibabacloudstack_ons_group", "ConsoleGroupCreate", AlibabacloudStackSdkGoERROR)
	}

	if err != nil {
		return WrapError(err)
	}

	log.Printf("groupid and instanceid %s %s", groupId, instanceId)
	d.SetId(groupId + COLON_SEPARATED + instanceId)

	return resourceAlibabacloudStackOnsGroupRead(d, meta)
}

func resourceAlibabacloudStackOnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}

	object, err := onsService.DescribeOnsGroup(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.Data[0].NamespaceID)
	d.Set("group_id", object.Data[0].GroupID)
	d.Set("remark", object.Data[0].Remark)

	return nil
}

func resourceAlibabacloudStackOnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudStackOnsGroupRead(d, meta)
}

func resourceAlibabacloudStackOnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	var requestInfo *ecs.Client
	check, err := onsService.DescribeOnsGroup(d.Id())
	parts, err := ParseResourceId(d.Id(), 2)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, parts[0], "IsGroupExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsGroupExist", check, requestInfo, map[string]string{"GroupId": parts[0]})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ons-inner",
			"Action":          "ConsoleGroupDelete",
			"Version":         "2018-02-05",
			"ProductName":     "Ons-inner",
			"PreventCache":    "",
			"GroupId":         parts[0],
			"OnsRegionId":     client.RegionId,
			"InstanceId":      parts[1],
		}

		request.Method = "POST"
		request.Product = "Ons-inner"
		request.Version = "2018-02-05"
		request.ServiceCode = "Ons-inner"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "ConsoleGroupDelete"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = onsService.DescribeOnsGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return nil
}
