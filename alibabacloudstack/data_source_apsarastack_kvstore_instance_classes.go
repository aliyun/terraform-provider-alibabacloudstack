package alibabacloudstack

import (
	"encoding/json"
	"log"
	"sort"
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
				ValidateFunc: validation.StringInSlice([]string{string(KVStoreMemcache), string(KVStoreRedis)}, false),
				Default:      string(KVStoreRedis),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(KVStore4Dot0), string(KVStore5Dot0), string(KVStore6Dot0)}, false),
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
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 4, 8, 16, 32, 64, 128, 256}),
			},
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cpu", "memory"}, false),
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
							Type:     schema.TypeInt,
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

func removeRepByMap(slc []string) []string {
	result := []string{}         // Store the returned non-duplicate slice
	tempMap := map[string]byte{} // Store non-duplicate keys
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 // When e exists in tempMap, it cannot be added again because keys are not allowed to be duplicated
		// If the above line is successfully added, the length changes and the element is definitely not duplicated
		if len(tempMap) != l { // After adding to the map, if the map length changes, the element is not duplicated
			result = append(result, e) // When the element is not duplicated, add the element to the result slice
		}
	}
	return result
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

	var Datas []KVInstanceClass
	var cpu, momroy int
	if v, ok := d.GetOk("cpu"); ok {
		cpu = v.(int)
	}
	if v, ok := d.GetOk("memory"); ok {
		momroy = v.(int)
	}

	for _, data := range response.Data {
		if cpu != 0 && momroy != 0 && (data.Cpu != cpu || data.Memory != momroy) {
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
				return Datas[i].Memory < Datas[j].Memory
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
	return nil
}
