package timestamp

import (
	"fmt"
	"net/http"
)

func init() {
	// Root = howto
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, HOME_TEMPLATE)
	})

	// Raw: spits text answer only
	http.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")
		t, err := parseUnknown(q)

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
		t, err := parseUnknown(q)

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
	"date": %v
}`


const HOME_TEMPLATE = `
<html>
<head>
	<title>Online timestamp stuff</title>
</head>
<body>
	<h1>Online timestamp &amp; date (approximate) codec</h1>
	<div>
		<dl>
			<dd>/</dd>
			<dt>display this message</dt>
			<dd>/raw?q=<b><i>123</b></i></dd>
			<dt>try to decode <b><i>123</b></i>, and gives plain text result</dt>
			<dd>/json?q=<b><i>123</b></i></dd>
			<dt>try to decode <b><i>123</b></i>, and gives JSON result</dt>
		</dl>
		
		<hr/>
		
		<form action="/raw">
			Raw <input type="text" name="q" value="123" /> <input type="submit" value="&gt;" />
		</form>
		<form action="/json">
			JSON <input type="text" name="q" value="123" /> <input type="submit" value="&gt;" />
		</form>
	</div>
</body>
</html>
`