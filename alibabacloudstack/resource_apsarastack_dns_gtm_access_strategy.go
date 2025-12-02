package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDnsGtmAccessStrategy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gtm_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_gtm_address_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_gtm_address_pool_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_min_available_addr_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"failover_gtm_address_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"failover_gtm_address_pool_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"DOMAIN", "IPV6", "IPV4"}, false),
			},
			"failover_min_available_addr_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"switch_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"BY_PROBE_RESULT", "BY_HAND"}, false),
			},
			"line_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"specified_gtm_address_pool": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "FAILOVER"}, false),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool{
					if d.Get("switch_mode").(string) == "BY_PROBE_RESULT" {
						return true
					}
					return oldValue == newValue
				},
			},
			"default_available_addr_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failover_available_addr_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_gtm_address_pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failover_gtm_address_pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"in_use_gtm_address_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"in_use_gtm_address_pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDnsGtmAccessStrategyCreate, resourceAlibabacloudStackDnsGtmAccessStrategyRead, resourceAlibabacloudStackDnsGtmAccessStrategyUpdate, resourceAlibabacloudStackDnsGtmAccessStrategyDelete)
	return resource
}

func resourceAlibabacloudStackDnsGtmAccessStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	gtm_instance_id := d.Get("gtm_instance_id").(string)
	request["GtmInstanceId"] = gtm_instance_id
	request["Name"] = d.Get("name")
	request["DefaultGtmAddressPoolId"] = d.Get("default_gtm_address_pool_id")
	request["DefaultGtmAddressPoolType"] = d.Get("default_gtm_address_pool_type")
	request["DefaultMinAvailableAddrNum"] = d.Get("default_min_available_addr_num")
	request["SwitchMode"] = d.Get("switch_mode")

	if v, ok := d.GetOk("failover_gtm_address_pool_id"); ok && v.(string) != "" {
		request["FailoverGtmAddressPoolId"] = v
	}

	if v, ok := d.GetOk("failover_gtm_address_pool_type"); ok && v.(string) != "" {
		request["FailoverGtmAddressPoolType"] = v
	}

	if v, ok := d.GetOk("failover_min_available_addr_num"); ok {
		request["FailoverMinAvailableAddrNum"] = v
	}

	if v, ok := d.GetOk("specified_gtm_address_pool"); ok && v.(string) != "" {
		request["SpecifiedGtmAddressPool"] = v
	}

	lineIds := d.Get("line_ids").(*schema.Set).List()
	if len(lineIds) > 0 {
		linesJson, _ := json.Marshal(lineIds)
		request["LineIds"] = string(linesJson)
	}

	action := "AddDnsGtmAccessStrategy"
	response, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", action, "", nil, nil, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	id, ok := response["Id"].(string)
	if !ok || id == "" {
		return errmsgs.WrapError(fmt.Errorf("failed to get Id from response"))
	}
	resourceId := fmt.Sprintf("%s:%s", gtm_instance_id, id)
	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackDnsGtmAccessStrategyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}

	object, err := dnsService.DescribeDnsGtmAccessStrategy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", object["Name"])
	d.Set("gtm_instance_id", object["GtmInstanceId"])
	d.Set("default_gtm_address_pool_id", object["DefaultGtmAddressPoolId"])
	d.Set("default_gtm_address_pool_type", object["DefaultGtmAddressPoolType"])
	d.Set("default_min_available_addr_num", object["DefaultMinAvailableAddrNum"])
	d.Set("failover_gtm_address_pool_id", object["FailoverGtmAddressPoolId"])
	d.Set("failover_gtm_address_pool_type", object["FailoverGtmAddressPoolType"])
	d.Set("failover_min_available_addr_num", object["FailoverMinAvailableAddrNum"])
	d.Set("switch_mode", object["SwitchMode"])
	d.Set("specified_gtm_address_pool", object["SpecifiedGtmAddressPool"])
	d.Set("default_available_addr_num", object["DefaultAvailableAddrNum"])
	d.Set("failover_available_addr_num", object["FailoverAvailableAddrNum"])
	d.Set("default_gtm_address_pool_name", object["DefaultGtmAddressPoolName"])
	d.Set("failover_gtm_address_pool_name", object["FailoverGtmAddressPoolName"])
	d.Set("in_use_gtm_address_pool_id", object["InUseGtmAddressPoolId"])
	d.Set("in_use_gtm_address_pool_name", object["InUseGtmAddressPoolName"])

	if lineIds, ok := object["LineIds"].([]interface{}); ok {
		lineIdStrings := make([]string, len(lineIds))
		for i, v := range lineIds {
			lineIdStrings[i] = v.(string)
		}
		d.Set("line_ids", lineIdStrings)
	}

	return nil
}

func resourceAlibabacloudStackDnsGtmAccessStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("name", "line_ids", "default_gtm_address_pool_type", "default_gtm_address_pool_id",
		"default_min_available_addr_num", "failover_gtm_address_pool_id", "failover_gtm_address_pool_type",
		"failover_min_available_addr_num", "switch_mode", "specified_gtm_address_pool") {
		params := strings.Split(d.Id(), ":")
		reqBody := make(map[string]interface{})
		reqBody["Id"] = params[1]
		reqBody["GtmInstanceId"] = params[0]
		reqBody["Name"] = d.Get("name")
		reqBody["DefaultGtmAddressPoolType"] = d.Get("default_gtm_address_pool_type")
		reqBody["DefaultGtmAddressPoolId"] = d.Get("default_gtm_address_pool_id")
		reqBody["DefaultMinAvailableAddrNum"] = d.Get("default_min_available_addr_num")
		reqBody["SwitchMode"] = d.Get("switch_mode")

		lineIds := d.Get("line_ids").(*schema.Set).List()
		if len(lineIds) > 0 {
			linesJson, _ := json.Marshal(lineIds)
			reqBody["LineIds"] = string(linesJson)
		}

		if v, ok := d.GetOk("failover_gtm_address_pool_id"); ok && v.(string) != "" {
			reqBody["FailoverGtmAddressPoolId"] = v
		}

		if v, ok := d.GetOk("failover_gtm_address_pool_type"); ok && v.(string) != "" {
			reqBody["FailoverGtmAddressPoolType"] = v
		}

		if v, ok := d.GetOk("failover_min_available_addr_num"); ok {
			reqBody["FailoverMinAvailableAddrNum"] = v
		}

		if v, ok := d.GetOk("specified_gtm_address_pool"); ok && v.(string) != "" {
			reqBody["SpecifiedGtmAddressPool"] = v
		}

		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "UpdateDnsGtmAccessStrategy", "", nil, nil, reqBody)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_gtm_access_strategy", "UpdateDnsGtmAccessStrategy", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return resourceAlibabacloudStackDnsGtmAccessStrategyRead(d, meta)
}

func resourceAlibabacloudStackDnsGtmAccessStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	params := strings.Split(d.Id(), ":")
	reqBody := map[string]interface{}{
		"Id":            params[1],
		"GtmInstanceId": params[0],
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DeleteDnsGtmAccessStrategy", "", nil, nil, reqBody)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteDnsGtmAccessStrategy", errmsgs.AlibabacloudStackSdkGoERROR, "")
			return resource.RetryableError(err)
		}
		addDebug("DeleteDnsGtmAccessStrategy", raw)
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
