"""
Define all interactions with Go library
"""

import json
import os
import platform
import subprocess
import sys

from .models import OxfordLookupError, OxfordLookupFailed


class GoBackend:
    def __init__(self) -> None:

        if sys.platform not in ["win32", "darwin", "linux"]:
            raise Exception("Platform isn't supported")

        _os = sys.platform
        if sys.platform == "win32":
            _os = "windows"

        _arch = platform.machine()
        if platform.machine() in ["aarch64", "arm64"]:
            _arch = "arm64"
        if platform.machine() in ["x86_64", "amd64"]:
            _arch = "amd64"

        self._bin = os.path.join(
            os.path.dirname(__file__), "bin", f"oxford-lookup-{_os}-{_arch}"
        )
        self._subprocess_timeout = 5

    def run_cmd(self, word: str, timeout: int) -> dict:
        """
        Lookup word information

        timeout:    int (seconds)
        word:       str
        """
        if not os.path.exists(self._bin):
            raise FileNotFoundError(self._bin)

        cmd = [str(self._bin), "-timeout", str(timeout), "-word", str(word)]
        resp = subprocess.run(
            cmd, capture_output=True, timeout=timeout + self._subprocess_timeout
        )
        exit_code = resp.returncode

        if exit_code != 0:
            if exit_code == 1:  # no_audio or not_found
                raise OxfordLookupError()
            if exit_code == 2:  # network
                raise OxfordLookupFailed()

        stdout = resp.stdout
        return json.loads(stdout)
