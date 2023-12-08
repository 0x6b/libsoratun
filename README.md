# libsoratun

The C library allows you to embed [Soracom Arc](https://users.soracom.io/ja-jp/docs/arc/) connectivity into your own program. You can send a message to the unified endpoint, with Soracom Arc, entirely from userspace (no root privilege is required).

## Tested Setup

- Go 1.21.5 darwin/arm64
- macOS Sonoma 14.1.2

## Prerequisites

1. You have to have a virtual SIM, along with `arc.json` which is a configuration file for [`soratun`](https://github.com/soracom/soratun/) locally. See documentation for detail.
   - [Soracom Arc Soratun Tool](https://developers.soracom.io/en/docs/arc/soratun/) (English)
   - [soratun を利用して接続する: soratun の概要と機能](https://users.soracom.io/ja-jp/docs/arc/soratun-overview/) (Japanese)
2. You have to enable the unified endpoint for your SIM group. See documentation for detail.
   - [Unified Endpoint Overview](https://developers.soracom.io/en/docs/unified-endpoint/) (English)
   - [Unified Endpoint](https://users.soracom.io/ja-jp/docs/unified-endpoint/) (Japanese)

## Build

```console
$ git clone https://github.com/soracom/libsoratun
$ cd libsoratun
$ make libs
```

## Run Examples

### Python

Tested with Python 3.12.0

- [`examples/python/main.py`](examples/python/main.py)

```console
$ cd examples/python
$ python3 main.py /path/to/arc.json '{"message": "hey"}'
```

### Node.js

Tested with Node.js v18.19.0.

- [`examples/nodejs/src/index.js`](examples/nodejs/src/index.js)

```console
$ cd examples/nodejs
$ npm install
$ node src/index.py /path/to/arc.json '{"message": "hey"}'
```

### Rust

Tested with Rust 1.74.1. In order to build Rust bindings, you have to install [bindgen-cli](https://rust-lang.github.io/rust-bindgen/command-line-usage.html) and its prerequisites.

- [`examples/rust/src/main.rs`](examples/rust/src/main.rs)

```console
$ make bindings
$ cd examples/rust
$ cargo run -- --config /path/to/arc.json '{"message": "hey"}'
```

### AWS Lambda (Python 3.11)

1. Build a shared library on target platform. Tested on arm64.

   ```console
   $ make libs
   ```
2. Place `examples/lamdada/lambda_function.py`, `lib/shared/libsoratun.so`, and your `arc.json` in a same directory.
3. Zip it up and upload it to AWS Lambda function.

   ```console
   $ zip -r lambda.zip lambda_function.py libsoratun.so arc.json
   ```

## License

MIT. See [LICENSE](LICENSE) for detail.
 