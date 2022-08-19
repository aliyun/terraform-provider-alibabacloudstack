package adb

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateDBCluster invokes the adb.CreateDBCluster API synchronously
// api document: https://help.aliyun.com/api/adb/createdbcluster.html
func (client *Client) CreateDBCluster(request *CreateDBClusterRequest) (response *CreateDBClusterResponse, err error) {
	response = CreateCreateDBClusterResponse()
	err = client.DoAction(request, response)
	return
}

// CreateDBClusterWithChan invokes the adb.CreateDBCluster API asynchronously
// api document: https://help.aliyun.com/api/adb/createdbcluster.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateDBClusterWithChan(request *CreateDBClusterRequest) (<-chan *CreateDBClusterResponse, <-chan error) {
	responseChan := make(chan *CreateDBClusterResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateDBCluster(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// CreateDBClusterWithCallback invokes the adb.CreateDBCluster API asynchronously
// api document: https://help.aliyun.com/api/adb/createdbcluster.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateDBClusterWithCallback(request *CreateDBClusterRequest, callback func(response *CreateDBClusterResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateDBClusterResponse
		var err error
		defer close(result)
		response, err = client.CreateDBCluster(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// CreateDBClusterRequest is the request struct for api CreateDBCluster
type CreateDBClusterRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DBClusterDescription string           `position:"Query" name:"DBClusterDescription"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	StorageType          string           `position:"Query" name:"StorageType"`
	Mode                 string           `position:"Query" name:"Mode"`
	ResourceGroupId      string           `position:"Query" name:"ResourceGroupId"`
	StorageResource      string           `position:"Query" name:"StorageResource"`
	DBClusterCategory    string           `position:"Query" name:"DBClusterCategory"`
	DBClusterNetworkType string           `position:"Query" name:"DBClusterNetworkType"`
	Period               string           `position:"Query" name:"Period"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	DBClusterVersion     string           `position:"Query" name:"DBClusterVersion"`
	DBClusterClass       string           `position:"Query" name:"DBClusterClass"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	DBNodeGroupCount     string           `position:"Query" name:"DBNodeGroupCount"`
	UsedTime             string           `position:"Query" name:"UsedTime"`
	VSwitchId            string           `position:"Query" name:"VSwitchId"`
	DBNodeStorage        string           `position:"Query" name:"DBNodeStorage"`
	ExecutorCount        string           `position:"Query" name:"ExecutorCount"`
	VPCId                string           `position:"Query" name:"VPCId"`
	ZoneId               string           `position:"Query" name:"ZoneId"`
	ComputeResource      string           `position:"Query" name:"ComputeResource"`
	PayType              string           `position:"Query" name:"PayType"`
	ClusterType          string           `position:"Query" name:"ClusterType"`
	CpuType              string           `position:"Query" name:"CpuType"`
}

// CreateDBClusterResponse is the response struct for api CreateDBCluster
type CreateDBClusterResponse struct {
	*responses.BaseResponse
	RequestId       string `json:"RequestId" xml:"RequestId"`
	DBClusterId     string `json:"DBClusterId" xml:"DBClusterId"`
	OrderId         string `json:"OrderId" xml:"OrderId"`
	ResourceGroupId string `json:"ResourceGroupId" xml:"ResourceGroupId"`
}

// CreateCreateDBClusterRequest creates a request to invoke CreateDBCluster API
func CreateCreateDBClusterRequest() (request *CreateDBClusterRequest) {
	request = &CreateDBClusterRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("adb", "2019-03-15", "CreateDBCluster", "ads", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateDBClusterResponse creates a response to parse from CreateDBCluster response
func CreateCreateDBClusterResponse() (response *CreateDBClusterResponse) {
	response = &CreateDBClusterResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
