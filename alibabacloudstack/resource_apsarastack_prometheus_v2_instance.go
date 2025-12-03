package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPrometheusV2Instance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_api": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_write_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"push_gateway_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"objid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPrometheusV2InstanceCreate, resourceAlibabacloudStackPrometheusV2InstanceRead, resourceAlibabacloudStackPrometheusV2InstanceUpdate, resourceAlibabacloudStackPrometheusV2InstanceDelete)
	return resource
}

func resourceAlibabacloudStackPrometheusV2InstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	body := map[string]interface{}{
		"clusterName": d.Get("cluster_name").(string),
	}

	resp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "InstallCommonPrometheusInstance", "/log/api/v2/prometheus/installCommonPrometheusInstance", nil, nil, body)
	if err != nil {
		return err
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse response data")
	}

	clusterId, ok := data["acosClusterId"].(string)
	if !ok {
		return fmt.Errorf("failed to get cluster ID from response")
	}
	d.SetId(clusterId)
	return nil
}

func resourceAlibabacloudStackPrometheusV2InstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	prometheusService := PrometheusService{client}
	object, err := prometheusService.DescribePrometheusV2Instance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DescribePrometheusV2InstanceReal", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.Set("cluster_id", object["clusterId"])
	d.Set("cluster_name", object["clusterName"])
	d.Set("http_api", object["httpApi"])
	d.Set("remote_write_url", object["remoteWriteUrl"])
	d.Set("push_gateway_url", object["pushGatewayUrl"])
	d.Set("objid", object["id"])
	real, err := prometheusService.DescribePrometheusV2InstanceReal(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DescribePrometheusV2InstanceReal", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if tagSet, ok := real["tagSet"].([]interface{}); ok {
		d.Set("tags", tagSet)
	}

	return nil
}

func resourceAlibabacloudStackPrometheusV2InstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	prometheusService := PrometheusService{client}
	if d.HasChange("tags") {
		objid := d.Get("objid")
		if objid == 0 {
			object, err := prometheusService.DescribePrometheusV2Instance(d.Id())
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DescribePrometheusV2Instance", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			objid = object["id"]
		}
		tags := d.Get("tags").(*schema.Set).List()
		reqBody := map[string]interface{}{
			"ids":    []interface{}{objid},
			"tagSet": tags,
			"type":   "update",
		}

		_, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "AddClusterTags", "/log/api/v2/cloud-native/cluster/add-tags", nil, nil, reqBody)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_prometheus_v2_instance", "AddClusterTags", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackPrometheusV2InstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := map[string]interface{}{
		"clusterId": d.Id(),
	}

	_, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "DeleteCluster", "/log/api/v2/cloud-native/cluster/delete", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteCluster", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
