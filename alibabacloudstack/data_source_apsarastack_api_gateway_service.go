package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackApiGatewayService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackApigatewayServiceRead,
		DeprecationMessage: "The 'alibabacloudstack_api_gateway_service' resource is unsupported on ApsaraStack and will be removed in version 3.21.0.",
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

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	response, err := client.DoTeaRequest("POST", "CloudAPI", "2016-07-14", "OpenApiGatewayService", "", nil, nil, request)
	addDebug("OpenApiGatewayService", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("ApiGatewayServicHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return err
	}

	d.SetId(fmt.Sprintf("%v", response["OrderId"]))
	d.Set("status", "Opened")

	return nil
}
