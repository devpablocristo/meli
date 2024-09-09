package graphqlkyc

import "time"

type response struct {
	User user `json:"user"`
}

type user struct {
	Identification identification `json:"identification"`
	DateCreated    time.Time      `json:"date_created"`
}

type identification struct {
	ID string `json:"id"`
}
