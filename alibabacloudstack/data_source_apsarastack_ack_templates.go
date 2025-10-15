package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAckTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAckTemplatesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"template_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_with_hist_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_hash_code_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ali_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAckTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request query params
	queryParams := make(map[string]string)
	queryParams["TemplateType"] = "kubernetes"

	// Call API to get templates list
	resp, err := client.DoTeaRequest("GET", "cs", "2015-12-15", "DescribeTemplates", "/templates", queryParams, nil, nil)
	if err != nil {
		return err
	}

	// Parse response
	templatesRaw, ok := resp["templates"]
	if !ok || templatesRaw == nil {
		templatesRaw = []interface{}{}
	}

	templatesList := templatesRaw.([]interface{})
	if len(templatesList) == 0 {
		templatesList = []interface{}{}
	}

	// Process filters
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		descriptionRegex = regexp.MustCompile(v.(string))
	}

	var filteredTemplates []interface{}
	for _, item := range templatesList {
		template := item.(map[string]interface{})

		// Filter by ids
		if len(idsMap) > 0 {
			id := template["id"].(string)
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}

		// Filter by name_regex
		if nameRegex != nil {
			name := template["name"].(string)
			if !nameRegex.MatchString(name) {
				continue
			}
		}

		if descriptionRegex != nil {
			description := template["description"].(string)
			if !descriptionRegex.MatchString(description) {
				continue
			}
		}

		filteredTemplates = append(filteredTemplates, template)
	}

	// Prepare result
	var ids []string
	var templates []map[string]interface{}

	for _, item := range filteredTemplates {
		template := item.(map[string]interface{})

		mapping := map[string]interface{}{
			"template":                   template["template"],
			"name":                       template["name"],
			"description":                template["description"],
			"template_type":              template["template_type"],
			"template_id":                template["id"],
			"template_with_hist_id":      template["template_with_hist_id"],
			"template_hash_code_version": template["template_hash_code_version"],
			"created":                    template["created"],
			"acl":                        template["acl"],
			"version":                    template["version"],
			"tags":                       template["tags"],
			"ali_uid":                    template["ali_uid"],
			"updated":                    template["updated"],
		}

		ids = append(ids, template["id"].(string))
		templates = append(templates, mapping)
	}

	// Set results
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return err
	}
	if err := d.Set("templates", templates); err != nil {
		return err
	}

	return nil
}
