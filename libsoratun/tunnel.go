package libsoratun

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

const (
	SoracomNameServer1 = "100.127.0.53"
	SoracomNameServer2 = "100.127.1.53"
)

type tunnel struct {
	device   *device.Device
	tunnel   tun.Device
	net      *netstack.Net
	resolver *net.Resolver
}

// Close frees tunnel related resources.
func (t *tunnel) Close() error {
	if t.device != nil {
		t.device.Close()
	}

	t.device, t.net, t.tunnel = nil, nil, nil
	return nil
}

// DialContext exposes internal net.DialContext for consumption.
func (t *tunnel) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return t.net.DialContext(ctx, network, addr)
}

// Resolver returns internal resolver for the tunnel. Since we use gVisor as TCP stack we have to implement DNS resolver by ourselves.
func (t *tunnel) Resolver() *net.Resolver {
	return t.resolver
}

func NewTunnel(config *Config) (*tunnel, error) {
	t, n, err := netstack.CreateNetTUN(
		[]netip.Addr{config.ArcSession.ArcClientPeerIpAddress},
		[]netip.Addr{netip.MustParseAddr(SoracomNameServer1), netip.MustParseAddr(SoracomNameServer2)},
		1420,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create a tunnel: %w", err)
	}

	logger := device.NewLogger(2, "(libsoratun) ")

	dev := device.NewDevice(t, conn.NewDefaultBind(), logger)

	conf := fmt.Sprintf(`private_key=%s
public_key=%s
endpoint=%s
allowed_ip=0.0.0.0/0
`,
		config.PrivateKey.AsHexString(),
		config.ArcSession.ArcServerPeerPublicKey.AsHexString(),
		config.ArcSession.ArcServerEndpoint.String(),
	)

	if err := dev.IpcSet(conf); err != nil {
		return nil, fmt.Errorf("failed to configure device: %w", err)
	}

	if err := dev.Up(); err != nil {
		return nil, fmt.Errorf("failed to setup device: %w", err)
	}

	return &tunnel{
		device: dev,
		tunnel: t,
		net:    n,
		resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return n.DialContext(ctx, network, "100.127.0.53:53")
			},
		},
	}, nil
}
