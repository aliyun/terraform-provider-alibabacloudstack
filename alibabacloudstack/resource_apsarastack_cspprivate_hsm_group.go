package alibabacloudstack

import (
	"encoding/json"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackCspprivateHsmGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hsm_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zone_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hsm_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hsm_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_level_tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCspprivateHsmGroupCreate, resourceAlibabacloudStackCspprivateHsmGroupRead, resourceAlibabacloudStackCspprivateHsmGroupUpdate, resourceAlibabacloudStackCspprivateHsmGroupDelete)
	return resource
}

func resourceAlibabacloudStackCspprivateHsmGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hsm_list, err := json.Marshal(d.Get("hsm_list").(*schema.Set).List())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	zoneIds, err := json.Marshal(d.Get("zone_ids").(*schema.Set).List())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := map[string]interface{}{
		"ZoneIds":  string(zoneIds),
		"HsmList":  string(hsm_list),
		"HsmCount": d.Get("hsm_list").(*schema.Set).Len(),
	}

	response, err := client.DoTeaRequest("POST", "cspprivate", "2022-02-17", "CreateHsmGroup", "", nil, request, nil)
	if err != nil {
		return err
	}

	// Get the GroupName from response as the resource ID
	groupName := response["GroupName"].(string)

	// Set the ID temporarily
	d.SetId(groupName)
	return nil
}

func resourceAlibabacloudStackCspprivateHsmGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cspprivateService := CspprivateService{client}

	object, err := cspprivateService.DescribeCspprivateHsmGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("vpc_id", object["VpcId"])
	d.Set("hsm_count", object["HsmCount"])
	d.Set("status", object["Status"])
	d.Set("group_name", object["GroupName"])
	d.Set("create_time", object["CreateTime"])
	d.Set("update_time", object["UpdateTime"])
	d.Set("security_level_tag", object["SecurityLevelTag"])

	// Set zone_ids and hsm_list from the response if they exist
	if v, exists := object["ZoneIds"]; exists && v.(string) != "" {
		zoneIds := make([]string, 0)
		err := json.Unmarshal([]byte(object["ZoneIds"].(string)), &zoneIds)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("zone_ids", zoneIds)
	}
	if hsmList, exists := object["HsmList"]; exists {
		d.Set("hsm_list", hsmList)
	}

	return nil
}

func resourceAlibabacloudStackCspprivateHsmGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("hsm_list", "zone_ids", "hsm_count") {
		request := map[string]interface{}{
			"GroupName":  d.Id(),
			"InstanceId": d.Id(),
		}

		if d.HasChange("hsm_list") {
			request["HsmList"] = d.Get("hsm_list").(string)
		}

		if d.HasChange("zone_ids") {
			request["ZoneIds"] = d.Get("zone_ids").(string)
		}

		if d.HasChange("hsm_count") {
			request["HsmCount"] = d.Get("hsm_count").(int)
		}

		_, err := client.DoTeaRequest("POST", "cspprivate", "2022-02-17", "UpdateHsmGroupHsms", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_cspprivate_hsm_group", "UpdateHsmGroupHsms", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackCspprivateHsmGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"GroupName":  d.Id(),
		"InstanceId": d.Id(),
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "DeleteHsmGroup", "", nil, reqQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"HsmGroup.NotFound", "InvalidGroupName.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteHsmGroup", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"HsmGroup.NotFound", "InvalidGroupName.NotFound"}) {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	return nil
}
