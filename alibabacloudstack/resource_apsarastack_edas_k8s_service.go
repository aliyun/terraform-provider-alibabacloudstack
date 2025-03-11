package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
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
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ClusterIP",
				ValidateFunc: validation.StringInSlice([]string{"ClusterIP", "NodePort", "LoadBalancer"}, false),
			},
			"port_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
						},
						"service_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"external_traffic_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Cluster"}, false),
			},
			"cluster_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inner_endpointer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodeip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"allow_edit": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackEdasK8sServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	service_type := d.Get("type").(string)
	request := client.NewCommonRequest("POST", "Edas", "2017-08-01", "CreateK8sService", "/pop/v5/k8s/service/service")
	request.QueryParams["AppId"] = d.Get("app_id").(string)

	request.QueryParams["Act"] = "1" // 创建接口默认值
	request.QueryParams["ServiceName"] = d.Get("service_name").(string)
	request.QueryParams["Type"] = service_type
	port_mappings := d.Get("port_mappings").([]interface{})
	k8s_port_mappings, err := edasService.GetK8sServicePorts(port_mappings)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.QueryParams["PortMappingsStrs"] = k8s_port_mappings
	if service_type != "ClusterIP" {
		request.QueryParams["ExternalTrafficPolicy"] = d.Get("external_traffic_policy").(string)
	}
	if v, ok := d.GetOk("annotations"); ok {
		AnnotationsStrs, err := json.Marshal(v.(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["AnnotationsStrs"] = string(AnnotationsStrs)
	}
	if v, ok := d.GetOk("labels"); ok {
		labelsStrs, err := json.Marshal(v.(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["LabelsStrs"] = string(labelsStrs)
	}
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("CreateK8sService", bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_application", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var response map[string]interface{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.WrapError(fmt.Errorf("Create edas k8s service failed for %s", response["Message"].(string)))
	}

	d.SetId(d.Get("app_id").(string) + ":" + d.Get("service_name").(string))
	return resourceAlibabacloudStackEdasK8sServiceUpdate(d, meta)
}

func resourceAlibabacloudStackEdasK8sServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	id := d.Id()
	service, err := edasService.DescribeEdasK8sService(id)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("app_id", strings.Split(id, ":")[0])
	d.Set("type", service.Type)
	d.Set("service_name", service.ServiceName)
	d.Set("allow_edit", service.AllowEdit)
	d.Set("inner_endpointer", service.InnerEndpointer)
	d.Set("namespace", service.Namespace)
	d.Set("nodeip_list", service.NodeIpList)
	d.Set("external_traffic_policy", service.ExternalTrafficPolicy)
	port_mappings := make([]map[string]interface{}, 0)
	for _, portMappings := range service.PortMappings {
		service_port, _ := strconv.Atoi(portMappings.ServicePort)
		target_port, _ := strconv.Atoi(portMappings.TargetPort)
		port_mappings = append(port_mappings, map[string]interface{}{
			"protocol":     portMappings.Protocol,
			"service_port": service_port,
			"target_port":  target_port,
		})
	}
	d.Set("port_mappings", port_mappings)
	if len(service.Labels) > 0 {
		labels := d.Get("labels").(map[string]interface{})
		new_labels := make(map[string]interface{})
		for k, _ := range labels {
			if v, ok := service.Labels[k]; ok {
				new_labels[k] = v
			}
		}
		d.Set("labels", new_labels)
	}
	d.Set("annotations", service.Annotations)
	return nil
}

func resourceAlibabacloudStackEdasK8sServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	request := client.NewCommonRequest("PUT", "Edas", "2017-08-01", "UpdateK8sService", "/pop/v5/k8s/acs/k8s_service")
	d.Partial(true)
	request.QueryParams["AppId"] = d.Get("app_id").(string)
	request.QueryParams["Name"] = d.Get("service_name").(string)
	request.QueryParams["Type"] = d.Get("type").(string)
	request.QueryParams["Act"] = "2" // 更新接口默认值
	port_mappings, err := edasService.GetK8sServicePorts(d.Get("port_mappings").([]interface{}))
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.QueryParams["PortMappingsStrs"] = port_mappings
	update := false
	if d.HasChange("service_name") {
		update = true
	}
	if d.HasChange("port_mappings") {
		update = true
	}
	if d.HasChange("type") {
		update = true
	}
	if d.HasChange("external_traffic_policy") {
		update = true
		request.QueryParams["ExternalTrafficPolicy"] = d.Get("external_traffic_policy").(string)
	}
	if d.HasChange("annotations") {
		AnnotationsStrs, err := json.Marshal(d.Get("annotations").(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["AnnotationsStrs"] = string(AnnotationsStrs)
	}
	if d.HasChange("labels") {
		service, err := edasService.DescribeEdasK8sService(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		labels := d.Get("labels").(map[string]interface{})
		for k, v := range service.Labels {
			if _, ok := labels[k]; !ok {
				labels[k] = v
			}
		}
		labelsStrs, err := json.Marshal(labels)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["LabelsStrs"] = string(labelsStrs)
	}
	if update && !d.IsNewResource() {
		bresponse, err := client.ProcessCommonRequest(request)
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
		if fmt.Sprint(response["Code"]) != "200" {
			return errmsgs.WrapError(errmsgs.Error("update edas k8s service failed for " + response["Message"].(string)))
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackEdasK8sServiceRead(d, meta)
}

func resourceAlibabacloudStackEdasK8sServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("DELETE", "Edas", "2017-08-01", "DeleteK8sService", "/pop/v5/k8s/service/service")
	parts := strings.Split(d.Id(), ":")
	app_id := parts[0]
	name := parts[1]
	request.QueryParams = map[string]string{
		"AppId":       app_id,
		"ServiceName": name,
	}
	bresponse, err := client.ProcessCommonRequest(request)
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
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.WrapError(errmsgs.Error("delete edas k8s service failed for " + response["Message"].(string)))
	}
	return nil
}
