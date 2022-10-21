package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackExpressConnectPhysicalConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackExpressConnectPhysicalConnectionCreate,
		Read:   resourceAlibabacloudStackExpressConnectPhysicalConnectionRead,
		Update: resourceAlibabacloudStackExpressConnectPhysicalConnectionUpdate,
		Delete: resourceAlibabacloudStackExpressConnectPhysicalConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_point_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"circuit_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"line_operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"device_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"physical_connection_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1000Base-LX", "1000Base-T", "100Base-T", "10GBase-LR", "10GBase-T", "40GBase-LR", "100GBase-LR"}, false),
			},
			"redundant_physical_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Canceled", "Enabled", "Terminated"}, false),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

var DeviceName interface{}

func resourceAlibabacloudStackExpressConnectPhysicalConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreatePhysicalConnection"
	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"Product":        "Vpc",
		"OrganizationId": client.Department,
	}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["AccessPointId"] = d.Get("access_point_id")
	DeviceName = d.Get("device_name")
	request["DeviceName"] = DeviceName
	if v, ok := d.GetOk("bandwidth"); ok {
		request["bandwidth"] = v
	}
	if v, ok := d.GetOk("circuit_code"); ok {
		request["CircuitCode"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["LineOperator"] = d.Get("line_operator")
	request["PeerLocation"] = d.Get("peer_location")
	if v, ok := d.GetOk("physical_connection_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("port_type"); ok {
		request["PortType"] = v
	}
	if v, ok := d.GetOk("redundant_physical_connection_id"); ok {
		request["RedundantPhysicalConnectionId"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	request["ClientToken"] = buildClientToken("CreatePhysicalConnection")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_express_connect_physical_connection", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PhysicalConnectionId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Allocated", "Confirmed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, vpcService.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{"Allocation Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlibabacloudStackExpressConnectPhysicalConnectionUpdate(d, meta)
}
func resourceAlibabacloudStackExpressConnectPhysicalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeExpressConnectPhysicalConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_express_connect_physical_connection vpcService.DescribeExpressConnectPhysicalConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("access_point_id", object["AccessPointId"])
	d.Set("bandwidth", fmt.Sprint(formatInt(object["Bandwidth"])))
	d.Set("circuit_code", object["CircuitCode"])
	d.Set("description", object["Description"])
	d.Set("line_operator", object["LineOperator"])
	d.Set("peer_location", object["PeerLocation"])
	d.Set("physical_connection_name", object["Name"])
	d.Set("port_type", object["PortType"])
	d.Set("redundant_physical_connection_id", object["RedundantPhysicalConnectionId"])
	d.Set("status", object["Status"])
	d.Set("type", object["Type"])
	d.Set("device_name", DeviceName)

	return nil
}
func resourceAlibabacloudStackExpressConnectPhysicalConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"PhysicalConnectionId": d.Id(),
		"Product":              "Vpc",
		"OrganizationId":       client.Department,
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		if v, ok := d.GetOk("bandwidth"); ok {
			request["bandwidth"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("circuit_code") {
		update = true
		if v, ok := d.GetOk("circuit_code"); ok {
			request["CircuitCode"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("line_operator") {
		update = true
		request["LineOperator"] = d.Get("line_operator")
	}
	if !d.IsNewResource() && d.HasChange("peer_location") {
		update = true
		request["PeerLocation"] = d.Get("peer_location")
	}
	if !d.IsNewResource() && d.HasChange("physical_connection_name") {
		update = true
		if v, ok := d.GetOk("physical_connection_name"); ok {
			request["Name"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("port_type") {
		update = true
		if v, ok := d.GetOk("port_type"); ok {
			request["PortType"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("redundant_physical_connection_id") {
		update = true
		if v, ok := d.GetOk("redundant_physical_connection_id"); ok {
			request["RedundantPhysicalConnectionId"] = v
		}
	}
	if update {
		action := "ModifyPhysicalConnectionAttribute"
		request["ClientToken"] = buildClientToken("ModifyPhysicalConnectionAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}

	}
	if d.HasChange("status") {
		object, err := vpcService.DescribeExpressConnectPhysicalConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Canceled" {
				request := map[string]interface{}{
					"PhysicalConnectionId": d.Id(),
					"Product":              "Vpc",
					"OrganizationId":       client.Department,
				}
				request["RegionId"] = client.RegionId
				action := "CancelPhysicalConnection"
				request["ClientToken"] = buildClientToken("CancelPhysicalConnection")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
				}
			}
			if target == "Enabled" {
				request := map[string]interface{}{
					"PhysicalConnectionId": d.Id(),
					"Product":              "Vpc",
					"OrganizationId":       client.Department,
				}
				request["RegionId"] = client.RegionId
				action := "EnablePhysicalConnection"
				request["ClientToken"] = buildClientToken("EnablePhysicalConnection")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
				}
			}
			if target == "Terminated" {
				request := map[string]interface{}{
					"PhysicalConnectionId": d.Id(),
					"Product":              "Vpc",
					"OrganizationId":       client.Department,
				}
				request["RegionId"] = client.RegionId
				action := "TerminatePhysicalConnection"
				request["ClientToken"] = buildClientToken("TerminatePhysicalConnection")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
				}
			}
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackExpressConnectPhysicalConnectionRead(d, meta)
}
func resourceAlibabacloudStackExpressConnectPhysicalConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	vpcService := VpcService{client}
	//Canceled 状态才可以删除
	object, err := vpcService.DescribeExpressConnectPhysicalConnection(d.Id())
	if object["Status"].(string) != "Canceled" {
		request := map[string]interface{}{
			"PhysicalConnectionId": d.Id(),
			"Product":              "Vpc",
			"OrganizationId":       client.Department,
		}
		request["RegionId"] = client.RegionId
		action := "CancelPhysicalConnection"
		request["ClientToken"] = buildClientToken("CancelPhysicalConnection")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}
	}

	action := "DeletePhysicalConnection"
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PhysicalConnectionId": d.Id(),
		"Product":              "Vpc",
		"OrganizationId":       client.Department,
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeletePhysicalConnection")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
