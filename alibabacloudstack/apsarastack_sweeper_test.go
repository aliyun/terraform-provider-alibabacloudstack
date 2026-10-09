package alibabacloudstack

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sweepFunc is the standard signature for all resource sweepers.
// Region is obtained from environment variables inside each sweeper via sharedClientForRegion("").
type sweepFunc func(prefixes []string) error

// broadPrefixes is the common test resource prefix pair used in broad sweep mode.
var broadPrefixes = []string{"tf-testAcc", "tf_testAcc"}

// sweeperEntry binds a sweeper function with its precise/broad prefixes and error handling policy.
type sweeperEntry struct {
	fn      sweepFunc // the sweeper function
	precise []string  // prefixes for precise mode (default)
	broad   []string  // prefixes for broad mode (default: nil, only first entry needs this)
	fatal   bool      // true: return err on failure; false: log warning only (default: false)
	waitSec int       // seconds to wait after this entry completes; 0 means default 5s
}

// SweeperFunc composes multiple sweeperEntries into a resource.Sweeper.F function.
//   - Precise mode (default): all entries run with their precise prefixes.
//   - Broad mode (ALIBABACLOUDSTACK_TF_SWEEP_BROAD=1): only the first entry runs, with its broad prefixes.
//   - fatal controls whether failure returns error (true) or logs warning (false).
//   - waitSec controls seconds to sleep after each entry; 0 means default 5s.
//
// Usage:
//
//	SweeperFunc(
//	    sweeperEntry{fn: testSweepCen, precise: []string{"tf-testAccCen"}, broad: broadPrefixes, fatal: true},
//	    sweeperEntry{fn: testSweepVSwitch, precise: []string{"tf-testAccCen_vsw"}},
//	    sweeperEntry{fn: testSweepVPC, precise: []string{"tf-testAccCen_vpc"}},
//	)
// SweeperFunc composes multiple sweeperEntries into a resource.Sweeper.F function.
//   - Precise mode (default): all entries run with their precise prefixes.
//   - Broad mode (ALIBABACLOUDSTACK_TF_SWEEP_BROAD=1): only entries with broad prefixes run, using their broad prefixes.
//   - Custom prefix mode (ALIBABACLOUDSTACK_TF_SWEEP_PREFIXES=prefix1,prefix2): only works in broad mode, all entries run with the specified custom prefixes instead of broad prefixes.
//   - fatal controls whether failure returns error (true) or logs warning (false).
//   - waitSec controls seconds to sleep after each entry; 0 means default 5s.
//
// Environment variable: ALIBABACLOUDSTACK_TF_SWEEP_BROAD=1 enables broad mode; ALIBABACLOUDSTACK_TF_SWEEP_PREFIXES overrides broad prefixes (requires BROAD=1).
func SweeperFunc(entries ...sweeperEntry) error {
	broad := os.Getenv("ALIBABACLOUDSTACK_TF_SWEEP_BROAD") == "1"

	// Custom prefix mode: only works in broad mode, use user-specified prefixes for all entries
	var customPrefixes []string
	if broad {
		if envPrefixes := os.Getenv("ALIBABACLOUDSTACK_TF_SWEEP_PREFIXES"); envPrefixes != "" {
			for _, p := range strings.Split(envPrefixes, ",") {
				p = strings.TrimSpace(p)
				if p != "" {
					customPrefixes = append(customPrefixes, p)
				}
			}
		}
	}

	for _, e := range entries {
		var prefixes []string
		if broad {
			if len(customPrefixes) > 0 {
				prefixes = customPrefixes
			} else if len(e.broad) == 0 {
				continue
			} else {
				prefixes = e.broad
			}
		} else {
			prefixes = e.precise
		}
		if err := e.fn(prefixes); err != nil {
			if e.fatal {
				return err
			}
			log.Printf("[WARN] Sweeper %s error: %s", prefixes, err)
		}
		wait := e.waitSec
		if wait == 0 {
			wait = 5
		}
		time.Sleep(time.Duration(wait) * time.Second)
	}
	return nil
}

// functions for a given region
func sharedClientForRegion(region string) (interface{}, error) {
	return sharedClientForRegionWithCredentials(region, false)
}

// sharedClientForRegionOrg creates a client using ORG credentials (skipping TEST_USER).
func sharedClientForRegionOrg(region string) (interface{}, error) {
	return sharedClientForRegionWithCredentials(region, true)
}

func sharedClientForRegionWithCredentials(region string, forceOrg bool) (interface{}, error) {
	var proxy, domain, popgwDomain, rgsName, rgid, dept, protocol string
	var insecure, is_center_region bool
	accessKey, secretKey := resolveTestCredentials(forceOrg)
	if accessKey == "" {
		return nil, fmt.Errorf("no valid access key found")
	}

	if secretKey == "" {
		return nil, fmt.Errorf("no valid secret key found")
	}
	insecure, _ = strconv.ParseBool(os.Getenv("ALIBABACLOUDSTACK_INSECURE"))
	if popgwDomain = os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"); popgwDomain == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_POPGW_DOMAIN")
	}
	if rgsName = os.Getenv("ALIBABACLOUDSTACK_RESOURCE_GROUP_SET"); rgsName == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_RESOURCE_GROUP_SET")
	}
	if rgid = os.Getenv("ALIBABACLOUDSTACK_RESOURCE_GROUP"); rgid == "" && rgsName == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_DOMAIN")
	}
	if dept = os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT"); dept == "" && rgsName == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_DOMAIN")
	}
	if protocol = os.Getenv("ALIBABACLOUDSTACK_PROTOCOL"); protocol == "" {
		protocol = "HTTP"
	}

	if region == "" {
		if v := os.Getenv("ALIBABACLOUDSTACK_REGION"); v == "" {
			return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_REGION")
		} else {
			region = v
		}
	}

	if is_center_region_str := os.Getenv("ALIBABACLOUDSTACK_CENTER_REGION"); is_center_region_str == "" {
		is_center_region = true
	} else {
		is_center_region, _ = strconv.ParseBool(os.Getenv("ALIBABACLOUDSTACK_CENTER_REGION"))
	}

	conf := connectivity.Config{
		Region:          connectivity.Region(region),
		RegionId:        region,
		AccessKey:       accessKey,
		SecretKey:       secretKey,
		Proxy:           proxy,
		Insecure:        insecure,
		Domain:          domain,
		Protocol:        protocol,
		Endpoints:       map[connectivity.ServiceCode]string{},
		ResourceGroup:   rgid,
		Department:      dept,
		ResourceSetName: rgsName,
		PopgwDomain:     popgwDomain,
	}
	if accountId := os.Getenv("ALIBABACLOUDSTACK_ACCOUNT_ID"); accountId != "" {
		conf.AccountId = accountId
	}

	if proxy = os.Getenv("ALIBABACLOUDSTACK_PROXY"); proxy != "" {
		conf.Proxy = proxy
	}

	for popcode := range connectivity.PopEndpoints {
		if domain != "" {
			conf.Endpoints[popcode] = domain
		} else {
			endpoint := connectivity.GeneratorEndpoint(popcode, region, popgwDomain, is_center_region)
			if endpoint != "" {
				conf.Endpoints[popcode] = endpoint
			}
		}
	}
	if asapiEndpoint := os.Getenv("ALIBABACLOUDSTACK_ASAPI_ENDPOINT"); asapiEndpoint != "" {
		conf.Endpoints[connectivity.ASAPICode] = asapiEndpoint
	}
	if conf.Department == "" || conf.ResourceGroup == "" {
		dept, resId, rgid, err := getResourceCredentials(&conf)
		if err != nil {
			return nil, err
		}
		conf.Department = dept
		conf.ResourceGroup = fmt.Sprintf("%d", rgid)
		conf.ResourceGroupId = resId
	}

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}
