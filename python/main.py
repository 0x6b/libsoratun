import ctypes
import sys

soratun = ctypes.cdll.LoadLibrary("../lib/shared/libsoratun.so")

with open(sys.argv[1], "r") as file:
    config = ctypes.c_char_p(file.read().encode("utf-8"))

method = ctypes.c_char_p(b"POST")
path = ctypes.c_char_p(b"/")
body = ctypes.c_char_p(" ".join(sys.argv[2:]).encode("utf-8"))

soratun.SendRequest(config, method, path, body)
