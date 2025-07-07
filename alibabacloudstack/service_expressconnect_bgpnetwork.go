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


	
	
type VpcDescribebgpnetworksResponse struct {
    BgpNetworks struct {
    BgpNetwork []struct {VpcId string `json:"VpcId"`
DstCidrBlock string `json:"DstCidrBlock"`
RouterId string `json:"RouterId"`
Status string `json:"Status"`
} `json:"BgpNetwork"`
} `json:"BgpNetworks"`
RequestId string `json:"RequestId"`
TotalCount int `json:"TotalCount"`
PageNumber int `json:"PageNumber"`
PageSize int `json:"PageSize"`
}

func (s *
ExpressconnectService)DoVpcDescribebgpnetworksRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*VpcDescribebgpnetworksResponse, error) {
	// api: Vpc - 2016-04-28 - DescribeBgpNetworks
    request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeBgpNetworks", "")
	VpcDescribebgpnetworksResponseObj := &VpcDescribebgpnetworksResponse{}

	
	
		
		
			//调用request_params_handler
    
        
        
    
        
        
            if v, ok := d.GetOk("page_number"); ok{
                request.QueryParams["PageNumber"] = strconv.Itoa(v.(int))
            }
        
    
        
    
        
        
    
        
        
            if v, ok := d.GetOk("page_size"); ok{
                request.QueryParams["PageSize"] = strconv.Itoa(v.(int))
            }
        
    
        
    
        
        
    
        
        
                request.QueryParams["RegionId"] = d.Get("region_id").(string)
        
    
        
    
        
        
    
        
        
            if v, ok := d.GetOk("router_id"); ok{
                request.QueryParams["RouterId"] = v.(string)
            }
        
    
        
    
		
	

	

		
        bresponse, err := s.client.ProcessCommonRequest(request)
        if err != nil {
            if bresponse == nil {
                return nil,errmsgs.WrapErrorf(err, "Process Common Request Failed")
            }
            errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
            return nil,errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,  "", "DescribeBgpNetworks", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
        }

		
			
			err = json.Unmarshal(bresponse.GetHttpContentBytes(),&VpcDescribebgpnetworksResponseObj)
			
			if err != nil {
				return nil,errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBgpNetworks", errmsgs.AlibabacloudStackSdkGoERROR)
			}
		

        
	return VpcDescribebgpnetworksResponseObj, nil
}

