package alibabacloudstack

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

type VpcService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *VpcService) DescribeEip(id string) (eip vpc.EipAddress, err error) {

	request := vpc.CreateDescribeEipAddressesRequest()
	request.RegionId = string(s.client.Region)
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.AllocationId = id
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeEipAddresses(request)
	})
	if err != nil {
		return eip, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.DescribeEipAddressesResponse)
	if len(response.EipAddresses.EipAddress) <= 0 || response.EipAddresses.EipAddress[0].AllocationId != id {
		return eip, WrapErrorf(Error(GetNotFoundMessage("Eip", id)), NotFoundMsg, ProviderERROR)
	}
	eip = response.EipAddresses.EipAddress[0]
	return
}

func (s *VpcService) DescribeEipAssociation(id string) (object vpc.EipAddress, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	object, err = s.DescribeEip(parts[0])
	if err != nil {
		err = WrapError(err)
		return
	}
	if object.InstanceId != parts[1] {
		err = WrapErrorf(Error(GetNotFoundMessage("Eip Association", id)), NotFoundMsg, ProviderERROR)
	}

	return
}

func (s *VpcService) DescribeNatGateway(id string) (nat vpc.NatGateway, err error) {
	request := vpc.CreateDescribeNatGatewaysRequest()
	request.RegionId = string(s.client.Region)
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.NatGatewayId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidNatGatewayId.NotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeNatGatewaysResponse)
		if len(response.NatGateways.NatGateway) <= 0 || response.NatGateways.NatGateway[0].NatGatewayId != id {
			return WrapErrorf(Error(GetNotFoundMessage("NatGateway", id)), NotFoundMsg, ProviderERROR)
		}
		nat = response.NatGateways.NatGateway[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeVpc(id string) (v vpc.Vpc, err error) {
	request := vpc.CreateDescribeVpcsRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.VpcId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpcs(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVpcsResponse)
		if len(response.Vpcs.Vpc) < 1 || response.Vpcs.Vpc[0].VpcId != id {
			return WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundMsg, ProviderERROR)
		}
		v = response.Vpcs.Vpc[0]
		return nil
	})
	return
}

func (s *VpcService) VpcStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpc(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}

		return object, object.Status, nil
	}
}

func (s *VpcService) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.VSwitchId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return WrapErrorf(Error(GetNotFoundMessage("vswitch", id)), NotFoundMsg, ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}

func (s *VpcService) VSwitchStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVSwitch(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}

		return object, object.Status, nil
	}
}

func (s *VpcService) DescribeSnatEntry(id string) (snat vpc.SnatTableEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return snat, WrapError(err)
	}
	request := vpc.CreateDescribeSnatTableEntriesRequest()
	request.RegionId = string(s.client.Region)
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.SnatTableId = parts[0]

	request.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		invoker := NewInvoker()
		var response *vpc.DescribeSnatTableEntriesResponse
		var raw interface{}
		err = invoker.Run(func() error {
			raw, err = s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeSnatTableEntries(request)
			})
			response, _ = raw.(*vpc.DescribeSnatTableEntriesResponse)
			return err
		})

		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		//so judge the snatEntries length priority
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidSnatTableId.NotFound", "InvalidSnatEntryId.NotFound"}) {
				return snat, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
			}
			return snat, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
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
			return snat, WrapError(err)
		}
		request.PageNumber = page
	}

	return snat, WrapErrorf(Error(GetNotFoundMessage("SnatEntry", id)), NotFoundMsg, ProviderERROR)
}

func (s *VpcService) DescribeForwardEntry(id string) (entry vpc.ForwardTableEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return entry, WrapError(err)
	}
	forwardTableId, forwardEntryId := parts[0], parts[1]
	request := vpc.CreateDescribeForwardTableEntriesRequest()
	request.RegionId = string(s.client.Region)
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.ForwardTableId = forwardTableId
	request.ForwardEntryId = forwardEntryId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeForwardTableEntries(request)
		})
		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		//so judge the snatEntries length priority
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidForwardEntryId.NotFound", "InvalidForwardTableId.NotFound"}) {
				return WrapErrorf(Error(GetNotFoundMessage("ForwardEntry", id)), NotFoundMsg, ProviderERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, "ForwardEntry", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeForwardTableEntriesResponse)

		if len(response.ForwardTableEntries.ForwardTableEntry) > 0 {
			entry = response.ForwardTableEntries.ForwardTableEntry[0]
			return nil
		}

		return WrapErrorf(Error(GetNotFoundMessage("ForwardEntry", id)), NotFoundMsg, ProviderERROR)
	})
	return
}

func (s *VpcService) QueryRouteTableById(routeTableId string) (rt vpc.RouteTable, err error) {
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.RouteTableId = routeTableId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTables(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, routeTableId, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTablesResponse)
		if len(response.RouteTables.RouteTable) == 0 ||
			response.RouteTables.RouteTable[0].RouteTableId != routeTableId {
			return WrapErrorf(Error(GetNotFoundMessage("RouteTable", routeTableId)), NotFoundMsg, ProviderERROR)
		}
		rt = response.RouteTables.RouteTable[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeRouteEntry(id string) (*vpc.RouteEntry, error) {
	v := &vpc.RouteEntry{}
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		return v, WrapError(err)
	}
	rtId, cidr, nexthop_type, nexthop_id := parts[0], parts[2], parts[3], parts[4]

	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.RouteTableId = rtId

	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			response, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			raw = response
			return err
		}); err != nil {
			return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTablesResponse)
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
			return v, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return v, WrapErrorf(Error(GetNotFoundMessage("RouteEntry", id)), NotFoundMsg, ProviderERROR)
}

func (s *VpcService) DescribeRouterInterface(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	request := vpc.CreateDescribeRouterInterfacesRequest()
	request.RegionId = regionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

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
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouterInterfacesResponse)
		if len(response.RouterInterfaceSet.RouterInterfaceType) <= 0 ||
			response.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != id {
			return WrapErrorf(Error(GetNotFoundMessage("RouterInterface", id)), NotFoundMsg, ProviderERROR)
		}
		ri = response.RouterInterfaceSet.RouterInterfaceType[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeRouterInterfaceConnection(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	ri, err = s.DescribeRouterInterface(id, regionId)
	if err != nil {
		return ri, WrapError(err)
	}

	if ri.OppositeInterfaceId == "" || ri.OppositeRouterType == "" ||
		ri.OppositeRouterId == "" || ri.OppositeInterfaceOwnerId == "" {
		return ri, WrapErrorf(Error(GetNotFoundMessage("RouterInterface", id)), NotFoundMsg, ProviderERROR)
	}
	return ri, nil
}

func (s *VpcService) DescribeCenInstanceGrant(id string) (rule vpc.CbnGrantRule, err error) {
	request := vpc.CreateDescribeGrantRulesToCenRequest()
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return rule, WrapError(err)
	}
	cenId := parts[0]
	instanceId := parts[1]
	instanceType, err := GetCenChildInstanceType(instanceId)
	if err != nil {
		return rule, WrapError(err)
	}

	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.InstanceId = instanceId
	request.InstanceType = instanceType

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeGrantRulesToCen(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeGrantRulesToCenResponse)
		ruleList := response.CenGrantRules.CbnGrantRule
		if len(ruleList) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("GrantRules", id)), NotFoundMsg, ProviderERROR)
		}

		for ruleNum := 0; ruleNum <= len(response.CenGrantRules.CbnGrantRule)-1; ruleNum++ {
			if ruleList[ruleNum].CenInstanceId == cenId {
				rule = ruleList[ruleNum]
				return nil
			}
		}

		return WrapErrorf(Error(GetNotFoundMessage("GrantRules", id)), NotFoundMsg, ProviderERROR)
	})
	return
}

func (s *VpcService) WaitForCenInstanceGrant(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[1]
	ownerId := parts[2]
	for {
		object, err := s.DescribeCenInstanceGrant(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.CenInstanceId == instanceId && fmt.Sprint(object.CenOwnerId) == ownerId && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.CenInstanceId, instanceId, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *VpcService) DescribeCommonBandwidthPackage(id string) (v vpc.CommonBandwidthPackage, err error) {
	request := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.BandwidthPackageId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCommonBandwidthPackages(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeCommonBandwidthPackagesResponse)
		//Finding the commonBandwidthPackageId
		for _, bandPackage := range response.CommonBandwidthPackages.CommonBandwidthPackage {
			if bandPackage.BandwidthPackageId == id {
				v = bandPackage
				return nil
			}
		}
		return WrapErrorf(Error(GetNotFoundMessage("CommonBandWidthPackage", id)), NotFoundMsg, ProviderERROR)
	})
	return
}

func (s *VpcService) DescribeCommonBandwidthPackageAttachment(id string) (v vpc.CommonBandwidthPackage, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return v, WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]

	object, err := s.DescribeCommonBandwidthPackage(bandwidthPackageId)
	if err != nil {
		return v, WrapError(err)
	}

	for _, ipAddresse := range object.PublicIpAddresses.PublicIpAddresse {
		if ipAddresse.AllocationId == ipInstanceId {
			v = object
			return
		}
	}
	return v, WrapErrorf(Error(GetNotFoundMessage("CommonBandWidthPackageAttachment", id)), NotFoundMsg, ProviderERROR)
}

func (s *VpcService) DescribeRouteTable(id string) (v vpc.RouterTableListType, err error) {
	request := vpc.CreateDescribeRouteTableListRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.RouteTableId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTableList(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTableListResponse)
		//Finding the routeTableId
		for _, routerTableType := range response.RouterTableList.RouterTableListType {
			if routerTableType.RouteTableId == id {
				v = routerTableType
				return nil
			}
		}
		return WrapErrorf(Error(GetNotFoundMessage("RouteTable", id)), NotFoundMsg, ProviderERROR)
	})
	return v, WrapError(err)
}

func (s *VpcService) DescribeRouteTableAttachment(id string) (v vpc.RouterTableListType, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return v, WrapError(err)
	}
	invoker := NewInvoker()
	routeTableId := parts[0]
	vSwitchId := parts[1]

	err = invoker.Run(func() error {
		object, err := s.DescribeRouteTable(routeTableId)
		if err != nil {
			return WrapError(err)
		}

		for _, id := range object.VSwitchIds.VSwitchId {
			if id == vSwitchId {
				v = object
				return nil
			}
		}
		return WrapErrorf(Error(GetNotFoundMessage("RouteTableAttachment", id)), NotFoundMsg, ProviderERROR)
	})
	return v, WrapError(err)
}

func (s *VpcService) WaitForVSwitch(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVSwitch(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForNatGateway(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNatGateway(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForRouteEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouteEntry(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForAllRouteEntriesAvailable(routeTableId string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		table, err := s.QueryRouteTableById(routeTableId)
		if err != nil {
			return WrapError(err)
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
			return WrapErrorf(err, WaitTimeoutMsg, routeTableId, GetFunc(1), timeout, Available, Null, ProviderERROR)
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
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForRouterInterfaceConnection(id, regionId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouterInterfaceConnection(id, regionId)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForEip(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEip(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForEipAssociation(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEipAssociation(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) DeactivateRouterInterface(interfaceId string) error {
	request := vpc.CreateDeactivateRouterInterfaceRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.RouterInterfaceId = interfaceId

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeactivateRouterInterface(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "RouterInterface", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *VpcService) ActivateRouterInterface(interfaceId string) error {
	request := vpc.CreateActivateRouterInterfaceRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}

	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.RouterInterfaceId = interfaceId
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ActivateRouterInterface(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "RouterInterface", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *VpcService) WaitForForwardEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeForwardEntry(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForSnatEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeSnatEntry(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}

	}
}

func (s *VpcService) WaitForCommonBandwidthPackage(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCommonBandwidthPackage(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
	}
}

func (s *VpcService) WaitForCommonBandwidthPackageAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCommonBandwidthPackageAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
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
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *VpcService) WaitForRouteTableAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouteTableAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *VpcService) DescribeNetworkAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeNetworkAclAttributes"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"NetworkAclId": id,
	}
	request["Product"] = "Vpc"
	request["OrganizationId"] = s.client.Department
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidNetworkAcl.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC:NetworkAcl", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NetworkAclAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetworkAclAttribute", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeNetworkAclAttachment(id string, resource []vpc.Resource) (err error) {

	invoker := NewInvoker()
	return invoker.Run(func() error {
		object, err := s.DescribeNetworkAcl(id)
		if err != nil {
			return WrapError(err)
		}
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		if len(resources) < 1 {
			return WrapErrorf(Error(GetNotFoundMessage("Network Acl Attachment", id)), NotFoundMsg, ProviderERROR)
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
				return WrapErrorf(Error(GetNotFoundMessage("Network Acl Attachment", id)), NotFoundMsg, ProviderERROR)
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
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
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
			return WrapErrorf(err, WaitTimeoutMsg, networkAclId, GetFunc(1), timeout, fmt.Sprint(object["Status"]), string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForNetworkAclAttachment(id string, resource []vpc.Resource, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		err := s.DescribeNetworkAclAttachment(id, resource)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
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
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(object["Status"]), string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []vpc.TagResource, err error) {
	request := vpc.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
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
			if IsExpectedErrors(err, []string{Throttling}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabacloudStackSdkGoERROR)
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
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = string(resourceType)
			if strings.ToLower(s.client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			request.Headers = map[string]string{"RegionId": s.client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err := s.client.WithVpcClient(func(client *vpc.Client) (interface{}, error) {
					return client.UnTagResources(request)
				})
				if err != nil {
					if IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
		}

		if len(create) > 0 {
			request := vpc.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			if strings.ToLower(s.client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.Tag = &create
			request.ResourceType = string(resourceType)
			request.RegionId = s.client.RegionId
			request.Headers = map[string]string{"RegionId": s.client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err := s.client.WithVpcClient(func(client *vpc.Client) (interface{}, error) {
					return client.TagResources(request)
				})
				if err != nil {
					if IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
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
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"VSwitchId": id,
	}
	request["Product"] = "Vpc"
	request["OrganizationId"] = s.client.Department
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	addDebug(action, response, request)
	if v, ok := response["VSwitchId"].(string); ok && v != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("vswitch", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (s *VpcService) NetworkAclStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNetworkAcl(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DeleteAclResources(id string) (object map[string]interface{}, err error) {
	acl, err := s.DescribeNetworkAcl(id)
	if err != nil {
		return object, WrapError(err)
	}
	deleteResources := make([]map[string]interface{}, 0)
	res, err := jsonpath.Get("$.Resources.Resource", acl)
	if err != nil {
		return object, WrapError(err)
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
		"RegionId":     s.client.RegionId,
	}
	request["Product"] = "Vpc"
	request["OrganizationId"] = s.client.Department
	action := "UnassociateNetworkAcl"
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	request["ClientToken"] = buildClientToken("UnassociateNetworkAcl")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, 10*time.Minute, 5*time.Second, s.NetworkAclStateRefreshFunc(id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return response, WrapErrorf(err, IdMsg, id)
	}
	return object, nil
}

func (s *VpcService) DescribeVpcIpv6EgressRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "DescribeIpv6EgressOnlyRules"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"PageSize":      PageSizeLarge,
		"PageNumber":    1,
		"Ipv6GatewayId": parts[0],
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			request["Product"] = "Vpc"
			request["OrganizationId"] = s.client.Department
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Ipv6EgressOnlyRules.Ipv6EgressOnlyRule", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Ipv6EgressOnlyRules.Ipv6EgressOnlyRule", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpcIpv6EgressRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpv6EgressRule(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {

			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcIpv6Gateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeIpv6GatewayAttribute"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"Ipv6GatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request["Product"] = "Vpc"
		request["OrganizationId"] = s.client.Department
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if Ipv6GatewayId, ok := object["Ipv6GatewayId"]; !ok || fmt.Sprint(Ipv6GatewayId) == "" {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *VpcService) VpcIpv6GatewayStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpv6Gateway(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcIpv6InternetBandwidth(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeIpv6Addresses"
	request := map[string]interface{}{
		"RegionId":                s.client.RegionId,
		"Ipv6InternetBandwidthId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request["Product"] = "Vpc"
		request["OrganizationId"] = s.client.Department
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Ipv6Addresses.Ipv6Address", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Ipv6Addresses.Ipv6Address", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["Ipv6InternetBandwidth"].(map[string]interface{})["Ipv6InternetBandwidthId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) setInstanceSecondaryCidrBlocks(d *schema.ResourceData) error {
	if d.HasChange("secondary_cidr_blocks") {
		oraw, nraw := d.GetChange("secondary_cidr_blocks")
		removed := oraw.([]interface{})
		added := nraw.([]interface{})
		conn, err := s.client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		if len(removed) > 0 {
			action := "UnassociateVpcCidrBlock"
			request := map[string]interface{}{
				"RegionId": s.client.RegionId,
				"VpcId":    d.Id(),
			}
			for _, item := range removed {
				request["SecondaryCidrBlock"] = item
				request["Product"] = "Vpc"
				request["OrganizationId"] = s.client.Department
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
				}
				addDebug(action, response, request)
			}
		}

		if len(added) > 0 {
			action := "AssociateVpcCidrBlock"
			request := map[string]interface{}{
				"RegionId": s.client.RegionId,
				"VpcId":    d.Id(),
			}
			for _, item := range added {
				request["SecondaryCidrBlock"] = item
				request["Product"] = "Vpc"
				request["OrganizationId"] = s.client.Department
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
				}
				addDebug(action, response, request)
			}
		}
		//d.SetPartial("secondary_cidr_blocks")
	}
	return nil
}

func (s *VpcService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				request["Product"] = "Vpc"
				request["OrganizationId"] = s.client.Department
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				request["Product"] = "Vpc"
				request["OrganizationId"] = s.client.Department
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
			}
		}
		//d.SetPartial("tags")
	}
	return nil
}

func (s *VpcService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"ResourceType":   resourceType,
		"ResourceId.1":   id,
		"Product":        "Vpc",
		"OrganizationId": s.client.Department,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *VpcService) DescribeExpressConnectPhysicalConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribePhysicalConnections"
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"Product":        "Vpc",
		"OrganizationId": s.client.Department,
	}
	filterMapList := make([]map[string]interface{}, 0)
	filterMapList = append(filterMapList, map[string]interface{}{
		"Key":   "PhysicalConnectionId",
		"Value": []string{id},
	})
	request["Filter"] = filterMapList
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.PhysicalConnectionSet.PhysicalConnectionType", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PhysicalConnectionSet.PhysicalConnectionType", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PhysicalConnectionId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) ExpressConnectPhysicalConnectionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectPhysicalConnection(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeExpressConnectVirtualBorderRouter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVirtualBorderRouters"
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"PageNumber":     1,
		"PageSize":       50,
		"Product":        "Vpc",
		"OrganizationId": s.client.Department,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
		}
		v, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) ExpressConnectVirtualBorderRouterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectVirtualBorderRouter(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
