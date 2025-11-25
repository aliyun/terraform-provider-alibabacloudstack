package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ApiGateWayV2Service struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *ApiGateWayV2Service) DescribeApiGatewayV2Instance(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"gwInstanceId": id,
	}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetInstanceInfo", "/gatewayInstance/getInstanceInfo", nil, request, request)
	if err != nil {
		return nil, err
	}
	data, ok := response["data"]
	if !ok {
		return nil, errmsgs.Error("GetInstanceInfo Failed! %v", response)
	}
	return data.(map[string]interface{}), nil
}

func (s *ApiGateWayV2Service) ApiGateWayV2InstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApiGatewayV2Instance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["status"]) == failState {
				return object, fmt.Sprint(object["status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["status"])))
			}
		}
		return object, fmt.Sprint(object["status"]), nil
	}
}

func (s *ApiGateWayV2Service) GetCustomDeployConfig(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"gwInstanceId": id,
	}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetCustomDeployConfig", "/gatewayInstance/getCustomDeployConfig", nil, request, request)
	if err != nil {
		return nil, err
	}
	data, ok := response["data"]
	if !ok {
		return nil, errmsgs.Error("GetCustomDeployConfig Failed! %v", response)
	}
	return data.(map[string]interface{}), nil
}

func (s *ApiGateWayV2Service) DescribeApigatewayv2K8sCluster(id string) (map[string]interface{}, error) {
	reqBody := map[string]interface{}{
		"current": 1,
		"size":    100,
	}

	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListClusters", "/k8s/listClusters", nil, nil, reqBody)
	if err != nil {
		return nil, err
	}

	if records, ok := response["data"].(map[string]interface{})["records"].([]interface{}); ok {
		for _, record := range records {
			item := record.(map[string]interface{})
			if k8sClusterCode, ok := item["k8sClusterCode"].(string); ok && k8sClusterCode == id {
				return item, nil
			}
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Apigatewayv2 K8s Cluster %s was not found", id))
}

func (s *ApiGateWayV2Service) DescribeApiGatewayV2Certificate(id string) (map[string]interface{}, error) {
	params := strings.Split(id, ":")
	request := map[string]interface{}{
		"gwInstanceId": params[0],
	}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListCertificates", "/certificate/listCertificates", nil, request, request)
	if err != nil {
		return nil, err
	}
	data, err := jsonpath.Get("$.data.records", response)
	if err != nil {
		return nil, errmsgs.Error("GetInstanceInfo Failed! %v", response)
	}
	for _, v := range data.([]interface{}) {
		certificate := v.(map[string]interface{})
		if certificate["certificateId"].(string) == params[1] {
			return certificate, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Not Found ApiGateway V2 Certificate " + id)
}

func (s *ApiGateWayV2Service) DescribeApiGatewayV2Domain(id string) (map[string]interface{}, error) {
	params := strings.Split(id, ":")
	request := map[string]interface{}{
		"gwInstanceId": params[0],
		"domainId":     params[1],
	}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetDomain", "/domain/getDomain", nil, nil, request)
	if err != nil {
		return nil, err
	}
	data, ok := response["data"]
	if !ok {
		return nil, errmsgs.Error("GetDomain Failed! %v", response)
	}
	return data.(map[string]interface{}), nil
}

func (s *ApiGateWayV2Service) DescribeApiGatewayV2Signature(id string) (map[string]interface{}, error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, err
	}
	idpre := parts[0]
	gwInstanceId := parts[1]
	sigSchemeId := parts[2]

	reqQuery := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
	}
	action := "ListSignatureSchemes"
	pattern := "/signatureScheme/listSignatureSchemes"
	if idpre == "sourceSig" {
		action = "GetSourceSigScheme"
		pattern = "/sourceSigScheme/getSourceSigScheme"
		reqQuery["sigSchemeId"] = sigSchemeId
	} else {
		reqQuery["current"] = 1
		reqQuery["size"] = 100
	}

	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}
	if action == "GetSourceSigScheme" {
		if data, ok := response["data"]; ok {
			return data.(map[string]interface{}), nil
		}
	} else {
		records, err := jsonpath.Get("$.data.records", response)
		if err != nil {
			return nil, errmsgs.Error("ListSignatureSchemes Failed! %v", response)
		}
		for _, v := range records.([]interface{}) {
			record := v.(map[string]interface{})
			if record["sigSchemeId"].(string) == sigSchemeId {
				return record, nil
			}
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("signature scheme %s not found", id))
}

func (s *ApiGateWayV2Service) DescribeAPIGatewayV2Consumer(id string) (map[string]interface{}, error) {

	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, expected idpre:gwInstanceId:appId")
	}

	idpre := parts[0]
	gwInstanceId := parts[1]
	appId := parts[2]

	reqQuery := map[string]interface{}{
		"appId":        appId,
		"gwInstanceId": gwInstanceId,
	}
	action := "GetApp"
	pattern := "/application/getApp"
	if idpre == "sourceApp" {
		action = "GetSourceApplication"
		pattern = "/sourceApplication/getSourceApplication"
	}

	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}

	if fmt.Sprint(response["code"]) != "200" {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("api gateway v2 consumer %s not found", id))
	}

	data := response["data"].(map[string]interface{})
	return data, nil
}

func (s *ApiGateWayV2Service) DescribeApiGatewayV2ServiceSource(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, expected gwInstanceId:sourceId")
	}
	gwInstanceId := parts[0]
	sourceId := parts[1]

	reqQuery := map[string]interface{}{
		"sourceId":     sourceId,
		"gwInstanceId": gwInstanceId,
	}

	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetSource", "/source/getSource", nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}
	data := response["data"].(map[string]interface{})
	return data, nil
}

func (s *ApiGateWayV2Service) DescribeRouteGroup(id string) (map[string]interface{}, error) {
	// Split the id into gwInstanceId and groupId
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid route group id: %s", id)
	}
	gwInstanceId := parts[0]
	groupId := parts[1]

	// Prepare the request parameters
	reqQuery := map[string]interface{}{
		"groupId":      groupId,
		"gwInstanceId": gwInstanceId,
	}

	// Call the API to get the route group details
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetGroup", "/group/getGroup", nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}

	// Check if the response contains data
	data, ok := response["data"]
	if !ok || data == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("route group %s not found", id))
	}

	// Convert the data to map[string]interface{}
	result, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return result, nil
}
func (s *ApiGateWayV2Service) DescribeCascadeLink(id string) (map[string]interface{}, error) {

	reqQuery := map[string]interface{}{
		"linkId": id,
	}

	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetCascadeLink", "/cascadeLink/getCascadeLink", nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}

	if fmt.Sprint(response["code"]) != "200" {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("cascade link %s not found", id))
	}

	data := response["data"].(map[string]interface{})
	return data, nil
}

func (s *ApiGateWayV2Service) DescribeCascadeInstance(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{"instanceId": id}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetCascadeInstance", "/cascadeInstance/getCascadeInstance", nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}

	if !response["asapiSuccess"].(bool) || response["data"] == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("CascadeInstance %s not found", id))
	}

	return response["data"].(map[string]interface{}), nil
}
func (s *ApiGateWayV2Service) DescribeApiGatewayV2Service(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, "^", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, expected gwInstanceId:serviceId")
	}
	gwInstanceId := parts[0]
	serviceId := parts[1]

	reqQuery := map[string]interface{}{
		"serviceId":    serviceId,
		"gwInstanceId": gwInstanceId,
	}

	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetService", "/microservice/getService", nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}

	if fmt.Sprint(response["code"]) != "200" {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("service %s not found", id))
	}

	data := response["data"].(map[string]interface{})
	return data, nil
}

func (s *ApiGateWayV2Service) ListApiGatewayV2Service(instanceId string) ([]interface{}, error) {
	request := map[string]interface{}{
		"gwInstanceId": instanceId,
		"current":      1,
		"size":         100, // Assuming a reasonable page size; adjust if needed
	}

	resp, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListServices", "/microservice/listServices", nil, nil, request)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_services", "ListServices", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response
	records, err := jsonpath.Get("$.data.records", resp)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_services", "ListServices", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return records.([]interface{}), nil
}

func (s *ApiGateWayV2Service) DescribeApigwV2Route(id string) (map[string]interface{}, error) {
	// Split the id into gwInstanceId and groupId
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid route group id: %s", id)
	}
	idpre := parts[0]
	gwInstanceId := parts[1]
	routeId := parts[2]
	action := "GetRoute"
	pattern := "/route/getRoute"
	if idpre == "sourceRoute" {
		action = "GetSourceRoute"
		pattern = "/sourceRoute/getSourceRoute"
	}
	reqQuery := map[string]interface{}{
		"routeId":      routeId,
		"gwInstanceId": gwInstanceId,
	}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqQuery)
	if err != nil {
		return nil, err
	}

	if data, ok := response["data"]; !ok {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("ApigwV2 route %s not found", id))
	} else if item, ok := data.(map[string]interface{}); !ok {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("ApigwV2 route %s not found", id))
	} else if _, existed := item["routeId"]; !existed {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("ApigwV2 route %s not found", id))
	} else {
		return item, nil
	}
}

func (s *ApiGateWayV2Service) DescribeMcpserver(id string) (map[string]interface{}, error) {
	// The resource ID is composed of gwInstanceId and name, separated by ":"
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid mcpserver id: %s", id)
	}
	gwInstanceId := parts[0]
	name := parts[1]

	reqBody := map[string]interface{}{
		"mcpServerName": name,
		"gwInstanceId":  gwInstanceId,
	}

	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetMcpServer", "/mcpServer/getMcpServer", nil, nil, reqBody)
	if err != nil {
		return nil, err
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Mcpserver %s not found", id))
	}

	return data, nil
}
