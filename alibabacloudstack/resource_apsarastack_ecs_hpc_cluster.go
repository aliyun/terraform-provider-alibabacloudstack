package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"log"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEcsHpcCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEcsHpcClusterCreate,
		Read:   resourceAlibabacloudStackEcsHpcClusterRead,
		Update: resourceAlibabacloudStackEcsHpcClusterUpdate,
		Delete: resourceAlibabacloudStackEcsHpcClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlibabacloudStackEcsHpcClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	//var response map[string]interface{}
	action := "CreateHpcCluster"
	//request := make(map[string]interface{})
	//conn, err := client.NewEcsClient()
	//if err != nil {
	//	return WrapError(err)
	//}
	var Description string
	if v, ok := d.GetOk("description"); ok {
		Description = fmt.Sprint(v.(string))
	}
	Name := d.Get("name").(string)
	//RegionId := client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	ClientToken := buildClientToken("CreateHpcCluster")
	request := requests.NewCommonRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Method = "POST"
	request.Product = "Ecs"
	request.Domain = client.Domain
	request.Version = "2014-05-26"
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "Ecs",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"Version":         "2014-05-26",
		"Name":            Name,
		"ClientToken":     ClientToken,
		"Description":     Description,
	}

	//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ecs_hpc_cluster", action, AlibabacloudStackSdkGoERROR)
	}
	addDebug(action, raw, request)
	//bresponse := raw.(*responses.CommonResponse)
	//err = json.Unmarshal(bresponse.GetHttpContentBytes(), EcsCreate)
	resp := &ecs.CreateHpcClusterResponse{}
	bresponse := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	//var response *ecs.CreateHpcClusterResponse
	//response, _ = raw.(*ecs.CreateHpcClusterResponse)
	d.SetId(fmt.Sprint(resp.HpcClusterId))
	return resourceAlibabacloudStackEcsHpcClusterRead(d, meta)
}
func resourceAlibabacloudStackEcsHpcClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsHpcCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ecs_hpc_cluster ecsService.DescribeEcsHpcCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object.HpcClusters.HpcCluster[0].Description)
	d.Set("name", object.HpcClusters.HpcCluster[0].Name)
	return nil
}
func resourceAlibabacloudStackEcsHpcClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	//var response map[string]interface{}
	update := false
	//request := map[string]interface{}{
	//	"HpcClusterId": d.Id(),
	//}
	request := requests.NewCommonRequest()
	//request["RegionId"] = client.RegionId
	Description := d.Get("description").(string)

	//var Description string
	//if d.HasChange("description") {
	//	update = true
	//	if v, ok := d.GetOk("description"); ok {
	//		Description = fmt.Sprint(v.(string))
	//	}
	//}
	var Name string
	if d.HasChange("name") {
		update = true
		if v, ok := d.GetOk("name"); ok {
			Name = fmt.Sprint(v.(string))
		}
	}
	HpcClusterId := d.Id()
	if update {
		action := "ModifyHpcClusterAttribute"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		ClientToken := buildClientToken("ModifyHpcClusterAttribute")
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Method = "POST"
		request.Product = "Ecs"
		request.Domain = client.Domain
		request.Version = "2014-05-26"
		request.ApiName = action
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "Ecs",
			"RegionId":        client.RegionId,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Action":          action,
			"Version":         "2014-05-26",
			"ClientToken":     ClientToken,
			"HpcClusterId":    HpcClusterId,
			"Description":     Description,
			"Name":            Name,
		}
		response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}
	}
	return resourceAlibabacloudStackEcsHpcClusterRead(d, meta)
}
func resourceAlibabacloudStackEcsHpcClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteHpcCluster"

	request := requests.NewCommonRequest()
	HpcClusterId := d.Id()
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	ClientToken := buildClientToken("DeleteHpcCluster")
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Method = "POST"
	request.Product = "Ecs"
	request.Domain = client.Domain
	request.Version = "2014-05-26"
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "Ecs",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"Version":         "2014-05-26",
		"ClientToken":     ClientToken,
		"HpcClusterId":    HpcClusterId,
	}
	response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
