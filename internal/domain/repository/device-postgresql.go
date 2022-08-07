package repository

import (
	"aspro/pkg/client/postgresql"
	"aspro/pkg/client/postgresql/model"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DevicePostgres struct {
	db *sqlx.DB
}

func NewDevicePostgres(db *sqlx.DB) *DevicePostgres {
	return &DevicePostgres{db: db}
}

func (r *DevicePostgres) Create(ctx context.Context, deviceId, userId, personId int, device model.Device) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	deviceQuery := fmt.Sprintf("INSERT INTO %s (name, age, category) VALUES ($1. $2, $3) RETURNING id", postgresql.DeviceTable)

	row := tx.QueryRow(deviceQuery, ctx)
	err = row.Scan(&deviceId)
	if err != nil {
		return 0, err
	}

	personDeviceQuery := fmt.Sprintf("INSERT INTO %s (person_id, device_id) VALUES ($1, $2)", postgresql.PersonDeviceTable)
	_, err = tx.Exec(personDeviceQuery, personId, deviceId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return deviceId, tx.Commit()
}

func (r *DevicePostgres) FindAll(ctx context.Context, userId, personId int) (device []model.Device, err error) {
	query := fmt.Sprintf("SELECT id, name, age, category FORM %s", postgresql.DeviceTable)

	row, err := r.db.Query(query, ctx)
	if err != nil {
		return nil, err
	}

	devices := make([]model.Device, 0)

	for row.Next() {
		var dv model.Device

		err = row.Scan(&dv.ID, &dv.Name, &dv.Age, &dv.Category)
		if err != nil {
			return nil, err
		}
		devices = append(devices, dv.Eq())
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *DevicePostgres) FindById(ctx context.Context, userId, deviceId int) (model.Device, error) {
	query := fmt.Sprintf("SELECT id, name, age, category FORM %s WHERE id = $1", postgresql.DeviceTable)

	var dv model.Device
	err := r.db.QueryRow(query, ctx, deviceId).Scan(&dv.ID, &dv.Name, &dv.Age, &dv.Category)
	if err != nil {
		return model.Device{}, err
	}
	return dv, nil
}

func (r *DevicePostgres) Update(ctx context.Context, userId, deviceId int) error {
	query := fmt.Sprintf("UPDATE %s SET name=$2 WHERE id = $1", postgresql.DeviceTable)
	_, err := r.db.Exec(query, 1, ctx, deviceId)
	if err != nil {
		panic(err)
	}
	return err
}

func (r *DevicePostgres) Delete(ctx context.Context, userId, deviceId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", postgresql.DeviceTable)
	_, err := r.db.Exec(query, ctx, userId, deviceId)
	if err != nil {
		panic(err)
	}
	return err
}
