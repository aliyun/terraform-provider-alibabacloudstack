package alibabacloudstack

type GroupRuleNicType string

type Direction string

const (
	DirectionIngress = Direction("ingress")
	DirectionEgress  = Direction("egress")
)
const (
	DiskTypeAll    = DiskType("all")
	DiskTypeSystem = DiskType("system")
	DiskTypeData   = DiskType("data")
)
const (
	GroupRuleInternet = GroupRuleNicType("internet")
	GroupRuleIntranet = GroupRuleNicType("intranet")
)
const (
	GroupInnerAccept = GroupInnerAccessPolicy("Accept")
	GroupInnerDrop   = GroupInnerAccessPolicy("Drop")
)
const (
	GroupRulePolicyAccept = GroupRulePolicy("accept")
	GroupRulePolicyDrop   = GroupRulePolicy("drop")
)

const AllPortRange = "-1/-1"

type GroupRulePolicy string

type GroupInnerAccessPolicy string

type SpotStrategyType string

const (
	NoSpot             = SpotStrategyType("NoSpot")
	SpotWithPriceLimit = SpotStrategyType("SpotWithPriceLimit")
	SpotAsPriceGo      = SpotStrategyType("SpotAsPriceGo")
)
const (
	RenewAutoRenewal = RenewalStatus("AutoRenewal")
	RenewNormal      = RenewalStatus("Normal")
	RenewNotRenewal  = RenewalStatus("NotRenewal")
)

type DestinationResource string

const (
	InstanceTypeResource = DestinationResource("InstanceType")
	ZoneResource         = DestinationResource("Zone")
)
const (
	DiskResizeTypeOffline = DiskResizeType("offline")
	DiskResizeTypeOnline  = DiskResizeType("online")
)

type RenewalStatus string

type DiskType string

type DiskCategory string

const (
	DiskAll             = DiskCategory("all") //Default
	DiskCloud           = DiskCategory("cloud")
	DiskEphemeralSSD    = DiskCategory("ephemeral_ssd")
	DiskCloudESSD       = DiskCategory("cloud_essd")
	DiskCloudEfficiency = DiskCategory("cloud_efficiency")
	DiskCloudSSD        = DiskCategory("cloud_ssd")
	DiskCloudPPERF      = DiskCategory("cloud_pperf")
	DiskCloudSPERF      = DiskCategory("cloud_sperf")
	DiskLocalDisk       = DiskCategory("local_disk")
)

type DiskResizeType string

type ImageOwnerAlias string

type SecurityEnhancementStrategy string

const (
	ActiveSecurityEnhancementStrategy   = SecurityEnhancementStrategy("Active")
	DeactiveSecurityEnhancementStrategy = SecurityEnhancementStrategy("Deactive")
)

type CreditSpecification string
