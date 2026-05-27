import os
import sys

if getattr(sys, "frozen", False) and not hasattr(sys, "_MEIPASS"):
    base = os.path.dirname(os.path.abspath(sys.executable))
    internal = os.path.join(base, "_internal")
    sys._MEIPASS = internal if os.path.isdir(internal) else base
