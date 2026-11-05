package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/url"

	"github.com/lqqyt2423/go-mitmproxy/cert"
	"github.com/lqqyt2423/go-mitmproxy/internal/helper"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	Debug             int
	Addr              string
	StreamLargeBodies int64 // 当请求或响应体大于此字节时，转为 stream 模式
	SslInsecure       bool
	CaRootPath        string
	NewCaFunc         func() (cert.CA, error) //创建 Ca 的函数
	Upstream          string
	LogFilePath       string // Path to write logs to file
}

type Proxy struct {
	Opts    *Options
	Version string
	Addons  []Addon

	entry            *entry
	attacker         *attacker
	webSocketHandler *webSocketHandler
	shouldIntercept  func(req *http.Request) bool              // req is received by proxy.server
	upstreamProxy    func(req *http.Request) (*url.URL, error) // req is received by proxy.server, not client request
	authProxy        func(res http.ResponseWriter, req *http.Request) (bool, error)
	dialCheck        func(host string) (string, error) // ARTEX: 拨号前门禁 + IP 钉住
}

// proxy.server req context key
var proxyReqCtxKey = new(struct{})

func NewProxy(opts *Options) (*Proxy, error) {
	if opts.StreamLargeBodies <= 0 {
		opts.StreamLargeBodies = 1024 * 1024 * 5 // default: 5mb
	}

	proxy := &Proxy{
		Opts:    opts,
		Version: "1.9.2",
		Addons:  make([]Addon, 0),
	}

	proxy.entry = newEntry(proxy)

	attacker, err := newAttacker(proxy)
	if err != nil {
		return nil, err
	}
	proxy.attacker = attacker

	proxy.webSocketHandler = newWebSocketHandler(proxy)

	return proxy, nil
}

func (proxy *Proxy) AddAddon(addon Addon) {
	proxy.Addons = append(proxy.Addons, addon)
}

func (proxy *Proxy) Start() error {
	go func() {
		if err := proxy.attacker.start(); err != nil {
			log.Error(err)
		}
	}()
	return proxy.entry.start()
}

func (proxy *Proxy) Close() error {
	return proxy.entry.close()
}

func (proxy *Proxy) Shutdown(ctx context.Context) error {
	return proxy.entry.shutdown(ctx)
}

func (proxy *Proxy) GetCertificate() x509.Certificate {
	return *proxy.attacker.ca.GetRootCA()
}

func (proxy *Proxy) GetCertificateByCN(commonName string) (*tls.Certificate, error) {
	return proxy.attacker.ca.GetCert(commonName)
}

func (proxy *Proxy) SetShouldInterceptRule(rule func(req *http.Request) bool) {
	proxy.shouldIntercept = rule
}

func (proxy *Proxy) SetUpstreamProxy(fn func(req *http.Request) (*url.URL, error)) {
	proxy.upstreamProxy = fn
}

// ARTEX 补丁：dialCheck 在真正拨号前判定目标并返回可钉住的 "ip:port"。
// 返回 ("", nil) 表示按原地址拨号；返回 error 表示拒绝；
// 返回非空地址表示用该地址拨号——它把「校验过的 IP」和「实际拨号的 IP」绑定，
// 消除 DNS rebinding（校验时解析到公网、拨号时解析回私网）的 TOCTOU 窗口。
func (proxy *Proxy) SetDialCheck(fn func(host string) (string, error)) {
	proxy.dialCheck = fn
}

func (proxy *Proxy) realUpstreamProxy() func(*http.Request) (*url.URL, error) {
	return func(cReq *http.Request) (*url.URL, error) {
		req := cReq.Context().Value(proxyReqCtxKey).(*http.Request)
		return proxy.getUpstreamProxyUrl(req)
	}
}

func (proxy *Proxy) getUpstreamProxyUrl(req *http.Request) (*url.URL, error) {
	if proxy.upstreamProxy != nil {
		return proxy.upstreamProxy(req)
	}
	if len(proxy.Opts.Upstream) > 0 {
		return url.Parse(proxy.Opts.Upstream)
	}
	cReq := &http.Request{URL: &url.URL{Scheme: "https", Host: req.Host}}
	return http.ProxyFromEnvironment(cReq)
}

func (proxy *Proxy) getUpstreamConn(ctx context.Context, req *http.Request) (net.Conn, error) {
	proxyUrl, err := proxy.getUpstreamProxyUrl(req)
	if err != nil {
		return nil, err
	}
	var conn net.Conn
	address := helper.CanonicalAddr(req.URL)
	if proxy.dialCheck != nil {
		// 传 address 而不是 req.Host：address 由 CanonicalAddr 生成，已按 scheme
		// 带上正确端口（http→80/https→443），否则明文 HTTP 会被误钉到 443。
		pinned, err := proxy.dialCheck(address)
		if err != nil {
			return nil, err
		}
		if pinned != "" {
			address = pinned // 钉住校验过的 IP，防 DNS rebinding
		}
	}
	if proxyUrl != nil {
		conn, err = helper.GetProxyConn(ctx, proxyUrl, address, proxy.Opts.SslInsecure)
	} else {
		conn, err = (&net.Dialer{}).DialContext(ctx, "tcp", address)
	}
	return conn, err
}

func (proxy *Proxy) SetAuthProxy(fn func(res http.ResponseWriter, req *http.Request) (bool, error)) {
	proxy.authProxy = fn
}
