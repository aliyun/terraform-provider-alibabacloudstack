package alibabacloudstack

import (
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDnsRecordDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(dataSourceAlibabacloudStackDnsRecord, getAccTestRandInt(1000000, 9999999)),
				Check:  resource.ComposeTestCheckFunc(

				//	testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_dns_records.default"),
				//	resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_records.default", "records.record_id"),
				//	resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_records.default", "records.domain_id"),
				//	resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_records.default", "records.host_record"),
				//	resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_records.default", "records.type"),
				//	resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_records.default", "records.rr_set"),
				//	resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_records.default", "records.ttl"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackDnsRecord = `

resource "alibabacloudstack_dns_domain" "default" {
 domain_name = "tf-testdnsrecordbasic-%d."
 remark = "tf-testdnsrecordbasic"
}
resource "alibabacloudstack_dns_record" "default" {
 zone_id   = alibabacloudstack_dns_domain.default.domain_id
lba_strategy = "ALL_RR",
 name = "testrecord"
 type        = "A"
 ttl         = 300
 rr_set      = ["192.168.2.4","192.168.2.7","10.0.0.4"]
}

data "alibabacloudstack_dns_records" "default"{
 zone_id         = alibabacloudstack_dns_record.default.zone_id
 name = alibabacloudstack_dns_record.default.name
}
`
