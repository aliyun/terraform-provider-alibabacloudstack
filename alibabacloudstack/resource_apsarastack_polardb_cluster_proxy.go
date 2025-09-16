package alibabacloudstack

import (
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbClusterProxy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_proxy_cluster_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_proxy_cluster_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"db_proxy_cluster_id": {
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
	setResourceFunc(resource, resourceAlibabacloudStackPolardbClusterProxyCreate,
		resourceAlibabacloudStackPolardbClusterProxyRead,
		resourceAlibabacloudStackPolardbClusterProxyUpdate,
		resourceAlibabacloudStackPolardbClusterProxyDelete)
	return resource
}

func resourceAlibabacloudStackPolardbClusterProxyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	// Prepare the request parameters for creating the proxy
	reqQuery := map[string]interface{}{
		"DBClusterId":         d.Get("db_cluster_id").(string),
		"DBProxyClusterClass": d.Get("db_proxy_cluster_class").(string),
	}

	// Call the CreateDBClusterProxy API
	if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "CreateDBClusterProxy", "", nil, reqQuery, nil); err != nil {
		return err
	}

	// The creation API does not return the resource ID, so we use the DBClusterId as the resource ID
	dbClusterId := d.Get("db_cluster_id").(string)
	d.SetId(dbClusterId)

	// Wait for the proxy to be in Running state
	stateConf := BuildStateConf([]string{"ProxyCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, polardbService.PolardbClusterInstanceStateRefreshFunc(dbClusterId, []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackPolardbClusterProxyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	object, err := polardbService.DescribePolardbClusterProxy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())
	d.Set("db_proxy_cluster_num", object["DBProxyClusterNum"])
	d.Set("db_proxy_cluster_id", object["DBProxyClusterId"])

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
	d.Set("proxy_instances", childInstances)

	return nil
}

func resourceAlibabacloudStackPolardbClusterProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	if d.IsNewResource() {
		return nil
	}

	if d.HasChange("db_proxy_cluster_class") {
		reqQuery := map[string]interface{}{
			"DBClusterId":         d.Get("db_cluster_id"),
			"DBProxyClusterClass": d.Get("db_proxy_cluster_class"),
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterProxyClass", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_cluster_proxy", "ModifyDBClusterProxyClass", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		// Wait for the proxy cluster to become Running after modification
		stateConf := BuildStateConf([]string{"ProxyModifyingClass"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbService.PolardbClusterInstanceStateRefreshFunc(d.Get("db_cluster_id").(string), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	return nil
}

func resourceAlibabacloudStackPolardbClusterProxyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	reqQuery := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	// Call the DeleteDBClusterProxy API
	_, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DeleteDBClusterProxy", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteDBClusterProxy", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
