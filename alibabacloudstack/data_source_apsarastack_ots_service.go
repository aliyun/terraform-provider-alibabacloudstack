package alibabacloudstack

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOtsService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOtsServiceRead,

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
func dataSourceAlibabacloudStackOtsServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("OtsServicHasNotBeenOpened")
		d.Set("status", "")
		request["enable"] = v
		return nil
	}
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}
	action := "OpenOtsService"
	request["RegionId"] = client.RegionId
	request["Product"] = "Ots"
	request["OrganizationId"] = client.Department
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}

		addDebug(action, response, nil)

		d.SetId(fmt.Sprintf("%v", response["OrderId"]))
		d.Set("status", "Opened")
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("OtsServicHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_ots_service", "OpenOtsService", AlibabacloudStackSdkGoERROR)
	}

	return nil
}
