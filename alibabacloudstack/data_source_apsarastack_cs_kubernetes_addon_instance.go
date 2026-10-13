package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCSKubernetesAddonInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCSKubernetesAddonInstanceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addon_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlibabacloudStackCSKubernetesAddonInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}

	clusterId := d.Get("cluster_id").(string)
	addonName := d.Get("addon_name").(string)

	addonInstance, err := csService.DescribeClusterAddonInstance(clusterId, addonName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_cs_kubernetes_addon_instance", "DescribeClusterAddonInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(clusterId + ":" + addonName)
	if err := d.Set("cluster_id", clusterId); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("name", addonInstance.Name); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("config", addonInstance.Config); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("version", addonInstance.Version); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("state", addonInstance.State); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
