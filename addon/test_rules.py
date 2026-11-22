import pytest

from .rules import should_lookup


@pytest.mark.parametrize(
    "word,audio,overwrite,expected",
    [
        ("", "", False, False),
        ("chromatic", "", False, True),
        ("chromatic", "url", False, False),
        ("chromatic", "url", True, True),
    ],
)
def test_should_lookup(word, audio, overwrite, expected):
    assert should_lookup(word, audio, overwrite) == expected
