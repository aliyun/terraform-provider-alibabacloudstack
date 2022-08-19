package apsarastack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApsaraStackZonesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackZonesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.foo"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_disk_categories.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccApsaraStackZonesDataSource_filter(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackZonesDataSourceFilter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.foo"),
					testCheckZoneLength("data.apsarastack_zones.foo"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_disk_categories.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
				),
			},

			{
				Config: testAccCheckApsaraStackZonesDataSourceFilterIoOptimized,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.foo"),
					testCheckZoneLength("data.apsarastack_zones.foo"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_disk_categories.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccApsaraStackZonesDataSource_unitRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackZonesDataSourceUnitRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.foo"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.available_disk_categories.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccApsaraStackZonesDataSource_multiZone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackZonesDataSourceMultiZone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.default"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.id"),
					//resource.TestMatchResourceAttr("data.apsarastack_zones.default", "zones.0.id", regexp.MustCompile(fmt.Sprintf(".%s.", MULTI_IZ_SYMBOL))),
					//resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.local_name", "a"),
					//resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.available_instance_types.#", "0"),
					//resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.available_resource_creation.#", "0"),
					//resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.available_disk_categories.#", "0"),

					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_disk_categories.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.multi_zone_ids.#"),
					//resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.multi_zone_ids.0"),
					//resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.multi_zone_ids.1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "ids.#"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.slb_slave_zone_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccApsaraStackZonesDataSource_chargeType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackZonesDataSourceChargeType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.default"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_disk_categories.#"),
					//resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.local_name", ""),
					//resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.available_instance_types.#", "0"),
					//resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.available_resource_creation.#", "0"),
					//resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.available_disk_categories.#", "0"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "ids.#"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.0.slb_slave_zone_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccApsaraStackZonesDataSource_slb(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackZonesDataSource_slb,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.default"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.available_disk_categories.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "ids.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "zones.0.slb_slave_zone_ids.#"),
				),
			},
		},
	})
}

func TestAccApsaraStackZonesDataSource_enable_details(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackZonesDataSourceEnableDetails,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.foo"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.local_name", ""),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.available_instance_types.#", "0"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.available_resource_creation.#", "0"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.available_disk_categories.#", "0"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
				),
			},
		},
	})
}
func TestAccApsaraStackZonesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackZonesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_zones.default"),
					resource.TestCheckResourceAttr("data.apsarastack_zones.default", "zones.#", "0"),
					resource.TestCheckNoResourceAttr("data.apsarastack_zones.default", "zones.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_zones.default", "zones.local_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_zones.default", "zones.available_instance_types"),
					resource.TestCheckNoResourceAttr("data.apsarastack_zones.default", "zones.available_resource_creation"),
					resource.TestCheckNoResourceAttr("data.apsarastack_zones.default", "zones.available_disk_categories"),
					resource.TestCheckResourceAttrSet("data.apsarastack_zones.default", "ids.#"),
				),
			},
		},
	})
}

// the zone length changed occasionally
// check by range to avoid test case failure
func testCheckZoneLength(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s", name)
		}

		i, err := strconv.Atoi(is.Attributes["zones.#"])

		if err != nil {
			return fmt.Errorf("convert zone length err: %#v", err)
		}

		if i <= 0 {
			return fmt.Errorf("zone length expected greater than 0 got err: %d", i)
		}

		return nil
	}
}

const testAccCheckApsaraStackZonesDataSourceBasicConfig = `
data "apsarastack_zones" "foo" {
	enable_details = true
}
`

const testAccCheckApsaraStackZonesDataSourceFilter = `
data "apsarastack_zones" "foo" {
	available_resource_creation= "VSwitch"
	available_disk_category= "cloud_efficiency"
	enable_details = true
}
`

const testAccCheckApsaraStackZonesDataSourceFilterIoOptimized = `
data "apsarastack_zones" "foo" {
	available_resource_creation= "IoOptimized"
	available_disk_category= "cloud_efficiency"
	enable_details = true
}
`

const testAccCheckApsaraStackZonesDataSourceUnitRegion = `
data "apsarastack_zones" "foo" {
	available_resource_creation= "VSwitch"
	enable_details = true
}
`

const testAccCheckApsaraStackZonesDataSourceMultiZone = `
data "apsarastack_zones" "default" {
  available_resource_creation= "Rds"
  multi = true
  enable_details = true
}`

const testAccCheckApsaraStackZonesDataSourceChargeType = `
data "apsarastack_zones" "default" {
  instance_charge_type = "PrePaid"
  available_resource_creation= "Rds"
  multi = true
  enable_details = true
}`

const testAccCheckApsaraStackZonesDataSource_slb = `
data "apsarastack_zones" "default" {
  available_resource_creation= "Slb"
  enable_details = true
  available_slb_address_ip_version= "ipv4"
  available_slb_address_type="Vpc"
}`

const testAccCheckApsaraStackZonesDataSourceEnableDetails = `
data "apsarastack_zones" "foo" {}
`
const testAccCheckApsaraStackZonesDataSourceEmpty = `
data "apsarastack_zones" "default" {
  available_instance_type = "ecs.n1.fake"
}
`
