package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDatahubTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDatahubTopicCreate,
		Read:   resourceAlibabacloudStackDatahubTopicRead,
		Update: resourceAlibabacloudStackDatahubTopicUpdate,
		Delete: resourceAlibabacloudStackDatahubTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 32),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"shard_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"life_cycle": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(1, 7),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) != "" && strings.ToLower(new) != strings.ToLower(old)
				},
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "topic added by terraform",
				ValidateFunc: validation.StringLenBetween(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"record_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "TUPLE",
				ValidateFunc: validation.StringInSlice([]string{"TUPLE", "BLOB"}, false),
			},
			"record_schema": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("record_type") != string(datahub.TUPLE)
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackDatahubTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	t := &datahub.GetTopicResult{
		ProjectName: d.Get("project_name").(string),
		TopicName:   d.Get("name").(string),
		ShardCount:  d.Get("shard_count").(int),
		LifeCycle:   d.Get("life_cycle").(int),
		Comment:     d.Get("comment").(string),
	}
	recordType := d.Get("record_type").(string)
	if recordType == string(datahub.TUPLE) {
		t.RecordType = datahub.TUPLE

		recordSchema := d.Get("record_schema").(map[string]interface{})
		if len(recordSchema) == 0 {
			recordSchema = getDefaultRecordSchemainMap()
		}
		t.RecordSchema = getRecordSchema(recordSchema)
	} else if recordType == string(datahub.BLOB) {
		t.RecordType = datahub.BLOB
	}

	request := client.NewCommonRequest("POST", "datahub", "2019-11-20", "CreateTopic", "")
	request.QueryParams["ProjectName"] = t.ProjectName
	request.QueryParams["Lifecycle"] = strconv.Itoa(t.LifeCycle)
	request.QueryParams["ShardCount"] = strconv.Itoa(t.ShardCount)
	request.QueryParams["TopicName"] = t.TopicName
	request.QueryParams["RecordType"] = recordType
	request.QueryParams["Comment"] = t.Comment
	request.QueryParams["SignatureVersion"] = "2.1"
	request.QueryParams["AcceptLanguage"] = "zh-CN"
	request.QueryParams["ExpandMode"] = "false"
	request.QueryParams["Forwardedregionid"] = client.RegionId

	if t.RecordSchema != nil {
		var record_schema []map[string]string
		for _, v := range t.RecordSchema.Fields {
			item := map[string]string{
				"Type":      v.Type.String(),
				"AllowNull": strconv.FormatBool(v.AllowNull),
				"Name":      v.Name,
			}
			record_schema = append(record_schema, item)
		}
		record_schema_json, _ := json.Marshal(record_schema)
		record_schema_str := string(record_schema_json)
		request.QueryParams["RecordSchema"] = record_schema_str
	}

	raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_topic", "CreateTopic", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	addDebug("CreateTopic", raw, request.Content, t)

	d.SetId(strings.ToLower(fmt.Sprintf("%s%s%s", t.ProjectName, COLON_SEPARATED, t.TopicName)))
	return resourceAlibabacloudStackDatahubTopicRead(d, meta)
}

func resourceAlibabacloudStackDatahubTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}
	object, err := datahubService.DescribeDatahubTopic(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object.TopicName)
	d.Set("project_name", object.ProjectName)
	d.Set("shard_count", object.ShardCount)
	d.Set("life_cycle", object.LifeCycle)
	d.Set("comment", object.Comment)
	d.Set("record_type", object.RecordType)

	var recordSchema []datahub.Field
	err = json.Unmarshal([]byte(object.RecordSchema), &recordSchema)

	if err == nil {
		d.Set("record_schema", recordSchemaToMap(recordSchema))
	}
	d.Set("create_time", strconv.FormatInt(object.CreateTime, 10))
	d.Set("last_modify_time", strconv.FormatInt(object.LastModifyTime, 10))
	return nil
}

func resourceAlibabacloudStackDatahubTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudStackDatahubTopicRead(d, meta)
}

func resourceAlibabacloudStackDatahubTopicDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]

	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}

	request := client.NewCommonRequest("GET", "datahub", "2019-11-20", "DeleteTopic", "")
	request.QueryParams["ProjectName"] = projectName
	request.QueryParams["TopicName"] = topicName

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
			return dataHubClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			if isRetryableDatahubError(err) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_topic", "DeleteTopic", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg))
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			requestMap["TopicName"] = topicName
			addDebug("DeleteTopic", raw, requestMap)
		}
		return nil
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteTopic", errmsgs.AlibabacloudStackDatahubSdkGo)
	}
	return errmsgs.WrapError(datahubService.WaitForDatahubTopic(d.Id(), Deleted, DefaultTimeout))
}
