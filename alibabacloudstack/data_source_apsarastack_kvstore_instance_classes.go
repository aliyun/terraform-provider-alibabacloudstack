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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
	result := []string{}         //存放返回的不重复切片
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 //当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		//如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e) //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

func dataSourceAlibabacloudStackKVStoreAvailableResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// TODO: 该接口为asapi接口，未对pop开放
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
		if cpu != 0 && momroy != 0 && ( data.Cpu != cpu || data.Memory != momroy ) {
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

	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), Datas)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}
