"""Exercise project documentation setup through its public Bash entrypoint."""

from pathlib import Path
import re
import subprocess
import tempfile
import unittest


HELPER = Path(__file__).resolve().parents[1] / "scripts/init-project-docs.sh"
EXPECTED_FILES = {
    "AGENTS.md",
    "CLAUDE.md",
    "docs/conventions/CODE.md",
    "docs/conventions/COMMITS.md",
    "docs/conventions/GOLANG.md",
    "docs/conventions/TDD.md",
    "docs/specs/FEATURES.md",
    "docs/specs/MISSION.md",
    "docs/specs/ROADMAP.md",
    "docs/specs/STACK.md",
    "docs/decisions/README.md",
    "docs/features/README.md",
}


class ProjectDocumentationSetup(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="project-docs-test-")
        self.addCleanup(self.temp.cleanup)
        self.workspace = Path(self.temp.name)
        self.project = self.workspace / "project with spaces"
        self.project.mkdir()

    def run_helper(self, *args):
        self.assertTrue(
            HELPER.is_file(),
            "Project documentation initializer not implemented: init-project-docs.sh",
        )
        return subprocess.run(
            ["bash", str(HELPER), *map(str, args)],
            cwd=self.workspace,
            capture_output=True,
            text=True,
        )

    def test_creates_documentation_in_selected_project(self):
        # Arrange: an empty root with spaces, distinct from the working directory.

        # Act.
        result = self.run_helper(self.project)

        # Assert: all requested documents exist, without application scaffolding.
        self.assertEqual(result.returncode, 0, result.stderr)
        files = {
            str(path.relative_to(self.project))
            for path in self.project.rglob("*")
            if path.is_file()
        }
        self.assertEqual(files, EXPECTED_FILES)
        for name in EXPECTED_FILES:
            with self.subTest(name=name):
                self.assertTrue((self.project / name).read_bytes())
        self.assertEqual((self.project / "CLAUDE.md").read_text(), "@AGENTS.md\n")
        self.assertEqual(set(self.workspace.iterdir()), {self.project})

    def test_preserves_existing_documents_and_rerun_keeps_all_bytes(self):
        # Arrange: user instructions and conventions already exist.
        existing = {
            "AGENTS.md": b"# My project rules\nKeep my custom commands.\n",
            "CLAUDE.md": b"@AGENTS.md\n\nKeep Claude-specific instructions.\n",
            "docs/conventions/CODE.md": b"# Our conventions\nPreserve this file.\n",
        }
        for name, contents in existing.items():
            destination = self.project / name
            destination.parent.mkdir(parents=True, exist_ok=True)
            destination.write_bytes(contents)

        # Act: complete the missing documents, then repeat setup.
        first = self.run_helper(self.project)
        self.assertEqual(first.returncode, 0, first.stderr)
        before = {name: (self.project / name).read_bytes() for name in EXPECTED_FILES}
        again = self.run_helper(self.project)

        # Assert: both user documents and generated files retain their bytes.
        self.assertEqual(again.returncode, 0, again.stderr)
        for name, contents in before.items():
            with self.subTest(name=name):
                self.assertEqual((self.project / name).read_bytes(), contents)
                first_label = "kept" if name in existing else "created"
                self.assertRegex(first.stdout, rf"(?m)\b{first_label}\s+{re.escape(name)}$")
                self.assertRegex(again.stdout, rf"(?m)\bkept\s+{re.escape(name)}$")
        for name, contents in existing.items():
            self.assertEqual(before[name], contents)

    def test_refuses_symlinks_before_creating_other_documents(self):
        cases = {
            "document": ("CLAUDE.md", False, False),
            "dangling document": ("docs/specs/STACK.md", False, True),
            "docs directory": ("docs", True, False),
        }
        for label, (name, directory, dangling) in cases.items():
            with self.subTest(case=label):
                # Arrange: each case has its own root and outside target.
                project = self.workspace / label
                project.mkdir()
                outside = self.workspace / (label + " outside")
                if directory:
                    outside.mkdir()
                elif not dangling:
                    outside.write_bytes(b"outside content\n")
                link = project / name
                link.parent.mkdir(parents=True, exist_ok=True)
                link.symlink_to(outside, target_is_directory=directory)
                before = set(project.rglob("*"))

                # Act.
                result = self.run_helper(project)

                # Assert: refuse the link before writing any project documents.
                self.assertNotEqual(result.returncode, 0)
                self.assertIn(name, result.stdout + result.stderr)
                self.assertEqual(set(project.rglob("*")), before)
                self.assertTrue(link.is_symlink())
                if directory:
                    self.assertEqual(list(outside.iterdir()), [])
                elif dangling:
                    self.assertFalse(outside.exists())
                else:
                    self.assertEqual(outside.read_bytes(), b"outside content\n")

    def test_refuses_file_directory_collisions_before_creating_documents(self):
        cases = {
            "document is directory": ("docs/specs/STACK.md", True),
            "parent is file": ("docs/features", False),
        }
        for label, (name, directory) in cases.items():
            with self.subTest(case=label):
                # Arrange.
                project = self.workspace / label
                project.mkdir()
                collision = project / name
                collision.parent.mkdir(parents=True, exist_ok=True)
                if directory:
                    collision.mkdir()
                else:
                    collision.write_bytes(b"keep this file\n")
                before = set(project.rglob("*"))

                # Act.
                result = self.run_helper(project)

                # Assert.
                self.assertNotEqual(result.returncode, 0)
                self.assertIn(name, result.stdout + result.stderr)
                self.assertEqual(set(project.rglob("*")), before)
                if not directory:
                    self.assertEqual(collision.read_bytes(), b"keep this file\n")

    def test_refuses_competing_layouts_before_creating_documents(self):
        cases = {
            "single spec": "docs/SPEC.md",
            "feature template": "docs/templates/FEATURE.md",
            "phase template": "docs/templates/PHASE.md",
            "feature record": "docs/features/0001-login/FEATURE.md",
        }
        for label, name in cases.items():
            with self.subTest(case=label):
                # Arrange: another documentation layout already owns specs or features.
                project = self.workspace / label
                project.mkdir()
                existing = project / name
                existing.parent.mkdir(parents=True)
                existing.write_bytes(b"existing owner\n")
                before = set(project.rglob("*"))

                # Act.
                result = self.run_helper(project)

                # Assert: name the competing owner and write nothing.
                self.assertNotEqual(result.returncode, 0)
                self.assertIn(name, result.stdout + result.stderr)
                self.assertEqual(set(project.rglob("*")), before)
                self.assertEqual(existing.read_bytes(), b"existing owner\n")

    def test_allowed_existing_layout_is_kept_beside_new_documents(self):
        # Arrange: the user agreed to keep an existing single spec.
        spec = self.project / "docs/SPEC.md"
        spec.parent.mkdir()
        spec.write_bytes(b"# Existing spec\n")

        # Act.
        result = self.run_helper("--allow-existing-layout", self.project)

        # Assert: every document is added and the existing spec keeps its bytes.
        self.assertEqual(result.returncode, 0, result.stderr)
        for name in EXPECTED_FILES:
            with self.subTest(name=name):
                self.assertTrue((self.project / name).is_file())
        self.assertEqual(spec.read_bytes(), b"# Existing spec\n")

    def test_refuses_missing_arguments_and_nonexistent_root(self):
        cases = {
            "missing argument": (),
            "extra argument": (self.project, "unexpected"),
            "unknown option": ("--force", self.project),
            "missing root": (self.workspace / "missing",),
        }
        for label, args in cases.items():
            with self.subTest(case=label):
                # Act.
                result = self.run_helper(*args)

                # Assert: explain invalid input without creating a project.
                self.assertNotEqual(result.returncode, 0)
                self.assertTrue((result.stdout + result.stderr).strip())
                self.assertEqual(list(self.project.iterdir()), [])
                self.assertFalse((self.workspace / "missing").exists())


if __name__ == "__main__":
    unittest.main()
