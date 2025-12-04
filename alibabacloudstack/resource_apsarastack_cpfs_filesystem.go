package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCpfsFileSystem() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "CPFS",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CPFS"}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"encrypt_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
				Default:      0,
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"bmcpfs"}, false),
				Default:      "bmcpfs",
			},
			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(20*1024, 1700*1024),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNasFileSystemCreate,
		resourceAlibabacloudStackNasFileSystemRead, resourceAlibabacloudStackNasFileSystemUpdate, resourceAlibabacloudStackNasFileSystemDelete)
	return resource
}
