#!/usr/bin/env python3

import argparse
from pathlib import Path
import sys


def main():
    parser = argparse.ArgumentParser(
        description="Print the bundled report template without writing files."
    )
    parser.parse_args()
    template = Path(__file__).resolve().parent.parent / "assets" / "report-template.md"
    try:
        content = template.read_text(encoding="utf-8")
    except (OSError, UnicodeError) as error:
        parser.exit(1, f"Cannot read {template}: {error}\n")
    sys.stdout.write(content)


if __name__ == "__main__":
    main()
