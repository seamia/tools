package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux"
	"github.com/seamia/tools/assets"
)

func main() {
	addr := ":3434"
	mux := httptreemux.NewContextMux()
	populateMux(mux)
	httpServer := &http.Server{Addr: addr, Handler: mux}
	fmt.Println("serving on :", addr)
	if err := httpServer.ListenAndServe(); err != nil {
		onError(err)
	}
}

func populateMux(mux *httptreemux.ContextMux) {
	mux.GET("/", constructHandler("/home.html"))
	mux.GET("/info", constructTemplateHandler("/two.html", process))
}

func constructHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := assets.ServeAsset(w, req, name); err != nil {
			w.WriteHeader(http.StatusNotExtended)
			w.Write([]byte(fmt.Sprintf("Error [%s] serving [%s]", err, req.URL.Path)))
		}
	}
}

type LiveDataFunc func() interface{}

func constructTemplateHandler(name string, lifeData LiveDataFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := assets.ServeAssetTemplate(w, req, name, lifeData()); err != nil {
			w.WriteHeader(http.StatusNotExtended)
			w.Write([]byte(fmt.Sprintf("Error [%s] serving [%s]", err, req.URL.Path)))
		}
	}
}

func process() interface{} {
	return map[string]string{
		"Name": "Example Service",
	}
}

func onError(err error) {
	fmt.Fprint(os.Stderr, "Error :", err)
}
