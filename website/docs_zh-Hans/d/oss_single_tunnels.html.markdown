---
page_title: "alibabacloudstack_oss_single_tunnels"
subcategory: "对象存储 OSS"
description: |-
  提供阿里云专有云中 OSS 单隧道的列表。
---

# alibabacloudstack_oss_single_tunnels

提供阿里云专有云中 OSS 单隧道的列表。这些隧道代表用于跨 VPC 连接 OSS 服务的虚拟 IP 地址。

## 用法示例

```hcl
# 声明数据源
data "alibabacloudstack_oss_single_tunnels" "tunnels" {
  ids        = ["cluster1:vpc-abc123:vip1", "cluster2:vpc-def456:vip2"]
  name_regex = "^test-.*"
}

output "first_tunnel_id" {
  value = data.alibabacloudstack_oss_single_tunnels.tunnels.tunnels.0.id
}

output "all_tunnel_ids" {
  value = [for tunnel in data.alibabacloudstack_oss_single_tunnels.tunnels.tunnels : tunnel.id]
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 隧道 ID 列表。每个 ID 的格式为 `cluster:vpc_id:vip`。
* `name_regex` - (可选) 用于按标签过滤结果的正则表达式字符串。仅返回标签与此正则表达式匹配的隧道。

## 属性参考

除了上述参数外，还导出以下属性：

* `tunnels` - 隧道列表。每个元素包含以下属性：
  * `id` - 隧道的 ID，格式为 `cluster:vpc_id:vip`。
  * `cluster` - 隧道所属的集群。
  * `label` - 隧道的标签。
  * `vip` - 隧道的虚拟 IP 地址。
  * `vpc_id` - 与隧道关联的 VPC 的 ID。
  * `shared` - 隧道是否共享（值为 0 或 1）。
