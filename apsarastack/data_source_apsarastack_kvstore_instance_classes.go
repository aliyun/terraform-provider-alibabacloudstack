package apsarastack

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackKVStoreInstanceClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackKVStoreAvailableResourceRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(KVStoreMemcache),
					string(KVStoreRedis),
				}, false),
				Default: string(KVStoreRedis),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				ValidateFunc: validation.StringInSlice([]string{"Community", "Enterprise"}, false),
			},
			"series_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"enhanced_performance_type", "hybrid_storage"}, false),
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"double", "single", "readone", "readthree", "readfive"}, false),
			},
			"shard_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 4, 8, 16, 32, 64, 128, 256}),
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PrePaid,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
			},
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Price"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price": {
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

func dataSourceApsaraStackKVStoreAvailableResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := r_kvstore.CreateSelectCommRequest()
	request.RegionId = client.RegionId
	request.ResourceType = "REDIS"
	request.Status = "Available"
	instanceChargeType := d.Get("instance_charge_type").(string)
	var response = &r_kvstore.DescribeCommSelectResponse{}
	err := resource.Retry(time.Minute*5, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeCommSelect(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling}) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*r_kvstore.DescribeCommSelectResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_kvstore_instance_classes", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	var instanceClasses []string
	var ids []string

	Datas2 := response.Data
	for _, Data := range Datas2 {
		instanceClasses = append(instanceClasses, Data.InstanceClass)
	}

	instanceClasses = removeRepByMap(instanceClasses)
	d.SetId(dataResourceIdHash(ids))

	var instanceClassPrices []map[string]interface{}
	sortedBy := d.Get("sorted_by").(string)
	if sortedBy == "Price" && len(instanceClasses) > 0 {
		bssopenapiService := BssopenapiService{client}
		priceList, err := getKVStoreInstanceClassPrice(bssopenapiService, instanceChargeType, instanceClasses)
		if err != nil {
			return WrapError(err)
		}
		for i, instanceClass := range instanceClasses {
			classPrice := map[string]interface{}{
				"instance_class": instanceClass,
				"price":          fmt.Sprintf("%.4f", priceList[i]),
			}
			instanceClassPrices = append(instanceClassPrices, classPrice)
		}
		sort.SliceStable(instanceClassPrices, func(i, j int) bool {
			iPrice, _ := strconv.ParseFloat(instanceClassPrices[i]["price"].(string), 64)
			jPrice, _ := strconv.ParseFloat(instanceClassPrices[j]["price"].(string), 64)
			return iPrice < jPrice
		})

		err = d.Set("classes", instanceClassPrices)
		if err != nil {
			return WrapError(err)
		}

		instanceClasses = instanceClasses[:0]
		for _, instanceClass := range instanceClassPrices {
			instanceClasses = append(instanceClasses, instanceClass["instance_class"].(string))
		}
	}

	err = d.Set("instance_classes", instanceClasses)
	if err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), instanceClassPrices)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func getKVStoreInstanceClassPrice(bssopenapiService BssopenapiService, instanceChargeType string, instanceClasses []string) ([]float64, error) {
	client := bssopenapiService.client
	var modules interface{}
	moduleCode := "InstanceClass"
	var payAsYouGo []bssopenapi.GetPayAsYouGoPriceModuleList
	var subsciption []bssopenapi.GetSubscriptionPriceModuleList
	for _, instanceClass := range instanceClasses {
		config := fmt.Sprintf("InstanceClass:%s,Region:%s", instanceClass, client.Region)
		if instanceChargeType == string(PostPaid) {
			payAsYouGo = append(payAsYouGo, bssopenapi.GetPayAsYouGoPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
				PriceType:  "Hour",
			})
		} else {
			subsciption = append(subsciption, bssopenapi.GetSubscriptionPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
			})

		}
	}

	if len(payAsYouGo) != 0 {
		modules = payAsYouGo
	} else {
		modules = subsciption
	}

	return bssopenapiService.GetInstanceTypePrice("redisa", "", modules)
}
