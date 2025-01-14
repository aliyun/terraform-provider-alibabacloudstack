package alibabacloudstack

import (
	"log"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackMaxcomputeUsers() *schema.Resource {
	return &schema.Resource{
		Read:	dataSourceAlibabacloudStackMaxcomputeUsersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:		schema.TypeList,
				Optional:	true,
				ForceNew:	true,
				Elem:		&schema.Schema{Type: schema.TypeString},
				Computed:	true,
				MinItems:	1,
			},
			"name_regex": {
				Type:		schema.TypeString,
				Optional:	true,
				ForceNew:	true,
			},
			"output_file": {
				Type:		schema.TypeString,
				Optional:	true,
			},
			"users": {
				Type:		schema.TypeList,
				Computed:	true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:		schema.TypeString,
							Computed:	true,
							ForceNew:	true,
						},
						"user_id": {
							Type:		schema.TypeString,
							Computed:	true,
							ForceNew:	true,
						},
						"user_pk": {
							Type:		schema.TypeString,
							Computed:	true,
							ForceNew:	true,
						},
						"user_name": {
							Type:		schema.TypeString,
							Required:	true,
						},
						"user_type": {
							Type:		schema.TypeString,
							Computed:	true,
							ForceNew:	true,
						},
						"organization_id": {
							Type:		schema.TypeInt,
							Optional:	true,
						},
						"organization_name": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"description": {
							Type:		schema.TypeString,
							Required:	true,
							ValidateFunc:	validation.StringLenBetween(2, 255),
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMaxcomputeUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	objects, err := maxcomputeService.DescribeMaxcomputeUser(d.Get("name_regex").(string))
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
	for _, object := range objects.Data {
		user := map[string]interface{}{
			"id":			strconv.Itoa(object.ID),
			"user_id":		object.UserID,
			"user_name":		object.UserName,
			"user_type":		object.UserType,
			"organization_id":	object.OrganizationId,
			"organization_name":	object.OrganizationName,
			"description":		object.Description,
			"user_pk":		object.UserPK,
		}
		t = append(t, user)
		ids = append(ids, user["id"].(string))

	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("users", t); err != nil {
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
