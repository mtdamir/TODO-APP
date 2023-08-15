package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/ToDo_App/delivery/deleveryparam"
	"net"
	"os"
)

func main() {
	fmt.Println("command", os.Args[0])

	message := "default message"
	if len(os.Args) > 1 {
		message = os.Args[1]
	}

	connection, err := net.Dial("tcp", "127.0.0.1:1986")
	if err != nil {
		log.Fatalln("can't dial the given address", err)
	}
	defer connection.Close()

	fmt.Println("local address", connection.LocalAddr())

	req := deleveryparam.Request{Command: message}

	if req.Command == "create-task" {
		req.CreateTaskRequest = deleveryparam.CreateTaskRequest{
			Title:      "test",
			DueDate:    "test",
			CategoryID: 1,
		}
	}

	serializedData, mErr := json.Marshal(&req)
	if mErr != nil {
		log.Fatalln("can' marshal reqeust", mErr)
	}

	numberOfWrittenBytes, wErr := connection.Write(serializedData)
	if wErr != nil {
		log.Fatalln("can't write data to connection", wErr)
	}

	fmt.Println("numberOfWrittenBytes", numberOfWrittenBytes)

	var data = make([]byte, 1024)
	_, rErr := connection.Read(data)
	if rErr != nil {
		log.Fatalln("can't read data from connection", rErr)
	}

	fmt.Println("server response:", string(data))
}
