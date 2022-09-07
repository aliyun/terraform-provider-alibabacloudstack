package alibabacloudstack

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackApiGatewayService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackApigatewayServiceRead,

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
func dataSourceAlibabacloudStackApigatewayServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("ApiGatewayServicHasNotBeenOpened")
		d.Set("status", "")

		return nil
	}
	if v, ok := d.GetOk("enable"); ok {
		request["status"] = v
	}

	action := "OpenApiGatewayService"
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	conn, err := meta.(*connectivity.AlibabacloudStackClient).NewTeaCommonClient(connectivity.OpenApiGatewayService)
	if err != nil {
		return WrapError(err)
	}
	request["Product"] = "CloudAPI"
	request["OrganizationId"] = client.Department
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-07-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

	addDebug("OpenApiGatewayService", response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("ApiGatewayServicHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_api_gateway_service", "OpenApiGatewayService", AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v", response["OrderId"]))
	d.Set("status", "Opened")

	return nil
}
