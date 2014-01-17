package broker

import (
	"net/http"
	"net/http/httputil"
)

type handler func(http.ResponseWriter, *http.Request) error

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1.Log request
	if dump, err := httputil.DumpRequest(r, true); err != nil {
		//log the error
	} else {
		//log the request
	}

	// 2.Verify version header
	// http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#api-version-header

	// 3.Verify authentication header
	//http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#authentication

	// 4. Execute and handle error
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func catalogHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func provisioningHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func bindingHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}
