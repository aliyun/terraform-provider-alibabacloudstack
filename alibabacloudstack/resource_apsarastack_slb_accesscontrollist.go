package alibabacloudstack

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlbAcl() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'acl_name' instead.",
				ConflictsWith: []string{"acl_name"},
			},
			"acl_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
				Deprecated:   "Field 'ip_version' is deprecated and will be removed in a future release. Please use new field 'address_ip_version' instead.",
				ConflictsWith: []string{"address_ip_version"},
			},
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
				ConflictsWith: []string{"ip_version"},
			},
			"entry_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 300,
				MinItems: 0,
			},
			//"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackSlbAclCreate, resourceAlibabacloudStackSlbAclRead, resourceAlibabacloudStackSlbAclUpdate, resourceAlibabacloudStackSlbAclDelete)
	return resource
}

func resourceAlibabacloudStackSlbAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := slb.CreateCreateAccessControlListRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AclName = strings.TrimSpace(connectivity.GetResourceData(d, "acl_name", "name").(string))
	if err := errmsgs.CheckEmpty(request.AclName, schema.TypeString, "acl_name", "name"); err != nil {
		return errmsgs.WrapError(err)
	}
	if v, ok := connectivity.GetResourceDataOk(d, "address_ip_version", "ip_version"); ok {
		request.AddressIPVersion = v.(string)
	} else {
		request.AddressIPVersion = "ipv4"
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateAccessControlList(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*slb.CreateAccessControlListResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_acl", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateAccessControlListResponse)

	d.SetId(response.AclId)
	return nil
}

func resourceAlibabacloudStackSlbAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

// 	slb-acl的tags由ascm实现，3162相关接口没有注册到pop
// 	tags, err := slbService.DescribeTags(d.Id(), nil, TagResourceAcl)
// 	if err != nil {
// 		return errmsgs.WrapError(err)
// 	}
// 	d.Set("tags", slbService.tagsToMap(tags))

	object, err := slbService.DescribeSlbAcl(d.Id())
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"AclNotExist"}) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	connectivity.SetResourceData(d, object.AclName, "acl_name", "name")
	connectivity.SetResourceData(d, object.AddressIPVersion, "address_ip_version", "ip_version")

	if err := d.Set("entry_list", slbService.FlattenSlbAclEntryMappings(object.AclEntrys.AclEntry)); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackSlbAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChanges("name", "acl_name") {
		request := slb.CreateSetAccessControlListAttributeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.AclId = d.Id()
		request.AclName = connectivity.GetResourceData(d, "acl_name", "name").(string)
		if err := errmsgs.CheckEmpty(request.AclName, schema.TypeString, "acl_name", "name"); err != nil {
			return errmsgs.WrapError(err)
		}

		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetAccessControlListAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*slb.SetAccessControlListAttributeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("entry_list") {
		o, n := d.GetChange("entry_list")
		oe := o.(*schema.Set)
		ne := n.(*schema.Set)
		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		if len(remove) > 0 {
			if err := slbService.SlbRemoveAccessControlListEntry(remove, d.Id()); err != nil {
				return errmsgs.WrapError(err)
			}
		}

		if len(add) > 0 {
			if err := slbService.SlbAddAccessControlListEntry(add, d.Id()); err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}

	d.Partial(false)

	return nil
}

func resourceAlibabacloudStackSlbAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	request := slb.CreateDeleteAccessControlListRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AclId = d.Id()

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteAccessControlList(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"AclInUsed"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*slb.DeleteAccessControlListResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if !errmsgs.IsExpectedErrors(err, []string{"AclNotExist"}) {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return errmsgs.WrapError(slbService.WaitForSlbAcl(d.Id(), Deleted, DefaultTimeoutMedium))
}
