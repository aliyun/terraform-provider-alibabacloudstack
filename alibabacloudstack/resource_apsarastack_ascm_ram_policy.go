package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"
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
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_ram_policy", "policy alreadyExist", AlibabacloudStackSdkGoERROR)
	}
	if len(check.Data) == 0 {
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "ascm",
			"Action":          "CreateRAMPolicy",
			"Version":         "2019-05-10",
			"policyName":      name,
			"description":     description,
			"policyDocument":  policyDoc,
		}
		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "CreateRAMPolicy"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" rsponse of CreateRAMPolicy : %s", raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_ram_policy", "CreateRAMPolicy", raw)
		}
		bresponse, _ := raw.(*responses.CommonResponse)

		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_ram_policy", "CreateRAMPolicy", AlibabacloudStackSdkGoERROR)
		}
		addDebug("CreateRAMPolicy", raw, requestInfo, bresponse.GetHttpContentString())
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
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsInstanceExist", AlibabacloudStackSdkGoERROR)
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
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":          client.RegionId,
		
		
		"Product":           "ascm",
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"Action":            "UpdateRAMPolicy",
		"Version":           "2019-05-10",
		"ProductName":       "ascm",
		"RamPolicyId":       did[1],
		"NewPolicyName":     name,
		"NewDescription":    description,
		"NewPolicyDocument": policydoc,
	}
	request.Domain = client.Domain
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.ApiName = "UpdateRAMPolicy"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	//check.Data[0].ID = d.Id()

	if attributeUpdate {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw UpdateRAMPolicy : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_ram_policy", "UpdateRAMPolicy", raw)
		}
		addDebug(request.GetActionName(), raw, request)
		log.Printf("total QueryParams and rampolicy %v %v", request.GetQueryParams(), name)
		//d.SetId(name + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))
		//
		//d.Set("ram_id",did[1])
		//d.Set("name", check.Data[0].NewPolicyName)
		//d.Set("description", check.Data[0].NewDescription)
		//d.Set("policy_document", check.Data[0].NewPolicyDocument)
		//return nil
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsPolicyExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsPolicyExist", check, requestInfo, map[string]string{"policyName": did[0]})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			
			
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "ascm",
			"Action":          "RemoveRAMPolicy",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"ramPolicyId":     did[1],
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "RemoveRAMPolicy"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

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
