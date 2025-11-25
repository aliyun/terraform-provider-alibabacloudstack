package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackMqttGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !strings.HasPrefix(value, "GID_") && !strings.HasPrefix(value, "GID-") {
						errors = append(errors, fmt.Errorf("%q must start with 'GID_' or 'GID-'", k))
					}
					return
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"independent_naming": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"channel_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackMqttGroupCreate, resourceAlibabacloudStackMqttGroupRead, nil, resourceAlibabacloudStackMqttGroupDelete)
	return resource
}

func resourceAlibabacloudStackMqttGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	instanceId := d.Get("instance_id").(string)
	groupId := d.Get("group_id").(string)

	reqQuery := map[string]interface{}{
		"MqttInstanceId": instanceId,
		"GroupId":        groupId,
		"OnsRegionId":    client.RegionId,
		"Platform":       "onsConsole",
		"PreventCache":   time.Now().UnixNano() / 1e6,
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances%2FinstanceDetail",
	}

	response, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttCreateGroupId", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "Create Mqtt Group", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error(fmt.Sprintf("Mqtt group create Failed!!! Response: %#v", response))
	}
	d.SetId(fmt.Sprintf("%s:%s", instanceId, groupId))

	return nil
}

func resourceAlibabacloudStackMqttGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	mqttService := OnsService{client}
	targetGroup, err := mqttService.DescribeOnsMqttGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("instance_id", targetGroup["namespaceId"])
	d.Set("group_id", targetGroup["groupId"])
	d.Set("create_time", targetGroup["createTime"])
	d.Set("update_time", targetGroup["updateTime"])
	d.Set("independent_naming", targetGroup["independentNaming"])
	d.Set("channel_name", targetGroup["channelName"])
	return nil
}

func resourceAlibabacloudStackMqttGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return err
	}
	instanceId := parts[0]
	groupId := parts[1]

	reqQuery := map[string]interface{}{
		"MqttInstanceId": instanceId,
		"GroupId":        groupId,
		"OnsRegionId":    client.RegionId,
		"PreventCache":   time.Now().UnixNano() / 1e6, // milliseconds
		"Platform":       "onsConsole",
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances%2FinstanceDetail",
	}

	response, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttDeleteGroupId", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "ConsoleMqttDeleteGroupId", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error(fmt.Sprintf("Mqtt group delete Failed!!! Response: %#v", response))
	}

	return nil
}
