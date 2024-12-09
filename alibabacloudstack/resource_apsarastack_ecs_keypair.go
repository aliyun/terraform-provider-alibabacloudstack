package alibabacloudstack

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"os"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackKeyPairCreate,
		Read:   resourceAlibabacloudStackKeyPairRead,
		Update: resourceAlibabacloudStackKeyPairUpdate,
		Delete: resourceAlibabacloudStackKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				Deprecated:   "Field 'key_name' is deprecated and will be removed in a future release. Please use new field 'key_pair_name' instead.",
				ConflictsWith: []string{"key_pair_name"},
			},
			"key_pair_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"key_name"},
			},
			"key_name_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return strings.TrimSpace(v.(string))
					default:
						return ""
					}
				},
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"finger_print": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var keyName string
	if v, ok := connectivity.GetResourceDataOk(d, "key_pair_name", "key_name"); ok {
		keyName = v.(string)
	} else if v, ok := d.GetOk("key_name_prefix"); ok {
		keyName = resource.PrefixedUniqueId(v.(string))
	} else {
		keyName = resource.UniqueId()
	}

	if publicKey, ok := d.GetOk("public_key"); ok {
		request := ecs.CreateImportKeyPairRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.KeyPairName = keyName
		request.PublicKeyBody = publicKey.(string)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ImportKeyPair(request)
		})
		response, ok := raw.(*ecs.ImportKeyPairResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "AlibabacloudStack_key_pair", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		d.SetId(response.KeyPairName)
	} else {
		request := ecs.CreateCreateKeyPairRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.KeyPairName = keyName
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CreateKeyPair(request)
		})
		response, ok := raw.(*ecs.CreateKeyPairResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "AlibabacloudStack_key_pair", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		d.SetId(response.KeyPairName)
		if file, ok := d.GetOk("key_file"); ok {
			ioutil.WriteFile(file.(string), []byte(response.PrivateKeyBody), 0600)
			os.Chmod(file.(string), 0400)
		}
	}

	return resourceAlibabacloudStackKeyPairUpdate(d, meta)
}

func resourceAlibabacloudStackKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	err := setTags(client, TagResourceKeypair, d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return resourceAlibabacloudStackKeyPairRead(d, meta)
}

func resourceAlibabacloudStackKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	keyPair, err := ecsService.DescribeKeyPair(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	connectivity.SetResourceData(d, keyPair.KeyPairName, "key_pair_name", "key_name")
	d.Set("finger_print", keyPair.KeyPairFingerPrint)
	tags := keyPair.Tags.Tag
	if len(tags) > 0 {
		err = d.Set("tags", ecsService.tagsToMap(tags))
	}
	return nil
}

func resourceAlibabacloudStackKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteKeyPairsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.KeyPairNames = convertListToJsonString(append(make([]interface{}, 0, 1), d.Id()))

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteKeyPairs(request)
		})
		response, ok := raw.(*ecs.DeleteKeyPairsResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
				return nil
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return errmsgs.WrapError(ecsService.WaitForKeyPair(d.Id(), Deleted, DefaultTimeoutMedium))
}
