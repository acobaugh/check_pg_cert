package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
)

var Usage = func() {
	fmt.Printf("Usage: %s postgresql://<host>:<port> [options]\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var insecureSkipVerify bool
	flag.BoolVar(&insecureSkipVerify, "insecure-skip-verify", false, "Skip peer certificate verification")
	flag.Parse()

	if len(flag.Args()) != 1 {
		Usage()
		os.Exit(1)
	}

	u, err := url.Parse(flag.Args()[0])
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	pgHost := u.Hostname()
	pgPort := u.Port()
	if pgPort == "" {
		pgPort = "5432"
	}
	addr := fmt.Sprintf("%s:%s", pgHost, pgPort)

	c, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	cn := &conn{}
	cn.c = c

	w := cn.writeBuf(0)
	w.int32(80877103)
	cn.sendStartupPacket(w)

	b := cn.scratch[:1]
	_, err = io.ReadFull(cn.c, b)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if b[0] != 'S' {
		fmt.Println("SSL not supported")
		os.Exit(1)
	}

	tlsConf := tls.Config{}
	tlsConf.InsecureSkipVerify = insecureSkipVerify
	tlsConf.ServerName = u.Hostname()

	client := tls.Client(cn.c, &tlsConf)

	err = client.Handshake()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	certs := client.ConnectionState().PeerCertificates
	for i := 0; i < len(certs); i++ {
		fmt.Printf("%d: Subject: %s\n", i, certs[i].Subject)
		fmt.Print("   DNSNames: ")
		for j := 0; j < len(certs[i].DNSNames); j++ {
			fmt.Printf("%s, ", certs[i].DNSNames[j])
		}
		fmt.Println()
		fmt.Printf("   Issuer: %s\n", certs[i].Issuer)
		fmt.Printf("   NotAfter: %s\n", certs[i].NotAfter)
	}

}
