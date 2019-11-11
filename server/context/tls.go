package context

type TLSService struct{}

type TLSPemBytes struct {
	PemBytes []byte
}

func (r *TLSService) RegisterNodeAndGetPemBytes(location string, reply *TLSPemBytes) error {
	//reply.Bytes = configuration.

	return nil
}

func setupTLSService() {
	//generate the tls stuff
	//set up service

	//reply to requests
}
