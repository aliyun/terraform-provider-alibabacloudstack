package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
	request.RegionId = client.RegionId
	request.ClusterId = d.Get("cs_cluster_id").(string)
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	if v, ok := d.GetOk("namespace_id"); ok {
		request.NamespaceId = v.(string)
	}
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ImportK8sCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ImportK8sClusterResponse)
	if response.Code != 200 {
		return WrapError(Error("import k8s cluster failed for " + response.Message))
	}
	if len(response.Data) == 0 {
		return WrapError(Error("null cluster id after import k8s cluster"))
	}
	d.SetId(response.Data)
	// Wait until import succeed
	req := edas.CreateGetClusterRequest()
	req.ClusterId = response.Data
	req.Headers["x-ascm-product-name"] = "Edas"
	req.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	req.RegionId = client.RegionId
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(req)
		})
		time.Sleep(120 * time.Second)
		response, _ := raw.(*edas.GetClusterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if response.Code != 200 {
			return resource.NonRetryableError(Error("Get cluster failed for " + response.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	return resourceAlibabacloudStackEdasK8sClusterRead(d, meta)
}

func resourceAlibabacloudStackEdasK8sClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	object, err := edasService.DescribeEdasListCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	regionId := client.RegionId

	request := edas.CreateDeleteClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteCluster(request)
		})
		response, _ := raw.(*edas.DeleteClusterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if response.Code != 200 {
			return resource.NonRetryableError(Error("Delete EDAS K8s cluster failed for " + response.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	reqGet := edas.CreateGetClusterRequest()
	reqGet.RegionId = regionId
	reqGet.ClusterId = clusterId
	reqGet.Headers["x-ascm-product-name"] = "Edas"
	reqGet.Headers["x-acs-organizationid"] = client.Department
	reqGet.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(reqGet)
		})
		response, _ := raw.(*edas.GetClusterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)

		if response.Code == 200 {
			return resource.RetryableError(Error("cluster deleting"))
		} else if response.Code == 601 && strings.Contains(response.Message, "does not exist") {
			return nil
		} else {
			return resource.NonRetryableError(Error("check cluster status failed for " + response.Message))
		}
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	return nil
}
