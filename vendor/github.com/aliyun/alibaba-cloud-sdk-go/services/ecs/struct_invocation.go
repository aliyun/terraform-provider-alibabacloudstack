package ecs

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

// Invocation is a nested struct in ecs response
type Invocation struct {
	Name              string                               `json:"Name" xml:"Name"`
	PageSize          int64                                `json:"PageSize" xml:"PageSize"`
	Timed             bool                                 `json:"Timed" xml:"Timed"`
	Frequency         string                               `json:"Frequency" xml:"Frequency"`
	Content           string                               `json:"Content" xml:"Content"`
	CommandContent    string                               `json:"CommandContent" xml:"CommandContent"`
	InvocationStatus  string                               `json:"InvocationStatus" xml:"InvocationStatus"`
	FileGroup         string                               `json:"FileGroup" xml:"FileGroup"`
	Description       string                               `json:"Description" xml:"Description"`
	Overwrite         string                               `json:"Overwrite" xml:"Overwrite"`
	PageNumber        int64                                `json:"PageNumber" xml:"PageNumber"`
	CommandId         string                               `json:"CommandId" xml:"CommandId"`
	TargetDir         string                               `json:"TargetDir" xml:"TargetDir"`
	FileMode          string                               `json:"FileMode" xml:"FileMode"`
	TotalCount        int64                                `json:"TotalCount" xml:"TotalCount"`
	ContentType       string                               `json:"ContentType" xml:"ContentType"`
	CreationTime      string                               `json:"CreationTime" xml:"CreationTime"`
	CommandName       string                               `json:"CommandName" xml:"CommandName"`
	Parameters        string                               `json:"Parameters" xml:"Parameters"`
	VmCount           int                                  `json:"VmCount" xml:"VmCount"`
	InvokeId          string                               `json:"InvokeId" xml:"InvokeId"`
	InvokeStatus      string                               `json:"InvokeStatus" xml:"InvokeStatus"`
	FileOwner         string                               `json:"FileOwner" xml:"FileOwner"`
	CommandType       string                               `json:"CommandType" xml:"CommandType"`
	InvokeInstances   InvokeInstancesInDescribeInvocations `json:"InvokeInstances" xml:"InvokeInstances"`
	InvocationResults InvocationResults                    `json:"InvocationResults" xml:"InvocationResults"`
}
