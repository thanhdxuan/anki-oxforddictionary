# import the main window object (mw) from aqt
from aqt import mw, gui_hooks, addcards
# import the "show info" tool from utils.py
from aqt.utils import showInfo
# import all of the Qt GUI library
from aqt.qt import *
from .backend import GoBackend
from PyQt6.QtWidgets import QLineEdit
import anki
from anki import hooks, collection, notes, decks

# We're going to add a menu item below. First we want to create a function to
# be called when the menu item is activated.
config = mw.addonManager.getConfig(__name__)

be = GoBackend()

def oxfordlookup_audiouk(col: collection.Collection, note: notes.Note, deck_id: decks.DeckId) -> None:
    """(col: anki.collection.Collection, note: anki.notes.Note, deck_id: anki.decks.DeckId)"""
    fid = note._field_index("Word")
    aid = note._field_index("Audio")
    word = note.fields[fid]
    res = be.run_cmd(word=word, timeout=config["timeout_seconds"])
    audio_uk = res["audio_uk"]
    note.fields[aid] = audio_uk

hooks.note_will_be_added.append(oxfordlookup_audiouk)


