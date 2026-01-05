package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type HpcCluster struct {
	Description  string `json:"Description"`
	HpcClusterId string `json:"HpcClusterId"`
	Name         string `json:"Name"`
}

type EcsDescribeEcsHpcClusterResult struct {
	HpcClusters struct {
		HpcCluster []HpcCluster `json:"HpcCluster"`
	} `json:"HpcClusters"`
}

func dataSourceAlibabacloudStackEcsHpcClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEcsHpcClustersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hpc_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEcsHpcClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var objects []HpcCluster
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	pageNumber := 1
	for {
		resp := &EcsDescribeEcsHpcClusterResult{}
		action := "DescribeHpcClusters"
		request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
		request.QueryParams["PageNumber"] = fmt.Sprintf("%d", pageNumber)
		request.QueryParams["PageSize"] = fmt.Sprintf("%d", PageSizeLarge)

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		log.Printf(" response of raw DescribeHpcClusters : %s", bresponse)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		result := resp.HpcClusters.HpcCluster
		for _, v := range result {
			if nameRegex != nil {
				if !nameRegex.MatchString(v.Name) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[v.HpcClusterId]; !ok {
					continue
				}
			}
			objects = append(objects, v)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNumber++
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"description":    object.Description,
			"id":             object.HpcClusterId,
			"hpc_cluster_id": object.HpcClusterId,
			"name":           object.Name,
		}
		ids = append(ids, object.HpcClusterId)
		names = append(names, object.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
