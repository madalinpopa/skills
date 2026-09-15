"""Exercise the local authoring commands through their Bash entrypoints."""
from pathlib import Path
import subprocess
import tempfile
import unittest

import yaml

SCRIPTS = Path(__file__).resolve().parents[1] / "scripts"


class AuthoringCommands(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="repo-skill-test-")
        self.addCleanup(self.temp.cleanup)
        self.repo = Path(self.temp.name) / "repository with spaces"
        (self.repo / "skills").mkdir(parents=True)
        (self.repo / "docs").mkdir()
        (self.repo / "AGENTS.md").write_text("# Instructions\n")
        (self.repo / "docs/SPEC.md").write_text("# Specification\n")

    def run_command(self, script, *args):
        # A missing entrypoint is missing feature behavior, not a skipped test.
        self.assertTrue((SCRIPTS / script).is_file(), f"Authoring command not implemented: {script}")
        return subprocess.run(["bash", str(SCRIPTS / script), *map(str, args)],
                              capture_output=True, text=True)

    def init(self, name="create-example", *args):
        return self.run_command("init-skill.sh", name, "--repo", self.repo,
                                "--description", 'Creates examples: safely, with "quotes".', *args)

    def write_skill(self, name="create-example", extra="", body="Complete the requested example.\n"):
        folder = self.repo / "skills" / name
        folder.mkdir()
        (folder / "SKILL.md").write_text(
            f"---\nname: {name}\ndescription: Creates examples.\nstatus: published\n{extra}---\n\n{body}")
        return folder

    def metadata(self, folder):
        return yaml.safe_load((folder / "SKILL.md").read_text().split("---", 2)[1])

    def test_init_published_and_requested_resources(self):
        result = self.init("create-example", "--resources", "scripts,references")
        self.assertEqual(result.returncode, 0, result.stderr)
        folder = self.repo / "skills/create-example"
        data = self.metadata(folder)
        self.assertEqual(data["status"], "published")
        self.assertEqual(data["description"], 'Creates examples: safely, with "quotes".')
        self.assertTrue((folder / "scripts").is_dir())
        self.assertTrue((folder / "references").is_dir())
        self.assertFalse((folder / "assets").exists())
        self.assertFalse((folder / "agents").exists())

    def test_init_draft_override_and_existing_skill_preserved(self):
        result = self.init("create-example", "--status", "draft")
        self.assertEqual(result.returncode, 0, result.stderr)
        folder = self.repo / "skills/create-example"
        before = (folder / "SKILL.md").read_bytes()
        self.assertEqual(self.metadata(folder)["status"], "draft")
        again = self.init()
        self.assertNotEqual(again.returncode, 0)
        self.assertEqual((folder / "SKILL.md").read_bytes(), before)

    def test_init_rejects_invalid_input_before_creating_files(self):
        cases = [("../escape", ()), ("create-example", ("--resources", "unexpected")),
                 ("create-example", ("--status", "hidden"))]
        for name, args in cases:
            with self.subTest(name=name, args=args):
                result = self.init(name, *args)
                self.assertNotEqual(result.returncode, 0)
                self.assertEqual(list((self.repo / "skills").iterdir()), [])
        outside = Path(self.temp.name) / "outside"
        outside.mkdir()
        (self.repo / "skills/create-example").symlink_to(outside, target_is_directory=True)
        result = self.init()
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(list(outside.iterdir()), [])

    def test_validate_preserves_source_and_supported_extensions(self):
        folder = self.write_skill(extra="tags: [examples]\ncompatibility: Requires Bash.\ncustom-setting: true\nx-claude:\n  argument-hint: example\n")
        before = (folder / "SKILL.md").read_bytes()
        result = self.run_command("validate-skill.sh", folder)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual((folder / "SKILL.md").read_bytes(), before)

    def test_validate_rejects_duplicate_keys_and_unfinished_instructions(self):
        cases = [("metadata:\n  owner: one\n  owner: two\n", "Complete the example.\n"),
                 ("tags: [1]\n", "Complete the example.\n"),
                 ("x-claude:\n  name: collision\n", "Complete the example.\n"),
                 ("", "[TODO: Finish these instructions.]\n")]
        for index, (extra, body) in enumerate(cases):
            with self.subTest(index=index):
                folder = self.write_skill(f"create-example-{index}", extra, body)
                result = self.run_command("validate-skill.sh", folder)
                self.assertNotEqual(result.returncode, 0)
                self.assertTrue(result.stderr.strip() or result.stdout.strip())

    def test_metadata_updates_preserve_policy_dependencies_and_other_fields(self):
        folder = self.write_skill()
        (folder / "agents").mkdir()
        path = folder / "agents/openai.yaml"
        existing = {"interface": {"icon_small": "./assets/icon.svg"},
                    "policy": {"allow_implicit_invocation": False},
                    "dependencies": {"tools": [{"type": "mcp", "value": "example"}]}}
        path.write_text(yaml.safe_dump(existing))
        result = self.run_command("write-agent-metadata.sh", folder,
            "--interface", "display_name=Example authoring",
            "--interface", "short_description=Create useful repository examples",
            "--interface", "default_prompt=Use $create-example to make an example.")
        self.assertEqual(result.returncode, 0, result.stderr)
        data = yaml.safe_load(path.read_text())
        self.assertEqual(data["policy"], existing["policy"])
        self.assertEqual(data["dependencies"], existing["dependencies"])
        self.assertEqual(data["interface"]["icon_small"], "./assets/icon.svg")
        self.assertIn("$create-example", data["interface"]["default_prompt"])

    def test_metadata_invalid_override_does_not_write(self):
        folder = self.write_skill()
        result = self.run_command("write-agent-metadata.sh", folder,
                                  "--interface", "short_description=short")
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse((folder / "agents").exists())


if __name__ == "__main__":
    unittest.main()
