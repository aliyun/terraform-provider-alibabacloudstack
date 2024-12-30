package alibabacloudstack

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasK8sCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEdasK8sClusterCreate,
		Read:   resourceAlibabacloudStackEdasK8sClusterRead,
		Delete: resourceAlibabacloudStackEdasK8sClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cs_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"network_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_import_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

type ImportK8sClusterResponse struct {
	*responses.BaseResponse
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	Data      string `json:"Data" xml:"Data"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

func resourceAlibabacloudStackEdasK8sClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "Edas", "2017-08-01", "ImportK8sCluster", "/pop/v5/import_k8s_cluster")
	request.QueryParams["ClusterId"] = d.Get("cs_cluster_id").(string)
	if v, ok := d.GetOk("namespace_id"); ok {
		request.QueryParams["NamespaceId"] = v.(string)
	}
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), bresponse, request)
	response := ImportK8sClusterResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	log.Printf("unmarshal response for read %v", &response)

	if response.Code != "200" {
		return errmsgs.WrapError(errmsgs.Error("import k8s cluster failed for " + response.Message))
	}
	if len(response.Data) == 0 {
		return errmsgs.WrapError(errmsgs.Error("null cluster id after import k8s cluster"))
	}
	d.SetId(response.Data)
	// Wait until import succeed
	request = client.NewCommonRequest("GET", "Edas", "2017-08-01", "GetCluster", "/pop/v5/resource/cluster")
	request.QueryParams["ClusterId"] = d.Get("cs_cluster_id").(string)
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		time.Sleep(120 * time.Second)
		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		response := edas.GetClusterResponse{}
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return resource.RetryableError(errmsgs.WrapError(err))
		}
		if response.Code != 200 {
			return resource.NonRetryableError(errmsgs.Error("Get cluster failed for " + response.Message))
		}

		addDebug(request.GetActionName(), response, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return resourceAlibabacloudStackEdasK8sClusterRead(d, meta)
}

func resourceAlibabacloudStackEdasK8sClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	object, err := edasService.DescribeEdasListCluster(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	region := object.RegionId
	pos := strings.Index(region, ":")
	// get ":", should intercept the string
	if pos != -1 {
		region = region[0:pos]
	}
	d.Set("cluster_name", object.ClusterName)
	d.Set("cluster_type", object.ClusterType)
	d.Set("network_mode", object.NetworkMode)
	d.Set("vpc_id", object.VpcId)
	d.Set("namespace_id", region)
	d.Set("cluster_import_status", object.ClusterImportStatus)
	d.Set("cs_cluster_id", object.CsClusterId)

	return nil
}

func resourceAlibabacloudStackEdasK8sClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	clusterId := d.Id()

	request := client.NewCommonRequest("DELETE", "Edas", "2017-08-01", "DeleteCluster", "/pop/v5/resource/cluster")
	request.QueryParams["ClusterId"] = clusterId
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	wait := incrementalWait(1*time.Second, 2*time.Second)
	response := edas.DeleteClusterResponse{}
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(err)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if response.Code != 200 {
			return resource.NonRetryableError(errmsgs.Error("Delete EDAS K8s cluster failed for " + response.Message))
		}

		addDebug(request.GetActionName(), response, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	request = client.NewCommonRequest("GET", "Edas", "2017-08-01", "GetCluster", "/pop/v5/resource/cluster")
	request.QueryParams["ClusterId"] = d.Get("cs_cluster_id").(string)
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(err)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), bresponse, request)
		response := edas.GetClusterResponse{}
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if response.Code == 200 {
			return resource.RetryableError(errmsgs.Error("cluster deleting"))
		} else if response.Code == 601 && strings.Contains(response.Message, "does not exist") {
			return nil
		} else {
			return resource.NonRetryableError(errmsgs.Error("check cluster status failed for " + response.Message))
		}
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
