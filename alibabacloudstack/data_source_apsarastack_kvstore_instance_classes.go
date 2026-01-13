package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackKVStoreInstanceClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackKVStoreAvailableResourceRead,
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(KVStoreMemcache), string(KVStoreRedis), string(KVStoreKVStore)}, false),
				Default:      string(KVStoreRedis),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(KVStore4Dot0), string(KVStore5Dot0), string(KVStore6Dot0), string(KVStore7Dot0)}, false),
			},
			"architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "cluster", "rwsplit"}, false),
			},
			"edition_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"community", "enterprise"}, false),
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"double", "single", "readone", "readthree", "readfive"}, false),
			},
			"cpu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 4, 8, 16, 32, 64, 128, 256}),
			},
			"memory": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CPU", "Memory"}, false),
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},
			"instance_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"edition_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackKVStoreAvailableResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// TODO: This interface is an asapi interface and is not open to pop
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "SelectCommonSpec", "")
	request.SetDomain(client.Config.Endpoints[connectivity.ASAPICode])
	mergeMaps(request.QueryParams, map[string]string{
		"PageSize":  "500",
		"saleType":  "new",
		"pageStart": "1",
		"status":    "Available",
	})
	if v, ok := d.GetOk("engine"); ok {
		request.QueryParams["resourceType"] = strings.ToUpper(v.(string))
	}
	if v, ok := d.GetOk("engine_version"); ok {
		request.QueryParams["engineVersion"] = v.(string)
	}
	if v, ok := d.GetOk("node_type"); ok {
		request.QueryParams["nodeType"] = v.(string)
	}
	if v, ok := d.GetOk("edition_type"); ok {
		request.QueryParams["series"] = v.(string)
	}
	if v, ok := d.GetOk("architecture"); ok {
		request.QueryParams["architecture"] = v.(string)
	}
	bresponse, err := client.ProcessCommonRequest(request)
	log.Printf("Response of ListBucketVpc: %s", bresponse)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		if ossNotFoundError(err) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "SelectCommonSpec", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}
	log.Printf("Bresponse SelectCommonSpec after error")
	addDebug("SelectCommonSpec", bresponse, nil, request)

	var response *GetKVInstanceClassResponse
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var datas []KVInstanceClass
	if string(response.Data) == "{}" {
		datas = []KVInstanceClass{}
	} else {
		if err := json.Unmarshal(response.Data, &datas); err != nil {
			// Optional: try to parse as single object? (not needed here)
			return fmt.Errorf("failed to unmarshal data as array: %w", err)
		}
	}
	var Datas []KVInstanceClass
	var cpu int
	var memory float64
	if v, ok := d.GetOk("cpu"); ok {
		cpu = v.(int)
	}
	if v, ok := d.GetOk("memory"); ok {
		memory = v.(float64)
	}

	for _, data := range datas {
		// Convert raw.Memory to float32
		switch v := data.Memory.(type) {
		case float64:
			data.Memory = v
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				data.Memory = float64(f)
			} else {
				return fmt.Errorf("cannot parse memory string %q as float32", v)
			}
		case int:
			data.Memory = float64(v)
		case float32:
			data.Memory = float64(v)
		default:
			return fmt.Errorf("unsupported type for memory: %T (value: %v)", v, v)
		}
		if data.Cpu == 0 && data.CpuCore != 0 {
			data.Cpu = data.CpuCore
		}
		if cpu != 0 && data.Cpu != cpu {
			continue
		}
		if memory != 0 && data.Memory != memory {
			continue
		}
		Datas = append(Datas, data)
	}

	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(Datas, func(i, j int) bool {
			switch sortedBy {
			case "CPU":
				return Datas[i].Cpu < Datas[j].Cpu
			case "Memory":
				return Datas[i].Memory.(float64) < Datas[j].Memory.(float64)
			}
			return false
		})
	}

	var ids []string
	var s []map[string]interface{}
	for _, t := range Datas {

		mapping := map[string]interface{}{
			"id":             t.InstanceClass,
			"engine":         t.Product,
			"engine_version": t.EngineVersion,
			"architecture":   t.Architecture,
			"edition_type":   t.Series,
			"node_type":      t.NodeType,
			"cpu":            t.Cpu,
			"memory":         t.Memory,
			"status":         t.Status,
		}

		ids = append(ids, t.InstanceClass)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	err = d.Set("instance_classes", s)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), Datas)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}
