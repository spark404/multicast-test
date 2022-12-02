package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		multicast          = flag.Bool("multicast", true, "Enable mulitcast listener")
		multicastInterface = flag.String("mcast-if", "", "Bind interface for multicast listener")
	)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	destination := "0.0.0.0:4242"
	if *multicast {
		destination = "239.42.42.42:4242"
	}

	addr, err := net.ResolveUDPAddr("udp4", destination)
	if err != nil {
		log.WithError(err).Fatalln("unable to resolve '%s'", destination)
	}

	var intf *net.Interface = nil
	if multicastInterface != nil && *multicastInterface != "" {
		intf, err = net.InterfaceByName(*multicastInterface)
		if err != nil {
			log.Errorf("No interface %s found", *multicastInterface)
			intf = nil
		}
	}

	var conn *net.UDPConn
	if addr.IP.IsMulticast() {
		log.Infof("Starting multicast listener for destination %s on interface %v", addr, intf)
		conn, err = net.ListenMulticastUDP("udp4", intf, addr)
	} else {
		log.Infof("Starting listener for destination %s", addr)
		conn, err = net.ListenUDP("udp4", addr)
	}
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	go func() {
		router := mux.NewRouter()
		router.Use(mux.CORSMethodMiddleware(router))

		router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			forwardedFor := request.Header.Get("X-Forwarded-For")
			if forwardedFor != "" {
				log.Infof("%s (%s)>  %s %s", request.RemoteAddr, forwardedFor, request.Method, request.URL.String())
			} else {
				log.Infof("%s>  %s %s", request.RemoteAddr, request.Method, request.URL.String())
			}

			writer.WriteHeader(200)
			_, _ = writer.Write([]byte("Hello World"))
		})

		log.Infof("Starting HTTP/TCP server on port 5001")
		err := http.ListenAndServe(fmt.Sprintf(":%d", 5001), router)
		if err != nil {
			log.WithError(err).Fatalf("failed to start http listener")
		}
	}()

	go func() {
		log.Infof("Listening for UDP traffic on %s", conn.LocalAddr().String())
		buffer := make([]byte, 1024)
		for {
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				log.WithError(err).Error("read failed")
				break
			}
			log.Infof("%s> %s", addr.String(), string(buffer[0:n]))
		}
	}()

	<-sigs

}
