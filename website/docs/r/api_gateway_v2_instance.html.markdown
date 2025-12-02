---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_instance"
sidebar_current: "docs-alibabacloudstack-resource-api-gateway-v2-instance"
description: |-
    Provides a Alibabacloudstack Api Gateway V2 Instance Resource.
---
Provides a Alibabacloudstack Api Gateway V2 Instance Resource.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "%s"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
      ignore_changes = [
		secondary_cidr_blocks,
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_cs_kubernetes" "default" {
	name						= var.name
	version						= "1.30.7-aliyun.1"
	os_type						= "linux"
	platform					= "AliyunLinux"
	num_of_nodes				= "3"
	master_count				= "3"
	master_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
	master_instance_types		= ["ecs.n4v2.large","ecs.n4v2.large","ecs.n4v2.large"]
	master_disk_category		= "cloud_ssd"
	vpc_id						= "${alibabacloudstack_vpc_vpc.default.id}"
	worker_instance_types		= ["ecs.n4v2.large"]
	worker_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}"]
	worker_disk_category		= "cloud_ssd"
	password					= random_password.password.0.result
	pod_cidr					= "172.20.0.0/16"
	service_cidr				= "172.21.0.0/20"
	worker_disk_size			= "40"
	master_disk_size			= "40"
	slb_internet_enabled		= "true"
	security_group_id			= alibabacloudstack_ecs_securitygroup.default.id
	runtime	 {
		name	= "containerd"
		version	= "1.6.28"
	}
}

resource "alibabacloudstack_edas_k8s_cluster" "default" {
	cs_cluster_id	= alibabacloudstack_cs_kubernetes.default.id
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alibabacloudstack_edas_namespace" "default" {
  debug_enable         = false
  description          = var.name
  namespace_logical_id = "${var.region}:${var.name}"
  namespace_name       = var.name
}

resource "alibabacloudstack_api_gateway_v2_instance" "example" {
  instance_name = "example-instance"
  deploy_mode   = "edas"
  broker_engine_type = "SCG"
  instance_class = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
  node_number    = 1
    
  edas_app_infos {
    edas_namespace = "${alibabacloudstack_edas_namespace.default.id}"
    k8s_cluster_id = "${alibabacloudstack_edas_k8s_cluster.default.id}"
  }
}
```

### Custom Deployment Mode

```hcl
resource "alibabacloudstack_api_gateway_v2_instance" "edas_example" {
  instance_name      = "example-instance"
  deploy_mode        = "custom"
  broker_engine_type = "SCG"
  instance_class     = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
  node_number        = 1
}
```

### APIGW K8s Deployment Mode

```hcl
resource "alibabacloudstack_api_gateway_v2_instance" "edas_example" {
	instance_name             = "${var.name}"
	node_number               = "1"
	instance_class            = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
	broker_engine_type        = "SCG"
	deploy_mode               = "apig_k8s"
	deploy_cluster_code       = "${alibabacloudstack_api_gateway_v2_k8s_cluster.default.id}"
	deploy_cluster_namespace  = "${var.name}-namespace"
	sls_enabled               = "true"
	prometheus_enabled        = "true"
}
```

### AiGateway Mode

```hcl
resource "alibabacloudstack_api_gateway_v2_instance" "custom_example" {
  instance_name         = "custom-instance"
  node_number               = "1"
  deploy_mode           = "k8s"
  broker_engine_type    = "HIGRESS"
  deploy_cluster_code   = "cluster-code"
  deploy_cluster_namespace = "testnamespace"
  ingress_class_name =       "test-class"
  instance_class     = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
  prometheus_enabled    = true
  sls_enabled           = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of the API Gateway instance.
* `deploy_mode` - (Optional, ForceNew) The deployment mode of the instance. Valid values: `k8s`, `edas`, `custom`, `apig_k8s`(APIGateWay K8s).
* `broker_engine_type` - (Optional, ForceNew) The broker engine type. Valid values: `HIGRESS`(AIGateWay), `SCG`(APIGateWay).
* `instance_class` - (Optional, ForceNew) The instance class/specification.
* `node_number` - (Optional, ForceNew) The number of nodes for the instance.
* `edas_namespace_id` - (Optional, ForceNew) The EDAS namespace ID.
* `deploy_cluster_code` - (Optional) The deployment cluster code.
* `deploy_cluster_namespace` - (Optional, ForceNew) The deployment cluster namespace.
* `ingress_class_name` - (Optional, ForceNew) The ingress class name.
* `prometheus_enabled` - (Optional, ForceNew) Whether Prometheus monitoring is enabled.
* `sls_enabled` - (Optional, ForceNew) Whether SLS logging is enabled.
* `edas_app_infos` - (Optional) A set of EDAS application information blocks. Each block contains:
  * `edas_namespace` - (Optional) The EDAS namespace.
  * `k8s_cluster_id` - (Optional) The Kubernetes cluster ID.
  * `k8s_namespace` - (Optional) The Kubernetes namespace. Defaults to `default`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance.
* `broker_engine_version` - The broker engine version.
* `access_mode` - The access mode of the instance.
* `tid` - The tenant ID.
* `create_time` - The creation time of the instance.
* `k8s_cluster_id` - The ID of the Kubernetes cluster.
* `edas_app_id` - The EDAS application ID.
* `status` - The status of the instance.
* `broker_latest_engine_version` - The latest broker engine version.
* `shared_instance` - Whether the instance is shared.
* `custom_deploy_config` - Custom deployment configuration (only available for custom deploy mode).

## Import

API Gateway V2 instances can be imported using the instance ID:

```shell
terraform import alibabacloudstack_api_gateway_v2_instance.example <instance_id>
```