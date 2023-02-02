package server

import (
	"fmt"
	"html/template"
	"net/http"
	"notionboy/internal/pkg/config"
	"strings"
)

var ipfsTemplate = `
<html>
<head>
	<title>NotionBoy IPFS Proxy</title>
	<style>
		body {
			font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
			font-size: 30px;
			line-height: 1.42857143;
			color: #333;
			background-color: #fff;
		}
		h1 {
			margin: 20px 0;
			font-size: 48px;
			font-weight: 500;
			line-height: 1.1;
			color: inherit;
		}
		p {
			margin: 0 0 10px;
		}
		ul {
			padding-left: 0;
			list-style: none;
		}
		li {
			display: block;
		}
		a {
			color: #337ab7;
			text-decoration: none;
		}
		a:hover {
			color: #23527c;
			text-decoration: underline;
		}
	</style>
</head>
<body>
	<h1>Download Links</h1>
	<p>Click the link to download the file.</p>
	<div>
		<ul>
			{{range .}} <li><a href="{{.Url}}">{{.Name}}</a></li> {{end}}
		</ul>
	</div>
</body>
</html>
`

type Link struct {
	Name string
	Url  string
}

func proxyIpfs(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/files/ipfs/")
	if path == "" {
		renderError(w, http.StatusBadRequest, "Bad request", nil)
		return
	}
	cid := strings.Split(path, ".")[0]
	gateways := config.GetConfig().Zlib.IpfsGateways
	links := make([]Link, 0)
	for _, gateway := range gateways {
		links = append(links, Link{
			Name: gateway,
			Url:  fmt.Sprintf("%s/ipfs/%s", gateway, cid),
		})
	}
	t, err := template.New("ipfsProxy").Parse(ipfsTemplate)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Parse template error", err)
		return
	}
	err = t.Execute(w, links)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Execute template error", err)
		return
	}
}
