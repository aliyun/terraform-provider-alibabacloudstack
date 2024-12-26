package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmRamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmRamPolicyCreate,
		Read:   resourceAlibabacloudStackAscmRamPolicyRead,
		Update: resourceAlibabacloudStackAscmRamPolicyUpdate,
		Delete: resourceAlibabacloudStackAscmRamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"policy_document": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ram_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackAscmRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	ascmService := AscmService{client}
	name := d.Get("name").(string)

	policyDoc := d.Get("policy_document").(string)
	description := d.Get("description").(string)
	//resp := RamPolicies{}
	check, err := ascmService.DescribeAscmRamPolicy(name)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ram_policy", "policy alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(check.Data) == 0 {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateRAMPolicy", "/ascm/auth/role/createRAMPolicy")
		request.QueryParams["policyName"] = name
		request.QueryParams["description"] = description
		request.QueryParams["policyDocument"] = policyDoc

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" rsponse of CreateRAMPolicy : %s", raw)
		if err != nil {
			errmsg := ""
			if raw != nil {
				bresponse, ok := raw.(*responses.CommonResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy", "CreateRAMPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		if bresponse, ok := raw.(*responses.CommonResponse); ok {
			if bresponse.GetHttpStatus() != 200 {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ram_policy", "CreateRAMPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			addDebug("CreateRAMPolicy", raw, requestInfo, bresponse.GetHttpContentString())
		}
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmRamPolicy(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	d.SetId(check.Data[0].PolicyName + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))

	return resourceAlibabacloudStackAscmRamPolicyUpdate(d, meta)
}

func resourceAlibabacloudStackAscmRamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	response, err := ascmService.DescribeAscmRamPolicy(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		// Handle exceptions
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", did[0])
	d.Set("ram_id", did[1])
	d.Set("description", response.Data[0].Description)
	d.Set("policy_document", response.Data[0].PolicyDocument)
	return nil
}

func resourceAlibabacloudStackAscmRamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	attributeUpdate := false
	check, err := ascmService.DescribeAscmRamPolicy(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsInstanceExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	var name, description, policydoc string

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data[0].PolicyName = name
		check.Data[0].NewPolicyName = name
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data[0].PolicyName = name
	}
	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		}
		check.Data[0].Description = description
		check.Data[0].NewDescription = description
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		}
		check.Data[0].Description = description
	}
	if d.HasChange("policy_document") {
		if v, ok := d.GetOk("policy_document"); ok {
			policydoc = v.(string)
		}
		check.Data[0].PolicyDocument = policydoc
		check.Data[0].NewPolicyDocument = policydoc
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("policy_document"); ok {
			policydoc = v.(string)
		}
		check.Data[0].PolicyDocument = policydoc
	}
	if attributeUpdate {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateRAMPolicy", "/ascm/auth/role/updateRAMPolicy")
		request.QueryParams["RamPolicyId"] = did[1]
		request.QueryParams["NewPolicyName"] = name
		request.QueryParams["NewDescription"] = description
		request.QueryParams["NewPolicyDocument"] = policydoc

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw UpdateRAMPolicy : %s", raw)

		if err != nil {
			errmsg := ""
			if raw != nil {
				bresponse, ok := raw.(*responses.CommonResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy", "UpdateRAMPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request)
		log.Printf("total QueryParams and rampolicy %v %v", request.GetQueryParams(), name)
	}
	d.SetId(name + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))

	return resourceAlibabacloudStackAscmRamPolicyRead(d, meta)
}

func resourceAlibabacloudStackAscmRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmRamPolicy(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsPolicyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsPolicyExist", check, requestInfo, map[string]string{"policyName": did[0]})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRAMPolicy", "/ascm/auth/role/removeRAMPolicy")
		request.QueryParams["ramPolicyId"] = did[1]

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = ascmService.DescribeAscmRamPolicy(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return nil
}
