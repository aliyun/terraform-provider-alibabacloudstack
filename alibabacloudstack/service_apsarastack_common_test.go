package alibabacloudstack

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"log"
	"time"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

/**
	This file aims to provide some const test cases and applied them for several specified resource or data source's test cases.
These common test cases are used to creating some dependence resources, like vpc, vswitch and security group.
*/

// be used to check attribute map value
const (
	NOSET      = "#NOSET"       // be equivalent to method "TestCheckNoResourceAttrSet"
	CHECKSET   = "#CHECKSET"    // "TestCheckResourceAttrSet"
	REMOVEKEY  = "#REMOVEKEY"   // remove checkMap key
	REGEXMATCH = "#REGEXMATCH:" // "TestMatchResourceAttr" ,the map name/key like `"attribute" : REGEXMATCH + "attributeString"`
	ForceSleep = "force_sleep"
)

const (
	// indentation symbol
	INDENTATIONSYMBOL = " "

	// child field indend number
	CHILDINDEND = 2
)

// get a function that change checkMap pairs for a series test step
type resourceAttrMapUpdate func(map[string]string) resource.TestCheckFunc

// get a function that change attributeMap pairs for a series test step
type ResourceTestAccConfigFunc func(map[string]interface{}) string

// check the existence of resource
type resourceCheck struct {
	// IDRefreshName, like "alibabacloudstack_instance.foo"
	resourceId string

	// The response of the service method DescribeXXX
	resourceObject interface{}

	// The resource service client type, like DnsService, VpcService
	serviceFunc func() interface{}

	// service describe method name
	describeMethod string
}

func resourceCheckInit(resourceId string, resourceObject interface{}, serviceFunc func() interface{}) *resourceCheck {
	return &resourceCheck{
		resourceId:     resourceId,
		resourceObject: resourceObject,
		serviceFunc:    serviceFunc,
	}
}

func resourceCheckInitWithDescribeMethod(resourceId string, resourceObject interface{}, serviceFunc func() interface{}, describeMethod string) *resourceCheck {
	return &resourceCheck{
		resourceId:     resourceId,
		resourceObject: resourceObject,
		serviceFunc:    serviceFunc,
		describeMethod: describeMethod,
	}
}

// check attribute only
type resourceAttr struct {
	resourceId string
	checkMap   map[string]string
}

func resourceAttrInit(resourceId string, checkMap map[string]string) *resourceAttr {
	if checkMap == nil {
		checkMap = make(map[string]string)
	}
	return &resourceAttr{
		resourceId: resourceId,
		checkMap:   checkMap,
	}
}

// check the existence and attribute of the resource at the same time
type resourceAttrCheck struct {
	*resourceCheck
	*resourceAttr
}

func resourceAttrCheckInit(rc *resourceCheck, ra *resourceAttr) *resourceAttrCheck {
	return &resourceAttrCheck{
		resourceCheck: rc,
		resourceAttr:  ra,
	}
}

// check the resource existence by invoking DescribeXXX method of service and assign *resourceCheck.resourceObject value,
// the service is returned by invoking *resourceCheck.serviceFunc
func (rc *resourceCheck) checkResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		rs, ok := s.RootModule().Resources[rc.resourceId]
		if !ok {
			return errmsgs.WrapError(fmt.Errorf("can't find resource by id: %s", rc.resourceId))

		}
		outValue, err := rc.callDescribeMethod(rs)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		errorValue := outValue[1]
		if !errorValue.IsNil() {
			return errmsgs.WrapError(fmt.Errorf("Checking resource %s %s exists error:%s ", rc.resourceId, rs.Primary.ID, errorValue.Interface().(error).Error()))
		}
		if reflect.TypeOf(rc.resourceObject).Elem().String() == outValue[0].Type().String() {
			reflect.ValueOf(rc.resourceObject).Elem().Set(outValue[0])
			return nil
		} else {
			return errmsgs.WrapError(fmt.Errorf("The response object type expected *%s, got %s \n outValue: %v, \n resourceObject: %v",
				outValue[0].Type().String(), reflect.TypeOf(rc.resourceObject).String(), outValue, rc.resourceObject))
		}
	}
}

// check the resource destroy
func (rc *resourceCheck) checkResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceId, ".")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "alibabacloudstack_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		if resourceType == "" {
			return errmsgs.WrapError(errmsgs.Error("The resourceId %s is not correct and it should prefix with alibabacloudstack_", rc.resourceId))
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			outValue, err := rc.callDescribeMethod(rs)
			// if len(outValue) == 0 {
			// 	continue
			// }
			errorValue := outValue[1]

			if !errorValue.IsNil() {
				err = errorValue.Interface().(error)
				if err != nil {
					if errmsgs.NotFoundError(err) {
						continue
					}
					return errmsgs.WrapError(err)
				}
			} else if outValue[0].IsNil() {
				// Return empty, and no error is also considered as data not found, deletion successful
				continue
			} else {
				return errmsgs.WrapError(errmsgs.Error("the resource %s %s was not destroyed ! ", rc.resourceId, rs.Primary.ID))
			}
		}
		return nil
	}
}

// invoking DescribeXXX method of service
func (rc *resourceCheck) callDescribeMethod(rs *terraform.ResourceState) ([]reflect.Value, error) {
	var err error
	if rs.Primary.ID == "" {
		return nil, errmsgs.WrapError(fmt.Errorf("resource ID is not set"))
	}
	serviceP := rc.serviceFunc()
	if rc.describeMethod == "" {
		rc.describeMethod, err = getResourceDescribeMethod(rc.resourceId)
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
	}
	value := reflect.ValueOf(serviceP)
	typeName := value.Type().String()
	value = value.MethodByName(rc.describeMethod)
	if !value.IsValid() {
		return nil, errmsgs.WrapError(errmsgs.Error("The service type %s does not have method %s", typeName, rc.describeMethod))
	}
	inValue := []reflect.Value{reflect.ValueOf(rs.Primary.ID)}
	return value.Call(inValue), nil
}

func getResourceDescribeMethod(resourceId string) (string, error) {
	start := strings.Index(resourceId, "alibabacloudstack_")
	if start < 0 {
		return "", errmsgs.WrapError(fmt.Errorf("the parameter \"name\" don't contain string \"alibabacloudstack_\""))
	}
	start += len("alibabacloudstack_")
	end := strings.Index(resourceId[start:], ".") + start
	if end < 0 {
		return "", errmsgs.WrapError(fmt.Errorf("the parameter \"name\" don't contain string \".\""))
	}
	strs := strings.Split(resourceId[start:end], "_")
	describeName := "Describe"
	for _, str := range strs {
		describeName = describeName + strings.Title(str)
	}
	return describeName, nil
}

// check attribute func and check resource exist
func (rac *resourceAttrCheck) resourceAttrMapCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := rac.resourceCheck.checkResourceExists()(s)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		return rac.resourceAttr.resourceAttrMapCheck()(s)
	}
}

// execute the callback before check attribute and check resource exist
func (rac *resourceAttrCheck) resourceAttrMapCheckWithCallback(callback func()) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := rac.resourceCheck.checkResourceExists()(s)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		return rac.resourceAttr.resourceAttrMapCheckWithCallback(callback)(s)
	}
}

// get resourceAttrMapUpdate for a series test step and check resource exist
func (rac *resourceAttrCheck) resourceAttrMapUpdateSet() resourceAttrMapUpdate {
	return func(changeMap map[string]string) resource.TestCheckFunc {
		callback := func() {
			rac.updateCheckMapPair(changeMap)
		}
		return rac.resourceAttrMapCheckWithCallback(callback)
	}
}

// make a new map and copy from the old field checkMap, then update it according to the changeMap
func (ra *resourceAttr) updateCheckMapPair(changeMap map[string]string) {
	if interval, ok := changeMap[ForceSleep]; ok {
		intervalInt, err := strconv.Atoi(interval)
		if err == nil {
			time.Sleep(time.Duration(intervalInt) * time.Second)
			delete(changeMap, ForceSleep)
		}
	}
	newCheckMap := make(map[string]string, len(ra.checkMap))
	for k, v := range ra.checkMap {
		newCheckMap[k] = v
	}
	ra.checkMap = newCheckMap
	if changeMap != nil && len(changeMap) > 0 {
		for rk, rv := range changeMap {
			if _, ok := ra.checkMap[rk]; ok {
				delete(ra.checkMap, rk)
			}
			if rv != REMOVEKEY {
				ra.checkMap[rk] = rv
			}
		}
	}
}

// check attribute func
func (ra *resourceAttr) resourceAttrMapCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[ra.resourceId]
		if !ok {
			return errmsgs.WrapError(fmt.Errorf("can't find resource by id: %s", ra.resourceId))
		}
		if rs.Primary.ID == "" {
			return errmsgs.WrapError(fmt.Errorf("resource ID is not set"))
		}

		if ra.checkMap == nil || len(ra.checkMap) == 0 {
			return errmsgs.WrapError(fmt.Errorf("the parameter \"checkMap\" is nil or empty"))
		}

		var errorStrSlice []string
		errorStrSlice = append(errorStrSlice, "")
		for key, value := range ra.checkMap {
			var err error
			if strings.HasPrefix(value, REGEXMATCH) {
				var regex *regexp.Regexp
				regex, err = regexp.Compile(value[len(REGEXMATCH):])
				if err == nil {
					err = resource.TestMatchResourceAttr(ra.resourceId, key, regex)(s)
				} else {
					err = nil
				}
			} else if value == NOSET {
				err = resource.TestCheckNoResourceAttr(ra.resourceId, key)(s)
			} else if value == CHECKSET {
				err = resource.TestCheckResourceAttrSet(ra.resourceId, key)(s)
			} else {
				err = resource.TestCheckResourceAttr(ra.resourceId, key, value)(s)
			}
			if err != nil {
				errorStrSlice = append(errorStrSlice, err.Error())
			}
		}
		if len(errorStrSlice) == 1 {
			return nil
		}
		return errmsgs.WrapError(fmt.Errorf("%s", strings.Join(errorStrSlice, "\n")))
	}
}

// execute the callback before check attribute
func (ra *resourceAttr) resourceAttrMapCheckWithCallback(callback func()) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		callback()
		return ra.resourceAttrMapCheck()(s)
	}
}

// get resourceAttrMapUpdate for a series test step
func (ra *resourceAttr) resourceAttrMapUpdateSet() resourceAttrMapUpdate {
	return func(changeMap map[string]string) resource.TestCheckFunc {
		callback := func() {
			ra.updateCheckMapPair(changeMap)
		}
		return ra.resourceAttrMapCheckWithCallback(callback)
	}
}

func resourceTestAccConfigFunc(resourceId string,
	name string,
	configDependence func(name string) string) ResourceTestAccConfigFunc {
	basicInfo := resourceConfig{
		name:             name,
		resourceId:       resourceId,
		attributeMap:     make(map[string]interface{}),
		configDependence: configDependence,
	}
	return basicInfo.configBuild(false)
}

func dataSourceTestAccConfigFunc(resourceId string,
	name string,
	configDependence func(name string) string) ResourceTestAccConfigFunc {
	basicInfo := resourceConfig{
		name:             name,
		resourceId:       resourceId,
		attributeMap:     make(map[string]interface{}),
		configDependence: configDependence,
	}
	return basicInfo.configBuild(true)
}

// be used for generate testcase step config
type resourceConfig struct {
	// the resource name
	name string

	resourceId string

	// store attribute value that primary resource
	attributeMap map[string]interface{}

	// generate assistant test config
	configDependence func(name string) string
}

// according to changeMap to change the attributeMap value
func (b *resourceConfig) configUpdate(changeMap map[string]interface{}) {
	newMap := make(map[string]interface{}, len(b.attributeMap))
	for k, v := range b.attributeMap {
		newMap[k] = v
	}
	b.attributeMap = newMap
	if changeMap != nil && len(changeMap) > 0 {
		for rk, rv := range changeMap {
			_, ok := b.attributeMap[rk]
			if strValue, isCost := rv.(string); ok && isCost && strValue == REMOVEKEY {
				delete(b.attributeMap, rk)
			} else if ok {
				delete(b.attributeMap, rk)
				b.attributeMap[rk] = rv
			} else {
				b.attributeMap[rk] = rv
			}
		}
	}
}

// get BasicConfigFunc for resource a series test step
// overwrite: if true ,the attributeMap will be replace by changMap , other will be update
func (b *resourceConfig) configBuild(overwrite bool) ResourceTestAccConfigFunc {
	return func(changeMap map[string]interface{}) string {
		if overwrite {
			b.attributeMap = changeMap
		} else {
			b.configUpdate(changeMap)
		}
		strs := strings.Split(b.resourceId, ".")
		assistantConfig := b.configDependence(b.name)
		var primaryConfig string
		if strings.Compare("data", strs[0]) == 0 {
			primaryConfig = fmt.Sprintf("\n\ndata \"%s\" \"%s\" ", strs[1], strs[2])
		} else {
			primaryConfig = fmt.Sprintf("\n\nresource \"%s\" \"%s\" ", strs[0], strs[1])
		}
		return assistantConfig + primaryConfig + valueConvert(0, reflect.ValueOf(b.attributeMap))
	}
}

// deal with the parameter common method
func valueConvert(indentation int, val reflect.Value) string {
	switch val.Kind() {
	case reflect.Interface:
		return valueConvert(indentation, reflect.ValueOf(val.Interface()))
	case reflect.String:
		return fmt.Sprintf("\"%s\"", val.String())
	case reflect.Int:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Bool:
		return fmt.Sprintf("%v", val.Bool())
	case reflect.Slice:
		return listValue(indentation, val)
	case reflect.Map:
		return mapValue(indentation, val)
	default:
		log.Panicf("the map value must be string  map or slice type! %s", val)
	}
	return ""
}

// deal with list parameter
func listValue(indentation int, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation+CHILDINDEND)+
			valueConvert(indentation+CHILDINDEND, val.Index(i)))
	}

	return fmt.Sprintf("[\n%s\n%s]", strings.Join(valList, ",\n"), addIndentation(indentation))
}

// deal with map parameter
func mapValue(indentation int, val reflect.Value) string {
	var valList []string
	for _, keyV := range val.MapKeys() {
		mapVal := getRealValueType(val.MapIndex(keyV))
		var line string
		if mapVal.Kind() == reflect.Slice && mapVal.Len() > 0 {
			eleVal := getRealValueType(mapVal.Index(0))
			if eleVal.Kind() == reflect.Map {
				line = fmt.Sprintf(`%s%s`, addIndentation(indentation),
					listValueMapChild(indentation+CHILDINDEND, keyV.String(), mapVal))
				valList = append(valList, line)
				continue
			}
		}
		line = fmt.Sprintf(`%s%s = %s`, addIndentation(indentation+CHILDINDEND), keyV.String(),
			valueConvert(indentation+len(keyV.String())+CHILDINDEND+3, val.MapIndex(keyV)))
		valList = append(valList, line)
	}
	return fmt.Sprintf("{\n%s\n%s}", strings.Join(valList, "\n"), addIndentation(indentation))
}

// deal with list parameter that child element is map
func listValueMapChild(indentation int, key string, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation)+key+" "+
			mapValue(indentation, getRealValueType(val.Index(i))))
	}

	return fmt.Sprintf("%s\n%s", strings.Join(valList, "\n"), addIndentation(indentation))
}

func getRealValueType(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Interface:
		return getRealValueType(reflect.ValueOf(value.Interface()))
	default:
		return value
	}
}

func addIndentation(indentation int) string {
	return strings.Repeat(INDENTATIONSYMBOL, indentation)
}

// in most cases, the TestCheckFunc list of dataSource test case is repeated, so we make an abstract in
// order to reduce redundant code.
// dataSourceAttr has 3 field ,incloud resourceId  existMapFunc fakeMapFunc, every dataSource test can use only one
type dataSourceAttr struct {
	// IDRefreshName, like "data.alibabacloudstack_dns_records.record"
	resourceId string

	// get existMap function
	existMapFunc func(rand int) map[string]string

	// get fakeMap function
	fakeMapFunc func(rand int) map[string]string
}

// get exist and empty resourceAttrMapUpdate function
func (dsa *dataSourceAttr) checkDataSourceAttr(rand int) (exist, empty resourceAttrMapUpdate) {
	exist = resourceAttrInit(dsa.resourceId, dsa.existMapFunc(rand)).resourceAttrMapUpdateSet()
	empty = resourceAttrInit(dsa.resourceId, dsa.fakeMapFunc(rand)).resourceAttrMapUpdateSet()
	return
}

// according to configs generate step list and execute the test
func (dsa *dataSourceAttr) dataSourceTestCheck(t *testing.T, rand int, configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		Steps:             steps,
	})
}

// according to configs generate step list and execute the test with preCheck
func (dsa *dataSourceAttr) dataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	ResourceTest(t, resource.TestCase{
		PreCheck:          preCheck,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		Steps:             steps,
	})
}

// per schema attribute test config
type dataSourceTestAccConfig struct {
	// be equal to testCase config string,but the result has only one record
	existConfig string

	// if the dataSourceAttr.existMapFunc returned map value not match we want, existChangMap can alter checkMap for existConfig
	existChangMap map[string]string

	// be equal to testCase config string,but the result is empty
	fakeConfig string

	// if the dataSourceAttr.fakeMapFunc returned map value not match we want, fakeChangMap can alter checkMap for fakeConfig
	fakeChangMap map[string]string
}

// build test cases for each attribute
func (conf *dataSourceTestAccConfig) buildDataSourceSteps(_ *testing.T, info *dataSourceAttr, rand int) []resource.TestStep {
	testAccCheckExist, testAccCheckEmpty := info.checkDataSourceAttr(rand)
	var steps []resource.TestStep
	if conf.existConfig != "" {
		step := resource.TestStep{
			Config: conf.existConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckExist(conf.existChangMap),
			),
		}
		steps = append(steps, step)
	}
	if conf.fakeConfig != "" {
		step := resource.TestStep{
			Config: conf.fakeConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckEmpty(conf.fakeChangMap),
			),
		}
		steps = append(steps, step)
	}
	return steps
}

func (s *VpcService) needSweepVpc(vpcId, vswitchId string) (bool, error) {
	if vpcId == "" && vswitchId != "" {
		object, err := s.DescribeVSwitch(vswitchId)
		if err != nil && !errmsgs.NotFoundError(err) {
			return false, errmsgs.WrapError(err)
		}
		name := strings.ToLower(object.VSwitchName)
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the vswitch (%s (%s)).", object.VSwitchId, object.VSwitchName)
			return true, nil
		}
		vpcId = object.VpcId
	}
	if vpcId != "" {
		object, err := s.DescribeVpc(vpcId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return false, nil
			}
			return false, errmsgs.WrapError(err)
		}
		name := strings.ToLower(object.VpcName)
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the VPC (%s (%s)).", object.VpcId, object.VpcName)
			return true, nil
		}
	}
	return false, nil
}

func (s *VpcService) sweepVpc(id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Vpc %s ...", id)
	request := vpc.CreateDeleteVpcRequest()

	request.VpcId = id
	_, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteVpc(request)
	})

	return errmsgs.WrapError(err)
}

func (s *VpcService) sweepVSwitch(id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Vswitch %s ...", id)
	request := vpc.CreateDeleteVSwitchRequest()
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "vpc", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.VSwitchId = id
	_, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteVSwitch(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errmsgs.WrapError(err)
}

func (s *VpcService) sweepNatGateway(id string) error {
	if id == "" {
		return nil
	}

	log.Printf("[INFO] Deleting Nat Gateway %s ...", id)
	request := vpc.CreateDeleteNatGatewayRequest()
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.NatGatewayId = id
	request.Force = requests.NewBoolean(true)
	_, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteNatGateway(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errmsgs.WrapError(err)
}

func (s *EcsService) sweepSecurityGroup(id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Security Group %s ...", id)
	request := ecs.CreateDeleteSecurityGroupRequest()
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.SecurityGroupId = id
	_, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteSecurityGroup(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errmsgs.WrapError(err)
}

func (s *SlbService) sweepSlb(id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Set SLB DeleteProtection to off before deleting %s ...", id)
	request := slb.CreateSetLoadBalancerDeleteProtectionRequest()

	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "slb", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.LoadBalancerId = id
	request.DeleteProtection = "off"
	_, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.SetLoadBalancerDeleteProtection(request)
	})
	if err != nil {
		log.Printf("[ERROR] Set SLB %s DeleteProtection to off failed.", id)
	}
	log.Printf("[DEBUG] Deleting SLB %s ...", id)
	delRequest := slb.CreateDeleteLoadBalancerRequest()

	delRequest.Headers = map[string]string{"RegionId": s.client.RegionId}
	delRequest.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "slb", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	delRequest.LoadBalancerId = id
	_, err = s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DeleteLoadBalancer(delRequest)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errmsgs.WrapError(err)
}

const DataAlibabacloudstackVswitchZones = `
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

`

const DataAlibabacloudstackInstanceTypes = `
data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  sorted_by         = "CPU"
}

data "alibabacloudstack_instance_types" "default" {
  count = 8  # Traverse 1-8 core CPU configurations

  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = count.index + 1  # 1-8
  sorted_by            = "Memory"
}

locals {
  filtered_default = [for d in data.alibabacloudstack_instance_types.default : d if length(d.ids) > 0]
  fallback_all     = length(data.alibabacloudstack_instance_types.all.ids) > 0 ? data.alibabacloudstack_instance_types.all.ids : []
  
  default_instance_type_id = coalesce(
    try(local.filtered_default[0].ids[0], null),
    try(local.fallback_all[0], null),
    "no-available-instance-type"
  )
}
`

const DataAlibabacloudstackInstanceTypes_Eni2 = `
data "alibabacloudstack_instance_types" "eni2" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  eni_amount        = 2
  sorted_by         = "Memory"
}
`

const DataAlibabacloudstackResizeableInstanceTypes = `

locals {
  resizeable_instance_type_families = toset(["ecs.e4","ecs.e4v2","ecs.mn4","ecs.mn4v2","ecs.n4","ecs.n4v2","ecs.xn4","ecs.xn4v2"])
}

data "alibabacloudstack_instance_types" "resizeable"{
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 2
  sorted_by         = "Memory"
  for_each = local.resizeable_instance_type_families 
  instance_type_family = each.key
 }

locals {
  resizeable_instance_ids = [for type_name, type in data.alibabacloudstack_instance_types.resizeable : type.instance_types.0.id if length(type.instance_types ) > 0 ]
  resizeable_instance_id = length(local.resizeable_instance_ids ) > 0 ? local.resizeable_instance_ids[0]:  local.default_instance_type_id
}
`

const DataAlibabacloudstackImages = `
data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
  most_recent = true
  owners      = "system"
}

`

func RdsMysqlCommonTestCase() string {
	return fmt.Sprintf(`
variable "rds_instance_type" {
  type      = string
  default   = "%s"
  sensitive = true
}

data "alibabacloudstack_rds_instance_types" "default" {
  ids                  = var.rds_instance_type != "" ? [var.rds_instance_type] : null
  engine               = "MySQL"
  engine_version       = "5.7"
  sorted_by            = "CPU"
  series               = "dual_ha"
}

resource "alibabacloudstack_db_instance" "default" {
  engine               = data.alibabacloudstack_rds_instance_types.default.instance_types.0.engine
  engine_version       = data.alibabacloudstack_rds_instance_types.default.instance_types.0.engine_version
  instance_type        = data.alibabacloudstack_rds_instance_types.default.instance_types.0.id
  instance_storage     = data.alibabacloudstack_rds_instance_types.default.instance_types.0.storage_min
  instance_charge_type = "Postpaid"
  instance_name        = "${var.name}"
  vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
  monitoring_period    = "60"
  storage_type         = data.alibabacloudstack_rds_instance_types.default.instance_types.0.storage_type
}
`, os.Getenv("ALIBABACLOUDSTACK_TEST_RDS_INSTNCE_TYPE"))
}

func PolarDBMysqlCommonTestCase(enableVpc bool) string {
	var vswtichId string
	if enableVpc {
		vswtichId = `vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"`
	}
	return fmt.Sprintf(`
variable "polardb_instance_type" {
  type      = string
  default   = "%s"
  sensitive = true
}

data "alibabacloudstack_polardb_instance_types" "default" {
  ids                  = var.polardb_instance_type != "" ? [var.polardb_instance_type] : null
  engine               = "MySQL"
  engine_version       = "5.7"
  sorted_by            = "CPU"
  series               = "dual_ha"
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
  engine               = "${data.alibabacloudstack_polardb_instance_types.default.instance_types.0.engine}"
  engine_version       = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.engine_version
  instance_type        = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.id
  instance_storage     = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.storage_min
  instance_name        = "${var.name}"
  %s
  storage_type         = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.storage_type
}
`, os.Getenv("ALIBABACLOUDSTACK_TEST_POLARDB_INSTNCE_TYPE"), vswtichId)
}

const PolardbxCommonTestCase = `
data "alibabacloudstack_polardbx_instance_types" "cn" {
	sorted_by = "CPU"
	spec_type = "CN"
}

data "alibabacloudstack_polardbx_instance_types" "dn" {
	sorted_by = "CPU"
	spec_type = "DN"
}

resource "alibabacloudstack_polardbx_instance" "default" {
	description    = "${var.name}"
	zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage        = 50
	vswitch_id     = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class  = "${data.alibabacloudstack_polardbx_instance_types.cn.instance_types.0.id}"
	cn_node_count  = "2"
	dn_node_class  = "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id}"
	dn_node_count  = "2"
}
`

func PolardbxReadOrCreateCommonTestCase() string {
	return fmt.Sprintf(`
variable "existed_polardbx_id" {
	type      = string
	default   = "%s"
}

data "alibabacloudstack_polardbx_instance_types" "cn" {
	sorted_by = "CPU"
	spec_type = "CN"
}

data "alibabacloudstack_polardbx_instance_types" "dn" {
	sorted_by = "CPU"
	spec_type = "DN"
}

data "alibabacloudstack_polardbx_instances" "default" {
	ids = var.existed_polardbx_id == "" ? [" ",] : ["${var.existed_polardbx_id}",]
}

resource "alibabacloudstack_polardbx_instance" "default" {
	count          = length(data.alibabacloudstack_polardbx_instances.default.polardbx_instances) == 0 ? 1 : 0
	zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage        = 50
	vswitch_id     = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class  = "${data.alibabacloudstack_polardbx_instance_types.cn.instance_types.0.id}"
	cn_node_count  = "2"
	dn_node_class  = "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id}"
	dn_node_count  = "2"
}

locals {
	polardbx_instance = length(data.alibabacloudstack_polardbx_instances.default.polardbx_instances) == 0 ? alibabacloudstack_polardbx_instance.default.0 : data.alibabacloudstack_polardbx_instances.default.polardbx_instances.0
}

`, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_POLARDBX_ID"))
}

func removeEOFMarkers(input string) string {
	startMarker := "<<EOF\n"
	if strings.HasPrefix(input, startMarker) {
		input = input[len(startMarker):]
	}

	endMarker := "EOF"
	if idx := strings.LastIndex(input, endMarker); idx != -1 {
		input = input[:idx] + input[idx+len(endMarker):]
	}

	return input
}

const AdbCommonTestCase = `
resource "alibabacloudstack_vpc" "default" {
 name = "${var.name}"
 cidr_block = "192.168.0.0/16"
}
data "alibabacloudstack_zones" "default" {
 available_resource_creation = "ADB"
}

data "alibabacloudstack_vswitches" "default" {
 vpc_id = "${alibabacloudstack_vpc.default.id}"
 zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_vswitch" "default" {
 name = "tf_testAccAdb_vpc"
 vpc_id = "${alibabacloudstack_vpc.default.id}"
 availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
 cidr_block = "192.168.0.0/16"
}
`

const DBMultiAZCommonTestCase = `
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
  multi = true
}
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "192.168.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "192.168.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.multi_zone_ids[0]}"
  name              = "${var.name}"
}
`

const EmrCommonTestCase = `
data "alibabacloudstack_emr_main_versions" "default" {
}

data "alibabacloudstack_emr_instance_types" "default" {
    destination_resource = "InstanceType"
    cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
    support_local_storage = false
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER", "CORE"]
}

data "alibabacloudstack_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alibabacloudstack_emr_instance_types.default.types.0.id
	zone_id = data.alibabacloudstack_emr_instance_types.default.types.0.zone_id
}

data "alibabacloudstack_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alibabacloudstack_emr_instance_types.default.types.0.id
	zone_id = data.alibabacloudstack_emr_instance_types.default.types.0.zone_id
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_emr_instance_types.default.types.0.zone_id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com", 
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
    description = "this is a role test."
    force = true
}
`

const EmrGatewayTestCase = `
data "alibabacloudstack_emr_main_versions" "default" {
}

data "alibabacloudstack_emr_instance_types" "default" {
    destination_resource = "InstanceType"
    cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
    support_local_storage = false
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER","CORE"]
}

data "alibabacloudstack_emr_instance_types" "gateway" {
    destination_resource = "InstanceType"
    cluster_type = "GATEWAY"
    support_local_storage = false
    instance_charge_type = "PostPaid"
    support_node_type = ["GATEWAY"]
}

data "alibabacloudstack_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alibabacloudstack_emr_instance_types.default.types.0.id
	zone_id = data.alibabacloudstack_emr_instance_types.default.types.0.zone_id
}

data "alibabacloudstack_emr_disk_types" "gateway_data_disk" {
	destination_resource = "DataDisk"
	cluster_type = "GATEWAY"
	instance_charge_type = "PostPaid"
	instance_type = data.alibabacloudstack_emr_instance_types.gateway.types.0.id
	zone_id = data.alibabacloudstack_emr_instance_types.gateway.types.0.zone_id
}

data "alibabacloudstack_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alibabacloudstack_emr_instance_types.default.types.0.id
	zone_id = data.alibabacloudstack_emr_instance_types.default.types.0.zone_id
}

data "alibabacloudstack_emr_disk_types" "gateway_system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = "GATEWAY"
	instance_charge_type = "PostPaid"
	instance_type = data.alibabacloudstack_emr_instance_types.gateway.types.0.id
	zone_id = data.alibabacloudstack_emr_instance_types.gateway.types.0.zone_id
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_emr_instance_types.default.types.0.zone_id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com", 
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
    description = "this is a role test."
    force = true
}
variable "password" {
}
resource "alibabacloudstack_emr_cluster" "default" {
    name = "${var.name}"

    emr_ver = data.alibabacloudstack_emr_main_versions.default.main_versions.0.emr_version

    cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0

    host_group {
        host_group_name = "master_group"
        host_group_type = "MASTER"
        node_count = "2"
        instance_type = data.alibabacloudstack_emr_instance_types.default.types.0.id
        disk_type = data.alibabacloudstack_emr_disk_types.data_disk.types.0.value
        disk_capacity = data.alibabacloudstack_emr_disk_types.data_disk.types.0.min > 160 ? data.alibabacloudstack_emr_disk_types.data_disk.types.0.min : 160
        disk_count = "1"
        sys_disk_type = data.alibabacloudstack_emr_disk_types.system_disk.types.0.value
		sys_disk_capacity = data.alibabacloudstack_emr_disk_types.system_disk.types.0.min > 160 ? data.alibabacloudstack_emr_disk_types.system_disk.types.0.min : 160
    }

	host_group {
        host_group_name = "core_group"
        host_group_type = "CORE"
        node_count = "2"
        instance_type = data.alibabacloudstack_emr_instance_types.default.types.0.id
        disk_type = data.alibabacloudstack_emr_disk_types.data_disk.types.0.value
        disk_capacity = data.alibabacloudstack_emr_disk_types.data_disk.types.0.min > 160 ? data.alibabacloudstack_emr_disk_types.data_disk.types.0.min : 160
        disk_count = "4"
        sys_disk_type = data.alibabacloudstack_emr_disk_types.system_disk.types.0.value
        sys_disk_capacity = data.alibabacloudstack_emr_disk_types.system_disk.types.0.min > 160 ? data.alibabacloudstack_emr_disk_types.system_disk.types.0.min : 160
    }

    high_availability_enable = true
    zone_id = data.alibabacloudstack_emr_instance_types.default.types.0.zone_id
    security_group_id = alibabacloudstack_security_group.default.id
    is_open_public_ip = true
    charge_type = "PostPaid"
    vswitch_id = alibabacloudstack_vswitch.default.id
    user_defined_emr_ecs_role = alibabacloudstack_ram_role.default.name
    ssh_enable = true
    master_pwd = var.password
}
`
const EmrLocalStorageTestCase = `
data "alibabacloudstack_emr_main_versions" "default" {
}

data "alibabacloudstack_emr_instance_types" "local_disk" {
    destination_resource = "InstanceType"
    cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
    support_local_storage = true
    instance_charge_type = "PostPaid"
    support_node_type = ["CORE"]
}

data "alibabacloudstack_emr_instance_types" "cloud_disk" {
    destination_resource = "InstanceType"
    cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER"]
    zone_id = data.alibabacloudstack_emr_instance_types.local_disk.types.0.zone_id
}

data "alibabacloudstack_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alibabacloudstack_emr_instance_types.cloud_disk.types.0.id
	zone_id = data.alibabacloudstack_emr_instance_types.cloud_disk.types.0.zone_id
}

data "alibabacloudstack_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = data.alibabacloudstack_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alibabacloudstack_emr_instance_types.cloud_disk.types.0.id
	zone_id = data.alibabacloudstack_emr_instance_types.cloud_disk.types.0.zone_id
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_emr_instance_types.cloud_disk.types.0.zone_id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com", 
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
    description = "this is a role test."
    force = true
}
`

const SlbListenerVserverCommonTestCase = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + SecurityGroupCommonTestCase + `
resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${alibabacloudstack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_slb_master_slave_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${alibabacloudstack_instance.default.0.id}"
      port = 80
      weight = 100
      server_type = "Master"
  }
  servers {
      server_id = "${alibabacloudstack_instance.default.1.id}"
      port = 80
      weight = 100
      server_type = "Slave"
  }
}
`

const DataZoneCommonTestCase = `
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

`

const NasCommonTestCase = `
data "alibabacloudstack_nas_zones" "default" {
}

resource "alibabacloudstack_nas_file_system" "default" {
  protocol_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type}"
  storage_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}"
  encrypt_type = "0"
  zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
  cluster_id ="${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}"
  description = "${var.name}"
}
`

const VpcCommonTestCase = `
resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
      ignore_changes = [
		secondary_cidr_blocks,
        tags
      ]
  }
}
`

func RandomPasswordTestCase(passwordLen, count int) string {
	return fmt.Sprintf(`
resource "random_password" "password" {
	count            = %d
	length           = %d
	special          = true
	override_special = "!@#$^&*()_"
	min_lower        = 1
	min_upper        = 1
	min_numeric      = 1
}`, count, passwordLen)
}

const VSwitchCommonTestCase = DataZoneCommonTestCase + VpcCommonTestCase + `
resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

`

const VpnGatewayCommonTestCase = VSwitchCommonTestCase + `
resource "alibabacloudstack_vpn_gateway" "default" {
 name                 = "${var.name}"
 vpc_id               = "${alibabacloudstack_vpc_vpc.default.id}"
 bandwidth            = 10
 instance_charge_type = "PostPaid"
 enable_ssl           = true
 enable_ipsec		  = true
 vswitch_id			  = "${alibabacloudstack_vpc_vswitch.default.id}"
}
`

const DBClusterCommonTestCase = VSwitchCommonTestCase + `
resource "alibabacloudstack_adb_db_cluster" "cluster" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Cluster"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  pay_type            = "PostPaid"
  vswitch_id          = ${alibabacloudstack_vpc_vswitch.default.id}
  description         = "${var.name}_am"
}

`

const EipCommonTestCase = `
resource "alibabacloudstack_eip" "example" {
  bandwidth            = "10"
}

`

const SecurityGroupCommonTestCase = VSwitchCommonTestCase + `
resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  	cidr_ip = "192.168.0.0/16"
}

`

const ECSInstanceCommonTestCase = SecurityGroupCommonTestCase + DataAlibabacloudstackImages + DataAlibabacloudstackInstanceTypes + `
resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

`

const CbwpCommonTestCase = `
resource "alibabacloudstack_common_bandwidth_package" "foo" {
  bandwidth            = "200"
  name                 = "${var.name}_cbwp"
  description          = "test-common-bandwidth-package"
}

`

func AckK8sCommonTestCase() string {
	return fmt.Sprintf(`
variable "existed_k8s_cluster_id" {
	default = "%s"
}

%s

%s

data "alibabacloudstack_cs_kubernetes_clusters" "default" {
	ids = var.existed_k8s_cluster_id == "" ? [] : [var.existed_k8s_cluster_id]
}

locals {
	create_count = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? 0 : 1
}

resource "alibabacloudstack_cs_kubernetes" "default" {
	count						= local.create_count
	name						= var.name
	version						= "1.30.7-aliyun.1"
	os_type						= "linux"
	platform					= "AliyunLinux"
	num_of_nodes				= "3"
	master_count				= "3"
	master_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
	master_instance_types		= ["ecs.n4v2.large","ecs.n4v2.large","ecs.n4v2.large"]
	master_disk_category		= "cloud_ssd"
	vpc_id						= "${alibabacloudstack_vpc_vpc.default.id}"
	worker_instance_types		= ["ecs.n4v2.large"]
	worker_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}"]
	worker_disk_category		= "cloud_ssd"
	password					= random_password.password.0.result
	pod_cidr					= "172.20.0.0/16"
	service_cidr				= "172.21.0.0/20"
	worker_disk_size			= "40"
	master_disk_size			= "40"
	slb_internet_enabled		= "true"
	security_group_id			= alibabacloudstack_ecs_securitygroup.default.id
	runtime	 {
		name	= "containerd"
		version	= "1.6.28"
	}
}

locals {
	k8s_cluster_id = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.ids.0 : alibabacloudstack_cs_kubernetes.default.0.id
	k8s_cluster_name = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.names.0 : alibabacloudstack_cs_kubernetes.default.0.name
}
`, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_K8S_ID"), SecurityGroupCommonTestCase, RandomPasswordTestCase(12, 1))
}

func EdasClusterCommonTestCase() string {
	return AckK8sCommonTestCase() + `
	
data "alibabacloudstack_edas_clusters" "default" {
	name_regex = "^${local.k8s_cluster_name}$"
}

resource "alibabacloudstack_edas_k8s_cluster" "default" {
	count			= length(data.alibabacloudstack_edas_clusters.default.ids) > 0 ? 0 : 1
	cs_cluster_id	= local.k8s_cluster_id
}

locals {
	edas_cluster_id = length(data.alibabacloudstack_edas_clusters.default.ids) > 0 ? data.alibabacloudstack_edas_clusters.default.ids.0 : alibabacloudstack_edas_k8s_cluster.default.0.id
}
`
}

const VrtCommonTestCase = `

data "alibabacloudstack_express_connect_physical_connections" "nameRegex" {
	
}

resource "alibabacloudstack_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alibabacloudstack_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = "${var.name}_vrt"
  vlan_id                    = 1
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

`

const PcCommonTestCase = `

resource "alibabacloudstack_express_connect_physical_connection" "domestic" {
  device_name              = "express_connect_foo"
  access_point_id          = "ap-cn-hangzhou-yh-B"
  line_operator            = "CT"
  peer_location            = "${var.name}_pc"
  physical_connection_name = "${var.name}_pc"
  type                     = "VPC"
  description              = "my domestic connection"
  port_type                = "1000Base-LX"
  bandwidth                = 100
}

`

const FtbCommonTestCase = VSwitchCommonTestCase + `

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
  name   = "${var.name}"
}                                        

`

const DiskCommonTestCase = DataZoneCommonTestCase + `

resource "alibabacloudstack_disk" "disk" { 
  zone_id = data.alibabacloudstack_zones.default.zones.0.id
  name              = "New-disk"
  description       = "ECS-Disk"
  category          = "cloud_efficiency"
  size              = "30"
}

`

const KVRInstanceClassCommonTestCase = `
data "alibabacloudstack_zones" "kv_zone" {
  available_resource_creation = "KVStore"
  enable_details = true
}
 
data alibabacloudstack_kvstore_instance_classes "default" {
  edition_type = "${var.kv_edition}"
  engine = "${var.kv_engine}"
  sorted_by = "cpu"
  architecture = "cluster"
}
`

const SlbCommonTestCase = VSwitchCommonTestCase + `

resource "alibabacloudstack_slb_loadbalancer" "default" {
  name          = "${var.name}_slb"
  vswitch_id    = "${alibabacloudstack_vpc_vswitch.default.id}"
  address_type  = "internet"
  specification = "slb.s2.small"
}

`

const KeyCommonTestCase = `

resource "alibabacloudstack_kms_key" "key" {
  description             = "Hello KMS"
  pending_window_in_days  = "7"
  key_state               = "Enabled"
}
`

func ApiGatwayV2K8sInstanceTestCase(engineType, deployMode string) string {
	ingressClass := ""
	if deployMode == "apig_k8s" {
		ingressClass = `ingress_class_name = "${var.name}-class"`
	}
	return fmt.Sprintf(`
	data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
		sorted_by = "CPU"
	}

	%s

	data "alibabacloudstack_api_gateway_v2_k8s_clusters" "default" {
		k8s_cluster_name = local.k8s_cluster_name
	}

	resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
		count = length(data.alibabacloudstack_api_gateway_v2_k8s_clusters.default.ids) > 0 ? 0 : 1
		cs_cluster_id =   local.k8s_cluster_id
		k8s_cluster_name = local.k8s_cluster_name
	}

	resource "alibabacloudstack_api_gateway_v2_instance" "default" {
	  instance_name      = "${var.name}-apigw"
	  node_number        = 1
	  instance_class     = "mini"
	  broker_engine_type = "%s"
	  deploy_mode        = "%s"
	  deploy_cluster_code = length(data.alibabacloudstack_api_gateway_v2_k8s_clusters.default.ids) > 0 ? "${data.alibabacloudstack_api_gateway_v2_k8s_clusters.default.ids.0}" : "${alibabacloudstack_api_gateway_v2_k8s_cluster.default.0.id}"
	  deploy_cluster_namespace = "${var.name}-namespace"
	  %s
	  sls_enabled = "true"
	  prometheus_enabled = "true"
	}`, ApiGatwayV2K8sInstanceTestCase("apig_k8s", "SCG"), engineType, deployMode, ingressClass)
}

func ServerCertificateTestCase() string {
	if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_TEST")); err == nil && v {
		if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_SENSITIVE")); err == nil && v {
			return `-----BEGIN CERTIFICATE-----\nMIIDRjCCAq*******<Your Server Certificate String>*****bJJyOm5LqoiA=\n-----END CERTIFICATE-----`
		}
	}
	return `<<EOF
-----BEGIN CERTIFICATE-----
MIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV
BAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP
MA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0
ZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow
djELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE
ChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG
9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ
AoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB
coG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook
KOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw
HQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy
+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC
QkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN
MAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ
AJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT
cQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1
Ofi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd
DUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=
-----END CERTIFICATE-----
EOF
`
}

func RsaPrivateKeyTestCase() string {
	if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_TEST")); err == nil && v {
		if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_SENSITIVE")); err == nil && v {
			return `-----BEGIN RSA PRIVATE KEY-----\nMIIDRjCCAq******<Your RSA PRIVATE KEY String>******bJJyOm5LqoiA=\n-----END RSA PRIVATE KEY-----`
		}
	}
	return `<<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV
kg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM
ywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB
AoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd
6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP
hwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4
MdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz
71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm
Ev9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE
qygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8
9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM
zWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe
DrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=
-----END RSA PRIVATE KEY-----
EOF
`
}

func UpdateCertificateTestCase() string {
	if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_TEST")); err == nil && v {
		if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_SENSITIVE")); err == nil && v {
			return `-----BEGIN CERTIFICATE-----\nMIIDRjCCAq*******<Your Server Certificate String>*****bJJyOm5LqoiA=\n-----END CERTIFICATE-----`
		}
	}
	return `<<EOF
-----BEGIN CERTIFICATE-----
MIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV
BAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP
MA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0
ZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow
djELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE
ChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG
9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ
AoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB
coG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook
KOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw
HQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy
+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC
QkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN
MAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ
AJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT
cQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1
Ofi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd
DUBTxOAvSYfXTtI90EAxEG/bJJyOm5Lxxxx=
-----END CERTIFICATE-----
EOF
`
}

func UpdateRsaPrivateKeyTestCase() string {
	if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_TEST")); err == nil && v {
		if v, err := stringToBool(os.Getenv("ALIBABACLOUDSTACK_DRYRUN_SENSITIVE")); err == nil && v {
			return `-----BEGIN RSA PRIVATE KEY-----\nMIIDRjCCAq******<Your RSA PRIVATE KEY String>******bJJyOm5LqoiA=\n-----END RSA PRIVATE KEY-----`
		}
	}
	return `<<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV
kg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM
ywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB
AoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd
6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP
hwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4
MdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz
71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm
Ev9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE
qygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8
9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM
zWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe
DrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaxxx=
-----END RSA PRIVATE KEY-----
EOF
`
}
