---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_service"
sidebar_current: "docs-alibabacloudstack-resource-edas-k8s-service"
description: |-
  Provides an EDAS K8s cluster resource.
---

# alibabacloudstack\_edas\_k8s\_service

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_edas_k8s_application" "default" {
  // package type is Image / FatJar / War
  package_type            = "Image"
  application_name        = "DemoApplication"
  application_description = "This is description of application"
  cluster_id              = var.cluster_id
  replicas                = 2

  // set 'image_url' and 'repo_id' when package_type is 'image'
  image_url = "registry-vpc.cn-beijing.aliyuncs.com/edas-demo-image/consumer:1.0"

  // set 'package_url','package_version' and 'jdk' when package_type is not 'image'
  package_url     = var.package_url
  package_version = var.package_version
  jdk             = var.jdk

  // set 'web_container' and 'edas_container' when package_type is 'war'
  web_container          = var.web_container
  edas_container_version = var.edas_container_version

  internet_target_port  = var.internet_target_port
  internet_slb_port     = var.internet_slb_port
  internet_slb_protocol = var.internet_slb_protocol
  internet_slb_id       = var.internet_slb_id
  limit_mem             = 2048
  requests_mem          = 0
  requests_m_cpu        = 0
  limit_m_cpu           = 4000
  command               = var.command
  command_args          = var.command_args
  envs                  = var.envs
  pre_stop              = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  post_start            = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  liveness              = var.liveness
  readiness             = var.readiness
  nas_id                = var.nas_id
  mount_descs           = var.mount_descs
  local_volume          = var.local_volume
  namespace             = "default"
  logical_region_id     = "cn-beijing"
}

resource "alibabacloudstack_edas_k8s_service" "default" {
  app_id = alibabacloudstack_edas_k8s_application.default.id
  name = "tf-testAccEdasK8sService"
  type = "ClusterIP"
  external_traffic_policy = "Local"
  service_ports {
    protocol = "TCP"
    service_port = 80
    target_port = 8080
  }
}

```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The id of the Edas k8s application
* `name` - (Required, ForceNew) The name of the Edas k8s service
* `type` - (Required, ForceNew) Edas K8s service types, Valid values are `ClusterIP`, `NodePort`, `LoadBalancer`
* `service_ports` - (Optional) K8s Service port mapping table, which needs to conform to the JsonArray format. The supported parameters are as follows:
  * `protocol` - (Required) the service protocol, supporting TCP and UDP.
  * `service_port` - (Required) the frontend service port, with a value range of 1~65535.
  * `target_port` - (Required) the backend container port, with a value range of 1~65535.
* `external_traffic_policy` - (Optional) Set the external traffic management policy. Valid values are `Local`, `Cluster`, Default to `Local`.
* `annotations` - (Optional) The annotations map of the application
* `labels` - (Optional) The labels map of the application
## Attributes Reference

The following attributes are exported:

* `cluster_ip` - The cluster ip of the application

## Import

EDAS k8s application can be imported as below, e.g.

```
$ terraform import alibabacloudstack_edas_k8s_service.new_k8s_service "app_id:name"
```
