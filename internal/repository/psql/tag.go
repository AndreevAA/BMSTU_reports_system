package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/report"
	"neatly/internal/model/tag"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
)

const (
	tagsTable      = "tags"
	reportsTagsTable = "tags_reports"
	usersTagsTable = "users_tags"
)

type TagPostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewTagPostgres(client *psqlclient.Client, logger logging.Logger) *TagPostgres {
	return &TagPostgres{db: client.DB, logger: logger}
}

func (r *TagPostgres) Create(userID, reportID int, t *tag.Tag) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createTagQuery := fmt.Sprintf(
		`INSERT INTO %s AS t (name, department) VALUES ($1, $2) RETURNING id`,
		tagsTable)

	r.logger.Infof("Tag with id %v created", t.ID)

	row := r.db.QueryRow(createTagQuery, t.Name, t.Department)
	err = row.Scan(&t.ID)

	if err != nil {
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return &report.CanNotCreateReportErr{}
		}
		return err
	}

	r.logger.Infof("Connecting tag with id %v and accounts with id with id %v", t.ID, userID)
	userTagQuery := fmt.Sprintf(
		`INSERT INTO %s
    			(users_id, tags_id)
				SELECT $1, $2
				WHERE
    			NOT EXISTS (
    			    SELECT users_id, tags_id FROM users_tags WHERE users_id = $1 AND tags_id = $2
    			);`, usersTagsTable)
	_, err = tx.Exec(userTagQuery, userID, t.ID)
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

func (r *TagPostgres) Assign(tagID, reportID, userID int) error {
	r.logger.Infof("Assigning tag with id %v to report with id with id %v", tagID, reportID)
	assignTagQuery := fmt.Sprintf(
		`INSERT INTO %s (reports_id, tags_id) VALUES ($1, $2)`, reportsTagsTable)
	_, err := r.db.Exec(assignTagQuery, reportID, tagID)
	if err != nil {
		r.logger.Info(err)
		return err
	}

	return nil
}

func (r *TagPostgres) GetAll(userID int) ([]tag.Tag, error) {
	var tags []tag.Tag
	tags = make([]tag.Tag, 0)

	query := fmt.Sprintf(`SELECT tags_id AS id, name, department FROM
								%s t INNER JOIN %s ut ON ut.tags_id = t.id  WHERE
								ut.users_id = $1`, tagsTable, usersTagsTable)

	err := r.db.Select(&tags, query, userID)
	if err != nil {
		r.logger.Info(err)
	}
	return tags, err
}

func (r *TagPostgres) GetAllByReport(userID, reportID int) ([]tag.Tag, error) {
	var tags []tag.Tag
	tags = make([]tag.Tag, 0)

	query := fmt.Sprintf(`SELECT t.id AS id, name, department FROM %s t
    							INNER JOIN %s ut ON ut.tags_id = t.id
    							INNER JOIN %s nt on t.id = nt.tags_id
    							WHERE users_id = $1 AND reports_id = $2`,
		tagsTable, usersTagsTable, reportsTagsTable)

	err := r.db.Select(&tags, query, userID, reportID)
	if err != nil {
		r.logger.Info(err)
	}
	return tags, err
}

func (r *TagPostgres) GetOne(userID, tagID int) (tag.Tag, error) {
	var t tag.Tag

	query := fmt.Sprintf(`SELECT t.id AS id, name, department FROM %s t
    							INNER JOIN %s ut ON ut.tags_id = t.id
    							INNER JOIN %s nt on t.id = nt.tags_id
    							WHERE users_id = $1 AND t.id = $2`,
		tagsTable, usersTagsTable, reportsTagsTable)

	err := r.db.Get(&t, query, userID, tagID)
	if err != nil {
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return t, &tag.TagNotFoundErr{}
		}
	}
	return t, err
}

func (r *TagPostgres) Delete(userID, tagID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s t USING %s ut WHERE 
              t.id = ut.tags_id AND ut.users_id = $1 AND ut.tags_id = $2`,
		tagsTable, usersTagsTable)
	_, err := r.db.Exec(query, userID, tagID)

	return err
}

func (r *TagPostgres) Update(userID, tagID int, t tag.Tag) error {
	query := fmt.Sprintf(
		`UPDATE %s t SET 
                name=$1, department=$2 FROM
                %s ut WHERE t.id = ut.tags_id AND 
				ut.tags_id = $3 AND ut.users_id = $4`,
		tagsTable, usersTagsTable)
	_, err := r.db.Exec(query, t.Name, t.Department, tagID, userID)

	return err
}

func (r *TagPostgres) Detach(userID, tagID, reportID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s USING %s ut WHERE
            	tags_reports.tags_id = ut.tags_id AND ut.users_id = $1 AND ut.tags_id = $2 AND reports_id = $3`,
		reportsTagsTable, usersTagsTable)
	_, err := r.db.Exec(query, userID, tagID, reportID)

	return err
}
