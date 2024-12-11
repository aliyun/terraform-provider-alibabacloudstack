package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_alikafka_sasl_user", &resource.Sweeper{
		Name: "alibabacloudstack_alikafka_sasl_user",
		F:    testSweepAlikafkaSaslUser,
		Dependencies: []string{
			"alibabacloudstack_alikafka_sasl_acl",
		},
	})
}

func testSweepAlikafkaSaslUser(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "error getting alibabacloudstack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = defaultRegionToTest

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)

	for _, v := range instanceListResp.InstanceList.InstanceVO {

		if v.ServiceStatus == 10 {
			log.Printf("[INFO] Skipping released alikafka instance id: %s ", v.InstanceId)
			continue
		}

		// Control the sasl user list request rate.
		time.Sleep(time.Duration(400) * time.Millisecond)

		request := alikafka.CreateDescribeSaslUsersRequest()
		request.InstanceId = v.InstanceId
		request.RegionId = defaultRegionToTest

		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DescribeSaslUsers(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka sasl users on instance (%s): %s", v.InstanceId, err)
			continue
		}

		saslUserListResp, _ := raw.(*alikafka.DescribeSaslUsersResponse)
		//saslUsers := saslUserListResp.SaslUserList.SaslUserVO
		saslUsers := saslUserListResp.SaslUserList.SaslUserVO
		for _, saslUser := range saslUsers {
			name := saslUser.Username
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping alikafka sasl username: %s ", name)
				continue
			}
			log.Printf("[INFO] delete alikafka sasl username: %s ", name)

			// Control the sasl username delete rate
			time.Sleep(time.Duration(400) * time.Millisecond)

			deleteUserReq := alikafka.CreateDeleteSaslUserRequest()
			deleteUserReq.InstanceId = v.InstanceId
			deleteUserReq.Username = saslUser.Username
			deleteUserReq.RegionId = defaultRegionToTest

			_, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.DeleteSaslUser(deleteUserReq)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete alikafka sasl username (%s): %s", name, err)
			}
		}
	}

	return nil
}

func TestAccAlibabacloudStackAlikafkaSaslUser_basic(t *testing.T) {

	//var v *alikafka.SaslUserVO
	var v *alikafka.SaslUserList
	resourceId := "alibabacloudstack_alikafka_sasl_user.default"
	ra := resourceAttrInit(resourceId, alikafkaSaslUserBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000,20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaSaslUserConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAlikafkaAclEnable(t)
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					//"instance_id": "${alibabacloudstack_alikafka_instance.default.id}",
					"instance_id": "cluster-private-paas-default",
					"username":    "${var.name}",
					"password":    "inputYourCodeHere",
					"type":        "scram",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand),
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
					"username": "newSaslUserName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": "newSaslUserName"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "inputYourCodeHere"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"username": "${var.name}",
					"password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand),
						"password": "inputYourCodeHere"}),
				),
			},
		},
	})

}

/*
func TestAccAlibabacloudStackAlikafkaSaslUser_multi(t *testing.T) {

	var v *alikafka.SaslUserVO
	resourceId := "alibabacloudstack_alikafka_sasl_user.default.1"
	ra := resourceAttrInit(resourceId, alikafkaSaslUserBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000,20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaSaslUserConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAlikafkaAclEnable(t)

			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "2",
					"instance_id": "${alibabacloudstack_alikafka_instance.default.id}",
					"username":    "${var.name}-${count.index}",
					"password":    "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v-1", rand),
						"password": "inputYourCodeHere",
					}),
				),
			},
		},
	})

}
*/
func resourceAlikafkaSaslUserConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

//data "alibabacloudstack_vpcs" "default" {
// name_regex = "^default-NODELETING"
//}
//data "alibabacloudstack_vswitches" "default" {
//  vpc_id = data.alibabacloudstack_vpcs.default.ids.0
//}
//
//resource "alibabacloudstack_security_group" "default" {
//  name   = var.name
//  vpc_id = data.alibabacloudstack_vpcs.default.ids.0
//}
//
//resource "alibabacloudstack_alikafka_instance" "default" {
//  name = "${var.name}"
//  topic_quota = "50"
//  disk_type = "1"
//  disk_size = "500"
//  deploy_type = "5"
//  io_max = "20"
//  vswitch_id = "${data.alibabacloudstack_vswitches.default.ids.0}"
//  security_group = alibabacloudstack_security_group.default.id
//}
`, name)
}

var alikafkaSaslUserBasicMap = map[string]string{
	"username": "${var.name}",
	"password": "inputYourCodeHere",
}
