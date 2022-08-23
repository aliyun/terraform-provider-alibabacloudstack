package apsarastack

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var component_ids map[string]interface{}
var component_id int

func resourceApsaraStackEdasApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEdasApplicationCreate,
		Update: resourceApsaraStackEdasApplicationUpdate,
		Read:   resourceApsaraStackEdasApplicationRead,
		Delete: resourceApsaraStackEdasApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"package_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"JAR", "WAR", "Image"}, false),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"build_pack_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"component_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"descriotion": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ecu_info": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"war_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackEdasApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}
	request := edas.CreateInsertApplicationRequest()

	request.ApplicationName = d.Get("application_name").(string)
	request.RegionId = client.RegionId
	request.PackageType = d.Get("package_type").(string)
	request.ClusterId = d.Get("cluster_id").(string)
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	if v, ok := d.GetOk("build_pack_id"); ok {
		request.BuildPackId = requests.NewInteger(v.(int))
	} else {
		request.BuildPackId = requests.NewInteger(-1)
	}

	if v, ok := d.GetOk("component_id"); ok {
		request.ComponentIds = strconv.Itoa(v.(int))
		component_id = v.(int)
	}

	if v, ok := d.GetOk("descriotion"); ok {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("health_check_url"); ok {
		request.HealthCheckUrl = v.(string)
	}

	if v, ok := d.GetOk("region_id"); ok {
		request.RegionId = v.(string)
	}

	if v, ok := d.GetOk("ecu_info"); ok {
		ecuInfo := v.([]interface{})
		aString := make([]string, len(ecuInfo))
		for i, v := range ecuInfo {
			if v != nil {
				aString[i] = v.(string)
				request.EcuInfo = strings.Join(aString, ",")
			}
		}
		//if aString != nil {
		//	request.EcuInfo = strings.Join(aString, ",")
		//}
	}

	var appId string
	var changeOrderId string

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_application", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.InsertApplicationResponse)
	appId = response.ApplicationInfo.AppId
	if _, ok := d.GetOk("component_id"); ok {
		if component_ids == nil {
			component_ids = make(map[string]interface{})
		}
		component_ids[appId] = component_id
	}
	changeOrderId = response.ApplicationInfo.ChangeOrderId
	d.SetId(appId)
	if response.Code != 200 {
		return WrapError(Error("create application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	// check url information
	var groupId string
	var warUrl string
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}
	if v, ok := d.GetOk("war_url"); ok {
		warUrl = v.(string)
	}
	if len(warUrl) != 0 && len(groupId) != 0 {
		// deploy application
		var packageVersion string
		if v, ok := d.GetOk("package_version"); ok {
			packageVersion = v.(string)
		} else {
			packageVersion = strconv.FormatInt(time.Now().Unix(), 10)
		}
		request := edas.CreateDeployApplicationRequest()
		request.RegionId = client.RegionId
		request.AppId = appId
		request.GroupId = groupId
		request.PackageVersion = packageVersion
		request.DeployType = "url"

		request.WarUrl = warUrl
		request.Headers["x-ascm-product-name"] = "Edas"
		request.Headers["x-acs-organizationid"] = client.Department
		request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeployApplication(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)

		response, _ := raw.(*edas.DeployApplicationResponse)
		changeOrderId := response.ChangeOrderId
		if response.Code != 200 {
			return WrapError(Error("deploy application failed for " + response.Message))
		}

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}

	return resourceApsaraStackEdasApplicationRead(d, meta)
}

func resourceApsaraStackEdasApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	if d.HasChange("application_name") || d.HasChange("descriotion") {
		request := edas.CreateUpdateApplicationBaseInfoRequest()
		request.AppId = d.Id()
		request.RegionId = client.RegionId
		request.AppName = d.Get("application_name").(string)
		request.Headers["x-ascm-product-name"] = "Edas"
		request.Headers["x-acs-organizationid"] = client.Department
		request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
		if v, ok := d.GetOk("descriotion"); ok {
			request.Desc = v.(string)
		}
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.UpdateApplicationBaseInfo(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	}

	time.Sleep(3 * time.Second)
	return resourceApsaraStackEdasApplicationRead(d, meta)
}

func resourceApsaraStackEdasApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := d.Id()

	request := edas.CreateGetApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetApplication(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, _ := raw.(*edas.GetApplicationResponse)
		d.Set("application_name", response.Applcation.Name)
		d.Set("cluster_id", response.Applcation.ClusterId)
		if response.Applcation.BuildPackageId != -1 {
			d.Set("build_pack_id", response.Applcation.BuildPackageId)
		}
		d.Set("descriotion", response.Applcation.Description)
		d.Set("health_check_url", response.Applcation.HealthCheckUrl)
		if len(response.Applcation.ApplicationType) > 0 && response.Applcation.ApplicationType == "FatJar" {
			d.Set("package_type", "JAR")
		} else {
			d.Set("package_type", "WAR")
		}

		if _, ok := component_ids[appId]; ok {
			component_id = component_ids[appId].(int)
			d.Set("component_id", component_id)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_application", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	return nil
}

func resourceApsaraStackEdasApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := d.Id()

	request := edas.CreateStopApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.StopApplication(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	response, _ := raw.(*edas.StopApplicationResponse)
	changeOrderId := response.ChangeOrderId

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	req := edas.CreateDeleteApplicationRequest()
	req.RegionId = regionId
	req.AppId = d.Id()
	req.Headers["x-ascm-product-name"] = "Edas"
	req.Headers["x-acs-organizationid"] = client.Department
	req.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteApplication(req)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(req.GetActionName(), raw, req.RoaRequest, req)
		rsp := raw.(*edas.DeleteApplicationResponse)
		if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
			err = Error("Operation cannot be processed because there are running instances.")
			return resource.RetryableError(err)
		}
		changeOrderId := response.ChangeOrderId
		//component_ids[request.AppId] = 0
		delete(component_ids, request.AppId)
		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return resource.NonRetryableError(WrapErrorf(err, IdMsg, d.Id()))
			}
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return nil
}
