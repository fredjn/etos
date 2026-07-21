# Copyright Axis Communications AB.
#
# For a full list of individual contributors, please see the commit history.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""Write file metadata to a JSON manifest."""

import json
from pathlib import Path


class ManifestWriter:  # pylint: disable=too-few-public-methods
    """Incrementally write file metadata to a workspace manifest."""

    def __init__(self, workspace: Path) -> None:
        """Initialize the manifest in the workspace root."""
        self.__path = workspace.joinpath("file_manifest.json")
        self.__entries: list[dict[str, str]] = []
        self.__write()

    def add(self, name: str, url: str) -> None:
        """Add a file and update the manifest."""
        self.__entries.append({"name": name, "url": url})
        self.__write()

    def __write(self) -> None:
        """Atomically write the current entries."""
        temporary_path = self.__path.with_suffix(".json.tmp")
        with temporary_path.open("w", encoding="utf-8") as manifest:
            json.dump({"files": self.__entries}, manifest, separators=(",", ":"))
        temporary_path.replace(self.__path)
