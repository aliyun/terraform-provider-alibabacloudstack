package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	action := "OpenOtsService"
	response, err := client.DoTeaRequest("POST", "Ots", "2016-06-20", action, "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("OtsServicHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return err
	}
	d.SetId(fmt.Sprintf("%v", response["OrderId"]))
	d.Set("status", "Opened")
	return nil
}
