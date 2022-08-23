package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
	"time"
)

func resourceApsaraStackAscmRamPolicyForRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmRamPolicyForRoleCreate,
		Read:   resourceApsaraStackAscmRamPolicyForRoleRead,
		Update: resourceApsaraStackAscmRamPolicyForRoleUpdate,
		Delete: resourceApsaraStackAscmRamPolicyForRoleDelete,
		Schema: map[string]*schema.Schema{
			"ram_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceApsaraStackAscmRamPolicyForRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	ram_id := d.Get("ram_policy_id").(string)
	roleid := d.Get("role_id").(int)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ascm",
		"Action":          "AddRAMPolicyToRole",
		"Version":         "2019-05-10",
		"ProductName":     "ascm",
		"RamPolicyId":     ram_id,
		"RoleId":          fmt.Sprint(roleid),
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "AddRAMPolicyToRole"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_policy_for_role", "AddRAMPolicyToRole", raw)
	}

	addDebug("AddRAMPolicyToRole", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_ram_policy_for_role", "AddRAMPolicyToRole", ApsaraStackSdkGoERROR)
	}
	addDebug("AddRAMPolicyToRole", raw, requestInfo, bresponse.GetHttpContentString())

	d.SetId(ram_id + COLON_SEPARATED + fmt.Sprint(roleid))

	return resourceApsaraStackAscmRamPolicyForRoleRead(d, meta)
}

func resourceApsaraStackAscmRamPolicyForRoleRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	_, err := ascmService.DescribeAscmRamPolicyForRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ram_policy_id", did[0])
	role_id, _ := strconv.Atoi(did[1])
	d.Set("role_id", role_id)

	return nil
}

func resourceApsaraStackAscmRamPolicyForRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmRamPolicyForRoleCreate(d, meta)

}

func resourceApsaraStackAscmRamPolicyForRoleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	did := strings.Split(d.Id(), COLON_SEPARATED)

	check, err := ascmService.DescribeAscmRamPolicyForRole(d.Id())

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBindingExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"ramPolicyId": did[0]})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "RemoveRAMPolicyFromRole",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"ramPolicyId":     did[0],
			"roleId":          did[1],
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
		request.ApiName = "RemoveRAMPolicyFromRole"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = ascmService.DescribeAscmRamPolicyForRole(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "RemoveRAMPolicyFromRole", ApsaraStackSdkGoERROR)
	}
	return nil
}
