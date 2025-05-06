package alibabacloudstack

import (
	"fmt"
	"log"
	"time"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasNamespace() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
// FIXME: edas缺少查询接口
// 			"debug_enable": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"namespace_logical_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 63),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasNamespaceCreate, resourceAlibabacloudStackEdasNamespaceRead, resourceAlibabacloudStackEdasNamespaceUpdate, resourceAlibabacloudStackEdasNamespaceDelete)
	return resource
}

func resourceAlibabacloudStackEdasNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	request := make(map[string]interface{})
	action := "InsertOrUpdateRegion"
	var err error
// 	if v, ok := d.GetOkExists("debug_enable"); ok {
// 		request["DebugEnable"] = StringPointer(strconv.FormatBool(v.(bool)))
// 	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	request["RegionTag"] = StringPointer(d.Get("namespace_logical_id").(string))
	request["RegionName"] = StringPointer(d.Get("namespace_name").(string))
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.DoTeaRequest("POST", "Edas", "2017-08-01", action, "/pop/v5/user_region_def", nil, request, nil)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if code , ok:= response["Code"]; !ok{
			return resource.NonRetryableError(errmsgs.Error("No Code in body of InsertOrUpdateRegion"))
		} else if v, ok := code.(string); ok && v != "200"{
			return resource.NonRetryableError(errmsgs.Error(response["Message"].(string)))
		} else if vv, ok := code.(json.Number); !ok {
			return resource.NonRetryableError(errmsgs.Error("Unknow Code type in body of InsertOrUpdateRegion"))
		} else if string(vv) != "200"{
			return resource.NonRetryableError(errmsgs.Error(response["Message"].(string)))
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_namespace", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	responseUserDefineRegionEntity := response["UserDefineRegionEntity"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseUserDefineRegionEntity["Id"]))

	return nil
}

func resourceAlibabacloudStackEdasNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	object, err := edasService.DescribeEdasNamespace(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_edas_namespace edasService.DescribeEdasNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
// 	d.Set("debug_enable", object["DebugEnable"])
	d.Set("description", object["Description"])
	d.Set("namespace_logical_id", object["RegionId"])
	d.Set("namespace_name", object["RegionName"])
	return nil
}

func resourceAlibabacloudStackEdasNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Id": StringPointer(d.Id()),
	}

	request["RegionTag"] = StringPointer(d.Get("namespace_logical_id").(string))
	if d.HasChange("namespace_name") {
		update = true
	}
	request["RegionName"] = StringPointer(d.Get("namespace_name").(string))
// 	if v, ok := d.GetOkExists("debug_enable"); ok {
// 		request["DebugEnable"] = StringPointer(strconv.FormatBool(v.(bool)))
// 	}
// 	if d.HasChange("debug_enable") || d.IsNewResource() {
// 		update = true
// 	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	if d.HasChange("description") {
		update = true
	}
	if update {
		action := "/pop/v5/user_region_def"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "Edas", "2017-08-01", "InsertOrUpdateRegion", "/pop/v5/user_region_def", nil, request, nil)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func resourceAlibabacloudStackEdasNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "/pop/v5/user_region_def"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Id": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.DoTeaRequest("DELETE", "Edas", "2017-08-01", "DeleteUserDefineRegion", "/pop/v5/user_region_def", nil, request, nil)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
