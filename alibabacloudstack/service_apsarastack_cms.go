package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type CmsService struct {
	client *connectivity.AlibabacloudStackClient
}

type IspCities []map[string]string

func (s *CmsService) DescribeCmsAlarm(id string) (alarm cms.AlarmInDescribeMetricRuleList, err error) {
	request := cms.CreateDescribeMetricRuleListRequest()
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alarm, errmsgs.WrapError(err)
	}
	request.RuleIds = parts[0]

	if len(parts) != 1 {
		request.RuleName = parts[1]
	}
	s.client.InitRpcRequest(*request.RpcRequest)
	request.QueryParams["Product"] = "cms"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response *cms.DescribeMetricRuleListResponse
	var raw interface{}
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeMetricRuleList(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	response, _ = raw.(*cms.DescribeMetricRuleListResponse)
	if err != nil {
		errmsg := ""
		if response != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return alarm, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeCmsAlarm", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	if len(response.Alarms.Alarm) < 1 {
		return alarm, errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("Alarm Rule", id))
	}
	return response.Alarms.Alarm[0], nil
}

func (s *CmsService) WaitForCmsAlarm(id string, enabled bool, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		alarm, err := s.DescribeCmsAlarm(id)
		if err != nil {
			return err
		}

		if alarm.EnableState == enabled {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return errmsgs.GetTimeErrorFromString(errmsgs.GetTimeoutMessage("Alarm", strconv.FormatBool(enabled)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *CmsService) BuildJsonWebhook(webhook string) string {
	if webhook != "" {
		return fmt.Sprintf("{\"method\":\"post\",\"url\":\"%s\"}", webhook)
	}
	return ""
}

func (s *CmsService) ExtractWebhookFromJson(webhookJson string) (string, error) {
	byt := []byte(webhookJson)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		return "", err
	}
	return dat["url"].(string), nil
}

func (s *CmsService) DoCmsDescribesitemonitorattributeRequest(id, keyword string) (siteMonitor cms.SiteMonitor, err error) {
	return s.DescribeSiteMonitor(id, keyword)
}

func (s *CmsService) DescribeSiteMonitor(id, keyword string) (siteMonitor cms.SiteMonitor, err error) {
	listRequest := cms.CreateDescribeSiteMonitorListRequest()
	s.client.InitRpcRequest(*listRequest.RpcRequest)
	listRequest.QueryParams["Product"] = "cms"
	listRequest.Keyword = keyword
	listRequest.TaskId = id
	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeSiteMonitorList(listRequest)
	})
	list ,ok := raw.(*cms.DescribeSiteMonitorListResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(list.BaseResponse)
		}
		return siteMonitor, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeSiteMonitor", listRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	
	if len(list.SiteMonitors.SiteMonitor) < 1 {
		return siteMonitor, errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("Site Monitor", id))
	}
	for _, v := range list.SiteMonitors.SiteMonitor {
		if v.TaskName == keyword || v.TaskId == id {
			return v, nil
		}
	}
	return siteMonitor, errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("Site Monitor", id))
}

func (s *CmsService) GetIspCities(id string) (ispCities IspCities, err error) {
	request := cms.CreateDescribeSiteMonitorAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.QueryParams["Product"] = "cms"
	request.TaskId = id

	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeSiteMonitorAttribute(request)
	})

	response, ok := raw.(*cms.DescribeSiteMonitorAttributeResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "GetIspCities", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	ispCity := response.SiteMonitors.IspCities.IspCity

	var list []map[string]string
	for _, element := range ispCity {
		list = append(list, map[string]string{"city": element.City, "isp": element.Isp})
	}

	return list, nil
}

func (s *CmsService) DoCmsDescribecontactlistRequest(id string) (object cms.Contact, err error) {
	return s.DescribeCmsAlarmContact(id)
}

func (s *CmsService) DescribeCmsAlarmContact(id string) (object cms.Contact, err error) {
	request := cms.CreateDescribeContactListRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.QueryParams["Product"] = "cms"
	request.ContactName = id

	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeContactList(request)
	})
	response, ok := raw.(*cms.DescribeContactListResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ContactNotExists", "errmsgs.ResourceNotfound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CmsAlarmContact", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	
	if response.Code != "200" {
		err = errmsgs.Error("DescribeContactList failed for " + response.Message)
		return
	}

	if len(response.Contacts.Contact) < 1 {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CmsAlarmContact", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
		return
	}
	return response.Contacts.Contact[0], nil
}

func (s *CmsService) DoCmsDescribecontactgrouplistRequest(id string) (object cms.ContactGroup, err error) {
	return s.DescribeCmsAlarmContactGroup(id)
}

func (s *CmsService) DescribeCmsAlarmContactGroup(id string) (object cms.ContactGroup, err error) {
	request := cms.CreateDescribeContactGroupListRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.QueryParams["Product"] = "cms"
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	for {

		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeContactGroupList(request)
		})
		response, ok := raw.(*cms.DescribeContactGroupListResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"ContactGroupNotExists", "errmsgs.ResourceNotfound"}) {
				err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CmsAlarmContactGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
				return object, err
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		
		if response.Code != "200" {
			err = errmsgs.Error("DescribeContactGroupList failed for " + response.Message)
			return object, err
		}

		if len(response.ContactGroupList.ContactGroup) < 1 {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CmsAlarmContactGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
			return object, err
		}
		for _, object := range response.ContactGroupList.ContactGroup {
			if object.Name == id {
				return object, nil
			}
		}
		if len(response.ContactGroupList.ContactGroup) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return object, errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CmsAlarmContactGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	return
}

func (s *CmsService) DescribeCmsMetricRuleTemplateList() (templates []cms.Template, err error) {
	request := cms.CreateDescribeMetricRuleTemplateListRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.QueryParams["IsDefault"] = "false"
	request.QueryParams["History"] = "true"
	PageNumber := 1
	request.QueryParams["PageNumber"] = fmt.Sprintf("%d", PageNumber)
	request.QueryParams["PageSize"] = "20"
	for {
		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeMetricRuleTemplateList(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, ok := raw.(*cms.DescribeMetricRuleTemplateListResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return templates, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeMetricRuleTemplateList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		
		if response.Code != 200 {
			err = errmsgs.Error("DescribeMetricRuleTemplateList failed for " + response.Message)
			return templates, err
		} else {
			templates = append(templates, response.Templates.Template...)
		}
		if len(templates) < int(response.Total) {
			PageNumber++
			request.QueryParams["PageNumber"] = fmt.Sprintf("%d", PageNumber)
		} else {
			break
		}
	}
	return templates, nil
}

func (s *CmsService) DescribeCmsMetricRuleTemplateDetail(id string) (object cms.Template, err error) {
	templates, err := s.DescribeCmsMetricRuleTemplateList()
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "DescribeCmsMetricRuleTemplateList", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, template := range templates {
		if fmt.Sprintf("%d", template.TemplateId) == id {
			return template, nil
		}
	}
	return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CmsMetricRuleTemplate", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}

func (s *CmsService) DescribeMetricRuleTemplateAttribute(id string) (object *cms.DescribeMetricRuleTemplateAttributeResponse, err error) {
	request := cms.CreateDescribeMetricRuleTemplateAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.TemplateId = id
	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeMetricRuleTemplateAttribute(request)
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	bresponse, ok := raw.(*cms.DescribeMetricRuleTemplateAttributeResponse)
	if err != nil || bresponse.Code != 200{
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"TemplateNotExists", "errmsgs.ResourceNotfound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DescribeMetricRuleTemplateAttribute", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return object, err
	}
	
		return bresponse, nil
}
