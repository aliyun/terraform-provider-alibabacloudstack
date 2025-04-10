---
subcategory: "AnalyticDB for MySQL (ADB)"  
layout: "alibabacloudstack"  
page_title: "Alibabacloudstack: alibabacloudstack_adb_zones"  
sidebar_current: "docs-alibabacloudstack-datasource-adb-zones"  
description: |-  
    查询ADB可用区  
---  

# alibabacloudstack_adb_zones  

根据指定过滤条件列出当前凭证权限可以访问的ADB可用区列表。  

## 示例用法  

```
# 声明数据源
data "alibabacloudstack_adb_zones" "zones_ids" {}
```  

## 参数说明  

支持以下参数：  

* `multi` - (可选) 指示这些可用区是否可以用于多AZ配置。默认值为`false`。当设置为`true`时，返回的可用区将支持多可用区部署，通常用于启动ADB实例。  

## 属性说明  

除了上述参数外，还导出以下属性：  

* `ids` - 区域ID列表。该列表包含所有符合条件的可用区ID。  
* `zones` - 可用区列表。每个元素包含以下属性：  
  * `id` - 区域的ID，唯一标识一个可用区。  
  * `multi_zone_ids` - 多区域中的区域ID列表。此属性列出支持多可用区部署的区域ID集合，仅在`multi`参数设置为`true`时有意义。  
