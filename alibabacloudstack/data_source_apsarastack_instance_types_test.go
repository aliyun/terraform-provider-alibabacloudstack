package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackInstanceTypesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackInstanceTypesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_instance_types.c4g8"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.burstable_instance.0.%", "2"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.local_storage.0.%", "3"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c4g8", "instance_types.0.local_storage.0.category", ""),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c4g8", "ids.#"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackInstanceTypesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackInstanceTypesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_instance_types.empty"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.price"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.empty", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.empty", "ids.#"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackInstanceTypesDataSource_k8sSpec(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackInstanceTypesDataSourceK8Sc1g2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_instance_types.c1g2"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.c1g2", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c1g2", "ids.#"),
				),
			},
			{
				Config: testAccCheckAlibabacloudStackInstanceTypesDataSourceK8Sc2g4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_instance_types.c2g4"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.cpu_core_count", "2"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.memory_size", "4"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.burstable_instance.0.%", "2"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.local_storage.0.%", "3"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.c2g4", "instance_types.0.local_storage.0.category", ""),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.c2g4", "ids.#"),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackInstanceTypesDataSource_k8sFamily(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackInstanceTypesDataSourceK8ST5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_instance_types.t5"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_instance_types.t5", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instance_types.t5", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackInstanceTypesDataSourceBasicConfig = `
data "alibabacloudstack_instance_types" "c4g8" {
	cpu_core_count = 4
	memory_size = 8
}
`

const testAccCheckAlibabacloudStackInstanceTypesDataSourceEmpty = `
data "alibabacloudstack_instance_types" "empty" {
	instance_type_family = "ecs.fake"
}
`

const testAccCheckAlibabacloudStackInstanceTypesDataSourceK8Sc1g2 = `
data "alibabacloudstack_instance_types" "c1g2" {
	cpu_core_count = 1
	memory_size = 2
	kubernetes_node_role = "Master"
}
`
const testAccCheckAlibabacloudStackInstanceTypesDataSourceK8Sc2g4 = `
data "alibabacloudstack_instance_types" "c2g4" {
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}
`
const testAccCheckAlibabacloudStackInstanceTypesDataSourceK8ST5 = `
data "alibabacloudstack_instance_types" "t5" {
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Master"
	instance_type_family = "ecs.t5"
}
`
