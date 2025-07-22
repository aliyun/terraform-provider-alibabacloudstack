package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_drds_instance", &resource.Sweeper{
		Name: "alibabacloudstack_drds_instance",
		F:    testSweepDrdsInstances,
	})
}

func testSweepDrdsInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

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
		log.Printf("[ERROR] Error retrieving DRDS Instances: %s", errmsgs.WrapError(err))
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
				log.Printf("[ERROR] Error retrieving DRDS Instance: %s. %s", id, errmsgs.WrapError(err))
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

func TestAccAlibabacloudStackDrdsInstance_Vpc(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alibabacloudstack_drds_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testacc-instance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsInstanceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${alibabacloudstack_vpc_vswitch.default.availability_zone}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${alibabacloudstack_vpc_vswitch.default.id}",
					"specification":        "${data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id}",
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
				ImportStateVerify: true,
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

func TestAccAlibabacloudStackDrdsInstance_Classic(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alibabacloudstack_drds_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testacc-instance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsInstanceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   "${var.name}",
					"zone_id":       "${alibabacloudstack_vpc_vswitch.default.availability_zone}",
					"specification": "${data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id}",
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
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   "${var.name}_update",
					"specification": "${data.alibabacloudstack_drds_instance_specifications.default.specifications.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name+"_update",
					}),
				),
			},
		},
	})
}

func resourceDrdsInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s"
	}
	
	data "alibabacloudstack_drds_instance_specifications" "default" {
		sorted_by = "CPU"
	}
	
	%s	
`, name, VSwitchCommonTestCase)
}

var drdsInstancebasicMap = map[string]string{
	"description":          CHECKSET,
	"zone_id":              CHECKSET,
	"instance_charge_type": "PostPaid",
	"specification":        CHECKSET,
}
