package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"fmt"
	"testing"
)

func TestAccApsaraStackDnsRecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(dataSourceApsaraStackDnsRecord, acctest.RandIntRange(1000000, 9999999)),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_dns_records.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.record_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.domain_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.host_record"),
					resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.rr_set"),
					resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.ttl"),
				),
			},
		},
	})
}

const dataSourceApsaraStackDnsRecord = `

resource "apsarastack_dns_domain" "default" {
 domain_name = "tf-testdnsrecordbasic-%d."
 remark = "tf-testdnsrecordbasic"
}
resource "apsarastack_dns_record" "default" {
 domain_id   = apsarastack_dns_domain.default.domain_id
 host_record = "testrecord"
 type        = "A"
 ttl         = 300
 rr_set      = ["192.168.2.4","192.168.2.7","10.0.0.4"]
}

data "apsarastack_dns_records" "default"{
 domain_id         = apsarastack_dns_record.default.domain_id
 host_record_regex = apsarastack_dns_record.default.host_record
}
`
