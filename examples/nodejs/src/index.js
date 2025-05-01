const { open, load, DataType, Result, unwrapErr } = require('node-ffi-rs');
const { readFileSync } = require("fs");
const { platform } = require("os");
open({
  library: 'libsoratun',
  path: '../../lib/shared/libsoratun' + (platform() === "win32" ? ".dll" : ".so")
})

const config = readFileSync(process.argv[2], "utf8");
const response = load({
  library: 'libsoratun',
  funcName: 'Send',
  retType: DataType.String,
  paramsType: [DataType.String, DataType.String, DataType.String, DataType.String],
  paramsValue: [config, "POST", "/", process.argv[3]]
});
console.log(response);
