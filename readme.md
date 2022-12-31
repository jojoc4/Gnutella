# Gnutella
This project simulates a [Gnutella](https://rfc-gnutella.sourceforge.net/src/rfc-0_6-draft.html) network.

The nodes or servers each have their own wearhouse. If the client wants a file, it send a request to the server of its choice and wait for response.

## How-to
### Requirements
You need to have a properly working Go environnment to use this simulator.
### Setup
Each node or client use an IP adress on the loopback interface. On MacOS you need to set these up using these commands:
```bash
sudo ifconfig lo0 alias 127.0.2.1
sudo ifconfig lo0 alias 127.0.1.1
sudo ifconfig lo0 alias 127.0.1.2
sudo ifconfig lo0 alias 127.0.1.3
sudo ifconfig lo0 alias 127.0.1.4
sudo ifconfig lo0 alias 127.0.1.5
sudo ifconfig lo0 alias 127.0.1.6
sudo ifconfig lo0 alias 127.0.1.7
sudo ifconfig lo0 alias 127.0.1.8
```
If you modify the default nodes or client configuration, you will need to adapt these commands.
### Start servers
Use the following command to start the servers:
```
cd Nodes && go run GnutellaNode.go
```
### Client usage
Open a terminal in the `Client` folder.
The base command is:
```
go run GnutellaClient.go searchedTerm <accessNodeIP> <timeout(s)>
```
By default, the access node is `172.0.1.1` and the timeout is set to 5 secondes. Timeout is the time the client is waiting for responses.
If no arguments are provided, an example request for ubuntu is done.
## Customization
You can add or remove file in the warehouse folder for each node, the folder will automatically be scaned during startup of the nodes.

you can also customize the network by editing the node-[1-8].yaml files.
Be carefull that your network remains consistant across all the files or you may encounter some strange behaviors.