package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasK8sCluster() *schema.Resource {
	resource := &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: resourceAlibabacloudStackEdasK8sClusterImportState,
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
	setResourceFunc(resource, resourceAlibabacloudStackEdasK8sClusterCreate, resourceAlibabacloudStackEdasK8sClusterRead, nil, resourceAlibabacloudStackEdasK8sClusterDelete)
	return resource
}

func upgradeEdasK8sClusterId(oldId, regionId string) string {
	if !strings.Contains(oldId, ":") && regionId != "" {
		return fmt.Sprintf("%s:%s", regionId, oldId)
	}
	return oldId
}

func resourceAlibabacloudStackEdasK8sClusterImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	if !strings.Contains(id, ":") {
		client := meta.(*connectivity.AlibabacloudStackClient)
		edasService := EdasService{client}
		object, err := edasService.DescribeEdasK8sCluster(id)
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
		newId := upgradeEdasK8sClusterId(id, object.RegionId)
		d.SetId(newId)
		if object.RegionId != "" {
			d.Set("namespace_id", object.RegionId)
		}
	}
	return []*schema.ResourceData{d}, nil
}

type ImportK8sClusterResponse struct {
	Message   string `json:"Message" xml:"Message"`
	Data      string `json:"Data" xml:"Data"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

func resourceAlibabacloudStackEdasK8sClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	retry := 5
	var id string
	var err error
	namespaceId := d.Get("namespace_id").(string)
	for retry > 0 {
		id, err = importK8sCluster(client, d.Get("cs_cluster_id").(string), namespaceId)
		if err == nil {
			if namespaceId != "" {
				d.SetId(fmt.Sprintf("%s:%s", namespaceId, id))
			} else {
				d.SetId(id)
			}
			break
		}
		retry -= 1
	}
	return err
}

func importK8sCluster(client *connectivity.AlibabacloudStackClient, clusterId, namesapceId string) (string, error) {
	var id string

	request := client.NewCommonRequest("POST", "Edas", "2017-08-01", "ImportK8sCluster", "/pop/v5/import_k8s_cluster")
	request.QueryParams["ClusterId"] = clusterId
	if namesapceId != "" {
		request.QueryParams["RegionId"] = namesapceId
	}
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return id, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return id, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response := ImportK8sClusterResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return id, errmsgs.WrapError(err)
	}

	log.Printf("unmarshal response for read %v", &response)

	if len(response.Data) == 0 {
		return id, errmsgs.WrapError(errmsgs.Error("null cluster id after import k8s cluster"))
	}
	id = response.Data
	// Wait until import succeed
	edasService := EdasService{client}
	stateConf := BuildStateConf([]string{"3"}, []string{"1"}, 30*time.Minute, 10*time.Second, edasService.ClusterImportK8sStateRefreshFunc(id, []string{"0", "2", "4"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return id, errmsgs.WrapErrorf(err, errmsgs.IdMsg, id)
	}

	return id, nil
}

func resourceAlibabacloudStackEdasK8sClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	object, err := edasService.DescribeEdasK8sCluster(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Compatible with the legacy id format (clusterid only); automatically upgrade it to the new format (namespaceid:clusterid)
	newId := upgradeEdasK8sClusterId(d.Id(), object.RegionId)
	if newId != d.Id() {
		d.SetId(newId)
	}

	d.Set("cluster_name", object.ClusterName)
	d.Set("cluster_type", object.ClusterType)
	d.Set("network_mode", object.NetworkMode)
	d.Set("vpc_id", object.VpcId)
	d.Set("namespace_id", object.RegionId)
	d.Set("cluster_import_status", object.ClusterImportStatus)
	d.Set("cs_cluster_id", object.CsClusterId)

	return nil
}

func resourceAlibabacloudStackEdasK8sClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	clusterId := d.Id()
	if idx := strings.LastIndex(clusterId, ":"); idx != -1 && idx+1 < len(clusterId) {
		clusterId = clusterId[idx+1:]
	}

	request := client.NewCommonRequest("DELETE", "Edas", "2017-08-01", "DeleteCluster", "/pop/v5/resource/cluster")
	request.QueryParams["ClusterId"] = clusterId
	// request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request)
	if err != nil {
		if bresponse == nil {
			return err
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	edasService := EdasService{client}
	stateConf := BuildStateConf([]string{"4"}, []string{"0"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, edasService.ClusterImportK8sStateRefreshFunc(d.Id(), []string{"1", "2", "3"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return nil
	}
	return nil
}
