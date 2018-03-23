package main

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	ocontext "github.com/aunem/transpose/pkg/context"

	log "github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/utils"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; http://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func getURLFromRequest(req *http.Request) *url.URL {
	// If the Request was created by Go via a real HTTP request,  RequestURI will
	// contain the original query string. If the Request was created in code, RequestURI
	// will be empty, and we will use the URL object instead
	u := req.URL
	if req.RequestURI != "" {
		parsedURL, err := url.ParseRequestURI(req.RequestURI)
		if err == nil {
			u = parsedURL
		} else {
			log.Warnf("error when parsing RequestURI: %s", err)
		}
	}
	return u
}

// changeTarget the request to handle the target URL
func changeTarget(outReq *http.Request, target *url.URL) {
	outReq.URL = utils.CopyURL(outReq.URL)
	outReq.URL.Scheme = target.Scheme
	outReq.URL.Host = target.Host

	u := getURLFromRequest(outReq)

	outReq.URL.Path = u.Path
	outReq.URL.RawPath = u.RawPath
	outReq.URL.RawQuery = u.RawQuery
	outReq.RequestURI = "" // Outgoing request should not have RequestURI

	// Do not pass client Host header unless optsetter PassHostHeader is set.
	// if !b.passHost {
	// 	outReq.Host = target.Host
	// }

	outReq.Proto = "HTTP/1.1"
	outReq.ProtoMajor = 1
	outReq.ProtoMinor = 1

	// if f.rewriter != nil {
	// 	f.rewriter.Rewrite(outReq)
	// }
}

// RequestToRoundtrip prepares an inbound request for its roundtrip
func RequestToRoundtrip(rc *ocontext.HTTPRequest) *ocontext.HTTPRequest {
	req := rc.Request
	ctx := req.Context()
	// if cn, ok := rw.(http.CloseNotifier); ok {   need to think this part through
	// 	var cancel context.CancelFunc
	// 	ctx, cancel = context.WithCancel(ctx)
	// 	defer cancel()
	// 	notifyChan := cn.CloseNotify()
	// 	go func() {
	// 		select {
	// 		case <-notifyChan:
	// 			cancel()
	// 		case <-ctx.Done():
	// 		}
	// 	}()
	// }

	outreq := new(http.Request)
	*outreq = *req // includes shallow copies of maps, but okay
	if req.ContentLength == 0 {
		outreq.Body = nil // Issue 16036: nil Body for http.Transport retries
	}
	outreq = outreq.WithContext(ctx)
	outreq.Close = false

	// We are modifying the same underlying map from req (shallow
	// copied above) so we only copy it if necessary.
	copiedHeaders := false

	// Remove hop-by-hop headers listed in the "Connection" header.
	// See RFC 2616, section 14.10.
	if c := outreq.Header.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				if !copiedHeaders {
					outreq.Header = make(http.Header)
					copyHeader(outreq.Header, req.Header)
					copiedHeaders = true
				}
				outreq.Header.Del(f)
			}
		}
	}

	// Remove hop-by-hop headers to the backend. Especially
	// important is "Connection" because we want a persistent
	// connection, regardless of what the client sent to us.
	for _, h := range hopHeaders {
		if outreq.Header.Get(h) != "" {
			if !copiedHeaders {
				outreq.Header = make(http.Header)
				copyHeader(outreq.Header, req.Header)
				copiedHeaders = true
			}
			outreq.Header.Del(h)
		}
	}

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := outreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outreq.Header.Set("X-Forwarded-For", clientIP)
	}
	rc.Request = outreq
	return rc
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
