package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_arms_alert_contact_group", &resource.Sweeper{
		Name: "alibabacloudstack_arms_alert_contact_group",
		F:    testSweepArmsAlertContactGroup,
	})
}

func testSweepArmsAlertContactGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "error getting AlibabacloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	action := "SearchAlertContactGroup"
	request := make(map[string]interface{})
	request["IsDetail"] = false
	request["RegionId"] = client.RegionId
	conn, err := client.NewArmsClient()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		log.Printf("[ERROR] %s failed error: %v", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.ContactGroups", response)
	if err != nil {
		log.Printf("[ERROR] %s error: %v", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["ContactGroupName"])
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping arms alert contact group: %s ", name)
			continue
		}
		log.Printf("[INFO] delete arms alert contact group: %s ", name)
		action = "DeleteAlertContactGroup"
		request = map[string]interface{}{
			"ContactGroupId": fmt.Sprint(item["ContactGroupId"]),
			"RegionId":       client.RegionId,
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			log.Printf("[ERROR] %s failed error: %v", action, err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackArmsAlertContactGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_arms_alert_contact_group.default"
	ra := resourceAttrInit(resourceId, ArmsAlertContactGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeArmsAlertContactGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsAlertContactGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ArmsAlertContactGroupBasicdependence)
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
					"alert_contact_group_name": "${var.name}",
					"contact_ids":              []string{"937", "938"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_group_name": name,
						"contact_ids.#":            "2",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"alert_contact_group_name": "${var.name}_update",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"alert_contact_group_name": name + "_update",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"contact_ids": []string{"937"},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"contact_ids.#": "1",
			//		}),
			//	),
			//},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var ArmsAlertContactGroupMap = map[string]string{}

func ArmsAlertContactGroupBasicdependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}
`, name)
}
