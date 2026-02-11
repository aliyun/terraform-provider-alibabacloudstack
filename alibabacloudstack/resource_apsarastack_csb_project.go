package alibabacloudstack

import (
	"encoding/json"
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
			"csb_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"owner_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner_phone_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
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
	reqQuery := map[string]interface{}{
		"CsbId": d.Get("csb_id"),
	}

	data := map[string]interface{}{
		"projectName":          d.Get("project_name"),
		"projectOwnerName":     d.Get("owner_name"),
		"projectOwnerEmail":    d.Get("owner_email"),
		"projectOwnerPhoneNum": d.Get("owner_phone_num"),
		"description":          d.Get("description"),
	}
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	body := map[string]interface{}{
		"Data": string(content),
	}

	response, err := client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, reqQuery, body)
	if err != nil {
		return err
	}

	if v, existed := response["Data"]; !existed {
		return errmsgs.Error("error CreateProject response %v", response)
	} else {
		data := v.(map[string]interface{})
		if v, existed := data["Id"]; !existed {
			return errmsgs.Error("error Id not in CreateProject response %v", data)
		} else {
			d.Set("project_id", v)
		}
	}

	d.SetId(fmt.Sprint(d.Get("csb_id"), ":", d.Get("project_name")))
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
	d.Set("csb_id", object["CsbId"])
	d.Set("project_id", object["Id"])
	d.Set("project_name", object["ProjectName"])
	d.Set("owner_name", object["ProjectOwnerName"])
	d.Set("owner_email", object["ProjectOwnerEmail"])
	d.Set("owner_phone_num", object["ProjectOwnerPhoneNum"])
	d.Set("description", object["Description"])
	d.Set("owner_id", object["OwnerId"])
	d.Set("api_num", object["ApiNum"])
	return nil
}

func resourceAlibabacloudStackCsbProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("project_name", "owner_name", "owner_email", "owner_phone_num", "description") {
		reqQuery := map[string]interface{}{
			"CsbId": d.Get("csb_id"),
		}

		data := map[string]interface{}{
			"id":                   d.Get("project_id"),
			"projectName":          d.Get("project_name"),
			"projectOwnerName":     d.Get("owner_name"),
			"projectOwnerEmail":    d.Get("owner_email"),
			"projectOwnerPhoneNum": d.Get("owner_phone_num"),
			"description":          d.Get("description"),
		}
		content, err := json.Marshal(data)
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"Data": string(content),
		}

		action := "UpdateProject"
		_, err = client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, reqQuery, body)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprint(d.Get("csb_id"), ":", d.Get("project_name")))
	}
	return nil
}

func resourceAlibabacloudStackCsbProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csbService := CsbService{client}
	object, err := csbService.DescribeCsbProjectDetail(d.Id())
	action := "DeleteProject"
	request := map[string]interface{}{
		"ProjectId": object["Id"],
		"CsbId":     object["CsbId"],
	}

	_, err = client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
