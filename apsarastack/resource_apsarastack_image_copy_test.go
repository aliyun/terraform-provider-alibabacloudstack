package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackImageCopyBasic(t *testing.T) {
	var v ecs.Image

	resourceId := "apsarastack_image_copy.default"
	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"apsarastack": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, testAccCopyImageCheckMap)
	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsCopyImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageCopyBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckImageDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":         "apsarastack.sh",
					"source_image_id":  "${apsarastack_image.default.id}",
					"source_region_id": "cn-hangzhou",
					"description":      fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"image_name":  name,
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"image_name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"image_name":  name,
					}),
				),
			},
		},
	})
}

func testAccCheckImageExistsWithProviders(n string, image *ecs.Image, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No image  ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.ApsaraStackClient)
			ecsService := EcsService{client}

			resp, err := ecsService.DescribeImageById(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}

			*image = resp
			return nil
		}
		return fmt.Errorf("image not found")
	}
}

func testAccCheckImageDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckImageDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckImageDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {

	client := provider.Meta().(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_copy_image" {
			continue
		}

		resp, err := ecsService.DescribeImageById(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("image still exist,  ID %s ", resp.ImageId)
		}
	}

	return nil
}

var testAccCopyImageCheckMap = map[string]string{}

func resourceImageCopyBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
provider "apsarastack" {
  alias = "sh"
  region = "cn-shanghai"
}
provider "apsarastack" {
  alias = "hz"
  region = "cn-hangzhou"
}
data "apsarastack_instance_types" "default" {
    provider = "apsarastack.hz"
 	cpu_core_count    = 1
	memory_size       = 2
}
data "apsarastack_images" "default" {
  provider = "apsarastack.hz"
  name_regex  = "^ubuntu_18.*64"
  owners      = "system"
}
resource "apsarastack_vpc" "default" {
  provider = "apsarastack.hz"
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  provider = "apsarastack.hz"
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  provider = "apsarastack.hz"
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "default" {
  provider = "apsarastack.hz"
  image_id = "${data.apsarastack_images.default.ids[0]}"
  instance_type = "${data.apsarastack_instance_types.default.ids[0]}"
  security_groups = "${[apsarastack_security_group.default.id]}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
  instance_name = "${var.name}"
}
resource "apsarastack_image" "default" {
  provider = "apsarastack.hz"
  instance_id = "${apsarastack_instance.default.id}"
  image_name        = "${var.name}"
}
`, name)
}
