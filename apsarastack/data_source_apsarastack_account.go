package apsarastack

import (
	"log"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApsaraStackAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackAccountRead,

		Schema: map[string]*schema.Schema{
			// Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApsaraStackAccountRead(d *schema.ResourceData, meta interface{}) error {
	accountId, err := meta.(*connectivity.ApsaraStackClient).AccountId()

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] apsarastack_account - account ID found: %#v", accountId)

	d.SetId(accountId)

	return nil
}
