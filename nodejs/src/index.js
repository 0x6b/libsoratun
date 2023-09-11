const { Library } = require("ffi-napi");
const { readFileSync } = require("fs");

// if you are on macOS, you have to rename the library `libsoratun.so` to `libsoratun.dylib`,
// since ffi-napi does not support `.so` extension on macOS.
const soratun = Library("libsoratun", {
  Send: ["string", ["string", "string", "string", "string"]],
});

const config = readFileSync(process.argv[2], "utf8");
soratun.Send(config, "POST", "/", process.argv[3]);
