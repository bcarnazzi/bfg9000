package scan

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

var ports chan int
var portsList []int

var hosts []string
var currentHost string

var mu sync.Mutex
var wg sync.WaitGroup
var scanMap = map[int]string{}

var bar *progressbar.ProgressBar

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
	hosts = h
}

func tcpWorker() {
	for i := range ports {
		p := strconv.Itoa(i)
		time.Sleep(time.Duration(rand.Intn(Timeout/4)) * time.Millisecond)
		c, err := net.DialTimeout("tcp", currentHost+":"+p, time.Duration(Timeout*int(time.Millisecond)))
		if err != nil {
			bar.Add(1)
			wg.Done()
			continue
		}

		srvName := "UNKNOWN"
		srv := netdb.ServiceByPort(i, "tcp")
		if srv != nil {
			srvName = srv.Name
		}

		s := fmt.Sprintf("%-10s %-10s %-20s", p+"/tcp", "open", srvName)
		mu.Lock()
		scanMap[i] = s
		mu.Unlock()

		c.Close()
		bar.Add(1)
		wg.Done()
	}
}

func TcpConnectScan() {

	initPorts()
	initCidrTarget()

	ports = make(chan int, Workers)
	for i := 0; i < Workers; i++ {
		go tcpWorker()
	}

	for _, h := range hosts {

		currentHost = h
		bar = progressbar.Default(int64(len(portsList)))
		for _, p := range portsList {
			wg.Add(1)
			ports <- p
		}

		wg.Wait()
		//close(ports)
		//ports = make(chan int, Workers)

		fmt.Println("\nScan report for", currentHost)
		fmt.Printf("%-10s %-10s %-10s\n", "PORT", "STATE", "SERVICE")
		for i := 1; i <= 65535; i++ {
			v := scanMap[i]
			if v != "" {
				fmt.Println(v)
			}
		}
		scanMap = make(map[int]string)

	}
	close(ports)

}
