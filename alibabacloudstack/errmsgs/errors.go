package errmsgs

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	sls "github.com/aliyun/aliyun-log-go-sdk"

	"fmt"

	"log"
	"runtime"

	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"

	//"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// common
	NotFound                = "NotFound"
	ResourceNotfound        = "ResourceNotfound"
	InstanceNotFound        = "Instance.Notfound"
	VSwitchIdNotFound       = "VSwitchId.Notfound"
	MessageInstanceNotFound = "instance is not found"
	Throttling              = "Throttling"
	ServiceUnavailable      = "ServiceUnavailable"

	// RAM Instance Not Found
	RamInstanceNotFound              = "Forbidden.InstanceNotFound"
	AlibabacloudStackGoClientFailure = "AlibabacloudStackGoClientFailure"
	DenverdinoAlibabacloudStackgo    = ErrorSource("[SDK denverdino/aliyungo ERROR]")
	ThrottlingUser                   = "Throttling.User"
	LogClientTimeout                 = "Client.Timeout exceeded while awaiting headers"
	AlibabacloudstackMaxComputeSdkGo = ErrorSource("[SDK aliyun-maxcompute-sdk-go ERROR]")
	InvalidFileSystemStatus_Ordering = "InvalidFileSystemStatus.Ordering"
)

var SlbIsBusy = []string{"SystemBusy", "OperationBusy", "ServiceIsStopping", "BackendServer.configuring", "ServiceIsConfiguring"}
var EcsNotFound = []string{"InvalidInstanceId.NotFound", "Forbidden.InstanceNotFound"}
var DiskInvalidOperation = []string{"IncorrectDiskStatus", "IncorrectInstanceStatus", "OperationConflict", "InternalError", "InvalidOperation.Conflict", "IncorrectDiskStatus.Initializing"}
var NetworkInterfaceInvalidOperations = []string{"InvalidOperation.InvalidEniState", "InvalidOperation.InvalidEcsState", "OperationConflict", "ServiceUnavailable", "InternalError"}
var SnapshotInvalidOperations = []string{"OperationConflict", "ServiceUnavailable", "InternalError", "SnapshotCreatedDisk", "SnapshotCreatedImage"}
var SnapshotPolicyInvalidOperations = []string{"OperationConflict", "ServiceUnavailable", "InternalError", "SnapshotCreatedDisk", "SnapshotCreatedImage"}
var DiskNotSupportOnlineChangeErrors = []string{"InvalidDiskCategory.NotSupported", "InvalidRegion.NotSupport", "IncorrectInstanceStatus", "IncorrectDiskStatus", "InvalidOperation.InstanceTypeNotSupport"}
var DBReadInstanceNotReadyStatus = []string{"OperationDenied.ReadDBInstanceStatus", "OperationDenied.MasterDBInstanceState", "ReadDBInstance.Mismatch"}

// An Error represents a custom error for Terraform failure response
type ProviderError struct {
	errorCode string
	message   string
}

// details at: https://help.aliyun.com/document_detail/27300.html
var OtsTableIsTemporarilyUnavailable = []string{"no such host", "OTSServerBusy", "OTSPartitionUnavailable", "OTSInternalServerError",
	"OTSTimeout", "OTSServerUnavailable", "OTSRowOperationConflict", "OTSTableNotReady", "OTSNotEnoughCapacityUnit", "Too frequent table operations."}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[ERROR] Terraform AlibabacloudStack Provider Error: Code: %s Message: %s", e.errorCode, e.message)
}

func (err *ProviderError) ErrorCode() string {
	return err.errorCode
}

func (err *ProviderError) Message() string {
	return err.message
}

func GetNotFoundErrorFromString(str string) error {
	return &ProviderError{
		errorCode: InstanceNotFound,
		message:   str,
	}
}
func GetNotFoundVPCError(str string) error {
	return &ProviderError{
		errorCode: VSwitchIdNotFound,
		message:   str,
	}
}
func NotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*ComplexError); ok {
		if e.Err != nil && strings.HasPrefix(e.Err.Error(), ResourceNotfound) {
			return true
		}
		return NotFoundError(e.Cause)
	}
	if err == nil {
		return false
	}

	if e, ok := err.(*sdkerrors.ServerError); ok {
		return e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)
	}

	if e, ok := err.(*ProviderError); ok {
		return e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)
	}

	if e, ok := err.(*common.Error); ok {
		return e.Code == InstanceNotFound || e.Code == RamInstanceNotFound || e.Code == NotFound || strings.Contains(strings.ToLower(e.Message), MessageInstanceNotFound)
	}

	if e, ok := err.(oss.ServiceError); ok {
		return e.StatusCode == 404 || strings.HasPrefix(e.Code, "NoSuch") || strings.HasPrefix(e.Message, "No Row found") || strings.HasPrefix(e.Message, "ResourceNotfound")
	}

	if e, ok := err.(*tea.SDKError); ok {
		return *e.StatusCode == 404 || strings.HasSuffix(*e.Code, ".NotFound")
	}

	return strings.HasSuffix(err.Error(), ".NotFound") || strings.HasPrefix(err.Error(), ResourceNotfound)

}

func NeedRetry(err error) bool {
	if err == nil {
		return false
	}

	postRegex := regexp.MustCompile("^Post [\"]*https://.*")
	if postRegex.MatchString(err.Error()) {
		return true
	}

	throttlingRegex := regexp.MustCompile("^Throttling.*")
	codeRegex := regexp.MustCompile("^code: 5[\\d]{2}")

	if e, ok := err.(*tea.SDKError); ok {
		if strings.Contains(*e.Message, "code: 500, 您已开通过") {
			return false
		}
		if strings.Contains(*e.Message, "The current status of the resource does not support this operation, please retry again.") {
			return true
		}
		if *e.Code == ServiceUnavailable || *e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(*e.Code) || codeRegex.MatchString(*e.Message) {
			return true
		}
	}

	if e, ok := err.(*sdkerrors.ServerError); ok {
		return e.ErrorCode() == ServiceUnavailable || e.ErrorCode() == "Rejected.Throttling" || throttlingRegex.MatchString(e.ErrorCode()) || codeRegex.MatchString(e.Message())
	}

	if e, ok := err.(*common.Error); ok {
		return e.Code == ServiceUnavailable || e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(e.Code) || codeRegex.MatchString(e.Message)
	}

	return false
}

func IsExpectedErrorCodes(code string, errorCodes []string) bool {
	if code == "" {
		return false
	}
	for _, v := range errorCodes {
		if v == code {
			return true
		}
	}
	return false
}

func IsExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ComplexError); ok {
		return IsExpectedErrors(e.Cause, expectCodes)
	}
	if err == nil {
		return false
	}

	if e, ok := err.(*sdkerrors.ServerError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*ProviderError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*common.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*sls.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) || strings.Contains(e.String(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(oss.ServiceError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*fc.ServiceError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode == code || strings.Contains(e.ErrorMessage, code) {
				return true
			}
		}
		return false
	}

	/*if e, ok := err.(datahub.DatahubError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}*/

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

func IsThrottling(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*sdkerrors.ServerError); ok {
		if e.ErrorCode() == Throttling {
			return true
		}
		return false
	}

	if e, ok := err.(*common.Error); ok {
		if e.Code == Throttling {
			return true
		}
		return false
	}
	return false
}

func GetTimeErrorFromString(str string) error {
	return &ProviderError{
		errorCode: "WaitForTimeout",
		message:   str,
	}
}

func GetNotFoundMessage(product string, id string) string {
	return fmt.Sprintf("The specified %s %s is not found.", product, id)
}
func GetNotVPCMessage() string {
	return fmt.Sprintf("The VSwitchId is not found.")
}
func GetTimeoutMessage(product string, status string) string {
	return fmt.Sprintf("Waitting for %s %s is timeout.", product, status)
}

type ErrorSource string

const (
	AlibabacloudStackSdkGoERROR    = ErrorSource("[SDK alibaba-cloud-sdk-go ERROR]")
	ProviderERROR                  = ErrorSource("[Provider ERROR]")
	AlibabacloudStackOssGoSdk      = ErrorSource("[SDK aliyun-oss-go-sdk ERROR]")
	AlibabacloudStackLogGoSdkERROR = ErrorSource("[SDK aliyun-log-go-sdk ERROR]")
	AliyunTablestoreGoSdk          = ErrorSource("[SDK aliyun-tablestore-go-sdk ERROR]")
	AlibabacloudStackDatahubSdkGo  = ErrorSource("[SDK aliyun-datahub-sdk-go ERROR]")
	DenverdinoAliyungo             = ErrorSource("[SDK denverdino/aliyungo ERROR]")
)

// ComplexError is a format error which including origin error, extra error message, error occurred file and line
// Cause: a error is a origin error that comes from SDK, some exceptions and so on
// Err: a new error is built from extra message
// Path: the file path of error occurred
// Line: the file line of error occurred
type ComplexError struct {
	Cause error
	Err   error
	Path  string
	Line  int
}

func (e ComplexError) Error() string {
	if e.Cause == nil {
		e.Cause = Error("<nil cause>")
	}
	if e.Err == nil {
		return fmt.Sprintf("[ERROR] %s:%d:\n%s", e.Path, e.Line, e.Cause.Error())
	}
	return fmt.Sprintf("[ERROR] %s:%d: %s:\n%s", e.Path, e.Line, e.Err.Error(), e.Cause.Error())
}

func Error(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

// Return a ComplexError which including error occurred file and path
func WrapError(cause error) error {
	if cause == nil {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in WrapError.")
		return WrapComplexError(cause, nil, "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	return WrapComplexError(cause, nil, filepath, line)
}

// Return a ComplexError which including extra error message, error occurred file and path
func WrapErrorf(cause error, msg string, args ...interface{}) error {
	if cause == nil && strings.TrimSpace(msg) == "" {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in WrapErrorf.")
		return WrapComplexError(cause, Error(msg), "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	// The second parameter of args is requestId, if the error message is NotFoundMsg the requestId need to be returned.
	if msg == NotFoundMsg && len(args) == 2 {
		msg += RequestIdMsg
	}
	return WrapComplexError(cause, fmt.Errorf(msg, args...), filepath, line)
}

func WrapComplexError(cause, err error, filepath string, fileline int) error {
	return &ComplexError{
		Cause: cause,
		Err:   err,
		Path:  filepath,
		Line:  fileline,
	}
}

func GetBaseResponseErrorMessage(err_response *responses.BaseResponse) (showMsg string) {
	raw_data := make(map[string]interface{})
	err := json.Unmarshal(err_response.GetHttpContentBytes(), &raw_data)
	if err != nil {
		return
	}
	return GetAsapiErrorMessage(raw_data)
}

func GetAsapiErrorMessage(raw_data map[string]interface{}) (showMsg string) {
	showfields := []string{"errorTitle", "errorMessage", "errorCode", "eagleEyeTraceId", "RequestId"}
	for _, field := range showfields {
		value := raw_data[field]
		if value != nil {
			showMsg = showMsg + fmt.Sprintf("\n%s: %s", field, value)
		}
	}
	return showMsg
}

func CheckEmpty(value interface{}, schemaType schema.ValueType, keys ...string) error {
	zero := schemaType.Zero()

	empty := false
	if eq, ok := value.(schema.Equal); ok {
		empty = eq.Equal(zero)
	} else {
		empty = reflect.DeepEqual(value, zero)
	}

	if !empty {
		return nil
	}
	errmsg := strings.Join(keys, " or ")
	return errors.New(errmsg + " can not be empty at the same time")
}

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
const IdMsg = "Resource id：%s "
const DefaultErrorMsg = "Resource %s %s Failed!!! %s"
const RequestV1ErrorMsg = "Resource %s %s Failed!!! %s%s"
const VPCErrorMsg = "Resource %s %s Failed!!! %s"
const RequestIdMsg = "RequestId: %s"
const NotFoundMsg = ResourceNotfound + "!!! %s"
const WaitTimeoutMsg = "Resource %s %s Timeout In %d Seconds. Got: %s Expected: %s !!! %s"
const DataDefaultErrorMsg = "Datasource %s %s Failed!!! %s"

var OperationDeniedDBStatus = []string{"OperationDenied.DBStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBClusterStatus", "InternalError", "OperationDenied.OutofUsage"}

const DefaultTimeoutMsg = "Resource %s %s Timeout!!! %s"
const DefaultDebugMsg = "\n*************** %s Response *************** \n%s\n%s******************************\n\n"
const FailedToReachTargetStatus = "Failed to reach target status. Current status is %s."
const FailedGetAttributeMsg = "Getting resource %s attribute by path %s failed!!! Body: %v."
const NotFoundWithResponse = ResourceNotfound + "!!! Response: %v"
