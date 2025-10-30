package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackLindormInstanceTypesDataSource0(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackLindormInstanceTypesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_lindorm_instance_types.c4g8"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.c4g8", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_lindorm_instance_types.c4g8", "instance_types.0.cpu", "4"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_lindorm_instance_types.c4g8", "instance_types.0.memory", "8"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.c4g8", "instance_types.0.name"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackLindormInstanceTypesDataSource1(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackLindormInstanceTypesDataSourceK8Sc32g64,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_lindorm_instance_types.c32g64"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.c32g64", "instance_types.0.id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.c32g64", "instance_types.0.cpu"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.c32g64", "instance_types.0.memory"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.c32g64", "instance_types.0.name"),
				),
			},
			{
				Config: testAccCheckAlibabacloudStackLindormInstanceTypesDataSourceK8Ssortbycpu,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_lindorm_instance_types.sortbycpu"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.sortbycpu", "instance_types.0.id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.sortbycpu", "instance_types.0.cpu"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.sortbycpu", "instance_types.0.memory"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_lindorm_instance_types.sortbycpu", "instance_types.0.name"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackLindormInstanceTypesDataSourceBasicConfig = `
data "alibabacloudstack_lindorm_instance_types" "c4g8" {
	cpu = 4
	memory = 8
	engine_type = "lindorm"
}
`

const testAccCheckAlibabacloudStackLindormInstanceTypesDataSourceK8Sc32g64 = `
data "alibabacloudstack_lindorm_instance_types" "c32g64" {
	cpu = 32
	memory = 64
}
`
const testAccCheckAlibabacloudStackLindormInstanceTypesDataSourceK8Ssortbycpu = `
data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lindorm"
}
`
