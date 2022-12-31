/*
	Project : Broadcast by waves demo for SDI course
	Author : Guillaume Riondet/Nabil Abdennadher
	Date : December 2021
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

	go run broadcastbywaves.go

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
	//long++

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

func sendToAllNeighbours(node yamlConfig, typem string) {
	for _, neigh := range node.Neighbours {
		go send(node.Address+";"+typem, neigh.Address)
	}
}
func sendToAllNeighboursExceptOne(node yamlConfig, id int, typem string) {
	for _, neigh := range node.Neighbours {
		if id!=neigh.ID{
			go send(node.Address+";"+typem, neigh.Address)
		}
	}
}
func sendToOneNeighbours(node yamlConfig, id int, typem string) {
	for _, neigh := range node.Neighbours {
		if id==neigh.ID{
			go send(node.Address+";"+typem, neigh.Address)
		}
	}
}

func server(neighboursFilePath string, isStartingPoint bool) {

	var node yamlConfig = initAndParseFileNeighbours(neighboursFilePath)
	var nbNei int = len(node.Neighbours)
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
	var parent []int
	var children []int
	var Nchildren[]int

	Log(file, "Starting algorithm ...\n")
	if isStartingPoint {
		Log(file, "Starting point\n")
		parent = append(parent, 0)
		go sendToAllNeighbours(node, "M")
		nbNei++
	}

	for len(children) + len(Nchildren) + len(parent) < nbNei {
		conn, _ := ln.Accept()
		message, _ := bufio.NewReader(conn).ReadString('\n')
		conn.Close()
		Log(file, "Message received : "+message+"\n")
		address := strings.Split(message, ";")[0]
		messageType := strings.Split(message, ";")[1]
		for _, neigh := range node.Neighbours {
			if neigh.Address == address{
				switch messageType{
				case "M":
					if len(parent) == 0{
						parent = append(parent, neigh.ID)
						Log(file, "Parent: "+strconv.Itoa(neigh.ID)+"\n")
						sendToOneNeighbours(node, neigh.ID, "P")
						sendToAllNeighboursExceptOne(node, neigh.ID, "M")
					}else {
						sendToOneNeighbours(node, neigh.ID, "R")
					}
				case "P":
					children = append(children, neigh.ID)
					Log(file, "New children: "+strconv.Itoa(neigh.ID)+"\n")
				case "R":
					Nchildren = append(Nchildren, neigh.ID)
				}
			}
		}
			

	}
	
}

func main() {
	
	go server("node-2.yaml", false)
	go server("node-3.yaml", false)
	go server("node-4.yaml", false)
	go server("node-5.yaml", false)
	go server("node-6.yaml", false)
	go server("node-7.yaml", false)
	go server("node-8.yaml", false)
	time.Sleep(2 * time.Second)
	server("node-1.yaml", true)
	time.Sleep(2 * time.Second)
}
