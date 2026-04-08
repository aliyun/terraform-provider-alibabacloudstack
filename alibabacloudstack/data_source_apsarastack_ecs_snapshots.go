package alibabacloudstack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSnapshotsRead,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				MinItems: 1,
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
			"shared": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to query resources shared from other organizations. If set to true, shared resources will be included in the results.",
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},
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

func dataSourceAlibabacloudStackSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateDescribeSnapshotsRequest()
	client.InitRpcRequest(*request.RpcRequest)

	if instanceId, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = instanceId.(string)
	}
	if diskId, ok := d.GetOk("disk_id"); ok {
		request.DiskId = diskId.(string)
		request.QueryParams["SourceDiskId"] = diskId.(string)
	}
	if ids, ok := d.GetOk("ids"); ok {
		request.SnapshotIds = convertListToJsonString(ids.([]interface{}))
		// request.QueryParams["SnapshotId"] = ids.(*schema.Set).List()[0].(string)
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

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	if v, ok := d.GetOk("shared"); ok {
		if v.(bool) {
			request.QueryParams["shared"] = "1"
		} else {
			request.QueryParams["shared"] = "0"
		}
	}
	var allSnapshots []ecs.Snapshot
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSnapshots(request)
		})
		bresponse, ok := raw.(*ecs.DescribeSnapshotsResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_snapshots", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		allSnapshots = append(allSnapshots, bresponse.Snapshots.Snapshot...)

		if len(bresponse.Snapshots.Snapshot) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	idsMap := getIdsStringFilter(d)
	var filteredSnapshots []ecs.Snapshot
	for _, snapshot := range allSnapshots {
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(snapshot.SnapshotName) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, exist := idsMap[snapshot.SnapshotId]; !exist {
				continue
			}
		}
		if status, ok := d.GetOk("status"); ok && !strings.EqualFold(status.(string), "all") && !strings.EqualFold(status.(string), snapshot.Status) {
			continue
		}
		if usage, ok := d.GetOk("usage"); ok && !strings.EqualFold(usage.(string), snapshot.Usage) {
			continue
		}
		if source_disk_type, ok := d.GetOk("source_disk_type"); ok && !strings.EqualFold(source_disk_type.(string), snapshot.SourceDiskType) {
			continue
		}
		filteredSnapshots = append(filteredSnapshots, snapshot)
	}

	return snapshotsDescriptionAttributes(d, filteredSnapshots)
}

func snapshotsDescriptionAttributes(d *schema.ResourceData, snapshots []ecs.Snapshot) error {
	var s []map[string]interface{}
	var snapshotsids []string
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
		snapshotsids = append(snapshotsids, snapshot.SnapshotId)
		names = append(names, snapshot.SnapshotName)
	}

	d.SetId(dataResourceIdHash(snapshotsids))
	if err := d.Set("snapshots", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", snapshotsids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
