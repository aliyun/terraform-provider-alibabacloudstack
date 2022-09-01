package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"log"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApsaraStackAdbDbCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAdbDbClusterCreate,
		Read:   resourceApsaraStackAdbDbClusterRead,
		Update: resourceApsaraStackAdbDbClusterUpdate,
		Delete: resourceApsaraStackAdbDbClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(72 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12, 24, 36}),
				Default:          1,
				DiffSuppressFunc: adbPostPaidAndRenewDiffSuppressFunc,
			},
			"compute_resource": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("mode"); ok && v.(string) == "reserver" {
						return true
					}
					return false
				},
			},
			/*"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},*/
			"db_cluster_category": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "Cluster", "basic", "cluster"}, false),
			},
			"db_cluster_class": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "It duplicates with attribute db_node_class and is deprecated from 1.121.2.",
			},
			"storage_resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_cluster_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"3.0"}, false),
				Default:      "3.0",
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"analyticdb", "AnalyticdbOnPanguHybrid"}, false),
			},
			"cpu_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"intel"}, false),
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"executor_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"db_node_storage": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			//"elastic_io_resource": {
			//	Type:     schema.TypeInt,
			//	Optional: true,
			//	Default:  0,
			//	DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//		if v, ok := d.GetOk("mode"); ok && v.(string) == "reserver" {
			//			return true
			//		}
			//		return false
			//	},
			//},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"reserver", "flexible"}, false),
			},
			"modify_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				ConflictsWith: []string{"pay_type"},
			},
			"pay_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				ConflictsWith: []string{"payment_type"},
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				DiffSuppressFunc: adbPostPaidDiffSuppressFunc,
				Optional:         true,
			},
			"renewal_status": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"AutoRenewal", "Normal", "NotRenewal"}, false),
				Default:          "NotRenewal",
				DiffSuppressFunc: adbPostPaidDiffSuppressFunc,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_inner_connection": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"instance_inner_port": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"instance_vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
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

func resourceApsaraStackAdbDbClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	adbService := AdbService{client}
	var response *adb.CreateDBClusterResponse

	request := adb.CreateCreateDBClusterRequest()

	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = v.(string)
	}

	if v, ok := d.GetOk("cpu_type"); ok {
		request.CpuType = v.(string)
	}

	request.DBClusterCategory = d.Get("db_cluster_category").(string)
	if v, ok := d.GetOk("db_node_class"); ok {
		request.DBClusterClass = v.(string)
	} else if v, ok := d.GetOk("db_cluster_class"); ok {
		request.DBClusterClass = v.(string)
	}

	if v, ok := d.GetOk("storage_resource"); ok {
		request.StorageResource = v.(string)
	}

	if v, ok := d.GetOk("compute_resource"); ok {
		request.ComputeResource = v.(string)
	}

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = v.(string)
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request.StorageType = v.(string)
	}

	request.DBClusterVersion = d.Get("db_cluster_version").(string)
	if v, ok := d.GetOk("db_node_count"); ok {
		request.DBNodeGroupCount = strconv.Itoa(v.(int))
		request.ExecutorCount = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("db_node_storage"); ok {
		request.DBNodeStorage = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.DBClusterDescription = v.(string)
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request.PayType = convertAdbDBClusterPaymentTypeRequest(v.(string))
		if request.PayType != string(Postpaid) {
			request.PayType = string(Prepaid)
			period := d.Get("period").(int)
			request.UsedTime = strconv.Itoa(period)
			request.Period = string(Month)
			if period > 9 {
				request.UsedTime = strconv.Itoa(period / 12)
				request.Period = string(Year)
			}
		}
	} else if v, ok := d.GetOk("pay_type"); ok {
		request.PayType = convertAdbDbClusterDBClusterPayTypeRequest(v.(string))
		if request.PayType != string(Postpaid) {
			request.PayType = string(Prepaid)
			period := d.Get("period").(int)
			request.UsedTime = strconv.Itoa(period)
			request.Period = string(Month)
			if period > 9 {
				request.UsedTime = strconv.Itoa(period / 12)
				request.Period = string(Year)
			}
		}
	} else {
		request.PayType = "Postpaid"
	}

	request.RegionId = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = v.(string)
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		//vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
		var vsw, err = vpcService.DescribeVSwitch(vswitchId)
		fmt.Sprint(vsw)
		if err != nil {
			return WrapError(err)
		}
		request.DBClusterNetworkType = "VPC"
		request.VPCId = vsw.VpcId
		request.VSwitchId = vswitchId

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		}
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request.ClientToken = buildClientToken("CreateDBCluster")
	request.Headers["x-ascm-product-name"] = "adb"
	request.Headers["x-acs-organizationId"] = client.Department

	raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.CreateDBCluster(request)
	})
	if err != nil {
		return WrapError(fmt.Errorf("[ERROR] CreateDBCluster got an error: %#v", err))
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*adb.CreateDBClusterResponse)

	d.SetId(fmt.Sprint(response.DBClusterId))
	stateConf := BuildStateConf([]string{"Preparing", "Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 900*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceApsaraStackAdbDbClusterUpdate(d, meta)
}
func resourceApsaraStackAdbDbClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbDbCluster(d.Id())
	info, err := adbService.DescribeAdbClusterNetInfo2(d.Id())

	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource ApsaraStack_analyticdb_for_mysql3.0_db_cluster adbService.DescribeAdbDbCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_inner_connection", info.ConnectionString)
	d.Set("instance_inner_port", info.Port)
	d.Set("instance_vpc_id", info.VPCId)
	d.Set("compute_resource", object["ComputeResource"])
	//d.Set("connection_string", object["ConnectionString"])
	d.Set("db_cluster_category", object["Category"])
	d.Set("db_node_class", object["DBNodeClass"])
	d.Set("db_node_count", object["DBNodeCount"])
	d.Set("executor_count", object["ExecutorCount"])
	d.Set("db_node_storage", object["DBNodeStorage"])
	d.Set("storage_type", object["StorageType"])
	d.Set("description", object["DBClusterDescription"])
	//d.Set("elastic_io_resource", formatInt(object["ElasticIOResource"]))
	d.Set("maintain_time", object["MaintainTime"])
	d.Set("mode", object["Mode"])
	d.Set("payment_type", convertAdbDBClusterPaymentTypeResponse(object["PayType"].(string)))
	d.Set("pay_type", convertAdbDbClusterDBClusterPayTypeResponse(object["PayType"].(string)))
	//d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["DBClusterStatus"])
	//d.Set("tags", tagsToMap(object["Tags"].(map[string]interface{})["Tag"]))
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])

	if object["PayType"].(string) == string(Prepaid) {
		describeAutoRenewAttributeObject, err := adbService.DescribeAutoRenewAttribute(d.Id())
		if err != nil {
			return WrapError(err)
		}
		renewPeriod := 1
		if describeAutoRenewAttributeObject != nil {
			renewPeriod = formatInt(describeAutoRenewAttributeObject["Duration"])
		}
		if describeAutoRenewAttributeObject != nil && describeAutoRenewAttributeObject["PeriodUnit"] == string(Year) {
			renewPeriod = renewPeriod * 12
		}
		d.Set("auto_renew_period", renewPeriod)
		//period, err := computePeriodByUnit(object["CreationTime"], object["ExpireTime"], d.Get("period").(int), "Month")
		//if err != nil {
		//	return WrapError(err)
		//}
		//d.Set("period", period)
		d.Set("renewal_status", describeAutoRenewAttributeObject["RenewalStatus"])
	}

	describeDBClusterAccessWhiteListObject, err := adbService.DescribeDBClusterAccessWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ips", strings.Split(describeDBClusterAccessWhiteListObject["SecurityIPList"].(string), ","))

	describeDBClustersObject, err := adbService.DescribeDBClusters(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("db_cluster_version", describeDBClustersObject["DBVersion"])
	return nil
}
func resourceApsaraStackAdbDbClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	adbService := AdbService{client}
	var response map[string]interface{}
	d.Partial(true)

	//专有云 没有 UntagResources 接口
	/*if d.HasChange("tags") {
		if err := adbService.SetResourceTags(d, "ALIYUN::ADB::CLUSTER"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}*/
	if !d.IsNewResource() && d.HasChange("description") {
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		request["DBClusterDescription"] = d.Get("description")
		action := "ModifyDBClusterDescription"
		conn, err := client.NewAdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["Product"] = "adb"
			request["OrganizationId"] = client.Department
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		//d.SetPartial("description")
	}
	//专有云 没有 ModifyDBClusterMaintainTime 接口
	/*if d.HasChange("maintain_time") {
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		request["MaintainTime"] = d.Get("maintain_time")
		request["Product"] = "adb"
		request["OrganizationId"] = client.Department
		action := "ModifyDBClusterMaintainTime"
		conn, err := client.NewAdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		d.SetPartial("maintain_time")
	}*/
	//专有云 没有 ModifyDBClusterResourceGroup 接口
	/*if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		request["NewResourceGroupId"] = d.Get("resource_group_id")
		request["Product"] = "adb"
		request["OrganizationId"] = client.Department
		action := "ModifyDBClusterResourceGroup"
		conn, err := client.NewAdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}*/

	update := false

	//专有云 316 版本没有 ModifyAutoRenewAttribute 接口
	/*request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	request["Product"] = "adb"
	request["OrganizationId"] = client.Department
	request["RegionId"] = client.RegionId
	if d.Get("pay_type").(string) == string(PrePaid) || d.Get("payment_type").(string) == "Subscription" && d.HasChange("auto_renew_period") {
		update = true
		if d.Get("renewal_status").(string) == string(RenewAutoRenewal) {
			period := d.Get("auto_renew_period").(int)
			request["Duration"] = strconv.Itoa(period)
			request["PeriodUnit"] = string(Month)
			if period > 9 {
				request["Duration"] = strconv.Itoa(period / 12)
				request["PeriodUnit"] = string(Year)
			}
		}
	}
	if d.Get("pay_type").(string) == string(PrePaid) || d.Get("payment_type").(string) == "Subscription" && d.HasChange("renewal_status") {
		update = true
		request["RenewalStatus"] = d.Get("renewal_status")
	}
	if update {
		action := "ModifyAutoRenewAttribute"
		conn, err := client.NewAdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		d.SetPartial("auto_renew_period")
		d.SetPartial("renewal_status")
	}*/

	update = false
	modifyDBClusterAccessWhiteListReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	modifyDBClusterAccessWhiteListReq["Product"] = "adb"
	modifyDBClusterAccessWhiteListReq["OrganizationId"] = client.Department
	modifyDBClusterAccessWhiteListReq["RegionId"] = client.RegionId
	if d.HasChange("security_ips") {
		update = true
	}
	modifyDBClusterAccessWhiteListReq["SecurityIps"] = convertListToCommaSeparate(d.Get("security_ips").(*schema.Set).List())
	if update {
		action := "ModifyDBClusterAccessWhiteList"
		conn, err := client.NewAdsClient()
		if err != nil {
			return WrapError(err)
		}
		if modifyDBClusterAccessWhiteListReq["SecurityIps"].(string) == "" {
			modifyDBClusterAccessWhiteListReq["SecurityIps"] = LOCAL_HOST_IP
		}
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, modifyDBClusterAccessWhiteListReq, &util.RuntimeOptions{})
		addDebug(action, response, modifyDBClusterAccessWhiteListReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		//d.SetPartial("security_ips")
	}
	update = false
	// 目前 316 版本页面 仅支持 变配 节点数量
	modifyDBClusterReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	//if !d.IsNewResource() && d.HasChange("compute_resource") {
	//	update = true
	//	modifyDBClusterReq["ComputeResource"] = d.Get("compute_resource")
	//}
	//if !d.IsNewResource() && d.HasChange("db_cluster_category") {
	//	update = true
	//	modifyDBClusterReq["DBClusterCategory"] = d.Get("db_cluster_category")
	//}
	//if !d.IsNewResource() && d.HasChange("db_node_class") {
	//	update = true
	//	modifyDBClusterReq["DBNodeClass"] = d.Get("db_node_class")
	//}
	if !d.IsNewResource() && d.HasChange("db_node_count") {
		update = true
		modifyDBClusterReq["DBNodeGroupCount"] = d.Get("db_node_count")
	}
	if !d.IsNewResource() && d.HasChange("executor_count") {
		update = true
		modifyDBClusterReq["ExecutorCount"] = d.Get("executor_count")
	}
	if !d.IsNewResource() && d.HasChange("db_node_storage") {
		update = true
		modifyDBClusterReq["DBNodeStorage"] = d.Get("db_node_storage")
	}
	//if d.HasChange("elastic_io_resource") {
	//	update = true
	//	modifyDBClusterReq["ElasticIOResource"] = d.Get("elastic_io_resource")
	//}
	modifyDBClusterReq["RegionId"] = client.RegionId
	if update {
		if _, ok := d.GetOk("mode"); ok {
			modifyDBClusterReq["Mode"] = d.Get("mode")
		}
		if _, ok := d.GetOk("modify_type"); ok {
			modifyDBClusterReq["ModifyType"] = d.Get("modify_type")
		}
		action := "ModifyDBCluster"
		conn, err := client.NewAdsClient()
		if err != nil {
			return WrapError(err)
		}
		modifyDBClusterReq["Product"] = "adb"
		modifyDBClusterReq["OrganizationId"] = client.Department
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, modifyDBClusterReq, &util.RuntimeOptions{})
		addDebug(action, response, modifyDBClusterReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 900*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		//d.SetPartial("compute_resource")
		//d.SetPartial("db_cluster_category")
		//d.SetPartial("db_node_class")
		//d.SetPartial("db_node_count")
		//d.SetPartial("db_node_storage")
		//d.SetPartial("elastic_io_resource")
	}
	d.Partial(false)
	return resourceApsaraStackAdbDbClusterRead(d, meta)
}
func resourceApsaraStackAdbDbClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	//adbService := AdbService{client}
	action := "DeleteDBCluster"
	var response map[string]interface{}
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "adb"
	request["OrganizationId"] = client.Department
	//var taskId string
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		//taskId = response["TaskId"].(json.Number).String()
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	//stateConf := BuildStateConf([]string{"Waiting", "Running", "Failed", "Retry", "Pause", "Stop"}, []string{"Finished", "Closed", "Cancel"}, d.Timeout(schema.TimeoutDelete), 10*time.Minute, adbService.AdbTaskStateRefreshFunc(d.Id(), taskId))
	//if _, err = stateConf.WaitForState(); err != nil {
	//	return WrapErrorf(err, IdMsg, d.Id())
	//}
	return nil
}
func convertAdbDbClusterDBClusterPayTypeRequest(source string) string {
	switch source {
	case "PostPaid":
		return "Postpaid"
	case "PrePaid":
		return "Prepaid"
	}
	return source
}

func convertAdbDbClusterDBClusterPayTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PostPaid"
	case "Prepaid":
		return "PrePaid"
	}
	return source
}

func convertAdbDBClusterPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}

func convertAdbDBClusterPaymentTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}

func convertAdbDBClusterCategoryResponse(source string) string {
	switch source {
	case "MIXED_STORAGE":
		return "MixedStorage"
	case "basic":
		return "Basic"
	case "cluster":
		return "Cluster"
	}
	return source
}
