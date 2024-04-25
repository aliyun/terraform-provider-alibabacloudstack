package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmResourceGroupCreate,
		Read:   resourceAlibabacloudStackAscmResourceGroupRead,
		Update: resourceAlibabacloudStackAscmResourceGroupUpdate,
		Delete: resourceAlibabacloudStackAscmResourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rg_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackAscmResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	name := d.Get("name").(string)
	check, err := ascmService.DescribeAscmResourceGroup(name)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_resource_group", "RG alreadyExist", AlibabacloudStackSdkGoERROR)
	}
	organizationid := d.Get("organization_id").(string)

	if len(check.Data) == 0 {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":            client.RegionId,
			"AccessKeySecret":     client.SecretKey,
			"Product":             "Ascm",
			"Action":              "CreateResourceGroup",
			"Version":             "2019-05-10",
			"ProductName":         "ascm",
			"resource_group_name": name,
			"organization_id":     organizationid,
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "CreateResourceGroup"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId, "x-acs-content-type": "application/json", "Content-Type": "application/json"}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf("response of raw CreateResourceGroup : %s", raw)
		addDebug("CreateResourceGroup", raw, request, request.QueryParams)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_resource_group", "CreateResourceGroup", raw)
		}

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_resource_group", "CreateResourceGroup", AlibabacloudStackSdkGoERROR)
		}
		addDebug("CreateResourceGroup", raw, requestInfo, bresponse.GetHttpContentString())
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmResourceGroup(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	d.SetId(check.Data[0].ResourceGroupName + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))

	return resourceAlibabacloudStackAscmResourceGroupUpdate(d, meta)

}

func resourceAlibabacloudStackAscmResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	name := d.Get("name").(string)
	attributeUpdate := false
	check, err := ascmService.DescribeAscmResourceGroup(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsResourceGroupExist", AlibabacloudStackSdkGoERROR)
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

	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":          client.RegionId,
		"AccessKeySecret":   client.SecretKey,
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"Product":           "ascm",
		"Action":            "UpdateResourceGroup",
		"Version":           "2019-05-10",
		"resourceGroupName": name,
		"id":                did[1],
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.SetHTTPSInsecure(true)
	request.ApiName = "UpdateResourceGroup"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId, "x-acs-content-type": "application/json", "Content-Type": "application/json"}

	if attributeUpdate {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw UpdateResourceGroup : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_resource_group", "UpdateResourceGroup", raw)
		}
		addDebug(request.GetActionName(), raw, request)

	}
	d.SetId(name + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))

	return resourceAlibabacloudStackAscmResourceGroupRead(d, meta)

}

func resourceAlibabacloudStackAscmResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmResourceGroup(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsResourceGroupExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsResourceGroupExist", check, requestInfo, map[string]string{"resourceGroupName": did[0]})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":          client.RegionId,
			"AccessKeySecret":   client.SecretKey,
			"Product":           "ascm",
			"Action":            "RemoveResourceGroup",
			"Version":           "2019-05-10",
			"ProductName":       "ascm",
			"resourceGroupName": did[0],
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "RemoveResourceGroup"
		request.Headers = map[string]string{"RegionId": client.RegionId, "x-acs-content-type": "application/json", "Content-Type": "application/json"}
		request.RegionId = client.RegionId
		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		log.Printf(" response of raw RemoveResourceGroup : %s", raw)
		_, err = ascmService.DescribeAscmResourceGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}
