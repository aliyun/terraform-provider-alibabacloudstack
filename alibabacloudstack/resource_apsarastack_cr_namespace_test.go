package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_cr_namespace", &resource.Sweeper{
		Name: "alibabacloudstack_cr_namespace",
		F:    testSweepCRNamespace,
	})
}

func testSweepCRNamespace(region string) error {
	// skip not supported region

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("error getting AlibabacloudStack client: %s", err))
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

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
	raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})

	if err != nil {
		log.Printf("[ERROR] %s ", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR))
		return nil
	}
	var resp crListResponse
	bresponse := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	if err != nil {
		log.Printf("[ERROR] %s", errmsgs.WrapError(err))
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
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		crService := CrService{client}
		resp1 := crDescribeNamespaceResponse{}
		raw, err = crService.DescribeCrNamespace(n)
		resp := raw.(*responses.CommonResponse)
		_ = json.Unmarshal(resp.GetHttpContentBytes(), &resp1)
		if resp1.Code != "NAMESPACE_NOT_EXIST" {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}
func testAccCheckNamespaceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_cr_namespace" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		crService := CrService{client}
		log.Printf("namespace ID %s", rs.Primary.ID)
		_, err := crService.DescribeCrNamespace(rs.Primary.ID)

		if err == nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}

	return nil
}
func TestAccAlibabacloudStackCRNamespace_Basic(t *testing.T) {
	var v *crDescribeNamespaceResponse
	resourceId := "alibabacloudstack_cr_namespace.default"
	ra := resourceAttrInit(resourceId, crNamespaceMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
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

func TestAccAlibabacloudStackCRNamespace_Multi(t *testing.T) {
	var v *crDescribeNamespaceResponse
	resourceId := "alibabacloudstack_cr_namespace.default.1"
	ra := resourceAttrInit(resourceId, crNamespaceMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
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

`)
}

var crNamespaceMap = map[string]string{
	"name":               CHECKSET,
	"auto_create":        CHECKSET,
	"default_visibility": CHECKSET,
}
