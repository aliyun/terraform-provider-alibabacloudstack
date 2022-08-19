package cdn

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

// DescribeDomainCcActivityLog invokes the cdn.DescribeDomainCcActivityLog API synchronously
func (client *Client) DescribeDomainCcActivityLog(request *DescribeDomainCcActivityLogRequest) (response *DescribeDomainCcActivityLogResponse, err error) {
	response = CreateDescribeDomainCcActivityLogResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainCcActivityLogWithChan invokes the cdn.DescribeDomainCcActivityLog API asynchronously
func (client *Client) DescribeDomainCcActivityLogWithChan(request *DescribeDomainCcActivityLogRequest) (<-chan *DescribeDomainCcActivityLogResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainCcActivityLogResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainCcActivityLog(request)
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

// DescribeDomainCcActivityLogWithCallback invokes the cdn.DescribeDomainCcActivityLog API asynchronously
func (client *Client) DescribeDomainCcActivityLogWithCallback(request *DescribeDomainCcActivityLogRequest, callback func(response *DescribeDomainCcActivityLogResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainCcActivityLogResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainCcActivityLog(request)
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

// DescribeDomainCcActivityLogRequest is the request struct for api DescribeDomainCcActivityLog
type DescribeDomainCcActivityLogRequest struct {
	*requests.RpcRequest
	RuleName      string           `position:"Query" name:"RuleName"`
	StartTime     string           `position:"Query" name:"StartTime"`
	TriggerObject string           `position:"Query" name:"TriggerObject"`
	PageNumber    requests.Integer `position:"Query" name:"PageNumber"`
	PageSize      requests.Integer `position:"Query" name:"PageSize"`
	Value         string           `position:"Query" name:"Value"`
	DomainName    string           `position:"Query" name:"DomainName"`
	EndTime       string           `position:"Query" name:"EndTime"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeDomainCcActivityLogResponse is the response struct for api DescribeDomainCcActivityLog
type DescribeDomainCcActivityLogResponse struct {
	*responses.BaseResponse
	RequestId   string    `json:"RequestId" xml:"RequestId"`
	PageIndex   int64     `json:"PageIndex" xml:"PageIndex"`
	PageSize    int64     `json:"PageSize" xml:"PageSize"`
	Total       int64     `json:"Total" xml:"Total"`
	ActivityLog []LogInfo `json:"ActivityLog" xml:"ActivityLog"`
}

// CreateDescribeDomainCcActivityLogRequest creates a request to invoke DescribeDomainCcActivityLog API
func CreateDescribeDomainCcActivityLogRequest() (request *DescribeDomainCcActivityLogRequest) {
	request = &DescribeDomainCcActivityLogRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "DescribeDomainCcActivityLog", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDomainCcActivityLogResponse creates a response to parse from DescribeDomainCcActivityLog response
func CreateDescribeDomainCcActivityLogResponse() (response *DescribeDomainCcActivityLogResponse) {
	response = &DescribeDomainCcActivityLogResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
