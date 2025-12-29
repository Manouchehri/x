package masque

import (
	"time"

	mdata "github.com/go-gost/core/metadata"
	mdutil "github.com/go-gost/x/metadata/util"
)

const (
	mdKeyHost                      = "host"
	mdKeyKeepAlive                 = "keepAlive"
	mdKeyKeepAlivePeriod           = "ttl"
	mdKeyHandshakeTimeout          = "handshakeTimeout"
	mdKeyMaxIdleTimeout            = "maxIdleTimeout"
	mdKeyMaxStreams                = "maxStreams"
	mdKeyConnectionPoolingDisabled = "connectionPoolingDisabled"
	mdKeyInitialPacketSize         = "initialPacketSize"
	mdKeyDisablePathMTUDiscovery   = "disablePathMTUDiscovery"

	// Default timeouts to prevent connections from hanging indefinitely
	defaultHandshakeTimeout = 30 * time.Second
	defaultMaxIdleTimeout   = 60 * time.Second
	defaultKeepAlivePeriod  = 20 * time.Second // Under typical NAT timeout of 30s
)

type metadata struct {
	host string

	// QUIC config options
	keepAlivePeriod  time.Duration
	maxIdleTimeout   time.Duration
	handshakeTimeout time.Duration
	maxStreams       int

	// MTU options
	initialPacketSize        int
	disablePathMTUDiscovery  bool

	// Connection pooling
	connectionPoolingDisabled bool
}

func (d *masqueDialer) parseMetadata(md mdata.Metadata) (err error) {
	d.md.host = mdutil.GetString(md, mdKeyHost)

	// Enable QUIC keepalive by default for reliability through NAT/firewalls
	// Can be disabled with keepAlive=false
	keepAlive := mdutil.GetString(md, mdKeyKeepAlive)
	if keepAlive != "false" && keepAlive != "0" {
		d.md.keepAlivePeriod = mdutil.GetDuration(md, mdKeyKeepAlivePeriod)
		if d.md.keepAlivePeriod <= 0 {
			d.md.keepAlivePeriod = defaultKeepAlivePeriod
		}
	}

	d.md.handshakeTimeout = mdutil.GetDuration(md, mdKeyHandshakeTimeout)
	if d.md.handshakeTimeout <= 0 {
		d.md.handshakeTimeout = defaultHandshakeTimeout
	}

	d.md.maxIdleTimeout = mdutil.GetDuration(md, mdKeyMaxIdleTimeout)
	if d.md.maxIdleTimeout <= 0 {
		d.md.maxIdleTimeout = defaultMaxIdleTimeout
	}

	d.md.maxStreams = mdutil.GetInt(md, mdKeyMaxStreams)

	// MTU options - for constrained networks (e.g., 1280 MTU requires initialPacketSize=1232)
	d.md.initialPacketSize = mdutil.GetInt(md, mdKeyInitialPacketSize)
	d.md.disablePathMTUDiscovery = mdutil.GetBool(md, mdKeyDisablePathMTUDiscovery)

	d.md.connectionPoolingDisabled = mdutil.GetBool(md, mdKeyConnectionPoolingDisabled)

	return nil
}
