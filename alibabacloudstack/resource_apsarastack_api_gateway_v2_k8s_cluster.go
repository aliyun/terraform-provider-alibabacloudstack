package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackApigatewayv2K8sCluster() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cs_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"k8s_cluster_name": {
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"config_contents": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackApigatewayv2K8sClusterCreate,
		resourceAlibabacloudStackApigatewayv2K8sClusterRead, nil,
		resourceAlibabacloudStackApigatewayv2K8sClusterDelete)
	return resource
}

func resourceAlibabacloudStackApigatewayv2K8sClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}

	// Prepare request parameters
	reqBody := make(map[string]interface{})
	if v, ok := d.GetOk("cs_cluster_id"); ok {
		reqBody["k8sClusterType"] = "container-service"
		reqBody["csClusterId"] = v.(string)
		if object, err := csService.DescribeCsKubernetes(v.(string)); err != nil {
			return err
		} else {
			reqBody["csClusterName"] = object.Name
		}
		if v, ok := d.GetOk("vpc_id"); ok {
			reqBody["vpcId"] = v.(string)
			reqBody["slbType"] = "intranet"
		} else {
			reqBody["slbType"] = "internet"
		}
		if v, ok := d.GetOk("config_contents"); ok && v.(string) != "" {
			reqBody["configContents"] = v.(string)
		} else {
			var privateAddress bool
			if reqBody["slbType"] == "intranet" {
				privateAddress = true
			} else {
				privateAddress = false
			}
			if configContents, err := csService.GetK8sCluterKubeConfig(reqBody["csClusterId"].(string), privateAddress); err != nil {
				return err
			} else {
				reqBody["configContents"] = configContents
			}
		}
	} else {
		reqBody["k8sClusterType"] = "self-built"
		if v, ok := d.GetOk("config_contents"); ok && v.(string) != "" {
			reqBody["configContents"] = v.(string)
		} else {
			return fmt.Errorf("configContents is necessory while cs_cluster_id not set")
		}
	}
	reqBody["k8sClusterName"] = d.Get("k8s_cluster_name").(string)

	// Call ImportCluster API
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ImportCluster", "/k8s/importCluster", nil, nil, reqBody)
	if err != nil {
		return err
	}

	// Extract k8s cluster ID from response
	data, ok := resp["data"]
	if !ok {
		return fmt.Errorf("failed to get data from ImportCluster response")
	}
	k8sClusterCode, ok := data.(map[string]interface{})["data"].(string)
	if !ok || k8sClusterCode == "" {
		return fmt.Errorf("failed to extract k8sClusterCode from ImportCluster response")
	}

	// Set the temporary ID
	d.SetId(k8sClusterCode)

	return nil
}

func resourceAlibabacloudStackApigatewayv2K8sClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apigatewayv2Service := ApiGateWayV2Service{client}

	object, err := apigatewayv2Service.DescribeApigatewayv2K8sCluster(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_apigatewayv2_k8s_cluster apigatewayv2Service.DescribeApigatewayv2K8sCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("k8s_cluster_name", object["k8sClusterName"])
	if attr, ok := object["k8sClusterAttribute"].(map[string]interface{}); ok {
		d.Set("cs_cluster_id", attr["csClusterId"])
		if v, exist := attr["vpcId"]; exist {
		d.Set("vpc_id",v)
		}
	}

	return nil
}

func resourceAlibabacloudStackApigatewayv2K8sClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare the request body for CancelImportCluster API
	reqBody := map[string]interface{}{
		"k8sClusterCode": d.Id(),
	}

	// Call the CancelImportCluster API
	_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CancelImportCluster", "/k8s/cancelImportCluster", nil, nil, reqBody)
	if err != nil {
		return fmt.Errorf("failed to cancel import cluster: %v", err)
	}

	return nil
}
