package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOssSingleTunnel() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"shared": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOssSingleTunnelCreate, resourceAlibabacloudStackOssSingleTunnelRead, nil, resourceAlibabacloudStackOssSingleTunnelDelete)
	return resource
}

func resourceAlibabacloudStackOssSingleTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "CreateVpcip"
	request.QueryParams["ProductName"] = "oss"
	reqQuery := map[string]interface{}{
		"Department":    client.Department,
		"ResourceGroup": client.ResourceGroup,
		"RegionId":      client.RegionId,
		"Cluster":       d.Get("cluster").(string),
		"Share":         d.Get("shared").(string),
		"Label":         d.Get("label").(string),
		"VpcId":         d.Get("vpc_id").(string),
		"VSwitchId":     d.Get("vswitch_id").(string),
	}
	if querybytes, err := json.Marshal(reqQuery); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "json Marshal", "CreateVpcip", errmsgs.AlibabacloudStackOssGoSdk)
	} else {
		request.QueryParams["Params"] = string(querybytes)
	}
	context := fmt.Sprintf("<CreateVpcip><Region>%s</Region><VSwitchId>%s</VSwitchId><Label>%s</Label><Cluster>%s</Cluster></CreateVpcip>", client.Region, reqQuery["VSwitchId"], reqQuery["Label"], reqQuery["Cluster"])
	request.QueryParams["Content"] = context
	response := make(map[string]interface{})
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("CreateVpcip", bresponse, request, request.QueryParams)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "json Unmarshal", "CreateVpcip", errmsgs.AlibabacloudStackOssGoSdk)
	}
	vpcIp, err := jsonpath.Get("$.Data.CreateVpcipResult.Vip", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, d.Id(), "$.Data.CreateVpcipResult.Vip", response)
	}

	// Set the ID as the VPC IP address according to the resource ID rule
	d.SetId(fmt.Sprintf("%s:%s:%s", d.Get("cluster").(string), d.Get("vpc_id").(string), vpcIp))

	return nil
}

func resourceAlibabacloudStackOssSingleTunnelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}

	object, err := ossService.DescribeOssSingleTunnel(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_oss_single_tunnel ossService.DescribeOssSingleTunnel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if object["shared"] != nil {
		d.Set("shared", fmt.Sprint(object["shared"]))
	}
	d.Set("label", object["Label"])
	d.Set("cluster", object["Cluster"])
	d.Set("vip", object["Vip"])
	d.Set("vpc_id", object["VpcId"])

	return nil
}

func resourceAlibabacloudStackOssSingleTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts := strings.Split(d.Id(), ":")

	content := fmt.Sprintf("<DeleteVpcip><Region>%s</Region><VpcId>%s</VpcId><Vip>%s</Vip></DeleteVpcip>", client.Region, parts[1], parts[2])

	request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	mergeMaps(request.QueryParams, map[string]string{
		"OpenApiAction": "DeleteVpcip",
		"ProductName":   "oss",
		"Content":       content,
		"Params":        "{\"Vip\":\"" + parts[2] + "\"}",
	})
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("DeleteVpcip", bresponse, request, bresponse.GetHttpContentString())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
