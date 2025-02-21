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

func resourceAlibabacloudStackGraphDatabaseDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackGraphDatabaseDbInstanceCreate,
		Read:   resourceAlibabacloudStackGraphDatabaseDbInstanceRead,
		Update: resourceAlibabacloudStackGraphDatabaseDbInstanceUpdate,
		Delete: resourceAlibabacloudStackGraphDatabaseDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_ip_array": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_ip_array_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"db_instance_ip_array_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_ips": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"db_instance_category": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HA"}, false),
			},
			"db_instance_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"vpc"}, false),
			},
			"db_instance_storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_ssd"}, false),
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc: validation.StringInSlice([]string{"gdb.r.xlarge", "gdb.r.2xlarge", "gdb.r.4xlarge", "gdb.r.8xlarge", "gdb.r.16xlarge"}, false),
			},
			"db_node_storage": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(20, 32000),
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// ValidateFunc: validation.StringInSlice([]string{"1.0", "1.0-OpenCypher"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PayAsYouGo",
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackGraphDatabaseDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateDBInstance"
	request := make(map[string]interface{})
	request["DBInstanceCategory"] = strings.ToLower(d.Get("db_instance_category").(string))
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	request["DBInstanceNetworkType"] = d.Get("db_instance_network_type")
	request["DBNodeStorageType"] = d.Get("db_instance_storage_type")
	request["DBInstanceClass"] = d.Get("db_node_class")
	request["DBNodeStorage"] = d.Get("db_node_storage")
	request["DBInstanceVersion"] = d.Get("db_version")
	request["PayType"] = convertGraphDatabaseDbInstancePaymentTypeRequest(d.Get("payment_type").(string))
	request["ClientToken"] = buildClientToken("CreateDBInstance")
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}

	response, err := client.DoTeaRequest("POST", "gdb", "2019-09-03", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	gdbService := GdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gdbService.GraphDatabaseDbInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackGraphDatabaseDbInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackGraphDatabaseDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gdbService := GdbService{client}
	object, err := gdbService.DescribeGraphDatabaseDbInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_graph_database_db_instance gdbService.DescribeGraphDatabaseDbInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("db_instance_category", object["Category"])
	d.Set("db_instance_description", object["DBInstanceDescription"])
	d.Set("db_instance_network_type", object["DBInstanceNetworkType"])
	d.Set("db_instance_storage_type", object["DBInstanceStorageType"])
	d.Set("db_node_class", object["DBNodeClass"])
	d.Set("db_node_storage", formatInt(object["DBNodeStorage"]))
	d.Set("db_version", object["DBVersion"])
	d.Set("payment_type", convertGraphDatabaseDbInstancePaymentTypeResponse(object["PayType"]))
	d.Set("status", object["DBInstanceStatus"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	if DBInstanceIPArray, ok := object["DBInstanceIPArray"]; ok {
		DBInstanceIPArrayAry, ok := DBInstanceIPArray.([]interface{})
		if ok && len(DBInstanceIPArrayAry) > 0 {
			DBInstanceIPArraySli := make([]map[string]interface{}, 0)
			for _, DBInstanceIPArrayArg := range DBInstanceIPArrayAry {
				DBInstanceIPArrayMap := make(map[string]interface{})
				DBInstanceIPArrayMap["security_ips"] = DBInstanceIPArrayArg.(map[string]interface{})["SecurityIps"]
				DBInstanceIPArrayMap["db_instance_ip_array_name"] = DBInstanceIPArrayArg.(map[string]interface{})["DBInstanceIPArrayName"]
				if v, ok := DBInstanceIPArrayArg.(map[string]interface{})["DBInstanceIPArrayAttribute"]; ok {
					DBInstanceIPArrayMap["db_instance_ip_array_attribute"] = v
				}
				DBInstanceIPArraySli = append(DBInstanceIPArraySli, DBInstanceIPArrayMap)
			}
			d.Set("db_instance_ip_array", DBInstanceIPArraySli)
		}
	}
	return nil
}

func resourceAlibabacloudStackGraphDatabaseDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gdbService := GdbService{client}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("db_instance_description") {
		update = true
	}
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	if update {
		action := "ModifyDBInstanceDescription"
		_, err := client.DoTeaRequest("POST", "gdb", "2019-09-03", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gdbService.GraphDatabaseDbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		//d.SetPartial("db_instance_description")
	}
	update = false
	modifyDBInstanceAccessWhiteListReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("db_instance_ip_array") {
		oraw, nraw := d.GetChange("db_instance_ip_array")
		o := oraw.(*schema.Set)
		n := nraw.(*schema.Set)
		remove := o.Difference(n).List()
		create := n.Difference(o).List()

		if len(remove) > 0 {
			for _, dBInstanceIPArray := range remove {
				dBInstanceIPArrayArg := dBInstanceIPArray.(map[string]interface{})

				if v, ok := dBInstanceIPArrayArg["db_instance_ip_array_name"]; !ok || v.(string) == "default" {
					continue
				}
				modifyDBInstanceAccessWhiteListReq["DBInstanceIPArrayName"] = dBInstanceIPArrayArg["db_instance_ip_array_name"]
				modifyDBInstanceAccessWhiteListReq["SecurityIps"] = "Empty"
				action := "ModifyDBInstanceAccessWhiteList"
				_, err := client.DoTeaRequest("POST", "gdb", "2019-09-03", action, "", nil, nil, modifyDBInstanceAccessWhiteListReq)
				if err != nil {
					return err
				}
			}
		}

		if len(create) > 0 {
			for _, dBInstanceIPArray := range create {
				dBInstanceIPArrayArg := dBInstanceIPArray.(map[string]interface{})

				modifyDBInstanceAccessWhiteListReq["DBInstanceIPArrayAttribute"] = dBInstanceIPArrayArg["db_instance_ip_array_attribute"]
				modifyDBInstanceAccessWhiteListReq["DBInstanceIPArrayName"] = dBInstanceIPArrayArg["db_instance_ip_array_name"]
				modifyDBInstanceAccessWhiteListReq["SecurityIps"] = dBInstanceIPArrayArg["security_ips"]
				action := "ModifyDBInstanceAccessWhiteList"
				_, err := client.DoTeaRequest("POST", "gdb", "2019-09-03", action, "", nil, nil, modifyDBInstanceAccessWhiteListReq)
				if err != nil {
					return err
				}
			}
		}

		//d.SetPartial("db_instance_ip_array")
	}

	modifyDBInstanceSpecReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	modifyDBInstanceSpecReq["DBInstanceClass"] = d.Get("db_node_class")
	if !d.IsNewResource() && d.HasChange("db_node_class") {
		update = true
	}
	modifyDBInstanceSpecReq["DBInstanceStorage"] = d.Get("db_node_storage")
	if !d.IsNewResource() && d.HasChange("db_node_storage") {
		update = true
	}
	if update {
		modifyDBInstanceSpecReq["DBInstanceStorageType"] = d.Get("db_instance_storage_type")
		action := "ModifyDBInstanceSpec"
		modifyDBInstanceSpecReq["ClientToken"] = buildClientToken("ModifyDBInstanceSpec")
		_, err := client.DoTeaRequest("POST", "gdb", "2019-09-03", action, "", nil, nil, modifyDBInstanceSpecReq)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gdbService.GraphDatabaseDbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		// d.SetPartial("db_instance_storage_type")
		// d.SetPartial("db_node_class")
		// d.SetPartial("db_node_storage")
	}
	d.Partial(false)
	return resourceAlibabacloudStackGraphDatabaseDbInstanceRead(d, meta)
}

func resourceAlibabacloudStackGraphDatabaseDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gdbService := GdbService{client}
	action := "DeleteDBInstance"
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	_, err := client.DoTeaRequest("POST", "gdb", "2019-09-03", action, "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			return nil
		}
		return err
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gdbService.GraphDatabaseDbInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func convertGraphDatabaseDbInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	}
	return source
}

func convertGraphDatabaseDbInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	}
	return source
}
