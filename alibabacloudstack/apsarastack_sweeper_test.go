package alibabacloudstack

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// functions for a given region
func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, proxy, domain, popgw_domain, rgsName, rgid, dept ,protocol string
	var insecure, is_center_region bool
	if accessKey = os.Getenv("ALIBABACLOUDSTACK_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALIBABACLOUDSTACK_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_SECRET_KEY")
	}
	insecure, _ = strconv.ParseBool(os.Getenv("ALIBABACLOUDSTACK_INSECURE"))

	//if proxy = os.Getenv("ALIBABACLOUDSTACK_PROXY"); proxy == "" {
	//	return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_PROXY")
	//}
	if domain = os.Getenv("ALIBABACLOUDSTACK_DOMAIN"); domain == "" {
		if popgw_domain = os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"); popgw_domain == "" {
			return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_DOMAIN or ALIBABACLOUDSTACK_POPGW_DOMAIN")
		}
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
	
	if is_center_region_str := os.Getenv("ALIBABACLOUDSTACK_CENTER_REGION"); is_center_region_str == "" {
		is_center_region = true
	}else {
		is_center_region, _ = strconv.ParseBool(os.Getenv("ALIBABACLOUDSTACK_CENTER_REGION"))
	}


	conf := connectivity.Config{
		Region:    connectivity.Region(region),
		RegionId:  region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Proxy:     proxy,
		Insecure:  insecure,
		Domain:    domain,
		Protocol:  protocol,
		Endpoints: map[connectivity.ServiceCode]string{},
		ResourceGroup:   rgid,
		Department:      dept,
		ResourceSetName: rgsName,
	}
	if accountId := os.Getenv("ALIBABACLOUDSTACK_ACCOUNT_ID"); accountId != "" {
		conf.AccountId = accountId
	}
	
	for popcode := range connectivity.PopEndpoints {
		if domain != "" {
			conf.Endpoints[popcode] = domain
		} else {
			endpoint := connectivity.GeneratorEndpoint(popcode, region, popgw_domain, is_center_region)
			if endpoint != "" {
				conf.Endpoints[popcode] = endpoint
			}
		}
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
