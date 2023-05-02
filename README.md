# libsoratun

The C library allows you to embed Soracom Arc connectivity into your own program. You can sent a message to the unified point, with Soracom Arc, entirely from userspace.

## Tested Setup

- Go 1.20.3
- Rust 1.68.2
- Python 3.9.6
- macOS Ventura 13.3

## Prerequisites

1. You have to have `arc.json`, configuration file for [`soratun`](https://github.com/soracom/soratun/) locally. See documentation for detail.
   - [Soracom Arc Soratun Tool](https://developers.soracom.io/en/docs/arc/soratun/) (English)
   - [soratun を利用して接続する: soratun の概要と機能](https://users.soracom.io/ja-jp/docs/arc/soratun-overview/) (Japanese)
2. You have to enable the unified endpoint for your SIM group. See documentation for detail.
   - [Unified Endpoint Overview](https://developers.soracom.io/en/docs/unified-endpoint/) (English)
   - [Unified Endpoint](https://users.soracom.io/ja-jp/docs/unified-endpoint/) (Japanese)

## Build and Run Examples

```console
$ git clone https://github.com/soracom/libsoratun
$ cd libsoratun
$ make all
$ # Rust
$ cd rust
$ cargo run -- --config /path/to/arc.json '{"message": "hey"}'
$ # Python
$ cd python
$ python3 main.py /path/to/arc.json '{"message": "hey"}'
```

See [`rust/src/main.rs`](rust/src/main.rs) or [`python/main.py`](python/main.py) for example usage.

## License

See [LICENSE](LICENSE) for detail.
