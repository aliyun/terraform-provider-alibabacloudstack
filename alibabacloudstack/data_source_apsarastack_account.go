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
			"organization_id" : {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region" : {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlibabacloudStackAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	accountId, err := client.AccountId()

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] alibabacloudstack_account - account ID found: %#v", accountId)

	d.SetId(accountId)
	d.Set("organization_id", client.Department)
	d.Set("region", client.RegionId)

	return nil
}
