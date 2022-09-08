package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackNetworkAclAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackNetworkAclAttachmentCreate,
		Read:   resourceAlibabacloudStackNetworkAclAttachmentRead,
		Update: resourceAlibabacloudStackNetworkAclAttachmentUpdate,
		Delete: resourceAlibabacloudStackNetworkAclAttachmentDelete,

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
		return WrapError(err)
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
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		return WrapError(err)
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
			request.Headers["x-ascm-product-name"] = "Vpc"
			request.Headers["x-acs-organizationId"] = client.Department
			request.RegionId = client.RegionId
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
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if err := vpcService.WaitForNetworkAclAttachment(request.NetworkAclId, vpcResource, Available, Timeout5Minute); err != nil {
				return WrapError(err)
			}
		}

		if len(create) > 0 {
			request := vpc.CreateAssociateNetworkAclRequest()
			request.Headers["x-ascm-product-name"] = "Vpc"
			request.Headers["x-acs-organizationId"] = client.Department
			request.RegionId = client.RegionId
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
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if err := vpcService.WaitForNetworkAclAttachment(request.NetworkAclId, vpcResource, Available, Timeout5Minute); err != nil {
				return WrapError(err)
			}
		}
		//d.SetPartial("resources")
	}

	return resourceAlibabacloudStackNetworkAclAttachmentRead(d, meta)
}

func resourceAlibabacloudStackNetworkAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	networkAclId := parts[0]

	resources := []vpc.UnassociateNetworkAclResource{}
	object, err := vpcService.DescribeNetworkAcl(networkAclId)
	vpcResource := []vpc.Resource{}
	request := vpc.CreateUnassociateNetworkAclRequest()
	request.Headers["x-ascm-product-name"] = "Vpc"
	request.Headers["x-acs-organizationId"] = client.Department
	request.RegionId = client.RegionId
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
		//Waiting for unassociate the network acl
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	return vpcService.WaitForNetworkAclAttachment(networkAclId, vpcResource, Deleted, Timeout5Minute)
}
