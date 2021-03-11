package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := flag.Int("port", 3333, "")
	metricsPort := flag.Int("metrics-port", 8080, "")
	instance := flag.String("instance", "server-0", "")

	flag.Parse()

	server(*port, *metricsPort, *instance)
}

type ActiveConnections struct {
	ActiveConnections int
}

type Msg struct {
	Instance          string
	ActiveConnections int
}

func server(port int, metricsPort int, instance string) {
	promActiveConnections := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "demo",
		Subsystem: "server",
		Name:      "active_connections",
		Help:      "Number of blob storage operations waiting to be processed.",
	})
	prometheus.MustRegister(promActiveConnections)

	activeConnections := ActiveConnections{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func(activeConnections *ActiveConnections) {
		<-sigs
		for {
			fmt.Println(activeConnections.ActiveConnections)
			if activeConnections.ActiveConnections == 0 {
				os.Exit(0)

			}
			time.Sleep(1 * time.Second)
		}
	}(&activeConnections)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+strconv.Itoa(metricsPort), nil)

	// Listen for incoming connections.
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on :" + strconv.Itoa(port))
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, instance, promActiveConnections, &activeConnections)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, instance string, promActiveConnections prometheus.Gauge, activeConnections *ActiveConnections) {
	promActiveConnections.Inc()
	activeConnections.ActiveConnections++
	for {
		rawData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading:", err.Error())
		}
		data := strings.TrimSpace(string(rawData))
		fmt.Println(data)
		msg := Msg{instance, activeConnections.ActiveConnections}
		msgJsonBin, err := json.Marshal(&msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(conn, string(msgJsonBin)+"\n")
		// conn.Write(msgJsonBin + []byte("\n"))
	}
	conn.Close()
	promActiveConnections.Dec()
	activeConnections.ActiveConnections--
}
