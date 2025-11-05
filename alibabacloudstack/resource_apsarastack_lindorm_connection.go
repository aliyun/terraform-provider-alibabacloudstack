package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLindormConnection() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"classic_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"net_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}

	setResourceFunc(resource,
		resourceAlibabacloudStackLindormConnectionCreate,
		resourceAlibabacloudStackLindormConnectionRead,
		resourceAlibabacloudStackLindormConnectionUpdate,
		nil,
	)

	return resource
}

func resourceAlibabacloudStackLindormConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	instanceId := d.Get("instance_id").(string)
	engineType := d.Get("engine_type").(string)
	d.SetId(fmt.Sprintf("%s:%s", instanceId, engineType))
	return nil
}

func resourceAlibabacloudStackLindormConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() && !d.Get("classic_network").(bool) {
		return nil
	}
	if d.HasChange("classic_network") {
		classic_network := d.Get("classic_network").(bool)
		actionType := "0"
		if classic_network {
			actionType = "1"
		}
		lindormService := LindormService{client}
		params := strings.Split(d.Id(), ":")
		request := map[string]interface{}{
			"InstanceId": params[0],
			"EngineType": params[1],
			"ActionType": actionType,
		}
		_, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "SwitchInstancePublicNetwork", "", nil, request, nil)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{"NET_MODIFYING", "NET_SWITCHING", "NET_CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, lindormService.LindormInstanceStateRefreshFunc(params[0], []string{"CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackLindormConnectionRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	lindormService := LindormService{client}

	instanceId := parts[0]
	d.Set("instance_id", parts[0])
	d.Set("engine_type", parts[1])
	engineData, err := lindormService.DescribeLindormConnection(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error describing Lindorm instance %s: %s", instanceId, err)
	}
	if username, ok := engineData["EngineUserName"].(string); ok {
		d.Set("engine_username", username)
	}

	if password, ok := engineData["EnginePassword"].(string); ok {
		d.Set("engine_password", password)
	}
	if netInfoList, ok := engineData["NetInfoList"].([]interface{}); ok {
		netInfos := make([]map[string]interface{}, 0)
		d.Set("classic_network", false)
		for _, netInfo := range netInfoList {
			netInfoData := netInfo.(map[string]interface{})
			netInfoMap := make(map[string]interface{})
			if port, ok := netInfoData["Port"].(float64); ok {
				netInfoMap["port"] = int(port)
			}
			var connectionString string
			if cs, ok := netInfoData["ConnectionString"].(string); ok {
				connectionString = cs
				netInfoMap["connection_string"] = cs
			}
			if netType, ok := netInfoData["NetType"].(string); ok {
				netInfoMap["net_type"] = netType
				if netType == "0" && strings.Contains(connectionString, "-external") {
					d.Set("classic_network", true)
				}
			}
			if accessType, ok := netInfoData["AccessType"].(float64); ok {
				netInfoMap["access_type"] = int(accessType)
			}
			netInfos = append(netInfos, netInfoMap)
		}

		d.Set("net_info_list", netInfos)
	}

	return nil
}
