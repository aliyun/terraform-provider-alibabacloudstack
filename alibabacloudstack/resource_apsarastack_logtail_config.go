package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLogtailConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackLogtailConfigCreate,
		Read:   resourceAlibabacloudStackLogtailConfigRead,
		Update: resourceAlibabacloudStackLogtailConfiglUpdate,
		Delete: resourceAlibabacloudStackLogtailConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
}

func resourceAlibabacloudStackLogtailConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var inputConfigInputDetail = make(map[string]interface{})
	data := d.Get("input_detail").(string)
	if jsonErr := json.Unmarshal([]byte(data), &inputConfigInputDetail); jsonErr != nil {
		return WrapError(jsonErr)
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
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		requestInfo = slsClient
		sls.AddNecessaryInputConfigField(inputConfigInputDetail)
		covertInput, covertErr := assertInputDetailType(inputConfigInputDetail, logconfig)
		if covertErr != nil {
			return nil, WrapError(covertErr)
		}
		logconfig.InputDetail = covertInput
		return nil, slsClient.CreateConfig(d.Get("project").(string), logconfig)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_logtail_config", "CreateConfig", AlibabacloudStackLogGoSdkERROR)
	}
	if debugOn() {
		addDebug("CreateConfig", raw, requestInfo, map[string]interface{}{
			"project": d.Get("project").(string),
			"config":  logconfig,
		})
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("logstore").(string), COLON_SEPARATED, d.Get("name").(string)))
	return resourceAlibabacloudStackLogtailConfigRead(d, meta)
}

func resourceAlibabacloudStackLogtailConfigRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	config, err := logService.DescribeLogtailConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		return WrapError(err)
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
		return WrapError(err)
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
			return WrapError(conver_err)
		}
		sls.AddNecessaryInputConfigField(inputConfigInputDetail)
		covertInput, covertErr := assertInputDetailType(inputConfigInputDetail, logconfig)
		if covertErr != nil {
			return WrapError(covertErr)
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
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.UpdateConfig(parts[0], params)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateConfig", AlibabacloudStackLogGoSdkERROR)
		}
		if debugOn() {
			addDebug("UpdateConfig", raw, requestInfo, map[string]interface{}{
				"project": parts[0],
				"config":  params,
			})
		}
	}
	return resourceAlibabacloudStackLogtailConfigRead(d, meta)
}

func resourceAlibabacloudStackLogtailConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteConfig(parts[0], parts[2])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
		if IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist", "ConfigNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteConfig", AlibabacloudStackLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogtailConfig(d.Id(), Deleted, DefaultTimeout))
}

// This function is used to assert and convert the type to sls.LogConfig
func assertInputDetailType(inputConfigInputDetail map[string]interface{}, logconfig *sls.LogConfig) (sls.InputDetailInterface, error) {
	if inputConfigInputDetail["logType"] == "json_log" {
		JSONConfigInputDetail, ok := sls.ConvertToJSONConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = JSONConfigInputDetail
	}
	if inputConfigInputDetail["logType"] == "apsara_log" {
		ApsaraLogConfigInputDetail, ok := sls.ConvertToApsaraLogConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = ApsaraLogConfigInputDetail
	}

	if inputConfigInputDetail["logType"] == "common_reg_log" {
		RegexConfigInputDetail, ok := sls.ConvertToRegexConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = RegexConfigInputDetail
	}
	if inputConfigInputDetail["logType"] == "delimiter_log" {
		DelimiterConfigInputDetail, ok := sls.ConvertToDelimiterConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = DelimiterConfigInputDetail
	}
	if logconfig.InputType == "plugin" {
		PluginLogConfigInputDetail, ok := sls.ConvertToPluginLogConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = PluginLogConfigInputDetail
	}
	return logconfig.InputDetail, nil
}
