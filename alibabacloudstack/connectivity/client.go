package connectivity

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/PaesslerAG/jsonpath"
	"github.com/google/uuid"

	roaCS "github.com/alibabacloud-go/cs-20151215/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	roa "github.com/alibabacloud-go/tea-roa/client"
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
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/fc-go-sdk"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"

	"sync"

	rpcutil "github.com/alibabacloud-go/tea-rpc-utils/service"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"fmt"
	"net/http"
	"net/url"
	"os"
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
	ResourceGroupId              string
	Config                       *Config
	teaRpcSdkConfig              rpc.Config
	teaRoaSdkConfig              roa.Config
	accountId                    string
	roleId                       int
	Conns                        map[ServiceCode]*sdk.Client
	connsMu                      sync.Mutex
	ascmconn                     *sdk.Client
	ecsconn                      *ecs.Client
	accountIdMutex               sync.RWMutex
	roleIdMutex                  sync.RWMutex
	OssEndpointOnce              sync.Once
	vpcconn                      *vpc.Client
	bastionhostprivateconn       *yundun_bastionhost.Client
	slbconn                      *slb.Client
	polarDBconn                  *polardb.Client
	cdnconn                      *cdn.Client
	kmsconn                      *kms.Client
	bssopenapiconn               *bssopenapi.Client
	rdsconn                      *rds.Client
	ramconn                      *ram.Client
	gpdbconn                     *gpdb.Client
	drdsconn                     *drds.Client
	elasticsearchconn            *elasticsearch.Client
	hbaseconn                    *hbase.Client
	adbconn                      *adb.Client
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

type ApiVersion string

// The main version number that is being run at the moment.
var ProviderVersion = "0.0.99"
var TerraformVersion = strings.TrimSuffix(schema.Provider{}.TerraformVersion, "-dev")
var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe
var loadSdkfromRemoteMutex = sync.Mutex{}
var loadSdkEndpointMutex = sync.Mutex{}

// Client for AlibabacloudStackClient
func (c *Config) Client() (*AlibabacloudStackClient, error) {
	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.

	teaRpcSdkConfig, err := c.getTeaRpcDslSdkConfig(true)
	if err != nil {
		return nil, err
	}
	teaRoaSdkConfig, err := c.getTeaRoaDslSdkConfig(true)
	if err != nil {
		return nil, err
	}

	return &AlibabacloudStackClient{
		Config:                       c,
		teaRpcSdkConfig:              teaRpcSdkConfig,
		teaRoaSdkConfig:              teaRoaSdkConfig,
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
	sdkConfig := client.teaRpcSdkConfig
	sdkConfig.SetEndpoint(endpoint).SetReadTimeout(client.Config.ClientReadTimeout * 1000) // Unit: milliseconds
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

	ramSupported := true
	if popcode == STSCode {
		ramSupported = false // TODO: STS does not support NewRamRoleArnWithPolicyCredential, need to investigate
	}
	conn, err := sdk.NewClientWithOptions(client.Config.RegionId, client.getSdkConfig(), client.Config.getAuthCredential(true, ramSupported))
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

func (client *AlibabacloudStackClient) WithElasticsearchClient(do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	if client.elasticsearchconn == nil {
		conn, error := client.WithProductSDKClient(ElasticsearchK8sCode)
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
	conn, err := client.getConnectClient("ESS")
	if err != nil {
		return nil, err
	}
	essconn := &ess.Client{
		Client: *conn,
	}
	return retryDo(func() (interface{}, error) {
		return do(essconn)
	})
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

func (client *AlibabacloudStackClient) WithBastionhostClient(do func(*yundun_bastionhost.Client) (interface{}, error)) (interface{}, error) {
	if client.bastionhostprivateconn == nil {
		conn, error := client.WithProductSDKClient(BastionHostCode)
		if error != nil {
			return nil, error
		}
		client.bastionhostprivateconn = &yundun_bastionhost.Client{
			Client: *conn,
		}
	}

	return do(client.bastionhostprivateconn)
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

func (client *AlibabacloudStackClient) WithKmsClient(do func(*kms.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the KMS client if necessary
	if client.kmsconn == nil {
		endpoint := client.Config.Endpoints[KmsCode]
		if endpoint == "" {
			return nil, fmt.Errorf("[ERROR] unable to initialize the KMS client: endpoint or domain is not provided")
		}

		// Create KMS client with AK/SK or STS
		tmpKmsConn, err := kms.NewClientWithOptions(client.Config.RegionId, client.getSdkConfig(), client.Config.getAuthCredential(true, true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the KMS client: %#v", err)
		}
		// Configure KMS client before assignment to avoid concurrent map writes
		tmpKmsConn.Domain = endpoint
		tmpKmsConn.SetReadTimeout(time.Duration(client.Config.ClientReadTimeout) * time.Hour)
		tmpKmsConn.SetConnectTimeout(time.Duration(client.Config.ClientConnectTimeout) * time.Hour)
		tmpKmsConn.SourceIp = client.Config.SourceIp
		tmpKmsConn.SecureTransport = client.Config.SecureTransport
		tmpKmsConn.AppendUserAgent(Terraform, TerraformVersion)
		tmpKmsConn.AppendUserAgent(Provider, ProviderVersion)
		tmpKmsConn.AppendUserAgent(Module, client.Config.ConfigurationSource)
		tmpKmsConn.SetHTTPSInsecure(client.Config.Insecure)
		if client.Config.Proxy != "" {
			tmpKmsConn.SetHttpsProxy(client.Config.Proxy)
			tmpKmsConn.SetHttpProxy(client.Config.Proxy)
		}
		client.kmsconn = tmpKmsConn
	}
	return do(client.kmsconn)
}

func (client *AlibabacloudStackClient) GetCallerInfo() (map[string]interface{}, error) {

	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "GetUserInfo", "/ascm/auth/user/getUserInfo", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := jsonpath.Get("$.data", response)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

func (client *AlibabacloudStackClient) GetCallerIdentity() (string, error) {

	resp, err := client.GetCallerInfo()
	if err != nil {
		return "", err
	}
	ownerId := resp["primaryKey"].(string)

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
	defaultRole, ok := resp["defaultRole"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("defaultRole is not a map")
	}

	var role string
	switch v := defaultRole["id"].(type) {
	case string:
		role = v
	case json.Number:
		role = v.String()
	default:
		return 0, fmt.Errorf("unexpected type for role id: %T", v)
	}

	roleId, err := strconv.Atoi(role)
	if err != nil {
		return 0, err
	}

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

func (client *AlibabacloudStackClient) GetUserAgent() string {
	return fmt.Sprintf("%s/%s %s/%s %s/%s", Terraform, TerraformVersion, Provider, ProviderVersion, Module, client.Config.ConfigurationSource)
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
			// FIXME: Modifying environment variables may pose risks
			os.Setenv("http_proxy", client.Config.Proxy)
			os.Setenv("https_proxy", client.Config.Proxy)
		}
		client.logconn = &sls.Client{
			// AccessKeyID:     client.Config.OrganizationAccessKey,
			// AccessKeySecret: client.Config.OrganizationSecretKey,
			AccessKeyID:     client.Config.AccessKey,
			AccessKeySecret: client.Config.SecretKey,
			Endpoint:        client.Config.Endpoints[SlSDataCode],
			SecurityToken:   client.Config.SecurityToken,
			UserAgent:       client.GetUserAgent(),
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
		edasconn, err := edas.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Edas client: %#v", err)
		}
		edasconn.SetReadTimeout(time.Duration(client.Config.ClientReadTimeout) * time.Millisecond)
		edasconn.SetConnectTimeout(time.Duration(client.Config.ClientConnectTimeout) * time.Millisecond)
		edasconn.SourceIp = client.Config.SourceIp
		edasconn.SecureTransport = client.Config.SecureTransport
		edasconn.Domain = endpoint
		edasconn.AppendUserAgent(Terraform, TerraformVersion)
		edasconn.AppendUserAgent(Provider, ProviderVersion)
		edasconn.AppendUserAgent(Module, client.Config.ConfigurationSource)
		if client.Config.Proxy != "" {
			edasconn.SetHttpsProxy(client.Config.Proxy)
			edasconn.SetHttpProxy(client.Config.Proxy)
		}
		client.edasconn = edasconn
	}

	return do(client.edasconn)
}

func (client *AlibabacloudStackClient) WithCrEeClient(do func(*cr_ee.Client) (interface{}, error)) (interface{}, error) {
	if client.creeconn == nil {
		conn, error := client.WithProductSDKClient(CREECode)
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

func (client *AlibabacloudStackClient) NewAdsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("ads", client.Config.Endpoints[ADBCode])
}

func (client *AlibabacloudStackClient) NewCmsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("Cms", client.Config.Endpoints[CMSCode])
}

func (client *AlibabacloudStackClient) WithTableStoreClient(instanceName string, do func(*tablestore.TableStoreClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the TABLESTORE client if necessary
	tableStoreClient, ok := client.tablestoreconnByInstanceName[instanceName]
	if !ok {
		endpoint := client.Config.Endpoints[OtsDataCode]
		if endpoint == "" {
			return nil, fmt.Errorf("[ERROR] missing the product Ots endpoint.")
		}

		transport := &http.Transport{
			MaxIdleConns:    2000,
			IdleConnTimeout: 90 * time.Second,
		}

		if client.Config.Insecure {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		if client.Config.Proxy != "" {
			if proxyURL, err := url.Parse(client.Config.Proxy); err == nil {
				transport.Proxy = http.ProxyURL(proxyURL)
			}
		}

		config := tablestore.NewDefaultTableStoreConfig()
		config.Transport = transport

		endpoint = fmt.Sprintf("%s://%s.%s", strings.ToLower(client.Config.Protocol), instanceName, endpoint)
		tableStoreClient = tablestore.NewClientWithConfig(endpoint, instanceName, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken, config)
		client.tablestoreconnByInstanceName[instanceName] = tableStoreClient
	}

	return do(tableStoreClient)
}
func (client *AlibabacloudStackClient) WithOtsClient(do func(*ots.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the OTS client if necessary
	if client.otsconn == nil {
		conn, error := client.WithProductSDKClient(OtsCode)
		if error != nil {
			return nil, error
		}
		client.otsconn = &ots.Client{
			Client: *conn,
		}
	}

	return do(client.otsconn)
}

func (client *AlibabacloudStackClient) NewVpcClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("vpc", client.Config.Endpoints[VPCCode])
}
func (client *AlibabacloudStackClient) NewEcsClient() (*rpc.Client, error) {
	//sdkConfig.SetEndpoint(endpoint).SetReadTimeout(60000)
	return client.NewTeaSDkClient("ecs", client.Config.Endpoints[EcsCode])
}

func (client *AlibabacloudStackClient) NewRosClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("ros", client.Config.Endpoints[RosCode])
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
		UserAgent:       tea.String(client.GetUserAgent()),
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

func (client *AlibabacloudStackClient) NewArmsClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("arms", client.Config.Endpoints[ARMSCode])
}

func (client *AlibabacloudStackClient) NewBastionhostClient() (*rpc.Client, error) {
	return client.NewTeaSDkClient("Bastionhostprivate", client.Config.Endpoints[BastionHostCode])
}

func (client *AlibabacloudStackClient) defaultHeaders(popcode string) map[string]string {
	return map[string]string{
		"RegionId":                client.RegionId, //	ASAPI
		"x-acs-organizationid":    client.Department,
		"x-acs-resourcegroupid":   client.ResourceGroup,
		"x-acs-regionid":          client.RegionId,
		"x-acs-request-version":   "v1",
		"x-acs-asapi-product":     popcode,
		"x-ascm-product-name":     popcode,
		"EagleEye-TraceId":        client.Eagleeye.GetTraceId(),
		"EagleEye-RpcId":          client.Eagleeye.GetRpcId(),
		"x-acs-territory":         "US",
		"x-acs-lang":              "EN",
		"x-acs-caller-sdk-source": "Terraform",
		//"x-acs-asapi-gateway-version": "3.0"  This specifies to use the ASAPI v3 gateway process, currently maintained is v4, by default it will use v4. Specifying to use v3 is not recommended unless there are compatibility issues that require it.
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
	if strings.ToLower(request.GetProduct()) == "kms" {
		request.Scheme = "https"
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

func buildClientToken(popcode, version, action string) string {
	token := strings.TrimSpace(fmt.Sprintf("TF_%s_%s_%s_%s", popcode, version, action, uuid.Must(uuid.NewV7()).String()))
	return token
}

func requestErrorHandler(api string, response map[string]interface{}, err error, retryTimes int) (*resource.RetryError, int) {
	if err == nil {
		success := true

		// Check HTTP status code (must be 2xx)
		if v, ok := response["HttpStatusCode"]; ok {
			code := strings.TrimSpace(fmt.Sprintf("%v", v))
			if len(code) == 3 && code[0] != '2' {
				success = false
			}
		}

		// Check asapiSuccess flag
		if v, ok := response["asapiSuccess"]; ok && fmt.Sprintf("%v", v) == "false" {
			success = false
		}

		// Check success flag
		if v, ok := response["success"]; ok && fmt.Sprintf("%v", v) == "false" {
			success = false
		}

		if !success {
			var errmsg string
			// Safely extract error message (avoid panic on type assert)
			if v, ok := response["asapiErrorMessage"]; ok {
				if msg, ok := v.(string); ok {
					errmsg = msg
				}
			} else if v, ok := response["errorMessage"]; ok {
				if msg, ok := v.(string); ok {
					errmsg = msg
				}
			}
			err = errmsgs.GetRequestFailedError(
				fmt.Sprintf(errmsgs.RequestV1ErrorMsg, "Request API", api, errmsgs.AlibabacloudStackSdkGoERROR, errmsg),
			)
		}
	}

	if err != nil {
		if errmsgs.NotFoundError(err) {
			return resource.NonRetryableError(err), retryTimes
		}

		if errmsgs.NeedRetry(err) {
			return resource.RetryableError(err), retryTimes
		}

		// Timeout or transient errors
		if errmsgs.IsExpectedErrors(err, errmsgs.ThrottlingUser, errmsgs.Throttling, errmsgs.LogClientTimeout, "ONS_SYSTEM_FLOW_CONTROL", "LockTimeout", "RequestTimeout", "asapi.server.timeout.socket") {
			return resource.RetryableError(err), retryTimes
		}

		// Auth or invalid action errors with retry budget
		if errmsgs.IsExpectedErrors(err, "Forbidden.RAM", "InvalidAction.NotFound", "ServiceUnavailable", "UnknownError") && retryTimes > 0 {
			retryTimes--
			return resource.RetryableError(err), retryTimes
		}

		// Fallback: non-retryable error
		errmsg := errmsgs.GetAsapiErrorMessage(response)
		wrappedErr := errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "Request API", api, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return resource.NonRetryableError(wrappedErr), retryTimes
	}

	return nil, retryTimes
}

func (client *AlibabacloudStackClient) DoTeaRequest(method, popcode, version, apiname, pathpattern string, headers map[string]string, query, body map[string]interface{}) (_result map[string]interface{}, _err error) {
	ServiceCodeStr := strings.ReplaceAll(strings.ToUpper(popcode), "-", "_")
	endpoint := client.Config.Endpoints[ServiceCode(ServiceCodeStr)]
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", popcode)
	}

	reqHeaders := make(map[string]*string)
	for key, value := range headers {
		v := value
		reqHeaders[key] = &v
	}
	for key, value := range client.defaultHeaders(popcode) {
		v := value
		reqHeaders[key] = &v
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
	query["ClientToken"] = buildClientToken(popcode, version, apiname)

	var protocol string
	if strings.ToLower(client.Config.Protocol) == "https" {
		protocol = "https"
	} else {
		protocol = "http"
	}
	switch popcode {
	case "CloudDns", "bms", "EasyAI":
		protocol = "http"
	case "CSB":
		protocol = "https"
	}
	authType := "AK"

	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	if client.Config.Proxy != "" {
		runtime.HttpProxy = &client.Config.Proxy
		runtime.HttpsProxy = &client.Config.Proxy
	}
	runtime.SetAutoretry(false) // When using ASAPI, the Tea package cannot retry, as it will modify the endpoint
	if client.Config.ClientReadTimeout > 0 {
		log.Printf("====================================================================== client.Config.ClientReadTimeout: %d", client.Config.ClientReadTimeout)
	}
	if client.Config.ClientConnectTimeout > 0 {
		log.Printf("====================================================================== client.Config.ClientConnectTimeout: %d", client.Config.ClientConnectTimeout)
	}
	readTimeout := client.Config.ClientReadTimeout * 1000
	connectTimeout := client.Config.ClientConnectTimeout * 1000
	// runtime.ConnectTimeout = &runtimeout
	runtime.SetConnectTimeout(connectTimeout)
	runtime.SetReadTimeout(readTimeout)

	var response map[string]interface{}
	wait := IncrementalWait(3*time.Second, 3*time.Second)
	retryTimes := 3
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {

		var err error
		// conn must be initialized before each request, as DoRequest will modify the content of conn, which will cause the configuration of conn to be invalid for the next request
		if pathpattern != "" {
			response, err = func() (map[string]interface{}, error) {
				sdkConfig := client.teaRoaSdkConfig
				sdkConfig.SetEndpoint(endpoint).SetReadTimeout(readTimeout)    // Unit: milliseconds
				sdkConfig.SetEndpoint(endpoint).SetReadTimeout(connectTimeout) // Unit: milliseconds
				sdkConfig.SetProtocol(protocol)
				conn, err := roa.NewClient(&sdkConfig)
				if err != nil {
					return nil, err
				}
				roa_query := rpcutil.Query(tea.ToMap(map[string]interface{}{
					"Format":         "json",
					"Timestamp":      tea.StringValue(rpcutil.GetTimestamp()),
					"Version":        tea.StringValue(&version),
					"SignatureNonce": tea.StringValue(util.GetNonce()),
				}, query))
				r, e := conn.DoRequestWithAction(&apiname, &version, &protocol, &method, &authType, &pathpattern, roa_query, reqHeaders, body, &runtime)
				if e != nil {
					return r, e
				} else {
					// Handle both map and slice responses
					bodyData, ok := r["body"].(map[string]interface{})
					if !ok {
						// If it's not a map, check if it's a slice
						if sliceData, ok := r["body"].([]interface{}); ok {
							// Convert slice to map with a "data" key
							return map[string]interface{}{"data": sliceData}, e
						}
						// If neither, return empty map
						return map[string]interface{}{}, e
					}
					return bodyData, e
				}
			}()
		} else {
			response, err = func() (map[string]interface{}, error) {
				sdkConfig := client.teaRpcSdkConfig
				sdkConfig.SetEndpoint(endpoint).SetReadTimeout(client.Config.ClientReadTimeout * 1000) // Unit: milliseconds
				sdkConfig.SetProtocol(protocol)
				conn, err := rpc.NewClient(&sdkConfig)
				if err != nil {
					return nil, err
				}
				conn.Headers = reqHeaders
				return conn.DoRequest(&apiname, &protocol, &method, &version, &authType, query, body, &runtime)
			}()
		}

		log.Printf(" ================================ %s ======================================\n query %#v \n request %#v \n response: %#v", apiname, query, body, response)
		var retryErr *resource.RetryError
		retryErr, retryTimes = requestErrorHandler(fmt.Sprintf("%s_%s_%s", popcode, version, apiname), response, err, retryTimes)
		if retryErr != nil {
			wait()
		}
		return retryErr
	})
	return response, err
}

func (client *AlibabacloudStackClient) getConnectClient(popcode ServiceCode) (*sdk.Client, error) {
	client.connsMu.Lock()
	defer client.connsMu.Unlock()
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

func (client *AlibabacloudStackClient) GetAccountInfo() string {
	accountMap := make(map[string]interface{})

	// User information
	// accountMap["aliyunPk"] = nil
	// accountMap["accountStructure"] = nil
	// accountMap["parentPk"] = nil
	accountMap["accessKeyId"] = client.Config.AccessKey
	accountMap["accessKeySecret"] = client.Config.SecretKey
	// accountMap["partnerPk"] = nil

	// // Local machine IP
	// accountMap["sourceIp"] = nil

	// STS required
	if client.Config.SecurityToken != "" {
		accountMap["securityToken"] = client.Config.SecurityToken
	}

	jsonData, _ := json.Marshal(accountMap)
	accountInfo := base64.StdEncoding.EncodeToString([]byte(jsonData))
	return accountInfo
}

func (client *AlibabacloudStackClient) ProcessCommonRequest(request *requests.CommonRequest) (*responses.CommonResponse, error) {
	popcode := ServiceCode(strings.ReplaceAll(strings.ToUpper(request.Product), "-", "_"))

	conn, err := client.getConnectClient(popcode)
	if err != nil {
		return nil, err
	}

	//request.Domain = conn.Domain
	domain := request.Domain
	if domain == "" {
		domain = conn.Domain
	}

	if popcode == OneRouterCode {
		// special logic, 3.16.2 mandatory, no longer required after 3.18.1
		// remove in 3.21.0
		request.QueryParams["AccountInfo"] = client.GetAccountInfo()
	}

	request.QueryParams["ClientToken"] = buildClientToken(request.Product, request.Version, request.ApiName)

	if strings.HasPrefix(domain, "internal.asapi.") || strings.HasPrefix(domain, "public.asapi.") {
		// asapi compatibility logic
		// # asapi When using common SDK, pathpattern cannot be concatenated, otherwise an error will be reported
		if request.PathPattern != "" {
			var r []string = strings.SplitN(domain, "/", 2)
			request.Domain = r[0]
			request.PathPattern = "/asapi/v3"
		}
		if len(request.Content) > 0 {
			request.QueryParams["x-acs-body"] = string(request.Content)
		}
		request.Method = "POST"
		if strings.HasPrefix(domain, "public.asapi.") {
			// If it's a public asapi gateway, force HTTPS
			request.SetScheme("https")
		} else {
			// If it's an internal asapi gateway, force HTTP
			request.SetScheme("http")
		}

	}

	if request.Product == "CloudDns" || request.Product == "bms" {
		// CloudDns / bms does not support HTTPS
		request.SetScheme("http")
	} else if popcode == "CSB" {
		request.SetScheme("https")
	}

	var response *responses.CommonResponse
	wait := IncrementalWait(3*time.Second, 3*time.Second)
	retryTimes := 3
	resource.Retry(5*time.Minute, func() *resource.RetryError {
		// Retry only when the request does not return normally
		response, err = conn.ProcessCommonRequest(request)
		resp := map[string]interface{}{}
		json.Unmarshal(response.GetHttpContentBytes(), &resp)
		if response == nil {
			retryTimes -= 1
			wait()
			return resource.RetryableError(err)
		}
		var retryErr *resource.RetryError
		retryErr, retryTimes = requestErrorHandler(fmt.Sprintf("%s_%s_%s", request.Product, request.Version, request.ApiName), resp, err, retryTimes)
		if retryErr != nil {
			wait()
		}
		return retryErr
	})
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

func retryDo(do func() (interface{}, error)) (interface{}, error) {
	var response interface{}
	var err error
	wait := IncrementalWait(3*time.Second, 3*time.Second)
	resource.Retry(5*time.Minute, func() *resource.RetryError {
		// Retry only when the request does not return normally
		response, err = do()
		if err == nil {
			return nil
		}
		if response == nil {
			wait()
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)

	})
	return response, err
}

func (client *AlibabacloudStackClient) GetRetryTimeout(defaultTimeout time.Duration) time.Duration {

	maxRetryTimeout := client.Config.MaxRetryTimeout
	if maxRetryTimeout != 0 {
		return time.Duration(maxRetryTimeout) * time.Second
	}

	return defaultTimeout
}
