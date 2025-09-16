package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbClusterDatabase() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"character_set_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"collate": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ctype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackPolardbClusterDatabaseCreate, resourceAlibabacloudStackPolardbClusterDatabaseRead, nil, resourceAlibabacloudStackPolardbClusterDatabaseDelete)
	return resource
}

func resourceAlibabacloudStackPolardbClusterDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	request["DBClusterId"] = d.Get("db_cluster_id")
	request["DBName"] = d.Get("db_name")
	request["CharacterSetName"] = d.Get("character_set_name")

	if v, ok := d.GetOk("db_description"); ok {
		request["DBDescription"] = v
	}

	if v, ok := d.GetOk("collate"); ok {
		request["Collate"] = v
	}

	if v, ok := d.GetOk("ctype"); ok {
		request["Ctype"] = v
	}

	_, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "CreateDatabase", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_polardb_cluster_database", "CreateDatabase", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Generate resource ID: {DBClusterId}:{DBName}
	dbClusterId := d.Get("db_cluster_id").(string)
	dbName := d.Get("db_name").(string)
	resourceId := fmt.Sprintf("%s:%s", dbClusterId, dbName)

	// Set the temporary ID
	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackPolardbClusterDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	object, err := polardbService.DescribePolardbClusterDatabase(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("db_cluster_id", strings.Split(d.Id(), ":")[0])
	d.Set("db_name", object["DBName"])
	d.Set("character_set_name", object["CharacterSetName"])
	d.Set("db_description", object["DBDescription"])
	d.Set("engine", object["Engine"])
	d.Set("db_status", object["DBStatus"])
	return nil
}

func resourceAlibabacloudStackPolardbClusterDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dbClusterId := d.Get("db_cluster_id").(string)
	dbName := d.Get("db_name").(string)

	reqQuery := map[string]interface{}{
		"DBClusterId": dbClusterId,
		"DBName":      dbName,
	}

	_, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DeleteDatabase", "", nil, reqQuery, nil)
	if err != nil {
		return fmt.Errorf("delete PolardbClusterDatabase failed: %v", err)
	}
	return nil
}
