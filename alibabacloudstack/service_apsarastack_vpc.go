package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type VpcService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *VpcService) DoVpcDescribeforwardtableentriesRequest(id string) (entry vpc.ForwardTableEntry, err error) {
	return s.DescribeForwardEntry(id)
}

func (s *VpcService) DoVpcDescribesnattableentriesRequest(id string) (snat vpc.SnatTableEntry, err error) {
	return s.DescribeSnatEntry(id)
}

func (s *VpcService) DoVpcDescribeeipaddressesRequest(id string) (eip vpc.EipAddress, err error) {
	return s.DescribeEip(id)
}

func (s *VpcService) DescribeEip(id string) (eip vpc.EipAddress, err error) {
	request := vpc.CreateDescribeEipAddressesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.AllocationId = id
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeEipAddresses(request)
	})
	bresponse, ok := raw.(*vpc.DescribeEipAddressesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return eip, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_eip", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(bresponse.EipAddresses.EipAddress) <= 0 || bresponse.EipAddresses.EipAddress[0].AllocationId != id {
		return eip, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Eip", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	eip = bresponse.EipAddresses.EipAddress[0]
	return
}

func (s *VpcService) DescribeEipAssociation(id string) (object vpc.EipAddress, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	object, err = s.DescribeEip(parts[0])
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	if object.InstanceId != parts[1] {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Eip Association", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return
}

func (s *VpcService) DescribeNatGateway(id string) (nat vpc.NatGateway, err error) {
	request := vpc.CreateDescribeNatGatewaysRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.NatGatewayId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		})
		bresponse, ok := raw.(*vpc.DescribeNatGatewaysResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"InvalidNatGatewayId.NotFound"}) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_nat_gateway", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(bresponse.NatGateways.NatGateway) <= 0 || bresponse.NatGateways.NatGateway[0].NatGatewayId != id {
			return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NatGateway", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		nat = bresponse.NatGateways.NatGateway[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeVpc(id string) (v vpc.Vpc, err error) {
	request := vpc.CreateDescribeVpcsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.VpcId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpcs(request)
		})
		bresponse, ok := raw.(*vpc.DescribeVpcsResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vpc", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(bresponse.Vpcs.Vpc) < 1 || bresponse.Vpcs.Vpc[0].VpcId != id {
			return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VPC", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		v = bresponse.Vpcs.Vpc[0]
		return nil
	})
	return
}

func (s *VpcService) VpcStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpc(id)
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

func (s *VpcService) DoVpcDescribevswitchattributesRequest(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	return s.DescribeVSwitch(id)
}
func (s *VpcService) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.VSwitchId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		bresponse, ok := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vswitch", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if bresponse.VSwitchId != id {
			return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("vswitch", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		v = *bresponse
		return nil
	})
	return
}

func (s *VpcService) VSwitchStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVSwitch(id)
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

func (s *VpcService) DescribeSnatEntry(id string) (snat vpc.SnatTableEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return snat, errmsgs.WrapError(err)
	}
	request := vpc.CreateDescribeSnatTableEntriesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.SnatTableId = parts[0]
	request.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		invoker := NewInvoker()
		var raw interface{}
		err = invoker.Run(func() error {
			raw, err = s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeSnatTableEntries(request)
			})

			return err
		})
		response, ok := raw.(*vpc.DescribeSnatTableEntriesResponse)
		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		//so judge the snatEntries length priority
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"InvalidSnatTableId.NotFound", "InvalidSnatEntryId.NotFound"}) {
				return snat, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return snat, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_snat_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if len(response.SnatTableEntries.SnatTableEntry) < 1 {
			break
		}

		for _, snat := range response.SnatTableEntries.SnatTableEntry {
			if snat.SnatEntryId == parts[1] {
				return snat, nil
			}
		}

		if len(response.SnatTableEntries.SnatTableEntry) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return snat, errmsgs.WrapError(err)
		}
		request.PageNumber = page
	}

	return snat, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SnatEntry", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *VpcService) DescribeForwardEntry(id string) (entry vpc.ForwardTableEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return entry, errmsgs.WrapError(err)
	}
	forwardTableId, forwardEntryId := parts[0], parts[1]
	request := vpc.CreateDescribeForwardTableEntriesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ForwardTableId = forwardTableId
	request.ForwardEntryId = forwardEntryId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeForwardTableEntries(request)
		})
		bresponse, ok := raw.(*vpc.DescribeForwardTableEntriesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"InvalidForwardEntryId.NotFound", "InvalidForwardTableId.NotFound"}) {
				return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ForwardEntry", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_forward_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(bresponse.ForwardTableEntries.ForwardTableEntry) > 0 {
			entry = bresponse.ForwardTableEntries.ForwardTableEntry[0]
			return nil
		}

		return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ForwardEntry", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	})
	return
}

func (s *VpcService) QueryRouteTableById(routeTableId string) (rt vpc.RouteTable, err error) {
	request := vpc.CreateDescribeRouteTablesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = routeTableId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTables(request)
		})
		bresponse, ok := raw.(*vpc.DescribeRouteTablesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_route_table", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(bresponse.RouteTables.RouteTable) == 0 ||
			bresponse.RouteTables.RouteTable[0].RouteTableId != routeTableId {
			return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RouteTable", routeTableId)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		rt = bresponse.RouteTables.RouteTable[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeRouteEntry(id string) (*vpc.RouteEntry, error) {
	v := &vpc.RouteEntry{}
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		return v, errmsgs.WrapError(err)
	}
	rtId, cidr, nexthop_type, nexthop_id := parts[0], parts[2], parts[3], parts[4]

	request := vpc.CreateDescribeRouteTablesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = rtId

	invoker := NewInvoker()
	for {
		var raw interface{}
		var err error
		err = invoker.Run(func() error {
			raw, err = s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			return err
		})
		response, ok := raw.(*vpc.DescribeRouteTablesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return v, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_route_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(response.RouteTables.RouteTable) < 1 {
			break
		}
		for _, table := range response.RouteTables.RouteTable {
			for _, entry := range table.RouteEntrys.RouteEntry {
				if entry.DestinationCidrBlock == cidr && entry.NextHopType == nexthop_type && entry.InstanceId == nexthop_id {
					return &entry, nil
				}
			}
		}
		if len(response.RouteTables.RouteTable) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return v, errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RouteEntry", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *VpcService) DoVpcDescriberouterinterfaceattributeRequest(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	return s.DescribeRouterInterface(id, regionId)
}
func (s *VpcService) DescribeRouterInterface(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	request := vpc.CreateDescribeRouterInterfacesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RegionId = regionId
	values := []string{id}
	filter := []vpc.DescribeRouterInterfacesFilter{
		{
			Key:   "RouterInterfaceId",
			Value: &values,
		},
	}
	request.Filter = &filter
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouterInterfaces(request)
		})
		bresponse, ok := raw.(*vpc.DescribeRouterInterfacesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_router_interface", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(bresponse.RouterInterfaceSet.RouterInterfaceType) <= 0 ||
			bresponse.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != id {
			return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RouterInterface", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		ri = bresponse.RouterInterfaceSet.RouterInterfaceType[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeRouterInterfaceConnection(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	ri, err = s.DescribeRouterInterface(id, regionId)
	if err != nil {
		return ri, errmsgs.WrapError(err)
	}

	if ri.OppositeInterfaceId == "" || ri.OppositeRouterType == "" ||
		ri.OppositeRouterId == "" || ri.OppositeInterfaceOwnerId == "" {
		return ri, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RouterInterface", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return ri, nil
}

func (s *VpcService) DescribeCenInstanceGrant(id string) (rule vpc.CbnGrantRule, err error) {
	request := vpc.CreateDescribeGrantRulesToCenRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return rule, errmsgs.WrapError(err)
	}
	cenId := parts[0]
	instanceId := parts[1]
	instanceType, err := GetCenChildInstanceType(instanceId)
	if err != nil {
		return rule, errmsgs.WrapError(err)
	}

	request.InstanceId = instanceId
	request.InstanceType = instanceType

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeGrantRulesToCen(request)
		})
		bresponse, ok := raw.(*vpc.DescribeGrantRulesToCenResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cen_instance_grant", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		ruleList := bresponse.CenGrantRules.CbnGrantRule
		if len(ruleList) <= 0 {
			return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GrantRules", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}

		for ruleNum := 0; ruleNum <= len(bresponse.CenGrantRules.CbnGrantRule)-1; ruleNum++ {
			if ruleList[ruleNum].CenInstanceId == cenId {
				rule = ruleList[ruleNum]
				return nil
			}
		}

		return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GrantRules", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	})
	return
}

func (s *VpcService) WaitForCenInstanceGrant(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instanceId := parts[1]
	ownerId := parts[2]
	for {
		object, err := s.DescribeCenInstanceGrant(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.CenInstanceId == instanceId && fmt.Sprint(object.CenOwnerId) == ownerId && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.CenInstanceId, instanceId, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *VpcService) DoVpcDescribecommonbandwidthpackagesRequest(id string) (v vpc.CommonBandwidthPackage, err error) {
	return s.DescribeCommonBandwidthPackage(id)
}

func (s *VpcService) DescribeCommonBandwidthPackage(id string) (v vpc.CommonBandwidthPackage, err error) {
	request := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.BandwidthPackageId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCommonBandwidthPackages(request)
		})
		bresponse, ok := raw.(*vpc.DescribeCommonBandwidthPackagesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_common_bandwidth_package", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//Finding the commonBandwidthPackageId
		for _, bandPackage := range bresponse.CommonBandwidthPackages.CommonBandwidthPackage {
			if bandPackage.BandwidthPackageId == id {
				v = bandPackage
				return nil
			}
		}
		return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CommonBandWidthPackage", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	})
	return
}

func (s *VpcService) DescribeCommonBandwidthPackageAttachment(id string) (v vpc.CommonBandwidthPackage, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return v, errmsgs.WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]

	object, err := s.DescribeCommonBandwidthPackage(bandwidthPackageId)
	if err != nil {
		return v, errmsgs.WrapError(err)
	}

	for _, ipAddresse := range object.PublicIpAddresses.PublicIpAddresse {
		if ipAddresse.AllocationId == ipInstanceId {
			v = object
			return
		}
	}
	return v, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CommonBandWidthPackageAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *VpcService) DoVpcDescriberoutetablelistRequest(id string) (v vpc.RouterTableListType, err error) {
	return s.DescribeRouteTable(id)
}

func (s *VpcService) DescribeRouteTable(id string) (v vpc.RouterTableListType, err error) {
	request := vpc.CreateDescribeRouteTableListRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTableList(request)
		})
		bresponse, ok := raw.(*vpc.DescribeRouteTableListResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_route_table", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//Finding the routeTableId
		for _, routerTableType := range bresponse.RouterTableList.RouterTableListType {
			if routerTableType.RouteTableId == id {
				v = routerTableType
				return nil
			}
		}
		return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RouteTable", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	})
	return v, errmsgs.WrapError(err)
}

func (s *VpcService) DescribeRouteTableAttachment(id string) (v vpc.RouterTableListType, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return v, errmsgs.WrapError(err)
	}
	invoker := NewInvoker()
	routeTableId := parts[0]
	vSwitchId := parts[1]

	err = invoker.Run(func() error {
		object, err := s.DescribeRouteTable(routeTableId)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		for _, id := range object.VSwitchIds.VSwitchId {
			if id == vSwitchId {
				v = object
				return nil
			}
		}
		return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RouteTableAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	})
	return v, errmsgs.WrapError(err)
}

func (s *VpcService) WaitForVSwitch(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVSwitch(id)
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

func (s *VpcService) WaitForNatGateway(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNatGateway(id)
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

func (s *VpcService) WaitForRouteEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouteEntry(id)
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
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForAllRouteEntriesAvailable(routeTableId string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		table, err := s.QueryRouteTableById(routeTableId)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		success := true
		for _, routeEntry := range table.RouteEntrys.RouteEntry {
			if routeEntry.Status != string(Available) {
				success = false
				break
			}
		}
		if success {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, routeTableId, GetFunc(1), timeout, Available, Null, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *VpcService) WaitForRouterInterface(id, regionId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouterInterface(id, regionId)
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

func (s *VpcService) WaitForRouterInterfaceConnection(id, regionId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouterInterfaceConnection(id, regionId)
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

func (s *VpcService) WaitForEip(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEip(id)
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

func (s *VpcService) WaitForEipAssociation(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEipAssociation(id)
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

func (s *VpcService) DeactivateRouterInterface(interfaceId string) error {
	request := vpc.CreateDeactivateRouterInterfaceRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RouterInterfaceId = interfaceId

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeactivateRouterInterface(request)
	})
	bresponse, ok := raw.(*vpc.DeactivateRouterInterfaceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_router_interface", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *VpcService) ActivateRouterInterface(interfaceId string) error {
	request := vpc.CreateActivateRouterInterfaceRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RouterInterfaceId = interfaceId
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ActivateRouterInterface(request)
	})
	bresponse, ok := raw.(*vpc.ActivateRouterInterfaceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_router_interface", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *VpcService) WaitForForwardEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeForwardEntry(id)
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

func (s *VpcService) WaitForSnatEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeSnatEntry(id)
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

	}
}

func (s *VpcService) WaitForCommonBandwidthPackage(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCommonBandwidthPackage(id)
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
	}
}

func (s *VpcService) WaitForCommonBandwidthPackageAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCommonBandwidthPackageAttachment(id)
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
	}
}

// Flattens an array of vpc.public_ip_addresses into a []map[string]string
func (s *VpcService) FlattenPublicIpAddressesMappings(list []vpc.PublicIpAddresse) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := map[string]interface{}{
			"ip_address":    i.IpAddress,
			"allocation_id": i.AllocationId,
		}
		result = append(result, l)
	}

	return result
}

func (s *VpcService) WaitForRouteTable(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouteTable(id)
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
		time.Sleep(3 * time.Second)
	}
}

func (s *VpcService) WaitForRouteTableAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouteTableAttachment(id)
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
		time.Sleep(3 * time.Second)
	}
}

func (s *VpcService) DoVpcDescribenetworkaclattributesRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeNetworkAcl(id)
}
func (s *VpcService) DescribeNetworkAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	request := map[string]interface{}{

		"NetworkAclId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", "DescribeNetworkAclAttributes", "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidNetworkAcl.NotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VPC:NetworkAcl", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	v, err := jsonpath.Get("$.NetworkAclAttribute", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.NetworkAclAttribute", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeNetworkAclAttachment(id string, resource []vpc.Resource) (err error) {

	invoker := NewInvoker()
	return invoker.Run(func() error {
		object, err := s.DescribeNetworkAcl(id)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		if len(resources) < 1 {
			return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Network Acl Attachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		success := true
		for _, source := range resources {
			success = false
			for _, res := range resource {
				item := source.(map[string]interface{})
				if fmt.Sprint(item["ResourceId"]) == res.ResourceId {
					success = true
				}
			}
			if success == false {
				return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Network Acl Attachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			}
		}
		return nil
	})
}

func (s *VpcService) WaitForNetworkAcl(networkAclId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNetworkAcl(networkAclId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		success := true
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		// Check Acl's binding resources
		for _, res := range resources {
			item := res.(map[string]interface{})
			if fmt.Sprint(item["Status"]) != string(BINDED) {
				success = false
			}
		}
		if fmt.Sprint(object["Status"]) == string(status) && success == true {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, networkAclId, GetFunc(1), timeout, fmt.Sprint(object["Status"]), string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForNetworkAclAttachment(id string, resource []vpc.Resource, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		err := s.DescribeNetworkAclAttachment(id, resource)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		object, err := s.DescribeNetworkAcl(id)
		success := true
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		// Check Acl's binding resources
		for _, res := range resources {
			item := res.(map[string]interface{})
			if fmt.Sprint(item["Status"]) != string(BINDED) {
				success = false
			}
		}
		if fmt.Sprint(object["Status"]) == string(status) && success == true {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(object["Status"]), string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []vpc.TagResource, err error) {
	request := vpc.CreateListTagResourcesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)

	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []vpc.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, vpc.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ListTagResources(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, resourceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		return
	}
	response, _ := raw.(*vpc.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *VpcService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := vpc.CreateUnTagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = string(resourceType)
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err := s.client.WithVpcClient(func(client *vpc.Client) (interface{}, error) {
					return client.UnTagResources(request)
				})
				if err != nil {
					if errmsgs.IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
			}
		}

		if len(create) > 0 {
			request := vpc.CreateTagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = string(resourceType)
			request.RegionId = s.client.RegionId

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err := s.client.WithVpcClient(func(client *vpc.Client) (interface{}, error) {
					return client.TagResources(request)
				})
				if err != nil {
					if errmsgs.IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
			}
		}

		//d.SetPartial("tags")
	}

	return nil
}

func (s *VpcService) tagsToMap(tags []vpc.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *VpcService) ignoreTag(t vpc.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found AlibabacloudStack Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *VpcService) tagToMap(tags []vpc.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.vpcTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}
func (s *VpcService) vpcTagIgnored(t vpc.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found AlibabacloudStack Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *VpcService) diffTags(oldTags, newTags []vpc.TagResourcesTag) ([]vpc.TagResourcesTag, []vpc.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []vpc.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *VpcService) tagsFromMap(m map[string]interface{}) []vpc.TagResourcesTag {
	result := make([]vpc.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, vpc.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *VpcService) DescribeVSwitchWithTeadsl(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"VSwitchId": id,
	}
	response, err := s.client.DoTeaRequest("POST", "VPC", "2016-04-28", "DescribeVSwitchAttributes", "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, err
	}
	if v, ok := response["VSwitchId"].(string); ok && v != id {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("vswitch", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return response, nil
}

func (s *VpcService) NetworkAclStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNetworkAcl(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DeleteAclResources(id string) (object map[string]interface{}, err error) {
	acl, err := s.DescribeNetworkAcl(id)
	if err != nil {
		return object, errmsgs.WrapError(err)
	}
	deleteResources := make([]map[string]interface{}, 0)
	res, err := jsonpath.Get("$.Resources.Resource", acl)
	if err != nil {
		return object, errmsgs.WrapError(err)
	}
	resources, _ := res.([]interface{})
	if resources != nil && len(resources) < 1 {
		return object, nil
	}
	for _, val := range resources {
		item, _ := val.(map[string]interface{})
		deleteResources = append(deleteResources, map[string]interface{}{
			"ResourceId":   item["ResourceId"],
			"ResourceType": item["ResourceType"],
		})
	}

	var response map[string]interface{}
	request := map[string]interface{}{
		"NetworkAclId": id,
		"Resource":     deleteResources,
	}
	request["ClientToken"] = buildClientToken("UnassociateNetworkAcl")
	response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", "UnassociateNetworkAcl", "", nil, request)
	if err != nil {
		return response, err
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, 10*time.Minute, 5*time.Second, s.NetworkAclStateRefreshFunc(id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return response, errmsgs.WrapErrorf(err, errmsgs.IdMsg, id)
	}
	return object, nil
}

func (s *VpcService) DescribeVpcIpv6EgressRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"PageSize":      PageSizeLarge,
		"PageNumber":    1,
		"Ipv6GatewayId": parts[0],
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
		runtime.SetAutoretry(true)
		response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", "DescribeIpv6EgressOnlyRules", "", nil, request)
		if err != nil {
			return object, err
		}
		v, err := jsonpath.Get("$.Ipv6EgressOnlyRules.Ipv6EgressOnlyRule", response)
		if err != nil {
			return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Ipv6EgressOnlyRules.Ipv6EgressOnlyRule", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VPC", id)), errmsgs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Ipv6EgressOnlyRuleId"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VPC", id)), errmsgs.NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpcIpv6EgressRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpv6EgressRule(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		for _, failState := range failStates {

			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcIpv6Gateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"Ipv6GatewayId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", "DescribeIpv6GatewayAttribute", "", nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if Ipv6GatewayId, ok := object["Ipv6GatewayId"]; !ok || fmt.Sprint(Ipv6GatewayId) == "" {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VPC", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *VpcService) VpcIpv6GatewayStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpv6Gateway(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcIpv6InternetBandwidth(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{

		"Ipv6InternetBandwidthId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", "DescribeIpv6Addresses", "", nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Ipv6Addresses.Ipv6Address", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Ipv6Addresses.Ipv6Address", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VPC", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["Ipv6InternetBandwidth"].(map[string]interface{})["Ipv6InternetBandwidthId"]) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("VPC", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) setInstanceSecondaryCidrBlocks(d *schema.ResourceData) error {
	var response map[string]interface{}
	var err error
	if d.HasChange("secondary_cidr_blocks") {
		oraw, nraw := d.GetChange("secondary_cidr_blocks")
		removed := oraw.([]interface{})
		added := nraw.([]interface{})
		if len(removed) > 0 {
			action := "UnassociateVpcCidrBlock"
			request := map[string]interface{}{
				"VpcId": d.Id(),
			}
			for _, item := range removed {
				request["SecondaryCidrBlock"] = item
				response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", "UnassociateVpcCidrBlock", "", nil, request)
				if err != nil {
					return err
				}
				addDebug(action, response, request)
			}
		}

		if len(added) > 0 {
			request := map[string]interface{}{
				"VpcId": d.Id(),
			}
			for _, item := range added {
				request["SecondaryCidrBlock"] = item
				response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", "AssociateVpcCidrBlock", "", nil, request)
				if err != nil {
					return err
				}
			}
		}
		//d.SetPartial("secondary_cidr_blocks")
	}
	return nil
}

func (s *VpcService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{

				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			_, err := s.client.DoTeaRequest("POST", "VPC", "2016-04-28", action, "", nil, request)
			if err != nil {
				return err
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{

				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			_, err := s.client.DoTeaRequest("POST", "VPC", "2016-04-28", action, "", nil, request)
			if err != nil {
				return err
			}
		}
		//d.SetPartial("tags")
	}
	return nil
}

func (s *VpcService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	action := "ListTagResources"
	request := map[string]interface{}{
		"ResourceType":   resourceType,
		"ResourceId.1":   id,
		"Product":        "Vpc",
		"OrganizationId": s.client.Department,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		_, err := s.client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, request)
		if err != nil {
			return nil, err
		}
		v, err := jsonpath.Get("$.TagResources.TagResource", response)
		if err != nil {
			return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response)
		}
		if v != nil {
			tags = append(tags, v.([]interface{})...)
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *VpcService) DoVpcDescribephysicalconnectionsRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeExpressConnectPhysicalConnection(id)
}

func (s *VpcService) DescribeExpressConnectPhysicalConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribePhysicalConnections"
	request := map[string]interface{}{
		"Product":        "Vpc",
		"OrganizationId": s.client.Department,
	}
	filterMapList := make([]map[string]interface{}, 0)
	filterMapList = append(filterMapList, map[string]interface{}{
		"Key":   "PhysicalConnectionId",
		"Value": []string{id},
	})
	request["Filter"] = filterMapList
	response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", action, "", nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.PhysicalConnectionSet.PhysicalConnectionType", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.PhysicalConnectionSet.PhysicalConnectionType", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ExpressConnect", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PhysicalConnectionId"]) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ExpressConnect", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) ExpressConnectPhysicalConnectionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectPhysicalConnection(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeExpressConnectVirtualBorderRouter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeVirtualBorderRouters"
	request := map[string]interface{}{
		"PageNumber":     1,
		"PageSize":       50,
		"Product":        "Vpc",
		"OrganizationId": s.client.Department,
	}
	idExist := false
	for {
		response, err = s.client.DoTeaRequest("POST", "VPC", "2016-04-28", action, "", nil, request)
		if err != nil {
			return object, err
		}
		v, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		if err != nil {
			return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ExpressConnect", id)), errmsgs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["VbrId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ExpressConnect", id)), errmsgs.NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) ExpressConnectVirtualBorderRouterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectVirtualBorderRouter(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
