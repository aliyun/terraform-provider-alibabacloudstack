package alibabacloudstack

import (
	"regexp"
	"sort"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCrEeNamespaces() *schema.Resource {
	return &schema.Resource{
		Read:	dataSourceAlibabacloudStackCrEeNamespacesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:		schema.TypeString,
				ForceNew:	true,
				Required:	true,
			},
			"name_regex": {
				Type:		schema.TypeString,
				Optional:	true,
				ValidateFunc:	validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:		schema.TypeString,
				Optional:	true,
			},

			// Computed values
			"ids": {
				Type:		schema.TypeList,
				Optional:	true,
				Computed:	true,
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:		schema.TypeList,
				Computed:	true,
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			"namespaces": {
				Type:		schema.TypeList,
				Computed:	true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"id": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"name": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"auto_create": {
							Type:		schema.TypeBool,
							Computed:	true,
						},
						"default_visibility": {
							Type:		schema.TypeString,
							Computed:	true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCrEeNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	pageNo := 1
	pageSize := 50
	instanceId := d.Get("instance_id").(string)

	var namespaces []cr_ee.NamespacesItem
	for {
		resp, err := crService.ListCrEeNamespaces(instanceId, pageNo, pageSize)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		namespaces = append(namespaces, resp.Namespaces...)
		if len(resp.Namespaces) < pageSize {
			break
		}
		pageNo++
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var idsMap map[string]string
	if v, ok := d.GetOk("ids"); ok {
		idsMap = make(map[string]string)
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var targetNamespaces []cr_ee.NamespacesItem
	for _, namespace := range namespaces {
		if nameRegex != nil && !nameRegex.MatchString(namespace.NamespaceName) {
			continue
		}

		if idsMap != nil && idsMap[namespace.NamespaceId] == "" {
			continue
		}

		targetNamespaces = append(targetNamespaces, namespace)
	}

	namespaces = targetNamespaces

	sort.SliceStable(namespaces, func(i, j int) bool {
		return namespaces[i].NamespaceName < namespaces[j].NamespaceName
	})

	var (
		ids		[]string
		names		[]string
		namespaceMaps	[]map[string]interface{}
	)

	for _, namespace := range namespaces {
		mapping := make(map[string]interface{})
		mapping["instance_id"] = namespace.InstanceId
		mapping["id"] = namespace.NamespaceId
		mapping["name"] = namespace.NamespaceName
		mapping["auto_create"] = namespace.AutoCreateRepo
		mapping["default_visibility"] = namespace.DefaultRepoType

		ids = append(ids, namespace.NamespaceId)
		names = append(names, namespace.NamespaceName)
		namespaceMaps = append(namespaceMaps, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("namespaces", namespaceMaps); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		if err := writeToFile(output.(string), namespaceMaps); err != nil {
			return err
		}
	}

	return nil
}
