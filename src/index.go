package timestamp

import (
	"fmt"
	"html/template"
	"net/http"
)

func init() {
	htmlTemplate, err := template.New("tmpl").Parse(HTML_TEMPLATE)
	if err!=nil{
		panic(err)
	}
	
	// Root = howto
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")
		t, unit, clean, err := parseUnknown(q)
		if err == nil {
			// We have a valid parameter
			data := Response{
				"Input": clean,
				"S": t.Unix(),
				"Ms": t.UnixNano()/1000000,
				"Mµs": t.UnixNano()/1000,
				"Ns": t.UnixNano(),
				"DateStr": t,
				"Unit": unit,
				"IsATimestamp": (unit=="seconds" || unit=="milliseconds" || unit=="microseconds" || unit=="nanoseconds"),
			}
			err := htmlTemplate.ExecuteTemplate(w, "tmpl", data)
			if err!=nil{
				panic(err)
			}
			return
		}

		fmt.Fprintln(w, HOME_TEMPLATE)
	})

	// Raw: spits text answer only
	http.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")
		t, _, _, err := parseUnknown(q)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprintf(w, "%v\n", t.Unix())
		fmt.Fprintf(w, "%v\n", t.UnixNano()/1000000)
		fmt.Fprintf(w, "%v\n", t.UnixNano())
		fmt.Fprintf(w, "%v\n", t)
	})

	// JSON
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")
		t, _, _, err := parseUnknown(q)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{ error:%q }", err.Error())
			return
		}
		fmt.Fprintf(w, JSON_TEMPLATE,
			t.Unix(), t.UnixNano()/1000000, t.UnixNano(), t)
	})
}

const JSON_TEMPLATE = `{
	"timestampSeconds": %v,
	"timestampMilliSeconds": %v,
	"timestampNanoseconds": %v,
	"date": "%v"
}`

const HOME_TEMPLATE = `
<html>
<head>
	<title>Online timestamp stuff</title>
</head>
<body>
	<h1>Online timestamp &amp; date (approximate) decoder</h1>
	<div>
		<dl>
			<dd>/</dd>
			<dt>display this message</dt>
			<dd>/raw?q=<b><i>1425979122</b></i></dd>
			<dt>try to decode <b><i>1425979122</b></i>, and gives plain text result</dt>
			<dd>/json?q=<b><i>1425979122</b></i></dd>
			<dt>try to decode <b><i>1425979122</b></i>, and gives JSON result</dt>
		</dl>
		
		<hr/>
		
		<form action="/raw">
			Raw <input type="text" name="q" value="1425979122" /> <input type="submit" value="&gt;" />
		</form>
		<form action="/json">
			JSON <input type="text" name="q" value="1425979122" /> <input type="submit" value="&gt;" />
		</form>
		<form action="/">
			HTML <input type="text" name="q" value="1425979122" /> <input type="submit" value="&gt;" />
		</form>
	</div>
</body>
</html>
`

const HTML_TEMPLATE = `
<html>
<head>
	<title>Online timestamp for value {.Input}}</title>
	<style>abbr{ border-bottom: 2px dotted; }</style>
</head>
<body>
	<h1>Date for value {{.Input}}</h1>
	<div>
		This value seems to be a timestamp in <b>{{.Unit}}</b> {{if .IsATimestamp}}since <abbr title="1970-01-01 00:00:00">epoch</abbr>.{{end}}
		<table>
			<tr><th>Seconds</th><td>{{.S}}</td></tr>
			<tr><th>Milliseconds</th><td>{{.Ms}}</td></tr>
			<tr><th>Microseconds</th><td>{{.Mµs}}</td></tr>
			<tr><th>Nanoseconds</th><td>{{.Ns}}</td></tr>
			<tr><th>Date</th><td>{{.DateStr}}</td></tr>
		</table>
		<hr/>
		
		<form action="/raw">
			Raw <input type="text" name="q" value="{{.Input}}" /> <input type="submit" value="&gt;" />
		</form>
		<form action="/json">
			JSON <input type="text" name="q" value="{{.Input}}" /> <input type="submit" value="&gt;" />
		</form>
	</div>
</body>
</html>
`

// See http://nesv.blogspot.fr/2012/09/super-easy-json-http-responses-in-go.html
type Response map[string]interface{}
