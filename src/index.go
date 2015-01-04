package timestamp

import (
	"fmt"
	"net/http"
)

const JSON_TEMPLATE = `{
	"timestampSeconds": %v,
	"timestampMilliSeconds": %v,
	"timestampNanoseconds": %v,
	"date": %v
}`

func init() {
	// Root = howto
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "/ : display this message")
		fmt.Fprintln(w, "/raw?q=123 : try to decode 123, and prints plain text result")
		fmt.Fprintln(w, "/json?q=123 : try to decode 123, and prints JSON result")
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
