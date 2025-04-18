package alibabacloudstack

import (
	"log"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackMaxcomputeProjects() *schema.Resource {
	return &schema.Resource{
		Read:	dataSourceAlibabacloudStackMaxcomputeProjectsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:		schema.TypeList,
				Optional:	true,
				ForceNew:	true,
				Elem:		&schema.Schema{Type: schema.TypeString},
				Computed:	true,
				MinItems:	1,
			},
			"name": {
				Type:		schema.TypeString,
				Optional:	true,
				ForceNew:	true,
			},
			"output_file": {
				Type:		schema.TypeString,
				Optional:	true,
			},
			"projects": {
				Type:		schema.TypeList,
				Computed:	true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:		schema.TypeString,
							Computed:	true,
							ForceNew:	true,
						},
						//						"vpc_tunnel_ids": {
						//							Type:     schema.TypeSet,
						//							Optional: true,
						//							Elem:     &schema.Schema{Type: schema.TypeString},
						//						},
						//						"cluster": {
						//							Type:     schema.TypeString,
						//							Required: true,
						//							ForceNew: true,
						//						},
						//						"external_table": {
						//							Type:     schema.TypeBool,
						//							ForceNew: true,
						//							Optional: true,
						//							Default:  false,
						//						},
						//						"enabled_mc_encrypt": {
						//							Type:     schema.TypeBool,
						//							Optional: true,
						//							Default:  false,
						//						},
						//						"mc_encrypt_algorithm": {
						//							Type:         schema.TypeString,
						//							Optional:     true,
						//							ValidateFunc: validation.StringInSlice([]string{"SM4", "RC4", "AES256", "AESCTR"}, false),
						//						},
						//						"mc_encrypt_key": {
						//							Type:     schema.TypeString,
						//							Optional: true,
						//						},
						//						"quota_id": {
						//							Type:     schema.TypeString,
						//							Required: true,
						//						},
						//						"disk": {
						//							Type:     schema.TypeInt,
						//							Required: true,
						//						},
						"name": {
							Type:		schema.TypeString,
							Required:	true,
						},
						//						"aliyun_account": {
						//							Type:     schema.TypeString,
						//							Required: true,
						//						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMaxcomputeProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	objects, err := maxcomputeService.DescribeMaxcomputeProject(d.Get("name").(string))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project_user maxcomputeService.DescribeMaxcomputeUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	var t []map[string]interface{}
	var ids []string
	for _, object := range objects.Data.CalcEngines {
		user := map[string]interface{}{
			"id":	strconv.Itoa(object.EngineId),
			"name":	object.Name,
		}
		t = append(t, user)
		ids = append(ids, user["id"].(string))

	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("projects", t); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), t); err != nil {
			return err
		}
	}
	return nil
}
