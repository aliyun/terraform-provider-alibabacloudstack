package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type CmsService struct {
	client *connectivity.AlibabacloudStackClient
}

type IspCities []map[string]string

func (s *CmsService) BuildCmsCommonRequest(region string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	return request
}

func (s *CmsService) BuildCmsAlarmRequest(id string) *requests.CommonRequest {

	request := s.BuildCmsCommonRequest(s.client.RegionId)
	request.QueryParams["Id"] = id

	return request
}

func (s *CmsService) DescribeCmsAlarm(id string) (alarm cms.Alarm, err error) {
	request := cms.CreateDescribeMetricRuleListRequest()
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alarm, WrapError(err)
	}
	request.RuleIds = parts[0]

	if len(parts) != 1 {
		request.RuleName = parts[1]
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "cms"}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response *cms.DescribeMetricRuleListResponse
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeMetricRuleList(request)
		})
		if err != nil && IsExpectedErrors(err, []string{Throttling}) {
			time.Sleep(10 * time.Second)
			return resource.RetryableError(err)
		}
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*cms.DescribeMetricRuleListResponse)
		return nil
	})
	if err != nil {
		return alarm, err
	}
	if len(response.Alarms.Alarm) < 1 {
		return alarm, GetNotFoundErrorFromString(GetNotFoundMessage("Alarm Rule", id))
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
			return GetTimeErrorFromString(GetTimeoutMessage("Alarm", strconv.FormatBool(enabled)))
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

func (s *CmsService) DescribeSiteMonitor(id, keyword string) (siteMonitor cms.SiteMonitor, err error) {
	listRequest := cms.CreateDescribeSiteMonitorListRequest()
	listRequest.Headers = map[string]string{"RegionId": s.client.RegionId}
	listRequest.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "cms", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	listRequest.Keyword = keyword
	listRequest.TaskId = id
	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeSiteMonitorList(listRequest)
	})
	if err != nil {
		return siteMonitor, err
	}
	list := raw.(*cms.DescribeSiteMonitorListResponse)
	if len(list.SiteMonitors.SiteMonitor) < 1 {
		return siteMonitor, GetNotFoundErrorFromString(GetNotFoundMessage("Site Monitor", id))

	}
	for _, v := range list.SiteMonitors.SiteMonitor {
		if v.TaskName == keyword || v.TaskId == id {
			return v, nil
		}
	}
	return siteMonitor, GetNotFoundErrorFromString(GetNotFoundMessage("Site Monitor", id))
}

func (s *CmsService) GetIspCities(id string) (ispCities IspCities, err error) {
	request := cms.CreateDescribeSiteMonitorAttributeRequest()
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "cms", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.TaskId = id

	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeSiteMonitorAttribute(request)
	})

	if err != nil {
		return nil, err
	}

	response := raw.(*cms.DescribeSiteMonitorAttributeResponse)
	ispCity := response.SiteMonitors.IspCities.IspCity

	var list []map[string]string
	for _, element := range ispCity {
		list = append(list, map[string]string{"city": element.City, "isp": element.Isp})
	}

	return list, nil
}

func (s *CmsService) DescribeCmsAlarmContact(id string) (object cms.Contact, err error) {
	request := cms.CreateDescribeContactListRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "cms", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.ContactName = id

	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeContactList(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ContactNotExists", "ResourceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContact", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cms.DescribeContactListResponse)
	if response.Code != "200" {
		err = Error("DescribeContactList failed for " + response.Message)
		return
	}

	if len(response.Contacts.Contact) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContact", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Contacts.Contact[0], nil
}

func (s *CmsService) DescribeCmsAlarmContactGroup(id string) (object cms.ContactGroup, err error) {
	request := cms.CreateDescribeContactGroupListRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "cms", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	for {

		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeContactGroupList(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ContactGroupNotExists", "ResourceNotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cms.DescribeContactGroupListResponse)
		if response.Code != "200" {
			err = Error("DescribeContactGroupList failed for " + response.Message)
			return object, err
		}

		if len(response.ContactGroupList.ContactGroup) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR, response.RequestId)
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
			return object, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR)
	return
}

func (s *CmsService) DescribeCmsMetricRuleTemplateList() (templates []cms.Template, err error) {
	request := cms.CreateDescribeMetricRuleTemplateListRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	PageNumber := 1
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "cms",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"IsDefault":       "false",
		"History":         "ture",
		"PageNumber":      fmt.Sprintf("%d", PageNumber),
		"PageSize":        "20",
	}
	for {
		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeMetricRuleTemplateList(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			return templates, WrapErrorf(err, DefaultErrorMsg, "DescribeMetricRuleTemplateList", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		response, _ := raw.(*cms.DescribeMetricRuleTemplateListResponse)
		if response.Code != 200 {
			err = Error("DescribeMetricRuleTemplateList failed for " + response.Message)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, "DescribeCmsMetricRuleTemplateList", AlibabacloudStackSdkGoERROR)
	}
	for _, template := range templates {
		if fmt.Sprintf("%d", template.TemplateId) == id {
			return template, nil
		}
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("CmsMetricRuleTemplate", id)), NotFoundMsg, AlibabacloudStackSdkGoERROR)
}

func (s *CmsService) DescribeMetricRuleTemplateAttribute(id string) (object *cms.DescribeMetricRuleTemplateAttributeResponse, err error) {
	request := cms.CreateDescribeMetricRuleTemplateAttributeRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "cms",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
	}
	request.TemplateId = id
	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeMetricRuleTemplateAttribute(request)
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"TemplateNotExists", "ResourceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DescribeMetricRuleTemplateAttribute", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return object, err
	}
	bresponse, _ := raw.(*cms.DescribeMetricRuleTemplateAttributeResponse)
	if bresponse.Code == 200 {
		return bresponse, nil
	} else {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DescribeMetricRuleTemplateAttribute", id)), NotFoundMsg, ProviderERROR)
	}
}
