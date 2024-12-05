package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackExpressConnectVirtualBorderRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackExpressConnectVirtualBorderRouterCreate,
		Read:   resourceAlibabacloudStackExpressConnectVirtualBorderRouterRead,
		Update: resourceAlibabacloudStackExpressConnectVirtualBorderRouterUpdate,
		Delete: resourceAlibabacloudStackExpressConnectVirtualBorderRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"associated_physical_connections": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"circuit_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detect_multiplier": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(3, 10),
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"local_gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local_ipv6_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"min_rx_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(200, 1000),
			},
			"min_tx_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(200, 1000),
			},
			"peer_gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_ipv6_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_ipv6_subnet_mask": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_subnet_mask": {
				Type:     schema.TypeString,
				Required: true,
			},
			"physical_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "deleting", "recovering", "terminated", "terminating", "unconfirmed"}, false),
			},
			"vbr_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual_border_router_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 2999),
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackExpressConnectVirtualBorderRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateVirtualBorderRouter"
	request := client.NewCommonRequest("POST", "Vpc", "2016-04-28", action, "")
	mergeMaps(request.QueryParams, map[string]string{
		"LocalGatewayIp":       d.Get("local_gateway_ip").(string),
		"PeerGatewayIp":        d.Get("peer_gateway_ip").(string),
		"PeeringSubnetMask":    d.Get("peering_subnet_mask").(string),
		"PhysicalConnectionId": d.Get("physical_connection_id").(string),
		"VlanId":               fmt.Sprintf("%d", d.Get("vlan_id").(int)),
	})
	if v, ok := d.GetOk("bandwidth"); ok {
		request.QueryParams["Bandwidth"] = v.(string)
	}
	if v, ok := d.GetOk("circuit_code"); ok {
		request.QueryParams["CircuitCode"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		request.QueryParams["Description"] = v.(string)
	}
	if v, ok := d.GetOkExists("enable_ipv6"); ok {
		request.QueryParams["EnableIpv6"] = v.(string)
	}
	if v, ok := d.GetOk("local_ipv6_gateway_ip"); ok {
		request.QueryParams["LocalIpv6GatewayIp"] = v.(string)
	}
	if v, ok := d.GetOk("peer_ipv6_gateway_ip"); ok {
		request.QueryParams["PeerIpv6GatewayIp"] = v.(string)
	}
	if v, ok := d.GetOk("peering_ipv6_subnet_mask"); ok {
		request.QueryParams["PeeringIpv6SubnetMask"] = v.(string)
	}
	if v, ok := d.GetOk("vbr_owner_id"); ok {
		request.QueryParams["VbrOwnerId"] = v.(string)
	}
	if v, ok := d.GetOk("virtual_border_router_name"); ok {
		request.QueryParams["Name"] = v.(string)
	}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		bresponse, ok := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(fmt.Errorf("CreateVirtualBorderRouter Failed!!! %s", errmsg))
		}
		err = json.Unmarshal([]byte(bresponse.GetHttpContentString()), &response)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		d.SetId(fmt.Sprint(response["VbrId"]))
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_express_connect_virtual_border_router", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return resourceAlibabacloudStackExpressConnectVirtualBorderRouterUpdate(d, meta)
}

func resourceAlibabacloudStackExpressConnectVirtualBorderRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeExpressConnectVirtualBorderRouter(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_express_connect_virtual_border_router vpcService.DescribeExpressConnectVirtualBorderRouter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("circuit_code", object["CircuitCode"])
	d.Set("description", object["Description"])
	if v, ok := object["DetectMultiplier"]; ok && fmt.Sprint(v) != "0" {
		d.Set("detect_multiplier", formatInt(v))
	}
	d.Set("enable_ipv6", object["EnableIpv6"])
	d.Set("local_gateway_ip", object["LocalGatewayIp"])
	d.Set("local_ipv6_gateway_ip", object["LocalIpv6GatewayIp"])
	if v, ok := object["MinRxInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_rx_interval", formatInt(v))
	}
	if v, ok := object["MinTxInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_tx_interval", formatInt(v))
	}
	d.Set("peer_gateway_ip", object["PeerGatewayIp"])
	d.Set("peer_ipv6_gateway_ip", object["PeerIpv6GatewayIp"])
	d.Set("peering_ipv6_subnet_mask", object["PeeringIpv6SubnetMask"])
	d.Set("peering_subnet_mask", object["PeeringSubnetMask"])
	d.Set("physical_connection_id", object["PhysicalConnectionId"])
	d.Set("status", object["Status"])
	d.Set("virtual_border_router_name", object["Name"])
	if v, ok := object["VlanId"]; ok && fmt.Sprint(v) != "0" {
		d.Set("vlan_id", formatInt(v))
	}
	d.Set("route_table_id", object["RouteTableId"])
	return nil
}

func resourceAlibabacloudStackExpressConnectVirtualBorderRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	d.Partial(true)

	update := false
	request := client.NewCommonRequest("POST", "Vpc", "2016-04-28", "ModifyVirtualBorderRouterAttribute", "")
	request.QueryParams["VbrId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("circuit_code") {
		update = true
		if v, ok := d.GetOk("circuit_code"); ok {
			request.QueryParams["CircuitCode"] = v.(string)
		}
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request.QueryParams["Description"] = v.(string)
		}
	}
	if d.HasChange("detect_multiplier") {
		update = true
		if v, ok := d.GetOk("detect_multiplier"); ok {
			request.QueryParams["DetectMultiplier"] = fmt.Sprint(v)
		} else if v, ok := d.GetOk("min_rx_interval"); ok && fmt.Sprint(v) != "" {
			if v, ok := d.GetOk("min_tx_interval"); ok && fmt.Sprint(v) != "" {
				return errmsgs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "detect_multiplier", "min_rx_interval", d.Get("min_rx_interval"), "min_tx_interval", d.Get("min_tx_interval")))
			}
		}
		request.QueryParams["MinRxInterval"] = fmt.Sprint(d.Get("min_rx_interval"))
		request.QueryParams["MinTxInterval"] = fmt.Sprint(d.Get("min_tx_interval"))
	}
	if !d.IsNewResource() && d.HasChange("enable_ipv6") {
		update = true
		if v, ok := d.GetOkExists("enable_ipv6"); ok {
			request.QueryParams["EnableIpv6"] = v.(string)
		}
	}
	if !d.IsNewResource() && d.HasChange("local_gateway_ip") {
		update = true
		request.QueryParams["LocalGatewayIp"] = d.Get("local_gateway_ip").(string)
	}
	if !d.IsNewResource() && d.HasChange("local_ipv6_gateway_ip") {
		update = true
		if v, ok := d.GetOk("local_ipv6_gateway_ip"); ok {
			request.QueryParams["LocalIpv6GatewayIp"] = v.(string)
		}
	}
	if d.HasChange("min_rx_interval") {
		update = true
		if v, ok := d.GetOk("min_rx_interval"); ok {
			request.QueryParams["MinRxInterval"] = fmt.Sprint(v)
		} else if v, ok := d.GetOk("detect_multiplier"); ok && fmt.Sprint(v) != "" {
			if v, ok := d.GetOk("min_tx_interval"); ok && fmt.Sprint(v) != "" {
				return errmsgs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "min_rx_interval", "detect_multiplier", d.Get("detect_multiplier"), "min_tx_interval", d.Get("min_tx_interval")))
			}
		}
		request.QueryParams["DetectMultiplier"] = fmt.Sprint(d.Get("detect_multiplier"))
		request.QueryParams["MinTxInterval"] = fmt.Sprint(d.Get("min_tx_interval"))
	}
	if d.HasChange("min_tx_interval") {
		update = true
		if v, ok := d.GetOk("min_tx_interval"); ok {
			request.QueryParams["MinTxInterval"] = fmt.Sprint(v)
		} else if v, ok := d.GetOk("detect_multiplier"); ok && fmt.Sprint(v) != "" {
			if v, ok := d.GetOk("min_rx_interval"); ok && fmt.Sprint(v) != "" {
				return errmsgs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "min_tx_interval", "detect_multiplier", d.Get("detect_multiplier"), "min_rx_interval", d.Get("min_rx_interval")))
			}
		}
		request.QueryParams["DetectMultiplier"] = fmt.Sprint(d.Get("detect_multiplier"))
		request.QueryParams["MinRxInterval"] = fmt.Sprint(d.Get("min_rx_interval"))
	}
	if !d.IsNewResource() && d.HasChange("peer_gateway_ip") {
		update = true
		request.QueryParams["PeerGatewayIp"] = d.Get("peer_gateway_ip").(string)
	}
	if !d.IsNewResource() && d.HasChange("peer_ipv6_gateway_ip") {
		update = true
		if v, ok := d.GetOk("peer_ipv6_gateway_ip"); ok {
			request.QueryParams["PeerIpv6GatewayIp"] = v.(string)
		}
	}
	if !d.IsNewResource() && d.HasChange("peering_ipv6_subnet_mask") {
		update = true
		if v, ok := d.GetOk("peering_ipv6_subnet_mask"); ok {
			request.QueryParams["PeeringIpv6SubnetMask"] = v.(string)
		}
	}
	if !d.IsNewResource() && d.HasChange("peering_subnet_mask") {
		update = true
		request.QueryParams["PeeringSubnetMask"] = d.Get("peering_subnet_mask").(string)
	}
	if !d.IsNewResource() && d.HasChange("virtual_border_router_name") {
		update = true
		if v, ok := d.GetOk("virtual_border_router_name"); ok {
			request.QueryParams["Name"] = v.(string)
		}
	}
	if !d.IsNewResource() && d.HasChange("vlan_id") {
		update = true
		request.QueryParams["VlanId"] = fmt.Sprint(d.Get("vlan_id"))
	}
	if update {
		if v, ok := d.GetOk("associated_physical_connections"); ok {
			request.QueryParams["AssociatedPhysicalConnections"] = v.(string)
		}
		if v, ok := d.GetOk("bandwidth"); ok {
			request.QueryParams["Bandwidth"] = v.(string)
		}
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request, request.QueryParams)
			bresponse, ok := raw.(*responses.CommonResponse)
			if bresponse.GetHttpStatus() != 200 {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.RetryableError(fmt.Errorf("ModifyVirtualBorderRouterAttribute Failed!!! %s", errmsg))
			}
			return nil
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	if d.HasChange("status") {
		object, err := vpcService.DescribeExpressConnectVirtualBorderRouter(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "active" {
				rqs := client.NewCommonRequest("POST", "Vpc", "2016-04-28", "RecoverVirtualBorderRouter", "")
				rqs.QueryParams["VbrId"] = d.Id()
				err := resource.Retry(3*time.Minute, func() *resource.RetryError {
					raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
						return vpcClient.ProcessCommonRequest(rqs)
					})
					if err != nil {
						return resource.NonRetryableError(err)
					}
					addDebug(rqs.GetActionName(), raw, rqs, rqs.QueryParams)
					bresponse, ok := raw.(*responses.CommonResponse)
					if bresponse.GetHttpStatus() != 200 {
						errmsg := ""
						if ok {
							errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
						}
						return resource.RetryableError(fmt.Errorf("RecoverVirtualBorderRouter Failed!!! %s", errmsg))
					}
					return nil
				})
				if err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), rqs.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.ExpressConnectVirtualBorderRouterStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
				}
			}
			if target == "terminated" {
				rqs := client.NewCommonRequest("POST", "Vpc", "2016-04-28", "TerminateVirtualBorderRouter", "")
				rqs.QueryParams["VbrId"] = d.Id()
				err := resource.Retry(3*time.Minute, func() *resource.RetryError {
					raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
						return vpcClient.ProcessCommonRequest(rqs)
					})
					if err != nil {
						return resource.NonRetryableError(err)
					}
					addDebug(rqs.GetActionName(), raw, rqs, rqs.QueryParams)
					bresponse, ok := raw.(*responses.CommonResponse)
					if bresponse.GetHttpStatus() != 200 {
						errmsg := ""
						if ok {
							errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
						}
						return resource.RetryableError(fmt.Errorf("TerminateVirtualBorderRouter Failed!!! %s", errmsg))
					}
					return nil
				})
				if err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), rqs.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"terminated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.ExpressConnectVirtualBorderRouterStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
				}
			}
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackExpressConnectVirtualBorderRouterRead(d, meta)
}

func resourceAlibabacloudStackExpressConnectVirtualBorderRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteVirtualBorderRouter"
	request := client.NewCommonRequest("POST", "Vpc", "2016-04-28", action, "")
	request.QueryParams["VbrId"] = d.Id()

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if errmsgs.NeedRetry(err) || errmsgs.IsExpectedErrors(err, []string{"DependencyViolation.BgpGroup"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		bresponse, ok := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(fmt.Errorf("DeleteVirtualBorderRouter Failed!!! %s", errmsg))
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
