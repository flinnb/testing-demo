package user

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"testing-demo/internal/db"
	"testing-demo/internal/logging"
)

type UserRequest struct {
	ID        int64  `json:"user_id"`
	Name      string `json:"name"`
	DOB       string `json:"DOB"`
	CreatedOn int64  `json:"created_on"`
}

type UserResponse struct {
	ID        int64     `json:"user_id"`
	Name      string    `json:"name"`
	DOB       string    `json:"DOB"`
	CreatedOn time.Time `json:"created_on"`
}

func (req *UserRequest) ToResponse() (res *UserResponse, err error) {
	logger := logging.GetLogger()
	logger.Infof("%#v", req)
	res = &UserResponse{
		ID:   req.ID,
		Name: req.Name,
	}
	// Deliberately ignoreing error here, because it is so unlikely.
	newYork, _ := time.LoadLocation("America/New_York")
	dob, dobErr := time.Parse("2006-01-02", req.DOB)
	if dobErr != nil {
		logger.Error(dobErr)
		err = errors.New(fmt.Sprintf("'%s' can't be parsed as a date", req.DOB))
		return
	}
	res.DOB = dob.Format("Monday")
	res.CreatedOn = time.Unix(req.CreatedOn, 0).In(newYork)

	return
}

type UserProfile struct {
	ID        int
	FirstName string
	LastName  string
	Username  string
	City      string
	ZipCode   string
	CreatedOn time.Time
	UpdatedOn time.Time
}

type PasswordHistory struct {
	ID        int
	UserID    int
	Password  string
	CreatedOn time.Time
	Active    bool
}

func (up *UserProfile) Create() error {
	tx, txErr := db.Pool().Begin()
	if txErr != nil {
		return txErr
	}
	row := tx.QueryRow(`INSERT INTO user_profiles (
	first_name,
	last_name,
	username,
	city,
	zip_code,
	created_on,
	updated_on
	)
	VALUES ($1, $2, $3, $4, $5, now(), now())
	RETURNING id, created_on, updated_on
	`,
		up.FirstName,
		up.LastName,
		up.Username,
		up.City,
		up.ZipCode,
	)
	sErr := row.Scan(&up.ID, &up.CreatedOn, &up.UpdatedOn)
	if sErr != nil {
		tx.Rollback()
		return sErr
	}
	tx.Commit()
	return nil

}

func (up *UserProfile) Delete() error {
	_, err := db.Pool().Exec("DELETE FROM user_profiles WHERE id = $1", up.ID)
	if err != nil {
		return err
	}
	return nil
}

func (up *UserProfile) PasswordHistory() ([]*PasswordHistory, error) {
	ph := make([]*PasswordHistory, 0)

	rows, err := db.Pool().Query(`SELECT id, user_id, created_on, is_active
	FROM user_password_history WHERE user_id = $1`, up.ID)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &PasswordHistory{}
		scanErr := rows.Scan(&item.ID, &item.UserID, &item.CreatedOn, &item.Active)
		if scanErr != nil {
			return nil, scanErr
		}
		ph = append(ph, item)
	}
	return ph, nil
}

// This should never return more than one record, but it is useful for testing purposes
func (up *UserProfile) ActivePasswordHistory() ([]*PasswordHistory, error) {
	ph := make([]*PasswordHistory, 0)

	rows, err := db.Pool().Query(`SELECT id, user_id, created_on, is_active
	FROM user_password_history WHERE user_id = $1 AND is_active`, up.ID)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &PasswordHistory{}
		scanErr := rows.Scan(&item.ID, &item.UserID, &item.CreatedOn, &item.Active)
		if scanErr != nil {
			return nil, scanErr
		}
		ph = append(ph, item)
	}
	return ph, nil
}

func (ph *PasswordHistory) Create() error {
	tx, txErr := db.Pool().Begin()
	if txErr != nil {
		return txErr
	}
	row1 := tx.QueryRow(`SELECT id, user_id, password_hash, created_on, is_active
	FROM user_password_history WHERE user_id = $1 AND is_active`, ph.UserID)
	oldPh := PasswordHistory{}
	oldPhErr := row1.Scan(&oldPh.ID, &oldPh.UserID, &oldPh.Password, &oldPh.CreatedOn, &oldPh.Active)
	if oldPhErr != nil && oldPhErr != sql.ErrNoRows {
		tx.Rollback()
		return oldPhErr
	}
	pwByteArray, pwErr := bcrypt.GenerateFromPassword([]byte(ph.Password), bcrypt.DefaultCost)
	if pwErr != nil {
		return pwErr
	}
	ph.Password = string(pwByteArray)
	insertPh := tx.QueryRow(`INSERT INTO user_password_history (user_id, password_hash, created_on, is_active)
	VALUES ($1, $2, now(), true)
	RETURNING id, created_on, is_active`,
		ph.UserID,
		ph.Password,
	)
	insertPhErr := insertPh.Scan(&ph.ID, &ph.CreatedOn, &ph.Active)
	if insertPhErr != nil {
		tx.Rollback()
		return insertPhErr
	}
	if oldPh.ID != 0 {
		updateOldPh := tx.QueryRow(`UPDATE user_password_history SET is_active=false
		WHERE id = $1 RETURNING is_active`, oldPh.ID)
		updateOldPhErr := updateOldPh.Scan(&oldPh.Active)
		if updateOldPhErr != nil {
			tx.Rollback()
			return updateOldPhErr
		}
	}
	tx.Commit()
	return nil
}
