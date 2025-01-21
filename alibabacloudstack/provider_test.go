package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	//	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	//
	//	"log"
	"os"
	"testing"
	"time"

	"strconv"
	"strings"
	"path/filepath"
	"io/ioutil"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var defaultRegionToTest = os.Getenv("ALIBABACLOUDSTACK_REGION")

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"alibabacloudstack": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALIBABACLOUDSTACK_ACCESS_KEY"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_SECRET_KEY"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_REGION"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_REGION must be set for acceptance tests")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_INSECURE"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_INSECURE must be set for acceptance tests")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_PROXY"); v == "" {
		//t.Fatal("ALIBABACLOUDSTACK_PROXY must be set for acceptance tests")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_POPGW_DOMAIN must be set for acceptance tests")
	}
	//	if v := os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT"); v == "" {
	//		t.Fatal("ALIBABACLOUDSTACK_DEPARTMENT must be set for acceptance tests")
	//	}
	//	if v := os.Getenv("ALIBABACLOUDSTACK_RESOURCE_GROUP"); v == "" {
	//		t.Fatal("ALIBABACLOUDSTACK_RESOURCE_GROUP must be set for acceptance tests")
	//	}
	//if v := os.Getenv("ALIBABACLOUDSTACK_RESOURCE_GROUP_SET"); v == "" {
	//	t.Fatal("ALIBABACLOUDSTACK_RESOURCE_GROUP_SET must be set for acceptance tests")
	//}

}

func testAccPreCheckWithAccountSiteType(t *testing.T, account AccountSite) {
	defaultAccount := string(DomesticSite)
	if v := strings.TrimSpace(os.Getenv("ALIBABACLOUDSTACK_ACCOUNT_SITE")); v != "" {
		defaultAccount = v
	}
	if defaultAccount != string(account) {
		t.Skipf("Skipping unsupported account type %s-Site. It only supports %s-Site.", defaultAccount, account)
		t.Skipped()
	}
}

func testAccPreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	if v := os.Getenv("ALIBABACLOUDSTACK_ACCESS_KEY"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_SECRET_KEY"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_REGION"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_REGION must be set for acceptance tests")
	}
}

// Skip automatically the sweep testcases which does not support some known regions.
// If supported is true, the regions should a list of supporting the service regions.
// If supported is false, the regions should a list of unsupporting the service regions.
func testSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}

func testAccCheckAlibabacloudStackDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("data source ID not set")
		}
		return nil
	}
}
func testAccPreCheckWithAPIIsNotSupport(t *testing.T) {
	t.Skipf("Skipping because of api is not support, the feature is not output to alibabacloudstack")
	t.Skipped()
}
func testAccPreCheckWithMultipleAccount(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALIBABACLOUDSTACK_ACCESS_KEY_2")); v == "" {
		t.Skipf("Skipping unsupported test with multiple account")
		t.Skipped()
	}
}

func testAccPreCheckOSSForImageImport(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALIBABACLOUDSTACK_OSS_BUCKET_FOR_IMAGE")); v == "" {
		t.Skipf("Skipping tests without OSS_Bucket set.")
		t.Skipped()
	}
	if v := strings.TrimSpace(os.Getenv("ALIBABACLOUDSTACK_OSS_OBJECT_FOR_IMAGE")); v == "" {
		t.Skipf("Skipping OSS_Object does not exist.")
		t.Skipped()
	}
}

func testAccPreCheckWithCmsContactGroupSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALIBABACLOUDSTACK_CMS_CONTACT_GROUP")); v == "" {
		t.Skipf("Skipping the test case with no cms contact group setting")
		t.Skipped()
	}
}

func testAccPreCheckWithSmartAccessGatewaySetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("SAG_INSTANCE_ID")); v == "" {
		t.Skipf("Skipping the test case with no sag instance id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithSmartAccessGatewayAppSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("SAG_APP_INSTANCE_ID")); v == "" {
		t.Skipf("Skipping the test case with no sag app instance id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithTime(t *testing.T) {
	if time.Now().Day() != 1 {
		t.Skipf("Skipping the test case with not the 1st of every month")
		t.Skipped()
	}
}
func testAccPreCheckWithTime2(t *testing.T, days []int) {
	skipped := true
	for _, d := range days {
		if time.Now().Day() == d {
			skipped = false
			break
		}
	}
	if skipped {
		t.Skipf("Skipping the test case when not in specified days %#v of every month", days)
		t.Skipped()
	}
}
func testAccPreCheckWithAlikafkaAclEnable(t *testing.T) {
	aclEnable := os.Getenv("ALIBABACLOUDSTACK_ALIKAFKA_ACL_ENABLE")

	if aclEnable != "true" && aclEnable != "TRUE" {
		t.Skipf("Skipping the test case because the acl is not enabled.")
		t.Skipped()
	}
}

func testAccPreCheckWithNoDefaultVpc(t *testing.T) {
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	request := vpc.CreateDescribeVpcsRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.PageNumber = requests.NewInteger(1)
	request.IsDefault = requests.NewBoolean(true)

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpcs(request)
	})
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	response, _ := raw.(*vpc.DescribeVpcsResponse)

	if len(response.Vpcs.Vpc) < 1 {
		t.Skipf("Skipping the test case with there is no default vpc")
		t.Skipped()
	}
}

func testAccPreCheckWithNoDefaultVswitch(t *testing.T) {
	region := os.Getenv("REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	request := vpc.CreateDescribeVSwitchesRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.PageNumber = requests.NewInteger(1)
	request.IsDefault = requests.NewBoolean(true)

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVSwitches(request)
	})
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	response, _ := raw.(*vpc.DescribeVSwitchesResponse)

	if len(response.VSwitches.VSwitch) < 1 {
		t.Skipf("Skipping the test case with there is no default vswitche")
		t.Skipped()
	}
}

func testAccPreCheckWithEnvVariable(t *testing.T, envVariableName string) {
	if v := strings.TrimSpace(os.Getenv(envVariableName)); v == "" {
		t.Skipf("Skipping the test case with no env variable %s", envVariableName)
		t.Skipped()
	}
}

// func checkoutSupportedRegions(t *testing.T, supported bool, regions []connectivity.Region) {
// 	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
// 	find := false
// 	backupRegion := string(connectivity.APSouthEast1)
// 	if region == string(connectivity.APSouthEast1) {
// 		backupRegion = string(connectivity.EUCentral1)
// 	}

// 	checkoutRegion := os.Getenv("CHECKOUT_REGION")
// 	if checkoutRegion == "true" {
// 		if region == string(connectivity.Hangzhou) {
// 			region = string(connectivity.EUCentral1)
// 			os.Setenv("ALIBABACLOUDSTACK_REGION", region)
// 		}
// 	}
// 	backupRegionFind := false
// 	hangzhouRegionFind := false
// 	for _, r := range regions {
// 		if region == string(r) {
// 			find = true
// 			break
// 		}
// 		if string(r) == backupRegion {
// 			backupRegionFind = true
// 		}
// 		if string(connectivity.Hangzhou) == string(r) {
// 			hangzhouRegionFind = true
// 		}
// 	}

// 	if (find && !supported) || (!find && supported) {
// 		if supported {
// 			if backupRegionFind {
// 				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, backupRegion)
// 				os.Setenv("ALIBABACLOUDSTACK_REGION", backupRegion)
// 				defaultRegionToTest = backupRegion
// 				return
// 			}
// 			if hangzhouRegionFind {
// 				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
// 				os.Setenv("ALIBABACLOUDSTACK_REGION", string(connectivity.Hangzhou))
// 				defaultRegionToTest = string(connectivity.Hangzhou)
// 				return
// 			}
// 			t.Skipf("Skipping unsupported region %s. Supported regions: %s.", region, regions)
// 		} else {
// 			if !backupRegionFind {
// 				t.Logf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
// 				os.Setenv("ALIBABACLOUDSTACK_REGION", backupRegion)
// 				defaultRegionToTest = backupRegion
// 				return
// 			}
// 			if !hangzhouRegionFind {
// 				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
// 				os.Setenv("ALIBABACLOUDSTACK_REGION", string(connectivity.Hangzhou))
// 				defaultRegionToTest = string(connectivity.Hangzhou)
// 				return
// 			}
// 			t.Skipf("Skipping unsupported region %s. Unsupported regions: %s.", region, regions)
// 		}
// 		t.Skipped()
// 	}
// }

var providerCommon = `

`

func getAccTestRandInt(min, max int) int {
	if v := os.Getenv("ALIBABACLOUDSTACK_ACCRANDINT"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return acctest.RandIntRange(min, max)
}

func ResourceTest(t *testing.T, c resource.TestCase) {
	if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_TEST")); err == nil && v {
		dateFolderName := "dryrun_" + time.Now().Format("2006_01_02") 
		err := os.MkdirAll(dateFolderName, 0755)
		if err != nil {
			t.Skipf("Failed to create date folder: %v", err)
		}
		subFolderPath := filepath.Join(dateFolderName, t.Name())
		err = os.MkdirAll(subFolderPath, 0755)
		if err != nil {
			t.Skipf("Failed to create sub folder: %v", err)
		}
		for index, step := range c.Steps{
			filePath := filepath.Join(subFolderPath, fmt.Sprintf("Step%d.tf", index))
			err = ioutil.WriteFile(filePath, []byte(step.Config), 0644)
			if err != nil {
				t.Skipf("Failed to write to file: %v", err)
			}
		}
		t.Skip("Print Terraform .tf file only!")
	} else {
		resource.Test(t, c)
	}
}
