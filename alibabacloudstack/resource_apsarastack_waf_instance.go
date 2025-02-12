package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlibabacloudstackWafInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudstackWafInstanceCreate,
		Read:   resourceAlibabacloudstackWafInstanceRead,
		Update: resourceAlibabacloudstackWafInstanceUpdate,
		Delete: resourceAlibabacloudstackWafInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"big_screen": {
				Type:     schema.TypeString,
				Required: true,
			},
			"exclusive_ip_package": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_domain_package": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_storage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modify_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"prefessional_service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"waf_log": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudstackWafInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var err error
	action := "CreateWAFInstance"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}

	request["ProductCode"] = "waf"
	request["ProductType"] = "waf"
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	}

	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	region := client.RegionId
	if v, ok := d.GetOk("region"); ok && v.(string) != "" {
		region = v.(string)
	}
	request["SubscriptionType"] = d.Get("subscription_type")
	request["Parameter"] = []map[string]string{
		{
			"Code":  "BigScreen",
			"Value": d.Get("big_screen").(string),
		},
		{
			"Code":  "ExclusiveIpPackage",
			"Value": d.Get("exclusive_ip_package").(string),
		},
		{
			"Code":  "ExtBandwidth",
			"Value": d.Get("ext_bandwidth").(string),
		},
		{
			"Code":  "ExtDomainPackage",
			"Value": d.Get("ext_domain_package").(string),
		},
		{
			"Code":  "LogStorage",
			"Value": d.Get("log_storage").(string),
		},
		{
			"Code":  "LogTime",
			"Value": d.Get("log_time").(string),
		},
		{
			"Code":  "PackageCode",
			"Value": d.Get("package_code").(string),
		},
		{
			"Code":  "PrefessionalService",
			"Value": d.Get("prefessional_service").(string),
		},
		{
			"Code":  "Region",
			"Value": region,
		},
		{
			"Code":  "WafLog",
			"Value": d.Get("waf_log").(string),
		},
	}
	response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", action, "", nil, request)
	if err != nil {
		return err
	}
	addDebug(action, response, request)
	response = response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["InstanceId"]))

	return resourceAlibabacloudstackWafInstanceUpdate(d, meta)
}
func resourceAlibabacloudstackWafInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	waf_openapiService := WafOpenapiService{client}
	object, err := waf_openapiService.DescribeWafInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_waf_instance waf_openapiService.DescribeWafInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("status", formatInt(object["Status"]))
	d.Set("subscription_type", object["SubscriptionType"])
	return nil
}
func resourceAlibabacloudstackWafInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("subscription_type") {
		update = true
	}
	request["ProductType"] = "waf"
	request["SubscriptionType"] = d.Get("subscription_type")
	request["ProductCode"] = "waf"
	request["ModifyType"] = d.Get("modify_type")
	request["Parameter"] = []map[string]string{
		{
			"Code":  "BigScreen",
			"Value": d.Get("big_screen").(string),
		},
		{
			"Code":  "ExclusiveIpPackage",
			"Value": d.Get("exclusive_ip_package").(string),
		},
		{
			"Code":  "ExtBandwidth",
			"Value": d.Get("ext_bandwidth").(string),
		},
		{
			"Code":  "ExtDomainPackage",
			"Value": d.Get("ext_domain_package").(string),
		},
		{
			"Code":  "LogStorage",
			"Value": d.Get("log_storage").(string),
		},
		{
			"Code":  "LogTime",
			"Value": d.Get("log_time").(string),
		},
		{
			"Code":  "PackageCode",
			"Value": d.Get("package_code").(string),
		},
		{
			"Code":  "PrefessionalService",
			"Value": d.Get("prefessional_service").(string),
		},
		{
			"Code":  "WafLog",
			"Value": d.Get("waf_log").(string),
		},
	}
	if update {
		action := "EditWAFInstance"
		response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", action, "", nil, request)
		if err != nil {
			return err
		}
		addDebug(action, response, request)
	}
	return resourceAlibabacloudstackWafInstanceRead(d, meta)
}
func resourceAlibabacloudstackWafInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteWAFInstance"
	var err error
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", action, "", nil, request)
	if err != nil {
		return err
	}
	addDebug(action, response, request)
	return nil
}
