def should_lookup(word: str, audio: str, overwrite: bool) -> bool:
    if not word:
        return False
    if not audio:
        return True
    return overwrite
