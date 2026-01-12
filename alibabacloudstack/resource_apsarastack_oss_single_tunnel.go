package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOssSingleTunnel() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"one_router": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
			},
			"department": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"resource_group": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"share": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"department_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ascm_create_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rm_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOssSingleTunnelCreate, resourceAlibabacloudStackOssSingleTunnelRead, resourceAlibabacloudStackOssSingleTunnelUpdate, resourceAlibabacloudStackOssSingleTunnelDelete)
	return resource
}

func resourceAlibabacloudStackOssSingleTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"$oneRouter":    d.Get("one_router").(string),
		"Department":    d.Get("department").(int),
		"ResourceGroup": d.Get("resource_group").(int),
		"RegionId":      d.Get("region_id").(string),
		"Cluster":       d.Get("cluster").(string),
		"Share":         d.Get("share").(int),
		"Label":         d.Get("label").(string),
		"VpcId":         d.Get("vpc_id").(string),
		"VSwitchId":     d.Get("vswitch_id").(string),
		"_inner":        map[string]interface{}{},
	}

	// Remove empty values from request
	for k, v := range reqQuery {
		if v == "" || v == 0 {
			delete(reqQuery, k)
		}
	}

	if _, err := client.DoTeaRequest("POST", "oss", "2019-09-01", "CreateVpcip", "", nil, reqQuery, nil); err != nil {
		return err
	}

	// Get the VPC IP by listing and matching the parameters
	listReq := map[string]interface{}{}
	response, err := client.DoTeaRequest("POST", "oss", "2019-09-01", "ListVpcip", "", nil, listReq, nil)
	if err != nil {
		return err
	}

	var vpcIp string
	if listResult, ok := response["ListVpcipResult"].(map[string]interface{}); ok {
		if vpcips, ok := listResult["Vpcip"].([]interface{}); ok {
			for _, vpcipItem := range vpcips {
				if item, ok := vpcipItem.(map[string]interface{}); ok {
					// Match based on the input parameters to find the created VPC IP
					if item["VpcId"] == d.Get("vpc_id") &&
						item["Label"] == d.Get("label") &&
						item["Cluster"] == d.Get("cluster") &&
						item["RegionId"] == d.Get("region_id") {
						vpcIp = item["Vip"].(string)
						break
					}
				}
			}
		}
	}

	if vpcIp == "" {
		return fmt.Errorf("failed to get VPC IP after creation")
	}

	// Set the ID as the VPC IP address according to the resource ID rule
	d.SetId(vpcIp)

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

	d.Set("shared", formatInt(object["shared"]))
	d.Set("department", formatInt(object["Department"]))
	d.Set("department_name", object["DepartmentName"])
	d.Set("region", object["Region"])
	d.Set("cluster", object["Cluster"])
	d.Set("region_id", object["RegionId"])
	d.Set("vip", object["Vip"])
	d.Set("ascm_create_user", object["AscmCreateUser"])
	d.Set("resource_group", formatInt(object["ResourceGroup"]))
	d.Set("resource_group_name", object["ResourceGroupName"])
	d.Set("rm_region_id", object["RMRegionId"])
	d.Set("label", object["Label"])
	d.Set("vpc_id", object["VpcId"])

	return nil
}

func resourceAlibabacloudStackOssSingleTunnelUpdate(d *schema.ResourceData, meta interface{}) error {
	// This resource does not support modification according to the API documentation
	return nil
}

func resourceAlibabacloudStackOssSingleTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	vip := d.Get("vip").(string)
	vpcId := d.Get("vpc_id").(string)
	regionId := d.Get("region_id").(string)

	if vip == "" || vpcId == "" || regionId == "" {
		return nil
	}

	content := fmt.Sprintf("<DeleteVpcip><Region>%s</Region><VpcId>%s</VpcId><Vip>%s</Vip></DeleteVpcip>", regionId, vpcId, vip)

	reqQuery := map[string]interface{}{
		"content": content,
		"params": map[string]interface{}{
			"Vip": vip,
		},
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "Oss", "2019-09-01", "DeleteVpcip", "", nil, reqQuery, nil)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteVpcip", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
