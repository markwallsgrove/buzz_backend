package domain

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Gender int

const (
	UnknownGender = iota
	Male
	Female
)

var stringToGender map[string]Gender = map[string]Gender{
	"Unknown": UnknownGender,
	"Male":    Male,
	"Female":  Female,
}

var Genders []Gender = []Gender{Male, Female}

func (g Gender) String() string {
	return [...]string{"Unknown", "Male", "Female"}[g]
}

// StringToGender convert a string into a gender enum
func StringToGender(gender string) Gender {
	g, ok := stringToGender[gender]
	if !ok {
		return UnknownGender
	}

	return g
}

type User struct {
	ID           int      `json:"id"`
	Email        string   `json:"email"`
	PasswordHash []byte   `json:"password_hash"`
	Name         string   `json:"name"`
	Gender       Gender   `json:"gender"`
	Age          int      `json:"age"`
	Location     Location `json:"location"`
}

type UserProfile struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

// See https://gorm.io/docs/create.html#Create-From-SQL-Expression-x2F-Context-Valuer for more details

// Location coordinates
type Location struct {
	X, Y float64
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	fmt.Println(fmt.Sprintf("%+v", v))
	return errors.New("oh no!")
}

func (loc Location) GormDataType() string {
	return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%f %f)", loc.X, loc.Y)},
	}
}
