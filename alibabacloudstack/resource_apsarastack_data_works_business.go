package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksBusiness() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksBusinessCreate, resourceAlibabacloudStackDataWorksBusinessRead, resourceAlibabacloudStackDataWorksBusinessUpdate, resourceAlibabacloudStackDataWorksBusinessDelete)
	return resource
}

func resourceAlibabacloudStackDataWorksBusinessCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateBusiness"
	request := map[string]interface{}{
		"ProjectId":    d.Get("project_id"),
		"BusinessName": d.Get("name"),
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	response, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	if id, err := toInt(response["BusinessId"]); err != nil {
		return err
	} else {
		d.SetId(fmt.Sprintf("%s:%d", d.Get("project_id"), id))
	}
	return nil
}

func resourceAlibabacloudStackDataWorksBusinessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksService{client}
	object, err := dataworksPublicService.DescribeDataWorksBusiness(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("project_id", object["ProjectId"])
	d.Set("name", object["BusinessName"])
	d.Set("description", object["Description"])
	return nil
}

func resourceAlibabacloudStackDataWorksBusinessUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("name", "description") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			err = errmsgs.WrapError(err)
			return nil
		}
		request := map[string]interface{}{
			"ProjectId":    parts[0],
			"BusinessId":   parts[1],
			"Description":  d.Get("description"),
			"BusinessName": d.Get("name"),
		}
		action := "UpdateBusiness"
		if _, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request); err != nil {
			return err
		}
	}

	return nil
}

func resourceAlibabacloudStackDataWorksBusinessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteBusiness"

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return nil
	}
	request := map[string]interface{}{
		"ProjectId":  parts[0],
		"BusinessId": parts[1],
	}
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	return err
}
