package model

import "time"

type Person struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json"created_at"`
}

type Device struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Category string   `json:"category"`
	People   []Person `json:"people"`
}

func (d *Device) Eq() Device {
	dv := Device{
		ID:   d.ID,
		Name: d.Name,
	}
	dv.Age = int(d.Age)

	return dv
}

type UserPerson struct {
	Id       int
	UserId   int
	PersonId int
}
