package pdf

import (
	"log"
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
  pdf,err:=h.service.ListPdf(r.Context())
  if err !=nil{
	log.Println(err)
	http.Error(w,err.Error(),http.StatusInternalServerError)
	return
    }
json.Write(w,http.StatusOK,pdf)

}
