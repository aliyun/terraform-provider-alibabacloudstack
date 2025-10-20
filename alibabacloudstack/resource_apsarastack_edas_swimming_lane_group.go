package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEdasSwimmingLaneGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entry_app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"apps": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"strategy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CONTENT", "PERCENT"}, false),
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasSwimmingLaneGroupCreate, resourceAlibabacloudStackEdasSwimmingLaneGroupRead, resourceAlibabacloudStackEdasSwimmingLaneGroupUpdate, resourceAlibabacloudStackEdasSwimmingLaneGroupDelete)
	return resource
}

func resourceAlibabacloudStackEdasSwimmingLaneGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	apps := d.Get("apps").(*schema.Set).List()
	applist := []string{}
	for _, app := range apps {
		applist = append(applist, app.(string))
	}
	logical_region_id := d.Get("logical_region_id").(string)
	if logical_region_id == "" {
		logical_region_id = client.RegionId
	}
	strategy_type := d.Get("strategy_type").(string)
	request := map[string]interface{}{
		"Name":            d.Get("name"),
		"EntryApp":        fmt.Sprintf("EDAS:%s", d.Get("entry_app_id").(string)),
		"AppIds":          strings.Join(applist, ","),
		"LogicalRegionId": logical_region_id,
		"StrategyType":    strategy_type,
	}
	response, err := client.DoTeaRequest("POST", "Edas", "2017-08-01", "InsertSwimmingLaneGroup", "/pop/v5/trafficmgnt/swimming_lane_groups", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane_group", "InsertSwimmingLaneGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error("create k8s swimming_lane_group failed for:  " + response["Message"].(string))
	}
	groupId, err := jsonpath.Get("$.Data.Id", response)
	if err != nil {
		errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane_group", "InsertSwimmingLaneGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%s:%s", logical_region_id, fmt.Sprint(groupId)))
	return nil
}

func resourceAlibabacloudStackEdasSwimmingLaneGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	object, err := edasService.DescribeEdasSwimmingLaneGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_edas_swimming_lane_group edasService.DescribeEdasSwimmingLaneGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	applist := object["ApplicationList"].([]interface{})
	apps := make([]string, 0)
	for _, v := range applist {
		app := v.(map[string]interface{})
		apps = append(apps, app["AppId"].(string))
	}
	entryApplication := object["EntryApplication"].(map[string]interface{})
	d.Set("name", object["Name"])
	d.Set("logical_region_id", object["NamespaceId"])
	d.Set("group_id", fmt.Sprint(object["Id"]))
	d.Set("entry_app_id", entryApplication["AppId"])
	d.Set("apps", apps)

	return nil
}

func resourceAlibabacloudStackEdasSwimmingLaneGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("name", "entry_app_id", "apps") {
		param := strings.Split(d.Id(), ":")
		apps := d.Get("apps").(*schema.Set).List()
		applist := []string{}
		for _, app := range apps {
			applist = append(applist, app.(string))
		}
		request := map[string]interface{}{
			"Name":     d.Get("name"),
			"EntryApp": fmt.Sprintf("EDAS:%s", d.Get("entry_app_id").(string)),
			"AppIds":   strings.Join(applist, ","),
			"GroupId":  param[1],
		}
		response, err := client.DoTeaRequest("PUT", "Edas", "2017-08-01", "UpdateSwimmingLaneGroup", "/pop/v5/trafficmgnt/swimming_lane_groups", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane_group", "UpdateSwimmingLaneGroup", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "200" {
			return errmsgs.Error("update k8s application failed for: " + response["Message"].(string))
		}
	}
	return nil
}

func resourceAlibabacloudStackEdasSwimmingLaneGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	param := strings.Split(d.Id(), ":")
	request := map[string]interface{}{
		"GroupId": param[1],
	}
	response, err := client.DoTeaRequest("DELETE", "Edas", "2017-08-01", "DeleteSwimmingLaneGroup", "/pop/v5/trafficmgnt/swimming_lane_groups", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane_group", "DeleteSwimmingLaneGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error("delete k8s application failed for: " + response["Message"].(string))
	}
	return nil
}
