package alibabacloudstack

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/google/uuid"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type InstanceNetWork string

const (
	ClassicNet = InstanceNetWork("classic")
	VpcNet     = InstanceNetWork("vpc")
)

type PayType string

const (
	PrePaid  = PayType("PrePaid")
	PostPaid = PayType("PostPaid")
	Prepaid  = PayType("Prepaid")
	Postpaid = PayType("Postpaid")
)

const (
	NormalMode = "normal"
	SafetyMode = "safety"
)

type DdosbgpInsatnceType string

const (
	Enterprise   = DdosbgpInsatnceType("Enterprise")
	Professional = DdosbgpInsatnceType("Professional")
)

type DdosbgpInstanceIpType string

const (
	IPv4 = DdosbgpInstanceIpType("IPv4")
	IPv6 = DdosbgpInstanceIpType("IPv6")
)

type NetType string

const (
	Internet = NetType("Internet")
	Intranet = NetType("Intranet")
)

type NetworkType string

const (
	Classic         = NetworkType("Classic")
	Vpc             = NetworkType("Vpc")
	ClassicInternet = NetworkType("classic_internet")
	ClassicIntranet = NetworkType("classic_intranet")
	PUBLIC          = NetworkType("PUBLIC")
	PRIVATE         = NetworkType("PRIVATE")
)

type NodeType string

const (
	WORKER = NodeType("WORKER")
	KIBANA = NodeType("KIBANA")
)

type ActionType string

const (
	OPEN  = ActionType("OPEN")
	CLOSE = ActionType("CLOSE")
)

type TimeType string

const (
	Hour  = TimeType("Hour")
	Day   = TimeType("Day")
	Week  = TimeType("Week")
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type IpVersion string

const (
	IPV4 = IpVersion("ipv4")
	IPV6 = IpVersion("ipv6")
)

type Status string

const (
	Pending     = Status("Pending")
	Creating    = Status("Creating")
	Running     = Status("Running")
	Available   = Status("Available")
	Unavailable = Status("Unavailable")
	Modifying   = Status("Modifying")
	Deleting    = Status("Deleting")
	Starting    = Status("Starting")
	Stopping    = Status("Stopping")
	Stopped     = Status("Stopped")
	Normal      = Status("Normal")
	Changing    = Status("Changing")
	Online      = Status("online")
	Configuring = Status("configuring")

	Associating   = Status("Associating")
	Unassociating = Status("Unassociating")
	InUse         = Status("InUse")
	DiskInUse     = Status("In_use")

	Active   = Status("Active")
	Inactive = Status("Inactive")
	Idle     = Status("Idle")

	SoldOut = Status("SoldOut")

	InService      = Status("InService")
	Removing       = Status("Removing")
	EnabledStatus  = Status("Enabled")
	DisabledStatus = Status("Disabled")

	Init            = Status("Init")
	Provisioning    = Status("Provisioning")
	Updating        = Status("Updating")
	FinancialLocked = Status("FinancialLocked")

	PUBLISHED   = Status("Published")
	NOPUBLISHED = Status("NonPublished")

	Deleted = Status("Deleted")
	Null    = Status("Null")

	Enable = Status("Enable")
	BINDED = Status("BINDED")
)

type IPType string

const (
	Inner   = IPType("Inner")
	Private = IPType("Private")
	Public  = IPType("Public")
)

type ResourceType string

const (
	ResourceTypeInstance      = ResourceType("Instance")
	ResourceTypeDisk          = ResourceType("Disk")
	ResourceTypeVSwitch       = ResourceType("VSwitch")
	ResourceTypeRds           = ResourceType("Rds")
	ResourceTypePolarDB       = ResourceType("PolarDB")
	IoOptimized               = ResourceType("IoOptimized")
	ResourceTypeRkv           = ResourceType("KVStore")
	ResourceTypeFC            = ResourceType("FunctionCompute")
	ResourceTypeElasticsearch = ResourceType("Elasticsearch")
	ResourceTypeSlb           = ResourceType("Slb")
	ResourceTypeMongoDB       = ResourceType("MongoDB")
	ResourceTypeGpdb          = ResourceType("Gpdb")
	ResourceTypeHBase         = ResourceType("HBase")
	ResourceTypeAdb           = ResourceType("ADB")
	ResourceTypeCassandra     = ResourceType("Cassandra")
)

type InternetChargeType string

const (
	PayByBandwidth = InternetChargeType("PayByBandwidth")
	PayByTraffic   = InternetChargeType("PayByTraffic")
	PayBy95        = InternetChargeType("PayBy95")
)

type AccountSite string

const (
	DomesticSite = AccountSite("Domestic")
	IntlSite     = AccountSite("International")
)

const (
	SnapshotCreatingInProcessing = Status("progressing")
	SnapshotCreatingAccomplished = Status("accomplished")
	SnapshotCreatingFailed       = Status("failed")

	SnapshotPolicyCreating  = Status("Creating")
	SnapshotPolicyAvailable = Status("available")
	SnapshotPolicyNormal    = Status("Normal")
)

// timeout for common product, ecs e.g.
const DefaultTimeout = 300
const Timeout5Minute = 300
const DefaultTimeoutMedium = 500

// timeout for long time progerss product, rds e.g.
const DefaultLongTimeout = 1000

const DefaultIntervalMini = 2

const DefaultIntervalShort = 5

const DefaultIntervalMedium = 10

const DefaultIntervalLong = 20

const (
	PageSizeSmall  = 10
	PageSizeMedium = 20
	PageSizeLarge  = 50
	PageSizeXLarge = 100
)

// Protocol represents network protocol
type Protocol string

// Constants of protocol definition
const (
	Http  = Protocol("http")
	Https = Protocol("https")
	Tcp   = Protocol("tcp")
	Udp   = Protocol("udp")
	All   = Protocol("all")
	Icmp  = Protocol("icmp")
	Gre   = Protocol("gre")
)

// ValidProtocols network protocol list
var ValidProtocols = []Protocol{Http, Https, Tcp, Udp}

// simple array value check method, support string type only
func isProtocolValid(value string) bool {
	res := false
	for _, v := range ValidProtocols {
		if string(v) == value {
			res = true
		}
	}
	return res
}

// default region for all resource
const DEFAULT_REGION = "cn-beijing"

const INT_MAX = 2147483647

// symbol of multiIZ
const MULTI_IZ_SYMBOL = "MAZ"

const COMMA_SEPARATED = ","

const COLON_SEPARATED = ":"

const SLASH_SEPARATED = "/"

const LOCAL_HOST_IP = "127.0.0.1"

func convertListStringToListInterface(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func addDebug(action, content interface{}, requestInfo ...interface{}) {
	if debugOn() {
		trace := "[DEBUG TRACE]:\n"
		for skip := 1; skip < 5; skip++ {
			_, filepath, line, _ := runtime.Caller(skip)
			trace += fmt.Sprintf("%s:%d\n", filepath, line)
		}

		if len(requestInfo) > 0 {
			var request = struct {
				Domain     string
				Version    string
				UserAgent  string
				ActionName string
				Method     string
				Product    string
				Region     string
				AK         string
			}{}
			switch requestInfo[0].(type) {
			case *requests.RpcRequest:
				tmp := requestInfo[0].(*requests.RpcRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *requests.RoaRequest:
				tmp := requestInfo[0].(*requests.RoaRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *requests.CommonRequest:
				tmp := requestInfo[0].(*requests.CommonRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *fc.Client:
				client := requestInfo[0].(*fc.Client)
				request.Version = client.Config.APIVersion
				request.Product = "FC"
				request.ActionName = fmt.Sprintf("%s", action)

			}

			requestContent := ""
			if len(requestInfo) > 1 {
				requestContent = fmt.Sprintf("%#v", requestInfo[1])
			}

			content = fmt.Sprintf("%vDomain:%v, Version:%v, ActionName:%v, Method:%v, Product:%v, Region:%v\n\n"+
				"*************** %s Request ***************\n%#v\n",
				content, request.Domain, request.Version, request.ActionName,
				request.Method, request.Product, request.Region, request.ActionName, requestContent)
		}

		//fmt.Printf(DefaultDebugMsg, action, content, trace)
		log.Printf(errmsgs.DefaultDebugMsg, action, content, trace)
	}
}

func debugOn() bool {
	for _, part := range strings.Split(os.Getenv("DEBUG"), ",") {
		if strings.TrimSpace(part) == "terraform" {
			return true
		}
	}
	return false
}

// Convert the result for an array and returns a Json string
func convertListToJsonString(configured []interface{}) string {
	if len(configured) < 1 {
		return ""
	}
	result := "["
	for i, v := range configured {
		result += "\"" + v.(string) + "\""
		if i < len(configured)-1 {
			result += ","
		}
	}
	result += "]"
	return result
}

func getNextpageNumber(number requests.Integer) (requests.Integer, error) {
	page, err := strconv.Atoi(string(number))
	if err != nil {
		return "", err
	}
	return requests.NewInteger(page + 1), nil
}

func incrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	//	迁移动作太大，使用重定向
	return connectivity.IncrementalWait(firstDuration, increaseDuration)
}

func GetFunc(level int) string {
	pc, _, _, ok := runtime.Caller(level)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in GetFuncName.")
		return ""
	}
	return strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")
}

func ParseResourceId(id string, length int) (parts []string, err error) {
	parts = strings.Split(id, ":")

	if len(parts) != length {
		err = errmsgs.WrapError(fmt.Errorf("Invalid Resource Id %s. Expected parts' length %d, got %d", id, length, len(parts)))
	}
	return parts, err
}

func BuildStateConf(pending, target []string, timeout, delay time.Duration, f resource.StateRefreshFunc) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    f,
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: 3 * time.Second,
	}
}
func BuildStateConfByTimes(pending, target []string, timeout, delay time.Duration, f resource.StateRefreshFunc, notFoundChecks int) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:        pending,
		Target:         target,
		Refresh:        f,
		Timeout:        timeout,
		Delay:          delay,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: notFoundChecks,
	}
}
func convertJsonStringToList(configured string) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(configured), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func StringPointer(s string) *string {
	return &s
}

func BoolPointer(b bool) *bool {
	return &b
}

func Int32Pointer(i int32) *int32 {
	return &i
}

func Int64Pointer(i int64) *int64 {
	return &i
}

func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type Catcher struct {
	Reason           string
	RetryCount       int
	RetryWaitSeconds int
}

var ClientErrorCatcher = Catcher{errmsgs.AlibabacloudStackGoClientFailure, 10, 5}
var ServiceBusyCatcher = Catcher{"ServiceUnavailable", 10, 5}
var ThrottlingCatcher = Catcher{errmsgs.Throttling, 50, 2}

func NewInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(ClientErrorCatcher)
	i.AddCatcher(ServiceBusyCatcher)
	i.AddCatcher(ThrottlingCatcher)
	return i
}

func userDataHashSum(user_data string) string {
	// Check whether the user_data is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(user_data)
	if base64DecodeError != nil {
		v = []byte(user_data)
	}
	return string(v)
}

const ServerSideEncryptionAes256 = "AES256"
const ServerSideEncryptionKMS = "KMS"

type OptimizedType string

const (
	IOOptimized   = OptimizedType("optimized")
	NoneOptimized = OptimizedType("none")
)

type TagResourceType string

const (
	TagResourceImage         = TagResourceType("image")
	TagResourceInstance      = TagResourceType("instance")
	TagResourceAcl           = TagResourceType("acl")
	TagResourceCertificate   = TagResourceType("certificate")
	TagResourceSnapshot      = TagResourceType("snapshot")
	TagResourceKeypair       = TagResourceType("keypair")
	TagResourceDisk          = TagResourceType("disk")
	TagResourceSecurityGroup = TagResourceType("securitygroup")
	TagResourceEni           = TagResourceType("eni")
	TagResourceCdn           = TagResourceType("DOMAIN")
	TagResourceVpc           = TagResourceType("VPC")
	TagResourceVSwitch       = TagResourceType("VSWITCH")
	TagResourceRouteTable    = TagResourceType("ROUTETABLE")
	TagResourceEip           = TagResourceType("EIP")
	TagResourcePlugin        = TagResourceType("plugin")
	TagResourceApiGroup      = TagResourceType("apiGroup")
	TagResourceApp           = TagResourceType("app")
	TagResourceTopic         = TagResourceType("topic")
	TagResourceConsumerGroup = TagResourceType("consumergroup")
	TagResourceCluster       = TagResourceType("cluster")
)

type KubernetesNodeType string

const (
	KubernetesNodeMaster = ResourceType("Master")
	KubernetesNodeWorker = ResourceType("Worker")
)

func GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Get current user got an error: %#v.", err)
	}
	return usr.HomeDir, nil
}

// writeToFile 函数
func writeToFile(filePath string, data interface{}) error {
	var out string
	switch v := data.(type) {
	case string:
		out = v
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %v", data, err)
		}
		out = string(bs)
	}

	// 替换 ~ 为用户主目录
	if strings.HasPrefix(filePath, "~") {
		home, err := GetUserHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// 获取用户主目录
	home, err := GetUserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	// 获取文件路径的绝对路径
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %v", filePath, err)
	}

	// 确保文件路径是相对于当前工作目录或用户主目录
	if !strings.HasPrefix(absFilePath, currentDir+string(filepath.Separator)) && !strings.HasPrefix(absFilePath, home+string(filepath.Separator)) {
		return fmt.Errorf("file path %s is not within the allowed directories: current directory %s or home directory %s", absFilePath, currentDir, home)
	}

	// 写入文件
	return ioutil.WriteFile(absFilePath, []byte(out), 0644)
}

type Invoker struct {
	catchers []*Catcher
}

func (a *Invoker) AddCatcher(catcher Catcher) {
	a.catchers = append(a.catchers, &catcher)
}

func (a *Invoker) Run(f func() error) error {
	err := f()

	if err == nil {
		return nil
	}

	for _, catcher := range a.catchers {
		if errmsgs.IsExpectedErrors(err, []string{catcher.Reason}) {
			catcher.RetryCount--

			if catcher.RetryCount <= 0 {
				return fmt.Errorf("Retry timeout and got an error: %#v.", err)
			} else {
				time.Sleep(time.Duration(catcher.RetryWaitSeconds) * time.Second)
				return a.Run(f)
			}
		}
	}
	return err
}
func Trim(v string) string {
	if len(v) < 1 {
		return v
	}
	return strings.Trim(v, " ")
}

func GetCenChildInstanceType(id string) (c string, e error) {
	if strings.HasPrefix(id, "vpc") {
		return ChildInstanceTypeVpc, nil
	} else if strings.HasPrefix(id, "vbr") {
		return ChildInstanceTypeVbr, nil
	} else if strings.HasPrefix(id, "ccn") {
		return ChildInstanceTypeCcn, nil
	} else {
		return c, fmt.Errorf("CEN child instance ID invalid. Now, it only supports VPC or VBR or CCN instance.")
	}
}
func ParseSlbListenerId(id string) (parts []string, err error) {
	parts = strings.Split(id, ":")
	if len(parts) != 2 && len(parts) != 3 {
		err = errmsgs.WrapError(fmt.Errorf("Invalid alibabacloudstack_slb_listener Id %s. Expected Id format is <slb id>:<protocol>:< frontend>.", id))
	}
	return parts, err
}

func buildClientToken(action string) string {
	token := strings.TrimSpace(fmt.Sprintf("TF-%s-%d-%s", action, time.Now().Unix(), strings.Trim(uuid.New().String(), "-")))
	if len(token) > 64 {
		token = token[0:64]
	}
	return token
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}

func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(int))
	}
	return vs
}

func computePeriodByUnit(createTime, endTime interface{}, currentPeriod int, periodUnit string) (int, error) {
	var createTimeStr, endTimeStr string
	switch value := createTime.(type) {
	case int64:
		createTimeStr = time.Unix(createTime.(int64), 0).Format(time.RFC3339)
		endTimeStr = time.Unix(endTime.(int64), 0).Format(time.RFC3339)
	case string:
		createTimeStr = createTime.(string)
		endTimeStr = endTime.(string)
	default:
		return 0, errmsgs.WrapError(fmt.Errorf("Unsupported time type: %#v", value))
	}
	// currently, there is time value does not format as standard RFC3339
	UnStandardRFC3339 := "2006-01-02T15:04Z07:00"
	create, err := time.Parse(time.RFC3339, createTimeStr)
	if err != nil {
		log.Printf("Parase the CreateTime %#v failed and error is: %#v.", createTime, err)
		create, err = time.Parse(UnStandardRFC3339, createTimeStr)
		if err != nil {
			return 0, errmsgs.WrapError(err)
		}
	}
	end, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		log.Printf("Parase the EndTime %#v failed and error is: %#v.", endTime, err)
		end, err = time.Parse(UnStandardRFC3339, endTimeStr)
		if err != nil {
			return 0, errmsgs.WrapError(err)
		}
	}
	var period int
	switch periodUnit {
	case "Month":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 30))
	case "Week":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 7))
	case "Year":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 365))
	default:
		err = fmt.Errorf("Unexpected period unit %s", periodUnit)
	}
	// The period at least is 1
	if period < 1 {
		period = 1
	}
	if period > 12 {
		period = 12
	}
	// period can not be modified and if the new period is changed, using the previous one.
	if currentPeriod > 0 && currentPeriod != period {
		period = currentPeriod
	}
	return period, errmsgs.WrapError(err)
}

func terraformToAPI(field string) string {
	var result string
	for _, v := range strings.Split(field, "_") {
		if len(v) > 0 {
			result = fmt.Sprintf("%s%s%s", result, strings.ToUpper(string(v[0])), v[1:])
		}
	}
	return result
}
func convertMaptoJsonString(m map[string]interface{}) (string, error) {
	sm := make(map[string]string, len(m))
	for k, v := range m {
		sm[k] = v.(string)
	}

	if result, err := json.Marshal(sm); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}
func formatInt(src interface{}) int {
	if src == nil {
		return 0
	}
	attrType := reflect.TypeOf(src)
	switch attrType.String() {
	case "float64":
		return int(src.(float64))
	case "float32":
		return int(src.(float32))
	case "int64":
		return int(src.(int64))
	case "int32":
		return int(src.(int32))
	case "int":
		return src.(int)
	case "string":
		v, err := strconv.Atoi(src.(string))
		if err != nil {
			panic(err)
		}
		return v
	case "json.Number":
		v, err := strconv.Atoi(src.(json.Number).String())
		if err != nil {
			panic(err)
		}
		return v
	default:
		panic(fmt.Sprintf("Not support type %s", attrType.String()))
	}
	return 0
}
func convertArrayObjectToJsonString(src interface{}) (string, error) {
	res, err := json.Marshal(&src)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
func convertListToCommaSeparate(configured []interface{}) string {
	if len(configured) < 1 {
		return ""
	}
	result := ""
	for i, v := range configured {
		rail := ","
		if i == len(configured)-1 {
			rail = ""
		}
		result += v.(string) + rail
	}
	return result
}
func compareJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	obj1 := make(map[string]interface{})
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	obj2 := make(map[string]interface{})
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Compare(canonicalJson1, canonicalJson2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}
	return equal, nil
}
func convertArrayToString(src interface{}, sep string) string {
	if src == nil {
		return ""
	}
	items := make([]string, 0)
	for _, v := range src.([]interface{}) {
		items = append(items, fmt.Sprint(v))
	}
	return strings.Join(items, sep)
}
func convertMapFloat64ToJsonString(m map[string]interface{}) (string, error) {
	sm := make(map[string]json.Number, len(m))

	for k, v := range m {
		sm[k] = v.(json.Number)
	}

	if result, err := json.Marshal(sm); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

// 合并两个 map，并在遇到相同键时覆盖第一个 map 的值
func mergeMaps(map1, map2 map[string]string) {

	// 将第二个 map 的所有键值对复制到 mergedMap，覆盖已存在的键
	for key, value := range map2 {
		map1[key] = value
	}
}

func mapMerge(target, merged map[string]interface{}) map[string]interface{} {
	for key, value := range merged {
		if _, exist := target[key]; !exist {
			target[key] = value
		} else {
			// key existed in both src,target
			switch merged[key].(type) {
			case []interface{}:
				sourceSlice := value.([]interface{})
				targetSlice := make([]interface{}, len(sourceSlice))
				copy(targetSlice, target[key].([]interface{}))

				for index, val := range sourceSlice {
					switch val.(type) {
					case map[string]interface{}:
						targetMap, ok := targetSlice[index].(map[string]interface{})
						if ok {
							targetSlice[index] = mapMerge(targetMap, val.(map[string]interface{}))
						} else {
							targetSlice[index] = mapMerge(map[string]interface{}{}, val.(map[string]interface{}))
						}
					default:
						targetSlice[index] = val
					}
				}
				target[key] = targetSlice
			case map[string]interface{}:
				target[key] = mapMerge(target[key].(map[string]interface{}), merged[key].(map[string]interface{}))
			default:
				target[key] = merged[key]
			}
		}
	}
	return target
}

func mapSort(target map[string]string) []string {
	result := make([]string, 0)
	for key := range target {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func newInstanceDiff(resourceName string, attributes, attributesDiff map[string]interface{}, state *terraform.InstanceState) (*terraform.InstanceDiff, error) {

	p := Provider().ResourcesMap
	dOld, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	dNew, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	for key, value := range attributes {
		err := dOld.Set(key, value)
		if err != nil {
			return nil, errmsgs.WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}
	for key, value := range attributesDiff {
		attributes[key] = value
	}

	for key, value := range attributes {
		err := dNew.Set(key, value)
		if err != nil {
			return nil, errmsgs.WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}

	diff := terraform.NewInstanceDiff()
	objectKey := ""
	for _, key := range mapSort(dNew.State().Attributes) {
		newValue := dNew.State().Attributes[key]
		if objectKey != "" && !strings.HasPrefix(key, objectKey) {
			objectKey = ""
		}
		if objectKey == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					objectKey = strings.TrimSuffix(key, suffix)
					break
				}
			}
		}
		oldValue, ok := dOld.State().Attributes[key]
		if ok && oldValue == newValue {
			continue
		}
		if oldValue == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					oldValue = "0"
				}
			}
		}
		// 使用 SetNew 和 SetOld 方法来设置属性差异
		if diff.Attributes == nil {
			diff.Attributes = make(map[string]*terraform.ResourceAttrDiff)
		}
		diff.Attributes[key] = &terraform.ResourceAttrDiff{
			Old: oldValue,
			New: newValue,
		}

		if objectKey != "" {
			for removeKey, removeValue := range dOld.State().Attributes {
				if strings.HasPrefix(removeKey, objectKey) {
					if _, ok := dNew.State().Attributes[removeKey]; !ok {
						// If the attribue has complex elements, there should remove the key, not setting it to empty
						if len(strings.Split(removeKey, ".")) > 2 {
							delete(diff.Attributes, removeKey)
						} else {
							diff.Attributes[removeKey] = &terraform.ResourceAttrDiff{
								Old: removeValue,
								New: "",
							}
						}
					}
				}
			}
			objectKey = ""
		}
	}
	return diff, nil
}

func setResourceFunc(resource *schema.Resource, createFunc schema.CreateFunc, readFunc schema.ReadFunc, updateFunc schema.UpdateFunc, deleteFunc schema.DeleteFunc) {
	resource.CreateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		var err error
		err = createFunc(d, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		waitSecondsIfWithTest(1)

		if updateFunc != nil {
			err = updateFunc(d, meta)
		}
		
		if err != nil {
			waitSecondsIfWithTest(3)
			// 如果创建成功但读取加载失败，tf不会终态，为方式残留资源，触发删除
			resource.DeleteContext(ctx, d, meta)
			return diag.FromErr(err)
		}
		
		waitSecondsIfWithTest(1)
		retry := 5
		for retry > 0 {
			// 大批量触发时asapi侧的资源同步会有一定的延迟，如果失败则重试
			err = readFunc(d, meta)
			if err != nil {
				time.Sleep(time.Second * 5)
				retry--
				continue
			}
			break
		}
		if err != nil {
			// 如果创建成功但读取加载失败，tf不会终态，为方式残留资源，触发删除
			waitSecondsIfWithTest(3)
			resource.DeleteContext(ctx, d, meta)
			return diag.FromErr(err)
		}
		
		return nil
	}

	resource.ReadContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		err := readFunc(d, meta)
		return diag.FromErr(err)
	}

	if updateFunc != nil {
		resource.UpdateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			update_err := updateFunc(d, meta)
			waitSecondsIfWithTest(1)
			read_err := readFunc(d, meta)
			if update_err != nil {
				return diag.FromErr(update_err)
			}
			return diag.FromErr(read_err)
		}
	}

	resource.DeleteContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		err := deleteFunc(d, meta)
		return diag.FromErr(err)
	}

	if resource.Importer == nil {
		resource.Importer = &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		}
	}
}

func noUpdatesAllowedCheck(d *schema.ResourceData, fields []string) error {
	if d.IsNewResource() {
		return nil
	}
	updatefields := make([]string, 0)
	for _, field := range fields {
		if d.HasChange(field) {
			updatefields = append(updatefields, field)
		}
	}
	if len(updatefields) > 0 {
		updatefieldsstr := fmt.Sprintf("%s", strings.Join(updatefields, ","))
		return errmsgs.Error(errmsgs.UpdateFailedErrorMsg, d.Id(), updatefieldsstr, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
