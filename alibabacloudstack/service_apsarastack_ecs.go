package alibabacloudstack

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EcsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *EcsService) JudgeRegionValidation(key, region string) error {
	request := ecs.CreateDescribeRegionsRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeRegions(request)
	})
	if err != nil {
		return fmt.Errorf("DescribeRegions got an error: %#v", err)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeRegionsResponse)
	if resp == nil || len(resp.Regions.Region) < 1 {
		return GetNotFoundErrorFromString("There is no any available region.")
	}

	var rs []string
	for _, v := range resp.Regions.Region {
		if v.RegionId == region {
			return nil
		}
		rs = append(rs, v.RegionId)
	}
	return fmt.Errorf("'%s' is invalid. Expected on %v.", key, strings.Join(rs, ", "))
}

// DescribeZone validate zoneId is valid in region
func (s *EcsService) DescribeZone(id string) (zone ecs.Zone, err error) {
	request := ecs.CreateDescribeZonesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeZonesResponse)
	if len(response.Zones.Zone) < 1 {
		return zone, WrapError(Error("There is no any availability zone in region %s.", s.client.RegionId))
	}

	zoneIds := []string{}
	for _, z := range response.Zones.Zone {
		if z.ZoneId == id {
			return z, nil
		}
		zoneIds = append(zoneIds, z.ZoneId)
	}
	return zone, WrapError(Error("availability_zone %s not exists in region %s, all zones are %s", id, s.client.RegionId, zoneIds))
}

func (s *EcsService) DescribeZones(d *schema.ResourceData) (zones []ecs.Zone, err error) {
	request := ecs.CreateDescribeZonesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	//request.SpotStrategy = d.Get("spot_strategy").(string)
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, "apsarastak_instance_type_families", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeZonesResponse)
	if len(response.Zones.Zone) < 1 {
		return zones, WrapError(Error("There is no any availability zone in region %s.", s.client.RegionId))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		zoneIds := []string{}
		for _, z := range response.Zones.Zone {
			if z.ZoneId == v.(string) {
				return []ecs.Zone{z}, nil
			}
			zoneIds = append(zoneIds, z.ZoneId)
		}
		return zones, WrapError(Error("availability_zone %s not exists in region %s, all zones are %s", v.(string), s.client.RegionId, zoneIds))
	} else {
		return response.Zones.Zone, nil
	}
}

func (s *EcsService) DescribeInstance(id string) (instance ecs.Instance, err error) {
	request := ecs.CreateDescribeInstancesRequest()
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.InstanceIds = convertListToJsonString([]interface{}{id})

	var response *ecs.DescribeInstancesResponse
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)

			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ecs.DescribeInstancesResponse)
		return nil
	})

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	if len(response.Instances.Instance) < 1 {
		return instance, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}
	log.Printf("[ECS Creation]: Getting Instance Details using DescribeInstances API: %s", response.Instances.Instance[0].Status)
	return response.Instances.Instance[0], nil
}

func (s *EcsService) DescribeInstanceAttribute(id string) (instance ecs.DescribeInstanceAttributeResponse, err error) {
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceAttribute(request)
	})
	if err != nil {
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeInstanceAttributeResponse)
	if response.InstanceId != id {
		return instance, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}

	return *response, nil
}

// func (s *EcsService) AssignIpv6Addresses() {}

func (s *EcsService) DescribeInstanceSystemDisk(id, rg string) (disk ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	request.InstanceId = id
	//request.DiskType = string(DiskTypeSystem)
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	var response *ecs.DescribeDisksResponse
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(request)
		})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ecs.DescribeDisksResponse)
		if len(response.Disks.Disk) < 1 {
			return resource.RetryableError(errors.New("Disk Not Found"))
		}
		return nil
	})
	if err != nil {
		return disk, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	log.Printf("[ECS Creation]: Getting Disks Query Params : %s ", request.GetQueryParams())
	log.Printf("[ECS Creation]: Getting Disks response : %s ", response)
	//log.Printf("[ECS Creation]: Getting Disks Details: %s, Instance ID: %s, Id_to_compare: %s ",response.Disks.Disk[0],response.Disks.Disk[0].InstanceId,id)
	if len(response.Disks.Disk) < 1 || response.Disks.Disk[0].InstanceId != id {
		return disk, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}
	for _, disk := range response.Disks.Disk {
		if disk.Type == "system" {
			return disk, nil
		}
	}
	return disk, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
}

// ResourceAvailable check resource available for zone
func (s *EcsService) ResourceAvailable(zone ecs.Zone, resourceType ResourceType) error {
	for _, res := range zone.AvailableResourceCreation.ResourceTypes {
		if res == string(resourceType) {
			return nil
		}
	}
	return WrapError(Error("%s is not available in %s zone of %s region", resourceType, zone.ZoneId, s.client.Region))
}

func (s *EcsService) DiskAvailable(zone ecs.Zone, diskCategory DiskCategory) error {
	for _, disk := range zone.AvailableDiskCategories.DiskCategories {
		if disk == string(diskCategory) {
			return nil
		}
	}
	return WrapError(Error("%s is not available in %s zone of %s region", diskCategory, zone.ZoneId, s.client.Region))
}

func (s *EcsService) JoinSecurityGroups(instanceId string, securityGroupIds []string) error {
	request := ecs.CreateJoinSecurityGroupRequest()
	request.InstanceId = instanceId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	for _, sid := range securityGroupIds {
		request.SecurityGroupId = sid
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.JoinSecurityGroup(request)
		})
		if err != nil && IsExpectedErrors(err, []string{"InvalidInstanceId.AlreadyExists"}) {
			return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *EcsService) LeaveSecurityGroups(instanceId string, securityGroupIds []string) error {
	request := ecs.CreateLeaveSecurityGroupRequest()
	request.InstanceId = instanceId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	for _, sid := range securityGroupIds {
		request.SecurityGroupId = sid
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.LeaveSecurityGroup(request)
		})
		if err != nil && IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *EcsService) DescribeSecurityGroup(id string) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.SecurityGroupId = id
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if response.SecurityGroupId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("Security Group", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}

	return *response, nil
}

func (s *EcsService) DescribeSecurityGroupRule(id string) (rule ecs.Permission, err error) {
	parts, err := ParseResourceId(id, 8)
	if err != nil {
		return rule, WrapError(err)
	}
	groupId, direction, ipProtocol, portRange, nicType, cidr_ip, policy := parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6]
	priority, err := strconv.Atoi(parts[7])
	if err != nil {
		return rule, WrapError(err)
	}
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.SecurityGroupId = groupId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Direction = direction
	request.NicType = nicType
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if response == nil {
		return rule, GetNotFoundErrorFromString(GetNotFoundMessage("Security Group", groupId))
	}

	for _, ru := range response.Permissions.Permission {
		if strings.ToLower(string(ru.IpProtocol)) == ipProtocol && ru.PortRange == portRange {
			cidr := ru.SourceCidrIp
			if direction == string(DirectionIngress) && cidr == "" {
				cidr = ru.SourceGroupId
			}
			if direction == string(DirectionEgress) {
				if cidr = ru.DestCidrIp; cidr == "" {
					cidr = ru.DestGroupId
				}
			}
			if cidr == cidr_ip && strings.ToLower(string(ru.Policy)) == policy && ru.Priority == strconv.Itoa(priority) {
				return ru, nil
			}
		}
	}

	return rule, WrapErrorf(Error(GetNotFoundMessage("Security Group Rule", id)), NotFoundMsg, ProviderERROR, response.RequestId)

}

func (s *EcsService) DescribeAvailableResources(d *schema.ResourceData, meta interface{}, destination DestinationResource) (zoneId string, validZones []ecs.AvailableZone, requestId string, err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	// Before creating resources, check input parameters validity according available zone.
	// If availability zone is nil, it will return all of supported resources in the current.
	request := ecs.CreateDescribeAvailableResourceRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.DestinationResource = string(destination)
	request.IoOptimized = string(IOOptimized)

	if v, ok := d.GetOk("availability_zone"); ok && strings.TrimSpace(v.(string)) != "" {
		zoneId = strings.TrimSpace(v.(string))
	} else if v, ok := d.GetOk("vswitch_id"); ok && strings.TrimSpace(v.(string)) != "" {
		vpcService := VpcService{s.client}
		if vsw, err := vpcService.DescribeVSwitch(strings.TrimSpace(v.(string))); err == nil {
			zoneId = vsw.ZoneId
		}
	}

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAvailableResource(request)
	})
	if err != nil {
		return "", nil, "", WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackGoClientFailure)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeAvailableResourceResponse)
	requestId = response.RequestId

	if len(response.AvailableZones.AvailableZone) < 1 {
		err = WrapError(Error("There are no availability resources in the region: %s. RequestId: %s.", client.RegionId, requestId))
		return
	}

	valid := false
	soldout := false
	var expectedZones []string
	for _, zone := range response.AvailableZones.AvailableZone {
		if zone.Status == string(SoldOut) {
			if zone.ZoneId == zoneId {
				soldout = true
			}
			continue
		}
		if zoneId != "" && zone.ZoneId == zoneId {
			valid = true
			validZones = append(make([]ecs.AvailableZone, 1), zone)
			break
		}
		expectedZones = append(expectedZones, zone.ZoneId)
		validZones = append(validZones, zone)
	}
	if zoneId != "" {
		if !valid {
			err = WrapError(Error("Availability zone %s status is not available in the region %s. Expected availability zones: %s. RequestId: %s.",
				zoneId, client.RegionId, strings.Join(expectedZones, ", "), requestId))
			return
		}
		if soldout {
			err = WrapError(Error("Availability zone %s status is sold out in the region %s. Expected availability zones: %s. RequestId: %s.",
				zoneId, client.RegionId, strings.Join(expectedZones, ", "), requestId))
			return
		}
	}

	if len(validZones) <= 0 {
		err = WrapError(Error("There is no availability resources in the region %s. Please choose another region. RequestId: %s.", client.RegionId, response.RequestId))
		return
	}

	return
}

func (s *EcsService) InstanceTypeValidation(targetType, zoneId string, validZones []ecs.AvailableZone) error {

	mapInstanceTypeToZones := make(map[string]string)
	var expectedInstanceTypes []string
	for _, zone := range validZones {
		if zoneId != "" && zoneId != zone.ZoneId {
			continue
		}
		for _, r := range zone.AvailableResources.AvailableResource {
			if r.Type == string(InstanceTypeResource) {
				for _, t := range r.SupportedResources.SupportedResource {
					if t.Status == string(SoldOut) {
						continue
					}
					if targetType == t.Value {
						return nil
					}

					if _, ok := mapInstanceTypeToZones[t.Value]; !ok {
						expectedInstanceTypes = append(expectedInstanceTypes, t.Value)
						mapInstanceTypeToZones[t.Value] = t.Value
					}
				}
			}
		}
	}
	if zoneId != "" {
		return WrapError(Error("The instance type %s is solded out or is not supported in the zone %s. Expected instance types: %s", targetType, zoneId, strings.Join(expectedInstanceTypes, ", ")))
	}
	return WrapError(Error("The instance type %s is solded out or is not supported in the region %s. Expected instance types: %s", targetType, s.client.RegionId, strings.Join(expectedInstanceTypes, ", ")))
}

func (s *EcsService) QueryInstancesWithKeyPair(instanceIdsStr, keyPair string) (instanceIds []string, instances []ecs.Instance, err error) {

	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.InstanceIds = instanceIdsStr
	request.KeyPairName = keyPair
	for {
		raw, e := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		if e != nil {
			err = WrapErrorf(e, DefaultErrorMsg, keyPair, request.GetActionName(), AlibabacloudStackSdkGoERROR)
			return
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		object, _ := raw.(*ecs.DescribeInstancesResponse)
		if len(object.Instances.Instance) < 0 {
			return
		}
		for _, inst := range object.Instances.Instance {
			instanceIds = append(instanceIds, inst.InstanceId)
			instances = append(instances, inst)
		}
		if len(instances) < PageSizeLarge {
			break
		}
		if page, e := getNextpageNumber(request.PageNumber); e != nil {
			err = WrapErrorf(e, DefaultErrorMsg, keyPair, request.GetActionName(), AlibabacloudStackSdkGoERROR)
			return
		} else {
			request.PageNumber = page
		}
	}
	return
}

func (s *EcsService) DescribeKeyPair(id string) (keyPair ecs.KeyPair, err error) {
	request := ecs.CreateDescribeKeyPairsRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeyId":     s.client.AccessKey,
		"AccessKeySecret": s.client.SecretKey,
		"SecurityToken":   s.client.Config.SecurityToken,
		"Product":         "ecs",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup}
	request.KeyPairName = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeKeyPairs(request)
	})
	if err != nil {
		return keyPair, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	object, _ := raw.(*ecs.DescribeKeyPairsResponse)
	if len(object.KeyPairs.KeyPair) < 1 || object.KeyPairs.KeyPair[0].KeyPairName != id {
		return keyPair, WrapErrorf(Error(GetNotFoundMessage("KeyPair", id)), NotFoundMsg, ProviderERROR, object.RequestId)
	}
	return object.KeyPairs.KeyPair[0], nil

}

func (s *EcsService) DescribeKeyPairAttachment(id string) (keyPair ecs.KeyPair, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return
	}
	keyPairName := parts[0]
	keyPair, err = s.DescribeKeyPair(keyPairName)
	if err != nil {
		return keyPair, WrapError(err)
	}
	if keyPair.KeyPairName != keyPairName {
		err = WrapErrorf(Error(GetNotFoundMessage("KeyPairAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return keyPair, nil

}

func (s *EcsService) DescribeDisk(id string) (disk ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	request.DiskIds = convertListToJsonString([]interface{}{id})
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(request)
	})
	if err != nil {
		return disk, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	response, _ := raw.(*ecs.DescribeDisksResponse)
	if len(response.Disks.Disk) < 1 || response.Disks.Disk[0].DiskId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("Disk", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return response.Disks.Disk[0], nil
}

func (s *EcsService) DescribeDiskAttachment(id string) (disk ecs.Disk, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return disk, WrapError(err)
	}
	disk, err = s.DescribeDisk(parts[0])
	if err != nil {
		return disk, WrapError(err)
	}

	if disk.InstanceId != parts[1] && disk.Status != string(InUse) {
		err = WrapErrorf(Error(GetNotFoundMessage("DiskAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *EcsService) DescribeDisksByType(instanceId string, diskType DiskType) (disk []ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	if instanceId != "" {
		request.InstanceId = instanceId
	}
	request.DiskType = string(diskType)

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(request)
	})
	if err != nil {
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeDisksResponse)
	if resp == nil {
		return
	}
	return resp.Disks.Disk, nil
}

func (s *EcsService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []ecs.Tag, err error) {
	request := ecs.CreateDescribeTagsRequest()
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.ResourceType = string(resourceType)
	request.ResourceId = resourceId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeTags(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeTagsResponse)
	if len(response.Tags.Tag) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("Tags", resourceId)), NotFoundMsg, ProviderERROR)
		return
	}

	return response.Tags.Tag, nil
}

func (s *EcsService) DescribeImageById(id string) (image ecs.Image, err error) {
	request := ecs.CreateDescribeImagesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.ImageId = id
	request.ImageOwnerAlias = "self"
	request.Status = fmt.Sprintf("%s,%s,%s,%s,%s", "Creating", "Waiting", "Available", "UnAvailable", "CreateFailed")
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImages(request)
	})
	if err != nil {
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeImagesResponse)
	if resp == nil || len(resp.Images.Image) < 1 {
		return image, GetNotFoundErrorFromString(GetNotFoundMessage("Image", id))
	}

	return resp.Images.Image[0], nil
}

func (s *EcsService) deleteImage(d *schema.ResourceData) error {

	object, err := s.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	request := ecs.CreateDeleteImageRequest()
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if force, ok := d.GetOk("force"); ok {
		request.Force = requests.NewBoolean(force.(bool))
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.ImageId = object.ImageId

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteImage(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 3*time.Second, s.ImageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func (s *EcsService) updateImage(d *schema.ResourceData) error {

	d.Partial(true)

	err := setTags(s.client, TagResourceImage, d)
	if err != nil {
		return WrapError(err)
	} else {
		//d.SetPartial("tags")
	}

	request := ecs.CreateModifyImageAttributeRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.ImageId = d.Id()

	if d.HasChange("description") || d.HasChange("name") || d.HasChange("image_name") {
		if description, ok := d.GetOk("description"); ok {
			request.Description = description.(string)
		}
		if imageName, ok := d.GetOk("image_name"); ok {
			request.ImageName = imageName.(string)
		} else {
			if imageName, ok := d.GetOk("name"); ok {
				request.ImageName = imageName.(string)
			}
		}
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyImageAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		//d.SetPartial("name")
		//d.SetPartial("image_name")
		//d.SetPartial("description")
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return nil
}

func (s *EcsService) DescribeNetworkInterface(id string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	request := ecs.CreateDescribeNetworkInterfacesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	eniIds := []string{id}
	request.NetworkInterfaceId = &eniIds
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeNetworkInterfaces(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeNetworkInterfacesResponse)
	//if len(response.NetworkInterfaceSets.NetworkInterfaceSet) < 1 ||
	//	response.NetworkInterfaceSets.NetworkInterfaceSet[0].NetworkInterfaceId != id {
	//	err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterface", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	//	return
	//}

	var found = bool(false)
	var result = &ecs.DescribeNetworkInterfacesResponse{}
	for _, k := range response.NetworkInterfaceSets.NetworkInterfaceSet {
		if /*len(response.NetworkInterfaceSets.NetworkInterfaceSet) < 1 */
		k.NetworkInterfaceId == id {
			found = true
			result.NetworkInterfaceSets.NetworkInterfaceSet = append(result.NetworkInterfaceSets.NetworkInterfaceSet, ecs.NetworkInterfaceSet{
				NetworkInterfaceId:   k.NetworkInterfaceId,
				Status:               k.Status,
				Type:                 k.Type,
				VpcId:                k.VpcId,
				VSwitchId:            k.VSwitchId,
				ZoneId:               k.ZoneId,
				PrivateIpAddress:     k.PrivateIpAddress,
				MacAddress:           k.MacAddress,
				NetworkInterfaceName: k.NetworkInterfaceName,
				Description:          k.Description,
				InstanceId:           k.InstanceId,
				CreationTime:         k.CreationTime,
				ResourceGroupId:      k.ResourceGroupId,
				ServiceID:            k.ServiceID,
				ServiceManaged:       k.ServiceManaged,
				QueueNumber:          k.QueueNumber,
				OwnerId:              k.OwnerId,
				SecurityGroupIds:     k.SecurityGroupIds,
				AssociatedPublicIp:   k.AssociatedPublicIp,
				Attachment:           k.Attachment,
				PrivateIpSets:        k.PrivateIpSets,
				Ipv6Sets:             k.Ipv6Sets,
				Tags:                 k.Tags,
			})
			return result.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
		}
	}
	if !found {
		err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterface", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}

	return response.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
}

func (s *EcsService) DescribeNetworkInterfaceAttachment(id string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	parts, err := ParseResourceId(id, 2)

	if err != nil {
		return networkInterface, WrapError(err)
	}
	eniId, instanceId := parts[0], parts[1]
	request := ecs.CreateDescribeNetworkInterfacesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.InstanceId = instanceId
	eniIds := []string{eniId}
	request.NetworkInterfaceId = &eniIds
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeNetworkInterfaces(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeNetworkInterfacesResponse)
	//if len(response.NetworkInterfaceSets.NetworkInterfaceSet) < 1 ||
	//	response.NetworkInterfaceSets.NetworkInterfaceSet[0].NetworkInterfaceId != eniId {
	//	err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterfaceAttachment", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	//	return
	//}
	var found = bool(false)
	var result = &ecs.DescribeNetworkInterfacesResponse{}
	for _, k := range response.NetworkInterfaceSets.NetworkInterfaceSet {
		if /*len(response.NetworkInterfaceSets.NetworkInterfaceSet) < 1 */
		k.NetworkInterfaceId == eniId {
			found = true
			result.NetworkInterfaceSets.NetworkInterfaceSet = append(result.NetworkInterfaceSets.NetworkInterfaceSet, ecs.NetworkInterfaceSet{
				NetworkInterfaceId:   k.NetworkInterfaceId,
				Status:               k.Status,
				Type:                 k.Type,
				VpcId:                k.VpcId,
				VSwitchId:            k.VSwitchId,
				ZoneId:               k.ZoneId,
				PrivateIpAddress:     k.PrivateIpAddress,
				MacAddress:           k.MacAddress,
				NetworkInterfaceName: k.NetworkInterfaceName,
				Description:          k.Description,
				InstanceId:           k.InstanceId,
				CreationTime:         k.CreationTime,
				ResourceGroupId:      k.ResourceGroupId,
				ServiceID:            k.ServiceID,
				ServiceManaged:       k.ServiceManaged,
				QueueNumber:          k.QueueNumber,
				OwnerId:              k.OwnerId,
				SecurityGroupIds:     k.SecurityGroupIds,
				AssociatedPublicIp:   k.AssociatedPublicIp,
				Attachment:           k.Attachment,
				PrivateIpSets:        k.PrivateIpSets,
				Ipv6Sets:             k.Ipv6Sets,
				Tags:                 k.Tags,
			})
			return result.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
		}
	}
	if !found {
		err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterfaceAttachment", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		//err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterface", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}

	return response.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
}

// WaitForInstance waits for instance to given status
func (s *EcsService) WaitForEcsInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeInstance(instanceId)
		if err != nil && !NotFoundError(err) {
			return err
		}
		if instance.Status == string(status) {
			//Sleep one more time for timing issues
			time.Sleep(DefaultIntervalMedium * time.Second)
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("ECS Instance", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}

// WaitForInstance waits for instance to given status
func (s *EcsService) InstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}

		return object, object.Status, nil
	}
}
func (s *EcsService) deleteImageforDest(d *schema.ResourceData, region string) error {

	object, err := s.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	request := ecs.CreateDeleteImageRequest()
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	if force, ok := d.GetOk("force"); ok {
		request.Force = requests.NewBoolean(force.(bool))
	}
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = region
	request.ImageId = object.ImageId

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteImage(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 3*time.Second, s.ImageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
func (s *EcsService) DescribeImage(id, region string) (image ecs.Image, err error) {
	request := ecs.CreateDescribeImagesRequest()
	request.RegionId = region
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ImageId = id
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs"}
	request.Status = fmt.Sprintf("%s,%s,%s,%s,%s", "Creating", "Waiting", "Available", "UnAvailable", "CreateFailed")
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImages(request)
	})
	log.Printf("[DEBUG] status %#v", raw)
	if err != nil {
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeImagesResponse)
	if resp == nil || len(resp.Images.Image) < 1 {
		return image, GetNotFoundErrorFromString(GetNotFoundMessage("Image", id))
	}
	return resp.Images.Image[0], nil
}

func (s *EcsService) ImageStateRefreshFuncforcopy(id string, region string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeImage(id, region)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}
func (s *EcsService) WaitForDisk(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeDisk(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		// Disk need 3-5 seconds to get ExpiredTime after the status is available
		if object.Status == string(status) && object.ExpiredTime != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForSecurityGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeSecurityGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForKeyPair(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeKeyPair(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForDiskAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDiskAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForNetworkInterface(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeNetworkInterface(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
	}
}

func (s *EcsService) QueryPrivateIps(eniId string) ([]string, error) {
	if eni, err := s.DescribeNetworkInterface(eniId); err != nil {
		return nil, fmt.Errorf("Describe NetworkInterface(%s) failed, %s", eniId, err)
	} else {
		filterIps := make([]string, 0, len(eni.PrivateIpSets.PrivateIpSet))
		for i := range eni.PrivateIpSets.PrivateIpSet {
			if eni.PrivateIpSets.PrivateIpSet[i].Primary {
				continue
			}
			filterIps = append(filterIps, eni.PrivateIpSets.PrivateIpSet[i].PrivateIpAddress)
		}
		return filterIps, nil
	}
}

func (s *EcsService) WaitForVpcAttributesChanged(instanceId, vswitchId, privateIp string) error {
	deadline := time.Now().Add(DefaultTimeout * time.Second)
	for {
		if time.Now().After(deadline) {
			return WrapError(Error("Wait for VPC attributes changed timeout"))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		instance, err := s.DescribeInstance(instanceId)
		if err != nil {
			return WrapError(err)
		}

		if instance.VpcAttributes.PrivateIpAddress.IpAddress[0] != privateIp {
			continue
		}

		if instance.VpcAttributes.VSwitchId != vswitchId {
			continue
		}

		return nil
	}
}

func (s *EcsService) WaitForPrivateIpsCountChanged(eniId string, count int) error {
	deadline := time.Now().Add(DefaultTimeout * time.Second)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Wait for private IP addrsses count changed timeout")
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		ips, err := s.QueryPrivateIps(eniId)
		if err != nil {
			return fmt.Errorf("Query private IP failed, %s", err)
		}
		if len(ips) == count {
			return nil
		}
	}
}

func (s *EcsService) WaitForPrivateIpsListChanged(eniId string, ipList []string) error {
	deadline := time.Now().Add(DefaultTimeout * time.Second)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Wait for private IP addrsses list changed timeout")
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		ips, err := s.QueryPrivateIps(eniId)
		if err != nil {
			return fmt.Errorf("Query private IP failed, %s", err)
		}

		if len(ips) != len(ipList) {
			continue
		}

		diff := false
		for i := range ips {
			exist := false
			for j := range ipList {
				if ips[i] == ipList[j] {
					exist = true
					break
				}
			}
			if !exist {
				diff = true
				break
			}
		}

		if !diff {
			return nil
		}
	}
}

func (s *EcsService) WaitForModifySecurityGroupPolicy(id, target string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSecurityGroup(id)
		if err != nil {
			return WrapError(err)
		}
		if object.InnerAccessPolicy == target {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InnerAccessPolicy, target, ProviderERROR)
		}
	}
}

func (s *EcsService) AttachKeyPair(keyName string, instanceIds []interface{}) error {
	request := ecs.CreateAttachKeyPairRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.KeyPairName = keyName
	request.InstanceIds = convertListToJsonString(instanceIds)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachKeyPair(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, keyName, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func (s *EcsService) SnapshotStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSnapshot(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *EcsService) ImageStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeImageById(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *EcsService) TaskStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTaskById(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if object.TaskStatus == failState {
				return object, object.TaskStatus, WrapError(Error(FailedToReachTargetStatus, object.TaskStatus))
			}
		}
		return object, object.TaskStatus, nil
	}
}

func (s *EcsService) DescribeTaskById(id string) (task *ecs.DescribeTaskAttributeResponse, err error) {
	request := ecs.CreateDescribeTaskAttributeRequest()
	request.TaskId = id
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeTaskAttribute(request)
	})
	if err != nil {
		return task, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	task, _ = raw.(*ecs.DescribeTaskAttributeResponse)

	if task.TaskId == "" {
		return task, GetNotFoundErrorFromString(GetNotFoundMessage("task", id))
	}
	return task, nil
}

func (s *EcsService) DescribeSnapshot(id string) (*ecs.Snapshot, error) {
	snapshot := &ecs.Snapshot{}
	request := ecs.CreateDescribeSnapshotsRequest()
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.SnapshotIds = fmt.Sprintf("[\"%s\"]", id)
	request.QueryParams["SnapshotId"] = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSnapshots(request)
	})
	if err != nil {
		return snapshot, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeSnapshotsResponse)
	for _, k := range response.Snapshots.Snapshot {
		if k.SnapshotId == id {
			return &k, nil
		}
	}
	if len(response.Snapshots.Snapshot) < 1 || response.Snapshots.Snapshot[0].SnapshotId != id {
		return snapshot, WrapErrorf(Error(GetNotFoundMessage("Snapshot", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}
	return &response.Snapshots.Snapshot[0], nil
}

func (s *EcsService) DescribeSnapshotPolicy(id string) (*ecs.AutoSnapshotPolicy, error) {
	policy := &ecs.AutoSnapshotPolicy{}
	request := ecs.CreateDescribeAutoSnapshotPolicyExRequest()
	request.AutoSnapshotPolicyId = id
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAutoSnapshotPolicyEx(request)
	})
	if err != nil {
		return policy, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*ecs.DescribeAutoSnapshotPolicyExResponse)
	if len(response.AutoSnapshotPolicies.AutoSnapshotPolicy) != 1 ||
		response.AutoSnapshotPolicies.AutoSnapshotPolicy[0].AutoSnapshotPolicyId != id {
		return policy, WrapErrorf(Error(GetNotFoundMessage("SnapshotPolicy", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}

	return &response.AutoSnapshotPolicies.AutoSnapshotPolicy[0], nil
}

func (s *EcsService) DescribeReservedInstance(id string) (reservedInstance ecs.ReservedInstance, err error) {
	request := ecs.CreateDescribeReservedInstancesRequest()
	var balance = &[]string{id}
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ReservedInstanceId = balance
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeReservedInstances(request)
	})
	if err != nil {
		return reservedInstance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeReservedInstancesResponse)

	if len(response.ReservedInstances.ReservedInstance) != 1 ||
		response.ReservedInstances.ReservedInstance[0].ReservedInstanceId != id {
		return reservedInstance, GetNotFoundErrorFromString(GetNotFoundMessage("PurchaseReservedInstance", id))
	}
	return response.ReservedInstances.ReservedInstance[0], nil
}

func (s *EcsService) WaitForReservedInstance(id string, status Status, timeout int) error {
	deadLine := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		reservedInstance, err := s.DescribeReservedInstance(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if reservedInstance.Status == string(status) {
			return nil
		}

		if time.Now().After(deadLine) {
			return WrapErrorf(GetTimeErrorFromString("ECS WaitForSnapshotPolicy"), WaitTimeoutMsg, id, GetFunc(1), timeout, reservedInstance.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) WaitForSnapshotPolicy(id string, status Status, timeout int) error {
	deadLine := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		snapshotPolicy, err := s.DescribeSnapshotPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}

		if snapshotPolicy.Status == string(status) {
			return nil
		}

		if time.Now().After(deadLine) {
			return WrapErrorf(GetTimeErrorFromString("ECS WaitForSnapshotPolicy"), WaitTimeoutMsg, id, GetFunc(1), timeout, snapshotPolicy.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) DescribeLaunchTemplate(id string) (set ecs.LaunchTemplateSet, err error) {

	request := ecs.CreateDescribeLaunchTemplatesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.LaunchTemplateId = &[]string{id}

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeLaunchTemplates(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeLaunchTemplatesResponse)
	if len(response.LaunchTemplateSets.LaunchTemplateSet) != 1 ||
		response.LaunchTemplateSets.LaunchTemplateSet[0].LaunchTemplateId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("LaunchTemplate", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}

	return response.LaunchTemplateSets.LaunchTemplateSet[0], nil

}

func (s *EcsService) DescribeLaunchTemplateVersion(id string, version int) (set ecs.LaunchTemplateVersionSet, err error) {

	request := ecs.CreateDescribeLaunchTemplateVersionsRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.LaunchTemplateId = id
	request.LaunchTemplateVersion = &[]string{strconv.FormatInt(int64(version), 10)}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeLaunchTemplateVersions(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidLaunchTemplate.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeLaunchTemplateVersionsResponse)
	if len(response.LaunchTemplateVersionSets.LaunchTemplateVersionSet) != 1 ||
		response.LaunchTemplateVersionSets.LaunchTemplateVersionSet[0].LaunchTemplateId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("LaunchTemplateVersion", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}

	return response.LaunchTemplateVersionSets.LaunchTemplateVersionSet[0], nil

}

func (s *EcsService) WaitForLaunchTemplate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeLaunchTemplate(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.LaunchTemplateId == id && string(status) != string(Deleted) {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) DescribeImageShareByImageId(id string) (imageShare *ecs.DescribeImageSharePermissionResponse, err error) {
	request := ecs.CreateDescribeImageSharePermissionRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return imageShare, WrapError(err)
	}
	request.ImageId = parts[0]
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImageSharePermission(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidImageId.NotFound"}) {
			return imageShare, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return imageShare, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeImageSharePermissionResponse)
	if len(resp.Accounts.Account) == 0 {
		return imageShare, WrapErrorf(Error(GetNotFoundMessage("ModifyImageSharePermission", id)), NotFoundMsg, ProviderERROR, resp.RequestId)
	}
	return resp, nil
}

func (s *EcsService) WaitForAutoProvisioningGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeAutoProvisioningGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
	}
}

func (s *EcsService) DescribeAutoProvisioningGroup(id string) (group ecs.AutoProvisioningGroup, err error) {
	request := ecs.CreateDescribeAutoProvisioningGroupsRequest()
	ids := []string{id}
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.AutoProvisioningGroupId = &ids
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, e := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAutoProvisioningGroups(request)
	})
	if e != nil {
		err = WrapErrorf(e, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeAutoProvisioningGroupsResponse)
	for _, v := range response.AutoProvisioningGroups.AutoProvisioningGroup {
		if v.AutoProvisioningGroupId == id {
			return v, nil
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("AutoProvisioningGroup", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	return
}

func (s *EcsService) tagsToMap(tags []ecs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ecsTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func (s *EcsService) ecsTagIgnored(t ecs.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found AlibabacloudStack Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *EcsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]ecs.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, ecs.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := ecs.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		if strings.ToLower(s.client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": s.client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "InvalidResourceId.NotFound", "InvalidResourceType.NotFound", "MissingParameter.RegionId", "MissingParameter.ResourceIds", "MissingParameter.ResourceType", "MissingParameter.TagOwnerBid", "MissingParameter.TagOwnerUid", "MissingParameter.Tags"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := ecs.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		if strings.ToLower(s.client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": s.client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "InvalidResourceId.NotFound", "InvalidResourceType.NotFound", "MissingParameter.RegionId", "MissingParameter.ResourceIds", "MissingParameter.ResourceType", "MissingParameter.TagOwnerBid", "MissingParameter.TagOwnerUid", "MissingParameter.Tags"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func (s *EcsService) SetResourceTagsNew(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				request["Product"] = "ascm"
				request["OrganizationId"] = s.client.Department
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)})
				if err != nil {
					if IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				request["Product"] = "ascm"
				request["OrganizationId"] = s.client.Department
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)})
				if err != nil {
					if IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *EcsService) SetSystemDiskTags(d *schema.ResourceData) error {
	resourceType := "disk"
	diskid := d.Get("system_disk_id").(string)
	if diskid == "" {
		disk, err := s.DescribeInstanceSystemDisk(d.Id(), s.client.ResourceGroup)
		if err != nil {
			return WrapError(err)
		}
		diskid = disk.DiskId
	}
	conn, err := s.client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	instance_tags := d.Get("tags").(map[string]interface{})
	tags := d.Get("system_disk_tags").(map[string]interface{})
	removed := make([]string, 0)
	for key, value := range instance_tags {
		old, ok := tags[key]
		if !ok || old != value {
			// Delete it!
			removed = append(removed, key)
		}
	}
	action := "UnTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": diskid,
	}
	for i, key := range removed {
		request[fmt.Sprintf("TagKey.%d", i+1)] = key
	}
	wait := incrementalWait(2*time.Second, 1*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		request["Product"] = "ascm"
		request["OrganizationId"] = s.client.Department
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)

			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	add_request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": diskid,
	}
	count := 1
	for key, value := range tags {
		add_request[fmt.Sprintf("Tag.%d.Key", count)] = key
		add_request[fmt.Sprintf("Tag.%d.Value", count)] = value
		count++
	}
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		add_request["Product"] = "ascm"
		add_request["OrganizationId"] = s.client.Department
		response, err := conn.DoRequest(StringPointer("TagResources"), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, add_request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)

			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, add_request, add_request)
		return nil
	})
	return err
}

func (s *EcsService) UpdateSystemDiskTags(d *schema.ResourceData) error {
	resourceType := "disk"
	diskid := d.Get("system_disk_id").(string)
	if diskid == "" {
		disk, err := s.DescribeInstanceSystemDisk(d.Id(), s.client.ResourceGroup)
		if err != nil {
			return WrapError(err)
		}
		diskid = disk.DiskId
	}
	oraw, nraw := d.GetChange("system_disk_tags")
	removedTags := oraw.(map[string]interface{})
	added := nraw.(map[string]interface{})
	// Build the list of what to remove
	removed := make([]string, 0)
	for key, value := range removedTags {
		old, ok := added[key]
		if !ok || old != value {
			// Delete it!
			removed = append(removed, key)
		}
	}
	conn, err := s.client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	removedTagKeys := make([]string, 0)
	for _, v := range removed {
		if !ignoredTags(v, "") {
			removedTagKeys = append(removedTagKeys, v)
		}
	}
	if len(removedTagKeys) > 0 {
		action := "UnTagResources"
		request := map[string]interface{}{
			"RegionId":     s.client.RegionId,
			"ResourceType": resourceType,
			"ResourceId.1": diskid,
		}
		for i, key := range removedTagKeys {
			request[fmt.Sprintf("TagKey.%d", i+1)] = key
		}
		wait := incrementalWait(2*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			request["Product"] = "ascm"
			request["OrganizationId"] = s.client.Department
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)

				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}
	}
	if len(added) > 0 {
		action := "TagResources"
		request := map[string]interface{}{
			"RegionId":     s.client.RegionId,
			"ResourceType": resourceType,
			"ResourceId.1": diskid,
		}
		count := 1
		for key, value := range added {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}

		wait := incrementalWait(2*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			request["Product"] = "ascm"
			request["OrganizationId"] = s.client.Department
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)

				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func (s *EcsService) DescribeEcsDedicatedHost(id string) (object ecs.DedicatedHost, err error) {
	request := ecs.CreateDescribeDedicatedHostsRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "ecs", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	response := new(ecs.DescribeDedicatedHostsResponse)
	for {

		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDedicatedHosts(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidLockReason.NotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("EcsDedicatedHost", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ecs.DescribeDedicatedHostsResponse)

		if len(response.DedicatedHosts.DedicatedHost) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("EcsDedicatedHost", id)), NotFoundMsg, ProviderERROR, response.RequestId)
			return object, err
		}
		for _, object := range response.DedicatedHosts.DedicatedHost {
			if object.DedicatedHostId == id {
				return object, nil
			}
		}
		if len(response.DedicatedHosts.DedicatedHost) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return object, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("EcsDedicatedHost", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	return
}

func (s *EcsService) EcsDedicatedHostStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcsDedicatedHost(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *EcsService) DescribeEcsEbsStorageSet(id string) (result *datahub.EcsDescribeEcsEbsStorageSetsResult, err error) {

	resp := &datahub.EcsDescribeEcsEbsStorageSetsResult{}
	action := "DescribeStorageSets"
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "Ecs"
	request.Domain = s.client.Domain
	request.Version = "2014-05-26"
	request.RegionId = s.client.RegionId
	request.ApiName = action
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers["x-ascm-product-name"] = "Ecs"
	request.Headers["x-acs-organizationId"] = s.client.Department

	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"AccessKeyId":     s.client.AccessKey,
		"RegionId":        s.client.RegionId,
		"Product":         "Ecs",
		"Version":         "2014-05-26",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Action":          action,
		"PageNumber":      "1",
	}
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	runtime.SetAutoretry(true)
	//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	raw, err := s.client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "Operation.Forbidden"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("EcsEbsStorageSet", id)), NotFoundMsg, ProviderERROR)
			return resp, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
		return resp, err
	}
	addDebug(action, raw, request)

	bresponse := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	//v, err := jsonpath.Get("$.Commands.Command", bresponse)
	if err != nil {
		return resp, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Commands.Command", raw)
	}
	//if len(v.([]interface{})) < 1 {
	//	return object, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, raw)
	//} else {
	//	if v.([]interface{})[0].(map[string]interface{})["CommandId"].(string) != id {
	//		return object, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, raw)
	//	}
	//}
	//object = v.([]interface{})[0].(map[string]interface{})
	return resp, nil
}

func (s *EcsService) DescribeEcsCommand(id string) (result *datahub.EcsDescribeEcsCommandResult, err error) {

	resp := &datahub.EcsDescribeEcsCommandResult{}
	action := "DescribeCommands"
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "Ecs"
	request.Domain = s.client.Domain
	request.Version = "2014-05-26"
	request.RegionId = s.client.RegionId
	request.ApiName = action
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers["x-ascm-product-name"] = "Ecs"
	request.Headers["x-acs-organizationId"] = s.client.Department

	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"AccessKeyId":     s.client.AccessKey,
		"RegionId":        s.client.RegionId,
		"Product":         "Ecs",
		"Version":         "2014-05-26",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Action":          action,
		"CommandId":       id,
	}
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	runtime.SetAutoretry(true)
	//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	raw, err := s.client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "Operation.Forbidden"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("EcsCommand", id)), NotFoundMsg, ProviderERROR)
			return resp, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
		return resp, err
	}
	addDebug(action, raw, request)

	bresponse := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	//v, err := jsonpath.Get("$.Commands.Command", bresponse)
	if err != nil {
		return resp, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Commands.Command", raw)
	}
	//if len(v.([]interface{})) < 1 {
	//	return object, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, raw)
	//} else {
	//	if v.([]interface{})[0].(map[string]interface{})["CommandId"].(string) != id {
	//		return object, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, raw)
	//	}
	//}
	//object = v.([]interface{})[0].(map[string]interface{})
	return resp, nil
}
func (s *EcsService) DescribeEcsHpcCluster(id string) (result *datahub.EcsDescribeEcsHpcClusterResult, err error) {
	//var response map[string]interface{}

	resp := &datahub.EcsDescribeEcsHpcClusterResult{}
	action := "DescribeHpcClusters"
	ids, err := json.Marshal([]string{id})
	if err != nil {
		return nil, err
	}
	//request := map[string]interface{}{
	//	"RegionId":      s.client.RegionId,
	//	"HpcClusterIds": string(ids),
	//}
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	runtime.SetAutoretry(true)
	ClientToken := buildClientToken("DescribeHpcClusters")
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Ecs"
	request.Domain = s.client.Domain
	request.Version = "2014-05-26"
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"AccessKeyId":     s.client.AccessKey,
		"RegionId":        s.client.RegionId,
		"Product":         "Ecs",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Action":          action,
		"HpcClusterIds":   string(ids),
		"ClientToken":     ClientToken,
	}

	raw, err := s.client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotExists.HpcCluster"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("EcsHpcCluster", id)), NotFoundMsg, ProviderERROR)
			return resp, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
		return resp, err
	}
	addDebug(action, raw, request)
	bresponse := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	//v, err := jsonpath.Get("$.HpcClusters.HpcCluster", raw)
	//if err != nil {
	//	return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HpcClusters.HpcCluster", raw)
	//}
	//if len(v.([]interface{})) < 1 {
	//	return object, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, raw)
	//} else {
	//	if v.([]interface{})[0].(map[string]interface{})["HpcClusterId"].(string) != id {
	//		return object, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, raw)
	//	}
	//}
	//object = v.([]interface{})[0].(map[string]interface{})
	return resp, nil
}
func (s *EcsService) DescribeEcsDeploymentSet(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewEcsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDeploymentSets"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"DeploymentSetIds": convertListToJsonString([]interface{}{id}),
		"Product":          "Ecs",
		"Department":       s.client.Department,
		"OrganizationId":   s.client.Department,
	}
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DeploymentSets.DeploymentSet", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DeploymentSets.DeploymentSet", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DeploymentSetId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}
