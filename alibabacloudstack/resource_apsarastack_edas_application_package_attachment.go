package alibabacloudstack

import (
	"strconv"
	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasApplicationPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEdasApplicationPackageAttachmentCreate,
		Read:   resourceAlibabacloudStackEdasApplicationPackageAttachmentRead,
		Delete: resourceAlibabacloudStackEdasApplicationPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"war_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_package_version": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackEdasApplicationPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	groupId := d.Get("group_id").(string)

	request := edas.CreateDeployApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = appId
	request.GroupId = groupId
	request.WarUrl = d.Get("war_url").(string)
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	var packageVersion string
	if v, ok := d.GetOk("package_version"); ok {
		packageVersion = v.(string)
	} else {
		packageVersion = strconv.FormatInt(time.Now().Unix(), 10)
	}
	request.DeployType = "url"
	request.PackageVersion = packageVersion

	if version, err := edasService.GetLastPackgeVersion(appId, groupId); err != nil {
		return WrapError(err)
	} else {
		d.Set("last_package_version", version)
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.DeployApplication(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
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
	d.SetId(appId + ":" + packageVersion)

	return resourceAlibabacloudStackEdasApplicationPackageAttachmentRead(d, meta)
}

func resourceAlibabacloudStackEdasApplicationPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := strings.Split(d.Id(), ":")[0]

	request := edas.CreateQueryApplicationStatusRequest()
	request.RegionId = regionId
	request.AppId = appId
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.QueryApplicationStatus(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_edas_application_package_attachment", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	response, _ := raw.(*edas.QueryApplicationStatusResponse)

	if response.Code != 200 {
		return WrapError(Error("QueryApplicationStatus failed for " + response.Message))
	}

	groupId := d.Get("group_id").(string)
	for _, group := range response.AppInfo.GroupList.Group {
		if group.GroupId == groupId {
			d.SetId(appId + ":" + group.PackageVersionId)
		}
	}

	return nil
}

func resourceAlibabacloudStackEdasApplicationPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := d.Get("app_id").(string)
	packageVersion := d.Get("last_package_version").(string)
	groupId := d.Get("group_id").(string)

	if len(packageVersion) == 0 {
		return nil
	}

	request := edas.CreateRollbackApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId
	request.HistoryVersion = packageVersion
	request.GroupId = groupId
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.RollbackApplication(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.RollbackApplicationResponse)
	changeOrderId := response.ChangeOrderId
	if response.Code != 200 && !strings.Contains(response.Message, "ex.app.deploy.group.empty") {
		return WrapError(Error("deploy application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return nil
}
