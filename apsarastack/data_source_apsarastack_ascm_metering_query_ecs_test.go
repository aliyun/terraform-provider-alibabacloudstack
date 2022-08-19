package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscmMetering_queryEcsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceapsarastack_metringqueryecs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_metering_query_ecs.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "data.private_ip_address"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.instance_type_family"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.memory"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.cpu"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.os_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.org_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.instance_network_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.eip_address"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.resource_g_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.instance_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.status"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.sys_disk_size"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.gpu_amount"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.instance_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.vpc_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.start_time"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.end_time"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.create_time"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_metering_query_ecs.default", "policies.data_disk_size"),
				),
			},
		},
	})
}

const datasourceapsarastack_metringqueryecs = `
data "apsarastack_ascm_metering_query_ecs" "default" {
  start_time  =  "2021-01-27T11:00:00Z"
  end_time  = "2021-01-27T12:00:00Z"
  product_name = "ECS"
  is_parent_id = 0
}
`
