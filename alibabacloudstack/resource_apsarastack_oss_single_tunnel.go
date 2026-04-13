package alibabacloudstack

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"

	oss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
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
	ossService := OssService{client}
	cluster := d.Get("cluster").(string)
	ossClient, err := ossService.GetOssClientForCluster(cluster)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_single_tunnel", "GetOssClientForCluster", errmsgs.AlibabacloudStackOssGoSdk)
	}
	label := d.Get("label").(string)
	vswitchId := d.Get("vswitch_id").(string)
	vpcId := d.Get("vpc_id").(string)
	shared := d.Get("shared").(string)

	xmlBody := fmt.Sprintf(
		"<CreateVpcip><Region>%s</Region><VSwitchId>%s</VSwitchId><Label>%s</Label><Cluster>%s</Cluster><VpcId>%s</VpcId><Share>%s</Share></CreateVpcip>",
		client.RegionId, vswitchId, label, cluster, vpcId, shared)

	input := &oss.OperationInput{
		OpName:     "CreateVpcip",
		Method:     "PUT",
		Parameters: map[string]string{"vpcip": ""},
		Headers:    map[string]string{"Content-Type": "application/xml"},
		Body:       strings.NewReader(xmlBody),
	}
	output, err := ossClient.InvokeOperation(context.Background(), input)
	addDebug("CreateVpcip", output, input, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "CreateVpcip", "InvokeOperation", errmsgs.AlibabacloudStackOssGoSdk)
	}
	defer output.Body.Close()
	body, err := io.ReadAll(output.Body)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "CreateVpcip", "ReadBody", errmsgs.AlibabacloudStackOssGoSdk)
	}

	// Parse XML response to get Vip
	var createResult struct {
		Vip string `xml:"Vip"`
	}
	if err = xml.Unmarshal(body, &createResult); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "CreateVpcip", "XMLUnmarshal", errmsgs.AlibabacloudStackOssGoSdk)
	}
	if createResult.Vip == "" {
		return errmsgs.WrapErrorf(fmt.Errorf("empty Vip in response"), errmsgs.DefaultErrorMsg, "CreateVpcip", "ParseVip", errmsgs.AlibabacloudStackOssGoSdk)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", cluster, vpcId, createResult.Vip))

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
	ossService := OssService{client}
	parts := strings.Split(d.Id(), ":")

	ossClient, err := ossService.GetOssClientForCluster(parts[0])
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_single_tunnel", "GetOssClientForCluster", errmsgs.AlibabacloudStackOssGoSdk)
	}

	xmlBody := fmt.Sprintf(
		"<DeleteVpcip><Region>%s</Region><VpcId>%s</VpcId><Vip>%s</Vip></DeleteVpcip>",
		client.RegionId, parts[1], parts[2])

	input := &oss.OperationInput{
		OpName:     "DeleteVpcip",
		Method:     "DELETE",
		Parameters: map[string]string{"vpcip": ""},
		Headers:    map[string]string{"Content-Type": "application/xml"},
		Body:       strings.NewReader(xmlBody),
	}
	output, err := ossClient.InvokeOperation(context.Background(), input)
	addDebug("DeleteVpcip", output, input, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
