/*
	Project : Gnutella Client
	Author : Jonatan Baumgartner
	Date : December 2022
	Based on a previous project by : Guillaume Riondet/Nabil Abdennadher
	The client is represented by its IP 127.0.2.1
	For MAC OS: to get loopback adresses other than 127.0.0.1, use this command:
	sudo ifconfig lo0 alias 127.0.2.1

	go run GnutellaClient.go

*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var PORT string = ":30000"
var results = make(map[string]string)

func Log(file *os.File, message string) {
	_, err := file.WriteString(message)
	if err != nil {
		panic(err)
	}

}

func send(nodeAddress string, neighAddress string) {

	outConn, err := net.Dial("tcp", neighAddress+PORT)
	if err != nil {
		log.Fatal(err)
		return
	}
	outConn.Write([]byte(nodeAddress))
	outConn.Close()
}

func sendToOne(ip string, msg string) {
	go send(msg, ip)
}

func receiver() {

	filename := "Log-client"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	Log(file, "Starting Client\n")
	ln, err := net.Listen("tcp", "127.0.2.1"+PORT)
	if err != nil {
		log.Fatal(err)
		return
	}

	//receive messages
	for true {
		conn, _ := ln.Accept()
		message, _ := bufio.NewReader(conn).ReadString('\n')
		conn.Close()
		Log(file, "Message received : "+message+"\n")
		decodedMessage := strings.Split(message, ";")
		messageType := decodedMessage[0]
		if messageType == "R" {
			if _, ok := results[decodedMessage[2]]; ok {
				results[decodedMessage[2]] = results[decodedMessage[2]] + ", " + decodedMessage[3]
			} else {
				results[decodedMessage[2]] = decodedMessage[3]
			}
		}
	}

}

func main() {

	sleeptime := 5 * time.Second
	accessNode := "127.0.1.1"
	searchedTerm := "ubuntu"

	args := os.Args[1:]
	switch len(args) {
	case 0:
		fmt.Println("Usage: go run GnutellaClient.go searchedTerm <accessNodeIP> <timeout(s)>")
		fmt.Println("Default values: accessNode = 127.0.1.1, timeout = 5 seconds")
		fmt.Println("Performing an example query with searchedTerm = ubuntu")
	case 1:
		searchedTerm = args[0]
	case 2:
		searchedTerm = args[0]
		accessNode = args[1]
	case 3:
		searchedTerm = args[0]
		accessNode = args[1]
		s, _ := strconv.Atoi(args[2])
		sleeptime = time.Duration(s) * time.Second
	}

	go receiver()
	time.Sleep(100 * time.Millisecond)
	sendToOne(accessNode, "C;"+uuid.New().String()+";"+searchedTerm+";"+"127.0.2.1")

	time.Sleep(sleeptime)

	fmt.Println("Results:")
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, " is available on: ", results[k])
	}
}
