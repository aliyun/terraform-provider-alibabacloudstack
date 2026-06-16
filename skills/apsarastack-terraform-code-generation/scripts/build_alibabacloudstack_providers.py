#!/usr/bin/env python3
"""Build references/alibabacloudstack-providers.md from the terraform-provider-alibabacloudstack repo.

Clones (or reuses) the provider repo, walks website/docs/{r,d}/*.html.markdown,
extracts subcategory + deprecation status, and writes a subcategory-grouped
markdown catalog to references/alibabacloudstack-providers.md.

Usage:
    scripts/build_alibabacloudstack_providers.py                 # clone or pull, regenerate
    scripts/build_alibabacloudstack_providers.py --no-refresh    # reuse existing clone as-is
    scripts/build_alibabacloudstack_providers.py --repo PATH     # use an existing checkout
    scripts/build_alibabacloudstack_providers.py --out PATH      # override output path

The SKILL's Step 4.1 reads the generated file. Re-run this script when
`aliyun/terraform-provider-alibabacloudstack` publishes a new release.
"""
from __future__ import annotations

import argparse
import re
import subprocess
import sys
from collections import defaultdict
from pathlib import Path

PROVIDER_REPO = "https://github.com/aliyun/terraform-provider-alibabacloudstack.git"
DEFAULT_CLONE = "/tmp/terraform-provider-alibabacloudstack"
BLOB_PREFIX = "https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs"

FRONTMATTER_RE = re.compile(r"\A---\n(.*?)\n---\n(.*)", re.DOTALL)
SUBCAT_RE = re.compile(r'^subcategory:\s*["\']?([^"\'\n]+?)["\']?\s*$', re.MULTILINE)
DEPRECATED_MARKER_RE = re.compile(
    r"->\s*\*\*DEPRECAT[A-Z]+(?:\s+NOTICE)?\*?[:\s\*]",
    re.IGNORECASE,
)
# Replacement extraction patterns
LINK_TEXT_RE = re.compile(r"\[(alibabacloudstack_[a-z0-9_]+)\]")
REGISTRY_LINK_SLUG_RE = re.compile(
    r"\]\([^)]*/(?:resources|data-sources|r|d)/([a-z0-9_]+)"
)
BACKTICK_REF_RE = re.compile(r"`(alibabacloudstack_[a-z0-9_]+)`")
BARE_REF_RE = re.compile(r"\b(alibabacloudstack_[a-z0-9_]+)\b")

def ensure_repo(clone_dir: Path) -> None:
    if (clone_dir / ".git").exists():
        print(f"[info] refreshing existing clone at {clone_dir}", file=sys.stderr)
        subprocess.run(
            ["git", "-C", str(clone_dir), "pull", "--ff-only"],
            check=False,
            timeout=300,
        )
        return
    clone_dir.parent.mkdir(parents=True, exist_ok=True)
    print(f"[info] cloning provider repo to {clone_dir}", file=sys.stderr)
    subprocess.run(
        [
            "git", "clone", "--depth", "1", "--filter=blob:none", "--sparse",
            PROVIDER_REPO, str(clone_dir),
        ],
        check=True,
        timeout=600,
    )
    subprocess.run(
        ["git", "-C", str(clone_dir), "sparse-checkout", "set", "website/docs"],
        check=True,
        timeout=60,
    )

def parse_doc(path: Path, kind: str) -> dict:
    """kind is 'r' (resource) or 'd' (data source)."""
    text = path.read_text(encoding="utf-8", errors="replace")
    subcategory = "Other"
    body = text
    m = FRONTMATTER_RE.match(text)
    if m:
        body = m.group(2)
        sm = SUBCAT_RE.search(m.group(1))
        if sm and sm.group(1).strip():
            subcategory = sm.group(1).strip()

    name_stem = path.name[: -len(".html.markdown")]
    name = f"alibabacloudstack_{name_stem}"

    deprecated = False
    replacement: str | None = None
    dm = DEPRECATED_MARKER_RE.search(body)
    if dm:
        deprecated = True
        start = dm.end()
        stops = [body.find("\n## ", start), body.find("\n\n-> **", start)]
        stops = [s for s in stops if s != -1]
        region = body[start:min(stops)] if stops else body[start:]

        def _first_non_self(iterator, transform=lambda m: m.group(1)):
            for m in iterator:
                candidate = transform(m)
                if not candidate.startswith("alibabacloudstack_"):
                    candidate = f"alibabacloudstack_{candidate}"
                if candidate != name:
                    return candidate
            return None

        replacement = (
            _first_non_self(LINK_TEXT_RE.finditer(region))
            or _first_non_self(REGISTRY_LINK_SLUG_RE.finditer(region))
            or _first_non_self(BACKTICK_REF_RE.finditer(region))
            or _first_non_self(BARE_REF_RE.finditer(region))
        )

    return {
        "name": name,
        "subcategory": subcategory,
        "deprecated": deprecated,
        "replacement": replacement,
        "kind": "resource" if kind == "r" else "data source",
        "url": f"{BLOB_PREFIX}/{kind}/{path.name}",
    }

def emit_markdown(entries: list[dict], out_path: Path) -> None:
    by_cat: dict[str, list[dict]] = defaultdict(list)
    for entry in entries:
        by_cat[entry["subcategory"]].append(entry)
    for cat_entries in by_cat.values():
        cat_entries.sort(key=lambda e: (e["kind"], e["name"]))

    total = len(entries)
    res_count = sum(1 for e in entries if e["kind"] == "resource")
    data_count = total - res_count
    dep_count = sum(1 for e in entries if e["deprecated"])

    lines: list[str] = [
        "# Alibaba Cloud ApsaraStack Terraform provider catalog",
        "",
        f"Total entries: **{total}** "
        f"(resources: {res_count}, data sources: {data_count}; "
        f"deprecated: {dep_count}).",
        "",
        "Built from `aliyun/terraform-provider-alibabacloudstack` provider documentation.",
        "Re-run `scripts/build_alibabacloudstack_providers.py` to refresh when provider updates.",
        "",
        "Columns — **type** (resource / data source), **name** (`alibabacloudstack_*`), ",
        "**status** (empty = supported; `⚠️ 弃用 → alibabacloudstack_X` = deprecated, use X), ",
        "**doc** (GitHub source, used by Step 4.2 WebFetch).",
        "",
    ]

    for cat in sorted(by_cat):
        lines.append(f"## {cat}")
        lines.append("")
        lines.append("| type | name | status | doc |")
        lines.append("| --- | --- | --- | --- |")
        for entry in by_cat[cat]:
            status = ""
            if entry["deprecated"]:
                status = (
                    f"⚠️ 弃用 → `{entry['replacement']}`"
                    if entry["replacement"]
                    else "⚠️ 弃用"
                )
            lines.append(
                f"| {entry['kind']} | `{entry['name']}` | {status} "
                f"| [doc]({entry['url']}) |"
            )
        lines.append("")

    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text("\n".join(lines), encoding="utf-8")

    dep = sum(1 for e in entries if e["deprecated"])
    print(
        f"[ok] wrote {out_path} ({len(entries)} entries, {dep} deprecated)",
        file=sys.stderr,
    )

def main() -> None:
    parser = argparse.ArgumentParser(
        description=__doc__,
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--repo",
        default=DEFAULT_CLONE,
        help=f"Local clone path (default: {DEFAULT_CLONE}).",
    )
    parser.add_argument(
        "--out",
        default=None,
        help="Output markdown path (default: <skill_root>/references/alibabacloudstack-providers.md).",
    )
    parser.add_argument(
        "--no-refresh",
        action="store_true",
        help="Skip git clone/pull; use the existing repo as-is.",
    )
    args = parser.parse_args()

    repo = Path(args.repo).expanduser().resolve()
    if not args.no_refresh:
        ensure_repo(repo)

    docs_r = repo / "website" / "docs" / "r"
    docs_d = repo / "website" / "docs" / "d"
    if not docs_r.is_dir() or not docs_d.is_dir():
        print(
            f"[error] docs dirs missing under {repo}; expected website/docs/{{r,d}}. "
            f"Re-clone with: rm -rf {repo} && python scripts/build_alibabacloudstack_providers.py",
            file=sys.stderr,
        )
        sys.exit(1)

    entries: list[dict] = []
    for md in sorted(docs_r.glob("*.html.markdown")):
        entries.append(parse_doc(md, "r"))
    for md in sorted(docs_d.glob("*.html.markdown")):
        entries.append(parse_doc(md, "d"))

    skill_root = Path(__file__).resolve().parent.parent
    out_path = (
        Path(args.out).expanduser().resolve()
        if args.out
        else skill_root / "references" / "alibabacloudstack-providers.md"
    )
    emit_markdown(entries, out_path)

    dep = sum(1 for e in entries if e["deprecated"])
    print(
        f"[ok] wrote {out_path} ({len(entries)} entries, {dep} deprecated)",
        file=sys.stderr,
    )

if __name__ == "__main__":
    main()
