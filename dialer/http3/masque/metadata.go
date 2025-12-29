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

	// Connection pooling
	connectionPoolingDisabled bool
}

func (d *masqueDialer) parseMetadata(md mdata.Metadata) (err error) {
	d.md.host = mdutil.GetString(md, mdKeyHost)

	if mdutil.GetBool(md, mdKeyKeepAlive) {
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
	d.md.connectionPoolingDisabled = mdutil.GetBool(md, mdKeyConnectionPoolingDisabled)

	return nil
}
