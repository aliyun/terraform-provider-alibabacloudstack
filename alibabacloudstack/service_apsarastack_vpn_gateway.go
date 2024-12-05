package alibabacloudstack

import (
	"time"

	"strings"

	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
