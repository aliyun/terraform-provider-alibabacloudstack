package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackBastionhostInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackBastionhostInstanceCreate,
		Read:   resourceAlibabacloudStackBastionhostInstanceRead,
		Update: resourceAlibabacloudStackBastionhostInstanceUpdate,
		Delete: resourceAlibabacloudStackBastionhostInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"license_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"asset": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"highavailability": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"disasterrecovery": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_public_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ad_auth_server": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					d1, d2 := d.GetChange("ad_auth_server")
					if len(d1.(*schema.Set).List()) == 0 || len(d2.(*schema.Set).List()) == 0 {
						return false
					}
					return compareMapWithIgnoreEquivalent(d1.(*schema.Set).List()[0].(map[string]interface{}), d2.(*schema.Set).List()[0].(map[string]interface{}), []string{"password"})
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account": {
							Type:     schema.TypeString,
							Required: true,
						},
						"base_dn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email_mapping": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"filter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_ssl": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"mobile_mapping": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name_mapping": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"server": {
							Type:     schema.TypeString,
							Required: true,
						},
						"standby_server": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ldap_auth_server": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					d1, d2 := d.GetChange("ldap_auth_server")
					if len(d1.(*schema.Set).List()) == 0 || len(d2.(*schema.Set).List()) == 0 {
						return false
					}
					return compareMapWithIgnoreEquivalent(d1.(*schema.Set).List()[0].(map[string]interface{}), d2.(*schema.Set).List()[0].(map[string]interface{}), []string{"password"})
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account": {
							Type:     schema.TypeString,
							Required: true,
						},
						"base_dn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email_mapping": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"filter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"login_name_mapping": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mobile_mapping": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name_mapping": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"server": {
							Type:     schema.TypeString,
							Required: true,
						},
						"standby_server": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.Any(validation.IntBetween(1, 9), validation.IntInSlice([]int{12, 24, 36})),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
			},
			"renewal_period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"M", "Y"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AutoRenewal", "ManualRenewal", "NotRenewal"}, false),
			},
		},
	}
}

func resourceAlibabacloudStackBastionhostInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	parameterMapList := make([]map[string]interface{}, 0)
	// conn, err := client.NewBastionhostClient()
	// if err != nil {
	// 	return errmsgs.WrapError(err)
	// }
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "NetworkType",
		"Value": "vpc",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "LicenseCode",
		"Value": d.Get("license_code").(string),
	})
	// parameterMapList = append(parameterMapList, map[string]interface{}{
	// 	"Code":  "PlanCode",
	// 	"Value": d.Get("plan_code").(string),
	// })
	// parameterMapList = append(parameterMapList, map[string]interface{}{
	// 	"Code":  "Storage",
	// 	"Value": d.Get("storage").(string),
	// })
	// parameterMapList = append(parameterMapList, map[string]interface{}{
	// 	"Code":  "Bandwidth",
	// 	"Value": d.Get("bandwidth").(string),
	// })
	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VswitchId"] = v
	}
	if v, ok := d.GetOk("asset"); ok {
		request["Asset"] = v
	}
	if v, ok := d.GetOk("highavailability"); ok {
		request["HighAvailability"] = v
	}
	if v, ok := d.GetOk("disasterrecovery"); ok {
		request["DisasterRecovery"] = v
	}

	if v, ok := d.GetOk("license_code"); ok {
		request["LicenseCode"] = v
	}

	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return errmsgs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renew_period", "renewal_status", d.Get("renewal_status")))
	}
	request["ProductCode"] = "bastionhost"
	request["ProductType"] = "bastionhost"
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "RegionId",
		"Value": client.RegionId,
	})
	request["Parameter"] = parameterMapList
	request["ClientToken"] = buildClientToken("CreateInstance")
	// response, err := client.DoTeaRequest("POST", "Bastionhostprivate", "2023-03-23", action, "", nil, request)
	response, err := client.DoTeaRequest("POST", "Bastionhostprivate", "2023-03-23", action, "", nil, request)
	addDebug(action, response, request)
	if err != nil {
		return err
	}
	// runtime := util.RuntimeOptions{}
	// runtime.SetAutoretry(true)
	// wait := incrementalWait(3*time.Second, 3*time.Second)
	// err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	// 	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
	// 	if err != nil {
	// 		if errmsgs.NeedRetry(err) {
	// 			wait()
	// 			return resource.RetryableError(err)
	// 		}
	// 		if errmsgs.IsExpectedErrors(err, []string{"NotApplicable"}) {
	// 			request["ProductType"] = "bastionhost_std_public_intl"
	// 			conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
	// 			return resource.RetryableError(err)
	// 		}
	// 		return resource.NonRetryableError(err)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudctack_bastionhost_instance", action, errmsgs.AlibabacloudStackSdkGoERROR)
	// }
	// if fmt.Sprint(response["Code"]) != "Success" {
	// 	return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	// }
	// responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["InstanceId"]))

	bastionhostService := YundunBastionhostService{client}

	// check RAM policy
	// if err := bastionhostService.ProcessRolePolicy(); err != nil {
	// 	return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	// }
	// wait for order complete
	stateConf := BuildStateConf([]string{}, []string{"PENDING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	rawSecurityGroupIds := d.Get("security_group_ids").(*schema.Set).List()
	securityGroupIds := make([]string, len(rawSecurityGroupIds))
	for index, rawSecurityGroupId := range rawSecurityGroupIds {
		securityGroupIds[index] = rawSecurityGroupId.(string)
	}
	// start instance
	if err := bastionhostService.StartBastionhostInstance(d.Id(), d.Get("vswitch_id").(string), securityGroupIds); err != nil {
		return errmsgs.WrapError(err)
	}
	// wait for pending
	stateConf = BuildStateConf([]string{"PENDING", "CREATING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 600*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return resourceAlibabacloudStackBastionhostInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackBastionhostInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	BastionhostService := YundunBastionhostService{client}
	instance, err := BastionhostService.DescribeBastionhostInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("description", instance["Description"])
	d.Set("license_code", instance["LicenseCode"])
	d.Set("vswitch_id", instance["VswitchId"])
	d.Set("security_group_ids", instance["AuthorizedSecurityGroups"])
	d.Set("enable_public_access", instance["PublicNetworkAccess"])
	d.Set("resource_group_id", instance["ResourceGroupId"])

	if fmt.Sprint(instance["PublicNetworkAccess"]) == "true" {
		d.Set("public_white_list", instance["PublicWhiteList"])
	}

	instance, err = BastionhostService.DescribeBastionhostInstances(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("plan_code", instance["PlanCode"])

	tags, err := BastionhostService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("tags", BastionhostService.tagsToMap(tags))

	adAuthServer, err := BastionhostService.DescribeBastionhostAdAuthServer(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	adAuthServerMap := map[string]interface{}{
		"account":        adAuthServer["Account"],
		"base_dn":        adAuthServer["BaseDN"],
		"domain":         adAuthServer["Domain"],
		"email_mapping":  adAuthServer["EmailMapping"],
		"filter":         adAuthServer["Filter"],
		"is_ssl":         adAuthServer["IsSSL"],
		"mobile_mapping": adAuthServer["MobileMapping"],
		"name_mapping":   adAuthServer["NameMapping"],
		"port":           formatInt(adAuthServer["Port"]),
		"server":         adAuthServer["Server"],
		"standby_server": adAuthServer["StandbyServer"],
	}
	d.Set("ad_auth_server", []map[string]interface{}{adAuthServerMap})

	ldapAuthServer, err := BastionhostService.DescribeBastionhostLdapAuthServer(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	ldapAuthServerMap := map[string]interface{}{
		"account":            ldapAuthServer["Account"],
		"base_dn":            ldapAuthServer["BaseDN"],
		"email_mapping":      ldapAuthServer["EmailMapping"],
		"filter":             ldapAuthServer["Filter"],
		"is_ssl":             ldapAuthServer["IsSSL"],
		"login_name_mapping": ldapAuthServer["LoginNameMapping"],
		"mobile_mapping":     ldapAuthServer["MobileMapping"],
		"name_mapping":       ldapAuthServer["NameMapping"],
		"port":               formatInt(ldapAuthServer["Port"]),
		"server":             ldapAuthServer["Server"],
		"standby_server":     ldapAuthServer["StandbyServer"],
	}
	d.Set("ldap_auth_server", []map[string]interface{}{ldapAuthServerMap})

	// can not set region when invoking QueryAvailableInstances for bastionhost instance

	d.Set("renewal_status", instance["RenewStatus"])
	d.Set("renew_period", formatInt(instance["RenewalDuration"]))
	d.Set("renewal_period_unit", instance["RenewalDurationUnit"])

	return nil
}

func resourceAlibabacloudStackBastionhostInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	bastionhostService := YundunBastionhostService{client}

	d.Partial(true)

	if d.HasChange("tags") {
		if err := bastionhostService.setInstanceTags(d, TagResourceInstance); err != nil {
			return errmsgs.WrapError(err)
		}

	}

	if d.HasChange("description") {
		if err := bastionhostService.UpdateBastionhostInstanceDescription(d.Id(), d.Get("description").(string)); err != nil {
			return errmsgs.WrapError(err)
		}

	}

	if d.HasChange("resource_group_id") {
		if err := bastionhostService.UpdateResourceGroup(d.Id(), d.Get("resource_group_id").(string)); err != nil {
			return errmsgs.WrapError(err)
		}

	}

	if !d.IsNewResource() && d.HasChange("license_code") {
		params := map[string]string{
			"LicenseCode": "license_code",
		}
		if err := bastionhostService.UpdateInstanceSpec(params, d, meta); err != nil {
			return errmsgs.WrapError(err)
		}
		stateConf := BuildStateConf([]string{"UPGRADING"}, []string{"PENDING", "RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"CREATING", "UPGRADE_FAILED", "CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

	}

	if !d.IsNewResource() && d.HasChange("security_group_ids") {
		securityGroupIds := d.Get("security_group_ids").(*schema.Set).List()
		sgs := make([]string, 0, len(securityGroupIds))
		for _, rawSecurityGroupId := range securityGroupIds {
			sgs = append(sgs, rawSecurityGroupId.(string))
		}
		if err := bastionhostService.UpdateBastionhostSecurityGroups(d.Id(), sgs); err != nil {
			return errmsgs.WrapError(err)
		}
		stateConf := BuildStateConf([]string{"UPGRADING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"CREATING", "UPGRADE_FAILED", "CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	_, ok := d.GetOkExists("enable_public_access")
	if d.HasChange("enable_public_access") || (d.IsNewResource() && ok) {
		client := meta.(*connectivity.AlibabacloudStackClient)
		BastionhostService := YundunBastionhostService{client}
		instance, err := BastionhostService.DescribeBastionhostInstance(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		target := strconv.FormatBool(d.Get("enable_public_access").(bool))
		if strconv.FormatBool(instance["PublicNetworkAccess"].(bool)) != target {
			if target == "false" {
				err := BastionhostService.DisableInstancePublicAccess(d.Id())
				if err != nil {
					return errmsgs.WrapError(err)
				}
			} else {
				err := BastionhostService.EnableInstancePublicAccess(d.Id())
				if err != nil {
					return errmsgs.WrapError(err)
				}
			}
		}

		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

	}

	if d.HasChange("ad_auth_server") {
		if v, ok := d.GetOk("ad_auth_server"); ok && len(v.(*schema.Set).List()) > 0 {
			var response map[string]interface{}
			modifyAdRequest := map[string]interface{}{
				"InstanceId": d.Id(),
				"RegionId":   client.RegionId,
			}
			adAuthServer := v.(*schema.Set).List()[0].(map[string]interface{})
			modifyAdRequest["Account"] = adAuthServer["account"]
			modifyAdRequest["BaseDN"] = adAuthServer["base_dn"]
			modifyAdRequest["Domain"] = adAuthServer["domain"]
			modifyAdRequest["IsSSL"] = adAuthServer["is_ssl"]
			modifyAdRequest["Port"] = adAuthServer["port"]
			modifyAdRequest["Server"] = adAuthServer["server"]
			modifyAdRequest["EmailMapping"] = adAuthServer["email_mapping"]
			modifyAdRequest["Filter"] = adAuthServer["filter"]
			modifyAdRequest["MobileMapping"] = adAuthServer["mobile_mapping"]
			modifyAdRequest["NameMapping"] = adAuthServer["name_mapping"]
			modifyAdRequest["Password"] = adAuthServer["password"]
			modifyAdRequest["StandbyServer"] = adAuthServer["standby_server"]

			action := "ModifyInstanceADAuthServer"
			_, err := client.DoTeaRequest("POST", "Bastionhostprivate", "2019-12-09", action, "", nil, modifyAdRequest)
			if err != nil {
				return err
			}
			addDebug(action, response, modifyAdRequest)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
			}

		}
	}

	if d.HasChange("ldap_auth_server") {
		if v, ok := d.GetOk("ldap_auth_server"); ok && len(v.(*schema.Set).List()) > 0 {
			var response map[string]interface{}
			modifyLdapRequest := map[string]interface{}{
				"InstanceId": d.Id(),
				"RegionId":   client.RegionId,
			}

			adAuthServer := v.(*schema.Set).List()[0].(map[string]interface{})
			modifyLdapRequest["Account"] = adAuthServer["account"]
			modifyLdapRequest["BaseDN"] = adAuthServer["base_dn"]
			modifyLdapRequest["Port"] = adAuthServer["port"]
			modifyLdapRequest["Server"] = adAuthServer["server"]
			modifyLdapRequest["Password"] = adAuthServer["password"]
			modifyLdapRequest["IsSSL"] = adAuthServer["is_ssl"]
			modifyLdapRequest["LoginNameMapping"] = adAuthServer["login_name_mapping"]
			modifyLdapRequest["EmailMapping"] = adAuthServer["email_mapping"]
			modifyLdapRequest["Filter"] = adAuthServer["filter"]
			modifyLdapRequest["MobileMapping"] = adAuthServer["mobile_mapping"]
			modifyLdapRequest["NameMapping"] = adAuthServer["name_mapping"]
			modifyLdapRequest["StandbyServer"] = adAuthServer["standby_server"]

			action := "ModifyInstanceLDAPAuthServer"
			// wait := incrementalWait(3*time.Second, 3*time.Second)
			// err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			// 	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, modifyLdapRequest, &util.RuntimeOptions{})
			// 	if err != nil {
			// 		if NeedRetry(err) {
			// 			wait()
			// 			return resource.RetryableError(err)
			// 		}
			// 		return resource.NonRetryableError(err)
			// 	}
			// 	return nil
			// })
			_, err := client.DoTeaRequest("POST", "Bastionhostprivate", "2019-12-09", action, "", nil, modifyLdapRequest)
			if err != nil {
				return err
			}
			addDebug(action, response, modifyLdapRequest)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
			}

		}
	}

	var setRenewalResponse map[string]interface{}
	update := false
	setRenewalReq := map[string]interface{}{
		"InstanceIDs":      d.Id(),
		"ProductCode":      "bastionhost",
		"ProductType":      "bastionhost",
		"SubscriptionType": "Subscription",
	}

	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		setRenewalReq["RenewalStatus"] = v
	}

	if !d.IsNewResource() && d.HasChange("renew_period") {
		update = true
		if v, ok := d.GetOk("renew_period"); ok {
			setRenewalReq["RenewalPeriod"] = v
		}
	}

	if d.HasChange("renewal_period_unit") {
		update = true
	}
	if v, ok := d.GetOk("renewal_period_unit"); ok {
		setRenewalReq["RenewalPeriodUnit"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return errmsgs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_period_unit", "renewal_status", d.Get("renewal_status")))
	}

	if update {
		request := map[string]interface{}{}
		action := "SetRenewal"
		_, err := client.DoTeaRequest("POST", "Bastionhostprivate", "2017-12-14", action, "", nil, request)
		if err != nil {
			return err
		}
		// conn, err := client.NewBssopenapiClient()
		// if err != nil {
		// 	return errmsgs.WrapError(err)
		// }
		// wait := incrementalWait(3*time.Second, 3*time.Second)
		// err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		// 	setRenewalResponse, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, setRenewalReq, &util.RuntimeOptions{})
		// 	if err != nil {
		// 		if NeedRetry(err) {
		// 			wait()
		// 			return resource.RetryableError(err)
		// 		}
		// 		if IsExpectedErrors(err, []string{"NotApplicable"}) {
		// 			conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
		// 			setRenewalReq["ProductType"] = "bastionhost_std_public_intl"
		// 			return resource.RetryableError(err)
		// 		}
		// 		return resource.NonRetryableError(err)
		// 	}
		// 	return nil
		// })
		addDebug(action, setRenewalResponse, setRenewalReq)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		if fmt.Sprint(setRenewalResponse["Code"]) != "Success" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, setRenewalResponse))
		}

	}

	update = false
	configInstanceWhiteListReq := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	if d.HasChange("public_white_list") {
		update = true
	}
	if v, ok := d.GetOk("public_white_list"); ok {
		configInstanceWhiteListReq["WhiteList"] = v
	}

	if update {
		request := map[string]interface{}{}
		action := "ConfigInstanceWhiteList"
		_, err := client.DoTeaRequest("POST", "Bastionhostprivate", "2019-12-09", action, "", nil, request)
		if err != nil {
			return err
		}
		// wait := incrementalWait(3*time.Second, 3*time.Second)
		// err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		// 	resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, configInstanceWhiteListReq, &util.RuntimeOptions{})
		// 	if err != nil {
		// 		if NeedRetry(err) {
		// 			wait()
		// 			return resource.RetryableError(err)
		// 		}
		// 		return resource.NonRetryableError(err)
		// 	}
		// 	addDebug(action, resp, configInstanceWhiteListReq)
		// 	return nil
		// })

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

	}

	d.Partial(false)
	// wait for order complete
	return resourceAlibabacloudStackBastionhostInstanceRead(d, meta)
}

func resourceAlibabacloudStackBastionhostInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceBastionhostInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
