package stores

import (
	"database/sql"
	"strconv"
	"time"

	"../apperrors"
	"../models"

	_ "github.com/mattn/go-sqlite3"
)

// DelegationID is the type of a delegation identifier (in sqlite land)
type DelegationID int64

const (
	statusNew            = "new"     // the delegation request has been submitted but it is not yet been worked on
	statusWorkInProgress = "wip"     // the delegation request is being worked on
	statusFailed         = "failed"  // the delegation request has failed (see error-message for the details)
	statusSuccess        = "success" // the delegation request has succeeded (see certificate-url)
)

const (
	maxDuration     time.Duration = time.Hour * 24 * 365 / time.Second
	maxCertLifetime time.Duration = time.Hour * 24 * 7 / time.Second
)

const (
	// TODO(tho) read from configuration
	dbFileName string = "./TODO.db"
)

// DelegationStore encapsulates the store logics.  In this case, we use a SQLite3 backend.
type DelegationStore struct {
	db *sql.DB
}

// NewDelegationStore creates a new instance of DelegationStore
func NewDelegationStore() (*DelegationStore, error) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	err = createDelegationTable(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &DelegationStore{db: db}, nil
}

func createDelegationTable(db *sql.DB) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS delegation (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		certURL         TEXT,
		created         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_update     DATETIME DEFAULT NULL,
		expires         DATETIME DEFAULT NULL,
		completed       DATETIME DEFAULT NULL,
		status          TEXT NOT NULL DEFAULT "new",
		details         TEXT,
		csr             BLOB,
		cert_lifetime   INTEGER,
		duration        INTEGER

		CHECK (status IN ("new", "wip", "done", "failed")),
		CHECK (cert_lifetime > 0),
		CHECK (duration > cert_lifetime)
	)
	`
	_, err := db.Exec(sqlQuery)

	return err
}

func normalizeDuration(dPtr *time.Duration, dMax time.Duration) time.Duration {
	if dPtr == nil || *dPtr > dMax || *dPtr < 0 {
		return dMax
	}
	return *dPtr
}

func sanitizeInDelegation(d *models.Delegation) (models.Delegation, error) {
	// CSR is mandatory
	if d.CSR == nil {
		return models.Delegation{}, apperrors.ErrDelegationMissingParameter
	}

	// If unspecified, or out-of-range, duration is set to 365 days
	// TODO(tho) read maxDuration from configuration
	duration := normalizeDuration(d.Duration, maxDuration)

	// If unspecified, or out-of-range, certificate lifetime is set to 7 days
	// TODO(tho) read maxCertLifetime from configuration
	certLifetime := normalizeDuration(d.CertLifetime, maxCertLifetime)

	if certLifetime > duration {
		return models.Delegation{}, apperrors.ErrDelegationBadParameter
	}

	return models.Delegation{
		CSR:          d.CSR,
		Duration:     &duration,
		CertLifetime: &certLifetime,
		Status:       statusNew,
	}, nil
}

func delegationIDToString(id DelegationID) string {
	return strconv.FormatUint(uint64(id), 10)
}

func stringToDelegationID(id string) (DelegationID, error) {
	i, err := strconv.ParseUint(id, 10, 64)
	return DelegationID(i), err
}

// saveDelegationToStore assumes that the supplied delegation has been already sanitised
// by sanitizeInDelegation()
func (us *DelegationStore) saveDelegationToStore(d models.Delegation) (*models.Delegation, error) {
	// Save the id into the delegation request and add a creation timestamp
	d.CreationDate = time.Now()

	// Save the delegation request into the sqlite backend
	sqlQuery := "INSERT INTO delegation(csr, duration, cert_lifetime) VALUES(?, ?, ?)"

	stmt, err := us.db.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(d.CSR, d.Duration, d.CertLifetime)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	d.ID = delegationIDToString(DelegationID(id))

	return &d, nil
}

// AddDelegation tries to add a new Delegation record and returns
// the /delegationID/ corresponding to the newly created object
func (us *DelegationStore) AddDelegation(delegation *models.Delegation) (*models.Delegation, error) {
	// Sanitize input
	saneDelegation, err := sanitizeInDelegation(delegation)
	if err != nil {
		return &models.Delegation{}, err
	}

	return us.saveDelegationToStore(saneDelegation)
}

// ReadDelegation tries to retrieve the Delegation record associated with the supplied /id/
// If no Delegation can be found, an apperrors.DelegationUnknown error is returned.
func (us *DelegationStore) ReadDelegation(id string) (*models.Delegation, error) {
	i, err := stringToDelegationID(id)
	if err != nil {
		return nil, apperrors.ErrDelegationUnknown
	}

	//r, ok := us.delegations[i] // TODO(tho) replace with getDelegationById()
	d, err := us.getDelegationByID(i)
	if err != nil {
		return nil, apperrors.ErrDelegationUnknown
	}

	return d, nil
}

// Return the (unique) identifier associated to the added record, or the empty
// string on error
func (us *DelegationStore) getDelegationByID(id DelegationID) (*models.Delegation, error) {
	sqlQuery := `
	SELECT id,
	       certURL,
		   created,  
		   last_update,
	       expires,
		   completed,
		   status,
		   details,
		   csr,
	       cert_lifetime,
	       duration
	  FROM delegation
	 WHERE id = ?
	`
	var d models.Delegation

	err := us.db.QueryRow(sqlQuery, id).
		Scan(&d.ID,
			&d.CertURL,
			&d.CreationDate,
			&d.LastUpdate,
			&d.ExpirationDate,
			&d.CompletionDate,
			&d.Status,
			&d.Details,
			&d.CSR,
			&d.CertLifetime,
			&d.Duration)

	if err != nil {
		return nil, err
	}

	return &d, nil
}
