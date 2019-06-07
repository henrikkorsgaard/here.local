package proximity

type RawDevice struct {
	Mac    string
	Signal int
}

type NmapDevice struct {
	Mac      string
	Ip       string
	Hostname string
	Vendor   string
}
