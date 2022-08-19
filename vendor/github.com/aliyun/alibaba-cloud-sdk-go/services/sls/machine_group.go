package sls

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	_ "github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)
// const MachineIDTypes
const (
	MachineIDTypeIP          = "ip"
	MachineIDTypeUserDefined = "userdefined"
)

// MachinGroupAttribute defines machine group attribute
type MachinGroupAttribute struct {
	ExternalName string `json:"externalName"`
	TopicName    string `json:"groupTopic"`
}
type baseRequest struct {
	Scheme         string
	Method         string
	Domain         string
	Port           string
	RegionId       string
	isInsecure     *bool

	userAgent map[string]string
	product   string
	version   string

	actionName string

	AcceptFormat string

	QueryParams map[string]string
	Headers     map[string]string
	FormParams  map[string]string
	Content     []byte

	locationServiceCode  string
	locationEndpointType string

	queries string

	stringToSign string
}
// MachineGroup defines machine group
type MachineGroup struct {
	*requests.RpcRequest
	Name           string               `json:"groupName"`
	Type           string               `json:"groupType"`
	MachineIDType  string               `json:"machineIdentifyType"`
	MachineIDList  []string             `json:"machineList"`
	Attribute      MachinGroupAttribute `json:"groupAttribute"`
	CreateTime     uint32               `json:"createTime,omitempty"`
	LastModifyTime uint32               `json:"lastModifyTime,omitempty"`
}
// CreateDescribeDBInstancesRequest creates a request to invoke DescribeDBInstances API

// Machine defines machine struct
type Machine struct {
	IP            string
	UniqueID      string `json:"machine-uniqueid"`
	UserdefinedID string `json:"userdefined-id"`
	LastHeartBeatTime int `json:"lastHeartbeatTime"`
}

// MachineList defines machine list
type MachineList struct {
	Total    int
	Machines []*Machine
}
func CreateCreateSLSMachineGroupRequest() (request *MachineGroup) {
	request = &MachineGroup{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("SLS", "2020-03-31", "CreateMachineGroup", "sls", "openAPI")
	request.Method = requests.POST
	return
}