package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackHsmCluster() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"master_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_nos": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_white_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackHsmClusterCreate, resourceAlibabacloudStackHsmClusterRead, resourceAlibabacloudStackHsmClusterUpdate, resourceAlibabacloudStackHsmClusterDelete)
	return resource
}

func resourceAlibabacloudStackHsmClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	reqQuery := make(map[string]interface{})
	reqQuery["ClusterName"] = d.Get("cluster_name").(string)
	reqQuery["MasterInstanceId"] = d.Get("master_instance_id").(string)
	reqQuery["VpcId"] = d.Get("vpc_id").(string)
	reqQuery["VSwitchIds"] = d.Get("vswitch_ids").(string)
	reqQuery["ZoneNos"] = d.Get("zone_nos").(string)

	if ipWhiteList, ok := d.GetOk("ip_white_list"); ok {
		reqQuery["IpWhiteList"] = ipWhiteList.(string)
	}

	resp, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "CreateCluster", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	clusterId, ok := resp["ClusterId"].(string)
	if !ok || clusterId == "" {
		return errmsgs.WrapError(fmt.Errorf("failed to get ClusterId from response"))
	}

	// Set the temporary ID
	d.SetId(clusterId)
	hsmService := HsmService{client}
	adminName, err := hsmService.DescribeQuickInitAdminName(clusterId, d.Get("master_instance_id").(string))
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := map[string]interface{}{
		"ClusterId": d.Id(),
		"Name":      adminName,
		"Password":  d.Get("password"),
	}
	_, err = client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "QuickInit", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	stateConf := BuildStateConf([]string{"0", "3"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hsmService.HsmClusterStateRefreshFunc(d.Id(), []string{"-1"}))
	_, err = stateConf.WaitForState()
	return nil
}

func resourceAlibabacloudStackHsmClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hsmService := HsmService{client}
	object, err := hsmService.DescribeHsmCluster(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("cluster_name", object["ClusterName"])
	d.Set("master_instance_id", object["ClusterMaster"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("ip_white_list", object["IpWhiteList"])
	vswitch_ids := make([]string, 0)
	zone_nos := make([]string, 0)
	if v, ok := object["ClusterZones"].([]interface{}); ok && len(v) > 0 {
		zoneInfo := v[0].(map[string]interface{})
		vswitch_ids = append(vswitch_ids, zoneInfo["VSwitchId"].(string))
		zone_nos = append(zone_nos, zoneInfo["ZoneNo"].(string))
	}
	if len(vswitch_ids) > 0 {
		d.Set("vswitch_ids", strings.Join(vswitch_ids, ","))
	} else {
		d.Set("vswitch_ids", "")
	}
	if len(zone_nos) > 0 {
		d.Set("zone_nos", strings.Join(zone_nos, ","))
	} else {
		d.Set("zone_nos", "")
	}
	subInstanceIds := make([]string, 0)
	hsmClusterItems, ok := object["HsmClusterItems"].([]interface{})
	if ok && len(hsmClusterItems) > 0 {
		for _, v := range hsmClusterItems {
			item := v.(map[string]interface{})
			if fmt.Sprint(item["IsMaster"]) == "1" {
				continue
			}
			subInstanceIds = append(subInstanceIds, item["InstanceId"].(string))
		}
	}

	d.Set("sub_instance_ids", subInstanceIds)
	return nil
}

func resourceAlibabacloudStackHsmClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hsmService := HsmService{client}
	if d.HasChange("sub_instance_ids") {
		o, n := d.GetChange("sub_instance_ids")
		old, new := o.(*schema.Set), n.(*schema.Set)
		add := new.Difference(old)
		remove := old.Difference(new)
		if add.Len() > 0 {
			subInstanceids := make([]string, 0)
			for _, instanceId := range add.List() {
				subInstanceids = append(subInstanceids, instanceId.(string))
			}
			reqQuery := map[string]interface{}{
				"ClusterId":      d.Id(),
				"HsmInstanceIds": strings.Join(subInstanceids, ","),
			}
			_, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "AddHsmToCluster", "", nil, reqQuery, nil)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
					"alibabacloudstack_hsm_cluster", "AddHsmToCluster", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{"0", "1", "3"}, []string{"2"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second,
				hsmService.HsmClusterStateRefreshFunc(d.Id(), []string{"4", "5"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
		if remove.Len() > 0 {
			subInstanceids := make([]string, 0)
			for _, instanceId := range remove.List() {
				subInstanceids = append(subInstanceids, instanceId.(string))
			}
			reqQuery := map[string]interface{}{
				"ClusterId":      d.Id(),
				"HsmInstanceIds": strings.Join(subInstanceids, ","),
			}
			_, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "RemoveHsmFromCluster", "", nil, reqQuery, nil)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
					"alibabacloudstack_hsm_cluster", "RemoveHsmFromCluster", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{"0", "1", "3"}, []string{"2"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second,
				hsmService.HsmClusterStateRefreshFunc(d.Id(), []string{"4", "5"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
	}

	if d.IsNewResource() {
		return nil
	}
	// AddHsmToCluster
	if d.HasChanges("cluster_name", "ip_white_list") {
		reqQuery := map[string]interface{}{
			"ClusterId": d.Id(),
		}

		if d.HasChange("cluster_name") {
			reqQuery["ClusterName"] = d.Get("cluster_name")
		}

		if d.HasChange("ip_white_list") {
			if ipWhiteList, ok := d.GetOk("ip_white_list"); ok {
				reqQuery["IpWhiteList"] = ipWhiteList
			} else {
				reqQuery["IpWhiteList"] = ""
			}
		}

		_, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "ModifyCluster", "", nil, reqQuery, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_hsm_cluster", "ModifyCluster", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackHsmClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hsmService := HsmService{client}
	sub_instance_ids := d.Get("sub_instance_ids").(*schema.Set)
	if sub_instance_ids.Len() > 0 {
		subInstanceids := make([]string, 0)
		for _, instanceId := range sub_instance_ids.List() {
			subInstanceids = append(subInstanceids, instanceId.(string))
		}
		reqQuery := map[string]interface{}{
			"ClusterId":      d.Id(),
			"HsmInstanceIds": strings.Join(subInstanceids, ","),
		}
		_, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "RemoveHsmFromCluster", "", nil, reqQuery, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_hsm_cluster", "RemoveHsmFromCluster", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"0", "1", "3"}, []string{"2"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second,
			hsmService.HsmClusterStateRefreshFunc(d.Id(), []string{"4", "5"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	reqQuery := map[string]interface{}{
		"ClusterId": d.Id(),
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "hsm-private", "2018-06-30", "DeleteCluster", "", nil, reqQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "ResourceNotExist") {
				return resource.NonRetryableError(err)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteCluster", errmsgs.AlibabacloudStackSdkGoERROR)
			return resource.RetryableError(err)
		}
		_ = raw
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(err)
	}
	time.Sleep(10 * time.Second)
	return nil
}
