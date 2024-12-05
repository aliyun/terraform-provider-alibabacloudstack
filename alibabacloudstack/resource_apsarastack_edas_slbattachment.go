package alibabacloudstack

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var component_ids map[string]interface{}
var component_id int

func resourceAlibabacloudStackEdasApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEdasApplicationCreate,
		Update: resourceAlibabacloudStackEdasApplicationUpdate,
		Read:   resourceAlibabacloudStackEdasApplicationRead,
		Delete: resourceAlibabacloudStackEdasApplicationDelete,
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

func resourceAlibabacloudStackEdasApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	request := edas.CreateInsertApplicationRequest()
	client.InitRoaRequest(*request.RoaRequest)

	request.ApplicationName = d.Get("application_name").(string)
	request.PackageType = d.Get("package_type").(string)
	request.ClusterId = d.Get("cluster_id").(string)

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
	}

	var appId string
	var changeOrderId string

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertApplication(request)
	})

	bresponse, ok := raw.(*edas.InsertApplicationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_application", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	appId = bresponse.ApplicationInfo.AppId
	if _, ok := d.GetOk("component_id"); ok {
		if component_ids == nil {
			component_ids = make(map[string]interface{})
		}
		component_ids[appId] = component_id
	}
	changeOrderId = bresponse.ApplicationInfo.ChangeOrderId
	d.SetId(appId)
	if bresponse.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("create application failed for " + bresponse.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
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
		client.InitRoaRequest(*request.RoaRequest)

		request.AppId = appId
		request.GroupId = groupId
		request.PackageVersion = packageVersion
		request.DeployType = "url"
		request.WarUrl = warUrl

		request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeployApplication(request)
		})
		bresponse, ok := raw.(*edas.DeployApplicationResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)

		changeOrderId = bresponse.ChangeOrderId
		if bresponse.Code != 200 {
			return errmsgs.WrapError(errmsgs.Error("deploy application failed for " + bresponse.Message))
		}

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
	}

	return resourceAlibabacloudStackEdasApplicationRead(d, meta)
}

func resourceAlibabacloudStackEdasApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	if d.HasChange("application_name") || d.HasChange("descriotion") {
		request := edas.CreateUpdateApplicationBaseInfoRequest()
		client.InitRoaRequest(*request.RoaRequest)

		request.AppId = d.Id()
		request.AppName = d.Get("application_name").(string)

		request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
		if v, ok := d.GetOk("descriotion"); ok {
			request.Desc = v.(string)
		}
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.UpdateApplicationBaseInfo(request)
		})
		bresponse, ok := raw.(*edas.UpdateApplicationBaseInfoResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	}

	time.Sleep(3 * time.Second)
	return resourceAlibabacloudStackEdasApplicationRead(d, meta)
}

func resourceAlibabacloudStackEdasApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	appId := d.Id()

	request := edas.CreateGetApplicationRequest()
	client.InitRoaRequest(*request.RoaRequest)

	request.AppId = appId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetApplication(request)
		})
		bresponse, ok := raw.(*edas.GetApplicationResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_application", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		d.Set("application_name", bresponse.Applcation.Name)
		d.Set("cluster_id", bresponse.Applcation.ClusterId)
		if bresponse.Applcation.BuildPackageId != -1 {
			d.Set("build_pack_id", bresponse.Applcation.BuildPackageId)
		}
		d.Set("descriotion", bresponse.Applcation.Description)
		d.Set("health_check_url", bresponse.Applcation.HealthCheckUrl)
		if len(bresponse.Applcation.ApplicationType) > 0 && bresponse.Applcation.ApplicationType == "FatJar" {
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
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_application", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackEdasApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	appId := d.Id()

	request := edas.CreateStopApplicationRequest()
	client.InitRoaRequest(*request.RoaRequest)

	request.AppId = appId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.StopApplication(request)
	})
	bresponse, ok := raw.(*edas.StopApplicationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	changeOrderId := bresponse.ChangeOrderId

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	req := edas.CreateDeleteApplicationRequest()
	client.InitRoaRequest(*req.RoaRequest)

	req.AppId = d.Id()

	req.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteApplication(req)
		})
		bresponse, ok := raw.(*edas.DeleteApplicationResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(req.GetActionName(), raw, req.RoaRequest, req)
		if bresponse.Code == 601 && strings.Contains(bresponse.Message, "Operation cannot be processed because there are running instances.") {
			err = errmsgs.Error("Operation cannot be processed because there are running instances.")
			return resource.RetryableError(err)
		}
		changeOrderId = bresponse.ChangeOrderId
		delete(component_ids, req.AppId)
		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id()))
			}
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
