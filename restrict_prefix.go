package restrictprefix

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(&RestrictPrefix{})
	httpcaddyfile.RegisterHandlerDirective("restrict_prefix", parseCaddyfileHandler)
}

// RestrictPrefix is middleware that restricts requests where any portion
// of the URI matches a given prefix.
type RestrictPrefix struct {
	logger *zap.Logger

	Prefix string `json:"prefix,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (*RestrictPrefix) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.restrict_prefix",
		New: func() caddy.Module { return new(RestrictPrefix) },
	}
}

// Provision a Zap logger to RestrictPrefix.
func (p *RestrictPrefix) Provision(ctx caddy.Context) error {
	p.logger = ctx.Logger(p)

	return nil
}

// ServeHTTP implements the caddyhttp.MiddlewareHandler interface.
func (p *RestrictPrefix) ServeHTTP(
	w http.ResponseWriter, r *http.Request, next caddyhttp.Handler,
) error {
	for _, part := range strings.Split(r.URL.Path, "/") {
		if strings.HasPrefix(part, p.Prefix) {
			http.Error(w, "Not Found", http.StatusNotFound)
			if p.logger != nil {
				p.logger.Debug(fmt.Sprintf(
					"restricted prefix: %q in %s", part, r.URL.Path))
			}
			return nil
		}
	}

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
//
// Syntax:
//
//	restrict_prefix <prefix>
func (p *RestrictPrefix) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&p.Prefix) {
			return d.ArgErr()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
	}

	return nil
}

// Validate the prefix from the module's configuration, setting the
// default prefix "." if necessary.
func (p *RestrictPrefix) Validate() error {
	if p.Prefix == "" {
		p.Prefix = "."
	}

	return nil
}

func parseCaddyfileHandler(
	h httpcaddyfile.Helper,
) (caddyhttp.MiddlewareHandler, error) {
	m := new(RestrictPrefix)
	if err := m.UnmarshalCaddyfile(h.Dispenser); err != nil {
		return nil, err
	}

	return m, nil
}

var (
	_ caddy.Module                = (*RestrictPrefix)(nil)
	_ caddy.Provisioner           = (*RestrictPrefix)(nil)
	_ caddy.Validator             = (*RestrictPrefix)(nil)
	_ caddyfile.Unmarshaler       = (*RestrictPrefix)(nil)
	_ caddyhttp.MiddlewareHandler = (*RestrictPrefix)(nil)
)
