package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenRouteMaps0(t *testing.T) {
	var v *CbnDescribeCenRouteMapsResponse

	resourceId := "alibabacloudstack_cen_route_map.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenRouteMapsCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribeCenRouteMapsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sroute_map%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%sroute_map%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenRouteMapsBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"cen_id":                        "${alibabacloudstack_cen_instance.default.cen_id}",
					"transit_router_route_table_id": "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}",
					"priority":                      "3",
					"transmit_direction":            "RegionIn",
					"map_result":                    "Deny",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cen_id": CHECKSET,

						"transit_router_route_table_id": CHECKSET,
						"priority":                      "3",
						"transmit_direction":            "RegionIn",
						"map_result":                    "Deny",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"source_instance_ids": []string{"aaaaa"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"source_instance_ids.#": "1",
						"source_instance_ids.0": "aaaaa",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"source_instance_ids_reverse_match": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"source_instance_ids_reverse_match": "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"source_instance_ids":                    REMOVEKEY,
					"source_instance_ids_reverse_match":      REMOVEKEY,
					"destination_instance_ids_reverse_match": "true",
					"destination_instance_ids":               []string{"bbbbb"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"destination_instance_ids_reverse_match": "true",
						"destination_instance_ids.#":             "1",
						"destination_instance_ids.0":             "bbbbb",
						"source_instance_ids.#":                 "0",
						"source_instance_ids.0": REMOVEKEY,
						"source_instance_ids_reverse_match": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"destination_instance_ids_reverse_match": REMOVEKEY,
					"destination_instance_ids":               REMOVEKEY,
					"destination_route_table_ids":            []string{"cccccc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"destination_route_table_ids.#": "1",
						"destination_route_table_ids.0": "cccccc",
						"destination_instance_ids_reverse_match": REMOVEKEY,
						"destination_instance_ids.#":      "0",
						"destination_instance_ids.0":      REMOVEKEY,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"destination_route_table_ids": REMOVEKEY,
					"source_child_instance_types": []string{"VBR", "VPC"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"source_child_instance_types.#": "2",
						"source_child_instance_types.0": "VBR",
						"source_child_instance_types.1": "VPC",
						"destination_route_table_ids.#": "0",
						"destination_route_table_ids.0": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"source_child_instance_types":      REMOVEKEY,
					"destination_child_instance_types": []string{"VBR", "VPC"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"destination_child_instance_types.#": "2",
						"destination_child_instance_types.0": "VBR",
						"destination_child_instance_types.1": "VPC",
						"source_child_instance_types.#":      "0",
						"source_child_instance_types.0": REMOVEKEY,
						"source_child_instance_types.1": REMOVEKEY,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"destination_child_instance_types": REMOVEKEY,
					"match_address_type":               "IPv6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"match_address_type": "IPv6",
						"destination_child_instance_types.#": "0",
						"destination_child_instance_types.0": REMOVEKEY,
						"destination_child_instance_types.1": REMOVEKEY,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"match_address_type":      REMOVEKEY,
					"cidr_match_mode":         "Include",
					"destination_cidr_blocks": []string{"11.11.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cidr_match_mode":           "Include",
						"destination_cidr_blocks.#": "1",
						"destination_cidr_blocks.0": "11.11.0.0/16",
						"match_address_type": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"cidr_match_mode": "Complete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cidr_match_mode": "Complete",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"cidr_match_mode":         REMOVEKEY,
					"destination_cidr_blocks": REMOVEKEY,
					"route_types":             []string{"Custom", "System"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"route_types.#": "2",
						"route_types.0": "Custom",
						"route_types.1": "System",
						"cidr_match_mode": REMOVEKEY,
						"destination_cidr_blocks.#": "0",
						"destination_cidr_blocks.0": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"route_types":        REMOVEKEY,
					"match_asns":         []string{"16100", "17100"},
					"as_path_match_mode": "Include",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"match_asns.#":       "2",
						"match_asns.0":       "16100",
						"match_asns.1":       "17100",
						"as_path_match_mode": "Include",
						"route_types.#": "0",
						"route_types.0": REMOVEKEY,
						"route_types.1": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"as_path_match_mode": "Complete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"as_path_match_mode": "Complete",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"as_path_match_mode":   REMOVEKEY,
					"match_asns":           REMOVEKEY,
					"match_community_set":  []string{"16100:111", "17100:111"},
					"community_match_mode": "Complete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"match_community_set.#": "2",
						"match_community_set.0": "16100:111",
						"match_community_set.1": "17100:111",
						"community_match_mode":  "Complete",
						"as_path_match_mode":   REMOVEKEY,
						"match_asns.#": "0",
						"match_asns.0":REMOVEKEY,
						"match_asns.1":REMOVEKEY,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"community_match_mode": "Include",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"community_match_mode": "Include",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description":   modify_name,
					"priority":      "10",
					"map_result":    "Permit",
					"preference":    "1",
					"next_priority": "11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description":   modify_name,
						"priority":      "10",
						"map_result":    "Permit",
						"preference":    "1",
						"next_priority": "11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preference":            REMOVEKEY,
					"operate_community_set": []string{"16100:111", "17100:111"},
					"community_operate_mode":  "Additive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operate_community_set.#": "2",
						"operate_community_set.0": "16100:111",
						"operate_community_set.1": "17100:111",
						"community_operate_mode":    "Additive",
						"preference":            REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"community_operate_mode": "Replace",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"community_operate_mode": "Replace",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"match_address_type", "next_priority", "source_route_table_ids", "source_instance_ids_reverse_match", "destination_instance_ids_reverse_match", "community_match_mode", "community_operate_mode", "cidr_match_mode", "map_result", "preference", "as_path_match_mode", "destination_child_instance_types", "source_region_ids", "source_instance_ids", "operate_community_set", "destination_instance_ids", "prepend_as_path", "destination_route_table_ids", "source_child_instance_types", "destination_cidr_blocks", "route_types", "match_asns", "match_community_set"},
			},
		},
	})
}

var AlibabacloudTestAccCenRouteMapsCheckmap = map[string]string{

	"cen_id": CHECKSET,

	"transit_router_route_table_id": CHECKSET,

	"priority": CHECKSET,

	"transmit_direction": CHECKSET,

	"map_result": CHECKSET,
}

func AlibabacloudTestAccCenRouteMapsBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_instance.default.transit_router_id}"
}

`, name)
}
