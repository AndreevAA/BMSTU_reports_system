package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/report"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
	"time"
)

const (
	reportsTable      = "reports"
	reportsBodyTable  = "reports_body"
	usersReportsTable = "users_reports"
)

type ReportPostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewReportPostgres(client *psqlclient.Client, logger logging.Logger) *ReportPostgres {
	return &ReportPostgres{
		db:     client.DB,
		logger: logger,
	}
}

func (r *ReportPostgres) Create(userID int, n *report.Report) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Info(err)
		return &report.CanNotCreateReportErr{}
	}

	createReportQuery := fmt.Sprintf(`
	INSERT INTO %s (header, short_body, department, edited)
	VALUES ($1, $2, $3, $4) RETURNING id`, reportsTable)
	row := tx.QueryRow(createReportQuery, n.Header, n.ShortBody, n.Department, time.Now())
	if err := row.Scan(&n.ID); err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return &report.CanNotCreateReportErr{}
	}
	createReportBodyQuery := fmt.Sprintf("INSERT INTO %s (id, body) VALUES ($1, $2)", reportsBodyTable)
	_, err = tx.Exec(createReportBodyQuery, n.ID, n.Body)
	if err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return &report.CanNotCreateReportErr{}
	}
	createUsersReportQuery := fmt.Sprintf("INSERT INTO %s (users_id, reports_id) VALUES ($1, $2)", usersReportsTable)
	_, err = tx.Exec(createUsersReportQuery, userID, n.ID)
	if err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return &report.CanNotCreateReportErr{}
	}

	return tx.Commit()
}

func (r *ReportPostgres) GetAll(userID int) ([]report.Report, error) {
	var reports []report.Report
	reports = make([]report.Report, 0)

	getReportsQuery := fmt.Sprintf(
		`SELECT n.id, n.header, n.short_body, n.department, n.edited FROM %s n
    			JOIN %s un ON n.id = un.reports_id
    			WHERE un.users_id = $1`,
		reportsTable,
		usersReportsTable,
	)

	err := r.db.Select(&reports, getReportsQuery, userID)
	if err != nil {
		r.logger.Info(err)
		return reports, err
	}

	return reports, err
}

func (r *ReportPostgres) GetOne(userID, reportID int) (report.Report, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return report.Report{}, err
	}
	var n report.Report

	selectReportQuery := fmt.Sprintf(
		`SELECT n.id, n.header, n.short_body, n.department, n.edited FROM
				%s n JOIN %s un ON n.id = un.reports_id
				WHERE un.users_id = $1 AND un.reports_id = $2`,
		reportsTable,
		usersReportsTable,
	)

	err = r.db.Get(&n, selectReportQuery, userID, reportID)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return report.Report{}, &report.ReportNotFoundErr{}
		}
		return report.Report{}, err
	}

	selectBodyQuery := fmt.Sprintf(
		`SELECT nb.body FROM %s nb JOIN %s n ON nb.id = n.id
				WHERE n.id = $1`,
		reportsBodyTable,
		reportsTable,
	)

	err = r.db.Get(&n.Body, selectBodyQuery, reportID)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return report.Report{}, &report.ReportNotFoundErr{}
		}
	}

	return n, tx.Commit()
}

func (r *ReportPostgres) Delete(userID, reportID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s n USING %s un WHERE 
              n.id = un.reports_id AND un.users_id = $1 AND un.reports_id = $2`,
		reportsTable, usersReportsTable)
	_, err := r.db.Exec(query, userID, reportID)

	return err
}

func (r *ReportPostgres) Update(userID int, n report.Report) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	reportQuery := fmt.Sprintf(
		`UPDATE %s n SET 
                header=$1, short_body=$2, department = $3, edited=$4 FROM
                %s un WHERE n.id = un.reports_id AND 
				un.reports_id = $5 AND un.users_id = $6`,
		reportsTable, usersReportsTable)
	_, err = r.db.Exec(
		reportQuery,
		n.Header,
		n.ShortBody,
		n.Department,
		time.Now().UTC().Format(time.RFC3339),
		n.ID,
		userID,
	)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return &report.ReportNotFoundErr{}
		}
		return err
	}

	bodyQuery := fmt.Sprintf(
		`UPDATE %s nb SET body=$2 WHERE nb.id = $1`,
		reportsBodyTable)
	_, err = r.db.Exec(bodyQuery, n.ID, n.Body)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return &report.ReportNotFoundErr{}
		}
	}

	return tx.Commit()
}
