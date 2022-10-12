package getcode

import (
	"context"
	"net/http"
	"time"
)

func GetCode() string {
	srv := &http.Server{Addr: ":8070"}
	code := ""

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")
		srv.Shutdown(context.Background())
	})

	srv.ListenAndServe()
	time.Sleep(1 * time.Second)

	return code
}
