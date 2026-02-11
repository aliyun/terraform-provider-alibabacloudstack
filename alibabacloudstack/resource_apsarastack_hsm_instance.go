package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackHsmInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vendor_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsm_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"evsm", "gvsm", "svsm"}, false),
			},
			"zone_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hsm_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"white_list": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_master": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"show_create_cluster": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"release_protection": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackHsmInstanceCreate, resourceAlibabacloudStackHsmInstanceRead, resourceAlibabacloudStackHsmInstanceUpdate, resourceAlibabacloudStackHsmInstanceDelete)
	return resource
}

func resourceAlibabacloudStackHsmInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"ProductCode":     d.Get("product_code"),
		"VendorCode":      d.Get("vendor_code"),
		"VsmType":         d.Get("vsm_type"),
		"VsmNum":          1,
		"ZoneNo":          d.Get("zone_no"),
		"IsAppointDevice": false,
	}
	if v, ok := d.GetOk("hsm_id"); ok {
		request["HsmType"] = v.(string)
		request["IsAppointDevice"] = true
	}

	response, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "CreateInstance", "", nil, request, nil)
	if err != nil {
		return err
	}

	instanceId, ok := response["InstanceId"].(string)
	if !ok || instanceId == "" {
		return fmt.Errorf("failed to get instance id from response")
	}

	d.SetId(instanceId)

	return nil
}

func resourceAlibabacloudStackHsmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hsmService := HsmService{client}
	instance, err := hsmService.DescribeHsmInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Set all computed fields from the response
	d.Set("instance_id", instance["InstanceId"])
	d.Set("product_code", instance["ProductCode"])
	d.Set("vendor_code", instance["VendorCode"])
	d.Set("vsm_type", instance["VsmType"])
	d.Set("zone_no", instance["ZoneNo"])
	d.Set("hsm_id", instance["HsmId"])
	d.Set("vpc_id", instance["VpcId"])
	d.Set("vswitch_id", instance["VswitchId"])
	d.Set("ip", instance["Ip"])
	d.Set("remark", instance["Remark"])
	d.Set("product_name", instance["ProductName"])
	d.Set("vendor_name", instance["VendorName"])
	d.Set("cluster_id", instance["ClusterId"])
	d.Set("cluster_name", instance["ClusterName"])
	d.Set("is_master", instance["IsMaster"])
	d.Set("show_create_cluster", instance["ShowCreateCluster"])
	d.Set("release_protection", instance["ReleaseProtection"])
	d.Set("status", instance["HsmStatus"])

	// Handle WhiteList field if present
	if whiteList, ok := instance["WhiteList"].([]interface{}); ok && len(whiteList) > 0 {
		whiteListStr := make([]string, len(whiteList))
		for i, item := range whiteList {
			whiteListStr[i] = fmt.Sprintf("%v", item)
		}
		d.Set("white_list", strings.Join(whiteListStr, ","))
	} else {
		d.Set("white_list", "")
	}

	return nil
}

func resourceAlibabacloudStackHsmInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChange("remark") {
		reqQuery := map[string]interface{}{
			"InstanceId": d.Id(),
			"Remark":     d.Get("remark"),
		}

		_, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "ModifyInstance", "", nil, reqQuery, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_hsm_instance", "ModifyInstance", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	if d.HasChanges("vpc_id", "vswitch_id", "ip", "white_list") {
		reqQuery := map[string]interface{}{
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("vpc_id"); ok && v != "" {
			reqQuery["VpcId"] = v
		}
		if v, ok := d.GetOk("vswitch_id"); ok && v != "" {
			reqQuery["VSwitchId"] = v
		}
		if v, ok := d.GetOk("ip"); ok && v != "" {
			reqQuery["Ip"] = v
		}
		if v, ok := d.GetOk("white_list"); ok && v != "" {
			reqQuery["WhiteList"] = v
		}

		_, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "ConfigNetwork", "", nil, reqQuery, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_hsm_instance", "ConfigNetwork", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func resourceAlibabacloudStackHsmInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "ReleaseInstance", "", nil, reqQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "ResourceNotExist") {
				return resource.NonRetryableError(err)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "ReleaseInstance", errmsgs.AlibabacloudStackSdkGoERROR, "")
			return resource.RetryableError(err)
		}
		addDebug("ReleaseInstance", raw)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ResourceNotExist") {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	return nil
}
