package broker

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type router struct {
	mux *mux.Router // TODO: Replace with own simpler regexp-based mux???
}

func NewRouter(h *rabbitHandler) *router {
	mux := mux.NewRouter()
	mux.Handle("/v2/catalog", reponseHandler(h.catalog)).Methods("GET")
	mux.Handle("/v2/service_instances/{iid}", reponseHandler(h.provision)).Methods("PUT")
	mux.Handle("/v2/service_instances/{iid}", reponseHandler(h.deprovision)).Methods("DELETE")
	mux.Handle("/v2/service_instances/{iid}/service_bindings/{bid}", reponseHandler(h.bind)).Methods("PUT")
	mux.Handle("/v2/service_instances/{iid}/service_bindings/{bid}", reponseHandler(h.unbind)).Methods("DELETE")
	return &router{mux}
}

// Log & verify request and then pass it to Gorilla to be dispatched approprietly.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	dump, _ := httputil.DumpRequest(req, true)
	log.Print(string(dump))

	// Verify version header
	// See http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#api-version-header
	if versions, _ := req.Header["X-Broker-Api-Version"]; len(versions) == 1 {
		//TODO
		log.Print(versions[0])
	}

	// Authenticate the request
	// See http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#authentication
	auths, _ := req.Header["Authorization"]
	if len(auths) != 1 {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	tokens := strings.SplitN(auths[0], " ", 2)
	if len(tokens) != 2 || tokens[0] != "Basic" {
		http.Error(w, "Unsupported authentication method", http.StatusUnauthorized)
		return
	}
	var credentials []string
	if raw, err := base64.StdEncoding.DecodeString(tokens[1]); err != nil {
		http.Error(w, "Unable to decode Authorization header", http.StatusBadRequest)
		return
	} else {
		credentials = strings.SplitN(string(raw), ":", 2)
	}
	//TODO
	log.Print(credentials)

	// Disptach
	r.mux.ServeHTTP(w, req)
}

type responseEntity struct {
	status int
	value  interface{}
}

type reponseHandler func(*http.Request) responseEntity

// Marshall the response entity as JSON and return the proper HTTP status code.
func (fn reponseHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	re := fn(req)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(re.status)
	if err := json.NewEncoder(w).Encode(re.value); err != nil {
		// TODO: Cannot do much just log it here
	}
}
