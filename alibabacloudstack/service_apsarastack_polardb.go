package alibabacloudstack

import (
	"slices"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolardbService struct {
	client *connectivity.AlibabacloudStackClient
}

type PolardbCheckaccountnameavailableResponse struct {
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DoPolardbCheckaccountnameavailableRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbCheckaccountnameavailableResponse, error) {
	// api: polardb - 2024-01-30 - CheckAccountNameAvailable
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "CheckAccountNameAvailable", "")
	PolardbCheckaccountnameavailableResponse := &PolardbCheckaccountnameavailableResponse{}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "CheckAccountNameAvailable", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbCheckaccountnameavailableResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "CheckAccountNameAvailable", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbCheckaccountnameavailableResponse, nil
}

type PolardbDescribeaccountsResponse struct {
	Accounts struct {
		DBInstanceAccount []struct {
			DatabasePrivileges struct {
				DatabasePrivilege []struct {
					DBName                 string `json:"DBName"`
					AccountPrivilege       string `json:"AccountPrivilege"`
					AccountPrivilegeDetail string `json:"AccountPrivilegeDetail"`
				} `json:"DatabasePrivilege"`
			} `json:"DatabasePrivileges"`
			DBInstanceId       string `json:"DBInstanceId"`
			AccountName        string `json:"AccountName"`
			AccountStatus      string `json:"AccountStatus"`
			AccountType        string `json:"AccountType"`
			AccountDescription string `json:"AccountDescription"`
			PrivExceeded       string `json:"PrivExceeded"`
			ValidUntil         string `json:"ValidUntil"`
			CreateDB           string `json:"CreateDB"`
			Replication        string `json:"Replication"`
			CreateRole         string `json:"CreateRole"`
			BypassRLS          string `json:"BypassRLS"`
		} `json:"DBInstanceAccount"`
	} `json:"Accounts"`
	RequestId                             string `json:"RequestId"`
	SystemAdminAccountStatus              string `json:"SystemAdminAccountStatus"`
	SystemAdminAccountFirstActivationTime string `json:"SystemAdminAccountFirstActivationTime"`
}

func (s *PolardbService) DescribeDBAccount(id string) (*PolardbDescribeaccountsResponse, error) {
	parts, _ := ParseResourceId(id, 2)
	// api: polardb - 2024-01-30 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeAccounts", "")
	PolardbDescribeaccountsResponse := &PolardbDescribeaccountsResponse{}

	request.QueryParams["AccountName"] = parts[1]

	// Regular parameter filling
	request.QueryParams["DBInstanceId"] = parts[0]

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
			return nil, nil
		}
		if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeAccounts", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribeaccountsResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeAccounts", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribeaccountsResponse.Accounts.DBInstanceAccount) < 1 {
		return PolardbDescribeaccountsResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("accounts", "")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribeaccountsResponse, nil
}

func (s *PolardbService) DoPolardbDescribeaccountsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribeaccountsResponse, error) {
	// api: polardb - 2024-01-30 - DescribeAccounts
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeAccounts", "")
	PolardbDescribeaccountsResponse := &PolardbDescribeaccountsResponse{}

	// Call request_params_handler

	// Regular parameter filling
	if v, ok := d.GetOk("account_name"); ok && v != "" {
		// Call requestin_handler
		request.QueryParams["AccountName"] = v.(string)
	}

	// Regular parameter filling
	if v, ok := d.GetOk("data_base_instance_id"); ok && v != "" {
		// Call requestin_handler
		request.QueryParams["DBInstanceId"] = v.(string)
	} else {
		return nil, fmt.Errorf("DataBaseInstanceId is required")
	}

	// Regular parameter filling
	if v, ok := d.GetOk("page_number"); ok {
		// Call requestin_handler
		request.QueryParams["PageNumber"] = strconv.Itoa(v.(int))
	}

	// Regular parameter filling
	if v, ok := d.GetOk("page_size"); ok {
		// Call requestin_handler
		request.QueryParams["PageSize"] = strconv.Itoa(v.(int))
	}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeAccounts", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribeaccountsResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeAccounts", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribeaccountsResponse.Accounts.DBInstanceAccount) < 1 {
		return PolardbDescribeaccountsResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("accounts", "")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return PolardbDescribeaccountsResponse, nil
}

type PolardbDescribedatabasesResponse struct {
	Databases struct {
		Database []struct {
			Accounts struct {
				AccountPrivilegeInfo []struct {
					Account          string `json:"Account"`
					AccountPrivilege string `json:"AccountPrivilege"`
				} `json:"AccountPrivilegeInfo"`
			} `json:"Accounts"`
			DBName           string `json:"DBName"`
			DBInstanceId     string `json:"DBInstanceId"`
			Engine           string `json:"Engine"`
			DBStatus         string `json:"DBStatus"`
			CharacterSetName string `json:"CharacterSetName"`
			DBDescription    string `json:"DBDescription"`
			Collate          string `json:"Collate"`
			Ctype            string `json:"Ctype"`
			ConnLimit        int    `json:"ConnLimit"`
			Tablespace       string `json:"Tablespace"`
		} `json:"Database"`
	} `json:"Databases"`
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DescribeDBDatabase(id string) (*PolardbDescribedatabasesResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDatabases
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDatabases", "")
	PolardbDescribedatabasesResponse := &PolardbDescribedatabasesResponse{}

	parts, err := ParseResourceId(id, 2)

	request.QueryParams["DBInstanceId"] = parts[0]

	request.QueryParams["DBName"] = parts[1]
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
			return nil, nil
		}
		if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDatabases", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedatabasesResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDatabases", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedatabasesResponse.Databases.Database) < 1 {
		return PolardbDescribedatabasesResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Databases", "")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedatabasesResponse, nil
}
func (s *PolardbService) DoPolardbDescribedatabasesRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribedatabasesResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDatabases
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDatabases", "")
	PolardbDescribedatabasesResponse := &PolardbDescribedatabasesResponse{}

	// Call request_params_handler

	// Regular parameter filling
	if v, ok := d.GetOk("data_base_instance_id"); ok && v != "" {
		// Call requestin_handler
		request.QueryParams["DBInstanceId"] = v.(string)
	} else {
		return nil, fmt.Errorf("DataBaseInstanceId is required")
	}

	// Regular parameter filling
	if v, ok := d.GetOk("data_base_name"); ok && v != "" {
		// Call requestin_handler
		request.QueryParams["DBName"] = v.(string)
	}

	// Regular parameter filling
	if v, ok := d.GetOk("page_number"); ok {
		// Call requestin_handler
		request.QueryParams["PageNumber"] = strconv.Itoa(v.(int))
	}

	// Regular parameter filling
	if v, ok := d.GetOk("page_size"); ok {
		// Call requestin_handler
		request.QueryParams["PageSize"] = strconv.Itoa(v.(int))
	}

	// Regular parameter filling
	if v, ok := d.GetOk("status"); ok && v != "" {
		// Call requestin_handler
		request.QueryParams["DBStatus"] = v.(string)
	}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDatabases", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedatabasesResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDatabases", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedatabasesResponse.Databases.Database) < 1 {
		return PolardbDescribedatabasesResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Databases", "")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedatabasesResponse, nil
}

type PolardbDescribebackuppolicyResponse struct {
	RequestId                     string `json:"RequestId"`
	BackupRetentionPeriod         int    `json:"BackupRetentionPeriod"`
	PreferredNextBackupTime       string `json:"PreferredNextBackupTime"`
	MaxRetention                  int    `json:"MaxRetention"`
	PreferredBackupTime           string `json:"PreferredBackupTime"`
	PreferredBackupPeriod         string `json:"PreferredBackupPeriod"`
	BackupLog                     string `json:"BackupLog"`
	LogBackupRetentionPeriod      string `json:"LogBackupRetentionPeriod"`
	EnableBackupLog               string `json:"EnableBackupLog"`
	LocalLogRetentionHours        string `json:"LocalLogRetentionHours"`
	LocalLogRetentionSpace        string `json:"LocalLogRetentionSpace"`
	CompressType                  string `json:"CompressType"`
	Duplication                   string `json:"Duplication"`
	DuplicationContent            string `json:"DuplicationContent"`
	HighSpaceUsageProtection      string `json:"HighSpaceUsageProtection"`
	LogBackupFrequency            string `json:"LogBackupFrequency"`
	ArchiveBackupRetentionPeriod  int    `json:"ArchiveBackupRetentionPeriod"`
	ArchiveBackupKeepPolicy       int    `json:"ArchiveBackupKeepPolicy"`
	ArchiveBackupKeepCount        int    `json:"ArchiveBackupKeepCount"`
	ReleasedKeepPolicy            string `json:"ReleasedKeepPolicy"`
	LogBackupLocalRetentionNumber int    `json:"LogBackupLocalRetentionNumber"`
	BackupMethod                  string `json:"BackupMethod"`

	DuplicationLocation struct {
		Sotrage string `json:"Sotrage"`

		Location struct {
			Endpoint string `json:"Endpoint"`
			Bucket   string `json:"Bucket"`
		} `json:"Location"`
	} `json:"DuplicationLocation"`
}

func (s *PolardbService) DoPolardbDescribebackuppolicyRequest(id string) (*PolardbDescribebackuppolicyResponse, error) {
	// api: polardb - 2024-01-30 - DescribeBackupPolicy
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeBackupPolicy", "")
	PolardbDescribebackuppolicyResponse := &PolardbDescribebackuppolicyResponse{}

	// Call requestin_handler
	request.QueryParams["DBInstanceId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBackupPolicy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribebackuppolicyResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBackupPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribebackuppolicyResponse, nil
}

type PolardbDescriberegionsResponse struct {
	Regions struct {
		RDSRegion []struct {
			RegionId  string `json:"RegionId"`
			ZoneId    string `json:"ZoneId"`
			SubDomain string `json:"SubDomain"`
			SubZoneId string `json:"SubZoneId"`
		} `json:"RDSRegion"`
	} `json:"Regions"`
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DoPolardbDescriberegionsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescriberegionsResponse, error) {
	// api: polardb - 2024-01-30 - DescribeRegions
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeRegions", "")
	PolardbDescriberegionsResponse := &PolardbDescriberegionsResponse{}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeRegions", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescriberegionsResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeRegions", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescriberegionsResponse, nil
}

type DBInstanceNetInfo struct {
	SecurityIPGroups struct {
		SecurityIPGroup []struct {
			SecurityIPGroupName string `json:"SecurityIPGroupName"`
			SecurityIPs         string `json:"SecurityIPs"`
		} `json:"securityIPGroup"`
	} `json:"SecurityIPGroups"`

	DBInstanceWeights struct {
		DBInstanceWeight []struct {
			DBInstanceId   string `json:"DBInstanceId"`
			DBInstanceType string `json:"DBInstanceType"`
			Availability   string `json:"Availability"`
			Weight         string `json:"Weight"`
		} `json:"DBInstanceWeight"`
	} `json:"DBInstanceWeights"`
	Upgradeable          string `json:"Upgradeable"`
	ExpiredTime          string `json:"ExpiredTime"`
	ConnectionString     string `json:"ConnectionString"`
	IPAddress            string `json:"IPAddress"`
	IPType               string `json:"IPType"`
	Port                 string `json:"Port"`
	VPCId                string `json:"VPCId"`
	VSwitchId            string `json:"VSwitchId"`
	ConnectionStringType string `json:"ConnectionStringType"`
	MaxDelayTime         string `json:"MaxDelayTime"`
	DistributionType     string `json:"DistributionType"`
}

type PolardbDescribedbinstancenetinfoResponse struct {
	DBInstanceNetInfos struct {
		DBInstanceNetInfo []DBInstanceNetInfo `json:"DBInstanceNetInfo"`
	} `json:"DBInstanceNetInfos"`
	RequestId           string `json:"RequestId"`
	InstanceNetworkType string `json:"InstanceNetworkType"`
	SecurityIPMode      string `json:"SecurityIPMode"`
}

func (s *PolardbService) DoPolardbDescribedbinstancenetinfoRequest(id string) (*PolardbDescribedbinstancenetinfoResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstanceNetInfo
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceNetInfo", "")
	PolardbDescribedbinstancenetinfoResponse := &PolardbDescribedbinstancenetinfoResponse{}

	// Call request_params_handler

	request.QueryParams["DBInstanceId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceNetInfo", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancenetinfoResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceNetInfo", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstancenetinfoResponse.DBInstanceNetInfos.DBInstanceNetInfo) < 1 {
		return PolardbDescribedbinstancenetinfoResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstanceNetInfo", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedbinstancenetinfoResponse, nil
}

type PolardbDescribedbinstanceattributeResponse struct {
	Items struct {
		DBInstanceAttribute []struct {
			SlaveZones struct {
				SlaveZone []struct {
					ZoneId string `json:"ZoneId"`
				} `json:"SlaveZone"`
			} `json:"SlaveZones"`

			ReadOnlyDBInstanceIds struct {
				ReadOnlyDBInstanceId []struct {
					DBInstanceId string `json:"DBInstanceId"`
				} `json:"ReadOnlyDBInstanceId"`
			} `json:"ReadOnlyDBInstanceIds"`
			IPType                            string      `json:"IPType"`
			DBInstanceDiskUsed                string      `json:"DBInstanceDiskUsed"`
			GuardDBInstanceName               string      `json:"GuardDBInstanceName"`
			CanTempUpgrade                    bool        `json:"CanTempUpgrade"`
			TempUpgradeTimeStart              string      `json:"TempUpgradeTimeStart"`
			TempUpgradeTimeEnd                string      `json:"TempUpgradeTimeEnd"`
			TempUpgradeRecoveryTime           string      `json:"TempUpgradeRecoveryTime"`
			TempUpgradeRecoveryClass          string      `json:"TempUpgradeRecoveryClass"`
			TempUpgradeRecoveryCpu            int         `json:"TempUpgradeRecoveryCpu"`
			TempUpgradeRecoveryMemory         int         `json:"TempUpgradeRecoveryMemory"`
			TempUpgradeRecoveryMaxIOPS        string      `json:"TempUpgradeRecoveryMaxIOPS"`
			SqlGrammarCompatibility           string      `json:"SqlGrammarCompatibility"`
			TempUpgradeRecoveryMaxConnections string      `json:"TempUpgradeRecoveryMaxConnections"`
			InsId                             int         `json:"InsId"`
			DBInstanceId                      string      `json:"DBInstanceId"`
			PayType                           string      `json:"PayType"`
			DBInstanceClassType               string      `json:"DBInstanceClassType"`
			DBInstanceType                    string      `json:"DBInstanceType"`
			RegionId                          string      `json:"RegionId"`
			ConnectionString                  string      `json:"ConnectionString"`
			SlaveConnectionString             string      `json:"SlaveConnectionString"`
			Port                              string      `json:"Port"`
			Engine                            string      `json:"Engine"`
			EngineVersion                     string      `json:"EngineVersion"`
			DBInstanceClass                   string      `json:"DBInstanceClass"`
			DBInstanceMemory                  int         `json:"DBInstanceMemory"`
			DBInstanceStorage                 int         `json:"DBInstanceStorage"`
			VpcCloudInstanceId                string      `json:"VpcCloudInstanceId"`
			DBInstanceNetType                 interface{} `json:"DBInstanceNetType"`
			DBInstanceStatus                  string      `json:"DBInstanceStatus"`
			DBInstanceDescription             string      `json:"DBInstanceDescription"`
			LockMode                          string      `json:"LockMode"`
			LockReason                        string      `json:"LockReason"`
			ReadDelayTime                     string      `json:"ReadDelayTime"`
			DBMaxQuantity                     int         `json:"DBMaxQuantity"`
			AccountMaxQuantity                int         `json:"AccountMaxQuantity"`
			CreationTime                      string      `json:"CreationTime"`
			ExpireTime                        string      `json:"ExpireTime"`
			MaintainTime                      string      `json:"MaintainTime"`
			AvailabilityValue                 string      `json:"AvailabilityValue"`
			MaxIOPS                           int         `json:"MaxIOPS"`
			MaxConnections                    int         `json:"MaxConnections"`
			MasterInstanceId                  string      `json:"MasterInstanceId"`
			DBInstanceCPU                     string      `json:"DBInstanceCPU"`
			IncrementSourceDBInstanceId       string      `json:"IncrementSourceDBInstanceId"`
			GuardDBInstanceId                 string      `json:"GuardDBInstanceId"`
			ReplicateId                       string      `json:"ReplicateId"`
			TempDBInstanceId                  string      `json:"TempDBInstanceId"`
			SecurityIPList                    string      `json:"SecurityIPList"`
			ZoneId                            string      `json:"ZoneId"`
			InstanceNetworkType               string      `json:"InstanceNetworkType"`
			DBInstanceStorageType             string      `json:"DBInstanceStorageType"`
			EncryptionKey                     string      `json:"EncryptionKey"`
			AdvancedFeatures                  string      `json:"AdvancedFeatures"`
			Category                          string      `json:"Category"`
			AccountType                       string      `json:"AccountType"`
			SupportUpgradeAccountType         string      `json:"SupportUpgradeAccountType"`
			SupportCreateSuperAccount         string      `json:"SupportCreateSuperAccount"`
			VpcId                             string      `json:"VpcId"`
			VSwitchId                         string      `json:"VSwitchId"`
			ConnectionMode                    string      `json:"ConnectionMode"`
			CurrentKernelVersion              string      `json:"CurrentKernelVersion"`
			LatestKernelVersion               string      `json:"LatestKernelVersion"`
			CurrentKernelShowVersion          string      `json:"CurrentKernelShowVersion"`
			ResourceGroupId                   string      `json:"ResourceGroupId"`
			ReadonlyInstanceSQLDelayedTime    string      `json:"ReadonlyInstanceSQLDelayedTime"`
			SecurityIPMode                    string      `json:"SecurityIPMode"`
			TimeZone                          string      `json:"TimeZone"`
			Collation                         string      `json:"Collation"`
			DispenseMode                      string      `json:"DispenseMode"`
			MasterZone                        string      `json:"MasterZone"`
			AutoUpgradeMinorVersion           string      `json:"AutoUpgradeMinorVersion"`
			ProxyType                         int         `json:"ProxyType"`
			ConsoleVersion                    string      `json:"ConsoleVersion"`
			CpuType                           string      `json:"CpuType"`
			Vip                               string      `json:"Vip"`
			Vip_v6                            string      `json:"Vip_v6"`
			Vport                             string      `json:"Vport"`

			Extra struct {
				DBInstanceId struct {
					DBInstanceId []string `json:"DBInstanceId"`
				} `json:"DBInstanceId"`
				ReplicaGroupID            string `json:"ReplicaGroupID"`
				ReplicaGroupStatus        string `json:"ReplicaGroupStatus"`
				ActiveReplicaDBInstanceID string `json:"ActiveReplicaDBInstanceID"`
			} `json:"Extra"`
		} `json:"DBInstanceAttribute"`
	} `json:"Items"`
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DoPolardbDescribedbinstanceattributeRequest(id string) (*PolardbDescribedbinstanceattributeResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstanceAttribute
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceAttribute", "")
	PolardbDescribedbinstanceattributeResponse := &PolardbDescribedbinstanceattributeResponse{}

	// Call request_params_handler

	request.QueryParams["DBInstanceId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstanceattributeResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstanceattributeResponse.Items.DBInstanceAttribute) < 1 {
		return PolardbDescribedbinstanceattributeResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("PolardbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedbinstanceattributeResponse, nil
}

type PolardbDescribedbinstancemonitorResponse struct {
	RequestId string `json:"RequestId"`
	Period    string `json:"Period"`
}

func (s *PolardbService) DoPolardbDescribedbinstancemonitorRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribedbinstancemonitorResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstanceMonitor
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceMonitor", "")
	PolardbDescribedbinstancemonitorResponse := &PolardbDescribedbinstancemonitorResponse{}

	request.QueryParams["DBInstanceId"] = d.Id()

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceMonitor", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancemonitorResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceMonitor", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribedbinstancemonitorResponse, nil
}

type PolardbDescribeparametersResponse struct {
	ConfigParameters struct {
		DBInstanceParameter []struct {
			ParameterName        string `json:"ParameterName"`
			ParameterValue       string `json:"ParameterValue"`
			ParameterDescription string `json:"ParameterDescription"`
		} `json:"DBInstanceParameter"`
	} `json:"ConfigParameters"`

	RunningParameters struct {
		DBInstanceParameter []struct {
			ParameterName        string `json:"ParameterName"`
			ParameterValue       string `json:"ParameterValue"`
			ParameterDescription string `json:"ParameterDescription"`
		} `json:"DBInstanceParameter"`
	} `json:"RunningParameters"`
	RequestId     string `json:"RequestId"`
	Engine        string `json:"Engine"`
	EngineVersion string `json:"EngineVersion"`
}

func (s *PolardbService) DoPolardbDescribeparametersRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribeparametersResponse, error) {
	// api: polardb - 2024-01-30 - DescribeParameters
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeParameters", "")
	PolardbDescribeparametersResponse := &PolardbDescribeparametersResponse{}

	// Call request_params_handler
	request.QueryParams["DBInstanceId"] = d.Id()
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeParameters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribeparametersResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeParameters", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribeparametersResponse, nil
}

type PolardbParametersTemplateRecord struct {
	ForceModify          string `json:"ForceModify"`
	CheckingCode         string `json:"CheckingCode"`
	ParameterValue       string `json:"ParameterValue"`
	ForceRestart         string `json:"ForceRestart"`
	ParameterName        string `json:"ParameterName"`
	ParameterDescription string `json:"ParameterDescription"`
}

func (s *PolardbService) DoPolardbDescribeParameterTemplatesRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) ([]PolardbParametersTemplateRecord, error) {
	// api: polardb - 2024-01-30 - DescribeParameters
	action := "DescribeParameterTemplates"
	requQuery := map[string]interface{}{
		"DBInstanceId":  d.Id(),
		"Engine":        d.Get("engine").(string),
		"EngineVersion": d.Get("engine_version").(string),
	}
	resp, err := client.DoTeaRequest("GET", "polardb", "2024-01-30", action, "", nil, requQuery, nil)
	if err != nil {
		return nil, err
	}

	templates, err := jsonpath.Get("$.Parameters.TemplateRecord	", resp)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.Parameters.TemplateRecord", resp)
	}

	result := make([]PolardbParametersTemplateRecord, 0)
	// Convert map to JSON
	jsonData, _ := json.Marshal(templates)

	// Parse JSON into struct
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type PolardbDescribedbinstanceiparraylistResponse struct {
	Items struct {
		DBInstanceIPArray []struct {
			DBInstanceIPArrayName      string `json:"DBInstanceIPArrayName"`
			DBInstanceIPArrayAttribute string `json:"DBInstanceIPArrayAttribute"`
			SecurityIPType             string `json:"SecurityIPType"`
			SecurityIPList             string `json:"SecurityIPList"`
			WhitelistNetworkType       string `json:"WhitelistNetworkType"`
		} `json:"DBInstanceIPArray"`
	} `json:"Items"`
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DoPolardbDescribedbinstanceiparraylistRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribedbinstanceiparraylistResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstanceIPArrayList
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceIPArrayList", "")
	PolardbDescribedbinstanceiparraylistResponse := &PolardbDescribedbinstanceiparraylistResponse{}

	// Call request_params_handler

	request.QueryParams["DBInstanceId"] = d.Id()

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceIPArrayList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstanceiparraylistResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceIPArrayList", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribedbinstanceiparraylistResponse, nil
}

// not finish
func (s *PolardbService) WaitForPolardbConnection(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		parts, err := ParseResourceId(d.Id(), 2)
		object, err := s.DoPolardbDescribedbinstancenetinfoRequest(parts[0])
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil && object.DBInstanceNetInfos.DBInstanceNetInfo[0].ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, d.Id(), GetFunc(1), timeout, object.DBInstanceNetInfos.DBInstanceNetInfo[0].ConnectionString, d.Id(), errmsgs.ProviderERROR)
		}
	}
}
func (s *PolardbService) WaitForDBInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DoPolardbDescribedbinstancesRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil {
			if status == Deleted && len(object.Items.DBInstance) == 0 {
				break
			} else if strings.EqualFold(object.Items.DBInstance[0].DBInstanceStatus, string(status)) {
				break
			}
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Items.DBInstance[0].DBInstanceStatus, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *PolardbService) DescribeModifyParameterLog(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	now := time.Now().UTC()
	StartTime := now.Add(-5 * time.Minute)
	EndTime := now.Add(5 * time.Minute)
	requQuery := map[string]interface{}{
		"DBInstanceId": d.Id(),
		"StartTime":    StartTime.Format("2006-01-02T15:04Z"),
		"EndTime":      EndTime.Format("2006-01-02T15:04Z"),
	}
	for {
		isSyncing := false
		resp, err := client.DoTeaRequest("GET", "polardb", "2024-01-30", "DescribeModifyParameterLog", "", nil, requQuery, nil)
		if err != nil {
			return err
		}

		for _, i := range resp["Items"].(map[string]interface{})["ParameterChangeLog"].([]interface{}) {
			item := i.(map[string]interface{})
			if item["Status"].(string) == "Syncing" {
				time.Sleep(DefaultIntervalShort * time.Second)
				if time.Now().After(deadline) {
					return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, d.Id(), GetFunc(1), timeout, item["ParameterName"].(string), status, errmsgs.ProviderERROR)
				}
				isSyncing = true
				break
			}
		}
		if !isSyncing {
			break
		}
	}
	return nil
}

func (s *PolardbService) GetSecurityIps(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) ([]string, error) {
	object, err := s.DoPolardbDescribedbinstanceiparraylistRequest(d, client)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	var ips, separator string
	ipsMap := make(map[string]string)
	for _, ip := range object.Items.DBInstanceIPArray {
		if ip.DBInstanceIPArrayAttribute == "hidden" {
			continue
		}
		ips += separator + ip.SecurityIPList
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ips, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *PolardbService) RefreshParameters(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) error {

	diffParameters := make([]map[string]string, 0)
	if object, err := s.DoPolardbDescribeparametersRequest(d, client); err != nil {
		return errmsgs.WrapError(err)
	} else {
		for _, i := range object.RunningParameters.DBInstanceParameter {
			if i.ParameterName != "" {
				param := map[string]string{
					"name":  i.ParameterName,
					"value": i.ParameterValue,
				}
				diffParameters = append(diffParameters, param)
			}
		}
	}

	if err := d.Set("parameters", diffParameters); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func (s *PolardbService) ModifyDBSecurityIps(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, ips string) error {
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifySecurityIps", "")
	PolardbModifysecurityipsResponse := PolardbModifysecurityipsResponse{}
	request.QueryParams["DBInstanceId"] = d.Id()
	request.QueryParams["SecurityIps"] = ips

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_polardb_db_instance", "ModifySecurityIps", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifysecurityipsResponse)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_db_instance", "ModifySecurityIps", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

type PolardbDescribeinstanceautorenewalattributeResponse struct {
	Items struct {
		Item []struct {
			DBInstanceId string `json:"DBInstanceId"`
			RegionId     string `json:"RegionId"`
			Duration     int    `json:"Duration"`
			Status       string `json:"Status"`
			AutoRenew    string `json:"AutoRenew"`
		} `json:"Item"`
	} `json:"Items"`
	RequestId        string `json:"RequestId"`
	PageNumber       int    `json:"PageNumber"`
	TotalRecordCount int    `json:"TotalRecordCount"`
	PageRecordCount  int    `json:"PageRecordCount"`
}

func (s *PolardbService) DoPolardbDescribeinstanceautorenewalattributeRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribeinstanceautorenewalattributeResponse, error) {
	// api: polardb - 2024-01-30 - DescribeInstanceAutoRenewalAttribute
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeInstanceAutoRenewalAttribute", "")
	PolardbDescribeinstanceautorenewalattributeResponse := &PolardbDescribeinstanceautorenewalattributeResponse{}

	request.QueryParams["DBInstanceId"] = d.Id()
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeInstanceAutoRenewalAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribeinstanceautorenewalattributeResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeInstanceAutoRenewalAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribeinstanceautorenewalattributeResponse, nil
}

func (s *PolardbService) ModifyParameters(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) error {

	changed := make(map[string]string)
	if _, ok := d.GetOk("parameters"); !ok {
		return nil
	} else {
		// FIXME: d.HasChange("parameters") is abnormal, start manual judgment
		rawConfig := d.GetRawConfig()
		if parametersVal := rawConfig.GetAttr("parameters"); !parametersVal.IsNull() {
			parametersSet := parametersVal.AsValueSet()
			for _, value := range parametersSet.Values() {
				item := value.AsValueMap()
				key := item["name"].AsString()
				value := item["value"].AsString()
				id := fmt.Sprintf("parameters.%d.value", hashcode.String(key))
				old, _ := d.GetChange(id)
				// FIXME: d.GetChange("parameters") cannot get the correct new value
				if old.(string) != value {
					changed[key] = value
				}
			}
		}
	}

	if !d.Get("force_restart").(bool) {
		templates, err := s.DoPolardbDescribeParameterTemplatesRequest(d, client)
		if err != nil {
			return err
		}
		for _, template := range templates {
			if template.ForceRestart != "true" {
				continue
			}
			key := template.ParameterName
			for k, _ := range changed {
				if key == k {
					return errmsgs.WrapError(fmt.Errorf("Modifying RDS instance's parameter '%s' requires setting 'force_restart = true'.", key))
				}
			}
		}
	}

	if len(changed) > 0 {
		request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "ModifyParameter", "")

		request.QueryParams["DBInstanceId"] = d.Id()
		if d.Get("force_restart").(bool) {
			request.QueryParams["Forcerestart"] = "true"
		} else {
			request.QueryParams["Forcerestart"] = "false"
		}
		cfg, _ := json.Marshal(changed)
		request.QueryParams["Parameters"] = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
			return errmsgs.WrapError(err)
		}
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "ModifyParameters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		if err := s.DescribeModifyParameterLog(d, client, Running, DefaultLongTimeout); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

type DescribeParameterTemplatesResponse struct {
	*responses.BaseResponse
	RequestId      string     `json:"RequestId" xml:"RequestId"`
	Engine         string     `json:"Engine" xml:"Engine"`
	ParameterCount string     `json:"ParameterCount" xml:"ParameterCount"`
	EngineVersion  string     `json:"EngineVersion" xml:"EngineVersion"`
	Parameters     Parameters `json:"Parameters" xml:"Parameters"`
}
type Parameters struct {
	TemplateRecord []TemplateRecord `json:"TemplateRecord" xml:"TemplateRecord"`
}
type TemplateRecord struct {
	CheckingCode         string `json:"CheckingCode" xml:"CheckingCode"`
	ParameterName        string `json:"ParameterName" xml:"ParameterName"`
	ParameterValue       string `json:"ParameterValue" xml:"ParameterValue"`
	ForceModify          string `json:"ForceModify" xml:"ForceModify"`
	ForceRestart         string `json:"ForceRestart" xml:"ForceRestart"`
	ParameterDescription string `json:"ParameterDescription" xml:"ParameterDescription"`
}

func (s *PolardbService) PolardbDBInstanceStateRefreshFunc(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoPolardbDescribedbinstancesRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"ServiceUnavailable"}) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		if slices.Contains(failStates, object.Items.DBInstance[0].DBInstanceStatus) {
				return object, object.Items.DBInstance[0].DBInstanceStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Items.DBInstance[0].DBInstanceStatus))
			}
		return object, object.Items.DBInstance[0].DBInstanceStatus, nil
	}
}

func (s *PolardbService) PolardbDBInstanceTdeStateRefreshFunc(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDBInstanceTDE(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		if slices.Contains(failStates, object["TDEStatus"].(string)) {
				return object, object["TDEStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["TDEStatus"].(string)))
			}
		return object, object["TDEStatus"].(string), nil
	}
}

func (s *PolardbService) DescribeDBInstanceTDE(id string) (map[string]interface{}, error) {
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceTDE", "")
	request.QueryParams["DBInstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceTDE", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	return result, nil
}

func (s *PolardbService) PolardbDBInstanceSslStateRefreshFunc(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDBInstanceSSL(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["SSLEnabled"].(string) == failState {
				return object, object["SSLEnabled"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["SSLEnabled"].(string)))
			}
		}
		return object, object["SSLEnabled"].(string), nil
	}
}

func (s *PolardbService) DescribeDBInstanceSSL(id string) (map[string]interface{}, error) {
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceSSL", "")
	request.QueryParams["DBInstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceSSL", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	return result, nil
}

func (s *PolardbService) WaitForDBConnection(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		response, err := s.DescribeDBConnection(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		data := response
		if data != nil && data.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, "", id, errmsgs.ProviderERROR)
		}
	}
}

func (s *PolardbService) WaitForConnectionDBInstance(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DoPolardbDescribedbinstancesRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil && strings.EqualFold(object.Items.DBInstance[0].DBInstanceStatus, string(status)) {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Items.DBInstance[0].DBInstanceStatus, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *PolardbService) DoPolardbDescribedbinstancesRequest(id string) (*PolardbDescribedbinstancesResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstances
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstances", "")
	PolardbDescribedbinstancesResponse := &PolardbDescribedbinstancesResponse{}
	request.QueryParams["DBInstanceId"] = id
	request.QueryParams["InstanceLevel"] = "1"
	request.QueryParams["PageNumber"] = "1"
	request.QueryParams["PageSize"] = "1"
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return PolardbDescribedbinstancesResponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancesResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstancesResponse.Items.DBInstance) < 1 {
		return PolardbDescribedbinstancesResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("PolardbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return PolardbDescribedbinstancesResponse, nil
}

func (s *PolardbService) Describedbinstances(id string) (*PolardbDescribedbinstancesResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstances
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstances", "")
	PolardbDescribedbinstancesResponse := &PolardbDescribedbinstancesResponse{}
	request.QueryParams["DBInstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return PolardbDescribedbinstancesResponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancesResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstancesResponse.Items.DBInstance) < 1 {
		return PolardbDescribedbinstancesResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("PolardbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedbinstancesResponse, nil
}

func (s *PolardbService) DescribeDBSecurityIps(instanceId string) (*PolardbDescribedbinstanceiparraylistResponse, error) {
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceIPArrayList", "")
	PolardbDescribedbinstanceiparraylistResponse := &PolardbDescribedbinstanceiparraylistResponse{}

	// Call request_params_handler

	request.QueryParams["DBInstanceId"] = instanceId

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceIPArrayList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstanceiparraylistResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceIPArrayList", errmsgs.AlibabacloudStackSdkGoERROR)

	}
	return PolardbDescribedbinstanceiparraylistResponse, nil

}

func (s *PolardbService) flattenDBSecurityIPs(resp *PolardbDescribedbinstanceiparraylistResponse) []map[string]interface{} {
	list := resp.Items.DBInstanceIPArray
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"security_ips": i.SecurityIPList,
		}
		result = append(result, l)
	}
	return result
}

func (s *PolardbService) DescribeDBConnection(id string) (*DBInstanceNetInfo, error) {
	parts, _ := ParseResourceId(id, 2)
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceNetInfo", "")
	PolardbDescribedbinstancenetinfoResponse := &PolardbDescribedbinstancenetinfoResponse{}

	// Call request_params_handler

	request.QueryParams["DBInstanceId"] = parts[0]

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidCurrentConnectionString.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, errmsgs.WrapError(err)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancenetinfoResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceNetInfo", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstancenetinfoResponse.DBInstanceNetInfos.DBInstanceNetInfo) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstanceNetInfo", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	object := PolardbDescribedbinstancenetinfoResponse.DBInstanceNetInfos.DBInstanceNetInfo
	for _, o := range object {
		if strings.HasPrefix(o.ConnectionString, parts[1]+".") {
			return &o, nil
		}
	}

	return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

type PolardbBackupJob struct {
	BackupProgressStatus string `json:"BackupProgressStatus"`
	BackupStatus         string `json:"BackupStatus"`
	JobMode              string `json:"JobMode"`
	Process              string `json:"Process"`
	TaskAction           string `json:"TaskAction"`
	BackupJobId          int    `json:"BackupJobId"`
	BackupId             string `json:"BackupId"`
}

type DescribeBackupTasksResponse struct {
	EagleEyeTraceId string `json:"eagleEyeTraceId"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	ResponseVersion string `json:"responseVersion"`
	RequestId       string `json:"RequestId"`
	Success         bool   `json:"success"`
	RequestID       string `json:"requestId"` // Note that field case may not match [RequestId](file://d:\terraform\terraform-provider-apsarastack\alibabacloudstack\resource_apsarastack_ack_cluster.go#L657-L657)
	Items           struct {
		BackupJob []PolardbBackupJob `json:"BackupJob"`
	} `json:"Items"`
}

func (s *PolardbService) DoPolardbDescribebackupTaskRequest(db_instance_id, job_id string) (PolardbBackupJob, error) {
	// api: polardb - 2024-01-30 - DescribeBackups
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeBackupTasks", "")
	DescribeBackupTasksResponseObj := &DescribeBackupTasksResponse{}
	var backup_job PolardbBackupJob
	request.QueryParams["DBInstanceId"] = db_instance_id
	request.QueryParams["BackupJobId"] = job_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return backup_job, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return backup_job, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBackups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DescribeBackupTasksResponseObj)

	if err != nil {
		return backup_job, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBackups", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, v := range DescribeBackupTasksResponseObj.Items.BackupJob {
		if fmt.Sprint(v.BackupJobId) == job_id {
			backup_job = v
			break
		}
	}
	return backup_job, nil
}

func (s *PolardbService) PolardbDescribebackupTaskStateRefreshFunc(db_instance_id string, job_id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoPolardbDescribebackupTaskRequest(db_instance_id, job_id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.BackupStatus == failState {
				return object, object.BackupStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.BackupStatus))
			}
		}
		return object, object.BackupStatus, nil
	}
}

type PolardbbackupData struct {
	BackupMethod              string `json:"BackupMethod"`
	BackupIntranetDownloadURL string `json:"BackupIntranetDownloadURL"`
	BackupMode                string `json:"BackupMode"`
	BackupSize                int    `json:"BackupSize"`
	BackupId                  int    `json:"BackupId"`
	SlaveStatus               string `json:"SlaveStatus"`
	HostInstanceID            int    `json:"HostInstanceID"`
	BackupDBNames             string `json:"BackupDBNames"`
	StoreStatus               string `json:"StoreStatus"`
	DBInstanceId              string `json:"DBInstanceId"`
	BackupDownloadURL         string `json:"BackupDownloadURL"`
	BackupEndTime             string `json:"BackupEndTime"`
	BackupStartTime           string `json:"BackupStartTime"`
	BackupType                string `json:"BackupType"`
	MetaStatus                string `json:"MetaStatus"`
	BackupScale               string `json:"BackupScale"`
	BackupStatus              string `json:"BackupStatus"`
	BackupLocation            string `json:"BackupLocation"`
}

type PolardbDescribebackupsResponse struct {
	Items struct {
		Backup []PolardbbackupData `json:"Backup"`
	} `json:"Items"`
	TotalRecordCount int    `json:"TotalRecordCount"`
	PageRecordCount  int    `json:"PageRecordCount"`
	RequestId        string `json:"RequestId"`
	PageNumber       int    `json:"PageNumber"`
}

func (s *PolardbService) DoPolardbDescribebackupsRequest(id string) (*PolardbbackupData, error) {
	// api: polardb - 2024-01-30 - DescribeBackups
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeBackups", "")
	PolardbDescribebackupsResponseObj := &PolardbDescribebackupsResponse{}
	param := strings.Split(id, ":")
	request.QueryParams["DBInstanceId"] = param[0]
	request.QueryParams["BackupId"] = param[1]

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBackups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribebackupsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBackups", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, v := range PolardbDescribebackupsResponseObj.Items.Backup {
		if fmt.Sprint(v.BackupId) == param[1] {
			return &v, nil
		}
	}

	// XXX: safety net logic
	delete(request.QueryParams, "BackupId")
	now := time.Now().UTC()
	oneDayAgo := now.Add(-24 * time.Hour)
	format := "2006-01-02T15:04Z"
	request.QueryParams["StartTime"] = oneDayAgo.Format(format)
	request.QueryParams["EndTime"] = now.Format(format)
	request.QueryParams["PageSize"] = "10"
	pageNumber := 1
	for {
		request.QueryParams["PageNumber"] = strconv.Itoa(pageNumber)
		bresponse, err := s.client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBackups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribebackupsResponseObj)

		if err != nil {
			return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBackups", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		for _, v := range PolardbDescribebackupsResponseObj.Items.Backup {
			if fmt.Sprint(v.BackupId) == param[1] {
				return &v, nil
			}
		}
		if len(PolardbDescribebackupsResponseObj.Items.Backup) < 10 {
			break
		}
		pageNumber += 1
	}

	return nil, errmsgs.Error(errmsgs.NotFoundWithResponse, id)
}

func (s *PolardbService) PolardbDescribebackupsStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoPolardbDescribebackupsRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.BackupStatus == failState {
				return object, object.BackupStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.BackupStatus))
			}
		}
		return object, object.BackupStatus, nil
	}
}

func (s *PolardbService) SetInstanceTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		remove := oraw.(map[string]interface{})
		add := nraw.(map[string]interface{})

		if len(remove) > 0 {
			b, err := json.Marshal(remove)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			reqQuery := map[string]interface{}{
				"DBInstanceId": d.Id(),
				"Tags":         string(b),
			}
			if _, err := s.client.DoTeaRequest(
				"POST", "polardb", "2024-01-30", "RemoveTagsFromResource", "", nil, reqQuery, nil); err != nil {
				return err
			}
		}

		if len(add) > 0 {
			b, err := json.Marshal(add)
			if err != nil {
				return errmsgs.WrapError(err)
			}

			reqQuery := map[string]interface{}{
				"DBInstanceId": d.Id(),
				"Tags":         string(b),
			}
			if _, err := s.client.DoTeaRequest(
				"POST", "polardb", "2024-01-30", "AddTagsToResource", "", nil, reqQuery, nil); err != nil {
				return err
			}
		}

	}

	return nil
}

func (s *PolardbService) describeTags(d *schema.ResourceData) ([]Tag, error) {
	reqQuery := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	resp, err := s.client.DoTeaRequest("GET", "polardb", "2024-01-30", "DescribeTags", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	tagSet := make([]Tag, 0)
	tags := resp["Items"].(map[string]interface{})["TagInfos"].([]interface{})
	for _, t := range tags {
		tag := t.(map[string]interface{})
		tagSet = append(tagSet, Tag{
			Key:   tag["TagKey"].(string),
			Value: tag["TagValue"].(string),
		})
	}
	return tagSet, nil
}

func (s *PolardbService) tagsToMap(tags []Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func (s *PolardbService) ignoreTag(t Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *PolardbService) tagsToString(tags []Tag) string {
	v, _ := json.Marshal(s.tagsToMap(tags))

	return string(v)
}

type PolardbDescribedbproxyResponse struct {
	DBProxyConnectStringItems struct {
		DBProxyConnectStringItems []struct {
			DBProxyEndpointId               int    `json:"DBProxyEndpointId"`
			DBProxyConnectString            string `json:"DBProxyConnectString"`
			DBProxyConnectStringPort        string `json:"DBProxyConnectStringPort"`
			DBProxyConnectStringNetType     string `json:"DBProxyConnectStringNetType"`
			DBProxyVpcInstanceId            string `json:"DBProxyVpcInstanceId"`
			DBProxyEndpointName             string `json:"DBProxyEndpointName"`
			DBProxyConnectStringNetWorkType int    `json:"DBProxyConnectStringNetWorkType"`
		} `json:"DBProxyConnectStringItems"`
	} `json:"DBProxyConnectStringItems"`

	DBProxyEndpointItems struct {
		DBProxyEndpointItems []struct {
			DBProxyEndpointName    string `json:"DBProxyEndpointName"`
			DBProxyEndpointType    string `json:"DBProxyEndpointType"`
			DBProxyEndpointAliases string `json:"DBProxyEndpointAliases"`
			DBProxyReadWriteMode   string `json:"DBProxyReadWriteMode"`
		} `json:"DBProxyEndpointItems"`
	} `json:"DBProxyEndpointItems"`
	RequestId                          string `json:"RequestId"`
	DBProxyServiceStatus               string `json:"DBProxyServiceStatus"`
	DBProxyInstanceType                string `json:"DBProxyInstanceType"`
	DBProxyInstanceNum                 int    `json:"DBProxyInstanceNum"`
	DBProxyInstanceStatus              string `json:"DBProxyInstanceStatus"`
	DBProxyInstanceCurrentMinorVersion string `json:"DBProxyInstanceCurrentMinorVersion"`
	DBProxyInstanceLatestMinorVersion  string `json:"DBProxyInstanceLatestMinorVersion"`
}

func (s *PolardbService) DoPolardbDescribedbproxyRequest(id string) (*PolardbDescribedbproxyResponse, error) {
	// api: Polardb - 2014-08-15 - DescribeDBProxy
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBProxy", "")
	PolardbDescribedbproxyResponseObj := &PolardbDescribedbproxyResponse{}

	request.QueryParams["DBInstanceId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBProxy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbproxyResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBProxy", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribedbproxyResponseObj, nil
}

func (s *PolardbService) PolardbProxyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoPolardbDescribedbproxyRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBProxyInstanceStatus == failState {
				return object, object.DBProxyInstanceStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.DBProxyInstanceStatus))
			}
		}
		return object, object.DBProxyInstanceStatus, nil
	}
}

type PolarDBProxyEndpoint struct {
	ReadOnlyInstanceDistributionType string `json:"ReadOnlyInstanceDistributionType"`
	DBProxyConnectString             string `json:"DBProxyConnectString"`
	DBProxyEndpointId                string `json:"DBProxyEndpointId"`
	DBProxyFeatures                  string `json:"DBProxyFeatures"`
	ReadOnlyInstanceWeight           string `json:"ReadOnlyInstanceWeight"`
	ReadOnlyInstanceMaxDelayTime     string `json:"ReadOnlyInstanceMaxDelayTime"`
	DBProxyConnectStringNetType      string `json:"DBProxyConnectStringNetType"`
	DBProxyConnectStringPort         string `json:"DBProxyConnectStringPort"`
}

func (s *PolardbService) DoDescribeDBProxyEndpointRequest(instanceId string) (*PolarDBProxyEndpoint, error) {
	// api: R-kvstore - 2015-01-01 - DescribeBackupTasks
	request := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBProxyEndpoint", "")
	dBProxyEndpoint := &PolarDBProxyEndpoint{}
	request.QueryParams["DBInstanceId"] = instanceId

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBackupTasks", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &dBProxyEndpoint)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBackupTasks", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return dBProxyEndpoint, nil
}

func (s *PolardbService) CheckCloudResourceAuthorized() (string, error) {
	req := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "CheckCloudResourceAuthorized", "")
	req.QueryParams["TargetRegionId"] = s.client.RegionId
	var arnresp RoleARN
	bresponse, err := s.client.ProcessCommonRequest(req)
	addDebug(req.GetActionName(), bresponse, req, req.QueryParams)
	if err != nil {
		if bresponse == nil {
			return "", errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return "", errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_polardb_db_instance", "CheckCloudResourceAuthorized", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &arnresp)
	if err != nil {
		return "", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_db_instance", "CheckCloudResourceAuthorized", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	arnrole := arnresp.RoleArn
	return arnrole, errmsgs.WrapError(err)
}

func (s *PolardbService) DescribeDBInstanceEncryptionKey(id string) string {
	var err error
	req := s.client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeDBInstanceEncryptionKey", "")
	req.QueryParams["DBInstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(req)
	addDebug(req.GetActionName(), bresponse, req, req.QueryParams)
	response := make(map[string]interface{})
	if err != nil {
		if bresponse == nil {
			return ""
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		log.Printf("[DEBUG] Polardb %s : DescribeDBInstanceEncryptionKey %s", id, errmsg)
		return ""
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		log.Printf("[DEBUG] Polardb %s : DescribeDBInstanceEncryptionKey json.Unmarshal Error! \n %s", id, bresponse.GetHttpContentString())
		return ""
	}
	encryptionKey, err := jsonpath.Get("$.EncryptionKey", response)
	if err != nil {
		return ""
	}
	return encryptionKey.(string)
}

func (s *PolardbService) DescribePolardbClusterInstance(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{"DBClusterId": id}
	response := make(map[string]interface{})
	var err error
	retry := 0
	for retry < 5 {
		response, err = s.client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeDBClusterAttribute", "", nil, reqQuery, nil)
		if err == nil {
			break
		} else if errmsgs.IsExpectedErrors(err, []string{"Forbidden.RAM"}) {
			time.Sleep(time.Duration(5) * time.Second)
			retry++
		} else {
			return response, err
		}
	}

	if response["DBClusterId"] == nil || response["DBClusterId"].(string) == "" {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("PolarDB shared instance %s was not found", id))
	}
	return response, nil
}

func (s *PolardbService) PolardbClusterInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolardbClusterInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["DBClusterStatus"] == failState {
				return object, object["DBClusterStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["DBClusterStatus"]))
			}
		}
		return object, object["DBClusterStatus"].(string), nil
	}
}

func (s *PolardbService) WaitPolardbClusterInstanceAllDbNodesRunning(d *schema.ResourceData) error {
	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		object, err := s.DescribePolardbClusterInstance(d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		readonlyNodeNum := d.Get("readonly_node_num").(int)
		dbNodeNum := readonlyNodeNum + 1 // sum a rw node
		if d.Get("hot_standby_cluster") == "standby" {
			dbNodeNum += 1
		}
		if len(object["DBNodes"].([]interface{})) != dbNodeNum {
			return resource.RetryableError(fmt.Errorf("The Number for Node in Instance %s is not enough", d.Id()))
		}

		for _, dbNode := range object["DBNodes"].([]interface{}) {
			nodeInfo := dbNode.(map[string]interface{})
			if nodeInfo["DBNodeStatus"].(string) != "Running" {
				return resource.RetryableError(fmt.Errorf("The Status for Node %s in Instance %s is %s", nodeInfo["DBNodeId"].(string), d.Id(), nodeInfo["DBNodeStatus"].(string)))
			}
		}
		return nil
	})
}

func (s *PolardbService) DescribeDBClusterEndpoints(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{"DBClusterId": id}
	response, err := s.client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeDBClusterEndpoints", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *PolardbService) ModifySecurityIps(instance_id string, old, new interface{}) error {
	var oldMap, newMap map[string]interface{}
	if old != nil {
		oldMap = old.(map[string]interface{})
	}
	if new != nil {
		newMap = new.(map[string]interface{})
	}
	if _, ok := newMap["default"]; ok {
		return errmsgs.WrapError(errmsgs.Error("Security IP group name cannot be `default`!"))
	}
	requests := make([]map[string]interface{}, 0)
	for old_group_name, old_ips := range oldMap {
		new_ips, ok := newMap[old_group_name]
		if !ok {
			requests = append(requests, map[string]interface{}{
				"DBClusterIPArrayName": old_group_name,
				"SecurityIps":          "",
				"DBClusterId":          instance_id,
			})
		} else {
			if old_ips != new_ips {
				requests = append(requests, map[string]interface{}{
					"DBClusterIPArrayName": old_group_name,
					"SecurityIps":          new_ips,
					"DBClusterId":          instance_id,
				})
			}
		}
	}
	for new_group_name, new_ips := range newMap {
		_, ok := oldMap[new_group_name]
		if !ok {
			requests = append(requests, map[string]interface{}{
				"DBClusterIPArrayName": new_group_name,
				"SecurityIps":          new_ips,
				"DBClusterId":          instance_id,
			})
		}
	}
	for _, request := range requests {
		_, err := s.client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterAccessWhiteList", "", nil, request, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PolardbService) ModifySecurityGroups(instance_id string, securityGroups []interface{}) error {
	var securityGroupIds string
	if len(securityGroups) > 0 {
		s := make([]string, len(securityGroups))
		for _, v := range securityGroups {
			s = append(s, v.(string))
		}
		securityGroupIds = strings.Join(s, ",")
	}
	request := map[string]interface{}{
		"DBClusterId":      instance_id,
		"WhiteListType":    "SecurityGroup",
		"SecurityGroupIds": securityGroupIds,
	}
	_, err := s.client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterAccessWhiteList", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_cluster_instance", "ModifyDBClusterAccessWhiteList => ModifySecurityGroups", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func (s *PolardbService) DoPolardbxDescribeClusterParametersRequest(id string) (map[string]interface{}, error) {
	// api: polardb - 2017-08-01 - DescribeParameters
	request := s.client.NewCommonRequest("GET", "polardb", "2017-08-01", "DescribeDBClusterParameters", "")
	request.QueryParams["DBClusterId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DoPolardbxDescribeClusterParametersRequest", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "DoPolardbxDescribeClusterParametersRequest", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return response, nil
}

func (s *PolardbService) DescribeClusterParameters(d *schema.ResourceData) (map[string]interface{}, error) {
	diffParameters := make(map[string]interface{})
	if response, err := s.DoPolardbxDescribeClusterParametersRequest(d.Id()); err != nil {
		return nil, errmsgs.WrapError(err)
	} else {
		tf_parameters := d.Get("parameters").(*schema.Set).List()
		tf_map := make(map[string]string)
		for _, v := range tf_parameters {
			m := v.(map[string]interface{})
			tf_map[m["name"].(string)] = m["value"].(string)
		}
		parameters, err := jsonpath.Get("$.RunningParameters.Parameter", response)
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
		for _, parameter := range parameters.([]interface{}) {
			parameter_map := parameter.(map[string]interface{})
			parameterName := parameter_map["ParameterName"].(string)
			if _, ok := tf_map[parameterName]; ok && parameterName != "" {
				diffParameters[parameterName] = parameter_map["ParameterValue"]
			}
		}
	}
	return diffParameters, nil
}

func (s *PolardbService) RefreshClusterParameters(d *schema.ResourceData) error {

	parameters := make([]map[string]string, 0)
	if diffParameters, err := s.DescribeClusterParameters(d); err != nil {
		return errmsgs.WrapError(err)
	} else {
		for name := range diffParameters {
			param := map[string]string{
				"name":  name,
				"value": diffParameters[name].(string),
			}
			parameters = append(parameters, param)
		}
	}
	if err := d.Set("parameters", parameters); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func (s *PolardbService) PolardbClusterParametersStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := s.DoPolardbxDescribeClusterParametersRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		parameters, err := jsonpath.Get("$.RunningParameters.Parameter", response)
		if err != nil {
			return nil, "", errmsgs.WrapError(err)
		}
		status := "Normal"
		var object map[string]interface{}
		for _, v := range parameters.([]interface{}) {
			parameter := v.(map[string]interface{})
			if parameter["ParameterStatus"].(string) == "Modifying" {
				status = "Modifying"
				object = parameter
				break
			}
		}

		for _, failState := range failStates {
			if status == failState {
				return object, status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, status))
			}
		}
		return object, status, nil
	}
}

func (s *PolardbService) ModifyClusterParameters(d *schema.ResourceData) error {
	request := s.client.NewCommonRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterParameters", "")
	request.QueryParams["DBClusterId"] = d.Id()
	config := make(map[string]string)
	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	var parameters map[string]interface{}
	if diffParameters, err := s.DescribeClusterParameters(d); err != nil {
		return errmsgs.WrapError(err)
	} else {
		parameters = diffParameters
	}
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			if v, exist := parameters[key]; exist && v.(string) == value {
				continue
			}
			config[key] = value
		}
		if len(config) == 0 {
			return nil
		}
		cfg, _ := json.Marshal(config)
		request.QueryParams["Parameters"] = string(cfg)
		// wait instance status is Normal before modifying
		bresponse, err := s.client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		stateConf := BuildStateConf([]string{"Modifying"}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.PolardbClusterParametersStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	//d.SetPartial(attribute)
	return nil
}

func (s *PolardbService) DescribeDBClusterTDE(id string) (string, string, error) {
	req := s.client.NewCommonRequest("GET", "polardb", "2017-08-01", "DescribeDBClusterTDE", "")
	req.QueryParams["DBClusterId"] = id
	bresponse, err := s.client.ProcessCommonRequest(req)
	addDebug(req.GetActionName(), bresponse, req, req.QueryParams)
	if err != nil {
		if bresponse == nil {
			return "", "", errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return "", "", errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"polardbService", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return "", "", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"polardbService", "DescribeDBClusterTDE", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	tdeStatus := response["TDEStatus"].(string)
	encryptionKey := response["EncryptionKey"].(string)
	return tdeStatus, encryptionKey, nil
}

func (s *PolardbService) DescribeDBClusterSSL(id string) (bool, error) {
	req := s.client.NewCommonRequest("GET", "polardb", "2017-08-01", "DescribeDBClusterSSL", "")
	req.QueryParams["DBClusterId"] = id
	bresponse, err := s.client.ProcessCommonRequest(req)
	addDebug(req.GetActionName(), bresponse, req, req.QueryParams)
	if err != nil {
		if bresponse == nil {
			return false, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return false, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"polardbService", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return false, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"polardbService", "DescribeDBClusterSSL", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	items := response["Items"].([]interface{})
	if len(items) == 0 {
		return false, nil
	}
	return true, nil
}

func (s *PolardbService) GetPolardbClusterClassData(dbType, dbVersion, dbClass string) (classData map[string]interface{}, err error) {
	reqQuery := map[string]interface{}{
		"pageStart":    1,
		"pageSize":     500,
		"label":        "true",
		"resourceType": "POLARDB",
		"status":       "Available",
		"dbVersion":    dbVersion,
		"dbType":       dbType,
	}
	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}
	response, err := s.client.DoTeaRequest("POST", "ascm", "2019-05-10", "SelectCommonSpec", "/ascm/manage/saleconf/commonSpec/select", reqHeader, reqQuery, nil)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "polardbService", "GetPolardbClusterClassData", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(response["data"].([]interface{})) == 0 {
		return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("polardb_cluster_instance_type", dbClass))
	}
	for _, v := range response["data"].([]interface{}) {
		data := v.(map[string]interface{})
		if data["dbNodeClass"].(string) == dbClass {
			return data, nil
		}
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("polardb_cluster_instance_type", dbClass))
}

func (s *PolardbService) DescribePolardbClusterAccount(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, expected DBClusterId:AccountName")
	}
	dbClusterId := parts[0]
	accountName := parts[1]

	reqQuery := map[string]interface{}{
		"DBClusterId": dbClusterId,
	}

	response, err := s.client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeAccounts", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	if accounts, ok := response["Accounts"].([]interface{}); ok {
		for _, acc := range accounts {
			account := acc.(map[string]interface{})
			if account["AccountName"].(string) == accountName {
				return account, nil
			}
		}
	}

	return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *PolardbService) PolardbClusterAccountStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolardbClusterAccount(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["AccountStatus"].(string) == failState {
				return object, object["AccountStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["AccountStatus"].(string)))
			}
		}
		return object, object["AccountStatus"].(string), nil
	}
}

func (s *PolardbService) CreatePolardbClusterAccount(clusterId, accountName, accountPassword string) (err error) {

	reqQuery := map[string]interface{}{
		"DBClusterId":     clusterId,
		"AccountName":     accountName,
		"AccountType":     "Normal",
		"AccountPassword": accountPassword,
	}
	_, err = s.client.DoTeaRequest("POST", "polardb", "2017-08-01", "CreateAccount", "", nil, reqQuery, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *PolardbService) DeletePolardClusterAccount(clusterId, accountName string) (err error) {

	reqQuery := map[string]interface{}{
		"AccountName": accountName,
		"DBClusterId": clusterId,
	}

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := s.client.DoTeaRequest("POST", "polardb", "2017-08-01", "DeleteAccount", "", nil, reqQuery, nil)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, accountName, "DeleteAccount", errmsgs.AlibabacloudStackSdkGoERROR, "")
			return resource.RetryableError(err)
		}
		return nil
	})
	return nil
}

func (s *PolardbService) DescribePolardbClusterDatabase(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, expected DBClusterId:DBName")
	}
	dbClusterId := parts[0]
	dbName := parts[1]

	reqQuery := map[string]interface{}{
		"DBClusterId": dbClusterId,
		"DBName":      dbName,
	}

	response, err := s.client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeDatabases", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	databases := response["Databases"].(map[string]interface{})["Database"].([]interface{})
	if len(databases) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("database %s not found in cluster %s", dbName, dbClusterId))
	}

	return databases[0].(map[string]interface{}), nil
}

func (s *PolardbService) DescribePolardbClusterProxy(id string) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"DBClusterId": id,
	}

	response, err := s.client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeDBClusterProxy", "", nil, query, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *PolardbService) PolardbClusterProxyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolardbClusterProxy(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		if _, exist := object["DBProxyClusterId"]; !exist {
			return nil, "", nil
		}

		for _, failState := range failStates {
			if object["DBProxyClusterStatus"] == failState {
				return object, object["DBProxyClusterStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["DBClusterStatus"]))
			}
		}
		return object, object["DBProxyClusterStatus"].(string), nil
	}
}

func (s *PolardbService) DescribePolardbClusterBackupPolicy(id string) (map[string]interface{}, error) {
	// DescribeBackupPolicy
	backupQuery := map[string]interface{}{
		"DBClusterId": id,
	}
	backupResp, err := s.client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeBackupPolicy", "", nil, backupQuery, nil)
	if err != nil {
		return nil, err
	}

	return backupResp, nil
}

func (s *PolardbService) DescribePolardbClusterLogBackupPolicy(id string) (map[string]interface{}, error) {
	// DescribeLogBackupPolicy
	logQuery := map[string]interface{}{
		"DBClusterId": id,
	}
	logResp, err := s.client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeLogBackupPolicy", "", nil, logQuery, nil)
	if err != nil {
		return nil, err
	}
	return logResp, nil
}

func (s *PolardbService) DescribeDBReadWriteSplittingConnection(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{
		"DBInstanceId": id,
	}
	response, err := s.client.DoTeaRequest("GET", "polardb", "2024-01-30", "DescribeDBProxyEndpoint", "", nil, reqQuery, nil)
	return response, err
}
