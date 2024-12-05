package alibabacloudstack

import (
	"strings"
	"time"

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

func resourceAlibabacloudStackEdasK8sClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	request := edas.CreateImportK8sClusterRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.ClusterId = d.Get("cs_cluster_id").(string)
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	if v, ok := d.GetOk("namespace_id"); ok {
		request.NamespaceId = v.(string)
	}
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ImportK8sCluster(request)
	})

	bresponse, ok := raw.(*edas.ImportK8sClusterResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if bresponse.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("import k8s cluster failed for " + bresponse.Message))
	}
	if len(bresponse.Data) == 0 {
		return errmsgs.WrapError(errmsgs.Error("null cluster id after import k8s cluster"))
	}
	d.SetId(bresponse.Data)
	// Wait until import succeed
	req := edas.CreateGetClusterRequest()
	client.InitRoaRequest(*req.RoaRequest)
	req.ClusterId = bresponse.Data
	req.Headers["x-acs-content-type"] = "application/json"
	req.Headers["Content-Type"] = "application/json"
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(req)
		})
		time.Sleep(120 * time.Second)
		bresponse, ok := raw.(*edas.GetClusterResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		if bresponse.Code != 200 {
			return resource.NonRetryableError(errmsgs.Error("Get cluster failed for " + bresponse.Message))
		}

		addDebug(req.GetActionName(), raw, req.RoaRequest, req)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
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
	edasService := EdasService{client}

	clusterId := d.Id()

	request := edas.CreateDeleteClusterRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.ClusterId = clusterId
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteCluster(request)
		})
		bresponse, ok := raw.(*edas.DeleteClusterResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		if bresponse.Code != 200 {
			return resource.NonRetryableError(errmsgs.Error("Delete EDAS K8s cluster failed for " + bresponse.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	reqGet := edas.CreateGetClusterRequest()
	client.InitRoaRequest(*reqGet.RoaRequest)
	reqGet.ClusterId = clusterId
	reqGet.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(reqGet)
		})
		bresponse, ok := raw.(*edas.GetClusterResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", reqGet.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(reqGet.GetActionName(), raw, reqGet.RoaRequest, reqGet)

		if bresponse.Code == 200 {
			return resource.RetryableError(errmsgs.Error("cluster deleting"))
		} else if bresponse.Code == 601 && strings.Contains(bresponse.Message, "does not exist") {
			return nil
		} else {
			return resource.NonRetryableError(errmsgs.Error("check cluster status failed for " + bresponse.Message))
		}
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), reqGet.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
