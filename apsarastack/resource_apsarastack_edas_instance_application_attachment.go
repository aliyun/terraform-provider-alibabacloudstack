package apsarastack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackEdasInstanceApplicationAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEdasInstanceApplicationAttachmentCreate,
		Read:   resourceApsaraStackEdasInstanceApplicationAttachmentRead,
		Delete: resourceApsaraStackEdasInstanceApplicationAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ecc_info": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deploy_group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force_status": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ecu_info": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackEdasInstanceApplicationAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	ecuInfo := d.Get("ecu_info").([]interface{})
	aString := make([]string, len(ecuInfo))
	for i, v := range ecuInfo {
		if v != nil {
			aString[i] = v.(string)
		}
	}

	request := edas.CreateScaleOutApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = appId
	request.DeployGroup = d.Get("deploy_group").(string)
	request.EcuInfo = strings.Join(aString, ",")
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	var changeOrderId string

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ScaleOutApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_instance_application_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ScaleOutApplicationResponse)
	changeOrderId = response.ChangeOrderId
	d.SetId(appId + ":" + strings.Join(aString, ","))
	if response.Code != 200 {
		return WrapError(Error("scaleOut application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceApsaraStackEdasInstanceApplicationAttachmentRead(d, meta)
}

func resourceApsaraStackEdasInstanceApplicationAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	strs, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	appId := strs[0]
	regionId := client.RegionId
	ecuInfo := strs[1]
	aString := strings.Split(ecuInfo, ",")
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
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_instance_application_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	var eccs []string
	response := raw.(*edas.QueryApplicationStatusResponse)
	if response.Code != 200 {
		return WrapError(Error("QueryApplicationStatus failed for " + response.Message))
	}
	for _, ecc := range response.AppInfo.EccList.Ecc {
		for _, ecu := range aString {
			if ecu == ecc.EcuId {
				eccs = append(eccs, ecc.EccId)
			}
		}

	}
	if eccs != nil {
		d.Set("ecc_info", strings.Join(eccs, ","))
	}
	return nil
}

func resourceApsaraStackEdasInstanceApplicationAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	request := edas.CreateScaleInApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = d.Get("app_id").(string)
	request.EccInfo = d.Get("ecc_info").(string)
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	if v, ok := d.GetOk("force_status"); ok {
		request.ForceStatus = requests.NewBoolean(v.(bool))
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ScaleInApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_instance_application_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	var changeOrderId string
	response, _ := raw.(*edas.ScaleInApplicationResponse)
	if response.Code != 200 {
		return WrapError(Error("scaleIn application failed for " + response.Message))
	}
	changeOrderId = response.ChangeOrderId

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return nil
}
