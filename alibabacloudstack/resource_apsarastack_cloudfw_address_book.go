package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCloudfwAddressBook() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ip", "port"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_list": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_add_tag_ecs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_list_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"reference_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"global": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCloudfwAddressBookCreate, resourceAlibabacloudStackCloudfwAddressBookRead, resourceAlibabacloudStackCloudfwAddressBookUpdate, resourceAlibabacloudStackCloudfwAddressBookDelete)
	return resource
}

func resourceAlibabacloudStackCloudfwAddressBookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	addlist := []string{}
	addressList := d.Get("address_list").(*schema.Set).List()

	for _, v := range addressList {
		addlist = append(addlist, v.(string))
	}
	request := map[string]interface{}{
		"SourceCode":  "yundun",
		"GroupName":   d.Get("group_name").(string),
		"GroupType":   d.Get("group_type").(string),
		"Description": d.Get("description").(string),
		"AddressList": strings.Join(addlist, ","),
	}

	response, err := client.DoTeaRequest("POST", "cloudfw", "2017-12-07", "AddAddressBook", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Get the GroupUuid from response
	groupUuid, ok := response["GroupUuid"]
	if !ok {
		return fmt.Errorf("failed to get GroupUuid from response")
	}
	resourceId := fmt.Sprintf("%s:%s", d.Get("group_type").(string), groupUuid)
	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackCloudfwAddressBookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeAddressBook(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("group_name", object["GroupName"])
	d.Set("group_type", object["GroupType"])
	d.Set("description", object["Description"])
	d.Set("address_list_count", object["AddressListCount"])
	d.Set("group_uuid", object["GroupUuid"])
	d.Set("reference_count", object["ReferenceCount"])
	d.Set("global", object["Global"])
	d.Set("auto_add_tag_ecs", object["AutoAddTagEcs"])
	d.Set("address_list", object["AddressList"])

	return nil
}

func resourceAlibabacloudStackCloudfwAddressBookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	// Prepare request parameters
	if d.HasChanges("description", "address_list") {
		addlist := []string{}
		addressList := d.Get("address_list").(*schema.Set).List()

		for _, v := range addressList {
			addlist = append(addlist, v.(string))
		}
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		requestInfo := map[string]interface{}{
			"SourceCode":  "yundun",
			"GroupUuid":   parts[1],
			"GroupName":   d.Get("group_name"),
			"GroupType":   parts[0],
			"Description": d.Get("description"),
			"AddressList": strings.Join(addlist, ","),
		}

		_, err = client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", "ModifyAddressBook", "", nil, requestInfo, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_cloudfw_address_book", "ModifyAddressBook", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackCloudfwAddressBookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	reqQuery := map[string]interface{}{
		"SourceCode": "yundun",
		"GroupUuid":  parts[1],
		"GroupType":  parts[0],
	}

	_, err = client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", "DeleteAddressBook", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_cloudfw_address_book", "DeleteAddressBook", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
