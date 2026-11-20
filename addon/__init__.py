# import the main window object (mw) from aqt
from aqt import mw, gui_hooks,addcards
# import the "show info" tool from utils.py
from aqt.utils import showInfo, qconnect
# import all of the Qt GUI library
from aqt.qt import *
from .backend import GoBackend
from PyQt6.QtWidgets import QLineEdit

# We're going to add a menu item below. First we want to create a function to
# be called when the menu item is activated.
config = mw.addonManager.getConfig(__name__)

be = GoBackend()
def oxfordlookup_audiouk(card: addcards.AddCards) -> None:
    # res = be.run_cmd(word=word, timeout=config["timeout_seconds"])
    # audio_uk = res["audio_uk"]
    # for c in childs:
    #     print(c.objectName())
    #return audio_uk




gui_hooks.add_cards_did_init.append(oxfordlookup_audiouk)

