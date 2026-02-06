package connectivity

import (
	"bytes"
	"text/template"
)

// ServiceCode Load endpoints from endpoints.xml or environment variables to meet specified application scenario, like private cloud.
type ServiceCode string

const (
	DcdnCode             = ServiceCode("DCDN")
	MseCode              = ServiceCode("MSE")
	ActiontrailCode      = ServiceCode("ACTIONTRAIL")
	OosCode              = ServiceCode("OOS")
	EcsCode              = ServiceCode("ECS")
	ASCMCode             = ServiceCode("ASCM")
	NasCode              = ServiceCode("NAS")
	EciCode              = ServiceCode("ECI")
	DdoscooCode          = ServiceCode("DDOSCOO")
	AlidnsCode           = ServiceCode("ALIDNS")
	ResourcemanagerCode  = ServiceCode("RESOURCEMANAGER")
	WafOpenapiCode       = ServiceCode("WAFOPENAPI")
	DmsEnterpriseCode    = ServiceCode("DMSENTERPRISE")
	DnsCode              = ServiceCode("ALIDNS")
	KmsCode              = ServiceCode("KMS")
	CbnCode              = ServiceCode("CBN")
	ESSCode              = ServiceCode("ESS")
	RAMCode              = ServiceCode("RAM")
	VPCCode              = ServiceCode("VPC")
	SLBCode              = ServiceCode("SLB")
	RDSCode              = ServiceCode("RDS")
	OSSCode              = ServiceCode("OSS")
	ONSCode              = ServiceCode("ONS")
	CONTAINCode          = ServiceCode("CS")
	CRCode               = ServiceCode("CR")
	CREECode             = ServiceCode("CR_EE")
	CDNCode              = ServiceCode("CDN")
	CMSCode              = ServiceCode("CMS")
	DNSCode              = ServiceCode("CLOUDDNS")
	PVTZCode             = ServiceCode("PVTZ")
	LOGCode              = ServiceCode("LOG")
	FCCode               = ServiceCode("FC")
	DDSCode              = ServiceCode("DDS")
	GPDBCode             = ServiceCode("GPDB")
	CENCode              = ServiceCode("CEN")
	KVSTORECode          = ServiceCode("R_KVSTORE") // Hyphens are not allowed, schema does not accept them, use _ instead
	POLARDBCode          = ServiceCode("POLARDB")
	MNSCode              = ServiceCode("MNS")
	CLOUDAPICode         = ServiceCode("CLOUDAPI")
	DRDSCode             = ServiceCode("DRDS")
	LOCATIONCode         = ServiceCode("LOCATION")
	ElasticsearchK8sCode = ServiceCode("ELASTICSEARCH_K8S")
	DDOSCOOCode          = ServiceCode("DDOSCOO")
	DDOSBGPCode          = ServiceCode("DDOSBGP")
	SAGCode              = ServiceCode("SAG")
	EMRCode              = ServiceCode("EMR")
	CasCode              = ServiceCode("CAS")
	YUNDUNDBAUDITCode    = ServiceCode("YUNDUNDBAUDIT")
	MARKETCode           = ServiceCode("MARKET")
	HBASECode            = ServiceCode("HBASE")
	ADBCode              = ServiceCode("ADB")
	EDASCode             = ServiceCode("EDAS")
	CassandraCode        = ServiceCode("CASSANDRA")
	OtsCode              = ServiceCode("OTS")
	DatahubCode          = ServiceCode("DATAHUB")
	STSCode              = ServiceCode("STS")
	ACMCode              = ServiceCode("ACM")
	// undefined code, add first
	GDBCode             = ServiceCode("GDB")
	ARMSCode            = ServiceCode("ARMS")
	CSBCode             = ServiceCode("CSB")
	DBSCode             = ServiceCode("DBS")
	DTSCode             = ServiceCode("DTS")
	SLSCode             = ServiceCode("SLS")
	HitsdbCode          = ServiceCode("HITSDB")
	RosCode             = ServiceCode("ROS")
	QuickbiCode         = ServiceCode("QUICKBI")
	DataworkspublicCode = ServiceCode("DATAWORKSPUBLIC")
	OneRouterCode       = ServiceCode("ONEROUTER")
	// Self-built gateway fake Code
	SlSDataCode  = ServiceCode("SLSDATA")
	ALIKAFKACode = ServiceCode("ALIKAFKADATA")
	BssDataCode  = ServiceCode("BSSDATA")

	// ASAPI
	ASAPICode = ServiceCode("ASAPI")
)

type Endpoints struct {
	Endpoint []Endpoint `xml:"Endpoint"`
}

type RegionIds struct {
	RegionId string `xml:"RegionId"`
}

type Products struct {
	Product []Product `xml:"Product"`
}

type Product struct {
	ProductName string `xml:"ProductName"`
	DomainName  string `xml:"DomainName"`
}

type Endpoint struct {
	Name      string    `xml:"name,attr"`
	RegionIds RegionIds `xml:"RegionIds"`
	Products  Products  `xml:"Products"`
}

var serviceCodeMapping = map[string]string{
	"cloudapi": "apigateway",
}

type PopEndpoint struct {
	CenterEndpoint string
	RegionEndpoint string
}

var PopEndpoints = map[ServiceCode]PopEndpoint{
	//vpc endpoint
	VPCCode: PopEndpoint{
		"vpc-internal.{{.domain}}",
		"vpc-internal.{{.region}}.{{.domain}}",
	},
	//slb endpoint
	SLBCode: PopEndpoint{
		"slb-vpc.{{.domain}}",
		"slb-vpc.{{.region}}.{{.domain}}",
	},
	//gdb endpoint
	GDBCode: PopEndpoint{
		"gdb.{{.domain}}",
		"gdb.{{.region}}.{{.domain}}",
	},
	//gpdb endpoint
	GPDBCode: PopEndpoint{
		"gpdb.{{.domain}}",
		"gpdb.{{.region}}.{{.domain}}",
	},
	//adb endpoint
	ADBCode: PopEndpoint{
		"adb.{{.domain}}",
		"adb.{{.region}}.{{.domain}}",
	},
	//apigateway endpoint
	//centralized deployment
	CLOUDAPICode: PopEndpoint{
		"apigateway.{{.region}}.{{.domain}}",
		"apigateway.{{.region}}.{{.domain}}",
	},
	//arms endpoint
	ARMSCode: PopEndpoint{
		"arms-api.console.{{.region}}.{{.domain}}",
		"arms-api.console.{{.region}}.{{.domain}}",
	},
	//ascm endpoint
	ASCMCode: PopEndpoint{
		"ascm.{{.domain}}",
		"ascm.{{.region}}.{{.domain}}",
	},
	//cloudfw endpoint
	WafOpenapiCode: PopEndpoint{
		"cloudfw.{{.domain}}",
		"cloudfw.{{.region}}.{{.domain}}",
	},
	//cr endpoint
	CRCode: PopEndpoint{
		"cr-biz.{{.region}}.{{.domain}}",
		"cr-biz.{{.region}}.{{.domain}}",
	},
	//cr endpoint
	CREECode: PopEndpoint{
		"cr-ee-biz.{{.region}}.{{.domain}}",
		"cr-ee-biz.{{.region}}.{{.domain}}",
	},
	//csb endpoint
	CSBCode: PopEndpoint{
		"csb.{{.domain}}",
		"csb.{{.region}}.{{.domain}}",
	},
	//datahub endpoint
	DatahubCode: PopEndpoint{
		"datahub.{{.region}}.api-pop.{{.domain}}",
		"datahub.{{.region}}.api-pop.{{.domain}}",
	},
	//dbs endpoint
	DBSCode: PopEndpoint{
		"dbs.{{.domain}}",
		"dbs.{{.region}}.{{.domain}}",
	},
	//dns endpoint
	DNSCode: PopEndpoint{
		"dns-control.pop.{{.domain}}",
		"dns-control.pop.{{.region}}.{{.domain}}",
	},
	//drds endpoint
	DRDSCode: PopEndpoint{
		"drds.{{.domain}}",
		"drds.{{.region}}.{{.domain}}",
	},
	//dts endpoint
	DTSCode: PopEndpoint{
		"dts.{{.domain}}",
		"dts.{{.region}}.{{.domain}}",
	},
	//edas-api.console endpoint
	EDASCode: PopEndpoint{
		"edas-api.console.{{.region}}.{{.domain}}",
		"edas-api.console.{{.region}}.{{.domain}}",
	},
	//ELASTICSEARCHCode endpoint
	ElasticsearchK8sCode: PopEndpoint{
		"elasticsearch.k8s.{{.region}}.{{.domain}}",
		"elasticsearch.k8s.{{.region}}.{{.domain}}",
	},
	//Ess endpoint
	ESSCode: PopEndpoint{
		"ess.{{.domain}}",
		"ess.{{.region}}.{{.domain}}",
	},
	//Ecs endpoint
	EcsCode: PopEndpoint{
		"ecs-internal.{{.domain}}",
		"ecs-internal.{{.region}}.{{.domain}}",
	},
	//Sts endpoint
	STSCode: PopEndpoint{
		"sts.{{.domain}}",
		"sts.{{.region}}.{{.domain}}",
	},

	POLARDBCode: PopEndpoint{
		"polardb-vpc.{{.domain}}",
		"polardb-vpc.{{.region}}.{{.domain}}",
	},

	SlSDataCode: PopEndpoint{
		"data.{{.region}}.sls-pub.{{.domain}}",
		"data.{{.region}}.sls-pub.{{.domain}}",
	},

	DmsEnterpriseCode: PopEndpoint{"", ""},
	OSSCode:           PopEndpoint{"", ""},
	DataworkspublicCode: PopEndpoint{
		"dataworks-public.{{.domain}}",
		"dataworks-public.{{.region}}.{{.domain}}",
	},
	DDSCode: PopEndpoint{
		"mongodb.{{.domain}}",
		"mongodb.{{.region}}.{{.domain}}",
	},
	RAMCode: PopEndpoint{
		"ram.{{.domain}}",
		"ram.{{.domain}}",
	},
	CMSCode: PopEndpoint{
		"metrics.open.{{.domain}}",
		"metrics.open.{{.region}}.{{.domain}}",
	},
	HitsdbCode: PopEndpoint{
		"hitsdb.{{.domain}}",
		"hitsdb.{{.region}}.{{.domain}}",
	},
	ALIKAFKACode: PopEndpoint{
		"kafka.openapi.{{.domain}}",
		"kafka.openapi.{{.region}}.{{.domain}}",
	},
	NasCode: PopEndpoint{
		"nas.{{.region}}.{{.domain}}",
		"nas.{{.region}}.{{.domain}}",
	},
	RosCode: PopEndpoint{
		"ros.{{.region}}.{{.domain}}",
		"ros.{{.region}}.{{.domain}}",
	},
	RDSCode: PopEndpoint{
		"rds.{{.domain}}",
		"rds.{{.region}}.{{.domain}}",
	},
	KVSTORECode: PopEndpoint{
		"kvstore-vpc.{{.domain}}",
		"kvstore-vpc.{{.region}}.{{.domain}}",
	},
	OosCode: PopEndpoint{
		"oos-public-inner.{{.domain}}",
		"oos-public-inner.{{.region}}.{{.domain}}",
	},
	CONTAINCode: PopEndpoint{
		"cs-intranet.{{.domain}}",
		"cs-intranet.{{.region}}.{{.domain}}",
	},
	HBASECode: PopEndpoint{
		"hbase.{{.domain}}",
		"hbase.{{.region}}.{{.domain}}",
	},
	ONSCode: PopEndpoint{
		"ons-biz.{{.region}}.{{.domain}}",
		"ons-biz.{{.region}}.{{.domain}}",
	},
	CDNCode:     PopEndpoint{"", ""},
	QuickbiCode: PopEndpoint{"", ""},
	BssDataCode: PopEndpoint{"", ""},
	OtsCode: PopEndpoint{
		"ots.{{.domain}}",
		"ots.{{.region}}.{{.domain}}",
	},
	ACMCode: PopEndpoint{
		"dncs-api.console.{{.region}}.{{.domain}}",
		"dncs-api.console.{{.region}}.{{.domain}}",
	},
	// New sites after 3.18.3 will no longer be open
	OneRouterCode: PopEndpoint{
		"public.asapi.{{.region}}.{{.domain}}",
		"public.asapi.{{.region}}.{{.domain}}",
	},
	ASAPICode: PopEndpoint{
		"public.asapi.{{.region}}.{{.domain}}",
		"public.asapi.{{.region}}.{{.domain}}",
	},
	KmsCode: PopEndpoint{
		"public.asapi.{{.region}}.{{.domain}}",
		"public.asapi.{{.region}}.{{.domain}}",
	},
	SLSCode: PopEndpoint{
		"public.asapi.{{.region}}.{{.domain}}",
		"public.asapi.{{.region}}.{{.domain}}",
	},
}

func GeneratorEndpoint(serviceCode ServiceCode, region string, domain string, isCenter bool) string {
	endpoints := PopEndpoints[serviceCode]

	var err error
	var tmp *template.Template
	if !isCenter {
		tmp, err = template.New(string(serviceCode)).Parse(endpoints.RegionEndpoint)
	} else {
		tmp, err = template.New(string(serviceCode)).Parse(endpoints.CenterEndpoint)
	}
	if err != nil {
		panic(err)
	}

	param := map[string]string{
		"domain": domain,
		"region": region,
	}

	var buffer bytes.Buffer
	if err = tmp.Execute(&buffer, param); err != nil {
		panic(err)
	}

	return buffer.String()
}
