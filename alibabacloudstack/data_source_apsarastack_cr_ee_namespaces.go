package alibabacloudstack

import (
	"regexp"
	"sort"

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

	var namespaces []interface{}
	for {
		resp, err := crService.ListCrEeNamespaces(instanceId, pageNo, pageSize)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		respNamespaces := resp["Namespaces"].([]interface{})
		namespaces = append(namespaces, respNamespaces...)
		if len(respNamespaces) < pageSize {
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

	var targetNamespaces []map[string]interface{}
	for _, namespaceItem := range namespaces {
		namespace := namespaceItem.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(namespace["NamespaceName"].(string)) {
			continue
		}

		namespaceId := crService.GenResourceId(namespace["InstanceId"].(string), namespace["NamespaceName"].(string))
		if idsMap != nil && idsMap[namespaceId] == "" {
			continue
		}

		targetNamespaces = append(targetNamespaces, namespace)
	}

	sort.SliceStable(targetNamespaces, func(i, j int) bool {
		return targetNamespaces[i]["NamespaceName"].(string) < targetNamespaces[j]["NamespaceName"].(string)
	})

	var (
		ids		[]string
		names		[]string
		namespaceMaps	[]map[string]interface{}
	)

	for _, namespace := range targetNamespaces {
		mapping := make(map[string]interface{})
		mapping["instance_id"] = namespace["InstanceId"].(string)
		mapping["id"] = crService.GenResourceId(namespace["InstanceId"].(string), namespace["NamespaceName"].(string))
		mapping["name"] = namespace["NamespaceName"].(string)
		mapping["auto_create"] = namespace["AutoCreateRepo"].(bool)
		mapping["default_visibility"] = namespace["DefaultRepoType"].(string)

		ids = append(ids, mapping["id"].(string))
		names = append(names, namespace["NamespaceName"].(string))
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
