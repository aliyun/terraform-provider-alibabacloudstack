package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEcsHpcCluster() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackEcsHpcClusterCreate, resourceAlibabacloudStackEcsHpcClusterRead, resourceAlibabacloudStackEcsHpcClusterUpdate, resourceAlibabacloudStackEcsHpcClusterDelete)
	return resource
}

func resourceAlibabacloudStackEcsHpcClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateHpcCluster"
	var Description string
	if v, ok := d.GetOk("description"); ok {
		Description = fmt.Sprint(v.(string))
	}
	Name := d.Get("name").(string)
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	ClientToken := buildClientToken("CreateHpcCluster")

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
	request.QueryParams["Name"] = Name
	request.QueryParams["ClientToken"] = ClientToken
	request.QueryParams["Description"] = Description

	raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ecs_hpc_cluster", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(action, raw, request)

	resp := &ecs.CreateHpcClusterResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	d.SetId(fmt.Sprint(resp.HpcClusterId))
	return nil
}

func resourceAlibabacloudStackEcsHpcClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsHpcCluster(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ecs_hpc_cluster ecsService.DescribeEcsHpcCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("description", object.HpcClusters.HpcCluster[0].Description)
	d.Set("name", object.HpcClusters.HpcCluster[0].Name)
	return nil
}

func resourceAlibabacloudStackEcsHpcClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	Description := d.Get("description").(string)
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
		runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
		runtime.SetAutoretry(true)
		ClientToken := buildClientToken("ModifyHpcClusterAttribute")

		request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
		request.QueryParams["HpcClusterId"] = HpcClusterId
		request.QueryParams["ClientToken"] = ClientToken
		request.QueryParams["Description"] = Description
		request.QueryParams["Name"] = Name

		response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := response.(*responses.CommonResponse)
		addDebug(action, response, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	return nil
}

func resourceAlibabacloudStackEcsHpcClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteHpcCluster"
	HpcClusterId := d.Id()
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	ClientToken := buildClientToken("DeleteHpcCluster")

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
	request.QueryParams["HpcClusterId"] = HpcClusterId
	request.QueryParams["ClientToken"] = ClientToken

	response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := response.(*responses.CommonResponse)
	addDebug(action, response, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}