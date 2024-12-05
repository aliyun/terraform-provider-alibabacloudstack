package alibabacloudstack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

func resourceAlibabacloudStackCmsSiteMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCmsSiteMonitorCreate,
		Read:   resourceAlibabacloudStackCmsSiteMonitorRead,
		Update: resourceAlibabacloudStackCmsSiteMonitorUpdate,
		Delete: resourceAlibabacloudStackCmsSiteMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{SiteMonitorHTTP, SiteMonitorDNS, SiteMonitorFTP, SiteMonitorPOP3, SiteMonitorPing, SiteMonitorSMTP, SiteMonitorTCP, SiteMonitorUDP}, false),
			},
			"alert_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 5, 15}),
			},
			"options_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isp_cities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Required: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"task_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackCmsSiteMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}

	taskName := d.Get("task_name").(string)
	request := cms.CreateCreateSiteMonitorRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.Address = d.Get("address").(string)
	request.TaskName = taskName
	request.TaskType = d.Get("task_type").(string)
	request.Interval = strconv.Itoa(d.Get("interval").(int))
	request.OptionsJson = d.Get("options_json").(string)
	alertIds := d.Get("alert_ids").([]interface{})
	alertId := getAlertId(alertIds)
	if alertId != "" {
		request.AlertIds = alertId
	}

	if isp_cities, ok := d.GetOk("isp_cities"); ok {
		var a []map[string]interface{}
		for _, element := range isp_cities.(*schema.Set).List() {
			isp_city := element.(map[string]interface{})
			a = append(a, isp_city)
		}
		b, err := json.Marshal(a)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.IspCities = bytes.NewBuffer(b).String()
	}

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.CreateSiteMonitor(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*cms.CreateSiteMonitorResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_site_monitor", "CreateSiteMonitor", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	siteMonitor, err := cmsService.DescribeSiteMonitor("", taskName)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.SetId(siteMonitor.TaskId)

	return resourceAlibabacloudStackCmsSiteMonitorRead(d, meta)
}

func resourceAlibabacloudStackCmsSiteMonitorRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}

	siteMonitor, err := cmsService.DescribeSiteMonitor(d.Id(), "")
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("address", siteMonitor.Address)
	d.Set("task_name", siteMonitor.TaskName)
	d.Set("task_type", siteMonitor.TaskType)
	d.Set("task_state", siteMonitor.TaskState)
	d.Set("interval", siteMonitor.Interval)
	d.Set("options_json", siteMonitor.OptionsJson)
	d.Set("create_time", siteMonitor.CreateTime)
	d.Set("update_time", siteMonitor.UpdateTime)

	ispCities, err := cmsService.GetIspCities(d.Id())
	var list []map[string]interface{}

	for _, e := range ispCities {
		list = append(list, map[string]interface{}{"city": e["city"], "isp": e["isp"]})
	}

	if err = d.Set("isp_cities", list); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackCmsSiteMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := cms.CreateModifySiteMonitorRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.TaskId = d.Id()
	request.Address = d.Get("address").(string)
	request.Interval = strconv.Itoa(d.Get("interval").(int))
	request.OptionsJson = d.Get("options_json").(string)
	request.TaskName = d.Get("task_name").(string)
	alertIds := d.Get("alert_ids").([]interface{})
	alertId := getAlertId(alertIds)
	if alertId != "" {
		request.AlertIds = alertId
	}

	if isp_cities, ok := d.GetOk("isp_cities"); ok {
		var a []map[string]interface{}
		for _, element := range isp_cities.(*schema.Set).List() {
			isp_city := element.(map[string]interface{})
			a = append(a, isp_city)
		}
		b, err := json.Marshal(a)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.IspCities = bytes.NewBuffer(b).String()
	}

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.ModifySiteMonitor(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*cms.ModifySiteMonitorResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_site_monitor", "ModifySiteMonitor", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	return resourceAlibabacloudStackCmsSiteMonitorRead(d, meta)
}

func resourceAlibabacloudStackCmsSiteMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}
	request := cms.CreateDeleteSiteMonitorsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.TaskIds = d.Id()
	request.IsDeleteAlarms = "false"

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteSiteMonitors(request)
		})

		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*cms.DeleteSiteMonitorsResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting site monitor got an error: %#v. %s", err, errmsg))
		}

		_, err = cmsService.DescribeSiteMonitor(d.Id(), "")
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("DescribeSiteMonitor got an error: %#v", err))
		}

		return resource.RetryableError(fmt.Errorf("Deleting site monitor got an error: %#v", err))

	})
}

func getAlertId(alertIds []interface{}) string {
	if alertIds != nil && len(alertIds) > 0 {
		alertId := strings.Join(expandStringList(alertIds), ",")
		return alertId
	}
	return ""
}
