package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEhpcJobTemplate() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"array_request": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clock_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_line": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"job_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mem": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"package_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"queue": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"re_runable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"runas_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stderr_redirect_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stdout_redirect_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"thread": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"variables": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEhpcJobTemplateCreate, resourceAlibabacloudStackEhpcJobTemplateRead, resourceAlibabacloudStackEhpcJobTemplateUpdate, resourceAlibabacloudStackEhpcJobTemplateDelete)
	return resource
}

func resourceAlibabacloudStackEhpcJobTemplateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateJobTemplate"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("array_request"); ok {
		request["ArrayRequest"] = v
	}
	if v, ok := d.GetOk("clock_time"); ok {
		request["ClockTime"] = v
	}
	request["CommandLine"] = d.Get("command_line")
	if v, ok := d.GetOk("gpu"); ok {
		request["Gpu"] = v
	}
	request["Name"] = d.Get("job_template_name")
	if v, ok := d.GetOk("mem"); ok {
		request["Mem"] = v
	}
	if v, ok := d.GetOk("node"); ok {
		request["Node"] = v
	}
	if v, ok := d.GetOk("package_path"); ok {
		request["PackagePath"] = v
	}
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}
	if v, ok := d.GetOk("queue"); ok {
		request["Queue"] = v
	}
	if v, ok := d.GetOkExists("re_runable"); ok {
		request["ReRunable"] = v
	}
	if v, ok := d.GetOk("runas_user"); ok {
		request["RunasUser"] = v
	}
	if v, ok := d.GetOk("stderr_redirect_path"); ok {
		request["StderrRedirectPath"] = v
	}
	if v, ok := d.GetOk("stdout_redirect_path"); ok {
		request["StdoutRedirectPath"] = v
	}
	if v, ok := d.GetOk("task"); ok {
		request["Task"] = v
	}
	if v, ok := d.GetOk("thread"); ok {
		request["Thread"] = v
	}
	if v, ok := d.GetOk("variables"); ok {
		request["Variables"] = v
	}

	response, err = client.DoTeaRequest("GET", "ECS", "2018-04-12", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return
}

func resourceAlibabacloudStackEhpcJobTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ehpcService := EhpcService{client}
	object, err := ehpcService.DescribeEhpcJobTemplate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ehpc_job_template ehpcService.DescribeEhpcJobTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("array_request", object["ArrayRequest"])
	d.Set("clock_time", object["ClockTime"])
	d.Set("command_line", object["CommandLine"])
	if v, ok := object["Gpu"]; ok && fmt.Sprint(v) != "0" {
		d.Set("gpu", formatInt(v))
	}
	d.Set("mem", object["Mem"])
	if v, ok := object["Node"]; ok && fmt.Sprint(v) != "0" {
		d.Set("node", formatInt(v))
	}
	d.Set("package_path", object["PackagePath"])
	if v, ok := object["Priority"]; ok {
		d.Set("priority", formatInt(v))
	}
	d.Set("queue", object["Queue"])
	d.Set("job_template_name", object["Name"])

	if v, ok := object["ReRunable"]; ok {
		v, _ := strconv.ParseBool(v.(string))
		err = d.Set("re_runable", v)
	}

	d.Set("runas_user", object["RunasUser"])
	d.Set("stderr_redirect_path", object["StderrRedirectPath"])
	d.Set("stdout_redirect_path", object["StdoutRedirectPath"])
	if v, ok := object["Task"]; ok && fmt.Sprint(v) != "0" {
		d.Set("task", formatInt(v))
	}
	if v, ok := object["Thread"]; ok && fmt.Sprint(v) != "0" {
		d.Set("thread", formatInt(v))
	}
	d.Set("variables", object["Variables"])
	return nil
}

func resourceAlibabacloudStackEhpcJobTemplateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}

	request["CommandLine"] = d.Get("command_line")
	request["Name"] = d.Get("job_template_name")

	if v, ok := d.GetOk("array_request"); ok {
		request["ArrayRequest"] = v
	}
	if v, ok := d.GetOk("clock_time"); ok {
		request["ClockTime"] = v
	}
	if v, ok := d.GetOk("gpu"); ok {
		request["Gpu"] = v
	}
	if v, ok := d.GetOk("mem"); ok {
		request["Mem"] = v
	}
	if v, ok := d.GetOk("node"); ok {
		request["Node"] = v
	}
	if v, ok := d.GetOk("package_path"); ok {
		request["PackagePath"] = v
	}
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}
	if v, ok := d.GetOk("queue"); ok {
		request["Queue"] = v
	}
	if v, ok := d.GetOk("re_runable"); ok {
		request["ReRunable"] = v
	}
	if v, ok := d.GetOk("runas_user"); ok {
		request["RunasUser"] = v
	}
	if v, ok := d.GetOk("stderr_redirect_path"); ok {
		request["StderrRedirectPath"] = v
	}
	if v, ok := d.GetOk("stdout_redirect_path"); ok {
		request["StdoutRedirectPath"] = v
	}
	if v, ok := d.GetOk("task"); ok {
		request["Task"] = v
	}
	if v, ok := d.GetOk("thread"); ok {
		request["Thread"] = v
	}
	if v, ok := d.GetOk("variables"); ok {
		request["Variables"] = v
	}

	action := "EditJobTemplate"
	_, err = client.DoTeaRequest("GET", "ECS", "2018-04-12", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	return
}

func resourceAlibabacloudStackEhpcJobTemplateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteJobTemplates"
	request := map[string]interface{}{
		"Templates": fmt.Sprintf("[{\"Id\":\"%s\"}]", d.Id()),
	}

	_, err = client.DoTeaRequest("GET", "ECS", "2018-04-12", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}