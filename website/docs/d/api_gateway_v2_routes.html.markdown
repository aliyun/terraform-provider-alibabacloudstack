# alibabacloudstack_api_gateway_v2_route

> API Gateway v2 Route Management

## Example Usage

```hcl

variable "name" {
  default = "tf-testacca-13005"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_service" "default" {
  name              = var.name
  description       = var.name
  protocol          = "HTTP"
  upstream_type     = "1"
  load_balance_type = "1"
  gw_instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  service_nodes {
    ip     = "127.0.0.1"
    port   = "80"
    weight = "100"
    enable = "true"
  }
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain      = "${var.name}.com"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol    = "HTTP"
}

resource "alibabacloudstack_api_gateway_v2_route" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  route_name     = var.name
  strip_prefix   = "2"
  order          = "100"
  route_path     = ["/testtc", "/test/aaa"]
  methods        = ["GET", "POST", "PUT", "DELETE"]
  header {
    key   = "header"
    value = "aaaaa"
  }
  cookie {
    key   = "cookie"
    value = "bbbbb"
  }
  query_param {
    key   = "query"
    value = "ccccc"
  }
  domain_ids = ["${alibabacloudstack_api_gateway_v2_domain.default.domain_id}"]
  service_id = alibabacloudstack_api_gateway_v2_service.default.service_id
}

data "alibabacloudstack_api_gateway_v2_routes" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_route.default.gw_instance_id
  name_regex     = "tf-testacca-13005"
  ids            = ["${alibabacloudstack_api_gateway_v2_route.default.id}"]
}

```

## Argument Reference

The following arguments are supported:

### Required

* `domainIds` (List): List of domain IDs, used to specify the domains bound to the route.

* `enableStatus` (Boolean): Whether to enable the route, true means enabled, false means disabled.

* `methods` (List): List of HTTP methods, such as ["GET", "POST"], etc.

* `openStripPrefix` (Boolean): Whether to enable path prefix stripping.

* `order` (Integer): Route priority, the smaller the value, the higher the priority.

* `routePath` (List): List of route paths, such as ["/test", "/api"], etc.

* `serviceIds` (List): List of service IDs, used to specify the backend services of the route.

* `serviceType` (String): Service type, such as "SINGLE" or "MULTI".

* `stripPrefix` (Integer): Path prefix stripping length.

### Force New

* `groupId` (String): Group ID, specifying the group to which the route belongs. (Force new when changed)

* `gwInstanceId` (String): Gateway instance ID, specifying the gateway instance to which the route belongs. (Force new when changed)

* `routeName` (String): Route name, used to identify the route. (Force new when changed)

### Optional

* `cookie` (List): Cookie configuration, used to set cookies in the request.

* `header` (List): Request header configuration, used to set headers in the request.

* `queryParam` (List): Query parameter configuration, used to set query parameters in the request.

* `serviceId` (String): Service ID, used when serviceType is "SINGLE".

## Attributes Reference

The following attributes are exported:

* `id` (String): Resource ID, in the format "{gwInstanceId}:{routeId}".

* `basePath` (String): Base path.

* `create_time` (String): Route creation time.

* `domainVO` (List): Domain information, including domain name, protocol, etc.

* `group_name` (String): Group name.

* `isOpenCross` (Boolean): Whether cross-origin is enabled.

* `isOpenFreeCert` (Boolean): Whether free certificate is enabled.

* `isOpenMock` (Boolean): Whether mock is enabled.

* `isOpenPlugin` (Boolean): Whether plugin is enabled.

* `isOpenResponseCache` (Boolean): Whether response caching is enabled.

* `isOpenSig` (Boolean): Whether signature is enabled.

* `isOpenTimeOut` (Boolean): Whether timeout is enabled.

* `isOpenTraffic` (Boolean): Whether traffic control is enabled.

* `loadBalanceType` (Integer): Load balancing type.

* `loadBalanceTypeName` (String): Load balancing type name.

* `serviceCreateTime` (String): Service creation time.

* `serviceDescription` (String): Service description.

* `serviceName` (String): Service name.

* `upstreamType` (Integer): Backend type.

* `upstreamTypeName` (String): Backend type name.