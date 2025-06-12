package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EbsService struct {
	client *connectivity.AlibabacloudStackClient
}

type EbsReplicaGroup struct {
	PairIds             []string `json:"PairIds"`
	ReplicaGroupId      string   `json:"ReplicaGroupId"`
	SourceRegionId      string   `json:"SourceRegionId"`
	SourceZoneId        string   `json:"SourceZoneId"`
	DestinationRegionId string   `json:"DestinationRegionId"`
	DestinationZoneId   string   `json:"DestinationZoneId"`
	GroupName           string   `json:"GroupName"`
	Description         string   `json:"Description"`
	Status              string   `json:"Status"`
	RPO                 int      `json:"RPO"`
	LastRecoverPoint    int      `json:"LastRecoverPoint"`
	Site                string   `json:"Site"`
	PairNumber          int      `json:"PairNumber"`
}

type EbsDescribediskreplicagroupsResponse struct {
	ReplicaGroups []EbsReplicaGroup `json:"ReplicaGroups"`
	RequestId     string            `json:"RequestId"`
	NextToken     string            `json:"NextToken"`
}

func (s *EbsService) DoEbsDescribediskreplicagroupsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*EbsDescribediskreplicagroupsResponse, error) {
	// api: ebs - 2021-07-30 - DescribeDiskReplicaGroups
	request := client.NewCommonRequest("POST", "ebs", "2021-07-30", "DescribeDiskReplicaGroups", "")
	EbsDescribediskreplicagroupsResponse := &EbsDescribediskreplicagroupsResponse{}
	request.QueryParams["GroupIds"] = d.Id()
	//调用request_params_handler

	if v, ok := d.GetOk("max_results"); ok {
		request.QueryParams["MaxResults"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("next_token"); ok {
		request.QueryParams["NextToken"] = v.(string)
	}

	if v, ok := d.GetOk("region_id"); ok {
		request.QueryParams["RegionId"] = v.(string)
	} else {
		return nil, fmt.Errorf("RegionId is required")
	}

	if v, ok := d.GetOk("site"); ok {
		request.QueryParams["Site"] = v.(string)
	}
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDiskReplicaGroups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &EbsDescribediskreplicagroupsResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDiskReplicaGroups", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if len(EbsDescribediskreplicagroupsResponse.ReplicaGroups) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ebs_diskreplicagroup", d.Id())), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return EbsDescribediskreplicagroupsResponse, nil
}

type EbsDiskReplicaPair struct {
	ReplicaPairId     string `json:"ReplicaPairId"`
	SourceRegion      string `json:"SourceRegion"`
	SourceZoneId      string `json:"SourceZoneId"`
	SourceDiskId      string `json:"SourceDiskId"`
	DestinationRegion string `json:"DestinationRegion"`
	DestinationZoneId string `json:"DestinationZoneId"`
	DestinationDiskId string `json:"DestinationDiskId"`
	PairName          string `json:"PairName"`
	Description       string `json:"Description"`
	Status            string `json:"Status"`
	RPO               int    `json:"RPO"`
	Bandwidth         int    `json:"Bandwidth"`
	StatusMessage     string `json:"StatusMessage"`
	LastRecoverPoint  int    `json:"LastRecoverPoint"`
	ReplicaGroupId    string `json:"ReplicaGroupId"`
	CreateTime        int64    `json:"CreateTime"`
	ReplicaGroupName  string `json:"ReplicaGroupName"`
	Site              string `json:"Site"`
	PrimaryRegion     string `json:"PrimaryRegion"`
	StandbyRegion     string `json:"StandbyRegion"`
	PrimaryZone       string `json:"PrimaryZone"`
	StandbyZone       string `json:"StandbyZone"`
}

type EbsDescribediskreplicapairsResponse struct {
	ReplicaPairs []EbsDiskReplicaPair `json:"ReplicaPairs"`
	RequestId    string               `json:"RequestId"`
	NextToken    string               `json:"NextToken"`
}

func (s *EbsService) DoEbsDescribediskreplicapairsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*EbsDescribediskreplicapairsResponse, error) {
	// api: ebs - 2021-07-30 - DescribeDiskReplicaPairs
	request := client.NewCommonRequest("POST", "ebs", "2021-07-30", "DescribeDiskReplicaPairs", "")
	EbsDescribediskreplicapairsResponse := &EbsDescribediskreplicapairsResponse{}

	//调用request_params_handler

	if v, ok := d.GetOk("max_results"); ok {
		request.QueryParams["MaxResults"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("next_token"); ok {
		request.QueryParams["NextToken"] = v.(string)
	}

	if v, ok := d.GetOk("region_id"); ok {
		request.QueryParams["RegionId"] = v.(string)
	} else {
		return nil, fmt.Errorf("RegionId is required")
	}

	if v, ok := d.GetOk("replica_group_id"); ok {
		request.QueryParams["ReplicaGroupId"] = v.(string)
	}

	if v, ok := d.GetOk("replica_pair_id"); ok {
		request.QueryParams["PairIds"] = v.(string)
	}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDiskReplicaPairs", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &EbsDescribediskreplicapairsResponse)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDiskReplicaPairs", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return EbsDescribediskreplicapairsResponse, nil
}

func (s *EbsService) Describediskreplicapairs(id string) (*EbsDiskReplicaPair, error) {
	// api: ebs - 2021-07-30 - DescribeDiskReplicaPairs
	request := s.client.NewCommonRequest("POST", "ebs", "2021-07-30", "DescribeDiskReplicaPairs", "")
	EbsDescribediskreplicapairsResponse := &EbsDescribediskreplicapairsResponse{}
	request.QueryParams["PairIds"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDiskReplicaPairs", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &EbsDescribediskreplicapairsResponse)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDiskReplicaPairs", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(EbsDescribediskreplicapairsResponse.ReplicaPairs) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ebs_diskreplicapair", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	replica_pair := &EbsDiskReplicaPair{}
	for _, v := range EbsDescribediskreplicapairsResponse.ReplicaPairs {
		if v.ReplicaPairId == id {
			replica_pair = &v
			break
		}
	}
	if replica_pair == nil {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ebs_diskreplicapair", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return replica_pair, nil
}

func (s *EbsService) DescribeEbsDiskreplicagroups(id string) (*EbsReplicaGroup, error) {
	// api: ebs - 2021-07-30 - DescribeDiskReplicaGroups
	request := s.client.NewCommonRequest("POST", "ebs", "2021-07-30", "DescribeDiskReplicaGroups", "")
	data := &EbsReplicaGroup{}
	EbsDescribediskreplicagroupsResponse := &EbsDescribediskreplicagroupsResponse{}

	//调用request_params_handler

	request.QueryParams["GroupIds"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDiskReplicaGroups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &EbsDescribediskreplicagroupsResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDiskReplicaGroups", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if len(EbsDescribediskreplicagroupsResponse.ReplicaGroups) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ebs_diskreplicagroup", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, v := range EbsDescribediskreplicagroupsResponse.ReplicaGroups {
		if v.ReplicaGroupId == id {
			data = &v
			break
		}
	}
	if data == nil {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ebs_diskreplicagroup", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return data, nil
}

func (s *EbsService) EbsDiskreplicagroupStatusRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEbsDiskreplicagroups(id)
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

func (s *EbsService) EbsDiskreplicapairStatusRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.Describediskreplicapairs(id)
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
