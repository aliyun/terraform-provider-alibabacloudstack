package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

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
	request := client.NewCommonRequest("GET", "datahub", "2019-11-20", "OpenDataHubService", "")

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
			return dataHubClient.ProcessCommonRequest(request)
		})
		response, ok := raw.(*responses.CommonResponse)
		addDebug(action, raw, nil)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || errmsgs.NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_service", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("DatahubServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return err
	}
	d.SetId("DatahubServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
