package alibabacloudstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSecurityGroupRule() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
				Description:  "Type of rule, ingress (inbound) or egress (outbound).",
			},

			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "gre", "all"}, false),
			},

			"nic_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},

			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      GroupRulePolicyAccept,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
			},

			"port_range": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//Default:          AllPortRange,//port_range has been set to Required in Alibabacloudstack and Default cannot be used with Required
				DiffSuppressFunc: ecsSecurityGroupRulePortRangeDiffSuppressFunc,
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"cidr_ip", "ipv6_cidr_ip", "source_security_group_id"},
			},
			"ipv6_cidr_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ip"},
			},

			"source_security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ip"},
			},

			"source_group_owner_account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackSecurityGroupRuleCreate, resourceAlibabacloudStackSecurityGroupRuleRead, resourceAlibabacloudStackSecurityGroupRuleUpdate, deleteSecurityGroupRule)
	return resource
}

func resourceAlibabacloudStackSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	direction := d.Get("type").(string)
	sgId := d.Get("security_group_id").(string)
	ptl := d.Get("ip_protocol").(string)
	port := d.Get("port_range").(string)
	if port == "" {
		return errmsgs.WrapError(fmt.Errorf("'port_range': required field is not set or invalid."))
	}
	nicType := d.Get("nic_type").(string)
	policy := d.Get("policy").(string)
	priority := d.Get("priority").(int)

	if _, ok := d.GetOk("cidr_ip"); !ok {
		if _, ok := d.GetOk("source_security_group_id"); !ok {
			return errmsgs.WrapError(fmt.Errorf("Either 'cidr_ip' or 'source_security_group_id' must be specified."))
		}
	}
	request, err := buildAlibabacloudStackSGRuleRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var cidr_ip string
	if ip, ok := d.GetOk("cidr_ip"); ok {
		cidr_ip = ip.(string)
	} else {
		cidr_ip = d.Get("source_security_group_id").(string)
	}
	if direction == string(DirectionIngress) {
		request.ApiName = "AuthorizeSecurityGroup"
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_security_group_rule", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	} else {
		request.ApiName = "AuthorizeSecurityGroupEgress"
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_security_group_rule", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	d.SetId(sgId + ":" + direction + ":" + ptl + ":" + port + ":" + nicType + ":" + cidr_ip + ":" + policy + ":" + strconv.Itoa(priority))

	return nil
}

func resourceAlibabacloudStackSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	parts := strings.Split(d.Id(), ":")
	policy := parseSecurityRuleId(d, meta, 6)
	strPriority := parseSecurityRuleId(d, meta, 7)
	var priority int
	if policy == "" || strPriority == "" {
		policy = d.Get("policy").(string)
		priority = d.Get("priority").(int)
		d.SetId(d.Id() + ":" + policy + ":" + strconv.Itoa(priority))
	} else {
		prior, err := strconv.Atoi(strPriority)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		priority = prior
	}
	sgId := parts[0]
	direction := parts[1]

	object, err := ecsService.DescribeSecurityGroupRule(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			//return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("type", object.Direction)
	d.Set("ip_protocol", strings.ToLower(string(object.IpProtocol)))
	d.Set("nic_type", object.NicType)
	d.Set("policy", strings.ToLower(string(object.Policy)))
	d.Set("port_range", object.PortRange)
	d.Set("description", object.Description)
	if pri, err := strconv.Atoi(object.Priority); err != nil {
		return errmsgs.WrapError(err)
	} else {
		d.Set("priority", pri)
	}
	d.Set("security_group_id", sgId)
	//support source and desc by type
	if direction == string(DirectionIngress) {
		d.Set("cidr_ip", object.SourceCidrIp)
		d.Set("source_security_group_id", object.SourceGroupId)
		d.Set("source_group_owner_account", object.SourceGroupOwnerAccount)
	} else {
		d.Set("cidr_ip", object.DestCidrIp)
		d.Set("source_security_group_id", object.DestGroupId)
		d.Set("source_group_owner_account", object.DestGroupOwnerAccount)
	}
	return nil
}

func resourceAlibabacloudStackSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	policy := parseSecurityRuleId(d, meta, 6)
	strPriority := parseSecurityRuleId(d, meta, 7)
	var priority int
	if policy == "" || strPriority == "" {
		policy = d.Get("policy").(string)
		priority = d.Get("priority").(int)
		d.SetId(d.Id() + ":" + policy + ":" + strconv.Itoa(priority))
	} else {
		prior, err := strconv.Atoi(strPriority)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		priority = prior
	}

	request, err := buildAlibabacloudStackSGRuleRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	direction := d.Get("type").(string)

	if direction == string(DirectionIngress) {
		request.ApiName = "ModifySecurityGroupRule"
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		addDebug(request.GetActionName(), raw, request.Headers, request)
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	} else {
		request.ApiName = "ModifySecurityGroupEgressRule"
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		addDebug(request.GetActionName(), raw, request.Headers, request)
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	return nil
}

func deleteSecurityGroupRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ruleType := d.Get("type").(string)
	request, err := buildAlibabacloudStackSGRuleRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if ruleType == string(DirectionIngress) {
		request.ApiName = "RevokeSecurityGroup"
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	} else {
		request.ApiName = "RevokeSecurityGroupEgress"
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	return nil
}

func resourceAlibabacloudStackSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	policy := parseSecurityRuleId(d, meta, 6)
	strPriority := parseSecurityRuleId(d, meta, 7)
	var priority int
	if policy == "" || strPriority == "" {
		policy = d.Get("policy").(string)
		priority = d.Get("priority").(int)
		d.SetId(d.Id() + ":" + policy + ":" + strconv.Itoa(priority))
	} else {
		prior, err := strconv.Atoi(strPriority)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		priority = prior
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := deleteSecurityGroupRule(d, meta)
		if err != nil {
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func buildAlibabacloudStackSGRuleRequest(d *schema.ResourceData, meta interface{}) (*requests.CommonRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	// Get product code from the built request
	request := client.NewCommonRequest("GET", "Ecs", "2014-05-26", "", "")

	direction := d.Get("type").(string)

	port_range := d.Get("port_range").(string)
	request.QueryParams["PortRange"] = port_range

	if v, ok := d.GetOk("ip_protocol"); ok {
		request.QueryParams["IpProtocol"] = v.(string)
		if v.(string) == string(Tcp) || v.(string) == string(Udp) {
			if port_range == AllPortRange {
				return nil, fmt.Errorf("'tcp' and 'udp' can support port range: [1, 65535]. Please correct it and try again.")
			}
		} else if port_range != AllPortRange {
			return nil, fmt.Errorf("'icmp', 'gre' and 'all' only support port range '-1/-1'. Please correct it and try again.")
		}
	}

	if v, ok := d.GetOk("policy"); ok {
		request.QueryParams["Policy"] = v.(string)
	}
	if v, ok := d.GetOk("nic_type"); ok {
		request.QueryParams["NicType"] = v.(string)
	}

	if v, ok := d.GetOk("priority"); ok {
		request.QueryParams["Priority"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("cidr_ip"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceCidrIp"] = v.(string)
		} else {
			request.QueryParams["DestCidrIp"] = v.(string)
		}
	} else if v, ok := d.GetOk("ipv6_cidr_ip"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["Ipv6SourceGroupId"] = v.(string)
		} else {
			request.QueryParams["Ipv6DestCidrIp"] = v.(string)
		}
	}

	var targetGroupId string
	if v, ok := d.GetOk("source_security_group_id"); ok {
		targetGroupId = v.(string)
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceGroupId"] = targetGroupId
		} else {
			request.QueryParams["DestGroupId"] = targetGroupId
		}
	}

	if v, ok := d.GetOk("source_group_owner_account"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceGroupOwnerAccount"] = v.(string)
		} else {
			request.QueryParams["DestGroupOwnerAccount"] = v.(string)
		}
	}

	sgId := d.Get("security_group_id").(string)

	group, err := ecsService.DescribeSecurityGroup(sgId)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	if v, ok := d.GetOk("nic_type"); ok {
		if group.VpcId != "" || targetGroupId != "" {
			if GroupRuleNicType(v.(string)) != GroupRuleIntranet {
				return nil, fmt.Errorf("When security group in the vpc or authorizing permission for source/destination security group, " + "the nic_type must be 'intranet'.")
			}
		}
		request.QueryParams["NicType"] = v.(string)
	}

	request.QueryParams["SecurityGroupId"] = sgId

	description := d.Get("description").(string)
	request.QueryParams["Description"] = description

	return request, nil
}

func parseSecurityRuleId(d *schema.ResourceData, meta interface{}, index int) (result string) {
	parts := strings.Split(d.Id(), ":")
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Panicing %s\r\n", e)
			result = ""
		}
	}()
	return parts[index]
}
