package alibabacloudstack

import (
	"encoding/json"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type BucketVpcService struct {
	client *connectivity.AlibabacloudStackClient
}

type VpcListResult struct {
	Api             string        `json:"api"`
	AsapiRequestId  string        `json:"asapiRequestId"`
	AsapiSuccess    bool          `json:"asapiSuccess"`
	HttpOk          bool          `json:"httpOk"`
	Success         bool          `json:"success"`
	Code            int64         `json:"code"`
	Domain          string        `json:"domain"`
	Message         string        `json:"message"`
	ServerRole      string        `json:"serverRole"`
	EagleEyeTraceId string        `json:"eagleEyeTraceId"`
	VpcList         []interface{} `json:"data"`
	PageModel       interface{}   `json:"pageModel"`
}

func (s *BucketVpcService) BucketVpcList(bucketName string) (vpclist *VpcListResult, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListBucketVpc", "/ascm/manage/saleconf/ossIsolationVpc/select")
	mergeMaps(request.QueryParams, map[string]string{
		"BucketName":       bucketName,
		"SignatureVersion": "2.1",
		"OpenApiAction":    "ListBucketVpc",
	})
	bresponse, err := s.client.ProcessCommonRequest(request)
	log.Printf("Response of ListBucketVpc: %s", bresponse)
	if err != nil {
		if bresponse == nil {
			return vpclist, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		if ossNotFoundError(err) {
			return vpclist, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return vpclist, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "ListBucketVpc", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}
	log.Printf("Bresponse ListBucketVpc after error")
	addDebug("ListBucketVpc", bresponse, nil, request)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &vpclist)
	if err != nil {
		return vpclist, errmsgs.WrapError(err)
	}
	if !vpclist.Success {
		return vpclist, errmsgs.WrapError(err)
	}
	return vpclist, nil
}

func (s *BucketVpcService) BindBucket(vpcId string, vpcName string, vLan string, bucket string) error {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "BindBucketPolicy", "/ascm/manage/saleconf/ossIsolationVpc/bind")
	mergeMaps(request.QueryParams, map[string]string{
		"bucketName":       bucket,
		"vpcName":          vpcName,
		"vLan":             vLan,
		"vpcId":            vpcId,
		"SignatureVersion": "2.1",
	})
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		if ossNotFoundError(err) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucket, "BindBucketPolicy", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}
	log.Printf("Bresponse BindBucketPolicy after error")
	addDebug("CreateBucketInfo", bresponse, nil, request)
	log.Printf("Bresponse BindBucketPolicy check")
	log.Printf("Bresponse BindBucketPolicy %s", bresponse)

	return nil
}

func (s *BucketVpcService) UnBindBucket(vpcId string, bucket string) error {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "UnBindBucketPolicy", "/ascm/manage/saleconf/ossIsolationVpc/unBind")
	mergeMaps(request.QueryParams, map[string]string{
		"bucketName":       bucket,
		"vpcId":            vpcId,
		"SignatureVersion": "2.1",
	})
	bresponse, err := s.client.ProcessCommonRequest(request)
	log.Printf("Bresponse UnBindBucketPolicy before error")
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		if ossNotFoundError(err) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucket, "UnBindBucketPolicy", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}
	log.Printf("Bresponse UnBindBucketPolicy after error")
	addDebug("CreateBucketInfo", bresponse, nil, request)
	log.Printf("Bresponse UnBindBucketPolicy check")
	log.Printf("Bresponse UnBindBucketPolicy %s", bresponse)

	return nil
}
