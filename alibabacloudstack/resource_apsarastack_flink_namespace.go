package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackFlinkNamespace() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"cu": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"cpu_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Intel"}, false),
			},
			"owner_uid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackFlinkNamespaceCreate, resourceAlibabacloudStackFlinkNamespaceRead,
		resourceAlibabacloudStackFlinkNamespaceUpdate, resourceAlibabacloudStackFlinkNamespaceDelete)
	return resource
}

func resourceAlibabacloudStackFlinkNamespaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	namespaceName := d.Get("name").(string)
	var userId string
	if v, ok := d.GetOk("owner_uid"); ok {
		userId = v.(string)
	} else if userId, err = client.AccountId(); err != nil {
		return errmsgs.WrapError(fmt.Errorf("Get user Owner Id Failed %v", err))
	}
	requestQuery := map[string]interface{}{
		"name":       namespaceName,
		"allocateCu": d.Get("cu").(int),
		"cpuBrand":   d.Get("cpu_type").(string),
		"ownerUid":   userId,
	}
	var response map[string]interface{}
	response, err = client.DoTeaRequest("POST", "ververica", "2020-05-01", "CreateNamespace", "/flink/namespace/create", nil, requestQuery, nil)
	if err != nil {
		return err
	}
	log.Printf("response for create %v", response)

	d.SetId(namespaceName)

	return nil
}

func resourceAlibabacloudStackFlinkNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("cu") {
		requestQuery := map[string]interface{}{
			"name":       d.Id(),
			"allocateCu": d.Get("cu").(int),
		}
		response, err := client.DoTeaRequest("PUT", "ververica", "2020-05-01", "UpdateNamespace", "/flink/namespace/update", nil, requestQuery, nil)
		if err != nil {
			return err
		}
		log.Printf("response for create %v", response)
	}

	time.Sleep(3 * time.Second)

	return nil
}

func resourceAlibabacloudStackFlinkNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	flinkService := FlinkService{client}

	object, err := flinkService.DescribeFlinkNamespace(d.Id())
	if err != nil || object == nil {
		if object == nil || errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object["Name"].(string))
	cu, err := strconv.Atoi(object["GuaranteeQuota"].(string))
	if err != nil {
		return err
	}
	d.Set("cu", cu)
	d.Set("cpu_type", object["CpuBrand"].(string))
	d.Set("owner_uid", object["Uid"].(string))

	return nil
}

func resourceAlibabacloudStackFlinkNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	requestQuery := map[string]interface{}{
		"name": d.Id(),
	}
	response, err := client.DoTeaRequest("DELETE", "ververica", "2020-05-01", "DeleteNamespace", "/flink/namespace/delete", nil, requestQuery, nil)
	if err != nil {
		return err
	}
	log.Printf("response for delete %v", response)
	return nil
}
