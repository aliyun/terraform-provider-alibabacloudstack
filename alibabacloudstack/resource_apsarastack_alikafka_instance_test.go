package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAlikafkaInstance_basic(t *testing.T) {
	var v *InstanceVO
	resourceId := "alibabacloudstack_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAlikafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-kafkainstan%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          "${var.name}",
					"vswitch_id":    "${alibabacloudstack_vpc_vswitch.default.id}",
					"zone_id":       "${data.alibabacloudstack_zones.default.zones.0.id}",
					"sasl":          "true",
					"plaintext":     "true",
					"sasl_ssl_port": "18080",
					"spec":          "Broker4C16G",
					"cup_type":      "Intel",
					"domain_refiex": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          name,
						"domain_refiex": name,
						"sasl":          "true",
						"sasl_ssl_port": "18080",
						"plaintext":     "true",
						"spec":          "Broker4C16G",
						"vpc_id":        CHECKSET,
						"vip_type":      CHECKSET,
						"status":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id": REMOVEKEY,
					"plaintext":  false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":     REMOVEKEY,
						"vswitch_id": REMOVEKEY,
						"plaintext":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"num_partitions":            "5",
					"auto_create_topics_enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"num_partitions":            "5",
						"auto_create_topics_enable": "true",
					}),
				),
			},
		},
	})
}

func resourceAlikafkaInstanceSimpleDependence(name string) string {
	return DataZoneCommonTestCase
}

func resourceAlikafkaInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	

%s
`, name, VSwitchCommonTestCase)
}

var alikafkaInstanceBasicMap = map[string]string{
	"cup_type":                   CHECKSET,
	"spec":                       CHECKSET,
	"replicas":                   CHECKSET,
	"disk_num":                   CHECKSET,
	"sasl":                       CHECKSET,
	"plaintext":                  CHECKSET,
	"message_max_bytes":          "10000000",
	"num_partitions":             "3",
	"auto_create_topics_enable":  "false",
	"num_io_threads":             "16",
	"queued_max_requests":        "80",
	"replica_fetch_wait_max_ms":  "500",
	"replica_lag_time_max_ms":    "30000",
	"num_network_threads":        "3",
	"log_retention_bytes":        "-1",
	"replica_fetch_max_bytes":    "10000000",
	"num_replica_fetchers":       "4",
	"default_replication_factor": "3",
	"offsets_retention_minutes":  "10080",
	"background_threads":         "10",
}
