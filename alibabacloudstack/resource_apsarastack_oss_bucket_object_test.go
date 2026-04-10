package alibabacloudstack

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackOssBucketObject_basic(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-oss-object-test-acc-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// first write some data to the tempfile just so it's not 0 bytes.
	err = ioutil.WriteFile(tmpFile.Name(), []byte("{anything will do }"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var v http.Header
	resourceId := "alibabacloudstack_oss_bucket_object.default"
	ra := resourceAttrInit(resourceId, ossBucketObjectBasicMap)
	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-object-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketObjectConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":       "${alibabacloudstack_oss_bucket.default.bucket}",
					"key":          "test-object-source-key",
					"source":       strings.Replace(tmpFile.Name(), "\\", "\\\\", -1),
					"content_type": "binary/octet-stream",
					"acl":          "public-read-write",
					"content_md5":  "ewBv9NcPaMxlBhrPL4Aubw==",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudOssBucketObjectExists(
						"alibabacloudstack_oss_bucket_object.default", name, v),
					testAccCheck(map[string]string{
						"bucket": name,
						"source": tmpFile.Name(),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// source is a local attribute, cannot be loaded from remote
				// acl requires special permissions, currently cannot be adjusted during testing
				ImportStateVerifyIgnore: []string{"source", "acl"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption": "AES256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption": "AES256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption": "KMS",
					"kms_key_id":             "${alibabacloudstack_kms_key.key.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption": "KMS",
						"kms_key_id":             CHECKSET,
					}),
				),
			},
			/*
				{
					Config: testAccConfig(map[string]interface{}{
						"source":  REMOVEKEY,
						"content": "some words for test oss object content",
						"acl":     "public-read-write",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAlicloudOssBucketObjectExists(
							"alibabacloudstack_oss_bucket_object.default", name, v),
						testAccCheck(map[string]string{
							"source":  REMOVEKEY,
							"content": "some words for test oss object content",
						}),
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
						"server_side_encryption": "KMS",
						"kms_key_id":             "${data.alibabacloudstack_kms_keys.enabled.ids.0}",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{}),
					),
				},
				{
					Config: testAccConfig(map[string]interface{}{
						"bucket":                 "${alibabacloudstack_oss_bucket.default.bucket}",
						"server_side_encryption": "AES256",
						"kms_key_id":             REMOVEKEY,
						"key":                    "test-object-source-key",
						"content":                REMOVEKEY,
						"source":                 tmpFile.Name(),
						"content_type":           "binary/octet-stream",
						"acl":                    REMOVEKEY,
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAlicloudOssBucketObjectExists(
							"alibabacloudstack_oss_bucket_object.default", name, v),
						testAccCheck(map[string]string{
							"bucket":       name,
							"key":          "test-object-source-key",
							"content":      REMOVEKEY,
							"source":       tmpFile.Name(),
							"content_type": "binary/octet-stream",
							"acl":          "public-read-write",
						}),
					),
				},

			*/
		},
	})
}

func resourceOssBucketObjectConfigDependence(name string) string {

	return fmt.Sprintf(`
resource "alibabacloudstack_oss_bucket" "default" {
	bucket = "%s"
	acl = "public-read-write"
}
data "alibabacloudstack_kms_keys" "enabled" {
	status = "%s"
}
%s
`, name, string(EnabledStatus), KeyCommonTestCase)
}

var ossBucketObjectBasicMap = map[string]string{
	"bucket":       CHECKSET,
	"key":          "test-object-source-key",
	"source":       CHECKSET,
	"content_type": "binary/octet-stream",
	"acl":          "public-read-write",
}

func testAccCheckAlicloudOssBucketObjectExists(n string, bucket string, obj http.Header) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOssBucketObjectExistsWithProviders(n, bucket, obj, &providers)
}
func testAccCheckOssBucketObjectExistsWithProviders(n string, bucketName string, obj http.Header, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}
			client := provider.Meta().(*connectivity.AlibabacloudStackClient)
			ossService := OssService{client}
			buck, err := ossService.GetBucketClient(bucketName)
			if err != nil {
				return fmt.Errorf("Error getting bucket: %#v", err)
			}
			id_info := strings.SplitN(rs.Primary.ID, ":", 2)
			key := id_info[1]
			request := &oss.GetObjectMetaRequest{
				Bucket: &bucketName,
				Key:    &key,
			}
			object, err := buck.GetObjectMeta(context.Background(), request)
			log.Printf("[WARN]get oss bucket object %#v", bucketName)
			if err == nil {
				if object != nil {
					obj = object.Headers
					return nil
				}
				continue
			} else {
				return err

			}
		}

		return fmt.Errorf("Bucket not found")
	}
}
