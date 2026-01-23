package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackNasNamespaceMountTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackNasNamespaceMountTargetsRead,

		Schema: map[string]*schema.Schema{
			"nas_namespace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"mount_targets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_target_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vsw_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackNasNamespaceMountTargetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	nasNamespaceId := d.Get("nas_namespace_id").(string)

	action := "DescribeNamespaceMountTargets"
	request := make(map[string]interface{})
	request["NasNamespaceId"] = nasNamespaceId

	// Add pagination parameters
	request["PageSize"] = 100

	idsMap := getIdsStringFilter(d)
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var objects []interface{}
	pageNumber := 1

	for {
		request["PageNumber"] = pageNumber
		response, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_nas_namespace_mount_targets", action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		totalCount, _ := response["TotalCount"].(json.Number).Int64()

		mountTargets, ok := response["MountTargets"].([]interface{})
		if ok {
			objects = append(objects, mountTargets...)
		}
		if int(totalCount) <= pageNumber*100 {
			break
		}
		pageNumber++
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)

	for _, object := range objects {
		obj := object.(map[string]interface{})
		mountTargetDomain := obj["MountTargetDomain"].(string)
		id := fmt.Sprintf("%s:%s", nasNamespaceId, mountTargetDomain)
		if nameRegex != nil && !nameRegex.MatchString(mountTargetDomain) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{
			"mount_target_domain": mountTargetDomain,
			"network_type":        obj["NetworkType"],
			"vpc_id":              obj["VpcId"],
			"vsw_id":              obj["VswId"],
			"access_group":        obj["AccessGroup"],
			"status":              obj["Status"],
		}

		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("mount_targets", s); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
