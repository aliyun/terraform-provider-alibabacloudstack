package alibabacloudstack

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

func resourceAlibabacloudStackEvpc() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"evpc_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"evpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"department": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"department_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ascm_create_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEvpcCreate, resourceAlibabacloudStackEvpcRead, resourceAlibabacloudStackEvpcUpdate, resourceAlibabacloudStackEvpcDelete)
	return resource
}

func resourceAlibabacloudStackEvpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"EvpcName": d.Get("evpc_name").(string),
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v.(string)
	}

	var response map[string]interface{}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "CreateEvpc", "", nil, request, nil)
		addDebug("CreateEvpc", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "TaskConflict", "UnknownError", errmsgs.Throttling) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_evpc_evpc", "CreateEvpc", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("CreateEvpc", raw, nil, request)
		response = raw
		return nil
	})
	if err != nil {
		return err
	}

	data := response["data"].(map[string]interface{})
	evpcId := data["EvpcId"].(string)
	d.SetId(evpcId)

	stateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, evpcStateRefreshFunc(client, evpcId))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, evpcId)
	}

	return nil
}

func resourceAlibabacloudStackEvpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"EvpcName": d.Get("evpc_name").(string),
	}

	var response map[string]interface{}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListEvpc", "", nil, request, nil)
		addDebug("ListEvpc", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.Throttling) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_evpc_evpc", "ListEvpc", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("ListEvpc", raw, nil, request)
		response = raw
		return nil
	})
	if err != nil {
		return err
	}

	evpcList := response["EvpcList"].([]interface{})
	if len(evpcList) == 0 {
		log.Printf("[DEBUG] Resource alibabacloudstack_evpc_evpc ListEvpc Failed!!! Evpc not found")
		d.SetId("")
		return nil
	}

	evpc := evpcList[0].(map[string]interface{})
	if evpc["EvpcId"].(string) != d.Id() {
		log.Printf("[DEBUG] Resource alibabacloudstack_evpc_evpc ListEvpc Failed!!! EvpcId mismatch")
		d.SetId("")
		return nil
	}

	d.Set("evpc_id", evpc["EvpcId"])
	d.Set("evpc_name", evpc["EvpcName"])
	d.Set("status", evpc["Status"])
	d.Set("description", evpc["Description"])
	d.Set("cidr", evpc["Cidr"])
	d.Set("tenant_id", evpc["TenantId"])
	d.Set("department", evpc["Department"])
	d.Set("department_name", evpc["DepartmentName"])
	d.Set("region_id", evpc["RegionId"])
	d.Set("resource_group", evpc["ResourceGroup"])
	d.Set("resource_group_name", evpc["ResourceGroupName"])
	d.Set("cluster_id", evpc["ClusterId"])
	d.Set("ascm_create_user", evpc["AscmCreateUser"])
	d.Set("create_time", evpc["CreateTime"])
	d.Set("update_time", evpc["UpdateTime"])

	return nil
}

func resourceAlibabacloudStackEvpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"EvpcId": d.Id(),
	}
	if d.HasChange("evpc_name") {
		request["EvpcName"] = d.Get("evpc_name").(string)
	}
	if d.HasChange("description") {
		request["Description"] = d.Get("description").(string)
	}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ModifyEVPC", "", nil, request, nil)
		addDebug("ModifyEVPC", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.Throttling) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "ModifyEVPC", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("ModifyEVPC", raw, nil, request)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceAlibabacloudStackEvpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"EvpcId": d.Id(),
		"Force":  true,
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "DeleteEVPC", "", nil, request, nil)
		addDebug("DeleteEVPC", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InvalidEvpcID.NotFound", "Forbidden.EvpcNotFound") {
				return nil
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteEVPC", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("DeleteEVPC", raw, nil, request)
		return nil
	})
	if err != nil {
		return err
	}

	stateConf := BuildStateConf([]string{"Pending"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, evpcStateRefreshFunc(client, d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func evpcStateRefreshFunc(client *connectivity.AlibabacloudStackClient, evpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		request := map[string]interface{}{
			"EvpcName": "",
		}

		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListEvpc", "", nil, request, nil)
		if err != nil {
			return nil, "", err
		}

		response := raw
		evpcList := response["EvpcList"].([]interface{})

		for _, evpc := range evpcList {
			evpcMap := evpc.(map[string]interface{})
			if evpcMap["EvpcId"].(string) == evpcId {
				return evpcMap, evpcMap["Status"].(string), nil
			}
		}

		return nil, "", nil
	}
}
