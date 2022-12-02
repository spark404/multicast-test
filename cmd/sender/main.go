package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	destination := "239.42.42.42:4242"
	//destination := "192.168.200.201:4242"
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	addr, err := net.ResolveUDPAddr("udp4", destination)
	if err != nil {
		log.WithError(err).Fatalln("unable to resolve '%s'", destination)
	}

	//laddr, err := net.ResolveUDPAddr("udp4", "192.168.200.254:0")
	//if err != nil {
	//	log.WithError(err).Fatalln("unable to resolve '%s'", destination)
	//}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.WithError(err).Fatalln("unable to resolve '%s'", destination)
	}

	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	go func() {
		for {
			_, err := conn.Write([]byte("hello world!"))
			if err != nil {
				log.WithError(err).Errorf("udp send failed")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	<-sigs
}
