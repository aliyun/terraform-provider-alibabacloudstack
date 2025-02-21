package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEdasK8sService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEdasK8sServiceCreate,
		Update: resourceAlibabacloudStackEdasK8sServiceUpdate,
		Read:   resourceAlibabacloudStackEdasK8sServiceRead,
		Delete: resourceAlibabacloudStackEdasK8sServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ClusterIP"}, false),
			},
			"service_ports": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"external_traffic_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Cluster"}, false),
			},
			"cluster_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackEdasK8sServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	request := client.NewCommonRequest("POST", "Edas", "2017-08-01", "CreateK8sService", "/pop/v5/k8s/acs/k8s_service")
	request.QueryParams["AppId"] = d.Get("app_id").(string)
	request.QueryParams["Name"] = d.Get("name").(string)
	request.QueryParams["Type"] = d.Get("type").(string)
	service_ports := d.Get("service_ports").([]interface{})
	k8s_service_ports, err := edasService.GetK8sServicePorts(service_ports)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.QueryParams["ServicePorts"] = k8s_service_ports
	if v, ok := d.GetOk("external_traffic_policy"); ok && v.(string) != "" {
		request.QueryParams["ExternalTrafficPolicy"] = v.(string)
	}
	bresponse, err := client.ProcessCommonRequestForOrganization(request)
	addDebug("CreateK8sService", bresponse, request.QueryParams, request)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_application", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var response map[string]interface{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if response["Code"].(float64) != 200 {
		return errmsgs.WrapError(fmt.Errorf("Create edas k8s service failed for %s", response["Message"].(string)))
	}

	d.SetId(d.Get("app_id").(string) + ":" + d.Get("name").(string))
	return resourceAlibabacloudStackEdasK8sServiceUpdate(d, meta)
}

func resourceAlibabacloudStackEdasK8sServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	id := d.Id()
	service, err := edasService.DescribeEdasService(id)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)

	}
	d.Set("type", service.Type)
	d.Set("name", service.Name)
	d.Set("cluster_ip", service.ClusterIP)
	service_ports := make([]map[string]interface{}, 0)
	for _, service_port := range service.ServicePorts {
		service_ports = append(service_ports, map[string]interface{}{
			"protocol":    service_port.Protocol,
			"port":        service_port.Port,
			"target_port": service_port.TargetPort,
		})
	}
	d.Set("service_ports", service_ports)

	return nil
}

func resourceAlibabacloudStackEdasK8sServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	request := client.NewCommonRequest("PUT", "Edas", "2017-08-01", "UpdateK8sService", "/pop/v5/k8s/acs/k8s_service")
	d.Partial(true)
	request.QueryParams["AppId"] = d.Get("app_id").(string)
	request.QueryParams["Name"] = d.Get("name").(string)
	request.QueryParams["Type"] = d.Get("type").(string)
	service_ports, err := edasService.GetK8sServicePorts(d.Get("service_ports").([]interface{}))
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.QueryParams["ServicePorts"] = service_ports
	update := false
	if d.HasChange("name") {
		update = true
	}
	if d.HasChange("service_ports") {
		update = true
	}
	if d.HasChange("type") {
		update = true
	}
	if d.HasChange("external_traffic_policy") {
		update = true
		request.QueryParams["ExternalTrafficPolicy"] = d.Get("external_traffic_policy").(string)
	}
	if update {
		bresponse, err := client.ProcessCommonRequestForOrganization(request)
		addDebug(request.GetActionName(), bresponse, request)

		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		response := make(map[string]interface{})
		_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if response["Code"].(float64) != 200 {
			return errmsgs.WrapError(errmsgs.Error("update edas k8s service failed for " + response["Message"].(string)))
		}
	}
	return resourceAlibabacloudStackEdasK8sServiceRead(d, meta)
}

func resourceAlibabacloudStackEdasK8sServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("DELETE", "Edas", "2017-08-01", "DeleteK8sService", "/pop/v5/k8s/acs/k8s_service")
	parts := strings.Split(d.Id(), ":")
	app_id := parts[0]
	name := parts[1]
	request.QueryParams = map[string]string{
		"AppId": app_id,
		"Name":  name,
	}
	bresponse, err := client.ProcessCommonRequestForOrganization(request)
	addDebug(request.GetActionName(), bresponse, request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if response["Code"].(float64) != 200 {
		return errmsgs.WrapError(errmsgs.Error("delete edas k8s service failed for " + response["Message"].(string)))
	}
	return nil
}
