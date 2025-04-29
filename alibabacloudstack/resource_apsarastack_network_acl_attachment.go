package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackNetworkAclAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{

			"network_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNetworkAclAttachmentCreate,
		resourceAlibabacloudStackNetworkAclAttachmentRead, resourceAlibabacloudStackNetworkAclAttachmentUpdate, resourceAlibabacloudStackNetworkAclAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackNetworkAclAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("network_acl_id").(string) + COLON_SEPARATED + resource.UniqueId())
	return resourceAlibabacloudStackNetworkAclAttachmentUpdate(d, meta)
}

func resourceAlibabacloudStackNetworkAclAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	networkAclId := parts[0]
	vpcResource := []vpc.Resource{}
	for _, e := range d.Get("resources").(*schema.Set).List() {
		resourceId := e.(map[string]interface{})["resource_id"]
		resourceType := e.(map[string]interface{})["resource_type"]
		vpcResource = append(vpcResource, vpc.Resource{
			ResourceId:   resourceId.(string),
			ResourceType: resourceType.(string),
		})
	}
	err = vpcService.DescribeNetworkAclAttachment(networkAclId, vpcResource)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("network_acl_id", networkAclId)
	d.Set("resources", vpcResource)
	return nil
}

func resourceAlibabacloudStackNetworkAclAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	networkAclId := parts[0]
	if d.HasChange("resources") {
		oraw, nraw := d.GetChange("resources")
		o := oraw.(*schema.Set)
		n := nraw.(*schema.Set)
		remove := o.Difference(n).List()
		create := n.Difference(o).List()

		if len(remove) > 0 {
			request := vpc.CreateUnassociateNetworkAclRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.NetworkAclId = networkAclId
			request.ClientToken = buildClientToken(request.GetActionName())
			var resources []vpc.UnassociateNetworkAclResource
			vpcResource := []vpc.Resource{}
			for _, t := range remove {
				s := t.(map[string]interface{})
				var resourceId, resourceType string
				if v, ok := s["resource_id"]; ok {
					resourceId = v.(string)
				}
				if v, ok := s["resource_type"]; ok {
					resourceType = v.(string)
				}
				resources = append(resources, vpc.UnassociateNetworkAclResource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
				vpcResource = append(vpcResource, vpc.Resource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
			}
			request.Resource = &resources
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.UnassociateNetworkAcl(request)
			})
			bresponse, ok := raw.(*vpc.UnassociateNetworkAclResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if err := vpcService.WaitForNetworkAclAttachment(request.NetworkAclId, vpcResource, Available, Timeout5Minute); err != nil {
				return errmsgs.WrapError(err)
			}
		}

		if len(create) > 0 {
			request := vpc.CreateAssociateNetworkAclRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.NetworkAclId = networkAclId
			request.ClientToken = buildClientToken(request.GetActionName())
			var resources []vpc.AssociateNetworkAclResource
			vpcResource := []vpc.Resource{}
			for _, t := range create {
				s := t.(map[string]interface{})
				var resourceId, resourceType string
				if v, ok := s["resource_id"]; ok {
					resourceId = v.(string)
				}
				if v, ok := s["resource_type"]; ok {
					resourceType = v.(string)
				}
				resources = append(resources, vpc.AssociateNetworkAclResource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
				vpcResource = append(vpcResource, vpc.Resource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
			}
			request.Resource = &resources
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.AssociateNetworkAcl(request)
			})
			bresponse, ok := raw.(*vpc.AssociateNetworkAclResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if err := vpcService.WaitForNetworkAclAttachment(request.NetworkAclId, vpcResource, Available, Timeout5Minute); err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackNetworkAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	networkAclId := parts[0]
	resources := []vpc.UnassociateNetworkAclResource{}
	object, err := vpcService.DescribeNetworkAcl(networkAclId)
	vpcResource := []vpc.Resource{}
	request := vpc.CreateUnassociateNetworkAclRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.NetworkAclId = networkAclId
	request.ClientToken = buildClientToken(request.GetActionName())
	res, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
	for _, e := range res {
		item := e.(map[string]interface{})
		resources = append(resources, vpc.UnassociateNetworkAclResource{
			ResourceId:   item["ResourceId"].(string),
			ResourceType: item["ResourceType"].(string),
		})
		vpcResource = append(vpcResource, vpc.Resource{
			ResourceId:   item["ResourceId"].(string),
			ResourceType: item["ResourceType"].(string),
		})
	}
	request.Resource = &resources
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateNetworkAcl(request)
		})
		bresponse, ok := raw.(*vpc.UnassociateNetworkAclResponse)
		//Waiting for unassociate the network acl
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	return vpcService.WaitForNetworkAclAttachment(networkAclId, vpcResource, Deleted, Timeout5Minute)
}