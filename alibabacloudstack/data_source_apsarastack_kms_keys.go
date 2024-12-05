package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackKmsKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackKmsKeysRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},

			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				// must contain a valid status, expected Enabled, Disabled, PendingDeletion
				ValidateFunc: validation.StringInSlice([]string{
					string(EnabledStatus),
					string(DisabledStatus),
					string(Pending),
				}, false),
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			//Computed value
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackKmsKeysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := kms.CreateListKeysRequest()
	client.InitRpcRequest(*request.RpcRequest)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, i := range v.([]interface{}) {
			idsMap[i.(string)] = i.(string)
		}
	}

	var s []map[string]interface{}
	var ids []string
	var keyIds []string

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for true {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.ListKeys(request)
		})
		response, ok := raw.(*kms.ListKeysResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kms_keys", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		for _, key := range response.Keys.Key {
			if len(idsMap) > 0 {
				if _, ok := idsMap[key.KeyId]; ok {
					keyIds = append(keyIds, key.KeyId)
					continue
				}
			} else {
				keyIds = append(keyIds, key.KeyId)
				continue
			}
		}
		if len(response.Keys.Key) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.PageNumber = page
	}

	descriptionRegex, ok := d.GetOk("description_regex")
	var r *regexp.Regexp
	if ok && descriptionRegex.(string) != "" {
		r = regexp.MustCompile(descriptionRegex.(string))
	}
	status, statusOk := d.GetOk("status")
	for _, k := range keyIds {

		request := kms.CreateDescribeKeyRequest()
		client.InitRpcRequest(*request.RpcRequest)

		request.KeyId = k
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.DescribeKey(request)
		})
		response, ok := raw.(*kms.DescribeKeyResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, k, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		if r != nil && !r.MatchString(response.KeyMetadata.Description) {
			continue
		}
		if statusOk && status != "" && status != response.KeyMetadata.KeyState {
			continue
		}
		mapping := map[string]interface{}{
			"id":              response.KeyMetadata.KeyId,
			"arn":             response.KeyMetadata.Arn,
			"description":     response.KeyMetadata.Description,
			"status":          response.KeyMetadata.KeyState,
			"creation_date":   response.KeyMetadata.CreationDate,
			"delete_date":     response.KeyMetadata.DeleteDate,
			"creator":         response.KeyMetadata.Creator,
		}

		s = append(s, mapping)
		ids = append(ids, response.KeyMetadata.KeyId)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("keys", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
