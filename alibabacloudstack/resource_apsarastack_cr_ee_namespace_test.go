package alibabacloudstack

/*import (
    "github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_cr_ee_namespace", &resource.Sweeper{
		Name: "alibabacloudstack_cr_ee_namespace",
		F:    testSweepCrEENamespace,
	})
}

var testaccCrEEInstanceId string

func setTestaccCrEEInstanceId(t *testing.T) {
	if testaccCrEEInstanceId != "" {
		return
	}

	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	resp, err := crService.ListCrEEInstances(1, 10)
	if err != nil {
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	if len(resp.Instances) == 0 {
		t.Skipf("Skipping cr ee test case without default instances")
	}
	testaccCrEEInstanceId = resp.Instances[0].InstanceId
}

func testSweepCrEENamespace(region string) error {
	if testaccCrEEInstanceId == "" {
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("error getting AlibabacloudStack client: %s", err))
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}

	pageNo := 1
	pageSize := 50
	var namespaces []cr_ee.NamespacesItem
	for {
		resp, err := crService.ListCrEENamespaces(testaccCrEEInstanceId, pageNo, pageSize)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		namespaces = append(namespaces, resp.Namespaces...)
		if len(resp.Namespaces) < pageSize {
			break
		}
		pageNo++
	}

	testPrefix := "tf-testacc"
	for _, namespace := range namespaces {
		if strings.HasPrefix(namespace.NamespaceName, testPrefix) {
			crService.DeleteCrEENamespace(namespace.InstanceId, namespace.NamespaceName)
		}
	}
	return nil
}

func TestAccAlibabacloudStackCrEENamespace_Basic(t *testing.T) {
	setTestaccCrEEInstanceId(t)
	var v *cr_ee.GetNamespaceResponse
	resourceId := "alibabacloudstack_cr_ee_namespace.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEENamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEENamespaceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithCrEE(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        testaccCrEEInstanceId,
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        testaccCrEEInstanceId,
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
					"auto_create": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_create": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_visibility": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_visibility": "PRIVATE",
					}),
				),
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

func TestAccAlibabacloudStackCrEENamespace_Multi(t *testing.T) {
	setTestaccCrEEInstanceId(t)
	var v *cr_ee.GetNamespaceResponse
	resourceId := "alibabacloudstack_cr_ee_namespace.default.4"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEENamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEENamespaceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithCrEE(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        testaccCrEEInstanceId,
					"name":               name + "${count.index}",
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
					"count":              "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name + fmt.Sprint(4),
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
					}),
				),
			},
		},
	})
}

func resourceCrEENamespaceConfigDependence(name string) string {
	return ""
}

func testAccPreCheckWithCrEE(t *testing.T) {
	testAccPreCheck(t)
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	resp, err := crService.ListCrEEInstances(1, 10)
	if err != nil {
		//Maybe crEE has not opened int the region
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	if len(resp.Instances) == 0 {
		t.Skipf("Skipping cr ee test case without default instances")
	}
}
*/
