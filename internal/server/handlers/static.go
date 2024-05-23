package handlers

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/weesvc/weesvc-gorilla/internal/config"
)

type StaticHandler struct {
	fs.FS
	Config    *config.Config
	Directory string
}

func NewStaticHandler(config *config.Config, fs fs.FS, dir string) *StaticHandler {
	return &StaticHandler{FS: fs, Config: config, Directory: dir}
}

func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	staticFS, err := fs.Sub(h.FS, h.Directory)
	if err != nil {
		panic(err)
	}

	if h.Config.ResourceCachingEnabled {
		cacheControlValue := fmt.Sprintf("public, max-age=%d", int((24 * time.Hour).Seconds()))
		w.Header().Set("Cache-Control", cacheControlValue)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	http.StripPrefix("/"+h.Directory, fileServer).ServeHTTP(w, r)
}
