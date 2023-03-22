package resources

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	detailURL = "/resource/detail"
	createURL = "/resource/create"
)

type ResourceError struct {
	Message string `json:"error"`
}

var _ Handler = &resourceHandler{}

type resourceHandler struct {
	ResourceService Service
}

type Handler interface {
	Register(router *mux.Router)
}

func GetHandler(resourceService Service) Handler {
	h := resourceHandler{
		ResourceService: resourceService,
	}
	return &h
}

func (h *resourceHandler) Register(router *mux.Router) {
	router.HandleFunc(detailURL, h.getResource).Methods(http.MethodGet)
	router.HandleFunc(createURL, h.createResource).Methods(http.MethodPost)
}

func (h *resourceHandler) createResource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var res resourceDTO

	err := dec.Decode(&res)

	if err != nil {
		writeResponse(w, &ResourceError{Message: "bad request"}, http.StatusTeapot)
		return
	}

	name := h.ResourceService.generateName(res.Type, res.Region)
	allocatedName := h.ResourceService.putAllocatedName(name, res.Type, res.Region)
	h.ResourceService.info(fmt.Sprintf("Type: %s", res.Type))
	h.ResourceService.info(fmt.Sprintf("Region: %s", res.Region))
	h.ResourceService.info(fmt.Sprintf("Allocated Name: %s", allocatedName))

	writeResponse(w, allocatedName, http.StatusOK)
}

func (h *resourceHandler) getResource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name, err := requestDetails(r, "name")
	if err != nil {
		writeResponse(w, &ResourceError{Message: "bad request"}, http.StatusBadRequest)
		return
	}

	o := h.ResourceService.getAllocatedName(name)
	if o == nil {
		writeResponse(w, &ResourceError{Message: "not found"}, http.StatusNotFound)
		return
	}
	writeResponse(w, o, http.StatusOK)
}

func requestDetails(r *http.Request, key string) (val string, err error) {
	q := r.URL.Query()
	if val := q.Get(key); val == "" {
		return "", fmt.Errorf("unable to fetch key %s from request", key)
	} else {
		return val, nil
	}
}

func writeResponse(w http.ResponseWriter, val interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(val)
}
