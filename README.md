# libsoratun

The C library allows you to embed Soracom Arc connectivity into your own program. You can send a message to the unified endpoint, with Soracom Arc, entirely from userspace (no root privilege is required).

## Tested Setup

- Go 1.21.0
- Rust 1.72.0
- Python 3.9.6
- Node.js v18.14.2
- macOS Ventura 13.4.1

## Prerequisites

1. You have to have `arc.json`, configuration file for [`soratun`](https://github.com/soracom/soratun/) locally. See documentation for detail.
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

- [`examples/rust/src/main.rs`](examples/rust/src/main.rs)

```console
$ cd examples/rust
$ cargo run -- --config /path/to/arc.json '{"message": "hey"}'
```

### Python

- [`examples/python/main.py`](examples/python/main.py)

```console
$ cd examples/python
$ python3 main.py /path/to/arc.json '{"message": "hey"}'
```

### Node.js

- [`examples/nodejs/src/index.js`](examples/nodejs/src/index.js)

```console
$ cd examples/nodejs
$ npm install
$ node src/index.py /path/to/arc.json '{"message": "hey"}'
```

## License

MIT. See [LICENSE](LICENSE) for detail.
