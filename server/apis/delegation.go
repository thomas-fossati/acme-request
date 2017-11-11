package apis

import (
	"encoding/json"
	"net/http"

	"../models"
	"github.com/gorilla/mux"
)

type (
	// ServiceInterface is the blueprint that the API layer expects the Service layer to implement
	ServiceInterface interface {
		ReadDelegation(id string) (*models.Delegation, error)
		AddDelegation(delegation *models.Delegation) (*models.Delegation, error)
	}

	delegationResource struct {
		serviceInterface ServiceInterface
	}
)

// SetupDelegationRoutes tells the supplied /router/ how to handle the "/delegation" resource(s)
func SetupDelegationRoutes(router *mux.Router, serviceInterface ServiceInterface) {
	ur := &delegationResource{serviceInterface}

	// TODO(tho) add DELETE as an acceptable method for "/star/delegation/{id}"
	router.HandleFunc("/star/delegation/{id}", RequestLogger(ur.readDelegation)).Methods("GET")
	router.HandleFunc("/star/delegation", RequestLogger(ur.addDelegation)).Methods("POST")
	// Add other delegation API methods here...
}

func (ur *delegationResource) addDelegation(w http.ResponseWriter, r *http.Request) error {
	var delegation models.Delegation

	jsonDecoder := json.NewDecoder(r.Body)

	err := jsonDecoder.Decode(&delegation)
	if err != nil {
		return err
	}

	out, err := ur.serviceInterface.AddDelegation(&delegation)
	if err != nil {
		return err
	}

	// Encode to body
	body, err := json.Marshal(out)
	if err != nil {
		return err
	}

	// Send reply
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Etag", models.ComputeETag(out))
	w.WriteHeader(http.StatusCreated)
	w.Write(body)

	return nil
}

func (ur *delegationResource) readDelegation(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	// Forward to the service layer
	delegation, err := ur.serviceInterface.ReadDelegation(id)
	if err != nil {
		return err
	}

	// Encode to body
	body, err := json.Marshal(delegation)
	if err != nil {
		return err
	}

	// Send reply
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Etag", models.ComputeETag(delegation))
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	return nil
}
