package handlers

import (
	"io/fs"
	"net/http"
)

type StaticHandler struct {
	fs.FS
	Directory string
}

func NewStaticHandler(fs fs.FS, dir string) *StaticHandler {
	return &StaticHandler{FS: fs, Directory: dir}
}

func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	staticFS, err := fs.Sub(h.FS, h.Directory)
	if err != nil {
		panic(err)
	}

	// TODO Enable caching of resources for 24 hours
	//cacheControlValue := fmt.Sprintf("public, max-age=%d", int((24 * time.Hour).Seconds()))
	//w.Header().Set("Cache-Control", cacheControlValue)

	fileServer := http.FileServer(http.FS(staticFS))
	http.StripPrefix("/"+h.Directory, fileServer).ServeHTTP(w, r)
}
