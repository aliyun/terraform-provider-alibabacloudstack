package alibabacloudstack

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackSlbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSlbListenersRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"slb_listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frontend_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backend_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_slave_server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"persistence_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"established_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sticky_session": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sticky_session_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cookie_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cookie": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_connect_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_connect_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_http_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gzip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_proto": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSlbListenersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = d.Get("load_balancer_id").(string)

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	response, ok := raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_listeners", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	var filteredListenersTemp []slb.ListenerPortAndProtocol
	port := -1
	if v, ok := d.GetOk("frontend_port"); ok && v.(int) != 0 {
		port = v.(int)
	}
	protocol := ""
	if v, ok := d.GetOk("protocol"); ok && v.(string) != "" {
		protocol = v.(string)
	}
	var r *regexp.Regexp
	if despRegex, ok := d.GetOk("description_regex"); ok && despRegex.(string) != "" {
		r = regexp.MustCompile(despRegex.(string))
	}
	if port != -1 || protocol != "" || r != nil {
		for _, listener := range response.ListenerPortsAndProtocol.ListenerPortAndProtocol {
			if port != -1 && listener.ListenerPort != port {
				continue
			}
			if protocol != "" && listener.ListenerProtocol != protocol {
				continue
			}
			if r != nil && !r.MatchString(listener.Description) {
				continue
			}

			filteredListenersTemp = append(filteredListenersTemp, listener)
		}
	} else {
		filteredListenersTemp = response.ListenerPortsAndProtocol.ListenerPortAndProtocol
	}

	return slbListenersDescriptionAttributes(d, filteredListenersTemp, meta)
}

func slbListenersDescriptionAttributes(d *schema.ResourceData, listeners []slb.ListenerPortAndProtocol, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var ids []string
	var s []map[string]interface{}

	for _, listener := range listeners {
		mapping := map[string]interface{}{
			"frontend_port": listener.ListenerPort,
			"protocol":      listener.ListenerProtocol,
			"description":   listener.Description,
		}

		loadBalancerId := d.Get("load_balancer_id").(string)
		switch Protocol(listener.ListenerProtocol) {
		case Http:
			request := slb.CreateDescribeLoadBalancerHTTPListenerAttributeRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.LoadBalancerId = loadBalancerId
			request.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerHTTPListenerAttribute(request)
			})
			response, ok := raw.(*slb.DescribeLoadBalancerHTTPListenerAttributeResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_listeners", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			mapping["backend_port"] = response.BackendServerPort
			mapping["status"] = response.Status
			mapping["bandwidth"] = response.Bandwidth
			mapping["scheduler"] = response.Scheduler
			mapping["server_group_id"] = response.VServerGroupId
			mapping["sticky_session"] = response.StickySession
			mapping["sticky_session_type"] = response.StickySessionType
			mapping["cookie_timeout"] = response.CookieTimeout
			mapping["cookie"] = response.Cookie
			mapping["health_check"] = response.HealthCheck
			mapping["health_check_domain"] = response.HealthCheckDomain
			mapping["health_check_uri"] = response.HealthCheckURI
			mapping["health_check_connect_port"] = response.HealthCheckConnectPort
			mapping["healthy_threshold"] = response.HealthyThreshold
			mapping["unhealthy_threshold"] = response.UnhealthyThreshold
			mapping["health_check_timeout"] = response.HealthCheckTimeout
			mapping["health_check_interval"] = response.HealthCheckInterval
			mapping["health_check_http_code"] = response.HealthCheckHttpCode
			mapping["gzip"] = response.Gzip
			mapping["x_forwarded_for"] = response.XForwardedFor
			mapping["x_forwarded_for_slb_ip"] = response.XForwardedForSLBIP
			mapping["x_forwarded_for_slb_id"] = response.XForwardedForSLBID
			mapping["x_forwarded_for_slb_proto"] = response.XForwardedForProto
		case Https:
			request := slb.CreateDescribeLoadBalancerHTTPSListenerAttributeRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.LoadBalancerId = loadBalancerId
			request.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerHTTPSListenerAttribute(request)
			})
			response, ok := raw.(*slb.DescribeLoadBalancerHTTPSListenerAttributeResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_listeners", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			mapping["backend_port"] = response.BackendServerPort
			mapping["status"] = response.Status
			mapping["security_status"] = response.SecurityStatus
			mapping["bandwidth"] = response.Bandwidth
			mapping["scheduler"] = response.Scheduler
			mapping["server_group_id"] = response.VServerGroupId
			mapping["sticky_session"] = response.StickySession
			mapping["sticky_session_type"] = response.StickySessionType
			mapping["cookie_timeout"] = response.CookieTimeout
			mapping["cookie"] = response.Cookie
			mapping["health_check"] = response.HealthCheck
			mapping["health_check_domain"] = response.HealthCheckDomain
			mapping["health_check_uri"] = response.HealthCheckURI
			mapping["health_check_connect_port"] = response.HealthCheckConnectPort
			mapping["healthy_threshold"] = response.HealthyThreshold
			mapping["unhealthy_threshold"] = response.UnhealthyThreshold
			mapping["health_check_timeout"] = response.HealthCheckTimeout
			mapping["health_check_interval"] = response.HealthCheckInterval
			mapping["health_check_http_code"] = response.HealthCheckHttpCode
			mapping["gzip"] = response.Gzip
			mapping["server_certificate_id"] = response.ServerCertificateId
			mapping["x_forwarded_for"] = response.XForwardedFor
			mapping["x_forwarded_for_slb_ip"] = response.XForwardedForSLBIP
			mapping["x_forwarded_for_slb_id"] = response.XForwardedForSLBID
			mapping["x_forwarded_for_slb_proto"] = response.XForwardedForProto
		case Tcp:
			request := slb.CreateDescribeLoadBalancerTCPListenerAttributeRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.LoadBalancerId = loadBalancerId
			request.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerTCPListenerAttribute(request)
			})
			response, ok := raw.(*slb.DescribeLoadBalancerTCPListenerAttributeResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_listeners", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			mapping["backend_port"] = response.BackendServerPort
			mapping["status"] = response.Status
			mapping["bandwidth"] = response.Bandwidth
			mapping["scheduler"] = response.Scheduler
			mapping["server_group_id"] = response.VServerGroupId
			mapping["master_slave_server_group_id"] = response.MasterSlaveServerGroupId
			mapping["persistence_timeout"] = response.PersistenceTimeout
			mapping["established_timeout"] = response.EstablishedTimeout
			mapping["health_check"] = response.HealthCheck
			mapping["health_check_type"] = response.HealthCheckType
			mapping["health_check_domain"] = response.HealthCheckDomain
			mapping["health_check_uri"] = response.HealthCheckURI
			mapping["health_check_connect_port"] = response.HealthCheckConnectPort
			mapping["health_check_connect_timeout"] = response.HealthCheckConnectTimeout
			mapping["healthy_threshold"] = response.HealthyThreshold
			mapping["unhealthy_threshold"] = response.UnhealthyThreshold
			mapping["health_check_interval"] = response.HealthCheckInterval
			mapping["health_check_http_code"] = response.HealthCheckHttpCode
		case Udp:
			request := slb.CreateDescribeLoadBalancerUDPListenerAttributeRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.LoadBalancerId = loadBalancerId
			request.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerUDPListenerAttribute(request)
			})
			response, ok := raw.(*slb.DescribeLoadBalancerUDPListenerAttributeResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_listeners", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			mapping["backend_port"] = response.BackendServerPort
			mapping["status"] = response.Status
			mapping["bandwidth"] = response.Bandwidth
			mapping["scheduler"] = response.Scheduler
			mapping["server_group_id"] = response.VServerGroupId
			mapping["master_slave_server_group_id"] = response.MasterSlaveServerGroupId
			mapping["persistence_timeout"] = response.PersistenceTimeout
			mapping["health_check"] = response.HealthCheck
			mapping["health_check_connect_port"] = response.HealthCheckConnectPort
			mapping["health_check_connect_timeout"] = response.HealthCheckConnectTimeout
			mapping["healthy_threshold"] = response.HealthyThreshold
			mapping["unhealthy_threshold"] = response.UnhealthyThreshold
			mapping["health_check_interval"] = response.HealthCheckInterval
		}

		ids = append(ids, strconv.Itoa(listener.ListenerPort))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_listeners", s); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
