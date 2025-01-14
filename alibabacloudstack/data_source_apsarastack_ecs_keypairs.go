package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackKeyPairsRead,

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
			"tags": tagsSchema(),
			"finger_print": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			//Computed value
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
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Resource{Schema: outputInstancesSchema()},
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var regex *regexp.Regexp
	if name, ok := d.GetOk("name_regex"); ok {
		regex = regexp.MustCompile(name.(string))
	}
	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}
	request := ecs.CreateDescribeKeyPairsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	if fingerPrint, ok := d.GetOk("finger_print"); ok {
		request.KeyPairFingerPrint = fingerPrint.(string)
	}

	tags := d.Get("tags").(map[string]interface{})
	if tags != nil && len(tags) > 0 {
		KeyPairsTags := make([]ecs.DescribeKeyPairsTag, 0, len(tags))
		for k, v := range tags {
			imageTag := ecs.DescribeKeyPairsTag{
				Key:   k,
				Value: v.(string),
			}
			KeyPairsTags = append(KeyPairsTags, imageTag)
		}
		request.Tag = &KeyPairsTags
	}
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	var keyPairs []ecs.KeyPair
	keyPairsAttach := make(map[string][]map[string]interface{})

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeKeyPairs(request)
		})
		bresponse, ok := raw.(*ecs.DescribeKeyPairsResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_key_pairs", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request, raw, request.RpcRequest, request)
		if len(bresponse.KeyPairs.KeyPair) < 1 {
			break
		}
		for _, key := range bresponse.KeyPairs.KeyPair {
			if regex != nil && !regex.MatchString(key.KeyPairName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[key.KeyPairName]; !ok {
					continue
				}
			}
			keyPairs = append(keyPairs, key)
			keyPairsAttach[key.KeyPairName] = make([]map[string]interface{}, 0)
		}
		if len(bresponse.KeyPairs.KeyPair) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.PageNumber = page
	}

	describeInstancesRequest := ecs.CreateDescribeInstancesRequest()
	client.InitRpcRequest(*describeInstancesRequest.RpcRequest)
	describeInstancesRequest.PageNumber = requests.NewInteger(1)
	describeInstancesRequest.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(describeInstancesRequest)
		})
		bresponse, ok := raw.(*ecs.DescribeInstancesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_key_pairs", describeInstancesRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(describeInstancesRequest.GetActionName(), raw)
		if bresponse == nil || len(bresponse.Instances.Instance) < 1 {
			break
		}
		for _, inst := range bresponse.Instances.Instance {
			if _, ok := keyPairsAttach[inst.KeyPairName]; ok {
				publicIp := inst.EipAddress.IpAddress
				if publicIp == "" && len(inst.PublicIpAddress.IpAddress) > 0 {
					publicIp = inst.PublicIpAddress.IpAddress[0]
				}
				var privateIp string
				if len(inst.InnerIpAddress.IpAddress) > 0 {
					privateIp = inst.InnerIpAddress.IpAddress[0]
				} else if len(inst.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
					privateIp = inst.VpcAttributes.PrivateIpAddress.IpAddress[0]
				}
				mapping := map[string]interface{}{
					"availability_zone": inst.ZoneId,
					"instance_id":       inst.InstanceId,
					"instance_name":     inst.InstanceName,
					"vswitch_id":        inst.VpcAttributes.VSwitchId,
					"public_ip":         publicIp,
					"private_ip":        privateIp,
				}
				if val, ok := keyPairsAttach[inst.KeyPairName]; ok {
					val = append(val, mapping)
					keyPairsAttach[inst.KeyPairName] = val
				} else {
					keyPairsAttach[inst.KeyPairName] = append(make([]map[string]interface{}, 0, 1), mapping)
				}
			}
		}
		if len(bresponse.Instances.Instance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(describeInstancesRequest.PageNumber)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		describeInstancesRequest.PageNumber = page
	}

	return keyPairsDescriptionAttributes(d, keyPairs, keyPairsAttach, meta)
}

func keyPairsDescriptionAttributes(d *schema.ResourceData, keyPairs []ecs.KeyPair, keyPairsAttach map[string][]map[string]interface{}, meta interface{}) error {
	var names []string
	var ids []string
	var s []map[string]interface{}
	for _, key := range keyPairs {
		mapping := map[string]interface{}{
			"id":           key.KeyPairName,
			"key_name":     key.KeyPairName,
			"finger_print": key.KeyPairFingerPrint,
			"instances":    keyPairsAttach[key.KeyPairName],
			"tags":         ecsTagsToMap(key.Tags.Tag),
		}

		names = append(names, string(key.KeyPairName))
		ids = append(ids, string(key.KeyPairName))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("key_pairs", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
