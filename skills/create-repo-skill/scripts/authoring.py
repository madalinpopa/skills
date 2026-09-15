"""Structured metadata operations for the local Bash authoring commands."""
# Scaffold-text checking adapts Apache-2.0 source listed in ../NOTICE.
import argparse
import os
from pathlib import Path
import re
import sys
import tempfile

try:
    import yaml
except ImportError:
    sys.exit("authoring: PyYAML is required; use an existing environment or uv run --with pyyaml bash <script> ...")


class UniqueLoader(yaml.SafeLoader):
    pass


def unique_mapping(loader, node, deep=False):
    # Check explicit keys before resolving YAML merge aliases.
    seen = set()
    for key_node, _ in node.value:
        if key_node.tag == "tag:yaml.org,2002:merge":
            key = "<<"
        else:
            key = loader.construct_object(key_node, deep=deep)
        if not isinstance(key, str):
            raise ValueError("YAML mapping keys must be strings")
        if key in seen:
            raise ValueError(f"duplicate YAML key: {key}")
        seen.add(key)
    loader.flatten_mapping(node)
    return yaml.SafeLoader.construct_mapping(loader, node, deep=deep)


UniqueLoader.add_constructor(yaml.resolver.BaseResolver.DEFAULT_MAPPING_TAG, unique_mapping)


def valid_name(name):
    if not isinstance(name, str) or len(name) > 64 or not re.fullmatch(r"[a-z0-9]+(?:-[a-z0-9]+)*", name):
        raise ValueError("name must be 1-64 lowercase letters, digits, and single hyphens")


def valid_description(value):
    if not isinstance(value, str) or not value.strip() or len(value) > 1024:
        raise ValueError("description must contain 1-1024 characters")
    if value.lstrip().startswith("[TODO:"):
        raise ValueError("description contains unfinished scaffold text")


def regular_folder(path):
    path = Path(path).absolute()
    if path.is_symlink() or not path.is_dir():
        raise ValueError(f"expected a regular skill directory: {path}")
    return path


def read_skill(folder):
    path = folder / "SKILL.md"
    if path.is_symlink() or not path.is_file():
        raise ValueError("SKILL.md must be a regular file")
    content = path.read_text()
    match = re.match(r"\A---\r?\n(.*?)\r?\n---(?:\r?\n|$)", content, re.S)
    if not match:
        raise ValueError("SKILL.md requires YAML frontmatter")
    data = yaml.load(match.group(1), Loader=UniqueLoader)
    if not isinstance(data, dict):
        raise ValueError("frontmatter must be a mapping")
    return data, content[match.end():]


def validate_metadata(data, folder):
    valid_name(data.get("name"))
    if data["name"] != folder.name:
        raise ValueError("name must match the skill directory")
    valid_description(data.get("description"))
    if data.get("status", "published") not in ("published", "draft"):
        raise ValueError("status must be published or draft")
    if "tags" in data and (not isinstance(data["tags"], list) or
                           any(not isinstance(tag, str) for tag in data["tags"])):
        raise ValueError("tags must be a list of strings")
    if "x-claude" in data:
        scoped = data["x-claude"]
        if not isinstance(scoped, dict):
            raise ValueError("x-claude must be a mapping")
        forbidden = {"name", "description", "status", "tags", "x-claude"} | set(data)
        collisions = forbidden & set(scoped)
        if collisions:
            raise ValueError("x-claude has reserved or colliding keys: " + ", ".join(sorted(collisions)))


def validate_body(body):
    marker = None
    length = 0
    for line in body.splitlines():
        fence = re.match(r"^[ \t]*(?:(?:[-+*]|\d+[.)])[ \t]+)?(`{3,}|~{3,})(.*)$", line)
        if fence:
            chars, trailing = fence.groups()
            if marker is None:
                marker, length = chars[0], len(chars)
            elif chars[0] == marker and len(chars) >= length and not trailing.strip():
                marker = None
            continue
        if marker is None and re.fullmatch(r"[ ]{0,3}\[TODO:[^\n]*\][ \t]*", line):
            raise ValueError("instructions contain unfinished scaffold text")


class Quoted(str):
    pass


class ValueDumper(yaml.SafeDumper):
    pass


ValueDumper.add_representer(Quoted, lambda dumper, value: dumper.represent_scalar(
    "tag:yaml.org,2002:str", value, style='"'))


def quote_values(value):
    if isinstance(value, str):
        return Quoted(value)
    if isinstance(value, list):
        return [quote_values(item) for item in value]
    if isinstance(value, dict):
        return {key: quote_values(item) for key, item in value.items()}
    return value


def metadata(args):
    folder = regular_folder(args.folder)
    skill, _ = read_skill(folder)
    validate_metadata(skill, folder)
    overrides = {}
    allowed = {"display_name", "short_description", "icon_small", "icon_large", "brand_color", "default_prompt"}
    for item in args.interface:
        key, sep, value = item.partition("=")
        if not sep or key not in allowed or not value.strip():
            raise ValueError(f"invalid interface override: {item}")
        if key in overrides:
            raise ValueError(f"duplicate interface override: {key}")
        if key == "short_description" and not 25 <= len(value) <= 64:
            raise ValueError("short_description must be 25-64 characters")
        if key == "default_prompt" and f"${skill['name']}" not in value:
            raise ValueError("default_prompt must mention the skill as $skill-name")
        if key == "brand_color" and not re.fullmatch(r"#[0-9a-fA-F]{6}", value):
            raise ValueError("brand_color must be a six-digit hex color")
        overrides[key] = value
    directory = folder / "agents"
    path = directory / "openai.yaml"
    if directory.is_symlink() or path.is_symlink():
        raise ValueError("agent metadata path must not be a symlink")
    data = yaml.load(path.read_text(), Loader=UniqueLoader) if path.exists() else {}
    if not isinstance(data, dict) or not isinstance(data.get("interface", {}), dict):
        raise ValueError("agent metadata and interface must be mappings")
    data.setdefault("interface", {}).update(overrides)
    output = yaml.dump(quote_values(data), Dumper=ValueDumper, sort_keys=False, allow_unicode=True)
    directory.mkdir(exist_ok=True)
    # Replace only after parsing and serialization succeed.
    fd, temporary = tempfile.mkstemp(prefix=".openai-", dir=directory)
    try:
        with os.fdopen(fd, "w") as stream:
            stream.write(output)
        if path.exists():
            os.chmod(temporary, path.stat().st_mode & 0o777)
        os.replace(temporary, path)
    finally:
        if os.path.exists(temporary):
            os.unlink(temporary)
    print(f"Updated {path}")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    init = commands.add_parser("frontmatter")
    init.add_argument("name")
    init.add_argument("description")
    init.add_argument("status", choices=("published", "draft"))
    validate = commands.add_parser("validate")
    validate.add_argument("folder")
    ui = commands.add_parser("metadata")
    ui.add_argument("folder")
    ui.add_argument("--interface", action="append", required=True)
    args = parser.parse_args()
    if args.command == "frontmatter":
        valid_name(args.name)
        valid_description(args.description)
        data = {"name": args.name, "description": args.description, "status": args.status}
        print("---\n" + yaml.safe_dump(data, sort_keys=False, allow_unicode=True) + "---")
    elif args.command == "validate":
        folder = regular_folder(args.folder)
        data, body = read_skill(folder)
        validate_metadata(data, folder)
        validate_body(body)
        print(f"Valid skill: {folder.name}")
    else:
        metadata(args)


if __name__ == "__main__":
    try:
        main()
    except (ValueError, OSError, yaml.YAMLError) as error:
        sys.exit(f"authoring: {error}")
