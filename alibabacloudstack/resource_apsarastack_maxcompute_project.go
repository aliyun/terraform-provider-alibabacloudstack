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
)

func resourceAlibabacloudStackMaxcomputeProject() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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
			// "encryption": {
			// 	Type:     schema.TypeBool,
			// 	Default:  false,
			// 	Optional: true,
			// },
			// "encrypt_algorithm": {
			// 	Type:         schema.TypeString,
			// 	Optional:     true,
			// 	Computed:     true,
			// 	ValidateFunc: validation.StringInSlice([]string{"SM4", "RC4", "AES256", "AESCTR"}, false),
			// },
			// "encryption_key": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			"vpc_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	disk_size := float64(d.Get("disk").(int)) / 1024.0
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
				"disk":      disk_size,
				"projectQuota": map[string]interface{}{
					"fileLength": d.Get("disk").(int) * 1024 * 1024 * 1024,
					"fileNumber": nil,
				},
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
		"ResourceGroup":  client.ResourceGroup,
		"CpuType":        d.Get("cpu_type").(string),
		"Clusters":       string(clusterStr),
		"ExternalTable":  fmt.Sprint(d.Get("external_table")),
		"OrganizationId": client.Department,
		"TaskPk":         d.Get("account_pk").(string),
		"CalcEngineType": "ODPS",
		"EnvType":        "PRD",
		"Name":           d.Get("name").(string),
		"RegionId":       client.RegionId,
		"EngineInfo":     string(engineInfoStr),
	}

	// if d.Get("encryption").(bool) {
	// 	body["EnabledMcEncrypt"] = "1"
	// 	body["McEncryptAlgorithm"] = d.Get("encrypt_algorithm").(string)
	// 	body["McEncryptKey"] = d.Get("encryption_key").(string)
	// }
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
	// d.SetId("21")
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
	d.Set("core_arch", project.EngineInfo.DefaultClusterArch)
	quota, err := maxcomputeService.ListOdpsEngineQuotaForAscm(d.Id(), project.Name)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("quota_id", fmt.Sprint(quota.QuotaId))
	d.Set("disk", int(quota.Disk*1024))
	account, err := maxcomputeService.DescribeMaxcomputeProjectEngine(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("account", account.EngineInfo.TaskAk.AliyunAccount)
	d.Set("account_pk", fmt.Sprint(account.EngineInfo.TaskAk.Kp))
	Properties, err := maxcomputeService.DescribeMaxProjectPropertiesForAscm(d.Id(), project.Name)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	// encryption := make(map[string]interface{})
	// err = json.Unmarshal([]byte(Properties["ENCRYPTION"].(string)), &encryption)
	// if err != nil {
	// 	return errmsgs.WrapError(err)
	// }
	vpcData := make([]string, 0)
	err = json.Unmarshal([]byte(Properties["odps.security.vpc.whitelist"].(string)), &vpcData)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	// d.Set("encryption", encryption["ENCRYPTION_ENABLE"].(string) == "true")
	// d.Set("encrypt_algorithm", encryption["ENCRYPTION_ALGORITHM"].(string))
	// d.Set("encryption_key", encryption["ENCRYPTION_KEY"].(string))
	vpcs := make([]string, 0)
	for _, v := range vpcData {
		vpc := strings.Split(v, "_")
		vpcs = append(vpcs, vpc[1])
	}
	d.Set("vpc_ids", vpcs)
	return nil
}

func resourceAlibabacloudStackMaxcomputeProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	// if !d.IsNewResource() && d.HasChanges("encryption", "encrypt_algorithm", "encryption_key") {
	// 	request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "UpdateOdpsProjectPropertiesForAscm", "")
	// 	encrypt_algorithm := d.Get("encrypt_algorithm").(string)
	// 	encryption_key := d.Get("encryption_key").(string)
	// 	encryption := fmt.Sprintf("{\"ENCRYPTION_ENABLE\":\"true\",\"ENCRYPTION_ALGORITHM\":\"%s\",\"ENCRYPTION_KEY\":\"%s\"}", encrypt_algorithm, encryption_key)
	// 	propertyMap := map[string]interface{}{"ENCRYPTION": encryption}
	// 	propertyStr, _ := json.Marshal(propertyMap)
	// 	mergeMaps(request.QueryParams, map[string]string{
	// 		"Region":      client.RegionId,
	// 		"Action":      "UpdateOdpsProjectPropertiesForAscm",
	// 		"AccessKeyId": client.AccessKey,
	// 		"EngineId":    d.Get("name").(string),
	// 		"ProjectName": d.Get("name").(string),
	// 		"PropertyMap": string(propertyStr),
	// 	})
	// 	bresponse, err := client.ProcessCommonRequest(request)
	// 	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	// 	if err != nil {
	// 		errmsg := ""
	// 		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
	// 		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "UpdateOdpsProjectPropertiesForAscm", errmsg)
	// 	}
	// }
	if d.HasChange("vpc_ids") && !(d.IsNewResource() && len(d.Get("vpc_ids").([]interface{})) == 0) {
		vpc := make([]string, 0)
		for _, vpc_id := range d.Get("vpc_ids").([]interface{}) {
			vpcitem := fmt.Sprintf("%s_%s", client.RegionId, vpc_id.(string))
			vpc = append(vpc, vpcitem)
		}
		vpcData, _ := json.Marshal(vpc)
		request := map[string]interface{}{
			"Region":      client.RegionId,
			"Action":      "UpdateOdpsProjectVpcForAscm",
			"AccessKeyId": client.AccessKey,
			"EngineId":    d.Get("name").(string),
			"VpcIdList":   string(vpcData),
			"ProjectName": d.Get("name").(string),
			"RegionIds":   client.RegionId,
		}
		response, err := client.DoTeaRequest("POST", "dataworks-private-cloud", "2019-01-17", "UpdateOdpsProjectVpcForAscm", "", nil, request, nil)
		if err != nil {
			errmsg := ""
			errmsg = errmsgs.GetAsapiErrorMessage(response)
			if errmsg != "" {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "UpdateOdpsProjectVpcForAscm", errmsg)
			}
		}
	}
	if !d.IsNewResource() && d.HasChange("disk") {
		return errmsgs.Error("[ERROR] The disk attribute cannot be modified after the project is created.")
	}
	return nil
}

func resourceAlibabacloudStackMaxcomputeProjectDelete(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*connectivity.AlibabacloudStackClient)
	// request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "DeleteCalcEngineForAscm", "")
	// request.QueryParams["EngineId"] = d.Id()
	// bresponse, err := client.ProcessCommonRequest(request)
	// addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	// if err != nil {
	// 	errmsg := ""
	// 	errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
	// 	return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "DeleteCalcEngineForAscm", errmsg)
	// }
	return nil
}
