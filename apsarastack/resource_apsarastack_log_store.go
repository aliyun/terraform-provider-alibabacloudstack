package apsarastack

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackLogStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackLogStoreCreate,
		Read:   resourceApsaraStackLogStoreRead,
		Update: resourceApsaraStackLogStoreUpdate,
		Delete: resourceApsaraStackLogStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
					return true
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
			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
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
		},
	}
}

func resourceApsaraStackLogStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	update := false
	if v, ok := d.GetOk("encryption"); ok {
		update = v.(bool)
	}

	//logstore := &sls.LogStore{
	//	Name:          d.Get("name").(string),
	//	TTL:           d.Get("retention_period").(int),
	//	ShardCount:    d.Get("shard_count").(int),
	//	WebTracking:   d.Get("enable_web_tracking").(bool),
	//	AutoSplit:     d.Get("auto_split").(bool),
	//	MaxSplitShard: d.Get("max_split_shard_count").(int),
	//	AppendMeta:    d.Get("append_meta").(bool),
	//}
	var requestinfo *sls.Client
	if update {
		logstore := &sls.LogStore{
			Name:          d.Get("name").(string),
			TTL:           d.Get("retention_period").(int),
			ShardCount:    d.Get("shard_count").(int),
			WebTracking:   d.Get("enable_web_tracking").(bool),
			AutoSplit:     d.Get("auto_split").(bool),
			MaxSplitShard: d.Get("max_split_shard_count").(int),
			AppendMeta:    d.Get("append_meta").(bool),
			Encrypt_conf: sls.Encrypt_conf{
				Enable:       true,
				Encrypt_type: "sm4_gcm",
			},
		}
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {

			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				requestinfo = slsClient
				return nil, slsClient.CreateLogStoreV2(d.Get("project").(string), logstore)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			if debugOn() {
				addDebug("CreateLogStoreV2", raw, requestinfo, map[string]interface{}{
					"project":  d.Get("project").(string),
					"logstore": logstore,
				})
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_log_store", "CreateLogStoreV2", ApsaraStackLogGoSdkERROR)
		}
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))
	}
	return resourceApsaraStackLogStoreUpdate(d, meta)
}

//err := resource.Retry(3*time.Minute, func() *resource.RetryError {
//
//	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
//		requestinfo = slsClient
//		return nil, slsClient.CreateLogStoreV2(d.Get("project").(string), logstore)
//	})
//	if err != nil {
//		if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
//			return resource.RetryableError(err)
//		}
//		return resource.NonRetryableError(err)
//	}
//	if debugOn() {
//		addDebug("CreateLogStoreV2", raw, requestinfo, map[string]interface{}{
//			"project":  d.Get("project").(string),
//			"logstore": logstore,
//		})
//	}
//	return nil
//})
//if err != nil {
//	return WrapErrorf(err, DefaultErrorMsg, "apsarastack_log_store", "CreateLogStoreV2", ApsaraStackLogGoSdkERROR)
//}
//d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))
//
//return resourceApsaraStackLogStoreUpdate(d, meta)

func resourceApsaraStackLogStoreRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := logService.DescribeLogStore(d.Id())
	if err != nil {
		if NotFoundError(err) {
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
			if IsExpectedErrors(err, []string{"InternalServerError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("ListShards", shards)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_log_store", "ListShards", ApsaraStackLogGoSdkERROR)
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

	return nil
}

func resourceApsaraStackLogStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	logService := LogService{client}

	if d.IsNewResource() {
		return resourceApsaraStackLogStoreRead(d, meta)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	if d.HasChange("retention_period") {
		update = true
		//d.SetPartial("retention_period")
	}
	if d.HasChange("max_split_shard_count") {
		update = true
		//d.SetPartial("max_split_shard_count")
	}
	if d.HasChange("enable_web_tracking") {
		update = true
		//d.SetPartial("enable_web_tracking")
	}
	if d.HasChange("append_meta") {
		update = true
		//d.SetPartial("append_meta")
	}
	if d.HasChange("auto_split") {
		update = true
		//d.SetPartial("auto_split")
	}

	if update {
		store, err := logService.DescribeLogStore(d.Id())
		if err != nil {
			return WrapError(err)
		}
		store.MaxSplitShard = d.Get("max_split_shard_count").(int)
		store.TTL = d.Get("retention_period").(int)
		store.WebTracking = d.Get("enable_web_tracking").(bool)
		store.AppendMeta = d.Get("append_meta").(bool)
		store.AutoSplit = d.Get("auto_split").(bool)
		var requestInfo *sls.Client
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.UpdateLogStoreV2(parts[0], store)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateLogStoreV2", ApsaraStackLogGoSdkERROR)
		}
		if debugOn() {
			addDebug("UpdateLogStoreV2", raw, requestInfo, map[string]interface{}{
				"project":  parts[0],
				"logstore": store,
			})
		}
	}
	d.Partial(false)

	return resourceApsaraStackLogStoreRead(d, meta)
}

func resourceApsaraStackLogStoreDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}
