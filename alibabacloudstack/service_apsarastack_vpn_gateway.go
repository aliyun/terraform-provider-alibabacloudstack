package alibabacloudstack

import (
	"strconv"
	"time"

	"strings"

	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type VpnGatewayService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *VpnGatewayService) DoVpcDescribevpngatewayRequest(id string) (v vpc.DescribeVpnGatewayResponse, err error) {
	return s.DescribeVpnGateway(id)
}
func (s *VpnGatewayService) DescribeVpnGateway(id string) (v vpc.DescribeVpnGatewayResponse, err error) {
	request := vpc.CreateDescribeVpnGatewayRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.VpnGatewayId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpnGateway(request)
	})
	response, ok := raw.(*vpc.DescribeVpnGatewayResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden", "InvalidVpnGatewayInstanceId.NotFound"}) {
			return v, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return v, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.VpnGatewayId != id {
		return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VpnGateway", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return *response, nil
}

func (s *VpnGatewayService) DoVpcDescribecustomergatewayRequest(id string) (v vpc.DescribeCustomerGatewayResponse, err error) {
	return s.DescribeVpnCustomerGateway(id)
}
func (s *VpnGatewayService) DescribeVpnCustomerGateway(id string) (v vpc.DescribeCustomerGatewayResponse, err error) {
	request := vpc.CreateDescribeCustomerGatewayRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.CustomerGatewayId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeCustomerGateway(request)
	})
	response, ok := raw.(*vpc.DescribeCustomerGatewayResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden", "InvalidCustomerGatewayInstanceId.NotFound"}) {
			return v, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return v, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.CustomerGatewayId != id {
		return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VpnCustomerGateway", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return *response, nil
}

func (s *VpnGatewayService) DescribeVpnConnection(id string) (v vpc.DescribeVpnConnectionResponse, err error) {
	request := vpc.CreateDescribeVpnConnectionRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.VpnConnectionId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpnConnection(request)
	})
	response, ok := raw.(*vpc.DescribeVpnConnectionResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden", "InvalidVpnConnectionInstanceId.NotFound"}) {
			return v, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return v, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.VpnConnectionId != id {
		return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VpnConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return *response, nil
}

func (s *VpnGatewayService) DescribeSslVpnServer(id string) (v vpc.SslVpnServer, err error) {
	request := vpc.CreateDescribeSslVpnServersRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.SslVpnServerId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeSslVpnServers(request)
	})
	response, ok := raw.(*vpc.DescribeSslVpnServersResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden", "InvalidSslVpnServerId.NotFound"}) {
			return v, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return v, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(response.SslVpnServers.SslVpnServer) == 0 || response.SslVpnServers.SslVpnServer[0].SslVpnServerId != id {
		return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SslVpnGateway", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return response.SslVpnServers.SslVpnServer[0], nil
}

func (s *VpnGatewayService) DescribeSslVpnClientCert(id string) (v vpc.DescribeSslVpnClientCertResponse, err error) {
	request := vpc.CreateDescribeSslVpnClientCertRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.SslVpnClientCertId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeSslVpnClientCert(request)
	})
	response, ok := raw.(*vpc.DescribeSslVpnClientCertResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden", "InvalidSslVpnClientCertId.NotFound"}) {
			return v, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return v, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.SslVpnClientCertId != id {
		return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SslVpnClientCert", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return *response, nil
}

func (s *VpnGatewayService) DescribeVpnRouteEntry(id string) (v vpc.VpnRouteEntry, err error) {
	request := vpc.CreateDescribeVpnRouteEntriesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return v, errmsgs.WrapError(err)
	}
	gatewayId := parts[0]

	request.VpnGatewayId = gatewayId
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpnRouteEntries(request)
	})
	response, ok := raw.(*vpc.DescribeVpnRouteEntriesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden", "InvalidVpnGatewayInstanceId.NotFound"}) {
			return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VpnRouterEntry", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return v, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	for _, routeEntry := range response.VpnRouteEntries.VpnRouteEntry {
		if id == gatewayId+":"+routeEntry.NextHop+":"+routeEntry.RouteDest {
			return routeEntry, nil
		}
	}
	return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VpnRouterEntry", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *VpnGatewayService) WaitForVpnGateway(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVpnGateway(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if strings.EqualFold(object.Status, string(status)) {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpnGatewayService) WaitForVpnConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVpnConnection(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if status != Deleted && object.VpnConnectionId == id {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpnGatewayService) WaitForVpnCustomerGateway(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVpnCustomerGateway(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.CustomerGatewayId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpnGatewayService) WaitForSslVpnServer(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSslVpnServer(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.SslVpnServerId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpnGatewayService) WaitForSslVpnClientCert(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSslVpnClientCert(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpnGatewayService) WaitForVpnRouteEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVpnRouteEntry(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		parts, err := ParseResourceId(id, 3)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		if object.NextHop == parts[1] && object.RouteDest == parts[2] && string(status) != string(Deleted) {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpnGatewayService) ParseIkeConfig(ike vpc.IkeConfig) (ikeConfigs []map[string]interface{}) {
	item := map[string]interface{}{
		"ike_auth_alg":  ike.IkeAuthAlg,
		"ike_enc_alg":   ike.IkeEncAlg,
		"ike_lifetime":  ike.IkeLifetime,
		"ike_local_id":  ike.LocalId,
		"ike_mode":      ike.IkeMode,
		"ike_pfs":       ike.IkePfs,
		"ike_remote_id": ike.RemoteId,
		"ike_version":   ike.IkeVersion,
		"psk":           ike.Psk,
	}

	ikeConfigs = append(ikeConfigs, item)
	return
}

func (s *VpnGatewayService) ParseIpsecConfig(ipsec vpc.IpsecConfig) (ipsecConfigs []map[string]interface{}) {
	item := map[string]interface{}{
		"ipsec_auth_alg": ipsec.IpsecAuthAlg,
		"ipsec_enc_alg":  ipsec.IpsecEncAlg,
		"ipsec_lifetime": ipsec.IpsecLifetime,
		"ipsec_pfs":      ipsec.IpsecPfs,
	}

	ipsecConfigs = append(ipsecConfigs, item)
	return
}

func (s *VpnGatewayService) AssembleIkeConfig(ikeCfgParam []interface{}) (string, error) {
	var ikeCfg IkeConfig
	v := ikeCfgParam[0]
	item := v.(map[string]interface{})
	ikeCfg = IkeConfig{
		IkeAuthAlg:  item["ike_auth_alg"].(string),
		IkeEncAlg:   item["ike_enc_alg"].(string),
		IkeLifetime: item["ike_lifetime"].(int),
		LocalId:     item["ike_local_id"].(string),
		IkeMode:     item["ike_mode"].(string),
		IkePfs:      item["ike_pfs"].(string),
		RemoteId:    item["ike_remote_id"].(string),
		IkeVersion:  item["ike_version"].(string),
		Psk:         item["psk"].(string),
	}

	data, err := json.Marshal(ikeCfg)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(data), nil
}

func (s *VpnGatewayService) AssembleIpsecConfig(ipsecCfgParam []interface{}) (string, error) {
	var ipsecCfg IpsecConfig
	v := ipsecCfgParam[0]
	item := v.(map[string]interface{})
	ipsecCfg = IpsecConfig{
		IpsecAuthAlg:  item["ipsec_auth_alg"].(string),
		IpsecEncAlg:   item["ipsec_enc_alg"].(string),
		IpsecLifetime: item["ipsec_lifetime"].(int),
		IpsecPfs:      item["ipsec_pfs"].(string),
	}

	data, err := json.Marshal(ipsecCfg)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(data), nil
}

func (s *VpnGatewayService) AssembleNetworkSubnetToString(list []interface{}) string {
	if len(list) < 1 {
		return ""
	}
	var items []string
	for _, id := range list {
		items = append(items, fmt.Sprintf("%s", id))
	}
	return fmt.Sprintf("%s", strings.Join(items, COMMA_SEPARATED))
}

func TimestampToStr(timestamp int64) string {
	tm := time.Unix(timestamp/1000, 0)
	timeString := tm.Format("2006-01-02T15:04:05Z")
	return timeString
}

type VpnGatewayVpnPbrRouteEntry struct {
	VpnInstanceId string `json:"VpnInstanceId"`
	RouteSource   string `json:"RouteSource"`
	RouteDest     string `json:"RouteDest"`
	NextHop       string `json:"NextHop"`
	Weight        int    `json:"Weight"`
	CreateTime    int    `json:"CreateTime"`
	State         string `json:"State"`
}

type VpcDescribevpnpbrrouteentriesResponse struct {
	VpnPbrRouteEntries struct {
		VpnPbrRouteEntry []VpnGatewayVpnPbrRouteEntry `json:"VpnPbrRouteEntry"`
	} `json:"VpnPbrRouteEntries"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
}

func (s *VpnGatewayService) DoVpcDescribevpnpbrrouteentriesRequest(id string) (*VpnGatewayVpnPbrRouteEntry, error) {
	// api: Vpc - 2016-04-28 - DescribeVpnPbrRouteEntries
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeVpnPbrRouteEntries", "")
	VpcDescribevpnpbrrouteentriesResponseObj := &VpcDescribevpnpbrrouteentriesResponse{}
	result := &VpnGatewayVpnPbrRouteEntry{}
	param := strings.Split(id, "_")
	request.QueryParams["VpnGatewayId"] = param[0]
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeVpnPbrRouteEntries", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribevpnpbrrouteentriesResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeVpnPbrRouteEntries", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, data := range VpcDescribevpnpbrrouteentriesResponseObj.VpnPbrRouteEntries.VpnPbrRouteEntry {
		if data.NextHop == param[3] && data.RouteDest == param[2] && data.RouteSource == param[1] && data.VpnInstanceId == param[0] {
			result = &data
		}
	}
	if result == nil {
		return nil, errmsgs.Error(errmsgs.NotFoundMsg, "VpnPbrRouteEntry")
	}

	return result, nil
}

func (s *VpnGatewayService) VpnGatewayStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnGateway(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}
	
type VpcDescribesslvpnserversResponse struct {
	SslVpnServers struct {
		SslVpnServer []struct {
			RegionId       string `json:"RegionId"`
			SslVpnServerId string `json:"SslVpnServerId"`
			VpnGatewayId   string `json:"VpnGatewayId"`
			Name           string `json:"Name"`
			LocalSubnet    string `json:"LocalSubnet"`
			ClientIpPool   string `json:"ClientIpPool"`
			CreateTime     int    `json:"CreateTime"`
			Cipher         string `json:"Cipher"`
			Proto          string `json:"Proto"`
			Port           int    `json:"Port"`
			Compress       bool   `json:"Compress"`
			Connections    int    `json:"Connections"`
			MaxConnections int    `json:"MaxConnections"`
			InternetIp     string `json:"InternetIp"`
			AscmCreateUser string `json:"AscmCreateUser"`
		} `json:"SslVpnServer"`
	} `json:"SslVpnServers"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
}

func (s *VpnGatewayService) DoVpcDescribesslvpnserversRequest(id string) (*VpcDescribesslvpnserversResponse, error) {
	// api: Vpc - 2016-04-28 - DescribeSslVpnServers
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeSslVpnServers", "")
	VpcDescribesslvpnserversResponseObj := &VpcDescribesslvpnserversResponse{}

	//调用request_params_handler
	request.QueryParams["SslVpnServerId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeSslVpnServers", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribesslvpnserversResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeSslVpnServers", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return VpcDescribesslvpnserversResponseObj, nil
}

type VpcDescribevpnsslserverlogsResponse struct {
	Data struct {
		Logs []string `json:"Logs"`
	} `json:"Data"`
	RequestId   string `json:"RequestId"`
	Count       int    `json:"Count"`
	IsCompleted bool   `json:"IsCompleted"`
	PageNumber  int    `json:"PageNumber"`
	PageSize    int    `json:"PageSize"`
}

func (s *VpnGatewayService) DoVpcDescribevpnsslserverlogsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*VpcDescribevpnsslserverlogsResponse, error) {
	// api: Vpc - 2016-04-28 - DescribeVpnSslServerLogs
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeVpnSslServerLogs", "")
	VpcDescribevpnsslserverlogsResponseObj := &VpcDescribevpnsslserverlogsResponse{}

	//调用request_params_handler

	if v, ok := d.GetOk("region_id"); ok {
		request.QueryParams["RegionId"] = v.(string)
	} else {
		return nil, fmt.Errorf("RegionId is required")
	}

	if v, ok := d.GetOk("ssl_vpn_server_id"); ok {
		request.QueryParams["VpnSslServerId"] = v.(string)
	} else {
		return nil, fmt.Errorf("SslVpnServerId is required")
	}

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeVpnSslServerLogs", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribevpnsslserverlogsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeVpnSslServerLogs", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return VpcDescribevpnsslserverlogsResponseObj, nil
}

type VpcDescribesslvpnclientcertResponse struct {
	RequestId          string `json:"RequestId"`
	RegionId           string `json:"RegionId"`
	SslVpnClientCertId string `json:"SslVpnClientCertId"`
	Name               string `json:"Name"`
	SslVpnServerId     string `json:"SslVpnServerId"`
	CaCert             string `json:"CaCert"`
	ClientCert         string `json:"ClientCert"`
	ClientKey          string `json:"ClientKey"`
	ClientConfig       string `json:"ClientConfig"`
	CreateTime         int    `json:"CreateTime"`
	EndTime            int    `json:"EndTime"`
	Status             string `json:"Status"`
}

func (s *VpnGatewayService) DoVpcDescribesslvpnclientcertRequest(id string) (*VpcDescribesslvpnclientcertResponse, error) {
	// api: Vpc - 2016-04-28 - DescribeSslVpnClientCert
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeSslVpnClientCert", "")
	VpcDescribesslvpnclientcertResponseObj := &VpcDescribesslvpnclientcertResponse{}

	//调用request_params_handler

	request.QueryParams["SslVpnClientCertId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeSslVpnClientCert", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribesslvpnclientcertResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeSslVpnClientCert", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return VpcDescribesslvpnclientcertResponseObj, nil
}

type VpcDescribesslvpnclientcertsResponse struct {
	SslVpnClientCertKeys struct {
		SslVpnClientCertKey []struct {
			RegionId           string `json:"RegionId"`
			SslVpnClientCertId string `json:"SslVpnClientCertId"`
			Name               string `json:"Name"`
			SslVpnServerId     string `json:"SslVpnServerId"`
			CreateTime         int    `json:"CreateTime"`
			EndTime            int    `json:"EndTime"`
			Status             string `json:"Status"`
		} `json:"SslVpnClientCertKey"`
	} `json:"SslVpnClientCertKeys"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
}

func (s *VpnGatewayService) DoVpcDescribesslvpnclientcertsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*VpcDescribesslvpnclientcertsResponse, error) {
	// api: Vpc - 2016-04-28 - DescribeSslVpnClientCerts
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeSslVpnClientCerts", "")
	VpcDescribesslvpnclientcertsResponseObj := &VpcDescribesslvpnclientcertsResponse{}

	//调用request_params_handler

	if v, ok := d.GetOk("page_number"); ok {
		request.QueryParams["PageNumber"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("page_size"); ok {
		request.QueryParams["PageSize"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("region_id"); ok {
		request.QueryParams["RegionId"] = v.(string)
	} else {
		return nil, fmt.Errorf("RegionId is required")
	}

	if v, ok := d.GetOk("ssl_vpn_client_cert_id"); ok {
		request.QueryParams["SslVpnClientCertId"] = v.(string)
	}

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeSslVpnClientCerts", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribesslvpnclientcertsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeSslVpnClientCerts", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return VpcDescribesslvpnclientcertsResponseObj, nil
}
