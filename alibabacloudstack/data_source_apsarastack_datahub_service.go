package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"strings"
	"time"

	_ "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDatahubService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDatahubServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlibabacloudStackDatahubServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("DatahubServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenDataHubService"
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "datahub"
	request.Version = "2019-11-20"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}

	request.RegionId = client.RegionId
	request.ApiName = "OpenDataHubService"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeyId":     client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"Product":         "datahub",
		"RegionId":        client.RegionId,
		"Action":          "OpenDataHubService",
		"Version":         "2019-11-20",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
			return dataHubClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("DatahubServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_datahub_service", action, AlibabacloudStackSdkGoERROR)
	}
	d.SetId("DatahubServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
