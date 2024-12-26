package connectivity

import (
	"encoding/json"
	"log"

	roaCS "github.com/alibabacloud-go/cs-20151215/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/maxcompute"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/fc-go-sdk"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"sync"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AlibabacloudStackClient struct {
	SourceIp                     string
	SecureTransport              string
	Region                       Region
	RegionId                     string
	Domain                       string
	AccessKey                    string
	SecretKey                    string
	Department                   string
	ResourceGroup                string
	ResourceGroupId              int
	Config                       *Config
	teaSdkConfig                 rpc.Config
	accountId                    string
	roleId                       int
	Conns                        map[ServiceCode]*sdk.Client
	ascmconn                     *sdk.Client
	ecsconn                      *ecs.Client
	accountIdMutex               sync.RWMutex
	roleIdMutex                  sync.RWMutex
	vpcconn                      *vpc.Client
	slbconn                      *slb.Client
	csconn                       *cs.Client
	polarDBconn                  *polardb.Client
	cdnconn                      *cdn.Client
	kmsconn                      *kms.Client
	bssopenapiconn               *bssopenapi.Client
	rdsconn                      *rds.Client
	ramconn                      *ram.Client
	essconn                      *ess.Client
	gpdbconn                     *gpdb.Client
	drdsconn                     *drds.Client
	elasticsearchconn            *elasticsearch.Client
	hbaseconn                    *hbase.Client
	adbconn                      *adb.Client
	ossconn                      *oss.Client
	rkvconn                      *r_kvstore.Client
	fcconn                       *fc.Client
	ddsconn                      *dds.Client
	onsconn                      *ons.Client
	logconn                      *sls.Client
	logpopconn                   *slsPop.Client
	dnsconn                      *alidns.Client
	edasconn                     *edas.Client
	creeconn                     *cr_ee.Client
	cmsconn                      *cms.Client
	maxcomputeconn               *maxcompute.Client
	alikafkaconn                 *alikafka.Client
	otsconn                      *ots.Client
	OtsInstanceName              string
	tablestoreconnByInstanceName map[string]*tablestore.TableStoreClient
	dhconn                       datahub.DataHubApi
	cloudapiconn                 *cloudapi.Client
	Eagleeye                     EagleEye
}

const (
	ApiVersion20140526 = ApiVersion("2014-05-26")
	ApiVersion20160815 = ApiVersion("2016-08-15")
	ApiVersion20140515 = ApiVersion("2014-05-15")
	ApiVersion20190510 = ApiVersion("2019-05-10")
)

const DefaultClientRetryCountSmall = 5

const Terraform = "HashiCorp-Terraform"

const Provider = "Terraform-Provider"

const Module = "Terraform-Module"

var providerVersion = "1.0.32"
var terraformVersion = strings.TrimSuffix(schema.Provider{}.TerraformVersion, "-dev")

type ApiVersion string

// The main version number that is being run at the moment.
var ProviderVersion = providerVersion
var TerraformVersion = strings.TrimSuffix(schema.Provider{}.TerraformVersion, "-dev")
var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe
var loadSdkfromRemoteMutex = sync.Mutex{}
var loadSdkEndpointMutex = sync.Mutex{}

// Client for AlibabacloudStackClient
func (c *Config) Client() (*AlibabacloudStackClient, error) {
	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.

	teaSdkConfig, err := c.getTeaDslSdkConfig(true)
	if err != nil {
		return nil, err
	}

	return &AlibabacloudStackClient{
		Config:                       c,
		teaSdkConfig:                 teaSdkConfig,
		Region:                       c.Region,
		RegionId:                     c.RegionId,
		AccessKey:                    c.AccessKey,
		SecretKey:                    c.SecretKey,
		Department:                   c.Department,
		ResourceGroup:                c.ResourceGroup,
		ResourceGroupId:              c.ResourceGroupId,
		Domain:                       c.Domain,
		OtsInstanceName:              c.OtsInstanceName,
		Conns:                        make(map[ServiceCode]*sdk.Client),
		tablestoreconnByInstanceName: make(map[string]*tablestore.TableStoreClient),
		Eagleeye:                     c.Eagleeye,
	}, nil
}

func (client *AlibabacloudStackClient) NewTeaSDkClient(productCode string, endpoint string) (*rpc.Client, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint).SetReadTimeout(client.Config.ClientReadTimeout * 1000) //单位毫秒
	conn, err := rpc.NewClient(&sdkConfig)
	for key, value := range client.defaultHeaders(productCode) {
		conn.Headers[key] = &value
	}
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AlibabacloudStackClient) WithProductSDKClient(popcode ServiceCode) (*sdk.Client, error) {
	endpoint := client.Config.Endpoints[popcode]
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] unable to initialize the %s client: endpoint or domain is not provided", string(popcode))
	}
	conn, err := sdk.NewClientWithOptions(client.Config.RegionId, client.getSdkConfig(), client.Config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", popcode, err)
	}

	conn.Domain = endpoint
	conn.SetReadTimeout(time.Duration(client.Config.ClientReadTimeout) * time.Hour)
	conn.SetConnectTimeout(time.Duration(client.Config.ClientConnectTimeout) * time.Hour)
	conn.SourceIp = client.Config.SourceIp
	conn.SecureTransport = client.Config.SecureTransport
	conn.AppendUserAgent(Terraform, TerraformVersion)
	conn.AppendUserAgent(Provider, ProviderVersion)
	conn.AppendUserAgent(Module, client.Config.ConfigurationSource)
	conn.SetHTTPSInsecure(client.Config.Insecure)
	if client.Config.Proxy != "" {
		conn.SetHttpsProxy(client.Config.Proxy)
		conn.SetHttpProxy(client.Config.Proxy)
	}
	return conn, nil
}

func (client *AlibabacloudStackClient) WithEcsClient(do func(*ecs.Client) (interface{}, error)) (interface{}, error) {
	if client.ecsconn == nil {
		conn, error := client.WithProductSDKClient(EcsCode)
		if error != nil {
			return nil, error
		}
		client.ecsconn = &ecs.Client{
			Client: *conn,
		}
	}
	return do(client.ecsconn)
}

func (client *AlibabacloudStackClient) WithAscmClient(do func(*sdk.Client) (interface{}, error)) (interface{}, error) {
	var err error
	if client.ascmconn == nil {
		client.ascmconn, err = client.WithProductSDKClient(ASCMCode)
		if err != nil {
			return nil, err
		}
	}
	return do(client.ascmconn)
}

func (client *AlibabacloudStackClient) WithElasticsearchClient(do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	if client.elasticsearchconn == nil {
		conn, error := client.WithProductSDKClient(ELASTICSEARCHCode)
		if error != nil {
			return nil, error
		}
		client.elasticsearchconn = &elasticsearch.Client{
			Client: *conn,
		}
	}

	return do(client.elasticsearchconn)
}

func (client *AlibabacloudStackClient) WithCloudApiClient(do func(*cloudapi.Client) (interface{}, error)) (interface{}, error) {
	if client.cloudapiconn == nil {
		conn, error := client.WithProductSDKClient(CLOUDAPICode)
		if error != nil {
			return nil, error
		}
		client.cloudapiconn = &cloudapi.Client{
			Client: *conn,
		}
	}
	return do(client.cloudapiconn)
}

func (client *AlibabacloudStackClient) WithEssClient(do func(*ess.Client) (interface{}, error)) (interface{}, error) {
	if client.essconn == nil {
		conn, error := client.WithProductSDKClient(ESSCode)
		if error != nil {
			return nil, error
		}
		client.essconn = &ess.Client{
			Client: *conn,
		}
	}
	return do(client.essconn)
}

func (client *AlibabacloudStackClient) WithOnsClient(do func(*ons.Client) (interface{}, error)) (interface{}, error) {
	if client.onsconn == nil {
		conn, error := client.WithProductSDKClient(ONSCode)
		if error != nil {
			return nil, error
		}
		client.onsconn = &ons.Client{
			Client: *conn,
		}
	}
	return do(client.onsconn)
}

func (client *AlibabacloudStackClient) WithRkvClient(do func(*r_kvstore.Client) (interface{}, error)) (interface{}, error) {
	if client.rkvconn == nil {
		conn, error := client.WithProductSDKClient(KVSTORECode)
		if error != nil {
			return nil, error
		}
		client.rkvconn = &r_kvstore.Client{
			Client: *conn,
		}
	}

	return do(client.rkvconn)
}

func (client *AlibabacloudStackClient) WithGpdbClient(do func(*gpdb.Client) (interface{}, error)) (interface{}, error) {
	if client.gpdbconn == nil {
		conn, error := client.WithProductSDKClient(GPDBCode)
		if error != nil {
			return nil, error
		}
		client.gpdbconn = &gpdb.Client{
			Client: *conn,
		}
	}

	return do(client.gpdbconn)
}
func (client *AlibabacloudStackClient) WithAdbClient(do func(*adb.Client) (interface{}, error)) (interface{}, error) {
	if client.adbconn == nil {
		conn, error := client.WithProductSDKClient(ADBCode)
		if error != nil {
			return nil, error
		}
		client.adbconn = &adb.Client{
			Client: *conn,
		}
	}

	return do(client.adbconn)
}
func (client *AlibabacloudStackClient) WithHbaseClient(do func(*hbase.Client) (interface{}, error)) (interface{}, error) {
	if client.hbaseconn == nil {
		conn, error := client.WithProductSDKClient(HBASECode)
		if error != nil {
			return nil, error
		}
		client.hbaseconn = &hbase.Client{
			Client: *conn,
		}
	}

	return do(client.hbaseconn)
}

func (client *AlibabacloudStackClient) WithVpcClient(do func(*vpc.Client) (interface{}, error)) (interface{}, error) {
	if client.vpcconn == nil {
		conn, error := client.WithProductSDKClient(VPCCode)
		if error != nil {
			return nil, error
		}
		client.vpcconn = &vpc.Client{
			Client: *conn,
		}
	}

	return do(client.vpcconn)
}

func (client *AlibabacloudStackClient) WithSlbClient(do func(*slb.Client) (interface{}, error)) (interface{}, error) {
	if client.slbconn == nil {
		conn, error := client.WithProductSDKClient(SLBCode)
		if error != nil {
			return nil, error
		}
		client.slbconn = &slb.Client{
			Client: *conn,
		}
	}

	return do(client.slbconn)
}
func (client *AlibabacloudStackClient) WithDdsClient(do func(*dds.Client) (interface{}, error)) (interface{}, error) {
	if client.ddsconn == nil {
		conn, error := client.WithProductSDKClient(DDSCode)
		if error != nil {
			return nil, error
		}
		client.ddsconn = &dds.Client{
			Client: *conn,
		}
	}

	return do(client.ddsconn)
}

func (client *AlibabacloudStackClient) WithOssNewClient(do func(*ecs.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ECS client if necessary
	if client.ecsconn == nil {
		endpoint := client.Config.Endpoints[OSSCode]
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the oss client: endpoint or domain is not provided for ecs service")
		}
		ecsconn, err := ecs.NewClientWithOptions(client.Config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.Config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ECS client: %#v", err)
		}

		ecsconn.Domain = endpoint
		ecsconn.AppendUserAgent(Terraform, TerraformVersion)
		ecsconn.AppendUserAgent(Provider, ProviderVersion)
		ecsconn.AppendUserAgent(Module, client.Config.ConfigurationSource)
		ecsconn.SetHTTPSInsecure(client.Config.Insecure)
		if client.Config.Proxy != "" {
			ecsconn.SetHttpsProxy(client.Config.Proxy)
			ecsconn.SetHttpProxy(client.Config.Proxy)
		}
		client.ecsconn = ecsconn
	}

	return do(client.ecsconn)
}

func (client *AlibabacloudStackClient) getSdkConfig() *sdk.Config {
	log.Printf("Protocol is set to %s", client.Config.Protocol)
	return sdk.NewConfig().
		WithMaxRetryTime(DefaultClientRetryCountSmall).
		WithTimeout(time.Duration(30) * time.Second).
		WithEnableAsync(true).
		WithGoRoutinePoolSize(100).
		WithMaxTaskQueueSize(10000).
		WithDebug(false).
		WithHttpTransport(client.getTransport()).
		WithScheme(strings.ToLower(client.Config.Protocol))
}

func (client *AlibabacloudStackClient) getTransport() *http.Transport {
	handshakeTimeout, err := strconv.Atoi(os.Getenv("TLSHandshakeTimeout"))
	if err != nil {
		handshakeTimeout = 120
	}
	transport := &http.Transport{}
	transport.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second

	return transport
}
func (client *AlibabacloudStackClient) AccountId() (string, error) {
	client.accountIdMutex.Lock()
	defer client.accountIdMutex.Unlock()

	if client.accountId == "" {
		log.Printf("[DEBUG] account_id not provided, attempting to retrieve it automatically...")
		identity, err := client.GetCallerIdentity()
		if err != nil {
			return "", err
		}
		if identity == "" {
			return "", fmt.Errorf("caller identity doesn't contain any AccountId")
		}
		client.accountId = identity
	}
	return client.accountId, nil
}

func (client *AlibabacloudStackClient) RoleIds() (int, error) {
	client.roleIdMutex.Lock()
	defer client.roleIdMutex.Unlock()

	if client.roleId == 0 {
		log.Printf("[DEBUG] role_ids not provided, attempting to retrieve it automatically...")
		roleId, err := client.GetCallerDefaultRole()
		if err != nil {
			return 0, err
		}
		if roleId == 0 {
			return 0, fmt.Errorf("caller identity doesn't contain default RoleId")
		}
		client.roleId = roleId
	}
	return client.roleId, nil
}

func (client *AlibabacloudStackClient) getHttpProxy() (proxy *url.URL, err error) {
	if client.Config.Protocol == "HTTPS" {
		if rawurl := os.Getenv("HTTPS_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("https_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	} else {
		if rawurl := os.Getenv("HTTP_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("http_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	}
	return proxy, err
}

func (client *AlibabacloudStackClient) skipProxy(endpoint string) (bool, error) {
	var urls []string
	if rawurl := os.Getenv("NO_PROXY"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	} else if rawurl := os.Getenv("no_proxy"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	}
	for _, value := range urls {
		if strings.HasPrefix(value, "*") {
			value = fmt.Sprintf(".%s", value)
		}
		noProxyReg, err := regexp.Compile(value)
		if err != nil {
			return false, err
		}
		if noProxyReg.MatchString(endpoint) {
			return true, nil
		}
	}
	return false, nil
}
func (client *AlibabacloudStackClient) WithKmsClient(do func(*kms.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the KMS client if necessary
	if client.kmsconn == nil {
		conn, error := client.WithProductSDKClient(KmsCode)
		if error != nil {
			return nil, error
		}
		client.kmsconn = &kms.Client{
			Client: *conn,
		}
	}
	return do(client.kmsconn)
}

func (client *AlibabacloudStackClient) GetCallerInfo() (*responses.BaseResponse, error) {

	endpoint := client.Config.Endpoints[ASCMCode]
	if endpoint == "" {
		return nil, fmt.Errorf("unable to initialize the ascm client: endpoint or domain is not provided for ascm service")
	}
	ascmClient, err := sdk.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the ascm client: %#v", err)
	}

	ascmClient.AppendUserAgent(Terraform, TerraformVersion)
	ascmClient.AppendUserAgent(Provider, ProviderVersion)
	ascmClient.AppendUserAgent(Module, client.Config.ConfigurationSource)
	ascmClient.SetHTTPSInsecure(client.Config.Insecure)
	ascmClient.Domain = endpoint
	if client.Config.Proxy != "" {
		ascmClient.SetHttpProxy(client.Config.Proxy)
	}
	if client.Config.Department == "" || client.Config.ResourceGroup == "" {
		return nil, fmt.Errorf("unable to initialize the ascm client: department or resource_group is not provided")
	}
	request := requests.NewCommonRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "GET"         // Set request method
	request.Product = "ascm"       // Specify product
	request.Domain = endpoint      // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2019-05-10" // Specify product version
	request.ApiName = "GetUserInfo"
	request.QueryParams = map[string]string{
		// 		"AccessKeySecret":  client.Config.SecretKey,
		// 		"SecurityToken":    client.Config.SecurityToken,
		// 		"Product":          "ascm",
		// 		"Department":       client.Config.Department,
		// 		"ResourceGroup":    client.Config.ResourceGroup,
		// 		"RegionId":         client.RegionId,
		// 		"Action":           "GetAllNavigationInfo",
		// 		"Version":          "2019-05-10",
		"SignatureVersion": "1.0",
	}
	resp := responses.BaseResponse{}
	request.TransToAcsRequest()
	err = ascmClient.DoAction(request, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
func (client *AlibabacloudStackClient) GetCallerIdentity() (string, error) {

	resp, err := client.GetCallerInfo()
	response := &AccountId{}
	err = json.Unmarshal(resp.GetHttpContentBytes(), response)
	ownerId := response.Data.PrimaryKey

	if ownerId == "" {
		return "", fmt.Errorf("ownerId not found")
	}
	return ownerId, err
}

func (client *AlibabacloudStackClient) GetCallerDefaultRole() (int, error) {

	resp, err := client.GetCallerInfo()
	if err != nil {
		return 1, err
	}
	response := &RoleId{}
	err = json.Unmarshal(resp.GetHttpContentBytes(), response)
	roleId := response.Data.DefaultRole.Id

	if roleId == 0 {
		return 0, fmt.Errorf("default roleId not found")
	}
	return roleId, err
}

type AccountId struct {
	Data struct {
		PrimaryKey string `json:"primaryKey"`
	} `json:"data"`
}

type RoleId struct {
	Data struct {
		DefaultRole struct {
			Id int `json:"id"`
		} `json:"defaultRole"`
	} `json:"data"`
}

func (client *AlibabacloudStackClient) WithBssopenapiClient(do func(*bssopenapi.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the bssopenapi client if necessary
	if client.bssopenapiconn == nil {
		conn, error := client.WithProductSDKClient(BssDataCode)
		if error != nil {
			return nil, error
		}
		client.bssopenapiconn = &bssopenapi.Client{
			Client: *conn,
		}
	}

	return do(client.bssopenapiconn)
}

func (client *AlibabacloudStackClient) NewNasClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("nas", client.Config.Endpoints[NasCode])
}

func (client *AlibabacloudStackClient) WithOssClientPutObject(do func(*oss.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the OSS client if necessary
	if client.ossconn == nil {
		schma := "http"
		endpoint := client.Config.Endpoints[OssDataCode]
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the oss client: endpoint or domain is not provided for OSS service")
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
		}

		clientOptions := []oss.ClientOption{oss.UserAgent(client.getUserAgent()),
			oss.SecurityToken(client.Config.SecurityToken)}
		if client.Config.Proxy != "" {
			clientOptions = append(clientOptions, oss.Proxy(client.Config.Proxy))
		}

		clientOptions = append(clientOptions, oss.UseCname(false))

		ossconn, err := oss.New(endpoint, client.Config.AccessKey, client.Config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the OSS client: %#v", err)
		}

		client.ossconn = ossconn
	}

	return do(client.ossconn)
}

func (client *AlibabacloudStackClient) WithOssClient(do func(*oss.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the OSS client if necessary
	if client.ossconn == nil {
		schma := "http"
		endpoint := client.Config.Endpoints[OSSCode]
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the oss client: endpoint or domain is not provided for OSS service")
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
		}

		clientOptions := []oss.ClientOption{oss.UserAgent(client.getUserAgent()),
			oss.SecurityToken(client.Config.SecurityToken)}
		if client.Config.Proxy != "" {
			clientOptions = append(clientOptions, oss.Proxy(client.Config.Proxy))
		}

		clientOptions = append(clientOptions, oss.UseCname(false))

		ossconn, err := oss.New(endpoint, client.Config.AccessKey, client.Config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the OSS client: %#v", err)
		}

		client.ossconn = ossconn
	}

	return do(client.ossconn)
}

func (client *AlibabacloudStackClient) WithRamClient(do func(*ram.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RAM client if necessary
	if client.ramconn == nil {
		conn, error := client.WithProductSDKClient(RAMCode)
		if error != nil {
			return nil, error
		}
		client.ramconn = &ram.Client{
			Client: *conn,
		}
	}

	return do(client.ramconn)
}

func (client *AlibabacloudStackClient) WithRdsClient(do func(*rds.Client) (interface{}, error)) (interface{}, error) {
	if client.rdsconn == nil {
		conn, error := client.WithProductSDKClient(RDSCode)
		if error != nil {
			return nil, error
		}
		client.rdsconn = &rds.Client{
			Client: *conn,
		}
	}

	return do(client.rdsconn)
}

func (client *AlibabacloudStackClient) WithCdnClient(do func(*cdn.Client) (interface{}, error)) (interface{}, error) {
	if client.cdnconn == nil {
		conn, error := client.WithProductSDKClient(CDNCode)
		if error != nil {
			return nil, error
		}
		client.cdnconn = &cdn.Client{
			Client: *conn,
		}
	}

	return do(client.cdnconn)
}
func (client *AlibabacloudStackClient) getUserAgent() string {
	return fmt.Sprintf("%s/%s %s/%s %s/%s", Terraform, TerraformVersion, Provider, ProviderVersion, Module, client.Config.ConfigurationSource)
}

func (client *AlibabacloudStackClient) WithCsClient(do func(*cs.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CS client if necessary
	if client.csconn == nil {
		csconn := cs.NewClientForAussumeRole(client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
		csconn.SetUserAgent(client.getUserAgent())
		endpoint := client.Config.Endpoints[CONTAINCode]
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the cs client: endpoint or domain is not provided for cs service")
		}
		if endpoint != "" {
			if strings.ToLower(client.Config.Protocol) == "https" {
				endpoint = fmt.Sprintf("https://%s", endpoint)
			} else {
				endpoint = fmt.Sprintf("http://%s", endpoint)
			}
			csconn.SetEndpoint(endpoint)
		}
		if client.Config.Proxy != "" {
			os.Setenv("http_proxy", client.Config.Proxy)
		}
		client.csconn = csconn
	}

	return do(client.csconn)
}

func (client *AlibabacloudStackClient) getHttpProxyUrl() *url.URL {
	for _, v := range []string{"HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy"} {
		value := strings.Trim(os.Getenv(v), " ")
		if value != "" {
			if !regexp.MustCompile(`^http(s)?://`).MatchString(value) {
				value = fmt.Sprintf("https://%s", value)
			}
			proxyUrl, err := url.Parse(value)
			if err == nil {
				return proxyUrl
			}
			break
		}
	}
	return nil
}

func (client *AlibabacloudStackClient) WithOssBucketByName(bucketName string, do func(*oss.Bucket) (interface{}, error)) (interface{}, error) {
	return client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		bucket, err := client.ossconn.Bucket(bucketName)

		if err != nil {
			return nil, fmt.Errorf("unable to get the bucket %s: %#v", bucketName, err)
		}
		return do(bucket)
	})
}

func (client *AlibabacloudStackClient) WithSlsClient(do func(*slsPop.Client) (interface{}, error)) (interface{}, error) {
	if client.logpopconn == nil {
		conn, error := client.WithProductSDKClient(SLSCode)
		if error != nil {
			return nil, error
		}
		client.logpopconn = &slsPop.Client{
			Client: *conn,
		}
	}

	return do(client.logpopconn)
}

func (client *AlibabacloudStackClient) WithSlsDataClient(do func(*sls.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the LOG client if necessary
	if client.logconn == nil {
		endpoint := client.Config.Endpoints[SlSDataCode]
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the log client: endpoint or domain is not provided for log service")
		}
		if client.Config.Proxy != "" {
			os.Setenv("http_proxy", client.Config.Proxy)
		}
		client.logconn = &sls.Client{
			// AccessKeyID:     client.Config.OrganizationAccessKey,
			// AccessKeySecret: client.Config.OrganizationSecretKey,
			AccessKeyID:     client.Config.AccessKey,
			AccessKeySecret: client.Config.SecretKey,
			Endpoint:        client.Config.Endpoints[SlSDataCode],
			SecurityToken:   client.Config.SecurityToken,
			UserAgent:       client.getUserAgent(),
		}
	}

	return do(client.logconn)
}

func (client *AlibabacloudStackClient) WithAlikafkaClient(do func(*alikafka.Client) (interface{}, error)) (interface{}, error) {
	if client.alikafkaconn == nil {
		conn, error := client.WithProductSDKClient(ALIKAFKACode)
		if error != nil {
			return nil, error
		}
		client.alikafkaconn = &alikafka.Client{
			Client: *conn,
		}
	}

	return do(client.alikafkaconn)
}

func (client *AlibabacloudStackClient) WithEdasClient(do func(*edas.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the edas client if necessary
	if client.edasconn == nil {
		endpoint := client.Config.Endpoints[EDASCode]
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the Edas client: endpoint or domain is not provided for Edas service")
		}
		// edasconn, err := edas.NewClientWithOptions(client.Config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.Config.getAuthCredential(true))
		var edasconn *edas.Client
		var err error
		if client.Config.OrganizationAccessKey != "" && client.Config.OrganizationSecretKey != "" {
			edasconn, err = edas.NewClientWithAccessKey(client.Config.RegionId, client.Config.OrganizationAccessKey, client.Config.OrganizationSecretKey)
		} else {
			edasconn, err = edas.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
		}
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Edas client: %#v", err)
		}
		edasconn.SetReadTimeout(time.Duration(client.Config.ClientReadTimeout) * time.Millisecond)
		edasconn.SetConnectTimeout(time.Duration(client.Config.ClientConnectTimeout) * time.Millisecond)
		edasconn.SourceIp = client.Config.SourceIp
		edasconn.SecureTransport = client.Config.SecureTransport
		edasconn.Domain = endpoint
		edasconn.AppendUserAgent(Terraform, terraformVersion)
		edasconn.AppendUserAgent(Provider, providerVersion)
		edasconn.AppendUserAgent(Module, client.Config.ConfigurationSource)
		if client.Config.Proxy != "" {
			edasconn.SetHttpsProxy(client.Config.Proxy)
			edasconn.SetHttpProxy(client.Config.Proxy)
		}
		client.edasconn = edasconn
	}

	return do(client.edasconn)
}

func (client *AlibabacloudStackClient) WithCrEEClient(do func(*cr_ee.Client) (interface{}, error)) (interface{}, error) {
	if client.creeconn == nil {
		conn, error := client.WithProductSDKClient(CRCode)
		if error != nil {
			return nil, error
		}
		client.creeconn = &cr_ee.Client{
			Client: *conn,
		}
	}

	return do(client.creeconn)
}

func (client *AlibabacloudStackClient) WithDnsClient(do func(*alidns.Client) (interface{}, error)) (interface{}, error) {
	if client.dnsconn == nil {
		conn, error := client.WithProductSDKClient(DNSCode)
		if error != nil {
			return nil, error
		}
		client.dnsconn = &alidns.Client{
			Client: *conn,
		}
	}

	return do(client.dnsconn)
}
func (client *AlibabacloudStackClient) WithCmsClient(do func(*cms.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CMS client if necessary
	if client.cmsconn == nil {
		conn, error := client.WithProductSDKClient(CMSCode)
		if error != nil {
			return nil, error
		}
		client.cmsconn = &cms.Client{
			Client: *conn,
		}
	}

	return do(client.cmsconn)
}

func (client *AlibabacloudStackClient) NewHitsdbClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("hitsdb", client.Config.Endpoints[HitsdbCode])
}

func (client *AlibabacloudStackClient) NewOdpsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("odps", client.Config.Endpoints[ASCMCode])
}
func (client *AlibabacloudStackClient) NewKmsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("kms", client.Config.Endpoints[KmsCode])
}

func (client *AlibabacloudStackClient) NewAscmClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("ascm", client.Config.Endpoints[ASCMCode])
}
func (client *AlibabacloudStackClient) NewCloudApiClient() (*rpc.Client, error) {
	//sdkConfig.SetEndpoint(endpoint).SetReadTimeout(60000)
	return client.NewTeaSDkClient("apigateway", client.Config.Endpoints[CLOUDAPICode])
}

func (client *AlibabacloudStackClient) NewAdsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("ads", client.Config.Endpoints[ADBCode])
}

func (client *AlibabacloudStackClient) NewCmsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("cms", client.Config.Endpoints[CMSCode])
}

func (client *AlibabacloudStackClient) WithTableStoreClient(instanceName string, do func(*tablestore.TableStoreClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the TABLESTORE client if necessary
	tableStoreClient, ok := client.tablestoreconnByInstanceName[instanceName]
	if !ok {
		endpoint := client.Config.Endpoints[OTSCode]
		if endpoint == "" {
			return nil, fmt.Errorf("[ERROR] missing the product Ots endpoint.")
		}
		// if !strings.HasPrefix(endpoint, "https") && !strings.HasPrefix(endpoint, "http") {
		// 	endpoint = fmt.Sprintf("https://%s", endpoint)
		// }
		// endpoint := "http://test1111.cn-wulan-env212-d01.ots-internal.inter.env212.shuguang.com"
		tableStoreClient = tablestore.NewClientWithConfig(endpoint, instanceName, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken, tablestore.NewDefaultTableStoreConfig())
		client.tablestoreconnByInstanceName[instanceName] = tableStoreClient
	}

	return do(tableStoreClient)
}
func (client *AlibabacloudStackClient) WithOtsClient(do func(*ots.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the OTS client if necessary
	if client.otsconn == nil {
		conn, error := client.WithProductSDKClient(OTSCode)
		if error != nil {
			return nil, error
		}
		client.otsconn = &ots.Client{
			Client: *conn,
		}
	}

	return do(client.otsconn)
}
func (client *AlibabacloudStackClient) WithDataHubClient(do func(api datahub.DataHubApi) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the DataHub client if necessary
	if client.dhconn == nil {
		endpoint := client.Config.Endpoints[DatahubCode]
		if endpoint == "" {
			return nil, fmt.Errorf("[ERROR] missing the product Ots endpoint.")
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}

		account := datahub.NewStsCredential(client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
		config := &datahub.Config{
			UserAgent: client.getUserAgent(),
		}

		client.dhconn = datahub.NewClientWithConfig(endpoint, config, account)
	}

	return do(client.dhconn)
}
func (client *AlibabacloudStackClient) NewVpcClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("vpc", client.Config.Endpoints[VPCCode])
}
func (client *AlibabacloudStackClient) NewEcsClient() (*rpc.Client, error) {
	//sdkConfig.SetEndpoint(endpoint).SetReadTimeout(60000)
	return client.NewTeaSDkClient("ecs", client.Config.Endpoints[EcsCode])
}
func (client *AlibabacloudStackClient) NewElasticsearchClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("elasticsearch", client.Config.Endpoints[ELASTICSEARCHCode])
}

func (client *AlibabacloudStackClient) NewRosClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("ros", client.Config.Endpoints[RosCode])
}

func (client *AlibabacloudStackClient) NewRdsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("rds", client.Config.Endpoints[RDSCode])
}

func (client *AlibabacloudStackClient) NewRoaCsClient() (*roaCS.Client, error) {
	productCode := "ros"
	endpoint := client.Config.Endpoints[RosCode]
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	// Initialize the CS client if necessary
	roaCSConn, err := roaCS.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(client.Config.AccessKey),
		AccessKeySecret: tea.String(client.Config.SecretKey),
		SecurityToken:   tea.String(client.Config.SecurityToken),
		RegionId:        tea.String(client.Config.RegionId),
		UserAgent:       tea.String(client.getUserAgent()),
		Endpoint:        tea.String(endpoint),
		ReadTimeout:     tea.Int(client.Config.ClientReadTimeout),
		ConnectTimeout:  tea.Int(client.Config.ClientConnectTimeout),
	})
	roaCSConn.Headers = map[string]*string{
		"x-acs-organizationid":  &client.Config.Department,
		"x-acs-resourcegroupid": &client.Config.ResourceGroup,
	}
	if err != nil {
		return nil, err
	}

	return roaCSConn, nil
}

func (client *AlibabacloudStackClient) NewDtsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("dts", client.Config.Endpoints[DTSCode])
}

func (client *AlibabacloudStackClient) NewDmsenterpriseClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("dmsenterprise", client.Config.Endpoints[DmsEnterpriseCode])
}

func (client *AlibabacloudStackClient) NewHbaseClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("hbase", client.Config.Endpoints[HBASECode])
}

func (client *AlibabacloudStackClient) WithDrdsClient(do func(*drds.Client) (interface{}, error)) (interface{}, error) {
	if client.drdsconn == nil {
		conn, error := client.WithProductSDKClient(DRDSCode)
		if error != nil {
			return nil, error
		}
		client.drdsconn = &drds.Client{
			Client: *conn,
		}
	}

	return do(client.drdsconn)
}
func (client *AlibabacloudStackClient) NewGpdbClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("gpdb", client.Config.Endpoints[GPDBCode])
}

func (client *AlibabacloudStackClient) NewQuickbiClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("quickbi", client.Config.Endpoints[QuickbiCode])
}
func (client *AlibabacloudStackClient) NewCsbClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("csb", client.Config.Endpoints[CSBCode])
}
func (client *AlibabacloudStackClient) NewGdbClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("gdb", client.Config.Endpoints[GDBCode])
}

func (client *AlibabacloudStackClient) NewDataworkspublicClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("dataworkspublic", client.Config.Endpoints[DataworkspublicCode])
}

func (client *AlibabacloudStackClient) NewDataworksPrivateClient() (*rpc.Client, error) {
	endpoint := client.Config.Endpoints[DataworkspublicCode]
	index := strings.Index(endpoint, ".")
	privateEndpoint := "dataworks" + endpoint[index:]
	return client.NewTeaSDkClient("dataworks-private-cloud", privateEndpoint)
}

func (client *AlibabacloudStackClient) NewDbsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("dbs", client.Config.Endpoints[DDSCode])
}
func (client *AlibabacloudStackClient) NewArmsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("arms", client.Config.Endpoints[ARMSCode])
}

func (client *AlibabacloudStackClient) NewOosClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("oos", client.Config.Endpoints[OosCode])
}
func (client *AlibabacloudStackClient) NewCloudfwClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("cloudfw", client.Config.Endpoints[WafOpenapiCode])
}

func (client *AlibabacloudStackClient) defaultHeaders(popcode string) map[string]string {
	return map[string]string{
		"RegionId":              client.RegionId, //	ASAPI
		"x-acs-organizationid":  client.Department,
		"x-acs-resourcegroupid": client.ResourceGroup,
		"x-acs-regionid":        client.RegionId,
		"x-acs-request-version": "v1",
		"x-acs-asapi-product":   popcode,
		"x-ascm-product-name":   popcode,
		"EagleEye-TraceId":      client.Eagleeye.GetTraceId(),
		"EagleEye-RpcId":        client.Eagleeye.GetRpcId(),
		//"x-acs-caller-sdk-source": "Terraform"
		//"x-acs-asapi-gateway-version": "3.0"  这个是指定走ASAPI的v3网关流程，目前在维护的是v4，默认会走v4，指定了走v3。不建议走v3，除非有不兼容的地方必须走
	}
}

func (client *AlibabacloudStackClient) defaultQueryParams() map[string]string {
	return map[string]string{
		"RegionId":       client.RegionId,
		"Department":     client.Department,
		"OrganizationId": client.Department,
		"ResourceGroup":  client.ResourceGroup,
	}
}

func (client *AlibabacloudStackClient) NewCommonRequest(method string, popcode string, version string, apiname string, pathpattern string) *requests.CommonRequest {
	request := requests.NewCommonRequest()

	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Method = method
	request.RegionId = client.RegionId
	request.Headers = client.defaultHeaders(popcode)
	request.QueryParams = client.defaultQueryParams()
	request.Product = popcode
	request.Version = version
	request.ApiName = apiname
	if pathpattern != "" {
		request.PathPattern = pathpattern
	}

	return request
}

func (client *AlibabacloudStackClient) InitRpcRequest(request requests.RpcRequest) {
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = client.defaultHeaders(request.GetProduct())
	request.QueryParams = client.defaultQueryParams()
}

func (client *AlibabacloudStackClient) InitRoaRequest(request requests.RoaRequest) {
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = client.defaultHeaders(request.GetProduct())
	request.QueryParams = client.defaultQueryParams()
}

func (client *AlibabacloudStackClient) DoTeaRequest(method string, popcode string, version string, apiname string, pathpattern string, query map[string]interface{}, body map[string]interface{}) (_result map[string]interface{}, _err error) {
	ServiceCodeStr := strings.ReplaceAll(strings.ToUpper(popcode), "-", "_")
	endpoint := client.Config.Endpoints[ServiceCode(ServiceCodeStr)]
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", popcode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint).SetReadTimeout(client.Config.ClientReadTimeout * 1000) //单位毫秒

	headers := make(map[string]*string)
	for key, value := range client.defaultHeaders(popcode) {
		v := value
		headers[key] = &v
	}

	if query == nil {
		query = make(map[string]interface{})
	}
	for key, value := range client.defaultQueryParams() {
		if _, exist := query[key]; !exist {
			query[key] = value
		}
	}
	query["Product"] = popcode

	var protocol string
	if strings.ToLower(client.Config.Protocol) == "https" {
		protocol = "https"
	} else {
		protocol = "http"
	}
	authType := "AK"

	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	if client.Config.Proxy != "" {
		runtime.HttpProxy = &client.Config.Proxy
		runtime.HttpsProxy = &client.Config.Proxy
	}
	runtime.SetAutoretry(false) // 使用ASAPI时，Tea包不能重试，他会修改endpoint

	var response map[string]interface{}
	wait := IncrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {

		// conn 必须在每次请求前初始化，因为DoRequest会修改conn的内容，会导致下次conn的配置失效
		conn, err := rpc.NewClient(&sdkConfig)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("unable to initialize the %s client: %#v", popcode, err))
		}
		conn.Headers = headers

		response, err = conn.DoRequest(&apiname, &protocol, &method, &version, &authType, query, body, &runtime)
		log.Printf(" ================================ %s ======================================\n query %v \n request %v \n response: %v", apiname, query, body, response)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			errmsg := errmsgs.GetAsapiErrorMessage(response)
			if errmsg != "" {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "popcode", apiname, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			wait()
			return resource.RetryableError(err)
		}
		return nil
	})
	return response, err
}

func (client *AlibabacloudStackClient) getConnectClient(popcode ServiceCode) (*sdk.Client, error) {
	var conn *sdk.Client
	var exists bool
	if conn, exists = client.Conns[popcode]; !exists {
		c, err := client.WithProductSDKClient(popcode)
		if err != nil {
			return nil, err
		}
		client.Conns[popcode] = c
		conn = c
	}
	return conn, nil
}

func (client *AlibabacloudStackClient) ProcessCommonRequest(request *requests.CommonRequest) (*responses.CommonResponse, error) {
	popcode := ServiceCode(strings.ToUpper(request.Product))

	conn, err := client.getConnectClient(popcode)
	if err != nil {
		return nil, err
	}

	//request.Domain = conn.Domain
	if strings.HasPrefix(conn.Domain, "internal.asapi.") || strings.HasPrefix(conn.Domain, "public.asapi.") {
		// asapi兼容逻辑
		// # asapi 使用common SDK时不能拼接pathpattern，否则会报错
		if request.PathPattern != "" {
			var r []string = strings.SplitN(conn.Domain, "/", 2)
			request.Domain = r[0]
			request.PathPattern = "/asapi/v3"
		}
		if len(request.Content) > 0 {
			request.QueryParams["x-acs-body"] = string(request.Content)
			request.SetContent([]byte("{}"))
		}
		request.Method = "POST"
	}

	response, err := conn.ProcessCommonRequest(request)
	return response, err
}

func IncrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	retryCount := 1
	return func() {
		var waitTime time.Duration
		if retryCount == 1 {
			waitTime = firstDuration
		} else if retryCount > 1 {
			waitTime += increaseDuration
		}
		time.Sleep(waitTime)
		retryCount++
	}
}

func GetResourceData(d *schema.ResourceData, keys ...string) interface{} {
	v, _ := GetResourceDataOk(d, keys...)
	return v
}

func GetResourceDataOk(d *schema.ResourceData, keys ...string) (interface{}, bool) {
	if d.IsNewResource() {
		for _, key := range keys {
			value, ok := d.GetOk(key)
			if ok {
				return value, true
			}
		}
	} else {
		for _, key := range keys {
			if d.HasChange(key) {
				return d.Get(key), true
			}
		}
	}
	return d.GetOk(keys[0])
}

func SetResourceData(d *schema.ResourceData, value interface{}, keys ...string) error {
	for _, key := range keys {
		if err := d.Set(key, value); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}
