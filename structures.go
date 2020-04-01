package toolucid

import (
	"time"
)

type loc struct {
	Name   string    `json:"type,omitempty"`
}

type relation struct {
	Type   string    `json:"type,omitempty"`
}

// If omitempty is not set, then edges with empty values (0 for int/float, "" for string, false
// for bool) would be created for values not specified explicitly.

type Person struct {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Relationship  relation   `json:"relation,omitempty"`
}

type Dream struct {
	Date 		*time.Time `json:"dob,omitempty"`
	Actors 		[]Person   `json:"friend,omitempty"`
	Location 	loc        `json:"loc,omitempty"`
	Intimity	int        `json:"intimity,omitempty"`
	Text		string	   `json:"text,omitempty"`
}