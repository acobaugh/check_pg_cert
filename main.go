package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s postgresql://<host>[:port]\n", os.Args[0])
		os.Exit(1)
	}

	u, err := url.Parse(os.Args[1])
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
	tlsConf.InsecureSkipVerify = false
	tlsConf.ServerName = u.Hostname()

	client := tls.Client(cn.c, &tlsConf)

	err = client.Handshake()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	certs := client.ConnectionState().PeerCertificates
	for i := 0; i < len(certs); i++ {
		fmt.Printf("%d:\n  Subject: %s\n", i, certs[i].Subject)
		fmt.Printf("  Issuer: %s\n", certs[i].Issuer)
		fmt.Printf("  NotAfter: %s\n", certs[i].NotAfter)
	}

}
