package vpc

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

// ModifyBandwidthPackageAttribute invokes the vpc.ModifyBandwidthPackageAttribute API synchronously
// api document: https://help.aliyun.com/api/vpc/modifybandwidthpackageattribute.html
func (client *Client) ModifyBandwidthPackageAttribute(request *ModifyBandwidthPackageAttributeRequest) (response *ModifyBandwidthPackageAttributeResponse, err error) {
	response = CreateModifyBandwidthPackageAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyBandwidthPackageAttributeWithChan invokes the vpc.ModifyBandwidthPackageAttribute API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifybandwidthpackageattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyBandwidthPackageAttributeWithChan(request *ModifyBandwidthPackageAttributeRequest) (<-chan *ModifyBandwidthPackageAttributeResponse, <-chan error) {
	responseChan := make(chan *ModifyBandwidthPackageAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyBandwidthPackageAttribute(request)
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

// ModifyBandwidthPackageAttributeWithCallback invokes the vpc.ModifyBandwidthPackageAttribute API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifybandwidthpackageattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyBandwidthPackageAttributeWithCallback(request *ModifyBandwidthPackageAttributeRequest, callback func(response *ModifyBandwidthPackageAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyBandwidthPackageAttributeResponse
		var err error
		defer close(result)
		response, err = client.ModifyBandwidthPackageAttribute(request)
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

// ModifyBandwidthPackageAttributeRequest is the request struct for api ModifyBandwidthPackageAttribute
type ModifyBandwidthPackageAttributeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Description          string           `position:"Query" name:"Description"`
	BandwidthPackageId   string           `position:"Query" name:"BandwidthPackageId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Name                 string           `position:"Query" name:"Name"`
}

// ModifyBandwidthPackageAttributeResponse is the response struct for api ModifyBandwidthPackageAttribute
type ModifyBandwidthPackageAttributeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyBandwidthPackageAttributeRequest creates a request to invoke ModifyBandwidthPackageAttribute API
func CreateModifyBandwidthPackageAttributeRequest() (request *ModifyBandwidthPackageAttributeRequest) {
	request = &ModifyBandwidthPackageAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "ModifyBandwidthPackageAttribute", "vpc", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyBandwidthPackageAttributeResponse creates a response to parse from ModifyBandwidthPackageAttribute response
func CreateModifyBandwidthPackageAttributeResponse() (response *ModifyBandwidthPackageAttributeResponse) {
	response = &ModifyBandwidthPackageAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
