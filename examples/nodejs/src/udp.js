const { open, load, DataType, Result, unwrapErr } = require('node-ffi-rs');

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

try {
  open({
    library: 'libsoratun', // key
    path: '../../lib/shared/libsoratun.dll' // path
  })

  const response = load({
    library: 'libsoratun',
    funcName: 'SendUDP',
    retType: DataType.String,
    paramsType: [DataType.String, DataType.U8Array, DataType.I64],
    paramsValue: [config, message, message.length]
  });
  console.log(response);
} catch (error) {
  console.error('Error:', error);
}
