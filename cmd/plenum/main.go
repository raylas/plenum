package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/raylas/plenum"
)

var (
	bindAddress string
	bindPort    int
	members     string
	membership  []string
)

func init() {
	flag.StringVar(&bindAddress, "address", "127.0.0.1", "127.0.0.1")
	flag.IntVar(&bindPort, "port", 0, "1111")
	flag.StringVar(&members, "members", "", "127.0.0.1:1111,127.0.0.1:2222")
}

func main() {
	flag.Parse()

	if members != "" {
		membership = strings.Split(members, ",")
	}

	p, err := plenum.Convene(nil, bindAddress, bindPort, membership)
	if err != nil {
		log.Fatal(err)
	}
	go p.Procede()

	// Wait for, and act on leader designation
	go func(*plenum.Plenum) {
		for <-p.IsLeader {
			asciiLeader()
		}
	}(p)

	http.HandleFunc("/members", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, p.Membership())
	})
	http.HandleFunc("/leader", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, p.Consensus())
	})
	http.ListenAndServe(fmt.Sprintf("%s:%d", bindAddress, bindPort+2), nil)
}

func asciiLeader() {
	fmt.Println(`
	██╗     ███████╗ █████╗ ██████╗ ███████╗██████╗ 
	██║     ██╔════╝██╔══██╗██╔══██╗██╔════╝██╔══██╗
	██║     █████╗  ███████║██║  ██║█████╗  ██████╔╝
	██║     ██╔══╝  ██╔══██║██║  ██║██╔══╝  ██╔══██╗
	███████╗███████╗██║  ██║██████╔╝███████╗██║  ██║
	╚══════╝╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝`)
}
