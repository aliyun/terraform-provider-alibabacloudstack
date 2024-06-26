package edas

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

// BindK8sSlb invokes the edas.BindK8sSlb API synchronously
func (client *Client) BindK8sSlb(request *BindK8sSlbRequest) (response *BindK8sSlbResponse, err error) {
	response = CreateBindK8sSlbResponse()
	err = client.DoAction(request, response)
	return
}

// BindK8sSlbWithChan invokes the edas.BindK8sSlb API asynchronously
func (client *Client) BindK8sSlbWithChan(request *BindK8sSlbRequest) (<-chan *BindK8sSlbResponse, <-chan error) {
	responseChan := make(chan *BindK8sSlbResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.BindK8sSlb(request)
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

// BindK8sSlbWithCallback invokes the edas.BindK8sSlb API asynchronously
func (client *Client) BindK8sSlbWithCallback(request *BindK8sSlbRequest, callback func(response *BindK8sSlbResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *BindK8sSlbResponse
		var err error
		defer close(result)
		response, err = client.BindK8sSlb(request)
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

// BindK8sSlbRequest is the request struct for api BindK8sSlb
type BindK8sSlbRequest struct {
	*requests.RoaRequest
	Scheduler        string `position:"Query" name:"Scheduler"`
	ServicePortInfos string `position:"Query" name:"ServicePortInfos"`
	SlbId            string `position:"Query" name:"SlbId"`
	SlbProtocol      string `position:"Query" name:"SlbProtocol"`
	Port             string `position:"Query" name:"Port"`
	AppId            string `position:"Query" name:"AppId"`
	Specification    string `position:"Query" name:"Specification"`
	ClusterId        string `position:"Query" name:"ClusterId"`
	Type             string `position:"Query" name:"Type"`
	TargetPort       string `position:"Query" name:"TargetPort"`
}

// BindK8sSlbResponse is the response struct for api BindK8sSlb
type BindK8sSlbResponse struct {
	*responses.BaseResponse
	ChangeOrderId string `json:"ChangeOrderId" xml:"ChangeOrderId"`
	Code          int    `json:"Code" xml:"Code"`
	Message       string `json:"Message" xml:"Message"`
	RequestId     string `json:"RequestId" xml:"RequestId"`
}

// CreateBindK8sSlbRequest creates a request to invoke BindK8sSlb API
func CreateBindK8sSlbRequest() (request *BindK8sSlbRequest) {
	request = &BindK8sSlbRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "BindK8sSlb", "/roa/pop/v5/k8s/acs/k8s_slb_binding", "Edas", "openAPI")
	request.Method = requests.POST
	return
}

// CreateBindK8sSlbResponse creates a response to parse from BindK8sSlb response
func CreateBindK8sSlbResponse() (response *BindK8sSlbResponse) {
	response = &BindK8sSlbResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
