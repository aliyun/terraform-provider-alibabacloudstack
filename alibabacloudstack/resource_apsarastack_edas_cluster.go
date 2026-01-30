package alibabacloudstack

import (
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasCluster() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"network_mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
			"logical_region_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"region_id"},
			},
			"region_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Deprecated:    "Field 'region_id' is deprecated and will be removed in a future release. Please use new field 'logical_region_id' instead.",
				ConflictsWith: []string{"logical_region_id"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasClusterCreate, resourceAlibabacloudStackEdasClusterRead, nil, resourceAlibabacloudStackEdasClusterDelete)
	return resource
}

func resourceAlibabacloudStackEdasClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	request := edas.CreateInsertClusterRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.ClusterName = d.Get("cluster_name").(string)
	request.ClusterType = requests.NewInteger(d.Get("cluster_type").(int))
	request.NetworkMode = requests.NewInteger(d.Get("network_mode").(int))
	if v, ok := connectivity.GetResourceDataOk(d, "logical_region_id", "region_id"); ok && v.(string) != "" {
		request.LogicalRegionId = v.(string)
	}
	request.OversoldFactor = requests.NewInteger(1)
	request.IaasProvider = "ALIYUN"

	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"

	if v, ok := d.GetOk("vpc_id"); !ok {
		if d.Get("network_mode") == 2 {
			return errmsgs.WrapError(errmsgs.Error("vpcId is required for vpc network mode"))
		}
	} else {
		request.VpcId = v.(string)
	}
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertCluster(request)
	})

	bresponse, ok := raw.(*edas.InsertClusterResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	log.Printf("request domain: %s", request.Domain)
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if bresponse.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("create cluster failed for " + bresponse.Message))
	}
	d.SetId(bresponse.Cluster.ClusterId)

	return nil
}

func resourceAlibabacloudStackEdasClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	clusterId := d.Id()
	cluster, err := edasService.DescribeEdasGetCluster(clusterId)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cluster_name", cluster.ClusterName)
	d.Set("cluster_type", cluster.ClusterType)
	d.Set("network_mode", cluster.NetworkMode)
	connectivity.SetResourceData(d, cluster.RegionId, "logical_region_id", "region_id")
	d.Set("vpc_id", cluster.VpcId)

	return nil
}

func resourceAlibabacloudStackEdasClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	clusterId := d.Id()

	request := edas.CreateDeleteClusterRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.ClusterId = clusterId

	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
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
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		if bresponse.Code != 200 {
			if strings.Contains(bresponse.Message, "there are still instances in it") {
				return resource.RetryableError(errmsgs.Error("delete cluster failed for " + bresponse.Message))
			}
			return resource.NonRetryableError(errmsgs.Error("delete cluster failed for " + bresponse.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
