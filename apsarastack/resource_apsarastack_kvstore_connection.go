package apsarastack

import (
	"fmt"
	"log"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackKvstoreConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackKvstoreConnectionCreate,
		Read:   resourceApsaraStackKvstoreConnectionRead,
		Update: resourceApsaraStackKvstoreConnectionUpdate,
		Delete: resourceApsaraStackKvstoreConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceApsaraStackKvstoreConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	r_kvstoreService := KvstoreService{client}

	request := r_kvstore.CreateAllocateInstancePublicConnectionRequest()
	request.ConnectionStringPrefix = d.Get("connection_string_prefix").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = d.Get("port").(string)
	request.RegionId = client.RegionId
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "R-kvstore",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "AllocateInstancePublicConnection",
		"Version":         "2015-01-01",
	}
	raw, err := client.WithRkvClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.AllocateInstancePublicConnection(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_kvstore_connection", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(fmt.Sprintf("%v", request.InstanceId))
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, r_kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceApsaraStackKvstoreConnectionUpdate(d, meta)
}
func resourceApsaraStackKvstoreConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	r_kvstoreService := KvstoreService{client}
	object, err := r_kvstoreService.DescribeKvstoreConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_kvstore_connection r_kvstoreService.DescribeKvstoreConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	for _, instanceNetInfo := range object {
		if instanceNetInfo.DBInstanceNetType == "0" {
			d.Set("connection_string", instanceNetInfo.ConnectionString)
			d.Set("port", instanceNetInfo.Port)
		}
	}
	return nil
}
func resourceApsaraStackKvstoreConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	r_kvstoreService := KvstoreService{client}
	update := false
	request := r_kvstore.CreateModifyDBInstanceConnectionStringRequest()
	request.DBInstanceId = d.Id()

	request.CurrentConnectionString = d.Get("connection_string").(string)
	if !d.IsNewResource() && d.HasChange("connection_string_prefix") {
		update = true
	}
	request.NewConnectionString = d.Get("connection_string_prefix").(string)
	request.IPType = "Public"
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "R-kvstore",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "ModifyDBInstanceConnectionString",
		"Version":         "2015-01-01",
	}
	if !d.IsNewResource() && d.HasChange("port") {
		update = true
		request.Port = d.Get("port").(string)
	}
	if update {
		raw, err := client.WithRkvClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyDBInstanceConnectionString(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, r_kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceApsaraStackKvstoreConnectionRead(d, meta)
}
func resourceApsaraStackKvstoreConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	r_kvstoreService := KvstoreService{client}
	request := r_kvstore.CreateReleaseInstancePublicConnectionRequest()
	request.InstanceId = d.Id()
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "R-kvstore",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "ReleaseInstancePublicConnection",
		"Version":         "2015-01-01",
	}
	request.CurrentConnectionString = d.Get("connection_string").(string)
	raw, err := client.WithRkvClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.ReleaseInstancePublicConnection(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutDelete), 30*time.Second, r_kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
