package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApiGatewayV2Consumer_Apikey(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_consumer.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAPIGatewayV2Consumer")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ConsumerDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_type":      "5",
					"app_name":       "${var.name}",
					"description":    "${var.name}",
					"key":            "a94231617c1a94cd900",
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"groups":         []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_type":      "5",
						"app_name":       name,
						"description":    name,
						"key":            "a94231617c1a94cd900",
						"gw_instance_id": CHECKSET,
						"groups.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "updated_description",
					"app_name":    "${var.name}_updated",
					"groups":      []string{"test", "updated_group"},
					"key":         "a94231617c1a94cd911",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "updated_description",
						"app_name":    name + "_updated",
						"groups.#":    "2",
						"key":         "a94231617c1a94cd911",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackApiGatewayV2Consumer_Basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_consumer.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAPIGatewayV2Consumer")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ConsumerDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_type":      "1",
					"app_name":       "${var.name}",
					"description":    "${var.name}",
					"key":            "root",
					"password":       "admin",
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"groups":         []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_type":      "1",
						"app_name":       name,
						"description":    name,
						"key":            "root",
						"password":       "admin",
						"gw_instance_id": CHECKSET,
						"groups.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "updated_description",
					"app_name":    "${var.name}_updated",
					"groups":      []string{"test", "updated_group"},
					"key":         "root1",
					"password":    "admin1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "updated_description",
						"app_name":    name + "_updated",
						"groups.#":    "2",
						"key":         "root1",
						"password":    "admin1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackApiGatewayV2Consumer_Jwt(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_consumer.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAPIGatewayV2Consumer")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ConsumerDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_type":   "3",
					"app_name":    "${var.name}",
					"description": "${var.name}",
					"payload": map[string]interface{}{
						"issuer":  "http://127.0.0.1:8000/test",
						"subject": "test.app",
					},
					"key":            "aaaaaaaaaaaaaaaaaa",
					"expire_time":    "86400000",
					"app_secret":     "bbbbbbbbbbbbbbbbbbbbb",
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"groups":         []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_type":      "3",
						"app_name":       name,
						"description":    name,
						"key":            "aaaaaaaaaaaaaaaaaa",
						"expire_time":    "86400000",
						"app_secret":     "bbbbbbbbbbbbbbbbbbbbb",
						"gw_instance_id": CHECKSET,
						"groups.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "updated_description",
					"app_name":    "${var.name}_updated",
					"groups":      []string{"test", "updated_group"},
					"payload": map[string]interface{}{
						"issuer":  "http://127.0.0.1:8000/test1",
						"subject": "test.app",
					},
					"key":         "aaaaaaaaaaaaaaaaaa1",
					"expire_time": REMOVEKEY,
					"app_secret":  "bbbbbbbbbbbbbbbbbbbbb1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "updated_description",
						"app_name":    name + "_updated",
						"groups.#":    "2",
						"key":         "aaaaaaaaaaaaaaaaaa1",
						"expire_time": REMOVEKEY,
						"app_secret":  "bbbbbbbbbbbbbbbbbbbbb1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"payload"},
			},
		},
	})
}

func TestAccAlibabacloudStackApiGatewayV2Consumer_Oauth2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_consumer.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAPIGatewayV2Consumer")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ConsumerDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_type":   "2",
					"app_name":    "${var.name}",
					"description": "${var.name}",
					"oauth2_payload": []map[string]interface{}{
						{
							"authorization_code":       true,
							"client_credentials":       false,
							"implicit_grant":           false,
							"password_grant":           false,
							"token_expiration":         7200,
							"refresh_token_expiration": 14,
							"pkce":                     true,
							"scopes":                   "aaa,bbb",
							"client_id":                "testclientid111",
							"client_secret":            "testclientsecret111",
							"redirect_uris":            "http://127.0.0.1:8000,https://127.0.0.1:7999",
						},
					},
					"key":            "aaaaaaaaaaaaaaaaaa",
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"groups":         []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_type":                                 "2",
						"key":                                       "aaaaaaaaaaaaaaaaaa",
						"app_name":                                  name,
						"description":                               name,
						"oauth2_payload.#":                          "1",
						"oauth2_payload.0.authorization_code":       "true",
						"oauth2_payload.0.client_credentials":       "false",
						"oauth2_payload.0.implicit_grant":           "false",
						"oauth2_payload.0.password_grant":           "false",
						"oauth2_payload.0.token_expiration":         "7200",
						"oauth2_payload.0.refresh_token_expiration": "14",
						"oauth2_payload.0.pkce":                     "true",
						"oauth2_payload.0.scopes":                   "aaa,bbb",
						"oauth2_payload.0.client_id":                "testclientid111",
						"oauth2_payload.0.client_secret":            "testclientsecret111",
						"oauth2_payload.0.redirect_uris":            "http://127.0.0.1:8000,https://127.0.0.1:7999",
						"gw_instance_id":                            CHECKSET,
						"groups.#":                                  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "updated_description",
					"app_name":    "${var.name}_updated",
					"groups":      []string{"test", "updated_group"},
					"key":         "aaaaaaaaaaaaaaaaaa1",
					"oauth2_payload": []map[string]interface{}{
						{
							"authorization_code":       true,
							"client_credentials":       false,
							"implicit_grant":           false,
							"password_grant":           false,
							"token_expiration":         3600,
							"refresh_token_expiration": 28,
							"pkce":                     false,
							"scopes":                   "aaa",
							"client_id":                "testclientid1222",
							"client_secret":            "testclientsecret1222",
							"redirect_uris":            "http://127.0.0.1:8001,https://127.0.0.1:7992",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                               "updated_description",
						"app_name":                                  name + "_updated",
						"groups.#":                                  "2",
						"key":                                       "aaaaaaaaaaaaaaaaaa1",
						"oauth2_payload.#":                          "1",
						"oauth2_payload.0.authorization_code":       "true",
						"oauth2_payload.0.client_credentials":       "false",
						"oauth2_payload.0.implicit_grant":           "false",
						"oauth2_payload.0.password_grant":           "false",
						"oauth2_payload.0.token_expiration":         "3600",
						"oauth2_payload.0.refresh_token_expiration": "28",
						"oauth2_payload.0.pkce":                     "false",
						"oauth2_payload.0.scopes":                   "aaa",
						"oauth2_payload.0.client_id":                "testclientid1222",
						"oauth2_payload.0.client_secret":            "testclientsecret1222",
						"oauth2_payload.0.redirect_uris":            "http://127.0.0.1:8001,https://127.0.0.1:7992",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackApiGatewayV2Consumer_Appgw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_consumer.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAPIGatewayV2Consumer")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ConsumerDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_type":      "6",
					"app_name":       "${var.name}",
					"description":    "${var.name}",
					"app_secret":     "vvvvvvvvvvvvvvvvvvvvv",
					"app_code":       "aaaaaaaaaaa",
					"key":            "aaaaaaaaaaaaaaaaaa",
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"groups":         []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_type":      "6",
						"app_name":       name,
						"description":    name,
						"key":            "aaaaaaaaaaaaaaaaaa",
						"app_secret":     "vvvvvvvvvvvvvvvvvvvv1",
						"app_code":       "aaaaaaaaaaa1",
						"gw_instance_id": CHECKSET,
						"groups.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "updated_description",
					"app_name":    "${var.name}_updated",
					"groups":      []string{"test", "updated_group"},
					"key":         "aaaaaaaaaaaaaaaaaa1",
					"app_secret":  "vvvvvvvvvvvvvvvvvvvv1",
					"app_code":    "aaaaaaaaaaa1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "updated_description",
						"app_name":    name + "_updated",
						"groups.#":    "2",
						"key":         "aaaaaaaaaaaaaaaaaa1",
						"app_secret":  "vvvvvvvvvvvvvvvvvvvv1",
						"app_code":    "aaaaaaaaaaa1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackApiGatewayV2Consumer_Csb(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_consumer.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAPIGatewayV2Consumer")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ConsumerDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_type":      "7",
					"app_name":       "${var.name}",
					"description":    "${var.name}",
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"groups":         []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_type":      "7",
						"app_name":       name,
						"description":    name,
						"gw_instance_id": CHECKSET,
						"groups.#":       "1",
						"access_key":     CHECKSET,
						"secret_key":     CHECKSET,
						"app_secret":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "updated_description",
					"app_name":    "${var.name}_updated",
					"groups":      []string{"test", "updated_group"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "updated_description",
						"app_name":    name + "_updated",
						"groups.#":    "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func ApiGatewayV2ConsumerDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  	instance_name = "${var.name}"
	node_number = "1"
	instance_class = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
	broker_engine_type = "SCG"
	deploy_mode = "custom"
}

`, name)
}
