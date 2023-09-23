const { Library } = require("ffi-napi");
const { readFileSync } = require("fs");

const soratun = Library("../../lib/shared/libsoratun", {
  Send: ["string", ["string", "string", "string", "string"]],
});

const config = readFileSync(process.argv[2], "utf8");
soratun.Send(config, "POST", "/", process.argv[3]);
