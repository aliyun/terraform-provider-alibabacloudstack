package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDatahubTopic() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 32),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
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
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "topic added by terraform",
				ValidateFunc: validation.StringLenBetween(3, 1024),
			},
			"expand_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"enable_schema_registry": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"record_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "TUPLE",
				ValidateFunc: validation.StringInSlice([]string{"TUPLE", "BLOB"}, false),
			},
			"record_schemas": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"STRING", "BOOLEAN", "TINYINT", "SMALLINT", "INTEGER", "BIGINT", "DECIMAL", "FLOAT", "DOUBLE", "TIMESTAMP"}, false),
						},
						"allow_null": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
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
	setResourceFunc(resource, resourceAlibabacloudStackDatahubTopicCreate,
		resourceAlibabacloudStackDatahubTopicRead,
		resourceAlibabacloudStackDatahubTopicUpdate,
		resourceAlibabacloudStackDatahubTopicDelete)
	return resource
}

func resourceAlibabacloudStackDatahubTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "datahub", "2019-11-20", "CreateTopic", "")
	request.QueryParams["ProjectName"] = d.Get("project_name").(string)
	request.QueryParams["Lifecycle"] = strconv.Itoa(d.Get("life_cycle").(int))
	request.QueryParams["ShardCount"] = strconv.Itoa(d.Get("shard_count").(int))
	request.QueryParams["TopicName"] = d.Get("name").(string)
	request.QueryParams["Comment"] = d.Get("comment").(string)
	request.QueryParams["ExpandMode"] = strconv.FormatBool(d.Get("expand_mode").(bool))
	request.QueryParams["EnableSchemaRegistry"] = strconv.FormatBool(d.Get("enable_schema_registry").(bool))
	recordType := d.Get("record_type").(string)
	if recordType == string(datahub.TUPLE) {
		request.QueryParams["RecordType"] = "TUPLE"

		if v, ok := d.GetOk("record_schemas"); ok {
			var record_schemas []map[string]interface{}
			for _, item := range v.(*schema.Set).List() {
				record_schema := item.(map[string]interface{})
				item := map[string]interface{}{
					"Type":      record_schema["type"],
					"Name":      record_schema["name"],
					"AllowNull": record_schema["allow_null"],
					"Comment":   record_schema["comment"],
				}
				record_schemas = append(record_schemas, item)
			}
			if content, err := json.Marshal(record_schemas); err != nil {
				return err
			} else {
				request.QueryParams["RecordSchema"] = string(content)
			}
		} else {
			return fmt.Errorf("record_type TUPLE must need record_schemas")
		}
	} else if recordType == string(datahub.BLOB) {
		request.QueryParams["RecordType"] = "BLOB"
		if v, ok := d.GetOk("record_schemas"); ok {
			if len(v.(*schema.Set).List()) > 0 {
				return fmt.Errorf("record_type BLOB does not support record_schemas")
			}
		}
	}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_topic", "CreateTopic", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}

	d.SetId(strings.ToLower(fmt.Sprintf("%s:%s", d.Get("project_name"), d.Get("name"))))
	return nil
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
	d.Set("expand_mode", object.ExpandMode)
	d.Set("enable_schema_registry", object.EnableSchemaRegistry)

	var recordSchemas []DataHubRecordSchema
	if len(object.RecordSchema) > 0 {
		if err = json.Unmarshal([]byte(object.RecordSchema), &recordSchemas); err != nil {
			return err
		}
	}

	data := []map[string]interface{}{}
	for _, recordSchema := range recordSchemas {
		data = append(data, map[string]interface{}{
			"name":       recordSchema.Name,
			"type":       recordSchema.Type,
			"allow_null": recordSchema.AllowNull,
			"comment":    recordSchema.Comment,
		})
	}
	d.Set("record_schemas", data)
	d.Set("create_time", strconv.FormatInt(object.CreateTime, 10))
	d.Set("last_modify_time", strconv.FormatInt(object.LastModifyTime, 10))
	return nil
}

func resourceAlibabacloudStackDatahubTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	noUpdateAllowedFields := []string{"life_cycle", "comment"}

	if err := noUpdatesAllowedCheck(d, noUpdateAllowedFields); err != nil {
		return err
	}

	if d.IsNewResource() {
		return nil
	}

	if d.HasChange("record_schemas") {
		o, n := d.GetChange("record_schemas")
		old_schemas := o.(*schema.Set).List()
		old_maps := map[string]interface{}{}
		for _, s := range old_schemas {
			schema := s.(map[string]interface{})
			old_maps[schema["name"].(string)] = schema
		}

		new_schemas := n.(*schema.Set).List()
		new_maps := map[string]interface{}{}
		for _, s := range new_schemas {
			schema := s.(map[string]interface{})
			new_maps[schema["name"].(string)] = schema
		}

		for k, v := range old_maps {
			if new_v, existed := new_maps[k]; !existed {
				return fmt.Errorf("Deleting schemas is currently not supported.")
			} else if !reflect.DeepEqual(new_v, v) {
				return fmt.Errorf("Updating schemas is currently not supported.")
			}
		}

		for k, v := range new_maps {
			if _, existed := old_maps[k]; !existed {
				schema := v.(map[string]interface{})
				if !schema["allow_null"].(bool) {
					return fmt.Errorf("For newly added fields, allow_null must be True.")
				}
				requestQuery := map[string]interface{}{
					"ProjectName": d.Get("project_name").(string),
					"TopicName":   d.Get("name").(string),
					"Fields": []map[string]interface{}{
						{"Type": schema["type"], "AllowNull": schema["allow_null"], "Name": schema["name"], "Comment": schema["comment"]},
					},
				}
				if v, err := json.Marshal(requestQuery["Fields"]); err != nil {
					return err
				} else {
					requestQuery["Fields"] = string(v)
				}
				if _, err := client.DoTeaRequest("POST", "datahub", "2019-11-20", "AppendField", "", nil, requestQuery, nil); err != nil {
					return err
				}
			}
		}

	}

	return nil
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

	bresponse, err := client.ProcessCommonRequest(request)
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		addDebug("DeleteTopic", bresponse, requestMap)
	}
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if isDatahubNotExistError(err) {
			return nil
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_topic", "DeleteTopic", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	return errmsgs.WrapError(datahubService.WaitForDatahubTopic(d.Id(), Deleted, DefaultTimeout))
}
