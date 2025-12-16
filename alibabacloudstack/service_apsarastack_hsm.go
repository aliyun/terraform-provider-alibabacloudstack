package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type HsmService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *HsmService) DescribeHsmInstance(id string) (object map[string]interface{}, err error) {

	request := make(map[string]interface{})
	request["InstanceId"] = id

	resp, err := s.client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "DescribeInstances", "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	instances, ok := resp["Instances"].([]interface{})
	if !ok || len(instances) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found: Hsm Instance " + id)
	}

	object = instances[0].(map[string]interface{})
	if fmt.Sprint(object["HsmStatus"]) == "3" {
		// instance is deleted
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found: Hsm Instance " + id)
	}
	return object, nil
}

func (s *HsmService) DescribeHsmCluster(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{
		"CurrentPage": 1,
		"PageSize":    100,
	}

	response, err := s.client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "DescribeClusters", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	clusters, ok := response["HsmClusters"].([]interface{})
	if !ok || len(clusters) == 0 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Hsm Cluster", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	for _, v := range clusters {
		cluster := v.(map[string]interface{})
		if clusterId, ok := cluster["ClusterId"].(string); ok && clusterId == id {
			return cluster, nil
		}
	}
	return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Hsm Cluster", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *HsmService) HsmClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHsmCluster(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ClusterStatus"]) == failState {
				return object, fmt.Sprint(object["ClusterStatus"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["ClusterStatus"])))
			}
		}

		return object, fmt.Sprint(object["ClusterStatus"]), nil
	}
}

func (s *HsmService) DescribeQuickInitAdminName(cluster_id, instance_id string) (string, error) {
	reqQuery := map[string]interface{}{
		"ClusterId":  cluster_id,
		"InstanceId": instance_id,
	}

	response, err := s.client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "DescribeQuickInitInfo", "", nil, reqQuery, nil)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}

	adminName, ok := response["AdminName"].(string)
	if !ok || adminName == "" {
		return "", errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Hsm Cluster:%s QuickInit AdminName not found", cluster_id))
	}
	return adminName, nil
}
