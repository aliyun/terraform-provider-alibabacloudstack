package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
			"vpc_tunnel_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_table": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			/*			"enabled_mc_encrypt": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"mc_encrypt_algorithm": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"SM4", "RC4", "AES256", "AESCTR"}, false),
						},
						"mc_encrypt_key": {
							Type:     schema.TypeString,
							Optional: true,
						},*/
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
			"pk": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aliyun_account": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackMaxcomputeProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}

	roleId, err := client.RoleIds()
	if err != nil {
		err = WrapErrorf(Error(GetNotFoundMessage("ASCM User", "defaultRoleId")), NotFoundMsg, ProviderERROR)
		return err
	}

	cluster_name := d.Get("cluster").(string)
	clusters, err := DescribeMaxcomputeProject(meta)
	if err != nil {
		return WrapError(err)
	}

	var cluster map[string]interface{}
	for _, object := range clusters {
		cluster = object.(map[string]interface{})
		if cluster["cluster"].(string) == cluster_name {
			break
		}
	}
	if cluster == nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_maxcompute_cluster", "getCluster", cluster_name)
	}

	disk_size := d.Get("disk").(int)

	name := d.Get("name").(string)

	pk := d.Get("pk").(string)

	request.QueryParams = map[string]string{
		"Action":          "CreateCalcEngineForAscm",
		"KmsRegion":       string(client.Region),
		"ResourceGroupId": client.ResourceGroup,
		"Product":         "dataworks-private-cloud",
		"CalcEngineType":  "ODPS", // 固定值
		"OrganizationId":  client.Department,
		//"Timestamp": "2021-11-25T10:14:13Z",
		"EnvType":    "PRD", // 固定值
		"Name":       name,
		"EngineInfo": "{\"taskAk\":{\"kp\":\"" + pk + "\",\"aliyunAccount\":\"ascm-dw-1637809230710\"},\"clusters\":[{\"name\":\"" + cluster_name + "\",\"quota\":" + d.Get("quota_id").(string) + ",\"disk\":" + fmt.Sprintf("%f", float64(disk_size)/1024) + ",\"isDefault\":1,\"projectQuota\":{\"fileLength\":" + strconv.Itoa(disk_size*1024*1024*1024) + ",\"fileNumber\":null}}],\"odpsProjectName\":\"" + name + "\",\"needToCreateOdpsProject\":true,\"defaultClusterArch\":\"" + cluster["core_arch"].(string) + "\",\"isOdpsDev\":false}",
		"Department": client.Department,
		//"Format": "JSON",
		//"XRealIp": "10.30.208.219",
		//"Language": "zh",
		//"VpcTunnelIdList": "cn-neimeng-env30-d01_vpc-3rqm203vt5n0nv1gv52fr,cn-neimeng-env30-d01_vpc-3rq45siuthzs2da1ras0k",
		"Version":     "2019-01-17",
		"ClusterItem": "{\"cluster\":\"" + cluster_name + "\",\"core_arch\":\"" + cluster["core_arch"].(string) + "\",\"project\":\"" + cluster["project"].(string) + "\",\"region\":\"" + cluster["region"].(string) + "\"}",
		//"McEncryptAlgorithm": "RC4",
		//"McEncryptKey":       "DEFAULT",
		"ClusterName":   cluster_name,
		"ResourceGroup": client.ResourceGroup,
		"ExternalTable": strconv.FormatBool(d.Get("external_table").(bool)),
		//"HaAlibabacloudStack": "false", //不需要
		"TaskPk":   pk,
		"OdpsName": name,
		//"CsrfToken": "jnzqHxYi-J7EDkvbD4tNwcpp19qoQV-K9B30",
		//"SignatureVersion": "2.1",
		//"EnabledMcEncrypt": "1",
		//"SignatureNonce": "5888fddc5ffaf5f62551544515ecde6e",
		//"WlProxyClientIp": "10.30.3.1",
		//"AccessKeyId": "O4TaXipwrLw6LPIV",
		//"XForwardedFor": "10.30.3.1, 10.30.208.219",
		//"SignatureMethod": "HMAC-SHA1",
		"RegionId":      client.RegionId,
		"CurrentRoleId": strconv.Itoa(roleId),
		//"XForwardedIp": "10.30.3.1"
	}

	if v, ok := d.GetOk("enabled_mc_encrypt"); ok && v.(bool) {
		request.QueryParams["EnabledMcEncrypt"] = "1"
		if _, ok := d.GetOk("mc_encrypt_algorithm"); !ok {
			log.Printf("mc_encrypt_algorithm not set while enable me encrypt")

			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_maxcompute_project", "mc_encrypt_algorithm", AlibabacloudStackSdkGoERROR)

		}
		request.QueryParams["McEncryptAlgorithm"] = d.Get("mc_encrypt_algorithm").(string)
		if _, ok := d.GetOk("mc_encrypt_key"); !ok {
			log.Printf("mc_encrypt_key not set while enable me encrypt")

			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_maxcompute_project", "mc_encrypt_algorithm", AlibabacloudStackSdkGoERROR)

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

	request.Method = "POST"
	request.Product = "dataworks-private-cloud"
	request.Version = "2019-01-17"
	request.ServiceCode = "dataworks-private-cloud"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "CreateCalcEngineForAscm"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw create maxcomputecluster is : %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_maxcompute_project", "Create", raw)
	}

	addDebug("MaxcomputeProjectCreate", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_maxcompute_project", "Create", AlibabacloudStackSdkGoERROR)
	}
	addDebug("MaxcomputeProjectCreate", raw, requestInfo, bresponse.GetHttpContentString())

	return resourceAlibabacloudStackMaxcomputeProjectRead(d, meta)
}
func resourceAlibabacloudStackMaxcomputeProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	objects, err := maxcomputeService.DescribeMaxcomputeProject(d.Get("name").(string))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project_user maxcomputeService.DescribeMaxcomputeUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	project := objects.Data.CalcEngines[0]
	d.SetId(strconv.Itoa(project.EngineId))
	return nil
}

func resourceAlibabacloudStackMaxcomputeProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChange("cluster") || d.HasChange("disk") {
		var requestInfo *ecs.Client
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		roleId, err := client.RoleIds()
		if err != nil {
			err = WrapErrorf(Error(GetNotFoundMessage("ASCM User", "defaultRoleId")), NotFoundMsg, ProviderERROR)
			return WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "UpdateOdpsQuota"
		request.Headers = map[string]string{
			"RegionId":              client.RegionId,
			"x-acs-roleid":          strconv.Itoa(roleId),
			"x-acs-resourcegroupid": client.ResourceGroup,
			"x-acs-regionid":        client.RegionId,
			"x-acs-organizationid":  client.Department,
		}
		request.QueryParams = map[string]string{
			"Action":            "UpdateOdpsQuota",
			"Cluster":           d.Get("cluster").(string),
			"Product":           "ascm",
			"Cu":                d.Get("quota_id").(string),
			"Format":            "JSON",
			"Forwardedregionid": client.RegionId,
			"Version":           "2019-05-10",
			"RegionId":          client.RegionId,
			"Id":                d.Get("id").(string),
			"Disk":              fmt.Sprintf("%f", float64(d.Get("disk").(int))/1024),
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ErrorOdpsQuota Not Found"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Get("user_name").(string), "UpdateOdpsQuota", AlibabacloudStackSdkGoERROR)

		}
		addDebug("UpdateOdpsQuota", raw, requestInfo, request)
	}

	return resourceAlibabacloudStackMaxcomputeProjectRead(d, meta)
}

func resourceAlibabacloudStackMaxcomputeProjectDelete(d *schema.ResourceData, meta interface{}) error {
	//不支持删除
	return nil
}
