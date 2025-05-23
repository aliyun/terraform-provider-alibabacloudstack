package alibabacloudstack

import (
	"fmt"

	"testing"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_ots_instance", &resource.Sweeper{
		Name: "alibabacloudstack_ots_instance",
		F:    testSweepOtsInstances,
	})
}

func testSweepOtsInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var insts []ots.InstanceInfo
	req := ots.CreateListInstanceRequest()
	req.RegionId = client.RegionId
	req.Method = "GET"
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNum = requests.NewInteger(1)
	for {
		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.ListInstance(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving OTS Instances: %s", err)
		}
		resp, _ := raw.(*ots.ListInstanceResponse)
		if resp == nil || len(resp.InstanceInfos.InstanceInfo) < 1 {
			break
		}
		insts = append(insts, resp.InstanceInfos.InstanceInfo...)

		if len(resp.InstanceInfos.InstanceInfo) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNum); err != nil {
			return err
		} else {
			req.PageNum = page
		}
	}
	sweeped := false

	for _, v := range insts {
		name := v.InstanceName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping OTS Instance: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting OTS Instance %s table stores.", name)
		raw, err := otsService.client.WithTableStoreClient(name, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.ListTable()
		})
		if err != nil {
			log.Printf("[ERROR] List OTS Instance %s table stores got an error: %#v.", name, err)
		}
		tables, _ := raw.(*tablestore.ListTableResponse)
		if tables != nil && len(tables.TableNames) > 0 {
			for _, t := range tables.TableNames {
				req := new(tablestore.DeleteTableRequest)
				req.TableName = t
				if _, err := otsService.client.WithTableStoreClient(name, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
					return tableStoreClient.DeleteTable(req)
				}); err != nil {
					log.Printf("[ERROR] Delete OTS Instance %s table store %s got an error: %#v.", name, t, err)
				}
			}
			time.Sleep(30 * time.Second)
		}
		log.Printf("[INFO] Deleting OTS Instance: %s", name)
		req := ots.CreateDeleteInstanceRequest()
		req.InstanceName = name
		_, err = client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete OTS Instance (%s): %s", name, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		time.Sleep(3 * time.Minute)
	}
	return nil
}

func TestAccAlibabacloudStackOtsInstance_basic(t *testing.T) {
	var v InstanceInfo

	resourceId := "alibabacloudstack_ots_instance.default"
	ra := resourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceConfigDependence)

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
					"name":          name,
					"description":   name,
					"instance_type": "Capacity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          name,
						"description":   name,
						"instance_type": "Capacity",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			/*
				{
					Config: testAccConfig(map[string]interface{}{
						"accessed_by": "Vpc",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"accessed_by": "Vpc",
						}),
					),
				},
			*/
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
						"tags.Updated": "TF",
					}),
				),
			},
			/*
				{
					Config: testAccConfig(map[string]interface{}{
						"accessed_by": REMOVEKEY,
						"tags":        REMOVEKEY,
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"accessed_by":  "Any",
							"tags.%":       "0",
							"tags.Created": REMOVEKEY,
							"tags.For":     REMOVEKEY,
							"tags.Updated": REMOVEKEY,
						}),
					),
				},
			*/
		},
	})
}

/*
	func TestAccAlibabacloudStackOtsInstanceHighPerformance(t *testing.T) {
		var v ots.InstanceInfo

		resourceId := "alibabacloudstack_ots_instance.default"
		ra := resourceAttrInit(resourceId, otsInstanceBasicMap)

		serviceFunc := func() interface{} {
			return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
		}
		rc := resourceCheckInit(resourceId, &v, serviceFunc)

		rac := resourceAttrCheckInit(rc, ra)

		testAccCheck := rac.resourceAttrMapUpdateSet()
		rand := getAccTestRandInt(10000, 99999)
		name := fmt.Sprintf("tf-testAcc%d", rand)
		testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceConfigDependence)

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
						"name":          name,
						"description":   name,
						"instance_type": "HighPerformance",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"name":          name,
							"description":   name,
							"instance_type": "HighPerformance",
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
						"accessed_by": "Vpc",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"accessed_by": "Vpc",
						}),
					),
				},
				{
					Config: testAccConfig(map[string]interface{}{
						"tags": map[string]string{
							"Created": "TF",
							"For":     "acceptance test",
						},
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"tags.%":       "2",
							"tags.Created": "TF",
							"tags.For":     "acceptance test",
						}),
					),
				},
				{
					Config: testAccConfig(map[string]interface{}{
						"tags": map[string]string{
							"Created": "TF",
							"For":     "acceptance test",
							"Updated": "TF",
						},
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"tags.%":       "3",
							"tags.Created": "TF",
							"tags.For":     "acceptance test",
							"tags.Updated": "TF",
						}),
					),
				},
				{
					Config: testAccConfig(map[string]interface{}{
						"accessed_by": REMOVEKEY,
						"tags":        REMOVEKEY,
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"accessed_by":  "Any",
							"tags.%":       "0",
							"tags.Created": REMOVEKEY,
							"tags.For":     REMOVEKEY,
							"tags.Updated": REMOVEKEY,
						}),
					),
				},
			},
		})
	}

	func TestAccAlibabacloudStackOtsInstance_multi(t *testing.T) {
		var v ots.InstanceInfo

		resourceId := "alibabacloudstack_ots_instance.default.4"
		ra := resourceAttrInit(resourceId, otsInstanceBasicMap)

		serviceFunc := func() interface{} {
			return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
		}
		rc := resourceCheckInit(resourceId, &v, serviceFunc)

		rac := resourceAttrCheckInit(rc, ra)

		testAccCheck := rac.resourceAttrMapUpdateSet()
		rand := getAccTestRandInt(10000, 99999)
		name := fmt.Sprintf("tf-testAcc%d", rand)
		testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceConfigDependence)

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
						"name":          name + "${count.index}",
						"description":   name + "${count.index}",
						"instance_type": "Capacity",
						"count":         "5",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(nil),
					),
				},
			},
		})
	}
*/
func resourceOtsInstanceConfigDependence(name string) string {
	return ""
}

var otsInstanceBasicMap = map[string]string{
	"name": CHECKSET,
	//"accessed_by":   "Any",
	"instance_type": CHECKSET,
	"description":   CHECKSET,
}

func testAccCheckOtsInstanceExist(n string, instance InstanceInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		otsService := OtsService{client}

		response, err := otsService.DescribeOtsInstance(rs.Primary.ID)

		if err != nil {
			return err
		}
		instance = response
		return nil
	}
}
