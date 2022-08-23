package apsarastack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				MaxItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"progressing", "accomplished", "failed", "all"}, false),
				Default:      "all",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"auto", "user", "all"}, false),
				Default:      "all",
			},
			"source_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Data"}, false),
			},
			"usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"image", "disk", "image_disk", "none"}, false),
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"snapshots": {
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
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remain_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := ecs.CreateDescribeSnapshotsRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if instanceId, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = instanceId.(string)
	}
	if diskId, ok := d.GetOk("disk_id"); ok {
		request.DiskId = diskId.(string)
		request.QueryParams["SourceDiskId"] = diskId.(string)
	}
	if ids, ok := d.GetOk("ids"); ok {
		request.SnapshotIds = convertListToJsonString(ids.(*schema.Set).List())
		request.QueryParams["SnapshotId"] = ids.(*schema.Set).List()[0].(string)
	}
	if status, ok := d.GetOk("status"); ok {
		request.Status = status.(string)
	}
	if typ, ok := d.GetOk("type"); ok {
		request.SnapshotType = typ.(string)
	}

	if diskType, ok := d.GetOk("source_disk_type"); ok {
		request.SourceDiskType = diskType.(string)
	}
	if usage, ok := d.GetOk("usage"); ok {
		request.Usage = usage.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeSnapshotsTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeSnapshotsTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allSnapshots []ecs.Snapshot
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSnapshots(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_snapshots", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response := raw.(*ecs.DescribeSnapshotsResponse)
		allSnapshots = append(allSnapshots, response.Snapshots.Snapshot...)

		if len(response.Snapshots.Snapshot) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredSnapshots []ecs.Snapshot
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, snapshot := range allSnapshots {
			if r != nil && !r.MatchString(snapshot.SnapshotName) {
				continue
			}

			filteredSnapshots = append(filteredSnapshots, snapshot)
		}
	} else {
		filteredSnapshots = allSnapshots
	}

	return snapshotsDescriptionAttributes(d, filteredSnapshots)
}

func snapshotsDescriptionAttributes(d *schema.ResourceData, snapshots []ecs.Snapshot) error {
	var s []map[string]interface{}
	var ids []string
	var names []string
	for _, snapshot := range snapshots {
		mapping := map[string]interface{}{
			"id":               snapshot.SnapshotId,
			"name":             snapshot.SnapshotName,
			"description":      snapshot.Description,
			"progress":         snapshot.Progress,
			"source_disk_id":   snapshot.SourceDiskId,
			"source_disk_type": snapshot.SourceDiskType,
			"source_disk_size": snapshot.SourceDiskSize,
			"product_code":     snapshot.ProductCode,
			"remain_time":      snapshot.RemainTime,
			"creation_time":    snapshot.CreationTime,
			"status":           snapshot.Status,
			"usage":            snapshot.Usage,
		}
		s = append(s, mapping)
		ids = append(ids, snapshot.SnapshotId)
		names = append(names, snapshot.SnapshotName)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("snapshots", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
