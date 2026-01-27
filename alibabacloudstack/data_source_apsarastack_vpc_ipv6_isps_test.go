package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpcIpv6IspsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_vpc_ipv6_isps.isps"
	name := fmt.Sprintf("tf-testacc-vpcipv6isps%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcIpv6IspsConfigDependence)

	defaultConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}
	
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_vpc_ipv6_isps.anyone.ids.0}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_vpc_ipv6_isps.anyone.ids.0}_fake"},
		}),
	}

	serviceProviderConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_provider": "BGP",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_provider": "fake_provider",
		}),
	}

	lockStatusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"lock_status": "unlocked",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"lock_status": "locked",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_provider": "BGP",
			"lock_status":      "unlocked",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_provider": "fake_provider",
			"lock_status":      "locked",
		}),
	}

	var existVpcIpv6IspsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ipv6_isps.#":                    CHECKSET,
			"ids.#":                          CHECKSET,
			"ipv6_isps.0.id":                 CHECKSET,
			"ipv6_isps.0.service_provider":   CHECKSET,
			"ipv6_isps.0.zone_id":            CHECKSET,
			"ipv6_isps.0.type":               CHECKSET,
			"ipv6_isps.0.cidr_block":         CHECKSET,
			"ipv6_isps.0.available_count":    CHECKSET,
			"ipv6_isps.0.in_use_count":       CHECKSET,
			"ipv6_isps.0.lock_status":        CHECKSET,
			"ipv6_isps.0.need_declare":       CHECKSET,
			"ipv6_isps.0.pool_id":            CHECKSET,
			"ipv6_isps.0.ula":                CHECKSET,
		}
	}

	var fakeVpcIpv6IspsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ipv6_isps.#": "0",
			"ids.#":       "0",
		}
	}

	var vpcIpv6IspsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcIpv6IspsMapFunc,
		fakeMapFunc:  fakeVpcIpv6IspsMapFunc,
	}
	vpcIpv6IspsCheckInfo.dataSourceTestCheck(t, rand, defaultConf, idsConf, serviceProviderConf, lockStatusConf, allConf)
}

func dataSourceVpcIpv6IspsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

data alibabacloudstack_vpc_ipv6_isps anyone {
}

`, name, DataZoneCommonTestCase)
}
