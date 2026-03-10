package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDataWorksProject() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				AtLeastOneOf:  []string{"name"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,26}$`),
					"project_name must be 3-27 characters long, start with a letter, and contain only letters, numbers, and underscores"),
			},
			"identifier": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"project_name"},
				AtLeastOneOf:  []string{"project_name"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,26}$`),
					"project_name must be 3-27 characters long, start with a letter, and contain only letters, numbers, and underscores"),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"USER", "PROJECT"}, false),
				Default:      "PROJECT",
			},
			"allow_download": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Available",
				ValidateFunc: validation.StringInSlice([]string{"Available", "Forbidden"}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return strings.EqualFold(oldValue, newValue)
				},
				DiffSuppressOnRefresh: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksProjectCreate,
		resourceAlibabacloudStackDataWorksProjectRead, resourceAlibabacloudStackDataWorksProjectUpdate, resourceAlibabacloudStackDataWorksProjectDelete)
	return resource
}

func resourceAlibabacloudStackDataWorksProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateProject"
	name := connectivity.GetResourceData(d, "name", "project_name")
	identifier := d.Get("identifier")
	if identifier == "" {
		identifier = name
	}
	request := map[string]interface{}{
		"ProjectName":       name,
		"ProjectIdentifier": identifier,
		"TaskAuthType":      d.Get("task_auth_type"),
		"ProjectDesc":       d.Get("description"),
	}

	response, err := client.DoTeaRequest("POST", "dataworks-private-cloud", "2019-01-17", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	if id, err := toInt(response["Data"]); err != nil {
		return err
	} else {
		d.SetId(strconv.Itoa(id))
	}

	dataworksService := DataworksService{client}
	stateConf := BuildStateConf([]string{"2"}, []string{"0"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, dataworksService.ProjectStateRefreshFunc(d.Id(), []string{"3", "6", "9"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackDataWorksProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksService := DataworksService{client}
	object, err := dataworksService.DescribeDataWorksProject(d.Id())
	log.Printf(fmt.Sprint(object))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		errmsg := ""
		if object != nil {
			errmsg = errmsgs.GetAsapiErrorMessage(object)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_data_works_project", "DescribeDataWorksProject", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	d.Set("project_id", object["ProjectId"])
	connectivity.SetResourceData(d, object["ProjectName"], "project_name", "name")
	d.Set("identifier", object["ProjectIdentifier"])
	d.Set("description", object["ProjectDescription"])
	if v, err := toInt(object["ProjectMode"]); err == nil {
		if v == 1 {
			d.Set("task_auth_type", "USER")
		} else {
			d.Set("task_auth_type", "PROJECT")
		}
	}
	d.Set("description", object["ProjectDescription"])
	if v, err := toInt(object["Status"]); err == nil {
		d.Set("status", DataworksProjectStatus[v])
	}
	if v, err := toInt(object["Status"]); err == nil {
		if v == 0 {
			d.Set("allow_download", false)
		} else {
			d.Set("allow_download", true)
		}
	}

	return nil
}

func resourceAlibabacloudStackDataWorksProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChanges("name", "project_name", "description", "allow_download", "status") {
		client := meta.(*connectivity.AlibabacloudStackClient)
		action := "UpdateProject"
		var status string
		if strings.EqualFold(d.Get("status").(string), "Available") {
			status = "0"
		} else if strings.EqualFold(d.Get("status").(string), "Forbidden") {
			status = "4"
		}
		request := map[string]interface{}{
			"ProjectId":          d.Id(),
			"ProjectDescription": d.Get("description"),
			"ProjectName":        connectivity.GetResourceData(d, "name", "project_name"),
			"Status":             status,
		}

		if d.Get("allow_download").(bool) {
			request["IsAllowDownload"] = 1
		} else {
			request["IsAllowDownload"] = 0
		}

		_, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		dataworksService := DataworksService{client}
		stateConf := BuildStateConf([]string{""}, []string{"0", "4"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, dataworksService.ProjectStateRefreshFunc(d.Id(), []string{"3", "6", "9"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	return nil
}

func resourceAlibabacloudStackDataWorksProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteProject"
	request := make(map[string]interface{})

	request["ProjectId"] = d.Id()
	if _, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request); err != nil {
		return err
	}
	dataworksService := DataworksService{client}
	stateConf := BuildStateConf([]string{"5"}, []string{}, d.Timeout(schema.TimeoutCreate), 3*time.Second, dataworksService.ProjectStateRefreshFunc(d.Id(), []string{"3", "6", "9"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
