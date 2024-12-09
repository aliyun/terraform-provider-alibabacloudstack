package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackMaxcomputeProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackMaxcomputeProjectCreate,
		Read:   resourceAlibabacloudStackMaxcomputeProjectRead,
		Update: resourceAlibabacloudStackMaxcomputeProjectUpdate,
		Delete: resourceAlibabacloudStackMaxcomputeProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlibabacloudStackMaxcomputeProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	roleId, err := client.RoleIds()
	if err != nil {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ASCM User", "defaultRoleId")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return err
	}

	cluster_name := d.Get("cluster").(string)
	clusters, err := DescribeMaxcomputeProject(meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var cluster map[string]interface{}
	for _, object := range clusters {
		cluster = object.(map[string]interface{})
		if cluster["cluster"].(string) == cluster_name {
			break
		}
	}
	if cluster == nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_cluster", "getCluster", cluster_name)
	}

	disk_size := d.Get("disk").(int)
	name := d.Get("name").(string)
	pk := d.Get("pk").(string)

	request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "CreateCalcEngineForAscm", "")
	mergeMaps(request.QueryParams, map[string]string{
		"KmsRegion":       string(client.Region),
		"ResourceGroupId": client.ResourceGroup,
		"Product":         "dataworks-private-cloud",
		"CalcEngineType":  "ODPS",
		"OrganizationId":  client.Department,
		"EnvType":         "PRD",
		"Name":            name,
		"EngineInfo":      "{\"taskAk\":{\"kp\":\"" + pk + "\",\"aliyunAccount\":\"ascm-dw-1637809230710\"},\"clusters\":[{\"name\":\"" + cluster_name + "\",\"quota\":" + d.Get("quota_id").(string) + ",\"disk\":" + fmt.Sprintf("%f", float64(disk_size)/1024) + ",\"isDefault\":1,\"projectQuota\":{\"fileLength\":" + strconv.Itoa(disk_size*1024*1024*1024) + ",\"fileNumber\":null}}],\"odpsProjectName\":\"" + name + "\",\"needToCreateOdpsProject\":true,\"defaultClusterArch\":\"" + cluster["core_arch"].(string) + "\",\"isOdpsDev\":false}",
		"Department":      client.Department,
		"Version":         "2019-01-17",
		"ClusterItem":     "{\"cluster\":\"" + cluster_name + "\",\"core_arch\":\"" + cluster["core_arch"].(string) + "\",\"project\":\"" + cluster["project"].(string) + "\",\"region\":\"" + cluster["region"].(string) + "\"}",
		"ClusterName":     cluster_name,
		"ResourceGroup":   client.ResourceGroup,
		"ExternalTable":   strconv.FormatBool(d.Get("external_table").(bool)),
		"TaskPk":          pk,
		"OdpsName":        name,
		"RegionId":        client.RegionId,
		"CurrentRoleId":   strconv.Itoa(roleId),
	})

	if v, ok := d.GetOk("enabled_mc_encrypt"); ok && v.(bool) {
		request.QueryParams["EnabledMcEncrypt"] = "1"
		if _, ok := d.GetOk("mc_encrypt_algorithm"); !ok {
			log.Printf("mc_encrypt_algorithm not set while enable me encrypt")
			return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_maxcompute_project", "mc_encrypt_algorithm", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		request.QueryParams["McEncryptAlgorithm"] = d.Get("mc_encrypt_algorithm").(string)
		if _, ok := d.GetOk("mc_encrypt_key"); !ok {
			log.Printf("mc_encrypt_key not set while enable me encrypt")
			return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_maxcompute_project", "mc_encrypt_algorithm", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		request.QueryParams["McEncryptKey"] = d.Get("mc_encrypt_key").(string)
	}

	if v, ok := d.GetOk("vpc_tunnel_ids"); ok {
		vpc_tunnel_ids := ""
		for _, id := range v.([]interface{}) {
			vpc_tunnel_ids += id.(string)
		}
		request.QueryParams["McEncryptKey"] = vpc_tunnel_ids
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw create maxcomputecluster is : %s", raw)

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "Create", errmsg)
	}

	addDebug("MaxcomputeProjectCreate", raw, request)

	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "Create", errmsg)
	}
	addDebug("MaxcomputeProjectCreate", raw, request, bresponse.GetHttpContentString())

	return resourceAlibabacloudStackMaxcomputeProjectRead(d, meta)
}

func resourceAlibabacloudStackMaxcomputeProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	objects, err := maxcomputeService.DescribeMaxcomputeProject(d.Get("name").(string))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project_user maxcomputeService.DescribeMaxcomputeUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	project := objects.Data.CalcEngines[0]
	d.SetId(strconv.Itoa(project.EngineId))
	return nil
}

func resourceAlibabacloudStackMaxcomputeProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChanges("cluster", "disk") {
		roleId, err := client.RoleIds()
		if err != nil {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ASCM User", "defaultRoleId")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateOdpsQuota", "")
		request.Headers["x-acs-roleid"] = strconv.Itoa(roleId)
		mergeMaps(request.QueryParams, map[string]string{
			"Cluster":           d.Get("cluster").(string),
			"Product":           "ascm",
			"Cu":                d.Get("quota_id").(string),
			"Format":            "JSON",
			"Forwardedregionid": client.RegionId,
			"Version":           "2019-05-10",
			"RegionId":          client.RegionId,
			"Id":                d.Get("id").(string),
			"Disk":              fmt.Sprintf("%f", float64(d.Get("disk").(int))/1024),
		})

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"ErrorOdpsQuota Not Found"}) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Get("user_name").(string), "UpdateOdpsQuota", errmsg)
		}
		addDebug("UpdateOdpsQuota", raw, request)
	}

	return resourceAlibabacloudStackMaxcomputeProjectRead(d, meta)
}

func resourceAlibabacloudStackMaxcomputeProjectDelete(d *schema.ResourceData, meta interface{}) error {
	// 不支持删除
	return nil
}
