package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type instanceTypeWithOriginalPrice struct {
	InstanceType ecs.InstanceType
	OriginalPrice float64
}

func dataSourceAlibabacloudStackInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type_family": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},
			"cpu_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"cpu_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"kp", "ft", "hg", "intel"}, false),
			},
			"memory_size": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"eni_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"kubernetes_node_role": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(KubernetesNodeMaster), string(KubernetesNodeWorker)}, false),
			},
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CPU", "Memory", "Price"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"burstable_instance": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_credit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"baseline_credit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},
						"eni_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"local_storage": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capacity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"amount": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	zoneId, validZones, _, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}

	mapInstanceTypes := make(map[string][]string)
	for _, zone := range validZones {
		if zoneId != "" && zoneId != zone.ZoneId {
			continue
		}
		for _, r := range zone.AvailableResources.AvailableResource {
			if r.Type == string(InstanceTypeResource) {
				for _, t := range r.SupportedResources.SupportedResource {
					if t.Status == string(SoldOut) {
						continue
					}

					zones, _ := mapInstanceTypes[t.Value]
					zones = append(zones, zone.ZoneId)
					mapInstanceTypes[t.Value] = zones
				}
			}
		}
	}

	cpu := d.Get("cpu_core_count").(int)
	mem := d.Get("memory_size").(float64)
	family := strings.TrimSpace(d.Get("instance_type_family").(string))

	req := ecs.CreateDescribeInstanceTypesRequest()
	client.InitRpcRequest(*req.RpcRequest)
	req.InstanceTypeFamily = family

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceTypes(req)
	})
	resp, ok := raw.(*ecs.DescribeInstanceTypesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_instance_types", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var instanceTypes []instanceTypeWithOriginalPrice

	if resp != nil {
		eniAmount := d.Get("eni_amount").(int)
		k8sNode := strings.TrimSpace(d.Get("kubernetes_node_role").(string))
		for _, types := range resp.InstanceTypes.InstanceType {
			if _, ok := mapInstanceTypes[types.InstanceTypeId]; !ok {
				continue
			}
			var ct string
			for _, cpu_type_string := range []string{"hg", "kp", "ft"} {
				if strings.Contains(types.InstanceTypeId, fmt.Sprintf("-%s-", cpu_type_string)) {
					ct = cpu_type_string
				}
			}
			if ct == "" {
				ct = "intel"
			}
			log.Printf("[DEBUG] cpu_type: %s, intance_type: %s", ct, types.InstanceTypeId)
			if cpu_type, ok := d.GetOk("cpu_type"); ok {
				if cpu_type.(string) != ct {
					continue
				}
			}

			if cpu > 0 && types.CpuCoreCount != cpu {
				continue
			}

			if mem > 0 && types.MemorySize != mem {
				continue
			}
			if eniAmount > types.EniQuantity {
				continue
			}

			// Kubernetes node does not support instance types which family is "ecs.t5" and spec less that c2g4
			// Kubernetes master node does not support gpu instance types which family prefixes with "ecs.gn"
			if k8sNode != "" {
				if types.InstanceTypeFamily == "ecs.t5" {
					continue
				}
				if types.CpuCoreCount < 2 || types.MemorySize < 4 {
					continue
				}
				if k8sNode == string(KubernetesNodeMaster) && strings.HasPrefix(types.InstanceTypeFamily, "ecs.gn") {
					continue
				}
			}

			instanceTypes = append(instanceTypes, instanceTypeWithOriginalPrice{
				InstanceType: types,
			})
		}
	}

	return instanceTypesDescriptionAttributes(d, instanceTypes, mapInstanceTypes)
}

func instanceTypesDescriptionAttributes(d *schema.ResourceData, types []instanceTypeWithOriginalPrice, mapTypes map[string][]string) error {
	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(types, func(i, j int) bool {
			switch sortedBy {
			case "Price":
				return types[i].OriginalPrice < types[j].OriginalPrice
			case "CPU":
				return types[i].InstanceType.CpuCoreCount < types[j].InstanceType.CpuCoreCount
			case "Memory":
				return types[i].InstanceType.MemorySize < types[j].InstanceType.MemorySize
			}
			return false
		})
	}

	var ids []string
	var s []map[string]interface{}
	for _, t := range types {

		mapping := map[string]interface{}{
			"id":            t.InstanceType.InstanceTypeId,
			"cpu_core_count": t.InstanceType.CpuCoreCount,
			"memory_size":   t.InstanceType.MemorySize,
			"family":        t.InstanceType.InstanceTypeFamily,
			"eni_amount":    t.InstanceType.EniQuantity,
		}
		if sortedBy == "Price" {
			mapping["price"] = fmt.Sprintf("%.4f", t.OriginalPrice)
		}
		zoneIds := mapTypes[t.InstanceType.InstanceTypeId]
		sort.Strings(zoneIds)
		mapping["availability_zones"] = zoneIds

		brust := []map[string]interface{}{{
			"initial_credit": strconv.Itoa(t.InstanceType.InitialCredit),
			"baseline_credit": strconv.Itoa(t.InstanceType.BaselineCredit),
		}}
		mapping["burstable_instance"] = brust
		local := []map[string]interface{}{{
			"capacity": strconv.FormatInt(t.InstanceType.LocalStorageCapacity, 10),
			"amount":   strconv.Itoa(t.InstanceType.LocalStorageAmount),
			"category": t.InstanceType.LocalStorageCategory,
		}}
		mapping["local_storage"] = local

		ids = append(ids, t.InstanceType.InstanceTypeId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instance_types", s); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
