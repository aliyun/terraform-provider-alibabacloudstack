package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					ok, err := compareJsonTemplateAreEquivalent(old, new)
					if err != nil {
						return old == new
					}
					return ok
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackLogtailConfigCreate,
		resourceAlibabacloudStackLogtailConfigRead, resourceAlibabacloudStackLogtailConfiglUpdate, resourceAlibabacloudStackLogtailConfigDelete)
	return resource
}

func resourceAlibabacloudStackLogtailConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	
	var inputConfigInputDetail = make(map[string]interface{})
	data := d.Get("input_detail").(string)
	if jsonErr := json.Unmarshal([]byte(data), &inputConfigInputDetail); jsonErr != nil {
		return errmsgs.WrapError(jsonErr)
	}
	
	logconfig := &sls.LogConfig{
		Name:       d.Get("name").(string),
		InputType:  d.Get("input_type").(string),
		OutputType: d.Get("output_type").(string),
		OutputDetail: sls.OutputDetail{
			ProjectName:  d.Get("project").(string),
			LogStoreName: d.Get("logstore").(string),
		},
	}
	
	sls.AddNecessaryInputConfigField(inputConfigInputDetail)
	covertInput, covertErr := assertInputDetailType(inputConfigInputDetail, logconfig)
	if covertErr != nil {
		return covertErr
	}
	logconfig.InputDetail = covertInput
	
	slsClient, err := logService.GetSlsDataClient(d.Get("project").(string))
	if err != nil {
		return errmsgs.WrapError(err)
	}
	
	err = slsClient.CreateConfig(d.Get("project").(string), logconfig)
	if err != nil {
		errmsg := ""
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_logtail_config", "CreateConfig", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}
	
	d.SetId(fmt.Sprintf("%s%s%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("logstore").(string), COLON_SEPARATED, d.Get("name").(string)))
	return nil
}

func resourceAlibabacloudStackLogtailConfigRead(d *schema.ResourceData, meta interface{}) error {
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
	var normalizedJson string
	inputDetail := d.Get("input_detail").(string)
	if inputDetail != "" {
		originalMap := make(map[string]interface{})
		resultMap := make(map[string]interface{})

		err = json.Unmarshal([]byte(inputDetail), &originalMap)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		serverInputDetail, ok := config.InputDetail.(map[string]interface{})
		if !ok {
			return errmsgs.WrapError(fmt.Errorf("failed to parse server input detail"))
		}

		for k := range originalMap {
			if v2, ok := serverInputDetail[k]; ok {
				resultMap[k] = v2
			} else {
				resultMap[k] = originalMap[k]
			}
		}
		inputDetailByte, err := json.Marshal(resultMap)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		normalizedJson, _ = normalizeJsonString(string(inputDetailByte))
		// normalizedJson = string(inputDetailByte)
	} else {
		serverInputDetail, ok := config.InputDetail.(map[string]interface{})
		if !ok {
			return errmsgs.WrapError(fmt.Errorf("failed to parse server input detail"))
		}
		serverInputDetailByte, err := json.Marshal(serverInputDetail)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		normalizedJson = string(serverInputDetailByte)
		// normalizedJson, _ = normalizeJsonString(string(serverInputDetailByte))
	}
	log.Printf("=====================================================%s", normalizedJson)
	d.Set("input_detail", normalizedJson)
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
	if d.IsNewResource() {
		return nil
	}
	if d.HasChanges("input_detail", "input_type") {
		inputConfigInputDetail := make(map[string]interface{})
		data := d.Get("input_detail").(string)
		conver_err := json.Unmarshal([]byte(data), &inputConfigInputDetail)
		if conver_err != nil {
			return errmsgs.WrapError(conver_err)
		}

		client := meta.(*connectivity.AlibabacloudStackClient)
		logService := LogService{client}
		params := &sls.LogConfig{
			Name:       parts[2],
			InputType:  d.Get("input_type").(string),
			OutputType: d.Get("output_type").(string),
			OutputDetail: sls.OutputDetail{
				ProjectName:  d.Get("project").(string),
				LogStoreName: d.Get("logstore").(string),
			},
		}
		sls.AddNecessaryInputConfigField(inputConfigInputDetail)
		covertInput, covertErr := assertInputDetailType(inputConfigInputDetail, params)
		if covertErr != nil {
			return covertErr
		}
		params.InputDetail = covertInput
		
		slsClient, err := logService.GetSlsDataClient(parts[0])
		if err != nil {
			return errmsgs.WrapError(err)
		}
		
		err = slsClient.UpdateConfig(parts[0], params)
		if err != nil {
			errmsg := ""
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "UpdateConfig", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
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
	
	slsClient, err := logService.GetSlsDataClient(parts[0])
	if err != nil {
		return errmsgs.WrapError(err)
	}
	
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := slsClient.DeleteConfig(parts[0], parts[2])
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.LogClientTimeout) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteConfig", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist", "LogStoreNotExist", "ConfigNotExist") {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteConfig", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	return errmsgs.WrapError(logService.WaitForLogtailConfig(d.Id(), Deleted, DefaultTimeout))
}

// This function is used to assert and convert the type to sls.LogConfig
func assertInputDetailType(inputConfigInputDetail map[string]interface{}, logconfig *sls.LogConfig) (sls.InputDetailInterface, error) {
	log.Printf("inputConfigInputDetail ============================================================ %#v", inputConfigInputDetail)
	if inputConfigInputDetail["logType"] == "json_log" {
		JSONConfigInputDetail, ok := sls.ConvertToJSONConfigInputDetail(inputConfigInputDetail)
		if !ok {
			return nil, errmsgs.WrapError(errmsgs.Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = JSONConfigInputDetail
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
		log.Printf("PluginLogConfigInputDetail type: %T, value: %+v", PluginLogConfigInputDetail, PluginLogConfigInputDetail)
	}
	return logconfig.InputDetail, nil
}
