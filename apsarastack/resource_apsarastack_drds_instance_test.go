package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
)

func init() {
	resource.AddTestSweepers("apsarastack_drds_instance", &resource.Sweeper{
		Name: "apsarastack_drds_instance",
		F:    testSweepDRDSInstances,
	})
}

func testSweepDRDSInstances(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DrdsSupportedRegions) {
		log.Printf("[INFO] Skipping DRDS Instance unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := drds.CreateDescribeDrdsInstancesRequest()
	request.Headers["x-ascm-product-name"] = "Drds"
	request.Headers["x-acs-organizationId"] = client.Department
	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving DRDS Instances: %s", WrapError(err))
	}
	response, _ := raw.(*drds.DescribeDrdsInstancesResponse)

	vpcService := VpcService{client}
	for _, v := range response.Instances.Instance {
		name := v.Description
		id := v.DrdsInstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			instanceDetailRequest := drds.CreateDescribeDrdsInstanceRequest()
			instanceDetailRequest.DrdsInstanceId = id
			instanceDetailRequest.Headers["x-ascm-product-name"] = "Drds"
			instanceDetailRequest.Headers["x-acs-organizationId"] = client.Department
			raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
				return drdsClient.DescribeDrdsInstance(instanceDetailRequest)
			})
			if err != nil {
				log.Printf("[ERROR] Error retrieving DRDS Instance: %s. %s", id, WrapError(err))
			}
			instanceDetailResponse, _ := raw.(*drds.DescribeDrdsInstanceResponse)
			for _, vip := range instanceDetailResponse.Data.Vips.Vip {
				if need, err := vpcService.needSweepVpc(vip.VpcId, ""); err == nil {
					skip = !need
					break
				}
			}

		}
		if skip {
			log.Printf("[INFO] Skipping DRDS Instance: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting DRDS Instance: %s (%s)", name, id)
		req := drds.CreateRemoveDrdsInstanceRequest()
		req.DrdsInstanceId = id
		req.Headers["x-ascm-product-name"] = "Drds"
		req.Headers["x-acs-organizationId"] = client.Department
		_, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.RemoveDrdsInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete DRDS Instance (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccApsaraStackDRDSInstance_Vpc(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "apsarastack_drds_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${apsarastack_vswitch.default.availability_zone}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${apsarastack_vswitch.default.id}",
					"specification":        "drds.sn2.4c16g.8C32G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackDRDSInstance_Multi(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "apsarastack_drds_instance.default.2"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithRegions(t, false, connectivity.DrdsClassicNoSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${apsarastack_vswitch.default.availability_zone}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${apsarastack_vswitch.default.id}",
					"specification":        "drds.sn2.4c16g.8C32G",
					"count":                "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
		},
	})
}

func resourceDRDSInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "apsarastack_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	
	variable "instance_series" {
		default = "drds.sn2.4c16g"
	}
	
	resource "apsarastack_vpc" "default" {
	  name       = "${var.name}"
	  cidr_block = "172.16.0.0/16"
	}
	resource "apsarastack_vswitch" "default" {
	  vpc_id            = "${apsarastack_vpc.default.id}"
	  cidr_block        = "172.16.0.0/24"
	  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	  name              = "${var.name}"
	}	
`, name)
}

var drdsInstancebasicMap = map[string]string{
	"description":          CHECKSET,
	"zone_id":              CHECKSET,
	"instance_series":      "drds.sn2.4c16g",
	"instance_charge_type": "PostPaid",
	"specification":        "drds.sn2.4c16g.8C32G",
}
