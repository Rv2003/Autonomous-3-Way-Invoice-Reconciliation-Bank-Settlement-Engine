package pdf

import (
	"encoding/json"
	"net/http"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/json"

	
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListPdf(w http.ResponseWriter,r *http.Request) {
pdf:=[]string{"hello","world"}

json.Write(w,http.StatusOK,pdf)

}