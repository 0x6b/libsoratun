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
$ make all
```

## Run Examples

### Rust

Tested with Rust 1.74.1.

- [`examples/rust/src/main.rs`](examples/rust/src/main.rs)

```console
$ cd examples/rust
$ cargo run -- --config /path/to/arc.json '{"message": "hey"}'
```

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

## License

MIT. See [LICENSE](LICENSE) for detail.
