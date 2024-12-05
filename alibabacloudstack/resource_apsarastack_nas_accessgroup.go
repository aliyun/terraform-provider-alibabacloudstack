package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackNasAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackNasAccessGroupCreate,
		Read:   resourceAlibabacloudStackNasAccessGroupRead,
		Update: resourceAlibabacloudStackNasAccessGroupUpdate,
		Delete: resourceAlibabacloudStackNasAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"access_group_name"},
			},
			"access_group_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Classic", "Vpc"}, false),
				ConflictsWith: []string{"type"},
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Classic", "Vpc"}, false),
				ConflictsWith: []string{"access_group_type"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"extreme", "standard"}, false),
				Default:      "standard",
			},
		},
	}
}

func resourceAlibabacloudStackNasAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateAccessGroup"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("access_group_name"); ok {
		request["AccessGroupName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["AccessGroupName"] = v
	} else {
		return errmsgs.WrapError(errmsgs.Error(`[ERROR] Argument "name" or "access_group_name" must be set one!`))
	}

	if v, ok := d.GetOk("access_group_type"); ok {
		request["AccessGroupType"] = v
	} else if v, ok := d.GetOk("type"); ok {
		request["AccessGroupType"] = v
	} else {
		return errmsgs.WrapError(errmsgs.Error(`[ERROR] Argument "type" or "access_group_type" must be set one!`))
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("file_system_type"); ok {
		request["FileSystemType"] = v
	}

	response, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, request)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["AccessGroupName"], ":", request["FileSystemType"]))

	return resourceAlibabacloudStackNasAccessGroupRead(d, meta)
}

func resourceAlibabacloudStackNasAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	object, err := nasService.DescribeNasAccessGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_nas_access_group nasService.DescribeNasAccessGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("access_group_name", parts[0])
	d.Set("name", parts[0])
	d.Set("file_system_type", parts[1])
	d.Set("access_group_type", object["AccessGroupType"])
	d.Set("type", object["AccessGroupType"])
	d.Set("description", object["Description"])
	return nil
}

func resourceAlibabacloudStackNasAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if d.HasChange("description") {
		request := map[string]interface{}{
			"AccessGroupName": parts[0],
			"FileSystemType":  parts[1],
			"Description":     d.Get("description"),
		}
		action := "ModifyAccessGroup"

		_, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, request)
		
		if err != nil {
			return err
		}
	}
	return resourceAlibabacloudStackNasAccessGroupRead(d, meta)
}

func resourceAlibabacloudStackNasAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteAccessGroup"
	request := map[string]interface{}{
		"AccessGroupName": parts[0],
		"FileSystemType":  parts[1],
	}

	_, err = client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.NasNotFound"}) {
			return nil
		}
		return err
	}
	return nil
}
