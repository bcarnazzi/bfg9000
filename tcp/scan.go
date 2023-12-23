/*
Copyright © 2023 Bruno Carnazzi <bcarnazzi@gmail.com>
*/
package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/thediveo/netdb"
)

var HostsSpec string
var PortsSpec string
var Workers int
var Timeout int
var Scanners int

var portsList []int
var hostsList []string
var bar *progressbar.ProgressBar

var mu sync.Mutex
var hwg sync.WaitGroup
var scanMap = map[string]string{}

func initPorts() {
	csv := strings.Split(PortsSpec, ",")
	for _, e := range csv {
		rng := strings.Split(e, "-")
		switch len(rng) {
		case 1:
			p, err := strconv.Atoi(rng[0])
			if err != nil {
				panic(err)
			}
			portsList = append(portsList, p)
		case 2:
			var start, end int
			var err error

			if rng[0] == "" {
				start, err = 1, nil
			} else {
				start, err = strconv.Atoi(rng[0])
			}
			if err != nil {
				panic(err)
			}

			if rng[1] == "" {
				end, err = 65535, nil
			} else {
				end, err = strconv.Atoi(rng[1])
			}
			if err != nil {
				panic(err)
			}

			for i := start; i <= end; i++ {
				portsList = append(portsList, i)
			}
		default:
			panic("wrong post specification")
		}
	}
}

func initCidrTarget() {
	h, err := Hosts(HostsSpec)
	if err != nil {
		panic(err)
	}
	hostsList = h
	if len(h) < Scanners {
		Scanners = len(hostsList)
	}
}

func tcpWorker(host string, ports chan int, wg *sync.WaitGroup) {
	for i := range ports {
		p := strconv.Itoa(i)
		time.Sleep(time.Duration(rand.Intn(Timeout/4)) * time.Millisecond)
		c, err := net.DialTimeout("tcp", host+":"+p, time.Duration(Timeout*int(time.Millisecond)))
		if err != nil {
			bar.Add(1)
			wg.Done()
			continue
		}

		srvName := "unknown"
		srv := netdb.ServiceByPort(i, "tcp")
		if srv != nil {
			srvName = srv.Name
		}

		s := fmt.Sprintf("%-10s %-10s %-20s", p+"/tcp", "open", srvName)
		mu.Lock()
		scanMap[fmt.Sprintf("%s:%s", host, p)] = s
		mu.Unlock()

		c.Close()
		bar.Add(1)
		wg.Done()
	}
}

func TcpScanHost(host string, limiter chan int) {
	var wg sync.WaitGroup

	ports := make(chan int, Workers)
	for i := 0; i < Workers; i++ {
		go tcpWorker(host, ports, &wg)
	}

	for _, p := range portsList {
		wg.Add(1)
		ports <- p
	}
	wg.Wait()
	close(ports)
	hwg.Done()
	<-limiter

}

func TcpConnectScan() {

	initPorts()
	initCidrTarget()

	bar = progressbar.Default(int64(len(portsList) * len(hostsList)))

	scanChan := make(chan int, Scanners)

	for _, h := range hostsList {
		hwg.Add(1)
		scanChan <- 1
		go TcpScanHost(h, scanChan)
	}
	hwg.Wait()

	for _, h := range hostsList {
		var hs strings.Builder
		var up = false
		fmt.Fprintln(&hs, "\nScan report for", h)
		fmt.Fprintf(&hs, "%-10s %-10s %-10s\n", "PORT", "STATE", "SERVICE")

		for p := 1; p <= 65535; p++ {
			v := scanMap[fmt.Sprintf("%s:%d", h, p)]
			if v != "" {
				up = true
				fmt.Fprintln(&hs, v)
			}
		}
		if up {
			fmt.Print(hs.String())
			fqdn, err := net.LookupAddr(h)
			if err == nil {
				fmt.Println("Hostname is", fqdn[0])
			}
		}
	}

}
