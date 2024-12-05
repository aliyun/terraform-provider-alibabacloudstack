package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

const ModulesSizeLimit = 50

type BssopenapiService struct {
	client *connectivity.AlibabacloudStackClient
}

func (b *BssopenapiService) GetInstanceTypePrice(productCode, productType string, modules interface{}) ([]float64, error) {
	var detailList []bssopenapi.ModuleDetail
	switch modules := modules.(type) {
	case []bssopenapi.GetPayAsYouGoPriceModuleList:
		request := bssopenapi.CreateGetPayAsYouGoPriceRequest()
		b.client.InitRpcRequest(*request.RpcRequest)
		request.ProductCode = productCode
		request.ProductType = productType
		for {
			if len(modules) < ModulesSizeLimit {
				tmp := modules
				request.ModuleList = &tmp
			} else {
				tmp := modules[:ModulesSizeLimit]
				modules = modules[ModulesSizeLimit:]
				request.ModuleList = &tmp
			}
			data, err := b.getPayAsYouGoData(request)

			if err != nil {
				return nil, errmsgs.WrapError(err)
			}

			detailList = append(detailList, data.ModuleDetails.ModuleDetail...)

			if len(*request.ModuleList) < ModulesSizeLimit {
				break
			}
		}

	case []bssopenapi.GetSubscriptionPriceModuleList:
		request := bssopenapi.CreateGetSubscriptionPriceRequest()
		b.client.InitRpcRequest(*request.RpcRequest)
		request.ProductCode = productCode
		request.ProductType = productType
		request.ServicePeriodQuantity = requests.NewInteger(1)
		request.ServicePeriodUnit = "Month"
		request.Quantity = requests.NewInteger(1)

		for {
			if len(modules) < ModulesSizeLimit {
				tmp := modules
				request.ModuleList = &tmp
			} else {
				tmp := modules[:ModulesSizeLimit]
				modules = modules[ModulesSizeLimit:]
				request.ModuleList = &tmp
			}
			data, err := b.getSubscriptionData(request)
			if err != nil {
				return nil, errmsgs.WrapError(err)
			}

			detailList = append(detailList, data.ModuleDetails.ModuleDetail...)

			if len(*request.ModuleList) < ModulesSizeLimit {
				break
			}
		}
	}

	var priceList []float64
	for _, module := range detailList {
		priceList = append(priceList, module.OriginalCost)
	}
	return priceList, nil
}

func (b *BssopenapiService) getSubscriptionData(request *bssopenapi.GetSubscriptionPriceRequest) (*bssopenapi.DataInGetSubscriptionPrice, error) {
	request.OrderType = "NewOrder"
	request.SubscriptionType = "Subscription"
	request.RegionId = b.client.RegionId
	raw, err := b.client.WithBssopenapiClient(func(client *bssopenapi.Client) (interface{}, error) {
		return client.GetSubscriptionPrice(request)
	})

	response, ok := raw.(*bssopenapi.GetSubscriptionPriceResponse)
	if err != nil || !response.Success {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bssopenapi", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)


	if len(response.Data.ModuleDetails.ModuleDetail) == 0 {
		return nil, errmsgs.WrapError(errmsgs.Error("Api:GetSubscriptionPrice  Modules:%v  RequestId:%s  the moduleDetails length is 0!",
			request.ModuleList, response.RequestId))
	}
	return &response.Data, nil
}

func (b *BssopenapiService) getPayAsYouGoData(request *bssopenapi.GetPayAsYouGoPriceRequest) (*bssopenapi.DataInGetPayAsYouGoPrice, error) {
	request.SubscriptionType = "PayAsYouGo"
	b.client.InitRpcRequest(*request.RpcRequest)
	request.RegionId = b.client.RegionId
	raw, err := b.client.WithBssopenapiClient(func(client *bssopenapi.Client) (interface{}, error) {
		return client.GetPayAsYouGoPrice(request)
	})

	response, ok := raw.(*bssopenapi.GetPayAsYouGoPriceResponse)
	if err != nil || !response.Success{
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bssopenapi", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(response.Data.ModuleDetails.ModuleDetail) == 0 {
		return nil, errmsgs.WrapError(errmsgs.Error("Api:GetPayAsYouGoPrice  Modules:%v  RequestId:%s  the moduleDetails length is 0!",
			request.ModuleList, response.RequestId))
	}
	return &response.Data, nil
}
