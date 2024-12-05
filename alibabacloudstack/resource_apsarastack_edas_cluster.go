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
	return &schema.Resource{
		Create: resourceAlibabacloudStackEdasClusterCreate,
		Read:   resourceAlibabacloudStackEdasClusterRead,
		Delete: resourceAlibabacloudStackEdasClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackEdasClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	request := edas.CreateInsertClusterRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.ClusterName = d.Get("cluster_name").(string)
	request.ClusterType = requests.NewInteger(d.Get("cluster_type").(int))
	request.NetworkMode = requests.NewInteger(d.Get("network_mode").(int))
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
	log.Printf("request domainaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa: %s", request.Domain)
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if bresponse.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("create cluster failed for " + bresponse.Message))
	}
	d.SetId(bresponse.Cluster.ClusterId)

	return resourceAlibabacloudStackEdasClusterRead(d, meta)
}

func resourceAlibabacloudStackEdasClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	clusterId := d.Id()

	request := edas.CreateGetClusterRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.ClusterId = clusterId

	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	bresponse, ok := raw.(*edas.GetClusterResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if bresponse.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("create cluster failed for " + bresponse.Message))
	}

	d.Set("cluster_name", bresponse.Cluster.ClusterName)
	d.Set("cluster_type", bresponse.Cluster.ClusterType)
	d.Set("network_mode", bresponse.Cluster.NetworkMode)
	//d.Set("region_id", bresponse.Cluster.RegionId)
	d.Set("vpc_id", bresponse.Cluster.VpcId)

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
