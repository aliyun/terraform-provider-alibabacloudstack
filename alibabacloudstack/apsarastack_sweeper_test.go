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
	var accessKey, secretKey, proxy, domain, ossEndpoint, essEndpoint, slbEndpoint, crEndpoint, vpcEndpoint, rdsEndpoint, ecsEndpoint, kvStoreEndpoint, rgsName, rgid, dept string
	var insecure bool
	if accessKey = os.Getenv("ALIBABACLOUDSTACK_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALIBABACLOUDSTACK_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_SECRET_KEY")
	}
	insecure, _ = strconv.ParseBool(os.Getenv("ALIBABACLOUDSTACK_INSECURE"))

	if proxy = os.Getenv("ALIBABACLOUDSTACK_PROXY"); proxy == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_PROXY")
	}
	if domain = os.Getenv("ALIBABACLOUDSTACK_DOMAIN"); domain == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_DOMAIN")
	}
	if ossEndpoint = os.Getenv("OSS_ENDPOINT"); ossEndpoint == "" {
		ossEndpoint = domain
	}
	if rdsEndpoint := os.Getenv("RDS_ENDPOINT"); rdsEndpoint == "" {
		rdsEndpoint = domain
	}
	if essEndpoint = os.Getenv("ESS_ENDPOINT"); essEndpoint == "" {
		essEndpoint = domain
	}
	if ecsEndpoint = os.Getenv("ECS_ENDPOINT"); ecsEndpoint == "" {
		ecsEndpoint = domain
	}
	if vpcEndpoint = os.Getenv("VPC_ENDPOINT"); vpcEndpoint == "" {
		vpcEndpoint = domain
	}
	if slbEndpoint = os.Getenv("SLB_ENDPOINT"); slbEndpoint == "" {
		slbEndpoint = domain
	}
	if crEndpoint = os.Getenv("CR_ENDPOINT"); crEndpoint == "" {
		crEndpoint = domain
	}
	if kvStoreEndpoint = os.Getenv("KVSTORE_ENDPOINT"); kvStoreEndpoint == "" {
		kvStoreEndpoint = domain
	}
	if rgid = os.Getenv("ALIBABACLOUDSTACK_RESOURCE_GROUP"); rgid == "" {
		//return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_DOMAIN")
	}
	if dept = os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT"); dept == "" {
		//return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_DOMAIN")
	}
	if rgsName = os.Getenv("ALIBABACLOUDSTACK_RESOURCE_GROUP_SET"); rgsName == "" {
		return nil, fmt.Errorf("empty ALIBABACLOUDSTACK_RESOURCE_GROUP_SET")
	}

	conf := connectivity.Config{
		Region:    connectivity.Region(region),
		RegionId:  region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Proxy:     proxy,
		Insecure:  insecure,
		Domain:    domain,
		Protocol:  "HTTP",
		Endpoints: map[connectivity.ServiceCode]string{
			connectivity.EcsCode:     ecsEndpoint,
			connectivity.VPCCode:     vpcEndpoint,
			connectivity.RDSCode:     rdsEndpoint,
			connectivity.ESSCode:     essEndpoint,
			connectivity.KVSTORECode: kvStoreEndpoint,
			connectivity.OSSCode:     ossEndpoint,
			connectivity.CRCode:      crEndpoint,
			connectivity.SLBCode:     slbEndpoint,
		},

		ResourceGroup:   rgid,
		Department:      dept,
		ResourceSetName: rgsName,
	}
	if accountId := os.Getenv("ALIBABACLOUDSTACK_ACCOUNT_ID"); accountId != "" {
		conf.AccountId = accountId
	}

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}
