package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackUniversalDnsRecord() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME", "MX", "TXT", "PTR", "SRV", "NAPTR", "CAA", "NS"}, false),
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"lba_strategy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_RR", "RATIO"}, false),
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rdatas": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"lba_weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"line_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackUniversalDnsRecordCreate, resourceAlibabacloudStackUniversalDnsRecordRead, resourceAlibabacloudStackUniversalDnsRecordUpdate, resourceAlibabacloudStackUniversalDnsRecordDelete)
	return resource
}

func resourceAlibabacloudStackUniversalDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request parameters
	reqBody := make(map[string]interface{})
	reqBody["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("name"); ok {
		reqBody["Name"] = v
	}
	reqBody["Type"] = d.Get("type")
	reqBody["Ttl"] = d.Get("ttl")
	reqBody["LbaStrategy"] = d.Get("lba_strategy")
	if v, ok := d.GetOk("remark"); ok {
		reqBody["Remark"] = v
	}

	// Handle rdatas as indexed list in the format RDatas.index.Key
	rdatas := d.Get("rdatas").([]interface{})
	for i, item := range rdatas {
		rdata := item.(map[string]interface{})
		index := i + 1 // index starts from 1
		if value, ok := rdata["value"]; ok {
			key := fmt.Sprintf("RDatas.%d.Value", index)
			reqBody[key] = value
		}
		if lbaWeight, ok := rdata["lba_weight"]; ok && lbaWeight != nil && d.Get("lba_strategy").(string) == "RATIO" {
			key := fmt.Sprintf("RDatas.%d.LbaWeight", index)
			reqBody[key] = lbaWeight
		}
	}

	// Handle line_ids as a JSON array string
	lineIds := d.Get("line_ids").([]interface{})
	lineIdStrings := make([]string, len(lineIds))
	for i, id := range lineIds {
		lineIdStrings[i] = id.(string)
	}
	lineIdsBytes, _ := json.Marshal(lineIdStrings)
	reqBody["LineIds"] = string(lineIdsBytes)

	// Call the API to create the resource
	resp, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "AddUniversalZoneRecord", "", nil, nil, reqBody)
	if err != nil {
		return err
	}

	// Extract record ID from response
	recordId, ok := resp["Id"].(string)
	if !ok || recordId == "" {
		return fmt.Errorf("failed to retrieve record ID from response")
	}

	// Construct resource ID using {ZoneId}:{Id}
	zoneId := d.Get("zone_id").(string)
	resourceId := fmt.Sprintf("%s:%s", zoneId, recordId)

	// Set the temporary ID
	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackUniversalDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	service := UniversalDnsService{client}

	object, err := service.DescribeUniversalDNSRecord(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("zone_id", object["ZoneId"])
	d.Set("name", object["Name"])
	d.Set("type", object["Type"])
	d.Set("ttl", object["Ttl"])
	d.Set("lba_strategy", object["LbaStrategy"])
	d.Set("remark", object["Remark"])

	if lineIds, ok := object["LineIds"].([]interface{}); ok {
		lines := make([]string, 0, len(lineIds))
		for _, v := range lineIds {
			lines = append(lines, v.(string))
		}
		d.Set("line_ids", lines)
	}

	if rdatas, ok := object["RDatas"].([]interface{}); ok {
		rdList := make([]map[string]interface{}, 0, len(rdatas))
		for _, v := range rdatas {
			rdata := v.(map[string]interface{})
			rd := map[string]interface{}{
				"value":      rdata["Value"],
				"lba_weight": rdata["LbaWeight"],
			}
			rdList = append(rdList, rd)
		}
		d.Set("rdatas", rdList)
	}

	return nil
}

func resourceAlibabacloudStackUniversalDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// If this is a new resource, do nothing and return early
	if d.IsNewResource() {
		return nil
	}

	// Handle remark update separately using dedicated API
	if d.HasChange("remark") {
		request := make(map[string]interface{})
		recordId := strings.Split(d.Id(), ":")[1] // Extract record ID from composite ID
		request["Id"] = recordId
		request["Remark"] = d.Get("remark").(string)
		request["ZoneId"] = d.Get("zone_id").(string)

		_, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "UpdateUniversalZoneRemark", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_universal_dns_record", "UpdateUniversalZoneRemark", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	// Handle updates for all other fields using the main update API
	if d.HasChanges("type", "ttl", "lba_strategy", "rdatas", "line_ids") {
		request := make(map[string]interface{})
		recordId := strings.Split(d.Id(), ":")[1] // Extract record ID from composite ID
		request["Id"] = recordId
		request["Name"] = d.Get("name").(string)
		request["Type"] = d.Get("type").(string)
		request["Ttl"] = d.Get("ttl").(int)
		request["LbaStrategy"] = d.Get("lba_strategy").(string)
		request["Remark"] = d.Get("remark").(string)
		request["ZoneId"] = d.Get("zone_id").(string)

		// Process rdatas with indexed keys
		rdatas := d.Get("rdatas").([]interface{})
		for i, rdata := range rdatas {
			rdataMap := rdata.(map[string]interface{})
			index := i + 1
			request[fmt.Sprintf("RDatas.%d.Value", index)] = rdataMap["value"].(string)
			if lbaWeight, ok := rdataMap["lba_weight"]; ok {
				request[fmt.Sprintf("RDatas.%d.LbaWeight", index)] = lbaWeight
			}
		}

		// Process line_ids as indexed list
		lineIds := d.Get("line_ids").([]interface{})
		lineIdStrings := make([]string, len(lineIds))
		for i, id := range lineIds {
			lineIdStrings[i] = id.(string)
		}
		lineIdsBytes, _ := json.Marshal(lineIdStrings)
		request["LineIds"] = string(lineIdsBytes)

		_, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "UpdateUniversalZoneRecord", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_universal_dns_record", "UpdateUniversalZoneRecord", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackUniversalDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid resource ID format")
	}
	zoneId := parts[0]
	recordId := parts[1]
	reqBody := map[string]interface{}{
		"Id":     recordId,
		"ZoneId": zoneId,
	}

	_, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "DeleteUniversalZoneRecord", "", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteUniversalZoneRecord", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
