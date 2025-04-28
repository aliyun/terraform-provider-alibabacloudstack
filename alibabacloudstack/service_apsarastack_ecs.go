package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/sdk_patch/datahub_patch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EcsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *EcsService) JudgeRegionValidation(key, region string) error {
	request := ecs.CreateDescribeRegionsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeRegions(request)
	})
	if err != nil {
		return fmt.Errorf("DescribeRegions got an error: %#v", err)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeRegionsResponse)
	if resp == nil || len(resp.Regions.Region) < 1 {
		return errmsgs.GetNotFoundErrorFromString("There is no any available region.")
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

func (s *EcsService) DescribeZone(id string) (zone ecs.Zone, err error) {
	request := ecs.CreateDescribeZonesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(request)
	})
	response, ok := raw.(*ecs.DescribeZonesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(response.Zones.Zone) < 1 {
		return zone, errmsgs.WrapError(errmsgs.Error("There is no any availability zone in region %s.", s.client.RegionId))
	}

	zoneIds := []string{}
	for _, z := range response.Zones.Zone {
		if z.ZoneId == id {
			return z, nil
		}
		zoneIds = append(zoneIds, z.ZoneId)
	}
	return zone, errmsgs.WrapError(errmsgs.Error("availability_zone %s not exists in region %s, all zones are %s", id, s.client.RegionId, zoneIds))
}

func (s *EcsService) DescribeZones(d *schema.ResourceData) (zones []ecs.Zone, err error) {
	request := ecs.CreateDescribeZonesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(request)
	})
	response, ok := raw.(*ecs.DescribeZonesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "apsarastak_instance_type_families", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(response.Zones.Zone) < 1 {
		return zones, errmsgs.WrapError(errmsgs.Error("There is no any availability zone in region %s.", s.client.RegionId))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		zoneIds := []string{}
		for _, z := range response.Zones.Zone {
			if z.ZoneId == v.(string) {
				return []ecs.Zone{z}, nil
			}
			zoneIds = append(zoneIds, z.ZoneId)
		}
		return zones, errmsgs.WrapError(errmsgs.Error("availability_zone %s not exists in region %s, all zones are %s", v.(string), s.client.RegionId, zoneIds))
	} else {
		return response.Zones.Zone, nil
	}
}

func (s *EcsService) DescribeInstance(id string) (instance ecs.Instance, err error) {
	request := ecs.CreateDescribeInstancesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceIds = convertListToJsonString([]interface{}{id})
	var raw interface{}
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		if err != nil {
			if errmsgs.IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	response, ok := raw.(*ecs.DescribeInstancesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	if len(response.Instances.Instance) < 1 {
		return instance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Instance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
	}
	log.Printf("[ECS Creation]: Getting Instance Details using DescribeInstances API: %s", response.Instances.Instance[0].Status)
	return response.Instances.Instance[0], nil
}

func (s *EcsService) DescribeInstanceAttribute(id string) (instance ecs.DescribeInstanceAttributeResponse, err error) {
	request := ecs.CreateDescribeInstanceAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceAttribute(request)
	})
	response, ok := raw.(*ecs.DescribeInstanceAttributeResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.InstanceId != id {
		return instance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Instance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
	}

	return *response, nil
}

func (s *EcsService) DescribeInstanceDisksByType(id string, rg string, disk_type string) (disks []ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id
	//request.DiskType = string(DiskTypeSystem)
	var response *ecs.DescribeDisksResponse
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(request)
		})
		if err != nil {
			if errmsgs.IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ecs.DescribeDisksResponse)
		if len(response.Disks.Disk) < 1 {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return disks, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	log.Printf("[ECS Creation]: Getting Disks Query Params : %s ", request.GetQueryParams())
	log.Printf("[ECS Creation]: Getting Disks response : %s ", response)
	//log.Printf("[ECS Creation]: Getting Disks Details: %s, Instance ID: %s, Id_to_compare: %s ",response.Disks.Disk[0],response.Disks.Disk[0].InstanceId,id)
	if len(response.Disks.Disk) < 1 || response.Disks.Disk[0].InstanceId != id {
		return disks, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Instance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
	}
	for _, diskdata := range response.Disks.Disk {
		if diskdata.InstanceId == id && diskdata.Type == string(disk_type) {
			disks = append(disks, diskdata)
		}
	}
	if len(disks) > 0 {
		return disks, nil
	}
	return disks, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Instance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
}

func (s *EcsService) ResourceAvailable(zone ecs.Zone, resourceType ResourceType) error {
	for _, res := range zone.AvailableResourceCreation.ResourceTypes {
		if res == string(resourceType) {
			return nil
		}
	}
	return errmsgs.WrapError(errmsgs.Error("%s is not available in %s zone of %s region", resourceType, zone.ZoneId, s.client.Region))
}

func (s *EcsService) DiskAvailable(zone ecs.Zone, diskCategory DiskCategory) error {
	for _, disk := range zone.AvailableDiskCategories.DiskCategories {
		if disk == string(diskCategory) {
			return nil
		}
	}
	return errmsgs.WrapError(errmsgs.Error("%s is not available in %s zone of %s region", diskCategory, zone.ZoneId, s.client.Region))
}

func (s *EcsService) JoinSecurityGroups(instanceId string, securityGroupIds []string) error {
	request := ecs.CreateJoinSecurityGroupRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	for _, sid := range securityGroupIds {
		request.SecurityGroupId = sid
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.JoinSecurityGroup(request)
		})
		response, ok := raw.(*ecs.JoinSecurityGroupResponse)
		if err != nil && errmsgs.IsExpectedErrors(err, []string{"InvalidInstanceId.AlreadyExists"}) {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *EcsService) LeaveSecurityGroups(instanceId string, securityGroupIds []string) error {
	request := ecs.CreateLeaveSecurityGroupRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	for _, sid := range securityGroupIds {
		request.SecurityGroupId = sid
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.LeaveSecurityGroup(request)
		})
		response, ok := raw.(*ecs.LeaveSecurityGroupResponse)
		if err != nil && errmsgs.IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *EcsService) DoEcsDescribesecuritygroupattributeRequest(id string) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	return s.DescribeSecurityGroup(id)
}
func (s *EcsService) DescribeSecurityGroup(id string) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.SecurityGroupId = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(request)
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			err = errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if response.SecurityGroupId != id {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Security Group", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
		return
	}

	return *response, nil
}

func (s *EcsService) DescribeSecurityGroupRule(id string) (rule ecs.Permission, err error) {
	parts, err := ParseResourceId(id, 8)
	if err != nil {
		return rule, errmsgs.WrapError(err)
	}
	groupId, direction, ipProtocol, portRange, nicType, cidr_ip, policy := parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6]
	priority, err := strconv.Atoi(parts[7])
	if err != nil {
		return rule, errmsgs.WrapError(err)
	}
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.SecurityGroupId = groupId
	request.Direction = direction
	request.NicType = nicType
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(request)
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			err = errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if response == nil {
		return rule, errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("Security Group", groupId))
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

	return rule, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Security Group Rule", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
}

func (s *EcsService) DescribeAvailableResources(d *schema.ResourceData, meta interface{}, destination DestinationResource) (zoneId string, validZones []ecs.AvailableZone, requestId string, err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateDescribeAvailableResourceRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
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
	response, ok := raw.(*ecs.DescribeAvailableResourceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return "", nil, "", errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	requestId = response.RequestId

	if len(response.AvailableZones.AvailableZone) < 1 {
		err = errmsgs.WrapError(errmsgs.Error("There are no availability resources in the region: %s. RequestId: %s.", client.RegionId, requestId))
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
			err = errmsgs.WrapError(errmsgs.Error("Availability zone %s status is not available in the region %s. Expected availability zones: %s. RequestId: %s.",
				zoneId, client.RegionId, strings.Join(expectedZones, ", "), requestId))
			return
		}
		if soldout {
			err = errmsgs.WrapError(errmsgs.Error("Availability zone %s status is sold out in the region %s. Expected availability zones: %s. RequestId: %s.",
				zoneId, client.RegionId, strings.Join(expectedZones, ", "), requestId))
			return
		}
	}

	if len(validZones) <= 0 {
		err = errmsgs.WrapError(errmsgs.Error("There is no availability resources in the region %s. Please choose another region. RequestId: %s.", client.RegionId, response.RequestId))
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
		return errmsgs.WrapError(errmsgs.Error("The instance type %s is solded out or is not supported in the zone %s. Expected instance types: %s", targetType, zoneId, strings.Join(expectedInstanceTypes, ", ")))
	}
	return errmsgs.WrapError(errmsgs.Error("The instance type %s is solded out or is not supported in the region %s. Expected instance types: %s", targetType, s.client.RegionId, strings.Join(expectedInstanceTypes, ", ")))
}

func (s *EcsService) QueryInstancesWithKeyPair(instanceIdsStr, keyPair string) (instanceIds []string, instances []ecs.Instance, err error) {
	request := ecs.CreateDescribeInstancesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.InstanceIds = instanceIdsStr
	request.KeyPairName = keyPair
	for {
		raw, e := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		object, ok := raw.(*ecs.DescribeInstancesResponse)
		if e != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(object.BaseResponse)
			}
			err = errmsgs.WrapErrorf(e, errmsgs.RequestV1ErrorMsg, keyPair, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(raw.(*responses.BaseResponse))
			}
			err = errmsgs.WrapErrorf(e, errmsgs.RequestV1ErrorMsg, keyPair, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return
		} else {
			request.PageNumber = page
		}
	}
	return
}

func (s *EcsService) DoEcsDescribekeypairsRequest(id string) (keyPair ecs.KeyPair, err error) {
	return s.DescribeKeyPair(id)
}
func (s *EcsService) DescribeKeyPair(id string) (keyPair ecs.KeyPair, err error) {
	request := ecs.CreateDescribeKeyPairsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.KeyPairName = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeKeyPairs(request)
	})
	object, ok := raw.(*ecs.DescribeKeyPairsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(object.BaseResponse)
		}
		return keyPair, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(object.KeyPairs.KeyPair) < 1 || object.KeyPairs.KeyPair[0].KeyPairName != id {
		return keyPair, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KeyPair", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, object.RequestId)
	}
	return object.KeyPairs.KeyPair[0], nil
}

func (s *EcsService) DescribeKeyPairAttachment(id string) (keyPair ecs.KeyPair, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
			err = errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return
	}
	keyPairName := parts[0]
	keyPair, err = s.DescribeKeyPair(keyPairName)
	if err != nil {
		return keyPair, errmsgs.WrapError(err)
	}
	if keyPair.KeyPairName != keyPairName {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KeyPairAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return keyPair, nil
}

func (s *EcsService) DoEcsDescribedisksRequest(id string) (disk ecs.Disk, err error) {
	return s.DescribeDisk(id)
}
func (s *EcsService) DescribeDisk(id string) (disk ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DiskIds = convertListToJsonString([]interface{}{id})
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(request)
	})
	response, ok := raw.(*ecs.DescribeDisksResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return disk, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	if len(response.Disks.Disk) < 1 || response.Disks.Disk[0].DiskId != id {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Disk", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return response.Disks.Disk[0], nil
}

func (s *EcsService) DescribeDiskAttachment(id string) (disk ecs.Disk, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return disk, errmsgs.WrapError(err)
	}
	disk, err = s.DescribeDisk(parts[0])
	if err != nil {
		return disk, errmsgs.WrapError(err)
	}

	if disk.InstanceId != parts[1] && disk.Status != string(InUse) {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DiskAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return
}

func (s *EcsService) DescribeDisksByType(instanceId string, diskType DiskType) (disk []ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	if instanceId != "" {
		request.InstanceId = instanceId
	}
	request.DiskType = string(diskType)

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(request)
	})
	resp, ok := raw.(*ecs.DescribeDisksResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeDisks", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if resp == nil {
		return
	}
	return resp.Disks.Disk, nil
}

func (s *EcsService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []ecs.Tag, err error) {
	request := ecs.CreateDescribeTagsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ResourceType = string(resourceType)
	request.ResourceId = resourceId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeTags(request)
	})
	response, ok := raw.(*ecs.DescribeTagsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resourceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(response.Tags.Tag) < 1 {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Tags", resourceId)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return
	}

	return response.Tags.Tag, nil
}

func (s *EcsService) DescribeImageById(id string) (image ecs.Image, err error) {
	request := ecs.CreateDescribeImagesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ImageId = id
	request.ImageOwnerAlias = "self"
	request.Status = fmt.Sprintf("%s,%s,%s,%s,%s", "Creating", "Waiting", "Available", "UnAvailable", "CreateFailed")
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImages(request)
	})
	resp, ok := raw.(*ecs.DescribeImagesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeImage", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if resp == nil || len(resp.Images.Image) < 1 {
		return image, errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("Image", id))
	}

	return resp.Images.Image[0], nil
}

func (s *EcsService) deleteImage(d *schema.ResourceData) error {
	object, err := s.DescribeImageById(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	request := ecs.CreateDeleteImageRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	if force, ok := d.GetOk("force"); ok {
		request.Force = requests.NewBoolean(force.(bool))
	}
	request.ImageId = object.ImageId

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteImage(request)
	})
	response, ok := raw.(*ecs.DeleteImageResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 3*time.Second, s.ImageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func (s *EcsService) updateImage(d *schema.ResourceData) error {
	d.Partial(true)

	err := setTags(s.client, TagResourceImage, d)
	if err != nil {
		return errmsgs.WrapError(err)
	} else {
		//d.SetPartial("tags")
	}

	request := ecs.CreateModifyImageAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ImageId = d.Id()

	if d.HasChanges("description", "name", "image_name") {
		if description, ok := d.GetOk("description"); ok {
			request.Description = description.(string)
		}
		if imageName, ok := connectivity.GetResourceDataOk(d, "image_name", "name"); ok {
			request.ImageName = imageName.(string)
		}
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyImageAttribute(request)
		})
		response, ok := raw.(*ecs.ModifyImageAttributeResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		//d.SetPartial("name")
		//d.SetPartial("image_name")
		//d.SetPartial("description")
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return nil
}

func (s *EcsService) DoEcsDescribenetworkinterfacesRequest(id string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	return s.DescribeNetworkInterface(id)
}

func (s *EcsService) DescribeNetworkInterface(id string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	request := ecs.CreateDescribeNetworkInterfacesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	eniIds := []string{id}
	request.NetworkInterfaceId = &eniIds
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeNetworkInterfaces(request)
	})
	response, ok := raw.(*ecs.DescribeNetworkInterfacesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return networkInterface, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var found = bool(false)
	var result = &ecs.DescribeNetworkInterfacesResponse{}
	for _, k := range response.NetworkInterfaceSets.NetworkInterfaceSet {
		if k.NetworkInterfaceId == id {
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
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NetworkInterface", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
		return
	}

	return response.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
}

func (s *EcsService) DescribeNetworkInterfaceAttachment(id string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	parts, err := ParseResourceId(id, 2)

	if err != nil {
		return networkInterface, errmsgs.WrapError(err)
	}
	eniId, instanceId := parts[0], parts[1]
	request := ecs.CreateDescribeNetworkInterfacesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	eniIds := []string{eniId}
	request.NetworkInterfaceId = &eniIds
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeNetworkInterfaces(request)
	})
	response, ok := raw.(*ecs.DescribeNetworkInterfacesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return networkInterface, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var found = bool(false)
	var result = &ecs.DescribeNetworkInterfacesResponse{}
	for _, k := range response.NetworkInterfaceSets.NetworkInterfaceSet {
		if k.NetworkInterfaceId == eniId {
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
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NetworkInterfaceAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
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
		if err != nil && !errmsgs.NotFoundError(err) {
			return err
		}
		if instance.Status == string(status) {
			//Sleep one more time for timing issues
			time.Sleep(DefaultIntervalMedium * time.Second)
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return errmsgs.GetTimeErrorFromString(errmsgs.GetTimeoutMessage("ECS Instance", string(status)))
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
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
			}
		}

		return object, object.Status, nil
	}
}

func (s *EcsService) deleteImageforDest(d *schema.ResourceData, region string) error {
	object, err := s.DescribeImageById(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	request := ecs.CreateDeleteImageRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	if force, ok := d.GetOk("force"); ok {
		request.Force = requests.NewBoolean(force.(bool))
	}
	request.RegionId = region
	request.ImageId = object.ImageId

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteImage(request)
	})
	response, ok := raw.(*ecs.DeleteImageResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 3*time.Second, s.ImageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func (s *EcsService) DescribeImage(id, region string) (image ecs.Image, err error) {
	request := ecs.CreateDescribeImagesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RegionId = region
	request.ImageId = id
	request.ImageOwnerAlias = "self"
	request.Status = fmt.Sprintf("%s,%s,%s,%s,%s", "Creating", "Waiting", "Available", "UnAvailable", "CreateFailed")
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImages(request)
	})
	log.Printf("[DEBUG] status %#v", raw)
	if err != nil {
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, ok := raw.(*ecs.DescribeImagesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return image, err
	}
	if resp == nil || len(resp.Images.Image) < 1 {
		return image, errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("Image", id))
	}
	return resp.Images.Image[0], nil
}

func (s *EcsService) ImageStateRefreshFuncforcopy(id string, region string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeImage(id, region)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
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
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		// Disk need 3-5 seconds to get ExpiredTime after the status is available
		if object.Status == string(status) && object.ExpiredTime != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (s *EcsService) WaitForSecurityGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeSecurityGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForKeyPair(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeKeyPair(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForDiskAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDiskAttachment(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (s *EcsService) WaitForNetworkInterface(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeNetworkInterface(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), errmsgs.ProviderERROR)
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
			return errmsgs.WrapError(errmsgs.Error("Wait for VPC attributes changed timeout"))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		instance, err := s.DescribeInstance(instanceId)
		if err != nil {
			return errmsgs.WrapError(err)
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
			return errmsgs.WrapError(err)
		}
		if object.InnerAccessPolicy == target {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.InnerAccessPolicy, target, errmsgs.ProviderERROR)
		}
	}
}

func (s *EcsService) AttachKeyPair(keyName string, instanceIds []interface{}) (err error) {
	request := ecs.CreateAttachKeyPairRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.KeyPairName = keyName
	request.InstanceIds = convertListToJsonString(instanceIds)
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachKeyPair(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"ServiceUnavailable"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	resp, ok := raw.(*ecs.AttachKeyPairResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, keyName, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return err
	}
	return nil
}

func (s *EcsService) SnapshotStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSnapshot(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *EcsService) ImageStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeImageById(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *EcsService) TaskStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTaskById(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		for _, failState := range failStates {
			if object.TaskStatus == failState {
				return object, object.TaskStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.TaskStatus))
			}
		}
		return object, object.TaskStatus, nil
	}
}

func (s *EcsService) DescribeTaskById(id string) (task *ecs.DescribeTaskAttributeResponse, err error) {
	request := ecs.CreateDescribeTaskAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.TaskId = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeTaskAttribute(request)
	})
	task, ok := raw.(*ecs.DescribeTaskAttributeResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(task.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return task, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if task.TaskId == "" {
		return task, errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("task", id))
	}
	return task, nil
}

func (s *EcsService) DoEcsDescribesnapshotsRequest(id string) (*ecs.Snapshot, error) {
	return s.DescribeSnapshot(id)
}
func (s *EcsService) DescribeSnapshot(id string) (*ecs.Snapshot, error) {
	snapshot := &ecs.Snapshot{}
	request := ecs.CreateDescribeSnapshotsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.SnapshotIds = fmt.Sprintf("[\"%s\"]", id)
	request.QueryParams["SnapshotId"] = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSnapshots(request)
	})
	response, ok := raw.(*ecs.DescribeSnapshotsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return snapshot, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	for _, k := range response.Snapshots.Snapshot {
		if k.SnapshotId == id {
			return &k, nil
		}
	}
	if len(response.Snapshots.Snapshot) < 1 || response.Snapshots.Snapshot[0].SnapshotId != id {
		return snapshot, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Snapshot", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
	}
	return &response.Snapshots.Snapshot[0], nil
}

func (s *EcsService) DoEcsDescribeautosnapshotpolicyexRequest(id string) (*ecs.AutoSnapshotPolicy, error) {
	return s.DescribeSnapshotPolicy(id)
}
func (s *EcsService) DescribeSnapshotPolicy(id string) (*ecs.AutoSnapshotPolicy, error) {
	policy := &ecs.AutoSnapshotPolicy{}
	request := ecs.CreateDescribeAutoSnapshotPolicyExRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.AutoSnapshotPolicyId = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAutoSnapshotPolicyEx(request)
	})
	response, ok := raw.(*ecs.DescribeAutoSnapshotPolicyExResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return policy, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(response.AutoSnapshotPolicies.AutoSnapshotPolicy) != 1 ||
		response.AutoSnapshotPolicies.AutoSnapshotPolicy[0].AutoSnapshotPolicyId != id {
		return policy, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SnapshotPolicy", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
	}

	return &response.AutoSnapshotPolicies.AutoSnapshotPolicy[0], nil
}

func (s *EcsService) DoEcsDescribereservedinstancesRequest(id string) (reservedInstance ecs.ReservedInstance, err error) {
	return s.DescribeReservedInstance(id)
}
func (s *EcsService) DescribeReservedInstance(id string) (reservedInstance ecs.ReservedInstance, err error) {
	request := ecs.CreateDescribeReservedInstancesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ReservedInstanceId = &[]string{id}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeReservedInstances(request)
	})
	response, ok := raw.(*ecs.DescribeReservedInstancesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return reservedInstance, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(response.ReservedInstances.ReservedInstance) != 1 ||
		response.ReservedInstances.ReservedInstance[0].ReservedInstanceId != id {
		return reservedInstance, errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("PurchaseReservedInstance", id))
	}
	return response.ReservedInstances.ReservedInstance[0], nil
}

func (s *EcsService) WaitForReservedInstance(id string, status Status, timeout int) error {
	deadLine := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		reservedInstance, err := s.DescribeReservedInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if reservedInstance.Status == string(status) {
			return nil
		}

		if time.Now().After(deadLine) {
			return errmsgs.WrapErrorf(errmsgs.GetTimeErrorFromString("ECS WaitForSnapshotPolicy"), errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, reservedInstance.Status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) WaitForSnapshotPolicy(id string, status Status, timeout int) error {
	deadLine := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		snapshotPolicy, err := s.DescribeSnapshotPolicy(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errmsgs.WrapError(err)
		}

		if snapshotPolicy.Status == string(status) {
			return nil
		}

		if time.Now().After(deadLine) {
			return errmsgs.WrapErrorf(errmsgs.GetTimeErrorFromString("ECS WaitForSnapshotPolicy"), errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, snapshotPolicy.Status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) DescribeLaunchTemplate(id string) (set ecs.LaunchTemplateSet, err error) {

	request := ecs.CreateDescribeLaunchTemplatesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.LaunchTemplateId = &[]string{id}

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeLaunchTemplates(request)
	})
	response, ok := raw.(*ecs.DescribeLaunchTemplatesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return set, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(response.LaunchTemplateSets.LaunchTemplateSet) != 1 ||
		response.LaunchTemplateSets.LaunchTemplateSet[0].LaunchTemplateId != id {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LaunchTemplate", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
		return
	}

	return response.LaunchTemplateSets.LaunchTemplateSet[0], nil

}

func (s *EcsService) DescribeLaunchTemplateVersion(id string, version int) (set ecs.LaunchTemplateVersionSet, err error) {

	request := ecs.CreateDescribeLaunchTemplateVersionsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.LaunchTemplateId = id
	request.LaunchTemplateVersion = &[]string{strconv.FormatInt(int64(version), 10)}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeLaunchTemplateVersions(request)
	})
	response, ok := raw.(*ecs.DescribeLaunchTemplateVersionsResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidLaunchTemplate.NotFound"}) {
			err = errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			return set, err
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return set, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(response.LaunchTemplateVersionSets.LaunchTemplateVersionSet) != 1 ||
		response.LaunchTemplateVersionSets.LaunchTemplateVersionSet[0].LaunchTemplateId != id {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LaunchTemplateVersion", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
		return
	}

	return response.LaunchTemplateVersionSets.LaunchTemplateVersionSet[0], nil

}

func (s *EcsService) WaitForLaunchTemplate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeLaunchTemplate(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.LaunchTemplateId == id && string(status) != string(Deleted) {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) DescribeImageShareByImageId(id string) (imageShare *ecs.DescribeImageSharePermissionResponse, err error) {
	request := ecs.CreateDescribeImageSharePermissionRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return imageShare, errmsgs.WrapError(err)
	}
	request.ImageId = parts[0]
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImageSharePermission(request)
	})
	resp, ok := raw.(*ecs.DescribeImageSharePermissionResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidImageId.NotFound"}) {
			return imageShare, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return imageShare, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(resp.Accounts.Account) == 0 {
		return imageShare, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ModifyImageSharePermission", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, resp.RequestId)
	}
	return resp, nil
}

func (s *EcsService) WaitForAutoProvisioningGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeAutoProvisioningGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), errmsgs.ProviderERROR)
		}
	}
}

func (s *EcsService) DescribeAutoProvisioningGroup(id string) (group ecs.AutoProvisioningGroup, err error) {
	request := ecs.CreateDescribeAutoProvisioningGroupsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.AutoProvisioningGroupId = &[]string{id}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAutoProvisioningGroups(request)
	})
	response, ok := raw.(*ecs.DescribeAutoProvisioningGroupsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return group, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range response.AutoProvisioningGroups.AutoProvisioningGroup {
		if v.AutoProvisioningGroupId == id {
			return v, nil
		}
	}
	err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AutoProvisioningGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
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
	filter := []string{"^aliyun", "^acs:", "^ascm:", "^http://", "^https://"}
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
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		response, ok := raw.(*ecs.UntagResourcesResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "InvalidResourceId.NotFound", "InvalidResourceType.NotFound", "MissingParameter.RegionId", "MissingParameter.ResourceIds", "MissingParameter.ResourceType", "MissingParameter.TagOwnerBid", "MissingParameter.TagOwnerUid", "MissingParameter.Tags"}) {
				return nil
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return err
		}
	}
	if len(added) > 0 {
		request := ecs.CreateTagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		response, ok := raw.(*ecs.TagResourcesResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "InvalidResourceId.NotFound", "InvalidResourceType.NotFound", "MissingParameter.RegionId", "MissingParameter.ResourceIds", "MissingParameter.ResourceType", "MissingParameter.TagOwnerBid", "MissingParameter.TagOwnerUid", "MissingParameter.Tags"}) {
				return nil
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return err
		}
	}
	return nil
}

func (s *EcsService) SetResourceTagsNew(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			_, err := s.client.DoTeaRequest("POST", "Ecs", "2019-05-10", action, "", nil, nil, request)
			if err != nil {
				return err
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			_, err := s.client.DoTeaRequest("POST", "Ecs", "2019-05-10", action, "", nil, nil, request)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *EcsService) DoEcsDescribededicatedhostautorenewRequest(id string) (object ecs.DedicatedHost, err error) {
	return s.DescribeEcsDedicatedHost(id)
}
func (s *EcsService) DescribeEcsDedicatedHost(id string) (object ecs.DedicatedHost, err error) {
	request := ecs.CreateDescribeDedicatedHostsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	for {

		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDedicatedHosts(request)
		})
		response, ok := raw.(*ecs.DescribeDedicatedHostsResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidLockReason.NotFound"}) {
				err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EcsDedicatedHost", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
				return object, err
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if len(response.DedicatedHosts.DedicatedHost) < 1 {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EcsDedicatedHost", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, response.RequestId)
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
			return object, errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EcsDedicatedHost", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, "")
	return
}

func (s *EcsService) EcsDedicatedHostStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcsDedicatedHost(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *EcsService) DoEcsDescribestoragesetdetailsRequest(id string) (result *datahub_patch.EcsDescribeEcsEbsStorageSetsResult, err error) {
	return s.DescribeEcsEbsStorageSet(id)
}
func (s *EcsService) DescribeEcsEbsStorageSet(id string) (result *datahub_patch.EcsDescribeEcsEbsStorageSetsResult, err error) {

	resp := &datahub_patch.EcsDescribeEcsEbsStorageSetsResult{}
	request := s.client.NewCommonRequest("GET", "Ecs", "2014-05-26", "DescribeStorageSets", "")

	request.QueryParams["PageNumber"] = "1"
	request.QueryParams["PageSize"] = "20"
	//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	raw, err := s.client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	addDebug("DescribeStorageSets", raw, request)
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "Operation.Forbidden"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EcsEbsStorageSet", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return resp, err
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return resp, err
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	//v, err := jsonpath.Get("$.Commands.Command", bresponse)
	if err != nil {
		return resp, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Commands.Command", raw)
	}
	//if len(v.([]interface{})) < 1 {
	//	return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ECS", id)), errmsgs.NotFoundWithResponse, raw)
	//} else {
	//	if v.([]interface{})[0].(map[string]interface{})["CommandId"].(string) != id {
	//		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ECS", id)), errmsgs.NotFoundWithResponse, raw)
	//	}
	//}
	//object = v.([]interface{})[0].(map[string]interface{})
	return resp, nil
}

func (s *EcsService) DoEcsDescribecommandsRequest(id string) (result *datahub_patch.EcsDescribeEcsCommandResult, err error) {
	return s.DescribeEcsCommand(id)
}
func (s *EcsService) DescribeEcsCommand(id string) (result *datahub_patch.EcsDescribeEcsCommandResult, err error) {

	action := "DescribeCommands"
	resp := &datahub_patch.EcsDescribeEcsCommandResult{}
	request := s.client.NewCommonRequest("GET", "Ecs", "2014-05-26", action, "")

	mergeMaps(request.QueryParams, map[string]string{
		"CommandId":  id,
		"PageNumber": "1",
		"PageSize":   "20",
	})

	//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	raw, err := s.client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "Operation.Forbidden"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EcsCommand", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return resp, err
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return resp, err
	}
	addDebug(action, raw, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	//v, err := jsonpath.Get("$.Commands.Command", bresponse)
	if err != nil {
		return resp, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Commands.Command", raw)
	}
	//if len(v.([]interface{})) < 1 {
	//	return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ECS", id)), errmsgs.NotFoundWithResponse, raw)
	//} else {
	//	if v.([]interface{})[0].(map[string]interface{})["CommandId"].(string) != id {
	//		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ECS", id)), errmsgs.NotFoundWithResponse, raw)
	//	}
	//}
	//object = v.([]interface{})[0].(map[string]interface{})
	return resp, nil
}
func (s *EcsService) DoEcsDescribehpcclustersRequest(id string) (result *datahub_patch.EcsDescribeEcsHpcClusterResult, err error) {
	return s.DescribeEcsHpcCluster(id)
}
func (s *EcsService) DescribeEcsHpcCluster(id string) (result *datahub_patch.EcsDescribeEcsHpcClusterResult, err error) {
	//var response map[string]interface{}

	resp := &datahub_patch.EcsDescribeEcsHpcClusterResult{}
	action := "DescribeHpcClusters"
	ids, err := json.Marshal([]string{id})
	if err != nil {
		return nil, err
	}
	//request := map[string]interface{}{
	//
	//	"HpcClusterIds": string(ids),
	//}

	ClientToken := buildClientToken("DescribeHpcClusters")
	request := s.client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
	request.QueryParams["HpcClusterIds"] = string(ids)
	request.QueryParams["ClientToken"] = ClientToken

	raw, err := s.client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"NotExists.HpcCluster"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EcsHpcCluster", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return resp, err
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return resp, err
	}
	addDebug(action, raw, request)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	//v, err := jsonpath.Get("$.HpcClusters.HpcCluster", raw)
	//if err != nil {
	//	return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.HpcClusters.HpcCluster", raw)
	//}
	//if len(v.([]interface{})) < 1 {
	//	return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ECS", id)), errmsgs.NotFoundWithResponse, raw)
	//} else {
	//	if v.([]interface{})[0].(map[string]interface{})["HpcClusterId"].(string) != id {
	//		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ECS", id)), errmsgs.NotFoundWithResponse, raw)
	//	}
	//}
	//object = v.([]interface{})[0].(map[string]interface{})
	return resp, nil
}
func (s *EcsService) DoEcsDescribedeploymentsetsRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeEcsDeploymentSet(id)
}
func (s *EcsService) DescribeEcsDeploymentSet(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"DeploymentSetIds": convertListToJsonString([]interface{}{id}),
	}
	response, err = s.client.DoTeaRequest("POST", "Ecs", "2014-05-26", "DescribeDeploymentSets", "", nil, nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.DeploymentSets.DeploymentSet", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.DeploymentSets.DeploymentSet", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ECS", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DeploymentSetId"].(string) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ECS", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}
