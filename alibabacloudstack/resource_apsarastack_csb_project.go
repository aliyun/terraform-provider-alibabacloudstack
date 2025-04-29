package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCsbProject() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csb_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gmt_create": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_flag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cs_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_owner_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCsbProjectCreate, 
		resourceAlibabacloudStackCsbProjectRead, resourceAlibabacloudStackCsbProjectUpdate, resourceAlibabacloudStackCsbProjectDelete)
	return resource
}

func resourceAlibabacloudStackCsbProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateProject"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("csb_id"); ok {
		request["CsbId"] = v
	}

	if v, ok := d.GetOk("data"); ok {
		request["Data"] = v
	}

	_, err := client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	request1 := make(map[string]interface{})
	if v, ok := d.GetOk("project_name"); ok {
		request1["ProjectName"] = v
	}
	d.SetId(fmt.Sprint(request["CsbId"], ":", request1["ProjectName"]))
	return nil
}

func resourceAlibabacloudStackCsbProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csbService := CsbService{client}
	object, err := csbService.DescribeCsbProjectDetail(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_csb_project csbService.DescribeCsbProjectDetail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	fmt.Sprint(object["ProjectName"])
	d.Set("project_name", fmt.Sprint(object["ProjectName"]))
	d.Set("csb_id", fmt.Sprint(object["CsbId"]))
	return nil
}

func resourceAlibabacloudStackCsbProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := map[string]interface{}{}
	if d.HasChange("project_name") {
		update = true
		if v, ok := d.GetOk("project_name"); ok {
			request["ProjectName"] = v
		}
	}
	if d.HasChange("data") {
		update = true
		if v, ok := d.GetOk("data"); ok {
			request["Data"] = v
		}
	}

	if update {
		action := "UpdateProject"
		_, err := client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, nil, request)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackCsbProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csbService := CsbService{client}
	object, err := csbService.DescribeCsbProjectDetail(d.Id())
	action := "DeleteProject"
	request := map[string]interface{}{
		"ProjectId":      object["Id"],
		"CsbId":          object["CsbId"],
	}

	_, err = client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
