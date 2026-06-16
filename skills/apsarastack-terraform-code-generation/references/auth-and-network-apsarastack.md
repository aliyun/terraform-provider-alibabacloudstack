# ApsaraStack Auth and network reference

Practical details for two environmental questions: **how the
alibabacloudstack provider finds credentials** and **what to do when
`terraform init` can't reach the upstream registry**.

## 1. Credential resolution

The alibabacloudstack provider requires explicit credentials. Two mechanisms:

| # | Mechanism | Detection signal |
| --- | --- | --- |
| 1 | Env AK/SK | `ALIBABACLOUDSTACK_ACCESS_KEY` + `ALIBABACLOUDSTACK_SECRET_KEY` both set. |
| 2 | Static in HCL | `access_key` / `secret_key` in `provider "alibabacloudstack"` block — **SKILL Hard Rule §1 forbids emitting this**. |

### Step 8 probe (non-reading)

Inside SKILL Step 8's pre-flight check, probe **presence** of credentials
without reading any value:

```bash
(
  [[ -n "${ALIBABACLOUDSTACK_ACCESS_KEY:-}" ]] && [[ -n "${ALIBABACLOUDSTACK_SECRET_KEY:-}" ]] && echo "ready:env-ak-sk"
) | head -1
```

If no line prints → `NO_CREDENTIALS`. Tell the user to set environment
variables or configure STS authentication via their ApsaraStack console.

## 2. `terraform init` network acceleration (Alibaba Cloud mirror)

### When it's needed

`terraform init` in China-mainland networks often can't reach
`registry.terraform.io`. Signatures in the init output that indicate a
network problem (not a config problem):

- `connection refused`
- `TLS handshake timeout`
- `network is unreachable`
- `no such host`
- `context deadline exceeded` fetching `registry.terraform.io`

### Configuration

**File location**:
- Linux / macOS: `~/.terraformrc`
- Windows: `%APPDATA%/terraform.rc`
- Custom: point `TF_CLI_CONFIG_FILE` at any `*.tfrc`.

**Content** (paste verbatim):

```hcl
provider_installation {
  network_mirror {
    url     = "https://mirrors.aliyun.com/terraform/"
    include = ["registry.terraform.io/aliyun/alibabacloudstack"]
  }
  direct {
    exclude = ["registry.terraform.io/aliyun/alibabacloudstack"]
  }
}
```

## 3. ApsaraStack-specific endpoints

The alibabacloudstack provider requires these ApsaraStack-specific settings:

| Parameter | Example | Description |
| --- | --- | --- |
| `popgw_domain` | `popgw.cn-hangzhou-env01-d01.intra.env01.shuguang.com` | POP gateway endpoint (required for V3.16+) |
| `protocol` | `HTTP` | Typically HTTP for internal networks |
| `insecure` | `true` | Usually true for self-signed certificates |
| `region` | `cn-hangzhou-env01-d01` | ApsaraStack region format |
| `resource_group_set_name` | `ResourceSet(wzw)` | Optional resource set name |

All MUST be configured via variables, never hardcoded.
