package alibabacloudstack

import "github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

const (
	Average          = "Average"
	Minimum          = "Minimum"
	Maximum          = "Maximum"
	Value            = "Value"
	ErrorCodeMaximum = "ErrorCodeMaximum"
)

const (
	MoreThan        = ">"
	MoreThanOrEqual = ">="
	LessThan        = "<"
	LessThanOrEqual = "<="
	Equal           = "=="
	NotEqual        = "!="
)

const (
	SiteMonitorHTTP = "HTTP"
	SiteMonitorPing = "Ping"
	SiteMonitorTCP  = "TCP"
	SiteMonitorUDP  = "UDP"
	SiteMonitorDNS  = "DNS"
	SiteMonitorSMTP = "SMTP"
	SiteMonitorPOP3 = "POP3"
	SiteMonitorFTP  = "FTP"
)

type CmsContact struct {
	Code string `json:"Code"`
	Cost int    `json:"Cost"`
	Data []struct {
		Cid  string `json:"Cid"`
		Name string `json:"Name"`
	} `json:"Data"`
	Message  string `json:"Message"`
	Redirect bool   `json:"Redirect"`
	Success  bool   `json:"Success"`
}

type MetaList struct {
	TotalCount int    `json:"TotalCount"`
	RequestID  string `json:"RequestId"`
	Resources  struct {
		Resource []struct {
			MetricName  string `json:"MetricName"`
			Periods     string `json:"Periods"`
			Description string `json:"Description"`
			Dimensions  string `json:"Dimensions"`
			Labels      string `json:"Labels"`
			Unit        string `json:"Unit"`
			Statistics  string `json:"Statistics"`
			Namespace   string `json:"Namespace"`
		} `json:"Resource"`
	} `json:"Resources"`
	Code    int  `json:"Code"`
	Success bool `json:"Success"`
}

type AlarmsData struct {
	RequestID string `json:"RequestId"`
	Total     int    `json:"Total"`
	Alarms    struct {
		Alarm []struct {
			GroupName           string `json:"GroupName"`
			NoEffectiveInterval string `json:"NoEffectiveInterval"`
			SilenceTime         int    `json:"SilenceTime"`
			ContactGroups       string `json:"ContactGroups"`
			MailSubject         string `json:"MailSubject"`
			SourceType          string `json:"SourceType"`
			RuleID              string `json:"RuleId"`
			Period              int    `json:"Period"`
			Dimensions          string `json:"Dimensions"`
			EffectiveInterval   string `json:"EffectiveInterval"`
			AlertState          string `json:"AlertState"`
			Namespace           string `json:"Namespace"`
			GroupID             string `json:"GroupId"`
			MetricName          string `json:"MetricName"`
			EnableState         bool   `json:"EnableState"`
			Escalations         struct {
				Critical struct {
					ComparisonOperator string `json:"ComparisonOperator"`
					Times              int    `json:"Times"`
					Statistics         string `json:"Statistics"`
					Threshold          string `json:"Threshold"`
				} `json:"Critical"`
				Info struct {
					ComparisonOperator string `json:"ComparisonOperator"`
					Times              int    `json:"Times"`
					Statistics         string `json:"Statistics"`
					Threshold          string `json:"Threshold"`
				} `json:"Info"`
				Warn struct {
					ComparisonOperator string `json:"ComparisonOperator"`
					Times              int    `json:"Times"`
					Statistics         string `json:"Statistics"`
					Threshold          string `json:"Threshold"`
				} `json:"Warn"`
			} `json:"Escalations"`
			Webhook   string `json:"Webhook"`
			Resources string `json:"Resources"`
			RuleName  string `json:"RuleName"`
		} `json:"Alarm"`
	} `json:"Alarms"`
	Code    string `json:"Code"`
	Success bool   `json:"Success"`
}

type DescribeMetricRuleTemplateAttributeResponse struct {
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	ResponseVersion string `json:"responseVersion"`
	RequestID       string `json:"RequestId"`
	Success         bool   `json:"success"`
	RequestId       string `json:"requestId"`
	Resource        struct {
		AlertTemplates struct {
			AlertTemplate []struct {
				MetricName  string `json:"MetricName"`
				Category    string `json:"Category"`
				Escalations struct {
					Critical Escalation `json:"Critical"`
					Info     Escalation `json:"Info"`
					Warn     Escalation `json:"Warn"`
				} `json:"Escalations"`
				RuleName  string      `json:"RuleName"`
				Namespace string      `json:"Namespace"`
				Webhook   string      `json:"Webhook"`
				Selector  interface{} `json:"Selector"` // Empty object {}, could be map[string]interface{}
			} `json:"AlertTemplate"`
		} `json:"AlertTemplates"`
		Description string `json:"Description"`
		RestVersion int    `json:"RestVersion"`
		TemplateID  int    `json:"TemplateId"`
		Name        string `json:"Name"`
	} `json:"Resource"`
	Code int `json:"Code"`
}

type Escalation struct {
	ComparisonOperator string `json:"ComparisonOperator,omitempty"`
	Times              int    `json:"Times,omitempty"`
	Statistics         string `json:"Statistics,omitempty"`
	Threshold          string `json:"Threshold,omitempty"`
}

type DescribeMetricRuleTemplateListResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Total     int64  `json:"Total" xml:"Total"`
	Success   bool   `json:"Success" xml:"Success"`
	Templates struct {
		Template []Template `json:"Template" xml:"Template"`
	} `json:"Templates" xml:"Templates"`
}

type Template struct {
	Description  string `json:"Description" xml:"Description"`
	GmtCreate    int64  `json:"GmtCreate" xml:"GmtCreate"`
	Name         string `json:"Name" xml:"Name"`
	RestVersion  int64  `json:"RestVersion" xml:"RestVersion"`
	GmtModified  int64  `json:"GmtModified" xml:"GmtModified"`
	TemplateId   int64  `json:"TemplateId" xml:"TemplateId"`
	ApplyRecords struct {
		ApplyItem []struct {
			GroupName string `json:"GroupName" xml:"GroupName"`
			ApplyTime int64  `json:"ApplyTime" xml:"ApplyTime"`
			GroupId   int    `json:"GroupId" xml:"GroupId"`
		} `json:"ApplyItem" xml:"ApplyItem"`
	} `json:"ApplyRecords" xml:"ApplyRecords"`
	ApplyHistories struct {
		ApplyHistory []struct {
			GroupId   int64  `json:"GroupId" xml:"GroupId"`
			GroupName string `json:"GroupName" xml:"GroupName"`
			ApplyTime int64  `json:"ApplyTime" xml:"ApplyTime"`
		} `json:"ApplyHistory" xml:"ApplyHistory"`
	} `json:"ApplyHistories" xml:"ApplyHistories"`
}
