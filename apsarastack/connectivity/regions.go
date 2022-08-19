package connectivity

// Region represents ECS region
type Region string

const (
	NeimengEnv30  = Region("cn-neimeng-env30-d01")
	NeimengEnv127 = Region("cn-neimeng-env127-d01")
	QingdaoEnv66  = Region("cn-qingdao-env66-d01")
	QingdaoEnv17  = Region("cn-qingdao-env17-d01")
	WulanEnv230   = Region("cn-wulan-env230-d01")
	WulanEnv205   = Region("cn-wulan-env205-d01")
	WulanEnv82    = Region("cn-wulan-env82-d01")
	QingdaoEnv15  = Region("cn-qingdao-env15-d01")
	Hangzhou      = Region("cn-hangzhou")
	Qingdao       = Region("cn-qingdao")
	Beijing       = Region("cn-beijing")
	Hongkong      = Region("cn-hongkong")
	Shenzhen      = Region("cn-shenzhen")
	Shanghai      = Region("cn-shanghai")
	Zhangjiakou   = Region("cn-zhangjiakou")
	Huhehaote     = Region("cn-huhehaote")
	ChengDu       = Region("cn-chengdu")
	HeYuan        = Region("cn-heyuan")
	WuLanChaBu    = Region("cn-wulanchabu")

	APSouthEast1 = Region("ap-southeast-1")
	APNorthEast1 = Region("ap-northeast-1")
	APSouthEast2 = Region("ap-southeast-2")
	APSouthEast3 = Region("ap-southeast-3")
	APSouthEast5 = Region("ap-southeast-5")

	APSouth1 = Region("ap-south-1")

	USWest1 = Region("us-west-1")
	USEast1 = Region("us-east-1")

	MEEast1 = Region("me-east-1")

	EUCentral1 = Region("eu-central-1")
	EUWest1    = Region("eu-west-1")

	ShenZhenFinance = Region("cn-shenzhen-finance-1")
	ShanghaiFinance = Region("cn-shanghai-finance-1")

	GuangZhou    = Region("cn-guangzhou")
	APSouthEast6 = Region("ap-southeast-6")
)

var ValidRegions = []Region{
	NeimengEnv30, QingdaoEnv66, QingdaoEnv17, QingdaoEnv15, WulanEnv82, Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, WuLanChaBu,
	USWest1, USEast1,
	APNorthEast1, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5,
	APSouth1,
	MEEast1,
	EUCentral1, EUWest1,
}
var OtsCapacityNoSupportedRegions = []Region{APSouthEast1, USWest1, USEast1}
var EcsClassicSupportedRegions = []Region{Shenzhen, Shanghai, Beijing, Qingdao, Hangzhou, Hongkong, USWest1, APSouthEast1}
var EcsSpotNoSupportedRegions = []Region{APSouth1}
var RdsClassicNoSupportedRegions = []Region{APSouth1, APSouthEast2, APSouthEast3, APNorthEast1, EUCentral1, EUWest1, MEEast1}
var RdsMultiAzNoSupportedRegions = []Region{Qingdao, APNorthEast1, APSouthEast5, MEEast1}
var RdsPPASNoSupportedRegions = []Region{Qingdao, USEast1, APNorthEast1, EUCentral1, MEEast1, APSouthEast2, APSouthEast3, APSouth1, APSouthEast5, ChengDu, EUWest1}
var RouteTableNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var SlbClassicNoSupportedRegions = []Region{APNorthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, MEEast1, EUCentral1, EUWest1, Huhehaote, Zhangjiakou}
var HttpHttpsHealthCheckMehtodSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, EUWest1, ChengDu, Qingdao, Hongkong, Shenzhen, APSouthEast5, Zhangjiakou, Huhehaote, MEEast1, APSouth1, EUCentral1, USWest1, APSouthEast3, APSouthEast2, APSouthEast1, APNorthEast1}
var OssVersioningSupportedRegions = []Region{APSouth1}
var OssSseSupportedRegions = []Region{Qingdao, Hangzhou, Beijing, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouth1, USEast1}
var CRNoSupportedRegions = []Region{Beijing, Hangzhou, Qingdao, Huhehaote, Zhangjiakou}
var KmsSkippedRegions = []Region{}
var KubernetesSupportedRegions = []Region{Beijing, Zhangjiakou, Huhehaote, Hangzhou, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouthEast1,
	APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, EUWest1, MEEast1, EUCentral1}
var GpdbClassicNoSupportedRegions = []Region{APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1, EUCentral1}
var MongoDBClassicNoSupportedRegions = []Region{Huhehaote, Zhangjiakou, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1}
var MongoDBMultiAzSupportedRegions = []Region{Hangzhou, Beijing, Shenzhen, EUCentral1}
var DatahubSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, Shenzhen, APSouthEast1}
var OtsHighPerformanceNoSupportedRegions = []Region{Qingdao, Zhangjiakou, Huhehaote, APSouthEast2, APSouthEast5, APNorthEast1, EUCentral1, MEEast1}
var EdasSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, Shenzhen, Zhangjiakou, Qingdao, Hongkong}
var AdbReserverUnSupportRegions = []Region{EUCentral1}
var VpcIpv6GatewaySupportRegions = []Region{QingdaoEnv15, QingdaoEnv66, WulanEnv205, WulanEnv230, NeimengEnv127, Qingdao, Beijing, Zhangjiakou, Huhehaote, WuLanChaBu, Hangzhou, Shanghai, Shenzhen, GuangZhou, Hongkong, ChengDu, HeYuan, APSouthEast1, APSouthEast6, USEast1, EUCentral1}
var VpcIpv6SupportRegions = []Region{QingdaoEnv15, Hangzhou, Shanghai, Shenzhen, Beijing, Huhehaote, Hongkong, APSouthEast1}
var ManagedKubernetesSupportedRegions = []Region{QingdaoEnv15, Beijing, Hangzhou, Shanghai, Shenzhen, ChengDu, Hongkong, APSouthEast1, APSouthEast2, EUCentral1, USWest1}
var AlikafkaSupportedRegions = []Region{Hangzhou, Qingdao, Beijing, Hongkong, Shenzhen, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, APNorthEast1, APSouthEast1, APSouthEast3, EUCentral1, EUWest1, USEast1, USWest1}
var NasClassicSupportedRegions = []Region{Hangzhou, Qingdao, Beijing, Hongkong, Shenzhen, Shanghai, Zhangjiakou, Huhehaote, ShenZhenFinance, ShanghaiFinance}
var NasNoSupportedRegions = []Region{Qingdao, APSouth1, APSouthEast3, APSouthEast5}
var ApiGatewayNoSupportedRegions = []Region{Zhangjiakou, Huhehaote, USEast1, USWest1, EUWest1, MEEast1}
var FcNoSupportedRegions = []Region{MEEast1}
var DrdsSupportedRegions = []Region{Beijing, Shenzhen, Hangzhou, Qingdao, Hongkong, Shanghai, Huhehaote, Zhangjiakou, APSouthEast1}
var DrdsClassicNoSupportedRegions = []Region{Hongkong}
var AlbSupportRegions = []Region{Hangzhou, Shanghai, Qingdao, Zhangjiakou, Beijing, WuLanChaBu, Shenzhen, ChengDu, Hongkong, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APNorthEast1, EUCentral1, USEast1, APSouth1}
