// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	fmt "fmt"
	io "io"
	strconv "strconv"
)

type MoodInput struct {
	Score   int     `json:"score"`
	Comment *string `json:"comment"`
}
type UserInput struct {
	Sex      Sex    `json:"sex"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Permission string

const (
	PermissionDebug Permission = "DEBUG"
)

func (e Permission) IsValid() bool {
	switch e {
	case PermissionDebug:
		return true
	}
	return false
}

func (e Permission) String() string {
	return string(e)
}

func (e *Permission) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Permission(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Permission", str)
	}
	return nil
}

func (e Permission) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Sex string

const (
	SexMale    Sex = "MALE"
	SexFemale  Sex = "FEMALE"
	SexUnknown Sex = "UNKNOWN"
)

func (e Sex) IsValid() bool {
	switch e {
	case SexMale, SexFemale, SexUnknown:
		return true
	}
	return false
}

func (e Sex) String() string {
	return string(e)
}

func (e *Sex) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Sex(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Sex", str)
	}
	return nil
}

func (e Sex) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
