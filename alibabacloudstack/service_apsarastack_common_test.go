package alibabacloudstack

import (
	"fmt"
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
			errorValue := outValue[1]
			if !errorValue.IsNil() {
				err = errorValue.Interface().(error)
				if err != nil {
					if errmsgs.NotFoundError(err) {
						continue
					}
					return errmsgs.WrapError(err)
				}
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
			_, ok := ra.checkMap[rk]
			if rv == REMOVEKEY && ok {
				delete(ra.checkMap, rk)
			} else if ok {
				delete(ra.checkMap, rk)
				ra.checkMap[rk] = rv
			} else {
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
		return errmsgs.WrapError(fmt.Errorf(strings.Join(errorStrSlice, "\n")))
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

// in most cases, the TestCheckFunc list of dataSource test case is repeated，so we make an abstract in
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps:     steps,
	})
}

// according to configs generate step list and execute the test with preCheck
func (dsa *dataSourceAttr) dataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps:     steps,
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
func (conf *dataSourceTestAccConfig) buildDataSourceSteps(t *testing.T, info *dataSourceAttr, rand int) []resource.TestStep {
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
}

`

const DataAlibabacloudstackInstanceTypes = `
data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 2
  sorted_by         = "Memory"
  
 }

 locals {
  instance_type_id = sort(data.alibabacloudstack_instance_types.default.ids)[0]
}
 
`
const DataAlibabacloudstackInstanceTypes_Eni2 = `
data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  eni_amount        = 2
  sorted_by         = "Memory"
}

locals {
  instance_type_id = sort(data.alibabacloudstack_instance_types.default.ids)[0]
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
  resizeable_instance_id = length(local.resizeable_instance_ids ) > 0 ? local.resizeable_instance_ids[0]:  local.instance_type_id
}
`

const DataAlibabacloudstackImages = `
data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
  most_recent = true
  owners      = "system"
}

`

const EcsInstanceCommonNoZonesTestCase = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + `
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = data.alibabacloudstack_zones.default.zones[0].id
  name              = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
`

const EcsInstanceCommonTestCase = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + `
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
`
const RdsCommonTestCase = `
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
`
const PolarDBCommonTestCase = `
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}
data "alibabacloudstack_vswitches" "default" {
  zone_id = data.alibabacloudstack_zones.default.ids[0]
  is_default = "true"
}
`
const AdbCommonTestCase = `
resource "alibabacloudstack_vpc" "default" {
 name = "${var.name}"
 cidr_block = "172.16.0.0/16"
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
 cidr_block = "172.16.0.0/24"
}
`

const KVStoreCommonTestCase = `
data "alibabacloudstack_zones" "default" {

}
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
`

const DBMultiAZCommonTestCase = `
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
  multi = true
}
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.multi_zone_ids[0]}"
  name              = "${var.name}"
}
`

const ElasticsearchInstanceCommonTestCase = `

`

const SlbVpcCommonTestCase = `
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
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
    master_pwd = "ABCtest1234!"
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

const SlbListenerCommonTestCase = `
variable "ip_version" {
  default = "ipv4"
}	
resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  address_type = "internet"
}
`

const SlbListenerVserverCommonTestCase = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + `
resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${local.instance_type_id}"
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

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

`

const VpcCommonTestCase = `

resource "alibabacloudstack_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/24"
}
`

const VSwichCommonTestCase = DataZoneCommonTestCase + VpcCommonTestCase + `

resource "alibabacloudstack_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

`

const DBClusterCommonTestCase = VSwichCommonTestCase + `

resource "alibabacloudstack_adb_db_cluster" "cluster" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Cluster"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  pay_type            = "PostPaid"
  vswitch_id          = ${alibabacloudstack_vswitch.default.id}
  description         = "${var.name}_am"
}

`

const EipCommonTestCase = `

resource "alibabacloudstack_eip" "example" {
  bandwidth            = "10"
}

`

const SecurityGroupCommonTestCase = VSwichCommonTestCase + `

resource "alibabacloudstack_security_group" "group" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc.vpc.id}"
}

`

const ECSInstanceCommonTestCase = SecurityGroupCommonTestCase + `

resource "alibabacloudstack_instance" "instance" {
  image_id             = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_type        = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_security_group.group.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vswitch.default.id
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
}

`

const CbwpCommonTestCase = `

resource "alibabacloudstack_common_bandwidth_package" "foo" {
  bandwidth            = "200"
  name                 = "${var.name}_cbwp"
  description          = "test-common-bandwidth-package"
}

`

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

const FtbCommonTestCase = VSwichCommonTestCase + `

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
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

const DBInstanceCommonTestCase = `

resource "alibabacloudstack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  storage_type     = "local_ssd"
  instance_name        = "testacctf-mysql"
  tde_status=false
  enable_ssl=false
}

`

const KVRInstanceCommonTestCase = `

resource "alibabacloudstack_kvstore_instance" "default" {
  instance_class = "redis.master.small.default"
  instance_name  = "testacctf-redis"
  private_ip     = "172.16.0.10"
  security_ips   = ["10.0.0.1"]
  instance_type  = "Redis"
  cpu_type       = "intel"
  architecture_type = "standard"
}

`

const SlbCommonTestCase = VSwichCommonTestCase + `

resource "alibabacloudstack_slb" "default" {
  name          = "${var.name}_slb"
  vswitch_id    = "${alibabacloudstack_vswitch.default.id}"
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
