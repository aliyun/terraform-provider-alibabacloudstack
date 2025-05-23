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

func resourceAlibabacloudStackVpcIpv6EgressRule() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Ipv6Address"}, false),
			},
			"ipv6_egress_rule_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z\u4E00-\u9FA5][\u4E00-\u9FA5A-Za-z0-9_-]{2,128}$"), "The name must be `2` to `128` characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`."),
			},
			"ipv6_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackVpcIpv6EgressRuleCreate, resourceAlibabacloudStackVpcIpv6EgressRuleRead, nil, resourceAlibabacloudStackVpcIpv6EgressRuleDelete)
	return resource
}

func resourceAlibabacloudStackVpcIpv6EgressRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateIpv6EgressOnlyRule"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("ipv6_egress_rule_name"); ok {
		request["Name"] = v
	}
	request["Ipv6GatewayId"] = d.Get("ipv6_gateway_id")
	request["ClientToken"] = buildClientToken("CreateIpv6EgressOnlyRule")
	response, err := client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(request["Ipv6GatewayId"], ":", response["Ipv6EgressRuleId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcIpv6EgressRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackVpcIpv6EgressRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcIpv6EgressRule(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipv6_egress_rule vpcService.DescribeVpcIpv6EgressRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("ipv6_gateway_id", parts[0])
	d.Set("description", object["Description"])
	d.Set("instance_id", object["InstanceId"])
	d.Set("instance_type", object["InstanceType"])
	d.Set("ipv6_egress_rule_name", object["Name"])
	d.Set("status", object["Status"])
	return nil
}

func resourceAlibabacloudStackVpcIpv6EgressRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	vpcService := VpcService{client}
	action := "DeleteIpv6EgressOnlyRule"
	request := map[string]interface{}{
		"Ipv6EgressOnlyRuleId": parts[1],
	}

	request["ClientToken"] = buildClientToken("DeleteIpv6EgressOnlyRule")
	_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcIpv6EgressRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
