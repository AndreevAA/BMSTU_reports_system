package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/label"
	"neatly/internal/model/report"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
)

const (
	labelsTable        = "labels"
	reportsLabelsTable = "labels_reports"
	usersLabelsTable   = "users_labels"
)

type LabelPostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewLabelPostgres(client *psqlclient.Client, logger logging.Logger) *LabelPostgres {
	return &LabelPostgres{db: client.DB, logger: logger}
}

func (r *LabelPostgres) Create(userID, reportID int, t *label.Label) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createLabelQuery := fmt.Sprintf(
		`INSERT INTO %s AS t (name, department) VALUES ($1, $2) RETURNING id`,
		labelsTable)

	r.logger.Infof("Label with id %v created", t.ID)

	row := r.db.QueryRow(createLabelQuery, t.Name, t.Department)
	err = row.Scan(&t.ID)

	if err != nil {
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return &report.CanNotCreateReportErr{}
		}
		return err
	}

	r.logger.Infof("Connecting label with id %v and accounts with id with id %v", t.ID, userID)
	userLabelQuery := fmt.Sprintf(
		`INSERT INTO %s
    			(users_id, labels_id)
				SELECT $1, $2
				WHERE
    			NOT EXISTS (
    			    SELECT users_id, labels_id FROM users_labels WHERE users_id = $1 AND labels_id = $2
    			);`, usersLabelsTable)
	_, err = tx.Exec(userLabelQuery, userID, t.ID)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return &report.CanNotCreateReportErr{}
		}
		return err
	}
	return tx.Commit()
}

func (r *LabelPostgres) Assign(labelID, reportID, userID int) error {
	r.logger.Infof("Assigning label with id %v to report with id with id %v", labelID, reportID)
	assignLabelQuery := fmt.Sprintf(
		`INSERT INTO %s (reports_id, labels_id) VALUES ($1, $2)`, reportsLabelsTable)
	_, err := r.db.Exec(assignLabelQuery, reportID, labelID)
	if err != nil {
		r.logger.Info(err)
		return err
	}

	return nil
}

func (r *LabelPostgres) GetAll(userID int) ([]label.Label, error) {
	var labels []label.Label
	labels = make([]label.Label, 0)

	query := fmt.Sprintf(`SELECT labels_id AS id, name, department FROM
								%s t INNER JOIN %s ut ON ut.labels_id = t.id  WHERE
								ut.users_id = $1`, labelsTable, usersLabelsTable)

	err := r.db.Select(&labels, query, userID)
	if err != nil {
		r.logger.Info(err)
	}
	return labels, err
}

func (r *LabelPostgres) GetAllByReport(userID, reportID int) ([]label.Label, error) {
	var labels []label.Label
	labels = make([]label.Label, 0)

	query := fmt.Sprintf(`SELECT t.id AS id, name, department FROM %s t
    							INNER JOIN %s ut ON ut.labels_id = t.id
    							INNER JOIN %s nt on t.id = nt.labels_id
    							WHERE users_id = $1 AND reports_id = $2`,
		labelsTable, usersLabelsTable, reportsLabelsTable)

	err := r.db.Select(&labels, query, userID, reportID)
	if err != nil {
		r.logger.Info(err)
	}
	return labels, err
}

func (r *LabelPostgres) GetOne(userID, labelID int) (label.Label, error) {
	var t label.Label

	query := fmt.Sprintf(`SELECT t.id AS id, name, department FROM %s t
    							INNER JOIN %s ut ON ut.labels_id = t.id
    							INNER JOIN %s nt on t.id = nt.labels_id
    							WHERE users_id = $1 AND t.id = $2`,
		labelsTable, usersLabelsTable, reportsLabelsTable)

	err := r.db.Get(&t, query, userID, labelID)
	if err != nil {
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return t, &label.LabelNotFoundErr{}
		}
	}
	return t, err
}

func (r *LabelPostgres) Delete(userID, labelID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s t USING %s ut WHERE 
              t.id = ut.labels_id AND ut.users_id = $1 AND ut.labels_id = $2`,
		labelsTable, usersLabelsTable)
	_, err := r.db.Exec(query, userID, labelID)

	return err
}

func (r *LabelPostgres) Update(userID, labelID int, t label.Label) error {
	query := fmt.Sprintf(
		`UPDATE %s t SET 
                name=$1, department=$2 FROM
                %s ut WHERE t.id = ut.labels_id AND 
				ut.labels_id = $3 AND ut.users_id = $4`,
		labelsTable, usersLabelsTable)
	_, err := r.db.Exec(query, t.Name, t.Department, labelID, userID)

	return err
}

func (r *LabelPostgres) Detach(userID, labelID, reportID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s USING %s ut WHERE
            	labels_reports.labels_id = ut.labels_id AND ut.users_id = $1 AND ut.labels_id = $2 AND reports_id = $3`,
		reportsLabelsTable, usersLabelsTable)
	_, err := r.db.Exec(query, userID, labelID, reportID)

	return err
}
