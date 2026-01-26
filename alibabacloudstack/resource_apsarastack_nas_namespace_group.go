package alibabacloudstack

import (
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackNasNamespaceGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nas_namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mapped_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Once mapped, the path will be permanently bound to this mount point",
			},
			"mount_target_domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Once the network type is set, it cannot be changed unless the list is cleared",
				ValidateFunc: validation.StringInSlice([]string{"Classic", "Vpc"}, false),
			},
			"member_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNasNamespaceGroupCreate, resourceAlibabacloudStackNasNamespaceGroupRead, resourceAlibabacloudStackNasNamespaceGroupUpdate, resourceAlibabacloudStackNasNamespaceGroupDelete)
	return resource
}

func resourceAlibabacloudStackNasNamespaceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"NasNamespaceId":     d.Get("nas_namespace_id").(string),
		"NasNamespaceRegion": client.RegionId,
		"MappedPath":         d.Get("mapped_path").(string),
		"MountTargetDomain":  d.Get("mount_target_domain").(string),
		"NetworkType":        d.Get("network_type").(string),
	}

	_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "AddNamespaceToGroup", "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(d.Get("mount_target_domain").(string))

	return nil
}

func resourceAlibabacloudStackNasNamespaceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}

	object, err := nasService.DescribeNasNamespaceGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_nas_namespace_group nasService.DescribeNasNamespaceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("nas_namespace_id", object["NasNamespaceId"])
	d.Set("mapped_path", object["MappedPath"])
	d.Set("mount_target_domain", object["MountTargetDomain"])
	d.Set("network_type", object["NetworkType"])
	d.Set("member_id", object["MemberId"])
	d.Set("status", object["Status"])
	d.Set("create_time", object["CreateTime"])

	return nil
}

func resourceAlibabacloudStackNasNamespaceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlibabacloudStackNasNamespaceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"MemberId": d.Get("member_id").(string),
	}

	_, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", "RemoveNamespaceFromGroup", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "RemoveNamespaceFromGroup", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	return nil
}
