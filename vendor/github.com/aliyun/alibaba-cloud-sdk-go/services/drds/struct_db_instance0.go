package drds

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

// DbInstance is a nested struct in drds response
type DbInstance struct {
	ReadWeight        int                                       `json:"ReadWeight" xml:"ReadWeight"`
	Port              int                                       `json:"Port" xml:"Port"`
	DbInstType        string                                    `json:"DbInstType" xml:"DbInstType"`
	Engine            string                                    `json:"Engine" xml:"Engine"`
	EngineVersion     string                                    `json:"EngineVersion" xml:"EngineVersion"`
	RemainDays        string                                    `json:"RemainDays" xml:"RemainDays"`
	PayType           string                                    `json:"PayType" xml:"PayType"`
	ReadMode          string                                    `json:"ReadMode" xml:"ReadMode"`
	RdsInstType       string                                    `json:"RdsInstType" xml:"RdsInstType"`
	DBInstanceStatus  string                                    `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
	ExpireTime        string                                    `json:"ExpireTime" xml:"ExpireTime"`
	ConnectUrl        string                                    `json:"ConnectUrl" xml:"ConnectUrl"`
	DBInstanceId      string                                    `json:"DBInstanceId" xml:"DBInstanceId"`
	DmInstanceId      string                                    `json:"DmInstanceId" xml:"DmInstanceId"`
	NetworkType       string                                    `json:"NetworkType" xml:"NetworkType"`
	Endpoints         Endpoints                                 `json:"Endpoints" xml:"Endpoints"`
	DBNodes           DBNodes                                   `json:"DBNodes" xml:"DBNodes"`
	ReadOnlyInstances ReadOnlyInstancesInDescribeDrdsDbInstance `json:"ReadOnlyInstances" xml:"ReadOnlyInstances"`
}
