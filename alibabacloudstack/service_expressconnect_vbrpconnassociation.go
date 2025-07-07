package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"strconv"
	"time"
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

)

type 
ExpressconnectService struct {
	client *connectivity.AlibabacloudStackClient
}


	
	
type VpcDescribevirtualborderroutersResponse struct {
    VirtualBorderRouterSet struct {
    VirtualBorderRouterType []struct {
    AssociatedPhysicalConnections struct {
    AssociatedPhysicalConnection []struct {CircuitCode string `json:"CircuitCode"`
VlanInterfaceId string `json:"VlanInterfaceId"`
LocalGatewayIp string `json:"LocalGatewayIp"`
PeerGatewayIp string `json:"PeerGatewayIp"`
PeeringSubnetMask string `json:"PeeringSubnetMask"`
PhysicalConnectionId string `json:"PhysicalConnectionId"`
PhysicalConnectionStatus string `json:"PhysicalConnectionStatus"`
PhysicalConnectionBusinessStatus string `json:"PhysicalConnectionBusinessStatus"`
PhysicalConnectionOwnerUid string `json:"PhysicalConnectionOwnerUid"`
VlanId string `json:"VlanId"`
LocalIpv6GatewayIp string `json:"LocalIpv6GatewayIp"`
PeerIpv6GatewayIp string `json:"PeerIpv6GatewayIp"`
PeeringIpv6SubnetMask string `json:"PeeringIpv6SubnetMask"`
Status string `json:"Status"`
EnableIpv6 bool `json:"EnableIpv6"`
} `json:"AssociatedPhysicalConnection"`
} `json:"AssociatedPhysicalConnections"`

    AssociatedCens struct {
    AssociatedCen []struct {CenId string `json:"CenId"`
CenOwnerId int `json:"CenOwnerId"`
CenStatus string `json:"CenStatus"`
} `json:"AssociatedCen"`
} `json:"AssociatedCens"`
VbrId string `json:"VbrId"`
CreationTime string `json:"CreationTime"`
ActivationTime string `json:"ActivationTime"`
TerminationTime string `json:"TerminationTime"`
RecoveryTime string `json:"RecoveryTime"`
Status string `json:"Status"`
VlanId int `json:"VlanId"`
CircuitCode string `json:"CircuitCode"`
RouteTableId string `json:"RouteTableId"`
VlanInterfaceId string `json:"VlanInterfaceId"`
LocalGatewayIp string `json:"LocalGatewayIp"`
PeerGatewayIp string `json:"PeerGatewayIp"`
PeeringSubnetMask string `json:"PeeringSubnetMask"`
PhysicalConnectionId string `json:"PhysicalConnectionId"`
PhysicalConnectionStatus string `json:"PhysicalConnectionStatus"`
PhysicalConnectionBusinessStatus string `json:"PhysicalConnectionBusinessStatus"`
PhysicalConnectionOwnerUid string `json:"PhysicalConnectionOwnerUid"`
AccessPointId string `json:"AccessPointId"`
Name string `json:"Name"`
Description string `json:"Description"`
PConnVbrExpireTime string `json:"PConnVbrExpireTime"`
EccId string `json:"EccId"`
Type string `json:"Type"`
MinTxInterval int `json:"MinTxInterval"`
MinRxInterval int `json:"MinRxInterval"`
DetectMultiplier int `json:"DetectMultiplier"`
LocalIpv6GatewayIp string `json:"LocalIpv6GatewayIp"`
PeerIpv6GatewayIp string `json:"PeerIpv6GatewayIp"`
PeeringIpv6SubnetMask string `json:"PeeringIpv6SubnetMask"`
EnableIpv6 bool `json:"EnableIpv6"`
} `json:"VirtualBorderRouterType"`
} `json:"VirtualBorderRouterSet"`
RequestId string `json:"RequestId"`
PageNumber int `json:"PageNumber"`
PageSize int `json:"PageSize"`
TotalCount int `json:"TotalCount"`
}

func (s *
ExpressconnectService)DoVpcDescribevirtualborderroutersRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*VpcDescribevirtualborderroutersResponse, error) {
	// api: Vpc - 2016-04-28 - DescribeVirtualBorderRouters
    request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeVirtualBorderRouters", "")
	VpcDescribevirtualborderroutersResponseObj := &VpcDescribevirtualborderroutersResponse{}

	
	
		
		
			//调用request_params_handler
    
        
        
    
        
        
                request.QueryParams["Filter[*].Key"] = d.Get("filter").(string)
        
    
        
    
        
        
    
        
        
            if v, ok := d.GetOk("page_number"); ok{
                request.QueryParams["PageNumber"] = strconv.Itoa(v.(int))
            }
        
    
        
    
        
        
    
        
        
            if v, ok := d.GetOk("page_size"); ok{
                request.QueryParams["PageSize"] = strconv.Itoa(v.(int))
            }
        
    
        
    
        
        
    
        
        
                request.QueryParams["RegionId"] = d.Get("region_id").(string)
        
    
        
    
        
        
    
        
        
            if v, ok := d.GetOk("vbr_id"); ok{
                request.QueryParams["Filter[*].Value"] = v.(string)
            }
        
    
        
    
		
	

	

		
        bresponse, err := s.client.ProcessCommonRequest(request)
        if err != nil {
            if bresponse == nil {
                return nil,errmsgs.WrapErrorf(err, "Process Common Request Failed")
            }
            errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
            return nil,errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,  "", "DescribeVirtualBorderRouters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
        }

		
			
			err = json.Unmarshal(bresponse.GetHttpContentBytes(),&VpcDescribevirtualborderroutersResponseObj)
			
			if err != nil {
				return nil,errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeVirtualBorderRouters", errmsgs.AlibabacloudStackSdkGoERROR)
			}
		

        
	return VpcDescribevirtualborderroutersResponseObj, nil
}

