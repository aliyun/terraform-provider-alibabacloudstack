package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDisksRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "data"}, false),
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud", "cloud_efficiency", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
				Default:      DiskAll,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),

			// Computed values
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attached_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detached_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDisksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := requests.NewCommonRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Method = "POST"
	request.Product = "Ecs"
	request.Version = "2014-05-26"

	request.ApiName = "DescribeDisks"
	request.Headers = map[string]string{"RegionId": client.RegionId, "Content-Type": requests.Json}
	PageNumber := 1
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"Product":         "ecs",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"Action":          "DescribeDisks",
		"Version":         "2014-05-26",
		"ResourceGroupId": client.ResourceGroup,
		"PageSize":        "50",
		"PageNumber":      strconv.Itoa(PageNumber),
	}
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		request.QueryParams["DiskIds"] = convertListToJsonString(v.([]interface{}))
	}
	DiskType := ""
	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		DiskType = v.(string)
		request.QueryParams["DiskType"] = DiskType
	}
	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		request.QueryParams["DiskCategory"] = v.(string)
	}
	if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" {
		request.QueryParams["InstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeDisksTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeDisksTag{
				Key:   key,
				Value: value.(string),
			})
		}
		tags_json, _ := json.Marshal(tags)
		request.QueryParams["Tag"] = string(tags_json)
	}

	var allDisks []interface{}
	resp := make(map[string]interface{})
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_disks", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		response, _ := raw.(*responses.CommonResponse)
		err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_disks", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		Disks := resp["Disks"].(map[string]interface{})["Disk"].([]interface{})
		log.Printf("ecsDescribeDisk25 %v", response)
		if response == nil || len(Disks) < 1 {
			break
		}
		for _, disk := range Disks {
			allDisks = append(allDisks, disk)

		}

		if len(Disks) < PageNumber*50 {
			break
		}

		PageNumber = PageNumber + 1
		if err != nil {
			return WrapError(err)
		}
		request.QueryParams["PageNumber"] = strconv.Itoa(PageNumber)
	}

	var filteredDisksTemp []interface{}

	nameRegex, ok := d.GetOk("name_regex")

	if DiskType == "system" || DiskType == "data" {
		log.Printf("entered ")
		for _, disk := range allDisks {

			if disk.(map[string]interface{})["Type"].(string) == DiskType {
				filteredDisksTemp = append(filteredDisksTemp, disk)
			}

		}
		log.Printf("filtereddisks %v", filteredDisksTemp)
		allDisks = filteredDisksTemp
	}
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, disk := range allDisks {
			if r != nil && !r.MatchString(disk.(map[string]interface{})["DiskName"].(string)) {
				continue
			}
			filteredDisksTemp = append(filteredDisksTemp, disk)
		}
	} else {
		filteredDisksTemp = allDisks
	}
	return disksDescriptionAttributes(d, filteredDisksTemp, meta)
}

func disksDescriptionAttributes(d *schema.ResourceData, disks []interface{}, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	var ids []string
	var s []map[string]interface{}
	for _, diskdata := range disks {
		disk := diskdata.(map[string]interface{})
		var tag []ecs.Tag
		tags := disk["Tags"].(map[string]interface{})
		if len(tags["Tag"].([]interface{})) > 0 {
			for _, v := range tags["Tag"].([]interface{}) {
				if v != nil {
					v = v.(map[string]interface{})

					tag = append(tag, ecs.Tag{
						TagKey:   v.(map[string]interface{})["TagKey"].(string),
						TagValue: v.(map[string]interface{})["TagValue"].(string),
					})
				}
			}
		}
		mapping := map[string]interface{}{
			"id":                disk["DiskId"].(string),
			"name":              disk["DiskName"].(string),
			"description":       disk["Description"].(string),
			"region_id":         disk["RegionId"].(string),
			"availability_zone": disk["ZoneId"].(string),
			"status":            disk["Status"].(string),
			"type":              disk["Type"].(string),
			"category":          disk["Category"].(string),
			"size":              disk["Size"].(float64),
			"image_id":          disk["ImageId"].(string),
			"snapshot_id":       disk["SourceSnapshotId"].(string),
			"instance_id":       disk["InstanceId"].(string),
			"creation_time":     disk["CreationTime"].(string),
			"attached_time":     disk["AttachedTime"].(string),
			"detached_time":     disk["DetachedTime"].(string),
			"expiration_time":   disk["ExpiredTime"].(string),
			"storage_set_id":    disk["StorageSetId"].(string),
			"kms_key_id":        disk["KMSKeyId"].(string),
			"tags":              ecsService.tagsToMap(tag),
		}
		ids = append(ids, disk["DiskId"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("disks", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
