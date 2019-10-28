package models

type Location struct {
	MAC     string
	IP      string
	Name    string
	Devices map[string]Device
	Signal  int
}
