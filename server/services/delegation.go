package services

import "../models"

type (
	StoreInterface interface {
		ReadDelegation(id string) (*models.Delegation, error)
		AddDelegation(delegation *models.Delegation) (*models.Delegation, error)
	}

	DelegationService struct {
		storeInterface StoreInterface
	}
)

// NewDelegationService instantiates a new DelegationService
func NewDelegationService(storeInterface StoreInterface) *DelegationService {
	return &DelegationService{storeInterface}
}

// ReadDelegation tries to return a Delegation record associated with the supplied /id/
func (us *DelegationService) ReadDelegation(id string) (*models.Delegation, error) {
	// No much to do in this case, except forwarding to the storage layer
	return us.storeInterface.ReadDelegation(id)
}

// AddDelegation adds a Delegation to the store
func (us *DelegationService) AddDelegation(delegation *models.Delegation) (*models.Delegation, error) {
	// TODO(tho) add CSR validation against template
	return us.storeInterface.AddDelegation(delegation)
}
