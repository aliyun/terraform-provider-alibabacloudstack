package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackApiGatewayV2CascadeLink() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_instance_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cascade_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"link_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"link_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cascade_service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cascade_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackApiGatewayV2CascadeLinkCreate, resourceAlibabacloudStackApiGatewayV2CascadeLinkRead, resourceAlibabacloudStackApiGatewayV2CascadeLinkUpdate, resourceAlibabacloudStackApiGatewayV2CascadeLinkDelete)
	return resource
}

func resourceAlibabacloudStackApiGatewayV2CascadeLinkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["sourceInstanceId"] = d.Get("source_instance_id")
	request["sourceInstanceAddress"] = d.Get("source_instance_address")
	request["cascadeInstanceId"] = d.Get("cascade_instance_id")
	request["linkName"] = d.Get("link_name")

	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateCascadeLink", "/cascadeLink/createCascadeLink", nil, nil, request)
	if err != nil {
		return err
	}

	linkId, ok := response["data"].(string)
	if !ok {
		return fmt.Errorf("failed to get linkId from response")
	}

	d.SetId(linkId)

	return nil
}

func resourceAlibabacloudStackApiGatewayV2CascadeLinkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	service := ApiGateWayV2Service{client}

	object, err := service.DescribeCascadeLink(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("source_instance_id", object["sourceInstanceId"])
	d.Set("source_instance_address", object["sourceInstanceAddress"])
	d.Set("cascade_instance_id", object["cascadeInstanceId"])
	d.Set("link_name", object["linkName"])
	d.Set("link_id", object["linkId"])
	d.Set("cascade_service_id", object["cascadeServiceId"])
	d.Set("source_instance_name", object["sourceInstanceName"])
	d.Set("cascade_instance_name", object["cascadeInstanceName"])

	return nil
}

func resourceAlibabacloudStackApiGatewayV2CascadeLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}
	if d.HasChange("source_instance_address") {
		request := make(map[string]interface{})
		request["linkId"] = d.Get("link_id")
		request["cascadeInstanceId"] = d.Get("cascade_instance_id")
		request["sourceInstanceAddress"] = d.Get("source_instance_address")
		request["sourceInstanceId"] = d.Get("source_instance_id")
		request["linkName"] = d.Get("link_name")

		if v, ok := d.GetOk("cascade_service_id"); ok {
			request["cascadeServiceId"] = v
		}
		if v, ok := d.GetOk("source_instance_name"); ok {
			request["sourceInstanceName"] = v
		}
		if v, ok := d.GetOk("cascade_instance_name"); ok {
			request["cascadeInstanceName"] = v
		}

		_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifyCascadeLink", "/cascadeLink/modifyCascadeLink", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_cascade_link", "ModifyCascadeLink", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackApiGatewayV2CascadeLinkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"linkId": d.Id(),
	}

	_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteCascadeLink", "/cascadeLink/deleteCascadeLink", nil, nil, reqQuery)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteCascadeLink", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
