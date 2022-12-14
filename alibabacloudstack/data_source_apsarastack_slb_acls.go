package alibabacloudstack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackSlbAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSlbAclsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"tags": tagsSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entry_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entry": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"comment": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MinItems: 0,
						},
						"related_listeners": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"frontend_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"acl_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MinItems: 0,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSlbAclsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := slb.CreateDescribeAccessControlListsRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	tags := d.Get("tags").(map[string]interface{})
	if tags != nil && len(tags) > 0 {
		KeyPairsTags := make([]slb.DescribeAccessControlListsTag, 0, len(tags))
		for k, v := range tags {
			keyPairsTag := slb.DescribeAccessControlListsTag{
				Key:   k,
				Value: v.(string),
			}
			KeyPairsTags = append(KeyPairsTags, keyPairsTag)
		}
		request.Tag = &KeyPairsTags
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAccessControlLists(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_slb_acls", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeAccessControlListsResponse)
	var filteredAclsTemp []slb.Acl
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, acl := range response.Acls.Acl {
			if r != nil && !r.MatchString(acl.AclName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[acl.AclId]; !ok {
					continue
				}
			}

			filteredAclsTemp = append(filteredAclsTemp, acl)
		}
	} else {
		filteredAclsTemp = response.Acls.Acl
	}

	return slbAclsDescriptionAttributes(d, filteredAclsTemp, client, meta)
}

func aclTagsMappings(d *schema.ResourceData, aclId string, meta interface{}) map[string]string {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	tags, err := slbService.DescribeTags(aclId, nil, TagResourceAcl)

	if err != nil {
		return nil
	}

	return slbTagsToMap(tags)
}

func slbAclsDescriptionAttributes(d *schema.ResourceData, acls []slb.Acl, client *connectivity.AlibabacloudStackClient, meta interface{}) error {

	var ids []string
	var names []string
	var s []map[string]interface{}
	slbService := SlbService{client}

	request := slb.CreateDescribeAccessControlListAttributeRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	for _, item := range acls {
		request.AclId = item.AclId
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeAccessControlListAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_slb_acls", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*slb.DescribeAccessControlListAttributeResponse)
		mapping := map[string]interface{}{
			"id":                response.AclId,
			"name":              response.AclName,
			"ip_version":        response.AddressIPVersion,
			"entry_list":        slbService.FlattenSlbAclEntryMappings(response.AclEntrys.AclEntry),
			"related_listeners": slbService.flattenSlbRelatedListenerMappings(response.RelatedListeners.RelatedListener),
			"tags":              aclTagsMappings(d, response.AclId, meta),
		}

		ids = append(ids, response.AclId)
		names = append(names, response.AclName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("acls", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
