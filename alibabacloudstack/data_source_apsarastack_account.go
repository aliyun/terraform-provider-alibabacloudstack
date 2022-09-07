package alibabacloudstack

import (
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAccountRead,

		Schema: map[string]*schema.Schema{
			// Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlibabacloudStackAccountRead(d *schema.ResourceData, meta interface{}) error {
	accountId, err := meta.(*connectivity.AlibabacloudStackClient).AccountId()

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] alibabacloudstack_account - account ID found: %#v", accountId)

	d.SetId(accountId)

	return nil
}
