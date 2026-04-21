package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOssBucket_Basic(t *testing.T) {
	var v *oss.BucketProperties

	resourceId := "alibabacloudstack_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":           name,
						"storage_capacity": "1024",
						"tags.%":           "2",
						"tags.Created":     "TF",
						"tags.For":         "Test",
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
					"logging": []map[string]interface{}{{
						"target_bucket": "${var.name}",
						"target_prefix": "oss-accesslog/",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{{
						"target_bucket": "${var.name}",
						"target_prefix": "oss-accesslog-update/",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "public-read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sse_algorithm": "AES256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sse_algorithm": "AES256",
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"sse_algorithm": "KMS",
			// 		// "kms_key_id":    "${local.kms_key}",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"sse_algorithm": "KMS",
			// 			"kms_key_id":    CHECKSET,
			// 		}),
			// 	),
			// },
			{
				Config: testAccConfig(map[string]interface{}{
					"sse_algorithm": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sse_algorithm": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_capacity": "2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_capacity": "2048",
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
					"tags": map[string]string{
						"For": "Test-reset",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.For":     "Test-reset",
						"tags.Created": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":   REMOVEKEY,
						"tags.For": REMOVEKEY,
					}),
				),
			},
			// In version v3.16.2, OSS does not temporarily support the delete method for tags, so tags cannot be deleted completely
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": REMOVEKEY,
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "0",
			// 			"tags.Created": REMOVEKEY,
			// 			"tags.For":     REMOVEKEY,
			// 		}),
			// 	),
			// },
		},
	})
}

func TestUatAlibabacloudStackOssBucket_Sync(t *testing.T) {
	var v *oss.BucketProperties

	resourceId := "alibabacloudstack_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	dual_sync_role := os.Getenv("ALIBABACLOUDSTACK_OSS_DUAL_SYNC_ROLE")
	if dual_sync_role == "" {
		dual_sync_role = "AliyunOSSPrivateCloudDrsSyncRole"
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketDualDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":      name,
					"oss_cluster": "${data.alibabacloudstack_oss_clusters.default.clusters.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
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
					"bucket_sync":    "true",
					"dual_kms_key":   "${local.kms_key}",
					"dual_sync_role": dual_sync_role,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_sync":    "true",
						"dual_kms_key":   CHECKSET,
						"dual_sync_role": dual_sync_role,
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackOssBucket_Nokmskey(t *testing.T) {
	var v *oss.BucketProperties

	resourceId := "alibabacloudstack_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketNokmskeyDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":           name,
						"storage_capacity": "1024",
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
					"logging": []map[string]interface{}{{
						"target_bucket": "${var.name}",
						"target_prefix": "oss-accesslog/",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{{
						"target_bucket": "${var.name}",
						"target_prefix": "oss-accesslog-update/",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "public-read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sse_algorithm": "AES256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sse_algorithm": "AES256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sse_algorithm": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sse_algorithm": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_capacity": "2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_capacity": "2048",
					}),
				),
			},
		},
	})
}

// func TestUatAlibabacloudStackOssBucket_Vpc(t *testing.T) {
// 	var v *oss.BucketProperties

// 	resourceId := "alibabacloudstack_oss_bucket.default"
// 	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

// 	serviceFunc := func() interface{} {
// 		return &OssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}
// 	rc := resourceCheckInit(resourceId, &v, serviceFunc)

// 	rac := resourceAttrCheckInit(rc, ra)

// 	testAccCheck := rac.resourceAttrMapUpdateSet()
// 	rand := getAccTestRandInt(1000000, 9999999)
// 	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketConfigDependence)
// 	ResourceTest(t, resource.TestCase{
// 		PreCheck: func() {
// 			testAccPreCheck(t)
// 		},
// 		// module name
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,
// 		CheckDestroy:  rac.checkResourceDestroy(),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"bucket":  name,
// 					"vpclist": []string{"${alibabacloudstack_vpc.vpc.id}", "${alibabacloudstack_vpc.vpc2.id}"},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"bucket":    name,
// 						"vpclist.0": CHECKSET,
// 						"vpclist.#": "2",
// 					}),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceId,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"bucket":  name,
// 					"vpclist": []string{"${alibabacloudstack_vpc.vpc.id}"},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"bucket":    name,
// 						"vpclist.0": CHECKSET,
// 						"vpclist.#": "1",
// 					}),
// 				),
// 			},
// 		},
// 	})
// }

func resourceOssBucketNokmskeyDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

%s

`, name)
}

func resourceOssBucketConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}
	
resource "alibabacloudstack_vpc" "vpc" {
	name = "${var.name}-v"
	cidr_block = "192.168.0.0/24"
}
resource "alibabacloudstack_vpc" "vpc2" {
	name = "${var.name}-v2"
	cidr_block = "192.168.0.0/24"
}

%s

`, name, KeyCommonTestCase)
}

func resourceOssBucketDualDependence(name string) string {
	return fmt.Sprintf(`

%s

data "alibabacloudstack_oss_clusters" "default" {
}
	
`, KeyCommonTestCase)
}

var ossBucketBasicMap = map[string]string{
	"creation_date":    CHECKSET,
	"lifecycle_rule.#": "0",
}
