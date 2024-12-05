package alibabacloudstack

import (
	"encoding/json"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackServiceClusterByProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackServiceClusterByProductRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_by_region": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackServiceClusterByProductRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("GET", "ascm", "2019-05-10", "GetClustersByProduct", "")
	productName := d.Get("product_name").(string)
	request.QueryParams["productName"] = productName

	response := ClustersByProduct1{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw GetClustersByProduct : %s", raw)
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_service_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Code == 200 || len(response.Body.ClusterList) < 1 {
			break
		}
	}

	log.Printf("ram %s", response.Body.ClusterList)
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Body.ClusterList {

		switch client.RegionId {
		case "cn-neimeng-env30-d01":
			s = append(s, map[string]interface{}{"cluster_by_region": rg.Region30})
		case "cn-qingdao-env66-d01":
			s = append(s, map[string]interface{}{"cluster_by_region": rg.Region66})
		case "cn-qingdao-env17-d01", "cn-wulan-env82-d01":
			s = append(s, map[string]interface{}{"cluster_by_region": rg.Region17})
		default:
			s = append(s, map[string]interface{}{"cluster_by_region": rg.Region17})
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("cluster_list", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
