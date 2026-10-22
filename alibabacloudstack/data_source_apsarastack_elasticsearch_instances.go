package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackElasticsearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackElasticsearchRead,

		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// There is an issue with Zone filtering
			// 			"zone_id": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},
			// Computed values
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// Basic instance information
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// Data node configuration

						"data_node_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"data_node_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"data_node_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"data_node_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// Kibana node configuration
						"kibana_node_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// Master node configuration
						"master_node_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"master_node_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"master_node_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"master_node_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// Client node configuration
						"client_node_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"client_node_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// network info
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
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

func dataSourceAlibabacloudStackElasticsearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var response map[string]interface{}
	var err error
	var ids []string
	var descriptions []string
	var instances []map[string]interface{}
	var filteredInstances []map[string]interface{}
	request := make(map[string]interface{})
	if v, ok := d.GetOk("version"); ok && v.(string) != "" {
		request["esVersion"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request["vpcId"] = v.(string)
	}
	// 	if v, ok := d.GetOk("zone_id"); ok && v.(string) != "" {
	// 		request["zoneId"] = v.(string)
	// 	}
	for {
		request["Size"] = requests.NewInteger(PageSizeLarge)
		request["Page"] = 1
		response, err = client.DoTeaRequest("GET", "elasticsearch-k8s", "2017-06-13", "ListInstance", "/openapi/instances", nil, request, nil)
		addDebug("ListInstance", response, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InstanceNotFound") {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListInstance", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", "ListInstance", response))
		}
		r, err := jsonpath.Get("$.Result", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "$.Result", response)
		}

		results := r.([]interface{})

		if len(results) < 1 {
			break
		}

		for _, result := range results {
			object := result.(map[string]interface{})
			mapping := map[string]interface{}{
				"id":          object["instanceId"],
				"version":     object["esVersion"],
				"description": object["description"],
			}
			if dataNode, ok := object["dataNode"].(bool); ok && dataNode {
				if nodeAmount, ok := object["nodeAmount"].(json.Number); ok {
					if v, err := nodeAmount.Int64(); err == nil {
						mapping["data_node_amount"] = int(v)
					} else {
						return errmsgs.WrapError(err)
					}
				}
				if nodeSpec, ok := object["nodeSpec"].(map[string]interface{}); ok {
					mapping["data_node_spec"] = nodeSpec["spec"]
					if diskSize, ok := nodeSpec["disk"].(json.Number); ok {
						if v, err := diskSize.Int64(); err == nil {
							mapping["data_node_disk_size"] = int(v)
						} else {
							return errmsgs.WrapError(err)
						}
					}
					mapping["data_node_disk_type"] = nodeSpec["storageClassName"]
				}
			}

			if haveKibana, ok := object["haveKibana"].(bool); ok && haveKibana {
				if kibanaConfig, ok := object["kibanaConfiguration"].(map[string]interface{}); ok {
					mapping["kibana_node_spec"] = kibanaConfig["spec"]
				}
			}

			if advancedDedicateMaster, ok := object["advancedDedicateMaster"].(bool); ok && advancedDedicateMaster {
				if masterConfiguration, ok := object["masterConfiguration"].(map[string]interface{}); ok {
					if masterAmount, ok := masterConfiguration["amount"].(json.Number); ok {
						if v, err := masterAmount.Int64(); err == nil {
							mapping["master_node_amount"] = int(v)
						} else {
							return errmsgs.WrapError(err)
						}
					}
					mapping["master_node_spec"] = masterConfiguration["spec"]
					if masterDisk, ok := masterConfiguration["disk"].(json.Number); ok {
						if v, err := masterDisk.Int64(); err == nil {
							mapping["master_node_disk_size"] = int(v)
						} else {
							return errmsgs.WrapError(err)
						}
					}
					mapping["master_node_disk_type"] = masterConfiguration["storageClassName"]
				}
			}

			if haveClientNode, ok := object["haveClientNode"].(bool); ok && haveClientNode {
				if clientNodeConfiguration, ok := object["clientNodeConfiguration"].(map[string]interface{}); ok {
					if clientAmount, ok := clientNodeConfiguration["amount"].(json.Number); ok {
						if v, err := clientAmount.Int64(); err == nil {
							mapping["client_node_amount"] = int(v)
						} else {
							return errmsgs.WrapError(err)
						}
					}
					mapping["client_node_spec"] = clientNodeConfiguration["spec"]
				}
			}
			if networkConfig, ok := object["networkConfig"].(map[string]interface{}); ok {
				mapping["vswitch_id"] = networkConfig["vswitchId"]
				mapping["vpc_id"] = networkConfig["vpcId"]
			}
			mapping["status"] = object["status"]

			instances = append(instances, mapping)

		}

		if len(results) < PageSizeLarge {
			break
		}

		request["Page"] = request["Page"].(int) + 1
	}

	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		descriptionRegex = r
	}

	// ids
	idsMap := getIdsStringFilter(d)

	for _, instance := range instances {
		if descriptionRegex != nil && !descriptionRegex.MatchString(instance["description"].(string)) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[instance["id"].(string)]; !ok {
				continue
			}
		}
		filteredInstances = append(filteredInstances, instance)
		ids = append(ids, instance["id"].(string))
		descriptions = append(descriptions, instance["description"].(string))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", filteredInstances); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("descriptions", descriptions); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), filteredInstances); err != nil {
			return err
		}
	}
	return nil
}
