package alibabacloudstack

import (
	"strings"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackKVStoreInstanceEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackKVStoreInstanceEnginesRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(KVStoreMemcache), string(KVStoreRedis)}, false),
				Default:      string(KVStoreRedis),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PrePaid,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"instance_engines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
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
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackKVStoreInstanceEnginesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := r_kvstore.CreateDescribeAvailableResourceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ZoneId = d.Get("zone_id").(string)
	instanceChargeType := d.Get("instance_charge_type").(string)
	request.InstanceChargeType = instanceChargeType
	request.Engine = d.Get("engine").(string)
	var response *r_kvstore.DescribeAvailableResourceResponse
	err := resource.Retry(time.Minute*5, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeAvailableResource(request)
		})
		response, ok := raw.(*r_kvstore.DescribeAvailableResourceResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kvstore_instance_engines", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}

	var infos []map[string]interface{}
	var ids []string

	engine, engineGot := d.GetOk("engine")
	engine = strings.ToLower(engine.(string))
	engineVersion, engineVersionGot := d.GetOk("engine_version")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		zondId := AvailableZone.ZoneId
		ids = append(ids, zondId)
		versions := make(map[string]interface{})
		for _, SupportedEngine := range AvailableZone.SupportedEngines.SupportedEngine {
			if engineGot && engine != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			if strings.ToLower(engine.(string)) == "memcache" {
				info := make(map[string]interface{})
				info["zone_id"] = AvailableZone.ZoneId
				info["engine"] = SupportedEngine.Engine
				info["engine_version"] = "2.8"
				ids = append(ids, "2.8")
				infos = append(infos, info)
			} else {
				for _, editionType := range SupportedEngine.SupportedEditionTypes.SupportedEditionType {
					for _, seriesType := range editionType.SupportedSeriesTypes.SupportedSeriesType {
						for _, SupportedEngineVersion := range seriesType.SupportedEngineVersions.SupportedEngineVersion {
							if engineVersionGot && engineVersion.(string) != SupportedEngineVersion.Version {
								continue
							}
							versions[SupportedEngineVersion.Version] = nil
						}
					}
				}
				for version := range versions {
					info := make(map[string]interface{})
					info["zone_id"] = AvailableZone.ZoneId
					info["engine"] = SupportedEngine.Engine
					info["engine_version"] = version
					ids = append(ids, version)
					infos = append(infos, info)
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_engines", infos)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), infos)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}
