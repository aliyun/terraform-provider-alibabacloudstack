package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCenVbrHealthCheck() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vbr_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_check_source_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_check_target_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 8),
			},
			"health_check_only": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"vbr_instance_owner_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Computed fields from Describe API
			"link_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"packet_loss": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCenVbrHealthCheckCreate,
		resourceAlibabacloudStackCenVbrHealthCheckRead,
		nil,
		resourceAlibabacloudStackCenVbrHealthCheckDelete)

	return resource
}

func resourceAlibabacloudStackCenVbrHealthCheckCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	// Prepare the request parameters for EnableCenVbrHealthCheck API
	reqQuery := make(map[string]interface{})
	reqQuery["CenId"] = d.Get("cen_id")
	reqQuery["VbrInstanceRegionId"] = client.RegionId
	reqQuery["VbrInstanceId"] = d.Get("vbr_instance_id")
	reqQuery["HealthCheckTargetIp"] = d.Get("health_check_target_ip")
	reqQuery["HealthCheckSourceIp"] = d.Get("health_check_source_ip")
	reqQuery["HealthCheckInterval"] = d.Get("health_check_interval")
	reqQuery["HealthyThreshold"] = d.Get("healthy_threshold")
	reqQuery["HealthCheckOnly"] = d.Get("health_check_only")
	// reqQuery["VbrInstanceOwnerId"] = d.Get("vbr_instance_owner_id")

	// Call the EnableCenVbrHealthCheck API to create the health check configuration
	if _, err := client.DoTeaRequest("POST", "cbn", "2017-09-12", "EnableCenVbrHealthCheck", "", nil, reqQuery, nil); err != nil {
		return err
	}

	// Generate resource ID based on {CenId:VbrInstanceId}
	cenId := d.Get("cen_id").(string)
	vbrInstanceId := d.Get("vbr_instance_id").(string)
	resourceId := fmt.Sprintf("%s:%s", cenId, vbrInstanceId)

	// Set the temporary ID
	d.SetId(resourceId)
	return nil
}

func resourceAlibabacloudStackCenVbrHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cenService := CenService{client}

	object, err := cenService.DescribeCenVbrHealthCheck(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cen_id", object["CenId"])
	d.Set("vbr_instance_id", object["VbrInstanceId"])
	d.Set("health_check_source_ip", object["HealthCheckSourceIp"])
	d.Set("health_check_target_ip", object["HealthCheckTargetIp"])
	d.Set("health_check_interval", object["HealthCheckInterval"])
	d.Set("healthy_threshold", object["HealthyThreshold"])
	d.Set("health_check_only", object["HealthCheckOnly"])
	if v, ok := object["LinkStatus"]; ok && v != nil {
		d.Set("link_status", v)
	}
	if v, ok := object["Delay"]; ok && v != nil {
		d.Set("delay", v)
	}
	if v, ok := object["PacketLoss"]; ok && v != nil {
		d.Set("packet_loss", v)
	}

	return nil
}

func resourceAlibabacloudStackCenVbrHealthCheckDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"CenId":               d.Get("cen_id"),
		"VbrInstanceRegionId": client.RegionId,
		"VbrInstanceId":       d.Get("vbr_instance_id"),
	}

	raw, err := client.DoTeaRequest("POST", "Cbn", "2017-09-12", "DisableCenVbrHealthCheck", "", nil, reqQuery, nil)
	addDebug("DisableCenVbrHealthCheck", raw, reqQuery)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DisableCenVbrHealthCheck", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
