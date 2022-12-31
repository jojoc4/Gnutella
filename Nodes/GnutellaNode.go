/*
	Project : Gnutella Node
	Author : Jonatan Baumgartner
	Date : December 2022
	Based on a previous project by : Guillaume Riondet/Nabil Abdennadher
	All nodes are executed on the same machine. Each of them is represented by an IP (127.0.1.X)
	For MAC OS: to get loopback adresses other than 127.0.0.1, use this command:
	sudo ifconfig lo0 alias 127.0.1.X (this program is using AP addresses: 127.0.1.X, X ranging from 1 to 8)

	sudo ifconfig lo0 alias 127.0.1.1
	sudo ifconfig lo0 alias 127.0.1.2
	sudo ifconfig lo0 alias 127.0.1.3
	sudo ifconfig lo0 alias 127.0.1.4
	sudo ifconfig lo0 alias 127.0.1.5
	sudo ifconfig lo0 alias 127.0.1.6
	sudo ifconfig lo0 alias 127.0.1.7
	sudo ifconfig lo0 alias 127.0.1.8

	go run GnutellaNode.go

*/

package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var PORT string = ":30000"

type yamlConfig struct {
	ID         int    `yaml:"id"`
	Address    string `yaml:"address"`
	Neighbours []struct {
		ID         int    `yaml:"id"`
		Address    string `yaml:"address"`
		EdgeWeight int    `yaml:"edge_weight"`
	} `yaml:"neighbours"`
}

func initAndParseFileNeighbours(filename string) yamlConfig {
	fullpath, _ := filepath.Abs("./" + filename)
	yamlFile, err := ioutil.ReadFile(fullpath)

	if err != nil {
		panic(err)
	}

	var data yamlConfig

	err = yaml.Unmarshal([]byte(yamlFile), &data)
	if err != nil {
		panic(err)
	}

	return data
}

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

func sendToAllNeighbours(node yamlConfig, msg string) {
	for _, neigh := range node.Neighbours {
		go send(msg, neigh.Address)
	}
}
func sendToAllNeighboursExceptOne(node yamlConfig, ip string, msg string) {
	for _, neigh := range node.Neighbours {
		if ip != neigh.Address {
			go send(msg, neigh.Address)
		}
	}
}
func sendToOne(ip string, msg string) {
	go send(msg, ip)
}

func server(neighboursFilePath string) {

	var node yamlConfig = initAndParseFileNeighbours(neighboursFilePath)
	filename := "Log-" + node.Address
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	Log(file, "Parsing done ...\n")
	Log(file, "Server starting ....\n")

	ln, err := net.Listen("tcp", node.Address+PORT)
	if err != nil {
		log.Fatal(err)
		return
	}
	filesInWarehouse, err := ioutil.ReadDir("Warehouses/" + strconv.Itoa(node.ID))
	if err != nil {
		log.Fatal(err)
		return
	}
	processedRequests := make(map[string]string)
	clientRequests := make(map[string]string)

	//receive messages
	for true {
		conn, _ := ln.Accept()
		message, _ := bufio.NewReader(conn).ReadString('\n')
		conn.Close()
		Log(file, "Message received : "+message+"\n")
		decodedMessage := strings.Split(message, ";")
		messageType := decodedMessage[0]
		switch messageType {
		case "Q":
			//Query
			if _, ok := processedRequests[decodedMessage[1]]; !ok {
				processedRequests[decodedMessage[1]] = decodedMessage[4]
				ttl, _ := strconv.Atoi(decodedMessage[3])
				if ttl > 1 {
					ttl--
					for _, f := range filesInWarehouse {
						if strings.Contains(strings.ToLower(f.Name()), strings.ToLower(decodedMessage[2])) {
							sendToOne(decodedMessage[4], "R;"+decodedMessage[1]+";"+f.Name()+";"+node.Address)
						}
					}
					sendToAllNeighboursExceptOne(node, decodedMessage[4], "Q;"+decodedMessage[1]+";"+decodedMessage[2]+";"+strconv.Itoa(ttl)+";"+node.Address)
				}
			}
		case "R":
			//Response
			if _, ok := processedRequests[decodedMessage[1]]; !ok {
				if _, ok := clientRequests[decodedMessage[1]]; ok {
					sendToOne(clientRequests[decodedMessage[1]], message)
				}
			} else {
				sendToOne(processedRequests[decodedMessage[1]], message)
			}
		case "C":
			//Client request
			clientRequests[decodedMessage[1]] = decodedMessage[3]
			for _, f := range filesInWarehouse {
				if strings.Contains(strings.ToLower(f.Name()), strings.ToLower(decodedMessage[2])) {
					sendToOne(decodedMessage[3], "R;"+decodedMessage[1]+";"+f.Name()+";"+node.Address)
				}
			}
			sendToAllNeighbours(node, "Q;"+decodedMessage[1]+";"+decodedMessage[2]+";5;"+node.Address)
		}

	}

}

func main() {

	go server("node-2.yaml")
	go server("node-3.yaml")
	go server("node-4.yaml")
	go server("node-5.yaml")
	go server("node-6.yaml")
	go server("node-7.yaml")
	go server("node-8.yaml")
	server("node-1.yaml")
	time.Sleep(2 * time.Second)
}
