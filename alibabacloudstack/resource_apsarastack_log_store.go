package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLogStore() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntBetween(1, 3650),
			},
			"shard_count": {
				Type:     schema.TypeInt,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "" {
						return false
					}
					return old == new
				},
			},
			"shards": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begin_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"auto_split": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"max_split_shard_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(1, 64),
			},
			"append_meta": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enable_web_tracking": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"encrypt_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "sm4_gcm",
				ValidateFunc: validation.StringInSlice([]string{"sm4_gcm", "aes_gcm"}, false),
			},
			"cmk_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackLogStoreCreate, resourceAlibabacloudStackLogStoreRead, resourceAlibabacloudStackLogStoreUpdate, resourceAlibabacloudStackLogStoreDelete)
	return resource
}

func resourceAlibabacloudStackLogStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slsservice := LogService{client}
	project := d.Get("project").(string)
	slsClient, err := slsservice.GetSlsDataClient(project)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	logstore := &sls.LogStore{
		Name:          d.Get("name").(string),
		TTL:           d.Get("retention_period").(int),
		ShardCount:    d.Get("shard_count").(int),
		WebTracking:   d.Get("enable_web_tracking").(bool),
		AutoSplit:     d.Get("auto_split").(bool),
		MaxSplitShard: d.Get("max_split_shard_count").(int),
		AppendMeta:    d.Get("append_meta").(bool),
	}

	if v, ok := d.GetOk("encryption"); ok && v.(bool) {
		logstore.EncryptConf = &sls.EncryptConf{
			Enable:      v.(bool),
			EncryptType: d.Get("encrypt_type").(string),
		}
		cmkKeyId, hasCmkKeyId := d.GetOk("cmk_key_id")
		arn, hasArn := d.GetOk("arn")
		if hasCmkKeyId && cmkKeyId.(string) != "" {
			if logstore.EncryptConf.UserCmkInfo == nil {
				logstore.EncryptConf.UserCmkInfo = &sls.EncryptUserCmkConf{}
			}
			logstore.EncryptConf.UserCmkInfo.CmkKeyId = cmkKeyId.(string)
			logstore.EncryptConf.UserCmkInfo.RegionId = client.RegionId
			if hasArn && arn.(string) != "" {
				logstore.EncryptConf.UserCmkInfo.Arn = arn.(string)
			} else {
				request := ram.CreateGetRoleRequest()
				client.InitRpcRequest(*request.RpcRequest)
				request.Scheme = "https"
				request.RoleName = "AliyunLogAccessingKMSRole"
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.GetRole(request)
				})
				if err != nil {
					if errmsgs.IsExpectedErrors(err, "EntityNotExist.Role") {
						return errmsgs.WrapErrorf(err, "cmk_key_id is set but arn is not provided, and default RAM role 'AliyunLogAccessingKMSRole' does not exist. Please either set arn explicitly or create the RAM role")
					}
					return errmsgs.WrapErrorf(err, "Failed to query RAM role 'AliyunLogAccessingKMSRole'")
				}
				response := raw.(*ram.GetRoleResponse)
				logstore.EncryptConf.UserCmkInfo.Arn = response.Role.Arn
			}
		} else if hasArn && arn.(string) != "" {
			if logstore.EncryptConf.UserCmkInfo == nil {
				logstore.EncryptConf.UserCmkInfo = &sls.EncryptUserCmkConf{}
			}
			logstore.EncryptConf.UserCmkInfo.Arn = arn.(string)
			logstore.EncryptConf.UserCmkInfo.RegionId = client.RegionId
		}
	}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := slsClient.CreateLogStoreV2(d.Get("project").(string), logstore)
		log.Printf("[DEBUG] SLS CreateLogStoreV2 %++v", logstore)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_log_store", "CreateLogStoreV2", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	// Wait for shards to be ready after creation to avoid ID-refresh race condition.
	// CreateLogStoreV2 returns before shards are fully distributed, causing Read()
	// to see 1 initial shard instead of the expected shard_count.
	expectedShardCount := d.Get("shard_count").(int)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		stored, readErr := slsClient.GetLogStore(d.Get("project").(string), d.Get("name").(string))
		if readErr != nil {
			if errmsgs.IsExpectedErrors(readErr, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(readErr)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(readErr, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_log_store", "GetLogStore", errmsgs.AlibabacloudStackLogGoSdkERROR, ""))
		}
		if stored.ShardCount < expectedShardCount {
			return resource.RetryableError(fmt.Errorf("logstore shard count not ready: expected %d, got %d", expectedShardCount, stored.ShardCount))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceAlibabacloudStackLogStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	object, err := logService.DescribeLogStore(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project", parts[0])
	d.Set("name", object.Name)
	d.Set("retention_period", object.TTL)
	d.Set("shard_count", object.ShardCount)
	var shards []*sls.Shard
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		shards, err = object.ListShards()
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError") {
				return resource.RetryableError(err)
			}
			errmsg := ""
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_log_store", "ListShards", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg))
		}
		addDebug("ListShards", shards)
		return nil
	})
	if err != nil {
		return err
	}
	var shardList []map[string]interface{}
	for _, s := range shards {
		mapping := map[string]interface{}{
			"id":        s.ShardID,
			"status":    s.Status,
			"begin_key": s.InclusiveBeginKey,
			"end_key":   s.ExclusiveBeginKey,
		}
		shardList = append(shardList, mapping)
	}
	d.Set("shards", shardList)
	d.Set("append_meta", object.AppendMeta)
	d.Set("auto_split", object.AutoSplit)
	d.Set("enable_web_tracking", object.WebTracking)
	d.Set("max_split_shard_count", object.MaxSplitShard)

	if object.EncryptConf != nil {
		d.Set("encryption", object.EncryptConf.Enable)
		d.Set("encrypt_type", object.EncryptConf.EncryptType)
		if object.EncryptConf.UserCmkInfo != nil {
			d.Set("cmk_key_id", object.EncryptConf.UserCmkInfo.CmkKeyId)
			d.Set("arn", object.EncryptConf.UserCmkInfo.Arn)
		}
	}

	return nil
}

func resourceAlibabacloudStackLogStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}

	if d.IsNewResource() {
		return nil
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if d.HasChanges("retention_period", "max_split_shard_count", "enable_web_tracking", "append_meta", "auto_split") {
		store, err := logService.DescribeLogStore(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		store.MaxSplitShard = d.Get("max_split_shard_count").(int)
		store.TTL = d.Get("retention_period").(int)
		store.WebTracking = d.Get("enable_web_tracking").(bool)
		store.AppendMeta = d.Get("append_meta").(bool)
		store.AutoSplit = d.Get("auto_split").(bool)

		slsClient, err := logService.GetSlsDataClient(parts[0])
		if err != nil {
			return errmsgs.WrapError(err)
		}

		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			err := slsClient.UpdateLogStoreV2(parts[0], store)
			if err != nil {
				if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
					return resource.RetryableError(err)
				}
				errmsg := ""
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "UpdateLogStoreV2", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceAlibabacloudStackLogStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	name := d.Get("name").(string)
	project := d.Get("project").(string)

	slsClient, err := logService.GetSlsDataClient(project)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := slsClient.DeleteLogStore(project, name)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteLogStore", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
