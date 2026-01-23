package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackNasNamespaceMountTarget() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nas_namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Classic", "Vpc"}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_target_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNasNamespaceMountTargetCreate, resourceAlibabacloudStackNasNamespaceMountTargetRead, resourceAlibabacloudStackNasNamespaceMountTargetUpdate, resourceAlibabacloudStackNasNamespaceMountTargetDelete)
	return resource
}

func resourceAlibabacloudStackNasNamespaceMountTargetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}
	networkType := d.Get("network_type").(string)
	request := map[string]interface{}{
		"NasNamespaceId":  d.Get("nas_namespace_id").(string),
		"AccessGroupName": d.Get("access_group_name").(string),
		"NetworkType":     networkType,
	}

	if v, ok := d.GetOk("vswitch_id"); ok && networkType == "Vpc" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request["VpcId"] = vsw["VpcId"]
		request["VSwitchId"] = v.(string)
	}

	raw, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", "CreateNamespaceMountTarget", "", nil, nil, request)
	if err != nil {
		return err
	}

	mountTargetDomain := raw["MountTargetDomain"].(string)
	nasNamespaceId := d.Get("nas_namespace_id").(string)

	id := fmt.Sprintf("%s:%s", nasNamespaceId, mountTargetDomain)
	d.SetId(id)
	stateConf := BuildStateConf([]string{"Pending"}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nasService.NasNamespaceMountTargetStateRefreshFunc(id, []string{"Inactive", "Deleting", "Deleted"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	d.SetId(id)

	return nil
}

func resourceAlibabacloudStackNasNamespaceMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}

	object, err := nasService.DescribeNasNamespaceMountTarget(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_nas_namespace_mount_target nasService.DescribeNasNamespaceMountTarget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid resource id format, expected NasNamespaceId:MountTargetDomain")
	}
	d.Set("nas_namespace_id", parts[0])
	d.Set("access_group_name", object["AccessGroup"])
	d.Set("network_type", object["NetworkType"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VswId"])
	d.Set("mount_target_domain", object["MountTargetDomain"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlibabacloudStackNasNamespaceMountTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}

	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid resource id format, expected NasNamespaceId:MountTargetDomain")
	}
	nasNamespaceId := parts[0]
	mountTargetDomain := parts[1]
	if d.HasChange("status") {
		requestInfo := make(map[string]interface{})
		requestInfo["NasNamespaceId"] = nasNamespaceId
		requestInfo["MountTargetDomain"] = mountTargetDomain
		requestInfo["Status"] = d.Get("status")
		if _, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", "ModifyNamespaceMountTarget", "", nil, nil, requestInfo); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_nas_namespace_mount_target", "ModifyNamespaceMountTarget", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"Pending"}, []string{d.Get("status").(string)}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasNamespaceMountTargetStateRefreshFunc(d.Id(), []string{"Failed", "Error"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	if d.IsNewResource() {
		return nil
	}
	if d.HasChange("access_group_name") {
		requestInfo := make(map[string]interface{})
		requestInfo["NasNamespaceId"] = nasNamespaceId
		requestInfo["MountTargetDomain"] = mountTargetDomain
		requestInfo["AccessGroupName"] = d.Get("access_group_name")
		if _, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", "ModifyNamespaceMountTarget", "", nil, nil, requestInfo); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_nas_namespace_mount_target", "ModifyNamespaceMountTarget", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"Pending"}, []string{d.Get("status").(string)}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasNamespaceMountTargetStateRefreshFunc(d.Id(), []string{"Failed", "Error"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	return nil
}

func resourceAlibabacloudStackNasNamespaceMountTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}

	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid resource id format, expected NasNamespaceId:MountTargetDomain")
	}
	nasNamespaceId := parts[0]
	mountTargetDomain := parts[1]

	requestInfo := map[string]interface{}{
		"NasNamespaceId":    nasNamespaceId,
		"MountTargetDomain": mountTargetDomain,
	}

	_, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", "DeleteNamespaceMountTarget", "", nil, nil, requestInfo)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteNamespaceMountTarget", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
	}

	stateConf := BuildStateConf([]string{"Active", "Pending", "Inactive", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nasService.NasNamespaceMountTargetStateRefreshFunc(d.Id(), []string{}))

	_, err = stateConf.WaitForState()
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
