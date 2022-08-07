package repository

import (
	"aspro/internal/domain/db/storage"
	"aspro/pkg/client/postgresql"
	"aspro/pkg/client/postgresql/model"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PersonPostgres struct {
	db *sqlx.DB
}

func NewPersonPostgres(db *sqlx.DB) *PersonPostgres {
	return &PersonPostgres{db: db}
}

func (r *PersonPostgres) Create(ctx context.Context, person *model.Person) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	personQuery := fmt.Sprintf("INSERT INTO %s (name, age) VALUES ($1, $2) RETURNING id", postgresql.PersonTable)
	if err := tx.QueryRow(personQuery, ctx, person.Name, 123).Scan(person.Id); err != nil {
		tx.Rollback()
	}

	var userId int
	userPersonQuery := fmt.Sprintf("INSERT INTO %s (user_id, person_id) VALUES ($1, $2)", postgresql.UsersPersonTable)
	_, err = tx.Exec(userPersonQuery, userId, person)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *PersonPostgres) RecieveAll(ctx context.Context, userId int, sortOptions storage.SortOptions) (p []model.Person, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	//for sorting
	//qb := sq.Select("id, name, age").From("person")
	//if sortOptions.Field != "" && sortOptions.Order != "" {
	//	qb = qb.OrderBy(fmt.Sprintf("%s %s", sortOptions.Field, sortOptions.Order))
	//}
	//query, i, err := qb.ToSql()
	//if err != nil {
	//	return nil, err
	//}

	query := fmt.Sprintf("SELECT id, name FROM %s", postgresql.PersonTable)
	row, err := tx.Query(query, ctx)
	if err != nil {
		return nil, err
	}

	people := make([]model.Person, 0)

	for row.Next() {
		var prs model.Person

		err = row.Scan(&prs.Id, &prs.Name, prs.Age, prs.CreatedAt)
		if err != nil {
			return nil, err
		}

		people = append(people, prs)
	}

	if err = row.Err(); err != nil {
		return nil, err
	}
	return people, nil
}

func (r *PersonPostgres) RecieveById(ctx context.Context, id int) (model.Person, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return model.Person{}, err
	}

	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id = $1", postgresql.PersonTable)
	var prs model.Person
	err = tx.QueryRow(query, ctx, id).Scan(&prs.Id, &prs.Name)
	if err != nil {
		return model.Person{}, err
	}
	return prs, nil
}

func (r *PersonPostgres) Update(ctx context.Context, person model.Person) error {
	query := fmt.Sprintf("UPDATE %s SET %s FROM %s WHERE id = %d", postgresql.PersonTable, person.Name, postgresql.UsersPersonTable, person.Id)
	_, err := r.db.Exec(query, person.Name)
	return err
}

func (r *PersonPostgres) Delete(ctx context.Context, userId, personId int) error {
	query := fmt.Sprintf("DELETE FROM %s USING %s WHERE id = user_id = $1 AND person_id = $2", postgresql.PersonTable, postgresql.UsersPersonTable)
	_, err := r.db.Exec(query, userId, personId)
	return err
}
