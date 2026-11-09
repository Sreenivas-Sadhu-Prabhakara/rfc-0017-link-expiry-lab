package main

import (
	"html/template"
	"strings"
)

var page = template.Must(template.New("report").Funcs(template.FuncMap{
	"names": func(items []string) string {
		if len(items) == 0 {
			return "—"
		}
		return strings.Join(items, ", ")
	},
}).Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>RFC 17 · Link expiry lab</title>
<style>
:root{color-scheme:light;background:#edf3f3;color:#152b35;font:17px/1.55 system-ui,sans-serif}
main{max-width:1000px;margin:auto;padding:40px 22px}
.kicker{font:700 .83rem monospace;letter-spacing:.1em;color:#1a605c}
h1{font-size:clamp(2.3rem,6vw,4.3rem);line-height:1.02;margin:.35em 0}
p{max-width:72ch} a{color:#115e67}
.table{overflow-x:auto;border:1px solid #aac1c0;border-radius:12px;background:#fff}
table{border-collapse:collapse;width:100%;text-align:left;font-size:.92rem}
caption{text-align:left;padding:14px;font-weight:700}
th,td{padding:9px 12px;border-top:1px solid #d8e5e4;white-space:nowrap}
thead{background:#dbeceb} tbody tr:nth-child(even){background:#f5f9f9}
</style>
</head>
<body><main>
<p class="kicker">INTERNET RFC SERIES · 017 · GO LAB</p>
<h1>What expires when a link goes quiet?</h1>
<p>RFC 17’s reply separated an idle table entry from a blocked entry with a transfer in flight. This educational trace shows how a bounded table behaves when active and idle entries compete for space.</p>
<div class="table" role="region" aria-label="Link table event trace" tabindex="0">
<table>
<caption>{{.Capacity}} table slots · {{.IdleSeconds}}-second idle expiry</caption>
<thead><tr><th scope="col">Time (s)</th><th scope="col">Event</th><th scope="col">Outcome</th><th scope="col">Expired</th><th scope="col">Active</th><th scope="col">Idle</th></tr></thead>
<tbody>{{range .Steps}}
<tr><th scope="row">{{.At}}</th><td>{{.Action}} {{.Link}}</td><td>{{.Outcome}}</td><td>{{names .Expired}}</td><td>{{names .Active}}</td><td>{{names .Idle}}</td></tr>{{end}}
</tbody></table></div>
<p>The Go CLI generated this static report. It runs locally with a synthetic event script.</p>
<p><a href="https://www.rfc-editor.org/info/rfc17/">Read RFC 17 and Kahn’s reply</a> · <a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0017-link-expiry-lab">Source and run instructions</a></p>
</main></body></html>
`))
