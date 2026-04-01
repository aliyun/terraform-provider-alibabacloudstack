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
	DmsEnterpriseCode    = ServiceCode("DMS_ENTERPRISE")
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
	KVSTORECode          = ServiceCode("R_KVSTORE") // Do not allow "-", schema does not accept it, use "_" instead
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
	CLOUDFWCode          = ServiceCode("CLOUDFW")
	FlinkCode            = ServiceCode("VERVERICA")
	EBSCode              = ServiceCode("EBS")
	YaochiOpsCode        = ServiceCode("YAOCHIOPS")
	ACMCode              = ServiceCode("ACM")
	POLARDBXCode         = ServiceCode("POLARDBX")
	HOLOGRAMCode         = ServiceCode("HOLOGRAM")
	HitsdbCode           = ServiceCode("HITSDB")
	CSB2Code             = ServiceCode("CSB2")
	OnsInnerCode         = ServiceCode("ONS_INNER")
	UniversalDnsCode     = ServiceCode("UNIVERSALDNS")
	SCHEDULERX2Code      = ServiceCode("SCHEDULERX2")
	EFSCode              = ServiceCode("EFS")
	HSMPRIVATECode       = ServiceCode("HSM_PRIVATE")
	CSPPRIVATECode       = ServiceCode("CSPPRIVATE")
	BMSCode              = ServiceCode("BMS")
	AEGISCode            = ServiceCode("AEGIS")
	// undefined code, add first
	GDBCode              = ServiceCode("GDB")
	ARMSCode             = ServiceCode("ARMS")
	CSBCode              = ServiceCode("CSB")
	DBSCode              = ServiceCode("DBS")
	DTSCode              = ServiceCode("DTS")
	SLSCode              = ServiceCode("SLS")
	RosCode              = ServiceCode("ROS")
	QuickbiCode          = ServiceCode("QUICKBI")
	DataworksPublicCode  = ServiceCode("DATAWORKS_PUBLIC")
	DataworksPrivateCode = ServiceCode("DATAWORKS_PRIVATE_CLOUD")
	OneRouterCode        = ServiceCode("ONEROUTER")
	BastionHostCode      = ServiceCode("BASTIONHOSTPRIVATE")
	WAFONECSCode         = ServiceCode("WAF_ONECS")
	ALIKAFKACode         = ServiceCode("ALIKAFKA")
	Prometheus2Code      = ServiceCode("PROMETHEUS2")
	TablestoreCode       = ServiceCode("TABLESTORE")
	EasyAICode           = ServiceCode("EASYAI")
	// Self-built gateway fake Code
	SlSDataCode = ServiceCode("SLSDATA")
	BssDataCode = ServiceCode("BSSDATA")
	OtsDataCode = ServiceCode("OTSDATA")

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
		"arms-api.{{.domain}}",
		"arms-api.{{.region}}.{{.domain}}",
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
	//cloudfw endpoint
	CLOUDFWCode: PopEndpoint{
		"cloudcfw-biz.{{.region}}.{{.domain}}",
		"cloudcfw-biz.{{.region}}.{{.domain}}",
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
	// waf-onecs endpoint
	WAFONECSCode: PopEndpoint{
		"waf-onecs-biz.{{.region}}.{{.domain}}",
		"waf-onecs-biz.{{.region}}.{{.domain}}",
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
		"sts-vpc.{{.domain}}",
		"sts-vpc.{{.region}}.{{.domain}}",
	},
	EBSCode: PopEndpoint{
		"ebsnext.{{.domain}}",
		"ebsnext.{{.region}}.{{.domain}}",
	},
	YaochiOpsCode: PopEndpoint{
		"yaochiops.{{.domain}}",
		"yaochiops.{{.region}}.{{.domain}}",
	},
	POLARDBCode: PopEndpoint{
		"polardb-vpc.{{.domain}}",
		"polardb-vpc.{{.region}}.{{.domain}}",
	},

	SlSDataCode: PopEndpoint{
		"data.{{.region}}.sls-pub.{{.domain}}",
		"data.{{.region}}.sls-pub.{{.domain}}",
	},

	OtsDataCode: PopEndpoint{
		"{{.region}}.ots-internal.{{.domain}}",
		"{{.region}}.ots-internal.{{.domain}}",
	},
	DmsEnterpriseCode: PopEndpoint{
		"newdms-api.{{.domain}}",
		"newdms-api.{{.region}}.{{.domain}}",
	},
	OSSCode: PopEndpoint{"", ""},
	DataworksPublicCode: PopEndpoint{
		"dataworks-public.{{.domain}}",
		"dataworks-public.{{.region}}.{{.domain}}",
	},
	DDSCode: PopEndpoint{
		"mongodb-vpc.{{.domain}}",
		"mongodb-vpc.{{.region}}.{{.domain}}",
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
		"kafka.biz.openapi.{{.domain}}",
		"kafka.biz.openapi.{{.region}}.{{.domain}}",
	},
	NasCode: PopEndpoint{
		"nas.{{.region}}.{{.domain}}",
		"nas.{{.region}}.{{.domain}}",
	},
	RosCode: PopEndpoint{
		"ros.{{.domain}}",
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
		"hbase-inc.{{.domain}}",
		"hbase-inc.{{.region}}.{{.domain}}",
	},
	ONSCode: PopEndpoint{
		"ons-biz.{{.region}}.{{.domain}}",
		"ons-biz.{{.region}}.{{.domain}}",
	},
	Prometheus2Code: PopEndpoint{
		"prometheus2-biz.pop.{{.domain}}",
		"prometheus2-biz.pop.{{.region}}.{{.domain}}",
	},
	CDNCode: PopEndpoint{"", ""},
	QuickbiCode: PopEndpoint{
		"quickbi-public.{{.domain}}",
		"quickbi-public.{{.domain}}",
	},
	BssDataCode: PopEndpoint{"", ""},
	BastionHostCode: PopEndpoint{
		"bastionhostprivate-biz.{{.region}}.{{.domain}}",
		"bastionhostprivate-biz.{{.region}}.{{.domain}}",
	},
	OnsInnerCode: PopEndpoint{
		"ons-biz.{{.region}}.{{.domain}}",
		"ons-biz.{{.region}}.{{.domain}}",
	},
	FlinkCode: PopEndpoint{
		"ververica.{{.region}}.{{.domain}}",
		"ververica.{{.region}}.{{.domain}}",
	},
	ACMCode: PopEndpoint{
		"dncs-api.console.{{.region}}.{{.domain}}",
		"dncs-api.console.{{.region}}.{{.domain}}",
	},
	POLARDBXCode: PopEndpoint{
		"polarx-vpc.{{.domain}}",
		"polarx-vpc.{{.region}}.{{.domain}}",
	},
	DataworksPrivateCode: PopEndpoint{
		"dataworks-vpc.{{.domain}}",
		"dataworks-vpc.{{.region}}.{{.domain}}",
	},
	//cbn endpoint
	CbnCode: PopEndpoint{
		"cbn.{{.domain}}",
		"cbn.{{.region}}.{{.domain}}",
	},
	HOLOGRAMCode: PopEndpoint{
		"hologram.{{.domain}}",
		"hologram.{{.region}}.{{.domain}}",
	},
	OtsCode: PopEndpoint{
		"ots.{{.domain}}",
		"ots.{{.region}}.{{.domain}}",
	},
	TablestoreCode: PopEndpoint{
		"tablestore.{{.domain}}",
		"tablestore.{{.region}}.{{.domain}}",
	},
	CSB2Code: PopEndpoint{
		"csb-cop-api-biz.{{.domain}}",
		"csb-cop-api-biz.{{.region}}.{{.domain}}",
	},
	UniversalDnsCode: PopEndpoint{
		"dns-universal-control.pop.{{.domain}}",
		"dns-universal-control.pop.{{.region}}.{{.domain}}",
	},
	SCHEDULERX2Code: PopEndpoint{
		"schedulerx2-api-{{.region}}.{{.domain}}",
		"schedulerx2-api-{{.region}}.{{.domain}}",
	},
	EFSCode: PopEndpoint{
		"efs-pub.{{.region}}.{{.domain}}",
		"efs-pub.{{.region}}.{{.domain}}",
	},
	HSMPRIVATECode: PopEndpoint{
		"hsmprivate.{{.region}}.{{.domain}}",
		"hsmprivate.{{.region}}.{{.domain}}",
	},
	CSPPRIVATECode: PopEndpoint{
		"kms-csp-pop-biz.{{.region}}.{{.domain}}",
		"kms-csp-pop-biz.{{.region}}.{{.domain}}",
	},
	BMSCode: PopEndpoint{
		"baremetalservice.{{.region}}.{{.domain}}",
		"baremetalservice.{{.region}}.{{.domain}}",
	},
	AEGISCode: PopEndpoint{
		"aegis-biz.{{.region}}.{{.domain}}",
		"aegis-biz.{{.region}}.{{.domain}}",
	},
	// 3.18.3 new sites will not be opened
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
	EasyAICode: PopEndpoint{
		"easy-ai.{{.region}}.{{.domain}}",
		"easy-ai.{{.region}}.{{.domain}}",
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
