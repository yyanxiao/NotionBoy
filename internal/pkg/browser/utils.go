package browser

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// forceIP tries to force the host component in urlstr to be an IP address.
//
// Since Chrome 66+, Chrome DevTools Protocol clients connecting to a browser
// must send the "Host:" header as either an IP address, or "localhost".
func forceIP(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		return urlstr
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return urlstr
	}
	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return urlstr
	}
	u.Host = net.JoinHostPort(addr.IP.String(), port)
	return u.String()
}

// detectURL detects the websocket debugger URL if the provided URL is not a
// valid websocket debugger URL.
//
// A valid websocket debugger URL is something like:
// ws://127.0.0.1:9222/devtools/browser/...
// The original URL with the following formats are accepted:
// * ws://127.0.0.1:9222/
// * http://127.0.0.1:9222/
func detectURL(urlstr string) string {
	if strings.Contains(urlstr, "/devtools/browser/") {
		return urlstr
	}

	// replace the scheme and path to construct the URL like:
	// http://127.0.0.1:9222/json/version
	u, err := url.Parse(urlstr)
	if err != nil {
		return urlstr
	}
	u.Scheme = "http"
	u.Path = "/json/version"

	// to get "webSocketDebuggerUrl" in the response
	resp, err := http.Get(forceIP(u.String()))
	if err != nil {
		return urlstr
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return urlstr
	}
	// the browser will construct the debugger URL using the "host" header of the /json/version request.
	// for example, run headless-shell in a container: docker run -d -p 9000:9222 chromedp/headless-shell:latest
	// then: curl http://127.0.0.1:9000/json/version
	// and the debugger URL will be something like: ws://127.0.0.1:9000/devtools/browser/...
	wsURL := result["webSocketDebuggerUrl"].(string)
	return wsURL
}
