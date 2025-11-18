# alibabacloudstack_api_gateway_v2_route

> api网关 v2 版本 路由管理

## 示例用法

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

## 参数说明

以下参数支持入参配置：

### 必填参数

* `domainIds` (列表)：域名ID列表，用于指定路由绑定的域名。

* `enableStatus` (布尔)：是否启用路由，true表示启用，false表示禁用。

* `methods` (列表)：HTTP方法列表，如["GET", "POST"]等。

* `openStripPrefix` (布尔)：是否开启路径前缀截断。

* `order` (整数)：路由优先级，数值越小优先级越高。

* `routePath` (列表)：路由路径列表，如["/test", "/api"]等。

* `serviceIds` (列表)：服务ID列表，用于指定路由后端服务。

* `serviceType` (字符串)：服务类型，如"SINGLE"或"MULTI"。

* `stripPrefix` (整数)：路径前缀截断长度。

### 变更时重建参数

* `groupId` (字符串)：分组ID，指定路由所属的分组。(变更时重建)

* `gwInstanceId` (字符串)：网关实例ID，指定路由所属的网关实例。(变更时重建)

* `routeName` (字符串)：路由名称，用于标识路由。(变更时重建)

### 可选参数

* `cookie` (列表)：Cookie配置，用于设置请求中的Cookie。

* `header` (列表)：请求头配置，用于设置请求中的Header。

* `queryParam` (列表)：查询参数配置，用于设置请求中的查询参数。

* `serviceId` (字符串)：服务ID，当serviceType为"SINGLE"时使用。

## 属性说明

以下属性被导出：

* `id` (字符串)：资源ID，格式为"{gwInstanceId}:{routeId}"。

* `basePath` (字符串)：基础路径。

* `create_time` (字符串)：路由创建时间。

* `domainVO` (列表)：域名信息，包含域名、协议等信息。

* `group_name` (字符串)：分组名称。

* `isOpenCross` (布尔)：是否开启跨域。

* `isOpenFreeCert` (布尔)：是否开启免费证书。

* `isOpenMock` (布尔)：是否开启Mock。

* `isOpenPlugin` (布尔)：是否开启插件。

* `isOpenResponseCache` (布尔)：是否开启响应缓存。

* `isOpenSig` (布尔)：是否开启签名。

* `isOpenTimeOut` (布尔)：是否开启超时。

* `isOpenTraffic` (布尔)：是否开启流量控制。

* `loadBalanceType` (整数)：负载均衡类型。

* `loadBalanceTypeName` (字符串)：负载均衡类型名称。

* `serviceCreateTime` (字符串)：服务创建时间。

* `serviceDescription` (字符串)：服务描述。

* `serviceName` (字符串)：服务名称。

* `upstreamType` (整数)：后端类型。

* `upstreamTypeName` (字符串)：后端类型名称。