package apsarastack

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"os"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackKeyPairCreate,
		Read:   resourceApsaraStackKeyPairRead,
		Update: resourceApsaraStackKeyPairUpdate,
		Delete: resourceApsaraStackKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
				//	Optional:      true,
				//	Computed:      true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				//ConflictsWith: []string{"key_name_prefix"},
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

func resourceApsaraStackKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	var keyName string
	if v, ok := d.GetOk("key_name"); ok {
		keyName = v.(string)
	} else if v, ok := d.GetOk("key_name_prefix"); ok {
		keyName = resource.PrefixedUniqueId(v.(string))
	} else {
		keyName = resource.UniqueId()
	}

	if publicKey, ok := d.GetOk("public_key"); ok {
		request := ecs.CreateImportKeyPairRequest()
		request.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{
			"RegionId":   client.RegionId,
			"Department": client.Department,
		}
		request.QueryParams = map[string]string{
			"AccessKeyId":     client.AccessKey,
			"AccessKeySecret": client.SecretKey,
			"SecurityToken":   client.Config.SecurityToken,
			"Product":         "ecs",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup}
		request.KeyPairName = keyName
		request.PublicKeyBody = publicKey.(string)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ImportKeyPair(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_key_pair", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		object, _ := raw.(*ecs.ImportKeyPairResponse)
		d.SetId(object.KeyPairName)
	} else {
		request := ecs.CreateCreateKeyPairRequest()
		request.RegionId = client.RegionId

		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeyId":     client.AccessKey,
			"AccessKeySecret": client.SecretKey,
			"SecurityToken":   client.Config.SecurityToken,
			"Product":         "ecs",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
		}
		request.KeyPairName = keyName
		fmt.Println(request)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CreateKeyPair(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_key_pair", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		keyPair, _ := raw.(*ecs.CreateKeyPairResponse)
		d.SetId(keyPair.KeyPairName)
		if file, ok := d.GetOk("key_file"); ok {
			ioutil.WriteFile(file.(string), []byte(keyPair.PrivateKeyBody), 0600)
			os.Chmod(file.(string), 0400)
		}
	}

	return resourceApsaraStackKeyPairUpdate(d, meta)
}
func resourceApsaraStackKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	err := setTags(client, TagResourceKeypair, d)
	if err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackKeyPairRead(d, meta)
}

func resourceApsaraStackKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	keyPair, err := ecsService.DescribeKeyPair(d.Id())
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("key_name", keyPair.KeyPairName)
	d.Set("finger_print", keyPair.KeyPairFingerPrint)
	tags := keyPair.Tags.Tag
	if len(tags) > 0 {
		err = d.Set("tags", ecsService.tagsToMap(tags))
	}
	return nil
}

func resourceApsaraStackKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteKeyPairsRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeyId":     client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"SecurityToken":   client.Config.SecurityToken,
		"Product":         "ecs",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup}
	request.KeyPairNames = convertListToJsonString(append(make([]interface{}, 0, 1), d.Id()))

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteKeyPairs(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(ecsService.WaitForKeyPair(d.Id(), Deleted, DefaultTimeoutMedium))
}
