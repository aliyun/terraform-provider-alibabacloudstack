package alibabacloudstack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
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
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "GET"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.ApiName = "GetClustersByProduct"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	productName := d.Get("product_name").(string)
	request.QueryParams = map[string]string{  "Product": "ascm", "RegionId": client.RegionId, "Department": client.Department, "ResourceGroup": client.ResourceGroup, "productName": productName, "Action": "GetClustersByProduct", "Version": "2019-05-10"}
	response := ClustersByProduct1{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw GetClustersByProduct : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_service_cluster", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Code == 200 || len(response.Body.ClusterList) < 1 {
			break
		}
	}

	log.Printf("ram %s", response.Body.ClusterList)
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Body.ClusterList {

		if client.RegionId == "cn-neimeng-env30-d01" {
			mapping := map[string]interface{}{
				"cluster_by_region": rg.Region30,
			}
			s = append(s, mapping)
		} else if client.RegionId == "cn-qingdao-env66-d01" {
			mapping := map[string]interface{}{
				"cluster_by_region": rg.Region66,
			}
			s = append(s, mapping)
		} else if client.RegionId == string(connectivity.QingdaoEnv17) {
			mapping := map[string]interface{}{
				"cluster_by_region": rg.Region17,
			}
			s = append(s, mapping)
		} else if client.RegionId == string(connectivity.WulanEnv82) {
			mapping := map[string]interface{}{
				"cluster_by_region": rg.Region17,
			}
			s = append(s, mapping)
		} else {
			mapping := map[string]interface{}{
				"cluster_by_region": rg.Region17,
			}
			s = append(s, mapping)
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("cluster_list", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
