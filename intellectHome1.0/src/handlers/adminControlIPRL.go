package handlers

import (
	"net/http"

	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/storage"
)

type AdminControlIPRl struct {
	stor *storage.Storage
}

func (a *AdminControlIPRl) ControlIPRLHandler(w http.ResponseWriter, r *http.Request) {

}
