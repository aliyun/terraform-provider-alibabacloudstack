package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackBcmpKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBcmpKeyPairsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"finger_print": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_pairs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
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
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackBcmpKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var regex *regexp.Regexp
	if name, ok := d.GetOk("name_regex"); ok {
		regex = regexp.MustCompile(name.(string))
	}

	// ids
	idsMap := getIdsStringFilter(d)

	var keyPairs []map[string]interface{}
	pageNumber := 1

	for {
		request := map[string]interface{}{
			"PageNumber": pageNumber,
			"PageSize":   PageSizeLarge,
		}

		response, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListKeyPair", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_keypairs", "ListKeyPair", errmsgs.AlibabacloudStackSdkGoERROR, "")
		}

		if !response["success"].(bool) {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_bcmp_keypairs", "ListKeyPair", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		data, ok := response["data"].([]interface{})
		if !ok {
			break
		}

		if len(data) < 1 {
			break
		}

		for _, item := range data {
			keyPair, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			name, ok := keyPair["name"].(string)
			if !ok {
				continue
			}

			if regex != nil && !regex.MatchString(name) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[name]; !ok {
					continue
				}
			}

			mapping := map[string]interface{}{
				"id":            name,
				"key_pair_name": name,
				"finger_print":  keyPair["keyPairFingerPrint"],
				"region":        keyPair["region"],
				"create_time":   keyPair["createTime"],
				"update_time":   keyPair["updateTime"],
			}

			keyPairs = append(keyPairs, mapping)
		}

		if len(data) < PageSizeLarge {
			break
		}

		pageNumber++
	}

	return bcmpKeyPairsDescriptionAttributes(d, keyPairs, meta)
}

func bcmpKeyPairsDescriptionAttributes(d *schema.ResourceData, keyPairs []map[string]interface{}, meta interface{}) error {
	var names []string
	var ids []string

	for _, key := range keyPairs {
		names = append(names, key["key_pair_name"].(string))
		ids = append(ids, key["id"].(string))
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("key_pairs", keyPairs); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
