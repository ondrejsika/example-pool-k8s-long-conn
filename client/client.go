package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	host := flag.String("host", "127.0.0.1", "")
	port := flag.Int("port", 3333, "")
	instance := flag.String("instance", "client-0", "")
	count := flag.Int("count", 100, "")
	sleep := flag.Int("sleep", 1, "")

	flag.Parse()

	client(*host, *port, *instance, *count, *sleep)
}

type Msg struct {
	Instance string
	Count    int
}

func client(host string, port int, instance string, count int, sleep int) {
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		msg := Msg{instance, count}
		msgJsonBin, err := json.Marshal(&msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(conn, string(msgJsonBin)+"\n")
		fmt.Println(string(msgJsonBin))
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(message)
		time.Sleep(time.Second)
		if count <= 0 {
			conn.Close()
			break
		}
		count--
	}
}
