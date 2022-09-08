package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"log"
	"testing"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var datahubProjectSuffixMin = 100000
var datahubProjectSuffixMax = 999999

func init() {
	resource.AddTestSweepers("alibabacloudstack_datahub_project", &resource.Sweeper{
		Name: "alibabacloudstack_datahub_project",
		F:    testSweepDatahubProject,
	})
}

func testSweepDatahubProject(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DatahubSupportedRegions) {
		log.Printf("[INFO] Skipping Datahub unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	// List projects
	request := requests.NewCommonRequest()
	request.Method = "GET"         // Set request method
	request.Product = "datahub"    // Specify product
	request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2019-11-20" // Specify product version
	request.Scheme = "http"        // Set request scheme. Default: http
	request.ApiName = "ListProjects"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "datahub",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"Action":          "ListProjects",
		"Version":         "2019-11-20",
	}

	raw, err := client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
		return slsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		// Now, only some region support Datahub
		log.Printf("[ERROR] Failed to list Datahub projects: %s", err)
	}
	projects, _ := raw.(*datahub.ListProjectResult)

	for _, projectName := range projects.ProjectNames {
		// a testing project?
		if !isTerraformTestingDatahubObject(projectName) {
			log.Printf("[INFO] Skipping Datahub project: %s", projectName)
			continue
		}
		log.Printf("[INFO] Deleting project: %s", projectName)

		// List topics
		request := requests.NewCommonRequest()
		request.Method = "GET"         // Set request method
		request.Product = "datahub"    // Specify product
		request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
		request.Version = "2019-11-20" // Specify product version
		request.Scheme = "http"        // Set request scheme. Default: http
		request.ApiName = "ListTopics"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "datahub",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"RegionId":        client.RegionId,
			"Action":          "ListTopics",
			"Version":         "2019-11-20",
			"ProjectName":     projectName,
		}
		raw, err := client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
			return slsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return fmt.Errorf("error listing Datahub topics: %s", err)
		}
		topics, _ := raw.(*datahub.ListTopicResult)

		for _, topicName := range topics.TopicNames {
			log.Printf("[INFO] Deleting topic: %s/%s", projectName, topicName)

			// List subscriptions
			request := requests.NewCommonRequest()
			request.Method = "GET"         // Set request method
			request.Product = "datahub"    // Specify product
			request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
			request.Version = "2019-11-20" // Specify product version
			request.Scheme = "http"        // Set request scheme. Default: http
			request.ApiName = "ListSubscription"
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.QueryParams = map[string]string{
				"AccessKeySecret": client.SecretKey,
				"AccessKeyId":     client.AccessKey,
				"Product":         "datahub",
				"Department":      client.Department,
				"ResourceGroup":   client.ResourceGroup,
				"RegionId":        client.RegionId,
				"Action":          "ListSubscription",
				"Version":         "2019-11-20",
				"ProjectName":     projectName,
				"TopicName":       topicName,
			}
			raw, err := client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
				return slsClient.ProcessCommonRequest(request)
			})

			if err != nil {
				return fmt.Errorf("error listing Datahub subscriptions: %s", err)
			}
			subscriptions, _ := raw.(*datahub.ListSubscriptionResult)

			for _, subscription := range subscriptions.Subscriptions {
				log.Printf("[INFO] Deleting subscription: %s/%s/%s", projectName, topicName, subscription.SubId)

				// Delete subscription
				request := requests.NewCommonRequest()
				request.Method = "GET"         // Set request method
				request.Product = "datahub"    // Specify product
				request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
				request.Version = "2019-11-20" // Specify product version
				request.Scheme = "http"        // Set request scheme. Default: http
				request.ApiName = "DeleteSubscription"
				request.Headers = map[string]string{"RegionId": client.RegionId}
				request.QueryParams = map[string]string{
					"AccessKeySecret": client.SecretKey,
					"AccessKeyId":     client.AccessKey,
					"Product":         "datahub",
					"Department":      client.Department,
					"ResourceGroup":   client.ResourceGroup,
					"RegionId":        client.RegionId,
					"Action":          "DeleteSubscription",
					"Version":         "2019-11-20",
					"ProjectName":     projectName,
					"TopicName":       topicName,
				}
				_, err := client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
					return slsClient.ProcessCommonRequest(request)
				})
				if err != nil {
					log.Printf("[ERROR] Failed to delete Datahub subscriptions: %s/%s/%s", projectName, topicName, subscription.SubId)
					return fmt.Errorf("error deleting  Datahub subscriptions: %s/%s/%s", projectName, topicName, subscription.SubId)
				}
			}

			// Delete topic
			topicrequest := requests.NewCommonRequest()
			topicrequest.Method = "GET"         // Set request method
			topicrequest.Product = "datahub"    // Specify product
			topicrequest.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
			topicrequest.Version = "2019-11-20" // Specify product version
			topicrequest.Scheme = "http"        // Set request scheme. Default: http
			topicrequest.ApiName = "DeleteTopic"
			topicrequest.Headers = map[string]string{"RegionId": client.RegionId}
			topicrequest.QueryParams = map[string]string{
				"AccessKeySecret": client.SecretKey,
				"AccessKeyId":     client.AccessKey,
				"Product":         "datahub",
				"Department":      client.Department,
				"ResourceGroup":   client.ResourceGroup,
				"RegionId":        client.RegionId,
				"Action":          "DeleteTopic",
				"Version":         "2019-11-20",
				"ProjectName":     projectName,
				"TopicName":       topicName,
			}
			_, err = client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
				return slsClient.ProcessCommonRequest(topicrequest)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Datahub topic: %s/%s", projectName, topicName)
				return fmt.Errorf("[ERROR] Failed to delete Datahub topic: %s/%s", projectName, topicName)
			}
		}

		// Delete project
		projectrequest := requests.NewCommonRequest()
		projectrequest.Method = "GET"         // Set request method
		projectrequest.Product = "datahub"    // Specify product
		projectrequest.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
		projectrequest.Version = "2019-11-20" // Specify product version
		projectrequest.Scheme = "http"        // Set request scheme. Default: http
		projectrequest.ApiName = "DeleteProject"
		projectrequest.Headers = map[string]string{"RegionId": client.RegionId}
		projectrequest.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "datahub",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"RegionId":        client.RegionId,
			"Action":          "DeleteProject",
			"Version":         "2019-11-20",
			"ProjectName":     projectName,
		}
		_, err = client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
			return slsClient.ProcessCommonRequest(projectrequest)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Datahub project: %s", projectName)
			return fmt.Errorf("[ERROR] Failed to delete Datahub project: %s", projectName)
		}
	}

	return nil
}

func TestAccAlibabacloudStackDatahubProject_basic(t *testing.T) {
	var v *datahub.GetProjectResult
	resourceId := "alibabacloudstack_datahub_project.default"
	ra := resourceAttrInit(resourceId, datahubProjectBasicMap)
	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testaccdatahubproject%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubProjectConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"comment": "project for basic.",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"comment": "project for basic.",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"comment": REMOVEKEY,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"comment": "project added by terraform",
			//		}),
			//	),
			//},
		},
	})
}
func TestAccAlibabacloudStackDatahubProject_multi(t *testing.T) {
	var v *datahub.GetProjectResult

	resourceId := "alibabacloudstack_datahub_project.default.4"
	ra := resourceAttrInit(resourceId, datahubProjectBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testaccdatahubproject%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubProjectConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":  name + "${count.index}",
					"count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func resourceDatahubProjectConfigDependence(name string) string {
	return ""
}

var datahubProjectBasicMap = map[string]string{
	"name":    CHECKSET,
	"comment": "project added by terraform",
}

func testAccCheckDatahubProjectExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found Datahub project: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Datahub project ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		_, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
			return dataHubClient.GetProject(rs.Primary.ID)
		})

		if err != nil {
			return err
		}
		return nil
	}
}
