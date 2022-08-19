package apsarastack

const (
	Average          = "Average"
	Minimum          = "Minimum"
	Maximum          = "Maximum"
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
