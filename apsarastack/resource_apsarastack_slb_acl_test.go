package apsarastack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("apsarastack_slb_acl", &resource.Sweeper{
		Name: "apsarastack_slb_acl",
		F:    testSweepSlbAcl,
	})
}
func testSweepSlbAcl(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	req := slb.CreateDescribeAccessControlListsRequest()
	req.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAccessControlLists(req)
	})
	if err != nil {
		return err
	}
	resp, _ := raw.(*slb.DescribeAccessControlListsResponse)

	for _, acl := range resp.Acls.Acl {
		name := acl.AclName
		id := acl.AclId

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Slb Acl: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting Slb Acl : %s (%s)", name, id)
		req := slb.CreateDeleteAccessControlListRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.AclId = id
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteAccessControlList(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Slb Acl (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccApsaraStackSlbAcl_basic(t *testing.T) {
	var acl *slb.DescribeAccessControlListAttributeResponse

	resourceId := "apsarastack_slb_acl.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &acl, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbAcl")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSLBAclDependence)
	resource.Test(t, resource.TestCase{
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
					"name":       "${var.name}",
					"ip_version": "ipv4",
					"entry_list": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": "80",
						},
						{
							"entry":   "168.10.10.0/24",
							"comment": "second",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl",
						"ip_version":   "ipv4",
						"entry_list.#": "2",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
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
					"name": "tf-testAccSlbAcl-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbAcl-name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entry_list": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": "80",
						},
						{
							"entry":   "168.10.10.0/24",
							"comment": "second",
						},
						{
							"entry":   "172.10.10.0/24",
							"comment": "third",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entry_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "acceptance test1231",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF1",
						"tags.For":     "acceptance test1231",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"entry_list": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": "80",
						},
						{
							"entry":   "168.10.10.0/24",
							"comment": "second",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl",
						"ip_version":   "ipv4",
						"entry_list.#": "2",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackSlbAcl_muilt(t *testing.T) {
	var acl *slb.DescribeAccessControlListAttributeResponse

	resourceId := "apsarastack_slb_acl.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &acl, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbAcl")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSLBAclDependence)
	resource.Test(t, resource.TestCase{
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
					"name":       "${var.name}-${count.index}",
					"count":      "10",
					"ip_version": "ipv4",
					"entry_list": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": "80",
						},
						{
							"entry":   "168.10.10.0/24",
							"comment": "second",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl-9",
						"ip_version":   "ipv4",
						"entry_list.#": "2",
					}),
				),
			},
		},
	})
}

func resourceSLBAclDependence(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}
variable "name" {
  default = "%s"
}
`, name)
}
