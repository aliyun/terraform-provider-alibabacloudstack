package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackPolardbClusterProxies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbClusterProxiesRead,

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_proxy_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_proxy_cluster_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_node_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbClusterProxiesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}
	dbClusterId := d.Get("db_cluster_id").(string)
	object, err := polardbService.DescribePolardbClusterProxy(dbClusterId)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	var childInstances []map[string]interface{}
	if v, ok := object["ChildInstances"].([]interface{}); ok && len(v) > 0 {
		for _, item := range v {
			childInstance := make(map[string]interface{})
			if val, ok := item.(map[string]interface{})["DBNodeStatus"]; ok {
				childInstance["db_node_status"] = val
			}
			if val, ok := item.(map[string]interface{})["DBNodeId"]; ok {
				childInstance["db_node_id"] = val
			}
			if val, ok := item.(map[string]interface{})["DBNodeClass"]; ok {
				childInstance["db_node_class"] = val
			}
			childInstances = append(childInstances, childInstance)
		}
	}
	d.Set("db_proxy_cluster_num", object["DBProxyClusterNum"])
	d.Set("db_proxy_cluster_id", object["DBProxyClusterId"])
	d.Set("proxy_instances", childInstances)
	d.SetId(dbClusterId)

	return nil
}
