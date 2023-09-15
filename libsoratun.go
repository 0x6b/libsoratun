package main

import (
	"C"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/0x6b/libsoratun/internal"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func main() {
}

const (
	UnifiedEndpointHostname = "100.127.69.42"
	UnifiedEndpointPort     = 80
	SoracomNameServer1      = "100.127.0.53"
	SoracomNameServer2      = "100.127.1.53"
)

type tunnel struct {
	device   *device.Device
	tunnel   tun.Device
	net      *netstack.Net
	resolver *net.Resolver
}

// Send sends a request to the unified endpoint with given Config and HTTP headers.
//
//export Send
func Send(configJson *C.char, method, path, body *C.char) *C.char {
	config, err := newConfig([]byte(C.GoString(configJson)))
	if err != nil {
		return nil
	}

	t, err := createTunnel(config)
	if err != nil {
		return nil
	}
	endpoint, _ := url.Parse(fmt.Sprintf("http://%s:%d", UnifiedEndpointHostname, UnifiedEndpointPort))

	c := &internal.UnifiedEndpointHTTPClient{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: t.DialContext,
			},
		},
		Endpoint: endpoint,
		Headers:  []string{"User-Agent: libsoratun/0.0.1"},
	}

	m := C.GoString(method)
	p := C.GoString(path)
	b := C.GoString(body)
	if !(m == http.MethodGet || m == http.MethodPost) {
		return nil
	}

	req, err := c.MakeRequest(&internal.Params{
		Path:    strings.TrimPrefix(p, "/"),
		Body:    strings.NewReader(b),
		Method:  C.GoString(method),
		Headers: c.Headers,
	})
	if err != nil {
		return nil
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("failed to read response from Unified Endpoint: %v\n", err)
		panic(err)
	}

	return C.CString(string(response))
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

func createTunnel(config *Config) (*tunnel, error) {
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

const arcServerEndpointDefaultPort string = "11010"

// UDPAddr represents the UDP address with keeping original endpoint.
type UDPAddr struct {
	IP          net.IP
	Port        int
	RawEndpoint []byte
}

// aliases to add custom un/marshaler to each type.
type (
	// Key is an alias for wgtypes.Key.
	Key wgtypes.Key
	// IPNet is an alias for net.IPNet.
	IPNet net.IPNet
)

// Config holds SORACOM Arc client configurations.
type Config struct {
	// PrivateKey is WireGuard private key.
	PrivateKey Key `json:"privateKey"`
	// ArcSession holds connection information provided from SORACOM Arc server.
	ArcSession *ArcSession `json:"arcSessionStatus,omitempty"`
}

// ArcSession holds SORACOM Arc configurations received from the server.
type ArcSession struct {
	// ArcServerPeerPublicKey is WireGuard public key of the SORACOM Arc server.
	ArcServerPeerPublicKey Key `json:"arcServerPeerPublicKey"`
	// ArcServerEndpoint is a UDP endpoint of the SORACOM Arc server.
	ArcServerEndpoint *UDPAddr `json:"arcServerEndpoint"`
	// ArcAllowedIPs holds IP addresses allowed for routing from the SORACOM Arc server.
	ArcAllowedIPs []*IPNet `json:"arcAllowedIPs"`
	// ArcClientPeerPrivateKey holds private key from SORACOM Arc server.
	ArcClientPeerPrivateKey Key `json:"arcClientPeerPrivateKey,omitempty"`
	// ArcClientPeerIpAddress is an IP address for this client.
	ArcClientPeerIpAddress netip.Addr `json:"arcClientPeerIpAddress,omitempty"`
}

// newConfig creates a new Config from a byte array of JSON.
func newConfig(configJson []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(configJson, &config)
	if err != nil {
		return nil, fmt.Errorf("error while reading config file: %s", err)
	}
	return &config, nil
}

// UnmarshalText decodes a byte array of private key to the Key. If text is invalid WireGuard key, UnmarshalText returns an error.
func (k *Key) UnmarshalText(text []byte) error {
	key, err := wgtypes.ParseKey(string(text))
	if err != nil {
		return err
	}
	copy(k[:], key[:])
	return nil
}

// AsWgKey converts Key back to wgtypes.Key.
func (k *Key) AsWgKey() *wgtypes.Key {
	key, _ := wgtypes.NewKey(k[:])
	return &key
}

// AsHexString returns hexadecimal encoding of Key.
func (k *Key) AsHexString() string {
	return hex.EncodeToString(k[:])
}

// UnmarshalText converts a byte array into UDPAddr. UnmarshalText returns error if the format is invalid (not "ip" or "ip:port"), IP address specified is invalid, or the port is not a 16-bit unsigned integer.
func (a *UDPAddr) UnmarshalText(text []byte) error {
	h, p, err := net.SplitHostPort(string(text))
	if err != nil {
		h = string(text)
		p = arcServerEndpointDefaultPort
	}

	var ip net.IP
	ip = net.ParseIP(h)
	if ip == nil {
		ips, err := net.LookupIP(h)
		if err != nil || len(ips) < 1 {
			return fmt.Errorf("invalid endpoint \"%s\": %s", h, err)
		}
		ip = ips[0]
	}

	port, err := strconv.Atoi(p)
	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("invalid serverEndpoint port number: %s, it should be a 16-bit unsigned integer", p)
	}

	a.IP, a.Port = ip, port
	a.RawEndpoint = text
	return nil
}

// String returns string representation of UDPAddr for WireGuard configuration.
func (a UDPAddr) String() string {
	if a.IP.To4() != nil {
		return fmt.Sprintf("%s:%d", a.IP, a.Port)
	} else if a.IP.To16() != nil {
		return fmt.Sprintf("[%s]:%d", a.IP, a.Port)
	} else {
		return ""
	}
}

// UnmarshalText converts a byte array into IPNet. UnmarshalText returns error if invalid CIDR is provided.
func (n *IPNet) UnmarshalText(text []byte) error {
	_, ipnet, err := net.ParseCIDR(string(text))
	if err != nil {
		return err
	}

	n.IP, n.Mask = ipnet.IP, ipnet.Mask
	return nil
}
