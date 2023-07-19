package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_slb", &resource.Sweeper{
		Name: "alibabacloudstack_slb",
		F:    testSweepSLBs,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alibabacloudstack_cs_cluster",
		},
	})
}

func testSweepSLBs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	k8sPrefix := "kubernetes"

	var slbs []slb.LoadBalancer
	req := slb.CreateDescribeLoadBalancersRequest()
	req.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeLoadBalancers(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving SLBs: %s", err)
		}
		resp, _ := raw.(*slb.DescribeLoadBalancersResponse)
		if resp == nil || len(resp.LoadBalancers.LoadBalancer) < 1 {
			break
		}
		slbs = append(slbs, resp.LoadBalancers.LoadBalancer...)

		if len(resp.LoadBalancers.LoadBalancer) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	service := SlbService{client}
	vpcService := VpcService{client}
	csService := CsService{client}
	for _, loadBalancer := range slbs {
		name := loadBalancer.LoadBalancerName
		id := loadBalancer.LoadBalancerId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(loadBalancer.VpcId, loadBalancer.VSwitchId); err == nil {
				skip = !need
			}

		}
		// If a slb tag key has prefix "kubernetes", this is a slb for k8s cluster and it should be deleted if cluster not exist.
		if skip {
			for _, t := range loadBalancer.Tags.Tag {
				if strings.HasPrefix(strings.ToLower(t.TagKey), strings.ToLower(k8sPrefix)) {
					_, err := csService.DescribeCsKubernetes(name)
					if NotFoundError(err) {
						skip = false
					} else {
						skip = true
						break
					}
				}
			}
		}
		if skip {
			log.Printf("[INFO] Skipping SLB: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting SLB: %s (%s)", name, id)
		if err := service.sweepSlb(id); err != nil {
			log.Printf("[ERROR] Failed to delete SLB (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackSlb_classictest(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alibabacloudstack_slb.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-test%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbClassicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.SlbClassicNoSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          name,
					"address_type":  "internet",
					"specification": "slb.s2.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          name,
						"address_type":  "internet",
						"specification": "slb.s2.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"specification": "slb.s1.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "slb.s1.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         name,
					"address_type": "internet",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"address_type": "internet",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
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
					"name": fmt.Sprintf("tf-testAccSlbClassicInstanceConfigSpot%d_change", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccSlbClassicInstanceConfigSpot%d_change", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         name,
					"address_type": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"address_type": "internet",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSlb_vpctest(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alibabacloudstack_slb.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-slbtest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbVpcConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       name,
					"vswitch_id": "${alibabacloudstack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("tf-testAccSlbVpcInstanceConfigSpot%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccSlbVpcInstanceConfigSpot%d", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSlb_vpcmulti(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alibabacloudstack_slb.default.2"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testaccslbvpcinstancemulticonfigspot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbVpcConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":        "3",
					"name":         name,
					"address_type": "intranet",
					"vswitch_id":   "${alibabacloudstack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
		},
	})
}

func resourceSlbVpcConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
provider "alibabacloudstack" {
	assume_role {}
}
	variable "name" {
		default = "%s"
	}
	`, SlbVpcCommonTestCase, name)
}

func resourceSlbClassicConfigDependence(name string) string {
	return fmt.Sprintf(`
provider "alibabacloudstack" {
	assume_role {}
}
	variable "name" {
		default = "%s"
	}
	`, name)
}
