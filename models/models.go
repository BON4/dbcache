package models

type Item struct {
	Id          string
	TransportId string
	Number      string
}

type Transport struct {
	Id     string
	Number string
}

type TransportItemView struct {
	Transport Transport
	Items     []Item
}
