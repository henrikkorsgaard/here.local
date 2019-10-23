package device

type Location struct {
	MAC     string
	IP      string
	Name    string
	Devices map[string]int
}
