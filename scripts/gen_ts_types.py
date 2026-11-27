#!/usr/bin/env python3
"""从 server/_proto/onij/*.proto 生成前端 TypeScript 类型。

后端 Hertz + protoc 输出的是 snake_case JSON（如 `json:"artist_types"`），
前端类型必须与之一致，否则字段读不到（静默 undefined，是线上真实 bug 来源）。

用法：
    python scripts/gen_ts_types.py          # 生成
    python scripts/gen_ts_types.py --check  # 只校验，漂移则退出码 1

产物：
    src/api/types/enums.ts   ← all_enum.proto
    src/api/types/models.ts  ← all_model.proto（领域模型）
    src/api/types/api.ts     ← all.proto（请求 / 响应）
"""
from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
PROTO_DIR = ROOT / "server" / "_proto" / "onij"
OUT_DIR = ROOT / "src" / "api" / "types"

SCALAR = {
    "double": "number",
    "float": "number",
    "int32": "number",
    "int64": "number",
    "uint32": "number",
    "uint64": "number",
    "sint32": "number",
    "sint64": "number",
    "fixed32": "number",
    "fixed64": "number",
    "sfixed32": "number",
    "sfixed64": "number",
    "bool": "boolean",
    "string": "string",
    "bytes": "string",
}

Field = tuple[str, str, bool, bool]  # name, protoType, repeated, optional


def header(source: str) -> str:
    return (
        "// 此文件由 scripts/gen_ts_types.py 自动生成，请勿手工修改。\n"
        f"// 源：server/_proto/onij/{source}\n"
        "// 重新生成：pnpm gen:types\n"
    )


def strip_comments(text: str) -> str:
    text = re.sub(r"/\*.*?\*/", "", text, flags=re.S)
    return re.sub(r"//[^\n]*", "", text)


def read_proto(name: str) -> str:
    return strip_comments((PROTO_DIR / name).read_text(encoding="utf-8"))


def match_block(text: str, open_idx: int) -> tuple[str, int]:
    """从 open_idx 处的 `{` 开始，返回 (块内文本, 右括号之后的位置)。"""
    depth = 0
    for i in range(open_idx, len(text)):
        if text[i] == "{":
            depth += 1
        elif text[i] == "}":
            depth -= 1
            if depth == 0:
                return text[open_idx + 1 : i], i + 1
    raise ValueError("未闭合的块")


def parse_enums(text: str) -> dict[str, list[tuple[str, int]]]:
    out: dict[str, list[tuple[str, int]]] = {}
    for m in re.finditer(r"\benum\s+(\w+)\s*\{", text):
        body, _ = match_block(text, m.end() - 1)
        out[m.group(1)] = [
            (vm.group(1), int(vm.group(2)))
            for vm in re.finditer(r"(\w+)\s*=\s*(-?\d+)\s*;", body)
        ]
    return out


def parse_fields(body: str, rename: dict[str, str] | None = None) -> list[Field]:
    rename = rename or {}
    fields: list[Field] = []
    for fm in re.finditer(
        r"(?:(repeated|optional)\s+)?([\w.]+)\s+(\w+)\s*=\s*\d+\s*[;\[]", body
    ):
        label, ftype, fname = fm.group(1), fm.group(2), fm.group(3)
        if ftype == "map":  # 未使用；避免生成错误类型
            continue
        fields.append((fname, rename.get(ftype, ftype), label == "repeated", label == "optional"))
    return fields


def parse_messages(text: str) -> dict[str, list[Field]]:
    """解析全部 message，嵌套的记为 Parent_Child，并重写父内对子类型的引用。"""
    messages: dict[str, list[Field]] = {}

    def walk(scope: str, body: str) -> None:
        nested: list[tuple[int, int, str, str]] = []
        for nm in re.finditer(r"message\s+(\w+)\s*\{", body):
            inner, end = match_block(body, nm.end() - 1)
            nested.append((nm.start(), end, nm.group(1), inner))

        flat, cursor = [], 0
        for start, end, _, _ in nested:
            flat.append(body[cursor:start])
            cursor = end
        flat.append(body[cursor:])

        # 父内引用子类型时用 Parent_Child 全名
        rename = {child: f"{scope}_{child}" for _, _, child, _ in nested}
        messages[scope] = parse_fields("".join(flat), rename)

        for _, _, child, child_body in nested:
            walk(f"{scope}_{child}", child_body)

    consumed: list[tuple[int, int]] = []
    for m in re.finditer(r"message\s+(\w+)\s*\{", text):
        if any(s <= m.start() < e for s, e in consumed):
            continue
        body, end = match_block(text, m.end() - 1)
        consumed.append((m.start(), end))
        walk(m.group(1), body)

    return messages


def ts_type(proto_type: str, repeated: bool, known: set[str]) -> str:
    base = SCALAR.get(proto_type) or (proto_type if proto_type in known else "any")
    return f"{base}[]" if repeated else base


def emit_enums(enums: dict[str, list[tuple[str, int]]]) -> str:
    out = [header("all_enum.proto")]
    for name, values in enums.items():
        out.append(f"export enum {name} {{")
        out.extend(f"  {vn} = {v}," for vn, v in values)
        out.append("}")
        out.append("")
    return "\n".join(out).rstrip() + "\n"


def emit_models(
    models: dict[str, list[Field]], enums: dict, known: set[str]
) -> str:
    out = [header("all_model.proto"), ""]
    # 模型引用的枚举必须可见
    used = {ft for fields in models.values() for _, ft, _, _ in fields if ft in enums}
    out.extend(f"import type {{ {e} }} from './enums'" for e in sorted(used))
    out.append("")
    out.extend(f"export type {{ {e} }} from './enums'" for e in sorted(used))
    out.append("")
    for name, fields in models.items():
        out.append(f"export interface {name} {{")
        for fn, ft, rp, op in fields:
            opt = "?" if op else ""
            out.append(f"  {fn}{opt}: {ts_type(ft, rp, known)}")
        out.append("}")
        out.append("")
    return "\n".join(out).rstrip() + "\n"


def emit_api(
    apis: dict[str, list[Field]], models: dict, enums: dict, known: set[str]
) -> str:
    out = [header("all.proto"), ""]
    out.append("import type {")
    out.extend(f"  {m}," for m in sorted(models))
    out.append("} from './models'")
    used = {ft for fields in apis.values() for _, ft, _, _ in fields if ft in enums}
    out.extend(f"import type {{ {e} }} from './enums'" for e in sorted(used))
    out.append("")
    for name, fields in apis.items():
        out.append(f"export interface {name} {{")
        for fn, ft, rp, op in fields:
            opt = "?" if op else ""
            out.append(f"  {fn}{opt}: {ts_type(ft, rp, known)}")
        out.append("}")
        out.append("")
    return "\n".join(out).rstrip() + "\n"


def generate() -> dict[Path, str]:
    enums = parse_enums(read_proto("all_enum.proto"))
    models = parse_messages(read_proto("all_model.proto"))
    apis = parse_messages(read_proto("all.proto"))
    known = set(enums) | set(models) | set(apis)

    apis = {k: v for k, v in apis.items() if k not in models}

    return {
        OUT_DIR / "enums.ts": emit_enums(enums),
        OUT_DIR / "models.ts": emit_models(models, enums, known),
        OUT_DIR / "api.ts": emit_api(apis, models, enums, known),
    }


def normalize(content: str) -> str:
    """统一空行与结尾换行，避免生成结果与磁盘内容因空白差异误报漂移。"""
    lines = [ln for ln in content.splitlines() if ln.strip() or ln == ""]
    return re.sub(r"\n{3,}", "\n\n", "\n".join(lines)) + "\n"


def main() -> int:
    check = "--check" in sys.argv
    outputs = generate()

    if check:
        drifted = []
        for path, content in outputs.items():
            current = path.read_text(encoding="utf-8") if path.exists() else ""
            if normalize(current) != normalize(content):
                drifted.append(path)
        if drifted:
            print("类型与 proto 不一致，请运行 pnpm gen:types：", file=sys.stderr)
            for p in drifted:
                print(f"  {p.relative_to(ROOT)}", file=sys.stderr)
            return 1
        print("类型与 proto 一致")
        return 0

    for path, content in outputs.items():
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(normalize(content), encoding="utf-8")
        print(f"generated {path.relative_to(ROOT)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
