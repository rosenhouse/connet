package handlers

import (
	"lib/marshal"
	"net/http"

	"github.com/pivotal-golang/lager"
)

type RulesList struct {
	Marshaler marshal.Marshaler
	Logger    lager.Logger
}

func (h *RulesList) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger := h.Logger.Session("list-rules")

	payload, err := h.Marshaler.Marshal([]string{"whatever"})
	if err != nil {
		logger.Error("marshal-failed", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(payload)
}
