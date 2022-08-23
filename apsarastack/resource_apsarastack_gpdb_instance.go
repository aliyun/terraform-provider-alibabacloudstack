package apsarastack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackGpdbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceApsaraStackGpdbInstanceRead,
		Create: resourceApsaraStackGpdbInstanceCreate,
		Update: resourceApsaraStackGpdbInstanceUpdate,
		Delete: resourceApsaraStackGpdbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_group_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid"}, false),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
			},
			"description": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(2, 256),
				Optional:     true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"engine": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"gpdb"}, false),
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceApsaraStackGpdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	gpdbService := GpdbService{client}

	instance, err := gpdbService.DescribeGpdbInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", instance.DBInstanceId)
	d.Set("region_id", instance.RegionId)
	d.Set("availability_zone", instance.ZoneId)
	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("status", instance.DBInstanceStatus)
	d.Set("description", instance.DBInstanceDescription)
	d.Set("instance_class", instance.DBInstanceClass)
	d.Set("instance_group_count", instance.DBInstanceGroupCount)
	d.Set("instance_network_type", instance.InstanceNetworkType)
	security_ips, err := gpdbService.DescribeGpdbSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ip_list", security_ips)
	//d.Set("create_time", instance.CreationTime)
	d.Set("instance_charge_type", instance.PayType)
	d.Set("tags", gpdbService.tagsToMap(instance.Tags.Tag))

	return nil
}

func resourceApsaraStackGpdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	gpdbService := GpdbService{client}

	request, err := buildGpdbCreateRequest(d, meta)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if err != nil {
		return WrapError(err)
	}
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.CreateDBInstance(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"SYSTEM.CONCURRENT_OPERATE"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	response, _ := raw.(*gpdb.CreateDBInstanceResponse)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_gpdb_instance", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	d.SetId(response.DBInstanceId)

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, gpdbService.GpdbInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceApsaraStackGpdbInstanceUpdate(d, meta)
}

func resourceApsaraStackGpdbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	gpdbService := GpdbService{client}

	// Begin Update
	d.Partial(true)

	// Update Instance Description
	if d.HasChange("description") {
		request := gpdb.CreateModifyDBInstanceDescriptionRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("description").(string)
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.ModifyDBInstanceDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("description")
	}

	// Update Security Ips
	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
		ipStr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipStr == "" {
			ipStr = LOCAL_HOST_IP
		}
		if err := gpdbService.ModifyGpdbSecurityIps(d.Id(), ipStr); err != nil {
			return WrapError(err)
		}
		//d.SetPartial("security_ip_list")
	}

	if err := gpdbService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	// Finish Update
	d.Partial(false)

	return resourceApsaraStackGpdbInstanceRead(d, meta)
}

func resourceApsaraStackGpdbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := gpdb.CreateDeleteDBInstanceRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DBInstanceId = d.Id()

	err := resource.Retry(10*5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.DeleteDBInstance(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	// because DeleteDBInstance is called synchronously, there is no wait or describe here.
	return nil
}

func buildGpdbCreateRequest(d *schema.ResourceData, meta interface{}) (*gpdb.CreateDBInstanceRequest, error) {
	client := meta.(*connectivity.ApsaraStackClient)
	request := gpdb.CreateCreateDBInstanceRequest()
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId, "ZoneId": Trim(d.Get("availability_zone").(string)), "EngineVersion": Trim(d.Get("engine").(string))}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "ZoneId": Trim(d.Get("availability_zone").(string)), "EngineVersion": Trim(d.Get("engine").(string)), "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ZoneId = Trim(d.Get("availability_zone").(string))
	request.PayType = d.Get("instance_charge_type").(string)
	request.VSwitchId = Trim(d.Get("vswitch_id").(string))
	request.DBInstanceDescription = d.Get("description").(string)
	request.DBInstanceClass = Trim(d.Get("instance_class").(string))
	request.DBInstanceGroupCount = Trim(d.Get("instance_group_count").(string))
	request.Engine = Trim(d.Get("engine").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))

	// Instance NetWorkType
	request.InstanceNetworkType = string(Classic)
	if request.VSwitchId != "" {
		// check vswitchId in zone
		vpcService := VpcService{client}
		object, err := vpcService.DescribeVSwitch(request.VSwitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = object.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zoneStr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zoneStr, string([]byte(object.ZoneId)[len(object.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s.", object.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != object.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", object.VSwitchId, request.ZoneId))
		}

		request.VPCId = object.VpcId
		request.InstanceNetworkType = strings.ToUpper(string(Vpc))
	}

	// Security Ips
	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	// ClientToken
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}
