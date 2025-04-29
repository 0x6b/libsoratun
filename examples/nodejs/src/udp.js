const { Library } = require("ffi-napi");
const { readFileSync } = require("fs");

const soratun = Library("../../lib/shared/libsoratun", {
  SendUDP: ["string", ["string","pointer","int"]]
});

const config = readFileSync(process.argv[2], "utf8");

// Send UDP Packet for IoT Button
const message = new Uint8Array(4);
message[0] = 0x4d;
message[1] = 1
message[2] = 3
message[3] = 0x4d + 1 + 3

const response = soratun.SendUDP(config,message,4)
console.log(response);
