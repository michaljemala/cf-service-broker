package broker

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type brokerError struct {
	Description string
}

type rabbitHandler struct {
	broker ServiceBroker
}

func NewHandler(b ServiceBroker) *rabbitHandler {
	return &rabbitHandler{b}
}

func (h *rabbitHandler) catalog(r *http.Request) responseEntity {
	if cat, err := h.broker.Catalog(); err != nil {
		return responseEntity{http.StatusInternalServerError, brokerError{err.Error()}}
	} else {
		return responseEntity{http.StatusOK, cat}
	}
}

func (h *rabbitHandler) provision(req *http.Request) responseEntity {
	vars := mux.Vars(req)
	preq := ProvisioningRequest{Id: vars["iid"]}

	if err := json.NewDecoder(req.Body).Decode(&preq); err != nil {
		return responseEntity{http.StatusBadRequest, brokerError{err.Error()}}
	}

	if url, err := h.broker.Provision(preq); err != nil {
		return responseEntity{http.StatusConflict, brokerError{err.Error()}}
	} else {
		return responseEntity{http.StatusCreated, struct {
			Dashboard_url string
		}{url}}
	}
}

func (h *rabbitHandler) deprovision(req *http.Request) responseEntity {
	vars := mux.Vars(req)
	preq := ProvisioningRequest{Id: vars["iid"]}

	if err := json.NewDecoder(req.Body).Decode(&preq); err != nil {
		return responseEntity{http.StatusBadRequest, brokerError{err.Error()}}
	}

	if err := h.broker.Deprovision(preq); err != nil {
		return responseEntity{http.StatusNotFound, struct{}{}}
	} else {
		return responseEntity{http.StatusOK, struct{}{}}
	}
}

func (h *rabbitHandler) bind(req *http.Request) responseEntity {
	vars := mux.Vars(req)
	breq := BindingRequest{InstanceId: vars["iid"], Id: vars["bid"]}

	if err := json.NewDecoder(req.Body).Decode(&breq); err != nil {
		return responseEntity{http.StatusBadRequest, brokerError{err.Error()}}
	}

	if cred, url, err := h.broker.Bind(breq); err != nil {
		return responseEntity{http.StatusConflict, struct{}{}}
	} else {
		return responseEntity{http.StatusCreated, struct {
			Credentials      interface{}
			Syslog_drain_url string
		}{cred, url}}
	}
}

func (h *rabbitHandler) unbind(req *http.Request) responseEntity {
	vars := mux.Vars(req)
	breq := BindingRequest{InstanceId: vars["iid"], Id: vars["bid"]}

	if err := json.NewDecoder(req.Body).Decode(&breq); err != nil {
		return responseEntity{http.StatusBadRequest, brokerError{err.Error()}}
	}

	if err := h.broker.Unbind(breq); err != nil {
		return responseEntity{http.StatusNotFound, struct{}{}}
	} else {
		return responseEntity{http.StatusOK, struct{}{}}
	}
}
