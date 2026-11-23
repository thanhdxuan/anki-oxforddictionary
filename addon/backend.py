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

        self._os = sys.platform
        _ext = ""

        self._popen_kwargs = {}
        if sys.platform == "win32":
            self._os = "windows"
            _ext = ".exe"
            self._popen_kwargs["creationflags"] = subprocess.CREATE_NO_WINDOW

        self._arch = platform.machine()
        if platform.machine() in ["aarch64", "arm64"]:
            self._arch = "arm64"
        if platform.machine() in ["x86_64", "amd64"]:
            self._arch = "amd64"
        self._bin = os.path.join(
            os.path.dirname(__file__),
            "bin",
            f"oxford-lookup-{self._os}-{self._arch}{_ext}",
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
            cmd,
            capture_output=True,
            timeout=timeout + self._subprocess_timeout,
            **self._popen_kwargs,
        )

        exit_code = resp.returncode

        if exit_code != 0:
            if exit_code == 1:  # no_audio or not_found
                raise OxfordLookupError()
            if exit_code == 2:  # network
                raise OxfordLookupFailed()

        stdout = resp.stdout
        return json.loads(stdout)
