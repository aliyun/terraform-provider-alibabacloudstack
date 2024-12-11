package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRedisTairinstance0(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "Tair_rdb_包年包月实例续费",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "Tair_rdb_包年包月实例续费",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_class": "tair.rdb.2g",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_class": "tair.rdb.2g",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccRedisTairinstanceCheckmap = map[string]string{

	"port": CHECKSET,

	"tair_instance_id": CHECKSET,

	"capacity": CHECKSET,

	"qps": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"engine_version": CHECKSET,

	"max_connections": CHECKSET,

	"vswitch_id": CHECKSET,

	"vpc_id": CHECKSET,

	"node_type": CHECKSET,

	"end_time": CHECKSET,

	"connection_domain": CHECKSET,

	"maintain_end_time": CHECKSET,

	"network_type": CHECKSET,

	"bandwidth": CHECKSET,

	"payment_type": CHECKSET,

	"instance_type": CHECKSET,

	"maintain_start_time": CHECKSET,

	"architecture_type": CHECKSET,

	"ssl_enabled": CHECKSET,

	"create_time": CHECKSET,

	"instance_class": CHECKSET,

	"secondary_zone_id": CHECKSET,

	"total_count": CHECKSET,

	"vpc_auth_mode": CHECKSET,

	"tair_instance_name": CHECKSET,

	"region_id": CHECKSET,
}

func AlibabacloudTestAccRedisTairinstanceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s


variable "zone_id" {
    default = cn-beijing-h
}

variable "region_id" {
    default = cn-beijing
}




`, name, DataZoneCommonTestCase)
}
func TestAccAlibabacloudStackRedisTairinstance1(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PostPaid",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "cn-hangzhou-h",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"vpc_id": "alibabacloudstack_vpc.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PostPaid",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "cn-hangzhou-h",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"vpc_id": "alibabacloudstack_vpc.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "test_new_name",

					"engine_version": "5.0",

					"instance_class": "tair.rdb.2g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "test_new_name",

						"engine_version": "5.0",

						"instance_class": "tair.rdb.2g",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance2(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance3(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "test_new_name",

					"engine_version": "5.0",

					"instance_class": "tair.rdb.2g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "test_new_name",

						"engine_version": "5.0",

						"instance_class": "tair.rdb.2g",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance4(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance5(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"instance_class": "tair.rdb.1g",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"tair_instance_name": "tags_test11_ljt",

					"engine_version": "5.0",

					"secondary_zone_id": "cn-huhehaote-b",

					"instance_type": "tair_rdb",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"instance_class": "tair.rdb.1g",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"tair_instance_name": "tags_test11_ljt",

						"engine_version": "5.0",

						"secondary_zone_id": "cn-huhehaote-b",

						"instance_type": "tair_rdb",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance6(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"engine_version": "1.0",

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_scm",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"secondary_zone_id": "cn-hangzhou-h",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"tair_instance_name": "my_scm_test",

					"instance_class": "tair.scm.standard.1m.4d",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"engine_version": "1.0",

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_scm",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"secondary_zone_id": "cn-hangzhou-h",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"tair_instance_name": "my_scm_test",

						"instance_class": "tair.scm.standard.1m.4d",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance7(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_scm",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.scm.standard.1m.4d",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_scm",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.scm.standard.1m.4d",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance8(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"instance_class": "tair.essd.standard.xlarge",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"instance_class": "tair.essd.standard.xlarge",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance9(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance10(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance11(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "redis_governance_test",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "redis_governance_test",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance12(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance13(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "theTestInstance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "theTestInstance",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance14(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "cn-hangzhou-h",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"vpc_id": "alibabacloudstack_vpc.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "cn-hangzhou-h",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"vpc_id": "alibabacloudstack_vpc.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "test_new_name",

					"engine_version": "5.0",

					"instance_class": "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "test_new_name",

						"engine_version": "5.0",

						"instance_class": "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance15(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"instance_class": "tair.rdb.with.proxy.1g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"instance_class": "tair.rdb.with.proxy.1g",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance16(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.2g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.2g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"engine_version": "5.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance17(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.2g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.2g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"engine_version": "5.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance18(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.2g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.2g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "tf-testacc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "tf-testacc",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance19(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.2g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.2g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "tf-testacc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "tf-testacc",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance20(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.2g",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.2g",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"instance_class": "tair.rdb.4g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"instance_class": "tair.rdb.4g",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance21(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.essd.standard.xlarge",

					"tair_instance_name": "tf_test_ins",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.essd.standard.xlarge",

						"tair_instance_name": "tf_test_ins",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance22(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

					"instance_class": "tair.essd.standard.xlarge",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVSwitch.VSwitchId)}}",

						"instance_class": "tair.essd.standard.xlarge",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance23(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"instance_class": "tair.rdb.1g",

					"instance_type": "tair_rdb",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"instance_class": "tair.rdb.1g",

						"instance_type": "tair_rdb",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance24(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.with.proxy.2g",

					"tair_instance_name": "test_tf_tair_rdb_rw",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.with.proxy.2g",

						"tair_instance_name": "test_tf_tair_rdb_rw",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance25(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_scm",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.scm.with.proxy.standard.2m.8d",

					"tair_instance_name": "test_tf_tair_scm",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "1.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_scm",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.scm.with.proxy.standard.2m.8d",

						"tair_instance_name": "test_tf_tair_scm",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "1.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"engine_version": "1.0",

					"instance_class": "tair.scm.with.proxy.standard.2m.8d",

					"tair_instance_name": "test_tair_scm3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"engine_version": "1.0",

						"instance_class": "tair.scm.with.proxy.standard.2m.8d",

						"tair_instance_name": "test_tair_scm3",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance26(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_scm",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.scm.standard.1m.4d",

					"tair_instance_name": "test_tf_tair_scm",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"engine_version": "1.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_scm",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.scm.standard.1m.4d",

						"tair_instance_name": "test_tf_tair_scm",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"engine_version": "1.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance27(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_scm",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.scm.standard.1m.4d",

					"tair_instance_name": "test_tf_tair_scm",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "1.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_scm",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.scm.standard.1m.4d",

						"tair_instance_name": "test_tf_tair_scm",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "1.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance28(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance29(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_规格变配",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_规格变配",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "Tair_rdb_规格变配_不强制",

					"instance_class": "tair.rdb.2g",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "Tair_rdb_规格变配_不强制",

						"instance_class": "tair.rdb.2g",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance30(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb单副本",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"node_type": "STAND_ALONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb单副本",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"node_type": "STAND_ALONE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance31(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.essd.standard.xlarge",

					"tair_instance_name": "test_tf_tair_essd",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "1.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.essd.standard.xlarge",

						"tair_instance_name": "test_tf_tair_essd",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "1.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance32(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "Tair_rdb_修改密码",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "STAND_ALONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "Tair_rdb_修改密码",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "STAND_ALONE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "test_rename3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "test_rename3",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance33(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_规格变配",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_规格变配",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"tair_instance_name": "test_rename3",

					"instance_class": "tair.rdb.2g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tair_instance_name": "test_rename3",

						"instance_class": "tair.rdb.2g",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance34(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_essd",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"node_type": "STAND_ALONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_essd",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"node_type": "STAND_ALONE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance35(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_essd",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"node_type": "STAND_ALONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_essd",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"node_type": "STAND_ALONE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance36(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "Tair_cluster_rdb_增加分片",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "Tair_cluster_rdb_增加分片",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance37(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period1",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period1",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance38(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period2",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period2",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance39(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period3",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period3",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance40(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period6",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period6",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance41(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period7",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period7",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance42(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period4",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period4",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance43(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period5",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period5",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance44(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period8",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period8",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance45(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period9",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period9",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance46(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period24",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period24",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance47(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period36",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period36",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance48(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_period60",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_period60",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance49(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.essd.standard.2xlarge",

					"tair_instance_name": "test_tf_tair_essd_pl2",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "1.0",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.essd.standard.2xlarge",

						"tair_instance_name": "test_tf_tair_essd_pl2",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "1.0",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance50(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.essd.standard.4xlarge",

					"tair_instance_name": "test_tf_tair_essd_pl3",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "1.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.essd.standard.4xlarge",

						"tair_instance_name": "test_tf_tair_essd_pl3",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "1.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance51(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"tair_instance_name": "Tair_rdb_升级大版本",

					"instance_type": "tair_rdb",

					"engine_version": "5.0",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"instance_class": "tair.rdb.1g",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"tair_instance_name": "Tair_rdb_升级大版本",

						"instance_type": "tair_rdb",

						"engine_version": "5.0",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"instance_class": "tair.rdb.1g",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"engine_version": "6.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"engine_version": "6.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance52(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_scm",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.scm.standard.1m.4d",

					"tair_instance_name": "test_tf_tair_scm",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "1.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "STAND_ALONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_scm",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.scm.standard.1m.4d",

						"tair_instance_name": "test_tf_tair_scm",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "1.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "STAND_ALONE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance53(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.with.proxy.1g",

					"tair_instance_name": "Tair_rdb_rw_双可用区_修改slave只读节点",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::secondaryvsw.ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.with.proxy.1g",

						"tair_instance_name": "Tair_rdb_rw_双可用区_修改slave只读节点",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::secondaryvsw.ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance54(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.2g",

					"tair_instance_name": "test_tf_cluster_rdb",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.2g",

						"tair_instance_name": "test_tf_cluster_rdb",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance55(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "STAND_ALONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "STAND_ALONE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"node_type": "MASTER_SLAVE",

					"instance_class": "tair.rdb.2g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"node_type": "MASTER_SLAVE",

						"instance_class": "tair.rdb.2g",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance56(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_scm",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.scm.standard.1m.4d",

					"tair_instance_name": "test_tf_tair_scm_关联安全组",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "1.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_scm",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.scm.standard.1m.4d",

						"tair_instance_name": "test_tf_tair_scm_关联安全组",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "1.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance57(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "Tair_rdb_包年包月_转按量付费",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "Tair_rdb_包年包月_转按量付费",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance58(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance59(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb_开启TLS加密",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",

					"ssl_enabled": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb_开启TLS加密",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",

						"ssl_enabled": "Disable",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"ssl_enabled": "Enable",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ssl_enabled": "Enable",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance60(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.with.proxy.1g",

					"tair_instance_name": "test_tf_tair_rdb_rw",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.with.proxy.1g",

						"tair_instance_name": "test_tf_tair_rdb_rw",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance61(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance62(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "6.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "6.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance63(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_type": "tair_scm",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.scm.standard.1m.4d",

					"tair_instance_name": "test_tf_tair_scm",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"engine_version": "1.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_type": "tair_scm",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.scm.standard.1m.4d",

						"tair_instance_name": "test_tf_tair_scm",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"engine_version": "1.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",

					"instance_class": "tair.scm.standard.2m.8d",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",

						"instance_class": "tair.scm.standard.2m.8d",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance64(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "test_tf_tair_rdb单副本",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"node_type": "STAND_ALONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "test_tf_tair_rdb单副本",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"node_type": "STAND_ALONE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance65(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "redis.shard.micro.ce",

					"tair_instance_name": "ProductType_Redis_shard_ce",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"node_type": "STAND_ALONE",

					"network_type": "VPC",

					"instance_type": "Redis",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "redis.shard.micro.ce",

						"tair_instance_name": "ProductType_Redis_shard_ce",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"node_type": "STAND_ALONE",

						"network_type": "VPC",

						"instance_type": "Redis",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance66(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"tair_instance_name": "tair_rdb_白名单",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",

					"secondary_zone_id": "${{ref(variable, ZoneId)}}",

					"node_type": "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"tair_instance_name": "tair_rdb_白名单",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",

						"secondary_zone_id": "${{ref(variable, ZoneId)}}",

						"node_type": "MASTER_SLAVE",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackRedisTairinstance67(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_tairinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisTairinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR-KvstoreDescribeinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistair_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisTairinstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "PayAsYouGo",

					"instance_type": "tair_rdb",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"instance_class": "tair.rdb.1g",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "PayAsYouGo",

						"instance_type": "tair_rdb",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"instance_class": "tair.rdb.1g",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultVSwitch.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultVpc.VpcId)}}",

						"engine_version": "5.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

