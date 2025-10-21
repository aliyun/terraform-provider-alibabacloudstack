package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEdasSwimmingLane() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"apps": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OR", "ADD"}, false),
			},
			"rest_items": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"cookie", "header", "param"}, false),
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cond": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"==", "!=", ">", "<", ">=", "<=", "list"}, false),
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "rawvalue",
							ValidateFunc: validation.StringInSlice([]string{"rawvalue", "mod", "list"}, false),
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"lane_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasSwimmingLaneCreate, resourceAlibabacloudStackEdasSwimmingLaneRead, resourceAlibabacloudStackEdasSwimmingLaneUpdate, resourceAlibabacloudStackEdasSwimmingLaneDelete)
	return resource
}

func resourceAlibabacloudStackEdasSwimmingLaneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logicalRegionId := d.Get("logical_region_id").(string)
	if logicalRegionId == "" {
		logicalRegionId = client.RegionId
	}
	groupId := d.Get("group_id").(string)
	apps := d.Get("apps").(*schema.Set).List()
	applist := make([]map[string]interface{}, 0)
	for _, app := range apps {
		applist = append(applist, map[string]interface{}{
			"appId": app,
		})
	}
	restItems := make([]map[string]interface{}, 0)
	for _, item := range d.Get("rest_items").(*schema.Set).List() {
		restItem := item.(map[string]interface{})
		restItems = append(restItems, map[string]interface{}{
			"type":     restItem["type"],
			"name":     restItem["name"],
			"value":    restItem["value"],
			"cond":     restItem["cond"],
			"operator": restItem["operator"],
		})
	}
	entryRule := map[string]interface{}{
		"priority":  d.Get("priority"),
		"path":      d.Get("path"),
		"condition": d.Get("condition"),
		"restItems": restItems,
	}
	request := map[string]interface{}{
		// "Tag":             "tag",
		"LogicalRegionId": logicalRegionId,
		"Name":            d.Get("name"),
		"GroupId":         groupId,
		"AppInfos":        applist,
		"EntryRules":      []interface{}{entryRule},
		"EnableRules":     d.Get("enabled"),
	}
	response, err := client.DoTeaRequest("POST", "Edas", "2017-08-01", "InsertSwimmingLane", "/pop/v5/trafficmgnt/swimming_lanes", nil, request, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane", "InsertSwimmingLane", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error("create k8s swimming_lane failed for: " + response["Message"].(string))
	}
	laneId, err := jsonpath.Get("$.Data.Id", response)
	if err != nil {
		errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane", "InsertSwimmingLane", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%s:%s:%s", logicalRegionId, groupId, fmt.Sprint(laneId)))
	return nil
}

func resourceAlibabacloudStackEdasSwimmingLaneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	object, err := edasService.DescribeEdasSwimmingLane(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_edas_swimming_lane edasService.DescribeEdasSwimmingLane Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	applist := object["SwimmingLaneAppRelationShipList"].([]interface{})
	apps := make([]string, 0)
	for _, v := range applist {
		app := v.(map[string]interface{})
		apps = append(apps, app["AppId"].(string))
	}
	entryRules := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(object["EntryRule"].(string)), &entryRules)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	restItems := make([]map[string]interface{}, 0)
	if len(entryRules) > 0 {
		entryRule := entryRules[0]
		d.Set("condition", entryRule["condition"])
		d.Set("enabled", entryRule["enable"])
		d.Set("path", entryRule["path"])
		d.Set("priority", entryRule["priority"])
		for _, v := range entryRule["restItems"].([]interface{}) {
			item := v.(map[string]interface{})
			restItems = append(restItems, map[string]interface{}{
				"type":     item["type"],
				"name":     item["name"],
				"value":    item["value"],
				"cond":     item["cond"],
				"operator": item["operator"],
			})
		}
		d.Set("rest_items", restItems)
	}
	d.Set("name", object["Name"])
	d.Set("logical_region_id", object["NamespaceId"])
	d.Set("group_id", fmt.Sprint(object["GroupId"]))
	d.Set("lane_id", fmt.Sprint(object["Id"]))
	d.Set("apps", apps)
	return nil
}

func resourceAlibabacloudStackEdasSwimmingLaneUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("name", "enabled", "apps", "path", "condition", "rest_items") {
		param := strings.Split(d.Id(), ":")
		logicalRegionId := param[0]
		groupId := param[1]
		laneId := param[2]
		apps := d.Get("apps").(*schema.Set).List()
		applist := make([]map[string]interface{}, 0)
		for _, app := range apps {
			applist = append(applist, map[string]interface{}{
				"appId": app,
			})
		}
		restItems := make([]map[string]interface{}, 0)
		for _, item := range d.Get("rest_items").(*schema.Set).List() {
			restItem := item.(map[string]interface{})
			restItems = append(restItems, map[string]interface{}{
				"type":     restItem["type"],
				"name":     restItem["name"],
				"value":    restItem["value"],
				"cond":     restItem["cond"],
				"operator": restItem["operator"],
			})
		}
		entryRule := map[string]interface{}{
			"priority":  d.Get("priority"),
			"path":      d.Get("path"),
			"condition": d.Get("condition"),
			"restItems": restItems,
		}
		request := map[string]interface{}{
			"LogicalRegionId": logicalRegionId,
			"Name":            d.Get("name"),
			"GroupId":         groupId,
			"LaneId":          laneId,
			"AppInfos":        applist,
			"EntryRules":      []interface{}{entryRule},
			"EnableRules":     d.Get("enabled"),
		}
		response, err := client.DoTeaRequest("PUT", "Edas", "2017-08-01", "UpdateSwimmingLane", "/pop/v5/trafficmgnt/swimming_lanes", nil, request, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane", "UpdateSwimmingLane", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "200" {
			return errmsgs.Error("update k8s application failed for: " + response["Message"].(string))
		}
	}
	return nil
}

func resourceAlibabacloudStackEdasSwimmingLaneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	param := strings.Split(d.Id(), ":")
	request := map[string]interface{}{
		"LaneId": param[2],
	}
	response, err := client.DoTeaRequest("DELETE", "Edas", "2017-08-01", "DeleteSwimmingLane", "/pop/v5/trafficmgnt/swimming_lanes", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane", "DeleteSwimmingLane", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error("delete k8s application failed for: " + response["Message"].(string))
	}
	return nil
}
