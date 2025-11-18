package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAPIGatewayV2CascadeInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cascade_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"console_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"console_ak": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"console_sk": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackAPIGatewayV2CascadeInstanceCreate, resourceAlibabacloudStackAPIGatewayV2CascadeInstanceRead, nil, resourceAlibabacloudStackAPIGatewayV2CascadeInstanceDelete)
	return resource
}

func resourceAlibabacloudStackAPIGatewayV2CascadeInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["instanceType"] = d.Get("instance_type")
	request["instanceName"] = d.Get("instance_name")
	request["cascadeInstanceId"] = d.Get("cascade_instance_id")
	request["consoleAddress"] = d.Get("console_address")
	request["consoleAk"] = d.Get("console_ak")
	request["consoleSk"] = d.Get("console_sk")

	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateCascadeInstance", "/cascadeInstance/createCascadeInstance", nil, nil, request)
	if err != nil {
		return err
	}

	data, ok := response["data"]
	if !ok {
		return fmt.Errorf("failed to get linkId from response")
	}

	instance_id, ok := data.(string)
	if !ok {
		return fmt.Errorf("invalid type of linkId in response")
	}

	d.SetId(instance_id)

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2CascadeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	service := ApiGateWayV2Service{client}

	object, err := service.DescribeCascadeInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("instance_type", object["instanceType"])
	d.Set("instance_name", object["instanceName"])
	d.Set("cascade_instance_id", object["cascadeInstanceId"])
	d.Set("console_address", object["consoleAddress"])
	d.Set("console_ak", object["consoleAk"])
	d.Set("console_sk", object["consoleSk"])

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2CascadeInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"instanceId": d.Id(),
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteCascadeInstance", "/cascadeInstance/deleteCascadeInstance", nil, nil, reqQuery)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteCascadeInstance", errmsgs.AlibabacloudStackSdkGoERROR, "")
			return resource.RetryableError(err)
		}
		// No need to check response data as per the documentation
		_ = raw
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
