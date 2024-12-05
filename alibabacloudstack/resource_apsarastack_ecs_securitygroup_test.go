package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsSecuritygroup0(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_securitygroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSecuritygroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesecuritygroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssecurity_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSecuritygroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "create",

					"security_group_name": "createname01",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"security_group_type": "normal",

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "create",

						"security_group_name": "createname01",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"security_group_type": "normal",

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "create",

					"security_group_name": "createname01",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"security_group_type": "normal",

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "create",

						"security_group_name": "createname01",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"security_group_type": "normal",

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccEcsSecuritygroupCheckmap = map[string]string{

	"ecs_count": CHECKSET,

	"description": CHECKSET,

	"resource_group_id": CHECKSET,

	"security_group_name": CHECKSET,

	"service_managed": CHECKSET,

	"create_time": CHECKSET,

	"security_group_id": CHECKSET,

	"security_group_references": CHECKSET,

	"security_group_type": CHECKSET,

	"available_instance_amount": CHECKSET,

	"service_id": CHECKSET,

	"vpc_id": CHECKSET,

	"permissions": CHECKSET,

	"inner_access_policy": CHECKSET,

	"region_id": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccEcsSecuritygroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s



`, name, VpcCommonTestCase)
}
func TestAccAlibabacloudStackEcsSecuritygroup1(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_securitygroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSecuritygroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesecuritygroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssecurity_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSecuritygroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultvx5yPv.VpcId)}}",

					"security_group_type": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultvx5yPv.VpcId)}}",

						"security_group_type": "normal",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultvx5yPv.VpcId)}}",

					"security_group_type": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultvx5yPv.VpcId)}}",

						"security_group_type": "normal",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackEcsSecuritygroup2(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_securitygroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSecuritygroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesecuritygroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssecurity_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSecuritygroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "testDescription",

					"security_group_name": "testSecurityGroupName",

					"security_group_type": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "testDescription",

						"security_group_name": "testSecurityGroupName",

						"security_group_type": "normal",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackEcsSecuritygroup3(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_securitygroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSecuritygroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesecuritygroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssecurity_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSecuritygroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "TestSg",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"security_group_name": "tf auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "TestSg",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"security_group_name": "tf auto",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "112233",

					"security_group_name": "tttttt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "112233",

						"security_group_name": "tttttt",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "112233",

					"security_group_name": "tttttt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "112233",

						"security_group_name": "tttttt",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "112233",

					"security_group_name": "tttttt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "112233",

						"security_group_name": "tttttt",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "112233",

					"security_group_name": "tttttt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "112233",

						"security_group_name": "tttttt",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "112233",

					"security_group_name": "tttttt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "112233",

						"security_group_name": "tttttt",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackEcsSecuritygroup4(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_securitygroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSecuritygroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesecuritygroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssecurity_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSecuritygroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "testDescription",

					"security_group_name": "testSecurityGroupName",

					"security_group_type": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "testDescription",

						"security_group_name": "testSecurityGroupName",

						"security_group_type": "normal",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackEcsSecuritygroup5(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_securitygroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSecuritygroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesecuritygroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssecurity_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSecuritygroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultvx5yPv.VpcId)}}",

					"security_group_type": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultvx5yPv.VpcId)}}",

						"security_group_type": "normal",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultvx5yPv.VpcId)}}",

					"security_group_type": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultvx5yPv.VpcId)}}",

						"security_group_type": "normal",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "alibabacloudstack_security_group.default.id",

					"security_group_name": "alibabacloudstack_security_group.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "alibabacloudstack_security_group.default.id",

						"security_group_name": "alibabacloudstack_security_group.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackEcsSecuritygroup6(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_securitygroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSecuritygroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesecuritygroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssecurity_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSecuritygroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "create",

					"security_group_name": "createname",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"security_group_type": "normal",

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "create",

						"security_group_name": "createname",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"security_group_type": "normal",

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "update",

					"security_group_name": "updatename",

					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "update",

						"security_group_name": "updatename",

						"inner_access_policy": "Accept",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackEcsSecuritygroup7(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_securitygroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSecuritygroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesecuritygroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssecurity_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSecuritygroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "create",

					"security_group_name": "createname02",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"security_group_type": "normal",

					"inner_access_policy": "Drop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "create",

						"security_group_name": "createname02",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"security_group_type": "normal",

						"inner_access_policy": "Drop",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "update",

					"security_group_name": "updatename02",

					"inner_access_policy": "Drop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "update",

						"security_group_name": "updatename02",

						"inner_access_policy": "Drop",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

