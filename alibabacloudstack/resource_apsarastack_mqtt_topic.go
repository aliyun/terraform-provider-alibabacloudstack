package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackMqttTopic() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"store_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"order_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"independent_naming": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"relation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"relation_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unit_flag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"channel_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"channel_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackMqttTopicCreate,
		resourceAlibabacloudStackMqttTopicRead,
		nil,
		resourceAlibabacloudStackMqttTopicDelete)
	return resource
}

func resourceAlibabacloudStackMqttTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	instanceId := d.Get("store_instance_id").(string)
	topic := d.Get("topic").(string)

	reqBody := map[string]interface{}{
		"Platform":       "onsConsole",
		"OnsRegionId":    client.RegionId,
		"PreventCache":   time.Now().UnixNano() / 1e6,
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances%2FinstanceDetail",
		"OrderType":      d.Get("order_type").(int),
		"Topic":          topic,
		"InstanceId":     instanceId,
	}

	if remark, ok := d.GetOk("remark"); ok {
		reqBody["Remark"] = remark.(string)
	}

	_, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleTopicCreate", "", nil, nil, reqBody)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", instanceId, topic))

	return nil
}

func resourceAlibabacloudStackMqttTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	object, err := onsService.DescribeOnsMqttTopic(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("store_instance_id", object["namespaceId"])
	d.Set("topic", object["topic"])
	d.Set("remark", object["remark"])
	d.Set("order_type", object["orderType"])
	d.Set("independent_naming", object["independentNaming"])
	d.Set("update_time", object["updateTime"])
	d.Set("relation", object["relation"])
	d.Set("relation_name", object["relationName"])
	d.Set("create_time", object["createTime"])
	d.Set("namespace_id", object["namespaceId"])
	d.Set("unit_flag", object["unitFlag"])
	d.Set("status_name", object["statusName"])
	d.Set("channel_name", object["channelName"])
	d.Set("channel_id", object["channelId"])
	d.Set("status", object["status"])

	return nil
}

func resourceAlibabacloudStackMqttTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	param := strings.Split(d.Id(), ":")
	instanceId := param[0]
	topic := param[1]

	reqQuery := map[string]interface{}{
		"InstanceId":     instanceId,
		"Topic":          topic,
		"Platform":       "onsConsole",
		"OnsRegionId":    client.RegionId,
		"PreventCache":   time.Now().UnixNano() / 1e6,
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances%2FinstanceDetail",
	}

	_, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleTopicDelete", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "ConsoleTopicDelete", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
