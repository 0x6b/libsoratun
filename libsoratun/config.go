package libsoratun

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

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
	// LogLevel specifies logging level, verbose (2), error (1), or silent (0).
	LogLevel int `json:"logLevel"`
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

// NewConfig creates a new Config from a byte array of JSON.
func NewConfig(configJson []byte) (*Config, error) {
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

// String returns string representation of Config. The private key is masked.
func (c *Config) String() string {
	var ips []string
	for _, ip := range c.ArcSession.ArcAllowedIPs {
		ips = append(ips, (*net.IPNet)(ip).String())
	}

	return fmt.Sprintf(`[Interface]
Address = %s/32
PrivateKey = <secret>

[Peer]
PublicKey = %s
AllowedIPs = %s
Endpoint = %s
`,
		c.ArcSession.ArcClientPeerIpAddress,
		c.ArcSession.ArcServerPeerPublicKey.AsHexString(),
		strings.Join(ips, ", "),
		c.ArcSession.ArcServerEndpoint.String(),
	)
}
