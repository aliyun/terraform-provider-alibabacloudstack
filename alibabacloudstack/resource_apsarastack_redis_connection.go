package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackKvstoreConnection() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackKvstoreConnectionCreate, resourceAlibabacloudStackKvstoreConnectionRead, resourceAlibabacloudStackKvstoreConnectionUpdate, resourceAlibabacloudStackKvstoreConnectionDelete)
	return resource
}

func resourceAlibabacloudStackKvstoreConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	r_kvstoreService := KvstoreService{client}

	request := r_kvstore.CreateAllocateInstancePublicConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ConnectionStringPrefix = d.Get("connection_string_prefix").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = d.Get("port").(string)

	raw, err := client.WithRkvClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.AllocateInstancePublicConnection(request)
	})
	response, ok := raw.(*r_kvstore.AllocateInstancePublicConnectionResponse)
	addDebug(request.GetActionName(), raw)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kvstore_connection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	d.SetId(fmt.Sprintf("%v", request.InstanceId))
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, r_kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackKvstoreConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	r_kvstoreService := KvstoreService{client}
	object, err := r_kvstoreService.DescribeKvstoreConnection(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_kvstore_connection r_kvstoreService.DescribeKvstoreConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
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

func resourceAlibabacloudStackKvstoreConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	r_kvstoreService := KvstoreService{client}
	update := false
	request := r_kvstore.CreateModifyDBInstanceConnectionStringRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()
	request.CurrentConnectionString = d.Get("connection_string").(string)

	if !d.IsNewResource() && d.HasChange("connection_string_prefix") {
		update = true
	}
	request.NewConnectionString = d.Get("connection_string_prefix").(string)
	request.IPType = "Public"

	if !d.IsNewResource() && d.HasChange("port") {
		update = true
		request.Port = d.Get("port").(string)
	}

	if update {
		raw, err := client.WithRkvClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyDBInstanceConnectionString(request)
		})
		response, ok := raw.(*r_kvstore.ModifyDBInstanceConnectionStringResponse)
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, r_kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	return nil
}

func resourceAlibabacloudStackKvstoreConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	r_kvstoreService := KvstoreService{client}
	request := r_kvstore.CreateReleaseInstancePublicConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Id()
	request.CurrentConnectionString = d.Get("connection_string").(string)

	raw, err := client.WithRkvClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.ReleaseInstancePublicConnection(request)
	})
	response, ok := raw.(*r_kvstore.ReleaseInstancePublicConnectionResponse)
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutDelete), 30*time.Second, r_kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}