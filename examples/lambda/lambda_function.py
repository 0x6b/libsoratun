import ctypes
import json
import sys

# Expected libsoratun.so, arc.json, lambda_function.py are all be located in the same directory
soratun = ctypes.cdll.LoadLibrary("libsoratun.so")

with open("arc.json", "r") as file:
    config = ctypes.c_char_p(file.read().encode("utf-8"))

def lambda_handler(event, context):
    method = ctypes.c_char_p(b"POST")
    path = ctypes.c_char_p(b"/")
    body = ctypes.c_char_p("hello from lambda".encode("utf-8"))
    soratun.Send(config, method, path, body)
    return {
        'statusCode': 200,
        'body': "Successfully sent to the unified endpoint"
    }
