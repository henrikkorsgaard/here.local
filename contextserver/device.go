package context

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"../proximity/proximity"
	"golang.org/x/crypto/bcrypt"
)

type Device struct {
	ID     string //macHashed with the device specific salt
	hash   string //hashed mac address
	Signal int    //device signal strength
}

var (
	Salt string
	src  rand.Source
)

func init() {
	rand.Seed(time.Now().UnixNano())
	Salt = generateSalt(32) //does not need to be long, as it is used to salt mac address to counter the finite address space (see https://en.wikipedia.org/wiki/MAC_address_anonymization)
}

func UpsertRawDevice(rd proximity.RawDevice) {
	/*
		Because of the finite mac address space, we need to salt the address
		before upserting into the database.

	*/
	mac := Salt + rd.Mac
	hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	device := Device{hash: string(hash), Signal: rd.Signal}
	fmt.Println(device)
}

func generateSalt(n int) string {
	//see https://stackoverflow.com/a/31832326
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}
