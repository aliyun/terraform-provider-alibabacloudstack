# ApsaraStack product-specific resource patterns

Positive patterns — "when the user asks X, use these attributes or
multi-resource idioms". This complements:

- `alibabacloudstack-providers.md` (catalog: does this resource exist?)
- `deprecated-fields-apsarastack.md` (field-level renames / splits)

Entries here are **product-specific conventions** that the provider doc
technically documents but does not emphasize.

---

## RDS cross-AZ primary/secondary HA

**Trigger phrases**: "High Availability / HA / primary-secondary architecture / multi-AZ / cross-AZ / primary-secondary / master-slave" applied to `alibabacloudstack_db_instance`.

**Non-obvious requirement**: the provider doc lists `zone_id_slave_a` as
*optional*. Agents often set `category = "HighAvailability"` alone and
assume ApsaraStack places the standby automatically in a different AZ. It does
not — without `zone_id_slave_a`, primary and standby land in the same AZ.

**Required attributes on `alibabacloudstack_db_instance`**:

| Attribute | Value | Why |
| --- | --- | --- |
| `category` | `"HighAvailability"` | Switches the edition. |
| `zone_id` | `data.alibabacloudstack_db_zones.<n>.zones[0].id` | Primary AZ. |
| `zone_id_slave_a` | `data.alibabacloudstack_db_zones.<n>.zones[1].id` | Secondary AZ. MUST differ from `zone_id`. |

**Sketch**:

```hcl
data "alibabacloudstack_db_zones" "mysql_ha" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
}

resource "alibabacloudstack_db_instance" "this" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  instance_type            = var.rds_instance_type
  instance_storage         = 100

  zone_id          = data.alibabacloudstack_db_zones.mysql_ha.zones[0].id
  zone_id_slave_a  = data.alibabacloudstack_db_zones.mysql_ha.zones[1].id

  # ... vswitch_id, security_group_ids, security_ips, etc.
}
```

---

## OSS lifecycle — current vs noncurrent versions

**Trigger phrases**: "old versions / historical versions / noncurrent / old object versions / N days later (transition to IA|Archive)" applied to an
`alibabacloudstack_oss_bucket` with a lifecycle rule.

**Non-obvious requirement**: the `lifecycle_rule` block has TWO
transition sub-blocks with different targets.

| Sub-block | Targets | When to use |
| --- | --- | --- |
| `transition { days = N, storage_class = … }` | *Current* object version | User says "files transition to IA after N days" (current objects) |
| `noncurrent_version_transition { days = N, storage_class = … }` | *Older* / noncurrent versions | User says "old versions / historical versions / noncurrent ..." |

Versioning MUST be enabled on the bucket (via
`alibabacloudstack_oss_bucket_versioning`) for
`noncurrent_version_transition` to have any effect.

**Sketch**:

```hcl
resource "alibabacloudstack_oss_bucket" "this" {
  bucket = var.bucket_name
  lifecycle_rule {
    id      = "archive-old-versions"
    prefix  = ""
    enabled = true

    # user said "old versions transition to IA after 90 days" → use noncurrent_version_transition
    noncurrent_version_transition {
      days          = 90
      storage_class = "IA"
    }
  }
}

resource "alibabacloudstack_oss_bucket_versioning" "this" {
  bucket = alibabacloudstack_oss_bucket.this.id
  status = "Enabled"
}
```

---

## ECS with data disk

**Trigger phrases**: "data disk / attach disk" applied to
`alibabacloudstack_instance`.

**Required pattern**:

```hcl
resource "alibabacloudstack_disk" "data" {
  disk_name = var.data_disk_name
  size      = var.data_disk_size
  category  = "cloud_efficiency"
  zone_id   = data.alibabacloudstack_zones.main.zones[0].id
}

resource "alibabacloudstack_disk_attachment" "data_attach" {
  disk_id   = alibabacloudstack_disk.data.id
  instance_id = alibabacloudstack_instance.main.id
}
```

---

## SLB with backend servers

**Trigger phrases**: "load balancer / backend servers" applied to
`alibabacloudstack_slb`.

**Required pattern**:

```hcl
resource "alibabacloudstack_slb" "main" {
  load_balancer_name = var.slb_name
  vpc_id             = alibabacloudstack_vpc.main.id
  address_type       = "intranet"
}

resource "alibabacloudstack_slb_backend_server_attachment" "main" {
  load_balancer_id = alibabacloudstack_slb.main.id

  backend_servers {
    server_id = alibabacloudstack_instance.web1.id
    weight    = 100
  }

  backend_servers {
    server_id = alibabacloudstack_instance.web2.id
    weight    = 100
  }
}
```
