package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackImageCopyBasic(t *testing.T) {
	var v ecs.Image

	resourceId := "alibabacloudstack_image_copy.default"
	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alibabacloudstack": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, testAccCopyImageCheckMap)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsCopyImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageCopyBasicConfigDependence)
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckImageDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_image_id":       "${alibabacloudstack_image.default.id}",
					"description":           fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":            name,
					"destination_region_id": region,
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

func TestAccApsaraStackImageCopyEncrypted(t *testing.T) {

	resourceId := "alibabacloudstack_image_copy.default"
	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alibabacloudstack": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, testAccCopyImageCheckMap)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsCopyImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageCopyBasicConfigDependenceEncrypted)
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckImageDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_image_id":       "m-ob6014mhpuyko0jkpvs6",
					"description":           fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"destination_region_id": region,
					"image_name":            name,
					"kms_key_id":            "3852c3cd-3ace-468d-8b9b-c301c33a32b2",
					"encrypted":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
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

			client := provider.Meta().(*connectivity.AlibabacloudStackClient)
			ecsService := EcsService{client}

			resp, err := ecsService.DescribeImageById(rs.Primary.ID)
			if err != nil {
				if errmsgs.NotFoundError(err) {
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

	client := provider.Meta().(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_image_copy" {
			continue
		}

		resp, err := ecsService.DescribeImageById(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
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
data "alibabacloudstack_instance_types" "default" {
    provider = "alibabacloudstack.hz"
 	cpu_core_count    = 1
	memory_size       = 2
}
data "alibabacloudstack_images" "default" {
  provider = "alibabacloudstack.hz"
  name_regex  = "^ubuntu_18.*64"
  owners      = "system"
}
resource "alibabacloudstack_vpc" "default" {
  provider = "alibabacloudstack.hz"
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  provider = "alibabacloudstack.hz"
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_instance_types.default.instance_types.0.availability_zones.0}"
  name              = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
  provider = "alibabacloudstack.hz"
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_instance" "default" {
  provider = "alibabacloudstack.hz"
  image_id = "${data.alibabacloudstack_images.default.ids[0]}"
  instance_type = "${data.alibabacloudstack_instance_types.default.ids[0]}"
  security_groups = "${[alibabacloudstack_security_group.default.id]}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
  instance_name = "${var.name}"
}
resource "alibabacloudstack_image" "default" {
  provider = "alibabacloudstack.hz"
  instance_id = "${alibabacloudstack_instance.default.id}"
  image_name        = "${var.name}"
}
`, name)
}

func resourceImageCopyBasicConfigDependenceEncrypted(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
}`, name)
}
