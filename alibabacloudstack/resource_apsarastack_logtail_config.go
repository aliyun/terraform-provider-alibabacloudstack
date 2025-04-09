package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLogtailConfig() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"input_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"file", "plugin"}, false),
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"input_detail": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yaml, _ := normalizeJsonString(v)
					return yaml
				},
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackLogtailConfigCreate,
		resourceAlibabacloudStackLogtailConfigRead, resourceAlibabacloudStackLogtailConfiglUpdate, resourceAlibabacloudStackLogtailConfigDelete)
	return resource
}

func resourceAlibabacloudStackLogtailConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var inputConfigInputDetail = make(map[string]interface{})
	data := d.Get("input_detail").(string)
	if jsonErr := json.Unmarshal([]byte(data), &inputConfigInputDetail); jsonErr != nil {
		return errmsgs.WrapError(jsonErr)
	}
	var requestInfo *sls.Client
	logconfig := &sls.LogConfig{
		Name:       d.Get("name").(string),
		InputType:  d.Get("input_type").(string),
		OutputType: d.Get("output_type").(string),
		OutputDetail: sls.OutputDetail{
			ProjectName:  d.Get("project").(string),
			LogStoreName: d.Get("logstore").(string),
		},
	}
	raw, err := client.WithSlsDataClient(func(slsClient *sls.Client) (interface{}, error) {
		requestInfo = slsClient
		sls.AddNecessaryInputConfigField(inputConfigInputDetail)
		covertInput, covertErr := assertInputDetailType(inputConfigInputDetail, logconfig)
		if covertErr != nil {
			return nil, covertErr
		}
		logconfig.InputDetail = covertInput
		return nil, slsClient.CreateConfig(d.Get("project").(string), logconfig)
	})
	if err != nil {
		errmsg := ""
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_logtail_config", "CreateConfig", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}
	if debugOn() {
		addDebug("CreateConfig", raw, requestInfo, map[string]interface{}{
			"project": d.Get("project").(string),
			"config":  logconfig,
		})
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("logstore").(string), COLON_SEPARATED, d.Get("name").(string)))
	return nil
}

func resourceAlibabacloudStackLogtailConfigRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	config, err := logService.DescribeLogtailConfig(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Because the server will return redundant parameters, we filter here
	inputDetail := d.Get("input_detail").(string)
	var oMap map[string]interface{}
	json.Unmarshal([]byte(inputDetail), &oMap)
	nMap := config.InputDetail.(map[string]interface{})
	if inputDetail != "" {
		for nk := range nMap {
			if _, ok := oMap[nk]; !ok {
				delete(nMap, nk)
			}
		}
	}
	nMapJson, err := json.Marshal(nMap)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("input_detail", string(nMapJson))
	d.Set("project", split[0])
	d.Set("name", config.Name)
	d.Set("logstore", split[1])
	d.Set("input_type", config.InputType)
	d.Set("output_type", config.OutputType)
	return nil
}

func resourceAlibabacloudStackLogtailConfiglUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	update := false
	if d.HasChange("input_detail") {
		update = true
	}
	if d.HasChange("input_type") {
		update = true
	}
	if update {
		logconfig := &sls.LogConfig{}
		inputConfigInputDetail := make(map[string]interface{})
		data := d.Get("input_detail").(string)
		conver_err := json.Unmarshal([]byte(data), &inputConfigInputDetail)
		if conver_err != nil {
			return errmsgs.WrapError(conver_err)
		}
		sls.AddNecessaryInputConfigField(inputConfigInputDetail)
		covertInput, covertErr := assertInputDetailType(inputConfigInputDetail, logconfig)
		if covertErr != nil {
			return errmsgs.WrapError(covertErr)
		}
		logconfig.InputDetail = covertInput

		client := meta.(*connectivity.AlibabacloudStackClient)
		var requestInfo *sls.Client
		params := &sls.LogConfig{
			Name:        parts[2],
			InputType:   d.Get("input_type").(string),
			OutputType:  d.Get("output_type").(string),
			InputDetail: logconfig.InputDetail,
			OutputDetail: sls.OutputDetail{
				ProjectName:  d.Get("project").(string),
				LogStoreName: d.Get("logstore").(string),
			},
		}
		raw, err := client.WithSlsDataClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.UpdateConfig(parts[0], params)
		})
		if err != nil {
			errmsg := ""
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "UpdateConfig", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		if debugOn() {
			addDebug("UpdateConfig", raw, requestInfo, map[string]interface{}{
				"project": parts[0],
				"config":  params,
			})
		}
	}
	return nil
}

func resourceAlibabacloudStackLogtailConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	var requestInfo *sls.Client
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlsDataClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteConfig(parts[0], parts[2])
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteConfig", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg))
		}
		if debugOn() {
			addDebug("DeleteConfig", raw, requestInfo, map[string]string{
				"project": parts[0],
				"config":  parts[2],
			})
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist", "ConfigNotExist"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteConfig", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	return errmsgs.WrapError(logService.WaitForLogtailConfig(d.Id(), Deleted, DefaultTimeout))
}

// This function is used to assert and convert the type to sls.LogConfig
func assertInputDetailType(inputConfigInputDetail map[string]interface{}, logconfig *sls.LogConfig) (sls.InputDetailInterface, error) {
	if inputConfigInputDetail["logType"] == "json_log" {
		JSONConfigInputDetail, ok := sls.ConvertToJSONConfigInputDetail(inputConfigInputDetail)
		if !ok {
			return nil, errmsgs.WrapError(errmsgs.Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = JSONConfigInputDetail
	}
	if inputConfigInputDetail["logType"] == "apsara_log" {
		ApsaraLogConfigInputDetail, ok := sls.ConvertToApsaraLogConfigInputDetail(inputConfigInputDetail)
		if !ok {
			return nil, errmsgs.WrapError(errmsgs.Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = ApsaraLogConfigInputDetail
	}

	if inputConfigInputDetail["logType"] == "common_reg_log" {
		RegexConfigInputDetail, ok := sls.ConvertToRegexConfigInputDetail(inputConfigInputDetail)
		if !ok {
			return nil, errmsgs.WrapError(errmsgs.Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = RegexConfigInputDetail
	}
	if inputConfigInputDetail["logType"] == "delimiter_log" {
		DelimiterConfigInputDetail, ok := sls.ConvertToDelimiterConfigInputDetail(inputConfigInputDetail)
		if !ok {
			return nil, errmsgs.WrapError(errmsgs.Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = DelimiterConfigInputDetail
	}
	if logconfig.InputType == "plugin" {
		PluginLogConfigInputDetail, ok := sls.ConvertToPluginLogConfigInputDetail(inputConfigInputDetail)
		if !ok {
			return nil, errmsgs.WrapError(errmsgs.Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = PluginLogConfigInputDetail
	}
	return logconfig.InputDetail, nil
}
