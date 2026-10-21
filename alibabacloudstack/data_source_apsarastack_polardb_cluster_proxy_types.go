package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPolardbClusterProxyTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbClusterProxyTypesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"db_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "PostgreSQL", "Oracle"}, false),
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"core_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Computed values.
			"proxy_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbClusterProxyTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	filterIds := getIdsStringFilter(d)

	existedId := map[string]string{}
	ids := []string{}
	types := []map[string]interface{}{}

	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}

	// Build base query parameters shared by both resourceType attempts.
	buildReqQuery := func(resourceType string) map[string]interface{} {
		reqQuery := map[string]interface{}{
			"pageStart":    1,
			"pageSize":     500,
			"label":        "true",
			"resourceType": resourceType,
			"status":       "Available",
			"groupFiled":   "proxyClass",
		}
		if v, ok := d.GetOk("db_version"); ok {
			reqQuery["dbVersion"] = v
		}
		if v, ok := d.GetOk("db_type"); ok {
			reqQuery["dbType"] = v
		}
		return reqQuery
	}

	// Try polardb-proxy first (new ASCM versions separate proxy specs from cluster specs).
	// Fall back to POLARDB for backward compatibility with older ASCM versions.
	var response map[string]interface{}
	var err error
	reqQuery := buildReqQuery("polardb-proxy")
	response, err = client.DoTeaRequest("POST", "ascm", "2019-05-10", "GroupCommonSpec", "/ascm/manage/saleconf/commonSpec/group", reqHeader, reqQuery, nil)
	if err != nil || response == nil || len(extractDataSlice(response)) == 0 {
		if err != nil {
			log.Printf("[WARN] GroupCommonSpec with resourceType=polardb-proxy failed: %s, falling back to POLARDB", err)
		} else {
			log.Printf("[WARN] GroupCommonSpec with resourceType=polardb-proxy returned empty data, falling back to POLARDB")
		}
		reqQuery = buildReqQuery("POLARDB")
		response, err = client.DoTeaRequest("POST", "ascm", "2019-05-10", "GroupCommonSpec", "/ascm/manage/saleconf/commonSpec/group", reqHeader, reqQuery, nil)
		if err != nil {
			return err
		}
	}

	for _, v := range extractDataSlice(response) {
		data := v.(map[string]interface{})
		idRaw, ok := data["proxyClass"].(string)
		if !ok || idRaw == "" {
			continue
		}
		if _, exists := filterIds[idRaw]; len(filterIds) > 0 && !exists {
			continue
		}
		if _, exists := existedId[idRaw]; exists {
			continue
		}
		classLabelRaw, ok := data["proxyClassLabel"].(string)
		if !ok {
			classLabelRaw = fmt.Sprintf("%v", data["proxyClassLabel"])
		}
		coreCount := parseProxyCoreCount(classLabelRaw, idRaw)
		if coreCount == 0 {
			log.Printf("[WARN] Unable to determine core_count from proxyClassLabel=%q proxyClass=%q, skipping", classLabelRaw, idRaw)
			continue
		}
		if v, ok := d.GetOk("core_count"); ok && coreCount != v.(int) {
			continue
		}

		types = append(types, map[string]interface{}{
			"id":         idRaw,
			"core_count": coreCount,
		})
		existedId[idRaw] = ""
		ids = append(ids, idRaw)
	}
	sort.SliceStable(types, func(i, j int) bool {
		return types[i]["core_count"].(int) < types[j]["core_count"].(int)
	})

	d.Set("ids", ids)
	d.Set("proxy_classes", types)
	d.SetId(dataResourceIdHash(ids))

	return nil
}

// parseProxyCoreCount extracts the CPU core count from the proxyClassLabel.
// Old ASCM format: label = "2 C" -> first token is the core count.
// New ASCM format: label = "polar.maxscale.g2.medium.c" (same as proxyClass) ->
// derive core count from the size tier in the spec ID.
func parseProxyCoreCount(label, proxyClass string) int {
	// Try old format first: "<number> C" or just "<number>"
	firstToken := strings.Split(strings.TrimSpace(label), " ")[0]
	if n, err := strconv.Atoi(firstToken); err == nil && n > 0 {
		return n
	}

	// New format: derive from size tier in spec ID
	// Patterns:
	//   polar.maxscale.g2.medium.c -> generation "g2" at index 2, size "medium" at index 3
	//   polar.o.maxscale.x2.medium -> generation "x2" at index 3, size "medium" at index 4
	// Find the generation identifier (g2/g4/g8/x2/etc.) and take the next part as size tier.
	sizeTier := ""
	parts := strings.Split(proxyClass, ".")
	genPattern := regexp.MustCompile(`^[gx]\d+$`)
	for i, p := range parts {
		if genPattern.MatchString(p) && i+1 < len(parts) {
			sizeTier = parts[i+1]
			break
		}
	}

	sizeToCore := map[string]int{
		"small":    1,
		"medium":   2,
		"large":    4,
		"xlarge":   8,
		"2xlarge":  16,
		"3xlarge":  24,
		"4xlarge":  32,
		"8xlarge":  64,
		"12xlarge": 96,
		"16xlarge": 128,
	}
	if n, ok := sizeToCore[sizeTier]; ok {
		return n
	}

	return 0
}

// extractDataSlice safely extracts the "data" field from response as []interface{}.
func extractDataSlice(response map[string]interface{}) []interface{} {
	if response == nil {
		return nil
	}
	data, ok := response["data"]
	if !ok || data == nil {
		return nil
	}
	if slice, ok := data.([]interface{}); ok {
		return slice
	}
	return nil
}
