package domain

import (
	"context"
	"encoding/binary"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/twpayne/go-geom/encoding/wkb"
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

// UserDistance user model with a the distance attribute
type UserDistance struct {
	User
	Distance float64 `json:"distance"`
}

type UserProfile struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Gender         string  `json:"gender"`
	Age            int     `json:"age"`
	DistanceFromMe float64 `json:"distanceFromMe"`
}

// See https://gorm.io/docs/create.html#Create-From-SQL-Expression-x2F-Context-Valuer for more details

// Location coordinates
type Location struct {
	X, Y float64
}

// Scan implements the sql.Scanner interface.
//
// Receive raw bytes from the sql driver which represents the value. The value
// is prefixed with prefix (well known text), which is followed by the coordinates.
// See https://stackoverflow.com/questions/60520863/working-with-spatial-data-with-gorm-and-mysql
// and https://dev.mysql.com/doc/refman/8.0/en/gis-data-formats.html for more details.
func (loc *Location) Scan(v interface{}) error {
	if v == nil {
		return nil
	}

	mysqlEncoding, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("did not scan: expected []byte but was %T", v)
	}

	var srid uint32 = binary.LittleEndian.Uint32(mysqlEncoding[0:4])

	var point wkb.Point
	if err := point.Scan(mysqlEncoding[4:]); err != nil {
		return err
	}

	point.SetSRID(int(srid))
	loc.X = point.Point.X()
	loc.Y = point.Point.Y()

	return nil
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
