package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type FlinkService struct {
	client *connectivity.AlibabacloudStackClient
}

func (f *FlinkService) DescribeFlinkNamespace(id string) (map[string]interface{}, error) {
	response, err := f.client.DoTeaRequest("GET", "ververica", "2020-05-01", "DescribeNamespaces", "/flink/namespace/list", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	vs, err := jsonpath.Get("$.data", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
	}

	var result map[string]interface{}
	data := vs.([]interface{})
	for _, v := range data {
		namespace := v.(map[string]interface{})
		if namespace["Name"].(string) == id {
			result = namespace
			break
		}
	}

	return result, nil
}
