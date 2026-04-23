package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"io/ioutil"
	"path/filepath"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	schemaHelper "github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testAccProviders map[string]*schema.Provider
var testAccExternalProviders map[string]resource.ExternalProvider
var testAccProvider, testYundunProvider *schema.Provider
var defaultRegionToTest = os.Getenv("ALIBABACLOUDSTACK_REGION")

func init() {
	rand.Seed(time.Now().UnixNano())
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"alibabacloudstack": testAccProvider,
	}
	testAccExternalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source: "hashicorp/random",
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testYunDunProviders() map[string]*schema.Provider {
	testYundunProvider = Provider()
	testYundunProvider.Schema["access_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_ACCESS_KEY", ""),
		Description: descriptions["access_key"],
	}
	testYundunProvider.Schema["secret_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_SECRET_KEY", ""),
		Description: descriptions["secret_key"],
	}
	testYundunProvider.Schema["role_arn"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: descriptions["assume_role_role_arn"],
		DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_ASSUME_ROLE_ARN", ""),
	}
	return map[string]*schema.Provider{
		"alibabacloudstack":        testYundunProvider,
		"alibabacloudstack-common": testAccProvider,
	}
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

func testAccPreCheckOssEndpointList(t *testing.T) {
	testAccPreCheck(t)
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping OSS test case with err: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	if _, err := ossService.GetOssEndpointList(); err != nil {
		if errmsgs.IsHostNotFound(err) {
			t.Skipf("Skipping OSS cluster test case: GetOssEndpointList API not support")
		}
		t.Fatalf("GetOssEndpointList failed: %s", err)
	}
}

func testAccPreYunCheck(t *testing.T) {
	if v := os.Getenv("ALIBABACLOUDSTACK_YUNDUN_ACCESS_KEY"); v == "" {
		t.Skipf("ALIBABACLOUDSTACK_YUNDUN_ACCESS_KEY must be set for acceptance tests")
		t.Skipped()
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_YUNDUN_SECRET_KEY"); v == "" {
		t.Skipf("ALIBABACLOUDSTACK_YUNDUN_SECRET_KEY must be set for acceptance tests")
		t.Skipped()
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_YUNDUN_ASSUME_ROLE_ARN"); v == "" {
		t.Skipf("ALIBABACLOUDSTACK_YUNDUN_ASSUME_ROLE_ARN must be set for acceptance tests")
		t.Skipped()
	}
}

func testAccApigwV2ServicePreCheck(t *testing.T) {
	if v := os.Getenv("ALIBABACLOUDSTACK_EDAS_ACCESS_KEY"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_EDAS_ACCESS_KEY must be set for create service with service source test")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_EDAS_SECRET_KEY"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_EDAS_ACCESS_KEY must be set for create service with service source test")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_EDAS_ENDPOINT_PORT"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_EDAS_ENDPOINT_PORT must be set for create service with service source test")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_EDAS_ENDPOINT"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_EDAS_ENDPOINT must be set for create service with service source test")
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_EDAS_NAMESPACE_ID"); v == "" {
		t.Fatal("ALIBABACLOUDSTACK_EDAS_NAMESPACE_ID must be set for create service with service source test")
	}
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

func testAccPreCheckWithTime2(t *testing.T, days []int) {
	if !slices.Contains(days, time.Now().Day()) {
		t.Skipf("Skipping the test case when not in specified days %#v of every month", days)
		t.Skipped()
	}
}

func testAccPreCheckWithMultiAZ(t *testing.T) {
	req := ecs.CreateDescribeZonesRequest()
	rawClient, err := sharedClientForRegion(os.Getenv("ALIBABACLOUDSTACK_REGION"))
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	client.InitRpcRequest(*req.RpcRequest)

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(req)
	})
	resp, ok := raw.(*ecs.DescribeZonesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_zones", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		t.Fatalf("%v", err)
		t.Failed()
	}
	addDebug(req.GetActionName(), raw, req.RpcRequest, req)
	if resp == nil {
		t.Fatalf("There are no availability zones in the region: %#v.", client.Region)
		t.Failed()
	}

	if len(resp.Zones.Zone) < 3 {
		t.Skipf("Insufficient number of zones")
		t.Skipped()
	}
}

func testAccPreCheckKmsServer(t *testing.T) {
	req := ecs.CreateDescribeZonesRequest()
	rawClient, err := sharedClientForRegion(os.Getenv("ALIBABACLOUDSTACK_REGION"))
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	client.InitRpcRequest(*req.RpcRequest)
	request := kms.CreateListKeysRequest()
	client.InitRpcRequest(*request.RpcRequest)
	_, err = client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.ListKeys(request)
	})
	if err != nil {
		t.Skipf("Skipping the test case with Kms Server err: %s", err)
	}

}

func testAccPreCheckWithEnvVariable(t *testing.T, envVariableName string) {
	if v := strings.TrimSpace(os.Getenv(envVariableName)); v == "" {
		t.Skipf("Skipping the test case with no env variable %s", envVariableName)
		t.Skipped()
	}
}

func testAccPreCheckWithCrEe(t *testing.T) {
	testAccPreCheck(t)
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	resp, err := crService.ListCrEeInstances(1, 10)
	if err != nil {
		//Maybe crEE has not opened int the region
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	instances := resp["Instances"].([]interface{})
	if len(instances) == 0 {
		t.Skipf("Skipping cr ee test case without default instances")
	}
}

func testAccPreCheckPhysicalConnection(t *testing.T) {
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	action := "DescribePhysicalConnections"
	request := map[string]interface{}{
		"PageSize":   100,
		"PageNumber": 1,
	}
	response, err := client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		t.Skipf("Skipping test case for PhysicalConnection because an error: %s", err)
	}
	resp, err := jsonpath.Get("$.PhysicalConnectionSet.PhysicalConnectionType", response)
	if err != nil {
		t.Skipf("Skipping test case for PhysicalConnection because an error: %s", err)
	}
	result, _ := resp.([]interface{})
	if len(result) == 0 {
		t.Skipf("Skipping test case for PhysicalConnection because Lacks example instances.")
	}
}

func getAccTestRandInt(min, max int) int {
	if v := os.Getenv("ALIBABACLOUDSTACK_ACCRANDINT"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return acctest.RandIntRange(min, max)
}

// shuffle function: randomly shuffle character order
func shuffle(chars []rune) []rune {
	copyChars := make([]rune, len(chars))
	copy(copyChars, chars)
	for i := range copyChars {
		j := rand.Intn(i + 1)
		copyChars[i], copyChars[j] = copyChars[j], copyChars[i]
	}
	return copyChars
}

func getAccTestPassword(length int) string {
	// Deprecated method, please use RandomPasswordTestCase instead
	if v := os.Getenv("ALIBABACLOUDSTACK_ACCRANDPWD"); v != "" {
		return v
	}
	if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_TEST")); err != nil && v {
		if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_SENSITIVE")); err == nil && v {
			return "<YOUR PASSWORD>"
		}
	}

	// Define character sets
	const (
		upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerLetters = "abcdefghijklmnopqrstuvwxyz"
		digits       = "0123456789"
		symbols      = "!@#$^&*()_"
	)

	// Convert to rune arrays
	var (
		upperRunes  = []rune(upperLetters)
		lowerRunes  = []rune(lowerLetters)
		digitRunes  = []rune(digits)
		symbolRunes = []rune(symbols)
		allRunes    = append(upperRunes, append(lowerRunes, append(digitRunes, symbolRunes...)...)...)
	)

	// First character: force uppercase letter
	firstChar := upperRunes[rand.Intn(len(upperRunes))]

	// Required characters (digits, lowercase letters, symbols)
	mustHave := []rune{
		digitRunes[rand.Intn(len(digitRunes))],   // digit
		lowerRunes[rand.Intn(len(lowerRunes))],   // lowercase letter
		symbolRunes[rand.Intn(len(symbolRunes))], // symbol
	}

	// Remaining characters (optional all types)
	remaining := length - 1 - 3 //  - 1(first char) - 3(required chars)
	others := make([]rune, remaining)
	for i := range others {
		others[i] = allRunes[rand.Intn(len(allRunes))]
	}

	// Merge required characters and random characters
	allChars := append(mustHave, others...)
	shuffled := shuffle(allChars)

	// Combine final password
	password := append([]rune{firstChar}, shuffled...)
	return string(password)
}

func getAccTestOsEnv(keyName string) string {
	if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_TEST")); err != nil && v {
		if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_SENSITIVE")); err == nil && v {
			return fmt.Sprintf("<%s>", keyName)
		}
	}
	return os.Getenv(keyName)
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
		for index, step := range c.Steps {
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

func TestProviderSchema(t *testing.T) {
	outputDir := filepath.Join("dryrun_provider")
	outputPath := filepath.Join(outputDir, "TestProviderSchema.json")

	// Prepare test environment
	require.NoError(t, os.MkdirAll(outputDir, 0755), "Failed to create directory")

	// Generate schema
	providerSchema := schemaHelper.ConvertAndWrapProviderSchema(DefaultProviderName, testAccProvider)

	data, err := json.MarshalIndent(providerSchema, "", "  ")
	require.NoError(t, err, "JSON formatting failed")
	require.NoError(t, os.WriteFile(outputPath, data, 0644), "Failed to write file")

	// Validate output
	fileInfo, err := os.Stat(outputPath)
	require.NoError(t, err, "File not found")
	require.True(t, fileInfo.Mode().IsRegular(), "Not a regular file")
	require.Greater(t, fileInfo.Size(), int64(100), "File content too short")

	jsonData, err := os.ReadFile(outputPath)
	require.NoError(t, err, "Failed to read JSON file")

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonData, &parsed), "JSON parsing failed")

	assert.Contains(t, parsed, "format_version", "Missing format_version field")
	assert.Equal(t, "1.0", parsed["format_version"], "Version mismatch")

	assert.Contains(t, parsed, "provider_schemas", "Missing provider_schemas")
	providerSchemas := parsed["provider_schemas"].(map[string]interface{})
	assert.Contains(t, providerSchemas, DefaultProviderName, "Provider entry missing")

}

func InitPreCreateMaxcomputeUser(t *testing.T, username string) (string, string) {
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skip("Failed to initialize MaxcomputeUser: unable to create sharedClient")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := &MaxcomputeService{client}
	userId, userPk, err := maxcomputeService.GetOrCreateMaxcomputeUser(username)
	if err != nil {
		t.Skip("Failed to initialize MaxcomputeUser: GetOrCreateMaxcomputeUser request failed")
	}
	return userId, userPk
}
