package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAccountRead,

		Schema: map[string]*schema.Schema{
			// Computed values
			"login_name" :{
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
	resp, err := client.GetCallerInfo()
	if err != nil {
		return  err
	}
	ownerId := resp["primaryKey"].(string)

	if ownerId == "" {
		return  fmt.Errorf("ownerId not found")
	}

	d.SetId(ownerId)
	d.Set("login_name", resp["loginName"])
	d.Set("organization_id", client.Department)
	d.Set("region", client.RegionId)

	return nil
}
