package broker

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
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
	if dump, err := httputil.DumpRequest(req, true); err != nil {
		log.Printf("Cannot log incoming request: %v", err)
	} else {
		log.Print(string(dump))
	}

	major, minor, err := extractVersion(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Router: Version check: [%v.%v]", major, minor)
	//TODO: Verify compatibility

	username, password, err := extractCredentials(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	log.Printf("Router: Authentication: [%v/%v]", username, password)
	//TODO: Authenticate based on the opts object

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
		log.Printf("Error occured while marshalling response entity: %v", err)
	}
}

// Helpers
func extractVersion(req *http.Request) (int, int, error) {
	versions, _ := req.Header["X-Broker-Api-Version"]
	if len(versions) != 1 {
		return 0, 0, errors.New("Missing Broker API version")
	}
	tokens := strings.Split(versions[0], ".")
	if len(tokens) != 2 {
		return 0, 0, errors.New("Invalid Broker API version")
	}
	major, err1 := strconv.Atoi(tokens[0])
	minor, err2 := strconv.Atoi(tokens[1])
	if err1 != nil || err2 != nil {
		return 0, 0, errors.New("Invalid Broker API version")
	}
	return major, minor, nil
}

func extractCredentials(req *http.Request) (string, string, error) {
	auths, _ := req.Header["Authorization"]
	if len(auths) != 1 {
		return "", "", errors.New("Unauthorized access")
	}
	tokens := strings.Split(auths[0], " ")
	if len(tokens) != 2 || tokens[0] != "Basic" {
		return "", "", errors.New("Unsupported authentication method")
	}
	raw, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		return "", "", errors.New("Unable to decode 'Authorization' header")
	}
	credentials := strings.Split(string(raw), ":")
	if len(credentials) != 2 {
		return "", "", errors.New("Missing credentials")
	}
	return credentials[0], credentials[1], nil
}
