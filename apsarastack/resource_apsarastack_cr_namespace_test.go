package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strings"
	"testing"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("apsarastack_cr_namespace", &resource.Sweeper{
		Name: "apsarastack_cr_namespace",
		F:    testSweepCRNamespace,
	})
}

func testSweepCRNamespace(region string) error {
	// skip not supported region
	for _, r := range connectivity.CRNoSupportedRegions {
		if region == string(r) {
			log.Printf("[INFO] testSweepCRNamespace skipped not supported region: %s", region)
			return nil
		}
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting ApsaraStack client: %s", err))
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "cr"
	request.Domain = client.Domain
	request.Version = "2016-06-07"
	request.Scheme = "http"
	request.ApiName = "GetNamespaceList"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "cr",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"Action":          "GetNamespaceList",
		"Version":         "2016-06-07",
	}
	raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})

	if err != nil {
		log.Printf("[ERROR] %s ", WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR))
		return nil
	}
	var resp crListResponse
	bresponse := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	if err != nil {
		log.Printf("[ERROR] %s", WrapError(err))
		return nil
	}

	var ns []string
	for _, n := range resp.Data.Namespaces {
		for _, p := range prefixes {
			if strings.HasPrefix(n.Namespace, strings.ToLower(p)) {
				ns = append(ns, n.Namespace)
			}
		}
	}
	log.Printf("namespace ray %v", ns)
	for _, n := range ns {
		request := requests.NewCommonRequest()
		request.Method = "POST"
		request.Product = "cr"
		request.Domain = client.Domain
		request.Version = "2016-06-07"
		request.Scheme = "http"
		request.ApiName = "DeleteNamespace"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "cr",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"RegionId":        client.RegionId,
			"Action":          "DeleteNamespace",
			"Version":         "2016-06-07",
			"Namespace":       n,
		}
		_, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
			return crClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		crService := CrService{client}
		resp1 := crDescribeNamespaceResponse{}
		raw, err = crService.DescribeCrNamespace(n)
		resp := raw.(*responses.CommonResponse)
		_ = json.Unmarshal(resp.GetHttpContentBytes(), &resp1)
		if resp1.Code != "NAMESPACE_NOT_EXIST" {
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cr_namespace", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	return nil
}
func testAccCheckNamespaceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_cr_namespace" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		crService := CrService{client}
		log.Printf("namespace ID %s", rs.Primary.ID)
		_, err := crService.DescribeCrNamespace(rs.Primary.ID)

		if err == nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}
func TestAccApsaraStackCRNamespace_Basic(t *testing.T) {
	var v *crDescribeNamespaceResponse
	resourceId := "apsarastack_cr_namespace.default"
	ra := resourceAttrInit(resourceId, crNamespaceMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRNamespaceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
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
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackCRNamespace_Multi(t *testing.T) {
	var v *crDescribeNamespaceResponse
	resourceId := "apsarastack_cr_namespace.default.1"
	ra := resourceAttrInit(resourceId, crNamespaceMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRNamespaceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name + "${count.index}",
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
					"count":              "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceCRNamespaceConfigDependence(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}
`)
}

var crNamespaceMap = map[string]string{
	"name":               CHECKSET,
	"auto_create":        CHECKSET,
	"default_visibility": CHECKSET,
}
