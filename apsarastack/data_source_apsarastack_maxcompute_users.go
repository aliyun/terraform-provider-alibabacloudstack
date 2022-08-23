package apsarastack

import (
	"log"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackMaxcomputeUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackMaxcomputeUsersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
						"user_pk": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_type": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
						"organization_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"organization_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(2, 255),
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackMaxcomputeUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	maxcomputeService := MaxcomputeService{client}
	objects, err := maxcomputeService.DescribeMaxcomputeUser(d.Get("name_regex").(string))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_maxcompute_project_user maxcomputeService.DescribeMaxcomputeUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var t []map[string]interface{}
	var ids []string
	for _, object := range objects.Data {
		user := map[string]interface{}{
			"id":                strconv.Itoa(object.ID),
			"user_id":           object.UserID,
			"user_name":         object.UserName,
			"user_type":         object.UserType,
			"organization_id":   object.OrganizationId,
			"organization_name": object.OrganizationName,
			"description":       object.Description,
			"user_pk":           object.UserPK,
		}
		t = append(t, user)
		ids = append(ids, user["id"].(string))

	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("users", t); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), t)
	}
	return nil
}
