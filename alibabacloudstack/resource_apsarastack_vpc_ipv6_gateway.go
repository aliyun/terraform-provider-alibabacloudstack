package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackVpcIpv6Gateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackVpcIpv6GatewayCreate,
		Read:   resourceAlibabacloudStackVpcIpv6GatewayRead,
		Update: resourceAlibabacloudStackVpcIpv6GatewayUpdate,
		Delete: resourceAlibabacloudStackVpcIpv6GatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"ipv6_gateway_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z\u4E00-\u9FA5][\u4E00-\u9FA5A-Za-z0-9_-]{2,128}$"), "The name must be `2` to `128` characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`."),
			},
			"spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Large", "Medium", "Small"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackVpcIpv6GatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateIpv6Gateway"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("ipv6_gateway_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("spec"); ok {
		request["Spec"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	request["ClientToken"] = buildClientToken("CreateIpv6Gateway")
	response, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["Ipv6GatewayId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcIpv6GatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackVpcIpv6GatewayRead(d, meta)
}

func resourceAlibabacloudStackVpcIpv6GatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcIpv6Gateway(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_vpc_ipv6_gateway vpcService.DescribeVpcIpv6Gateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("ipv6_gateway_name", object["Name"])
	d.Set("spec", object["Spec"])
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["VpcId"])
	return nil
}

func resourceAlibabacloudStackVpcIpv6GatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"Ipv6GatewayId": fmt.Sprintf("[\"%s\"]", d.Id()),
	}
	if d.HasChange("spec") {
		update = true
	}
	if v, ok := d.GetOk("spec"); ok {
		request["Spec"] = v
	}

	if update {
		action := "ModifyIpv6GatewaySpec"
		request["ClientToken"] = buildClientToken("ModifyIpv6GatewaySpec")
		response, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
		addDebug(action, response, request)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpcIpv6GatewayStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	update = false
	modifyIpv6GatewayAttributeReq := map[string]interface{}{
		"Ipv6GatewayId": d.Id(),
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			modifyIpv6GatewayAttributeReq["Description"] = v
		}
	}
	if d.HasChange("ipv6_gateway_name") {
		update = true
		if v, ok := d.GetOk("ipv6_gateway_name"); ok {
			modifyIpv6GatewayAttributeReq["Name"] = v
		}
	}
	if update {
		action := "ModifyIpv6GatewayAttribute"
		response, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, modifyIpv6GatewayAttributeReq)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpcIpv6GatewayStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackVpcIpv6GatewayRead(d, meta)
}

func resourceAlibabacloudStackVpcIpv6GatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	action := "DeleteIpv6Gateway"
	request := map[string]interface{}{
		"Ipv6GatewayId": d.Id(),
	}
	_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcIpv6GatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
