---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_service_roles"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ram-service-roles"
description: |-
    查询RAM服务角色列表。
---

# alibabacloudstack_ascm_ram_service_roles  

根据指定过滤条件列出当前凭证权限可以访问的RAM角色列表。  

## 示例用法  

```
data "alibabacloudstack_ascm_ram_service_roles" "role" {  
  product = "ECS"  
}  
output "role" {  
  value = data.alibabacloudstack_ascm_ram_service_roles.role.*  
}  
```  

## 参数说明  

支持以下参数：  

* `ids` - (可选) RAM 角色 ID 列表。用于筛选特定的角色。  
* `product` - (可选) 按其产品过滤结果的正则表达式字符串。有效值为 "ECS"，表示只返回与 ECS 产品相关的角色。  
* `description` - (可选) 关于 RAM 角色的描述。可以通过此字段进一步筛选符合特定描述的角色。  

## 属性参考

除了上述列出的参数外，还导出以下属性：  

* `roles` - 角色列表。每个元素包含以下属性：  
    * `id` - 角色的唯一标识符，用于区分不同的角色。  
    * `name` - 角色名称，通常用于标识该角色的功能或用途。  
    * `description` - 角色的详细描述信息，帮助用户理解角色的具体功能和权限范围。  
    * `role_type` - 角色类型，表示该角色的分类或用途（例如系统预定义角色、自定义角色等）。  
    * `product` - 角色所属的产品类型，例如 "ECS" 表示该角色与弹性计算服务相关。  
    * `organization_name` - 组织名称，表示该角色所属的组织或部门。  
    * `aliyun_user_id` - 阿里云用户 ID，表示创建或管理该角色的阿里云账户。  
    * `computed_property_example` - (计算属性) 根据提供的模式自动生成的计算属性示例，具体含义取决于实际实现逻辑。  

**补充说明：**  
- `ids` 参数允许用户通过提供一个角色 ID 列表来精确匹配需要查询的角色。  
- `product` 参数用于限定查询范围，确保返回的角色仅与指定产品相关联。  
- `description` 参数可以帮助用户进一步缩小查询范围，仅返回描述信息符合特定条件的角色。  

在 `roles` 属性中：  
- `id` 是每个角色的唯一标识符，可用于后续操作中的角色引用。  
- `name` 提供了角色的直观名称，便于用户识别其功能或用途。  
- `description` 字段详细描述了角色的职责和权限范围，有助于用户更好地理解角色的作用。  
- `role_type` 明确了角色的类型（如系统预定义或用户自定义），方便用户分类管理。  
- `product` 字段指明了角色所属的产品类别，例如 ECS、RDS 等。  
- `organization_name` 和 `aliyun_user_id` 分别表示角色所属的组织和创建者信息，有助于多租户环境下的权限管理。  
- `computed_property_example` 是一个计算属性，其值由系统根据内部逻辑动态生成，具体含义需结合实际场景理解。