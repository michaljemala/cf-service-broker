package broker

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var empty struct{} = struct{}{}

type handler struct {
	brokerService BrokerService
}

func newHandler(bs BrokerService) *handler {
	return &handler{bs}
}

func (h *handler) catalog(r *http.Request) responseEntity {
	if cat, err := h.brokerService.Catalog(); err != nil {
		return responseEntity{http.StatusInternalServerError, BrokerError{err.Error()}}
	} else {
		return responseEntity{http.StatusOK, cat}
	}
}

func (h *handler) provision(req *http.Request) responseEntity {
	vars := mux.Vars(req)
	preq := ProvisioningRequest{Id: vars[instanceId]}

	if err := json.NewDecoder(req.Body).Decode(&preq); err != nil {
		return responseEntity{http.StatusBadRequest, BrokerError{err.Error()}}
	}

	url, err := h.brokerService.Provision(preq)
	if err != nil {
		return handleServiceError(err)
	}

	return responseEntity{http.StatusCreated, struct {
		Dashboard_url string
	}{url}}
}

func (h *handler) deprovision(req *http.Request) responseEntity {
	vars := mux.Vars(req)
	preq := ProvisioningRequest{Id: vars[instanceId]}

	if err := json.NewDecoder(req.Body).Decode(&preq); err != nil {
		return responseEntity{http.StatusBadRequest, BrokerError{err.Error()}}
	}

	if err := h.brokerService.Deprovision(preq); err != nil {
		return handleServiceError(err)
	}

	return responseEntity{http.StatusOK, empty}
}

func (h *handler) bind(req *http.Request) responseEntity {
	vars := mux.Vars(req)
	breq := BindingRequest{InstanceId: vars[instanceId], Id: vars[serviceId]}

	if err := json.NewDecoder(req.Body).Decode(&breq); err != nil {
		return responseEntity{http.StatusBadRequest, BrokerError{err.Error()}}
	}

	cred, url, err := h.brokerService.Bind(breq)
	if err != nil {
		return handleServiceError(err)
	}

	return responseEntity{http.StatusCreated, struct {
		Credentials      interface{}
		Syslog_drain_url string
	}{cred, url}}
}

func (h *handler) unbind(req *http.Request) responseEntity {
	vars := mux.Vars(req)
	breq := BindingRequest{InstanceId: vars[instanceId], Id: vars[serviceId]}

	if err := json.NewDecoder(req.Body).Decode(&breq); err != nil {
		return responseEntity{http.StatusBadRequest, BrokerError{err.Error()}}
	}

	if err := h.brokerService.Unbind(breq); err != nil {
		return handleServiceError(err)
	}

	return responseEntity{http.StatusOK, empty}
}

func handleServiceError(err error) responseEntity {
	log.Printf("Service error occured: %v", err.Error())

	switch err := err.(type) {
	case BrokerServiceError:
		switch err.Code() {
		case ErrCodeConflict:
			return responseEntity{http.StatusConflict, empty}
		case ErrCodeGone:
			return responseEntity{http.StatusGone, empty}
		}
	}
	return responseEntity{http.StatusInternalServerError, BrokerError{err.Error()}}
}
