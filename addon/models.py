import dataclasses

@dataclasses.dataclass()
class Result:
    word:       str
    audio_uk:   str
    source_url: str

class OxfordLookupError(Exception):
    def __init__(self) -> None:
        super().__init__("OxfordLookup failed with error: NotFoundError or NoAudioError")


class OxfordLookupFailed(Exception):
    def __init__(self) -> None:
        super().__init__("OxfordLookup failed with error: NetworkError")
