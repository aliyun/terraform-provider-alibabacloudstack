package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alibabacloudstack_dns_group",
		&resource.Sweeper{
			Name: "alibabacloudstack_dns_group",
			F:    testSweepDnsGroup,
		})
}

func testSweepDnsGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	request := alidns.CreateDescribeDomainGroupsRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup

	var allGroups []alidns.DomainGroup
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if strings.HasPrefix(domainGroup.GroupName, "tf-testacc") {
				allGroups = append(allGroups, domainGroup)
			} else {
				log.Printf("Skip %#v.", domainGroup)
			}
		}
		if len(groups) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	removeRequest := alidns.CreateDeleteDomainGroupRequest()
	removeRequest.Headers = map[string]string{"RegionId": client.RegionId}
	removeRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
	removeRequest.QueryParams["Department"] = client.Department
	removeRequest.QueryParams["ResourceGroup"] = client.ResourceGroup

	for _, group := range allGroups {
		removeRequest.GroupId = group.GroupId
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainGroup(removeRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		addDebug(request.GetActionName(), raw)
	}
	return nil
}

func TestAccAlibabacloudStackDnsGroup_basic(t *testing.T) {
	testAccPreCheckWithAPIIsNotSupport(t)
	resourceId := "alibabacloudstack_dns_group.default"
	var v alidns.DomainGroup
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 99999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testaccdns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDnsGroupConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("tf-testaccdns%d", rand-1),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testaccdns%d", rand-1),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func resourceDnsGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}
