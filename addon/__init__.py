from anki import collection, decks, hooks, notes, utils
from aqt import mw

from .backend import GoBackend
from .models import OxfordLookupError, OxfordLookupFailed
from .rules import should_lookup

config = mw.addonManager.getConfig(__name__)
logger = mw.addonManager.get_logger(__name__)

be = GoBackend()


def oxfordlookup_audiouk(
    col: collection.Collection, note: notes.Note, deck_id: decks.DeckId
) -> None:
    try:
        wf = config["word_field"]
        af = config["audio_field"]
        to = config["timeout_seconds"]
        ow_if_existing = config["overwrite_existing"]
        word = utils.strip_html(note[wf])
        audio = utils.strip_html(note[af])
    except KeyError as e:
        logger.error("error when getting config or fields: %s", e)
        return

    # keep existing audio unless overwrite is enabled
    if not should_lookup(
        word=word.strip(), audio=audio.strip(), overwrite=ow_if_existing
    ):
        return

    try:
        res = be.run_cmd(word=word, timeout=to)
        audio_uk = res["audio_uk"]
        note[af] = audio_uk
    except FileNotFoundError as e:
        logger.error("not found oxford backend binary %s", e)
    except OxfordLookupFailed:
        logger.exception("error when looking up word=%r", word)
    except OxfordLookupError:
        logger.warning("not found or no audio for word=%s", word)
    except Exception:
        logger.exception("bug while looking up %r", word)


hooks.note_will_be_added.append(oxfordlookup_audiouk)
