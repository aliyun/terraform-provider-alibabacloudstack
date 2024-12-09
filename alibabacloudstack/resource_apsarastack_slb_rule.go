package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlbRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSlbRuleCreate,
		Read:   resourceAlibabacloudStackSlbRuleRead,
		Update: resourceAlibabacloudStackSlbRuleUpdate,
		Delete: resourceAlibabacloudStackSlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Required:     true,
				ForceNew:     true,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'rule_name' instead.",
				ConflictsWith: []string{"rule_name"},
			},

			"rule_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ConflictsWith: []string{"name"},
			},

			"listener_sync": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
				Optional:     true,
				Default:      string(OnFlag),
			},
			"scheduler": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"wrr", "wlc", "rr"}, false),
				Optional:         true,
				Default:          WRRScheduler,
				DiffSuppressFunc: slbRuleListenerSyncDiffSuppressFunc,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cookie": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringLenBetween(1, 200),
				Optional:         true,
				DiffSuppressFunc: slbRuleCookieDiffSuppressFunc,
			},
			"cookie_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 86400),
				Optional:         true,
				DiffSuppressFunc: slbRuleCookieTimeoutDiffSuppressFunc,
			},
			"health_check": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				Optional:         true,
				Default:          OnFlag,
				DiffSuppressFunc: slbRuleListenerSyncDiffSuppressFunc,
			},
			"health_check_http_code": {
				Type:             schema.TypeString,
				ValidateFunc:     validateAllowedSplitStringValue([]string{string(HTTP_2XX), string(HTTP_3XX), string(HTTP_4XX), string(HTTP_5XX)}, ","),
				Optional:         true,
				Default:          HTTP_2XX,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_interval": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 50),
				Optional:         true,
				Default:          2,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_domain": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringLenBetween(1, 80),
				Optional:         true,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_uri": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringLenBetween(1, 80),
				Optional:         true,
				Default:          "/",
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_connect_port": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.Any(validation.IntBetween(1, 65535), validation.IntInSlice([]int{-520})),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 300),
				Optional:         true,
				Default:          5,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"healthy_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"unhealthy_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			//http & https
			"sticky_session": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				Optional:         true,
				Default:          OffFlag,
				DiffSuppressFunc: slbRuleListenerSyncDiffSuppressFunc,
			},
			//http & https
			"sticky_session_type": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{string(InsertStickySessionType), string(ServerStickySessionType)}, false),
				Optional:         true,
				DiffSuppressFunc: slbRuleStickySessionTypeDiffSuppressFunc,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlibabacloudStackSlbRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slb_id := d.Get("load_balancer_id").(string)
	port := d.Get("frontend_port").(int)
	name := strings.Trim(connectivity.GetResourceData(d, "rule_name", "name").(string), " ")
	group_id := strings.Trim(d.Get("server_group_id").(string), " ")

	var domain, url, rule string
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		url = v.(string)
	}

	if domain == "" && url == "" {
		return errmsgs.WrapError(errmsgs.Error("At least one 'domain' or 'url' must be set."))
	} else if domain == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Url':'%s','VServerGroupId':'%s'}]", name, url, group_id)
	} else if url == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','VServerGroupId':'%s'}]", name, domain, group_id)
	} else {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','Url':'%s','VServerGroupId':'%s'}]", name, domain, url, group_id)
	}

	request := slb.CreateCreateRulesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = slb_id
	request.ListenerPort = requests.NewInteger(port)
	request.RuleList = rule

	var raw interface{}
	var err error
	if err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateRules(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"BackendServer.configuring", "OperationFailed.ListenerStatusNotSupport"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*slb.CreateRulesResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_rule", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response, _ := raw.(*slb.CreateRulesResponse)
	d.SetId(response.Rules.Rule[0].RuleId)

	return resourceAlibabacloudStackSlbRuleUpdate(d, meta)
}

func resourceAlibabacloudStackSlbRuleRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbRule(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.RuleName, "rule_name", "name")
	d.Set("load_balancer_id", object.LoadBalancerId)
	if port, err := strconv.Atoi(object.ListenerPort); err != nil {
		return errmsgs.WrapError(err)
	} else {
		d.Set("frontend_port", port)
	}
	d.Set("domain", object.Domain)
	d.Set("url", object.Url)
	d.Set("server_group_id", object.VServerGroupId)
	d.Set("sticky_session", object.StickySession)
	d.Set("sticky_session_type", object.StickySessionType)
	d.Set("unhealthy_threshold", object.UnhealthyThreshold)
	d.Set("healthy_threshold", object.HealthyThreshold)
	d.Set("health_check_timeout", object.HealthCheckTimeout)
	d.Set("health_check_connect_port", object.HealthCheckConnectPort)
	d.Set("health_check_uri", object.HealthCheckURI)
	d.Set("health_check", object.HealthCheck)
	d.Set("health_check_http_code", object.HealthCheckHttpCode)
	d.Set("health_check_interval", object.HealthCheckInterval)
	d.Set("scheduler", object.Scheduler)
	d.Set("listener_sync", object.ListenerSync)
	d.Set("cookie_timeout", object.CookieTimeout)
	d.Set("cookie", object.Cookie)
	d.Set("health_check_domain", object.HealthCheckDomain)
	return nil
}

func resourceAlibabacloudStackSlbRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	fullUpdate := false
	request := slb.CreateSetRuleRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RuleId = d.Id()

	if listenerSync, ok := d.GetOk("listener_sync"); ok && listenerSync == string(OffFlag) {
		if stickySession := d.Get("sticky_session"); stickySession == string(OnFlag) {
			if _, ok := d.GetOk("sticky_session_type"); !ok {
				return errmsgs.WrapError(errmsgs.Error(`'sticky_session_type': required field is not set when the sticky_session is 'on'.`))
			}
		}
		if stickySessionType := d.Get("sticky_session_type"); stickySessionType == string(InsertStickySessionType) {
			if _, ok := d.GetOk("cookie_timeout"); !ok {
				return errmsgs.WrapError(errmsgs.Error(`'cookie_timeout': required field is not set when the sticky_session_type is 'insert'.`))
			}
		}
		if stickySessionType := d.Get("sticky_session_type"); stickySessionType == string(ServerStickySessionType) {
			if _, ok := d.GetOk("cookie"); !ok {
				return errmsgs.WrapError(errmsgs.Error(`'cookie': required field is not set when the sticky_session_type is 'server'.`))
			}
		}
	}

	if d.HasChange("server_group_id") {
		request.VServerGroupId = d.Get("server_group_id").(string)
		update = true
	}

	if d.HasChanges("rule_name", "name") {
		request.RuleName = connectivity.GetResourceData(d, "rule_name", "name").(string)
		update = true
	}

	fullUpdate = d.HasChanges("listener_sync","scheduler","cookie","cookie_timeout","health_check", "health_check_http_code",
		"health_check_interval","health_check_domain","health_check_uri","health_check_connect_port", "health_check_timeout",
		"healthy_threshold","unhealthy_threshold", "sticky_session", "sticky_session_type")

	if fullUpdate {
		request.ListenerSync = d.Get("listener_sync").(string)
		if listenerSync, ok := d.GetOk("listener_sync"); ok && listenerSync == string(OffFlag) {
			request.Scheduler = d.Get("scheduler").(string)
			request.HealthCheck = d.Get("health_check").(string)
			request.StickySession = d.Get("sticky_session").(string)
			if request.HealthCheck == string(OnFlag) {
				request.HealthCheckTimeout = requests.NewInteger(d.Get("health_check_timeout").(int))
				request.HealthCheckURI = d.Get("health_check_uri").(string)
				request.HealthyThreshold = requests.NewInteger(d.Get("healthy_threshold").(int))
				request.UnhealthyThreshold = requests.NewInteger(d.Get("unhealthy_threshold").(int))
				request.HealthCheckInterval = requests.NewInteger(d.Get("health_check_interval").(int))
				request.HealthCheckHttpCode = d.Get("health_check_http_code").(string)
				if healthCheckDomain, ok := d.GetOk("health_check_domain"); ok {
					request.HealthCheckDomain = healthCheckDomain.(string)
				}
				if healthCheckConnectPort, ok := d.GetOk("health_check_connect_port"); ok {
					request.HealthCheckConnectPort = requests.NewInteger(healthCheckConnectPort.(int))
				}
			}
			if request.StickySession == string(OnFlag) {
				request.StickySessionType = d.Get("sticky_session_type").(string)
				if request.StickySessionType == string(InsertStickySessionType) {
					request.CookieTimeout = requests.NewInteger(d.Get("cookie_timeout").(int))
				}
				if request.StickySessionType == string(ServerStickySessionType) {
					request.Cookie = d.Get("cookie").(string)
				}
			}
		}
	}

	if update || fullUpdate {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetRule(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*slb.SetRuleResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlibabacloudStackSlbRuleRead(d, meta)
}

func resourceAlibabacloudStackSlbRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return errmsgs.WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return errmsgs.WrapError(fmt.Errorf("Current rule's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the rule.", lbId))
		}
	}

	request := slb.CreateDeleteRulesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RuleIds = fmt.Sprintf("['%s']", d.Id())

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteRules(request)
		})
		response, ok := raw.(*slb.DeleteRulesResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"OperationFailed.ListenerStatusNotSupport"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_rule", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidRuleId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, "")
	}
	return errmsgs.WrapError(slbService.WaitForSlbRule(d.Id(), Deleted, DefaultTimeoutMedium))
}
