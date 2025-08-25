package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackMaxcomputeProject() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_pk": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"core_arch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "x86_64",
			},
			"cpu_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Intel",
			},
			"external_table": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"encryption": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"SM4", "RC4", "AES256", "AESCTR"}, false),
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackMaxcomputeProjectCreate, resourceAlibabacloudStackMaxcomputeProjectRead, resourceAlibabacloudStackMaxcomputeProjectUpdate, resourceAlibabacloudStackMaxcomputeProjectDelete)
	return resource
}

func resourceAlibabacloudStackMaxcomputeProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "CreateCalcEngineForAscm", "")

	maxcomputeService := MaxcomputeService{client}
	quote_id := d.Get("quota_id").(string)
	cu, err := maxcomputeService.DescribeMaxcomputeCu(quote_id)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_project", "DescribeMaxcomputeCu", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	clusters := []map[string]string{
		{
			"cluster":   cu["cluster"].(string),
			"core_arch": d.Get("core_arch").(string),
			"project":   "odps",
			"region":    client.RegionId,
			"aliasName": cu["cluster"].(string),
		},
	}
	clusterStr, _ := json.Marshal(clusters)
	engineInfo := map[string]interface{}{
		"taskAk": map[string]interface{}{
			"kp":            d.Get("account_pk").(string),
			"aliyunAccount": d.Get("account").(string),
		},
		"clusters": []map[string]interface{}{
			{
				"name":      cu["cluster"].(string),
				"quota":     quote_id,
				"isDefault": 1,
			},
		},
		"odpsProjectName":         d.Get("name").(string),
		"needToCreateOdpsProject": true,
		"defaultClusterArch":      d.Get("core_arch").(string),
		"isOdpsDev":               false,
	}
	engineInfoStr, _ := json.Marshal(engineInfo)
	body := map[string]string{
		"Region":         client.RegionId,
		"Action":         "CreateCalcEngineForAscm",
		"AccessKeyId":    client.AccessKey,
		"cpuType":        "intel",
		"Clusters":       string(clusterStr),
		"isNewFeature":   "true",
		"ExternalTable":  fmt.Sprint(d.Get("external_table")),
		"organizationId": client.Department,
		"Department":     client.Department,
		"odpsName":       "testtf",
		"clusterName":    fmt.Sprintf("[\"%s\"]", cu["cluster"].(string)),
		"taskPk":         d.Get("account_pk").(string),
		"CalcEngineType": "ODPS",
		"EnvType":        "PRD",
		"Name":           d.Get("name").(string),
		"RegionId":       "cn-wulan-env149-d01",
		"diskForQuota":   fmt.Sprint(d.Get("disk")),
		"EngineInfo":     string(engineInfoStr),
	}

	if d.Get("encryption").(bool) {
		body["EnabledMcEncrypt"] = "1"
		body["McEncryptAlgorithm"] = d.Get("encryption_algorithm").(string)
		body["McEncryptKey"] = d.Get("encryption_key").(string)
	}
	mergeMaps(request.QueryParams, body)
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "Create", errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	id, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "alibabacloudstack_maxcompute_project", "$", response)
	}
	d.SetId(fmt.Sprint(id))
	// d.SetId("53")
	return nil
}

func resourceAlibabacloudStackMaxcomputeProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	project, err := maxcomputeService.DescribeMaxcomputeProject(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project_user maxcomputeService.DescribeMaxcomputeUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("name", project.Name)
	d.Set("quota_id", fmt.Sprint(project.Data[0].QuotaId))
	d.Set("disk", project.Data[0].Disk)
	d.Set("core_arch", project.EngineInfo.DefaultClusterArch)
	d.Set("external_table", project.EngineInfo.ExternalProjectCnt == 1)
	Properties, err := maxcomputeService.DescribeMaxProjectPropertiesForAscm(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	encryption := make(map[string]interface{})
	err = json.Unmarshal(Properties["ENCRYPTION"].([]byte), &encryption)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	vpcData := make([]string, 0)
	err = json.Unmarshal(Properties["odps.security.vpc.whitelist"].([]byte), &vpcData)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("encryption", encryption["ENCRYPTION_ENABLE"].(bool))
	d.Set("encrypt_algorithm", encryption["ENCRYPTION_ALGORITHM"].(string))
	d.Set("encryption", encryption["ENCRYPTION_KEY"].(string))
	vpc := strings.Split(vpcData[0], "_")
	if len(vpc) == 2 {
		d.Set("encryption", vpc[0])
	}
	return nil
}

func resourceAlibabacloudStackMaxcomputeProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if !d.IsNewResource() && d.HasChanges("encryption", "encrypt_algorithm", "encryption_key") {
		request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "UpdateOdpsProjectPropertiesForAscm", "")
		encrypt_algorithm := d.Get("encrypt_algorithm").(string)
		encryption_key := d.Get("encryption_key").(string)
		encryption := fmt.Sprintf("{\"ENCRYPTION_ENABLE\":\"true\",\"ENCRYPTION_ALGORITHM\":\"%s\",\"ENCRYPTION_KEY\":\"%s\"}", encrypt_algorithm, encryption_key)
		propertyMap := map[string]interface{}{"ENCRYPTION": encryption}
		propertyStr, _ := json.Marshal(propertyMap)
		mergeMaps(request.QueryParams, map[string]string{
			"engineId":    d.Get("name").(string),
			"ProjectName": d.Get("name").(string),
			"propertyMap": string(propertyStr),
		})
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			errmsg := ""
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "UpdateOdpsProjectPropertiesForAscm", errmsg)
		}
	}
	if d.HasChange("vpc_id") && d.Get("vpc_id") != "" {
		request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "UpdateOdpsProjectVpcForAscm", "")
		vpcIdList := fmt.Sprintf("[\"%s_%s\"]", client.RegionId, d.Get("vpc_id").(string))
		mergeMaps(request.QueryParams, map[string]string{
			"engineId":    d.Get("name").(string),
			"vpcIdList":   vpcIdList,
			"ProjectName": d.Get("name").(string),
			"RegionIds":   client.RegionId,
		})
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			errmsg := ""
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "UpdateOdpsProjectVpcForAscm", errmsg)
		}
	}
	return nil
}

func resourceAlibabacloudStackMaxcomputeProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "DeleteCalcEngineForAscm", "")
	request.QueryParams["EngineId"] = d.Id()
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "DeleteCalcEngineForAscm", errmsg)
	}
	return nil
}
