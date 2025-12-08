package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type PrometheusService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *PrometheusService) DescribePrometheusV2Instance(id string) (map[string]interface{}, error) {
	now := time.Now()
	endTime := now.Format("2006-01-02 15:04:05")
	startTime := now.Add(-5 * time.Minute).Format("2006-01-02 15:04:05")
	period := fmt.Sprintf("%s~%s", startTime, endTime)
	reqQuery := map[string]interface{}{
		"period":     period,
		"sumMetrics": false,
		"timeType":   2,
	}
	resp, err := s.client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "ListCluster", "/log/api/v2/cloud-native/cluster/list", nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}

	if data, ok := resp["data"].([]interface{}); ok {

		for _, item := range data {
			cluster := item.(map[string]interface{})
			if cluster["clusterId"].(string) == id {
				return cluster, nil
			}
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Prometheus instance with ClusterId %s not found", id))
}

func (s *PrometheusService) DescribePrometheusV2InstanceReal(id string) (map[string]interface{}, error) {
	now := time.Now()
	endTime := now.Format("2006-01-02 15:04:05")
	startTime := now.Add(-5 * time.Minute).Format("2006-01-02 15:04:05")
	period := fmt.Sprintf("%s~%s", startTime, endTime)
	reqQuery := map[string]interface{}{
		"period":     period,
		"sumMetrics": false,
		"timeType":   2,
	}
	resp, err := s.client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "ListClusterReal", "/log/api/v2/cloud-native/cluster/list2", nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}

	if data, ok := resp["data"].([]interface{}); ok {
		for _, item := range data {
			cluster := item.(map[string]interface{})
			if cluster["clusterId"].(string) == id {
				return cluster, nil
			}
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Prometheus instance with ClusterId %s not found", id))
}
func (s *PrometheusService) DescribePrometheusV2NotifyGroup(id string) (map[string]interface{}, error) {
	clusterId := id

	reqBody := map[string]interface{}{
		"keyword":         "",
		"pageNumber":      1,
		"pageSize":        10,
		"direction":       "desc",
		"page":            1,
		"selectedRowKeys": []interface{}{},
	}

	response, err := s.client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "PageNotifyGroup", "/log/api/v2/alert/group/page", nil, nil, reqBody)
	if err != nil {
		return nil, err
	}

	if success, ok := response["success"].(bool); !ok || !success {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("prometheus notify group %s not found", clusterId))
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("prometheus notify group %s not found", clusterId))
	}

	dataList, ok := data["data"].([]interface{})
	if !ok || len(dataList) == 0 {
		return nil, errmsgs.Error(errmsgs.NotFoundMsg, fmt.Sprintf("prometheus notify group %s not found", clusterId))
	}

	for _, item := range dataList {
		group, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if groupId, ok := group["id"]; ok && fmt.Sprintf("%v", groupId) == clusterId {
			return group, nil
		}
	}

	return nil, errmsgs.Error(errmsgs.NotFoundMsg, fmt.Sprintf("prometheus notify group %s not found", clusterId))
}
