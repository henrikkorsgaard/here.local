package context

type CertService struct{}

type TLSCertBytes struct {
	Bytes []byte
}

func (r *CertService) GetCertBytes(msg string, reply *TLSCertBytes) error {
	//reply.Bytes = configuration.

	return nil
}

func setupRPCsimpleServer() {

}
