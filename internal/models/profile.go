package models

type Profile struct {
	UUID   string
	Name   string
	Orders []*Order
}
