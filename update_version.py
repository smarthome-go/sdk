#!/usr/bin/env python3
# This file is used to update the version number in all relevant places
# The SemVer (https://semver.org) versioning system is used.
import re

main_go_path = "main.go"
changelog_path = "./docs/CHANGELOG.md"
makefile_path = "Makefile"

with open(main_go_path, "r") as main_go:
    content = main_go.read()
    old_version = content.split("const Version = \"")[1].split("\"\n")[0]
    print(f"Found old version in {main_go_path}: {old_version}")

try:
    VERSION = input(
        f"Current version: {old_version}\nNew version (without 'v' prefix): ")
except KeyboardInterrupt:
    print("\nCanceled by user")
    quit()

if VERSION == "":
    VERSION = old_version

if not re.match(r"^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$", VERSION):
    print(
        f"\x1b[31mThe version: '{VERSION}' is not a valid SemVer version.\x1b[0m")
    quit()


with open(main_go_path, "w") as main_go:
    main_go.write(content.replace(old_version, VERSION))

# The `CHANGELOG.md`
with open(changelog_path, "r") as changelog:
    content = changelog.read()
    old_version = content.split("## Changelog for v")[1].split("\n")[0]
    print(f"Found old version in {changelog_path}: {old_version}")

with open(changelog_path, "w") as changelog:
    changelog.write(content.replace(old_version, VERSION))

# The `Makefile`
with open(makefile_path, "r") as makefile:
    content = makefile.read()
    old_version = content.split("version := ")[1].split("\n")[0]
    print(f"Found old version in {makefile_path}: {old_version}")

with open(makefile_path, "w") as makefile:
    makefile.write(content.replace(old_version, VERSION))

print(f"Version has been changed from '{old_version}' -> '{VERSION}'")
