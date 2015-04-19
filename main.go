package main

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    "github.com/kylelemons/go-gypsy/yaml"
    zmq "github.com/pebbe/zmq4"
)

var (
    wsaddr     = ":8080"
    zmqaddr    = "tcp://localhost:5563"
    zmqsubject = "message"
    responder  *zmq.Socket
)

func main() {
    // parse config file
    file := "config.yaml"
    config, err := yaml.ReadFile(file)
    if err != nil {
        log.Fatal("Error load config", err)
    }

    wsaddr, _ = config.Get("ws.address")
    zmqaddr, _ = config.Get("zmq.address")
    zmqsubject, _ = config.Get("zmq.subject")

    // connect to ZMQ broker
    responder, _ = zmq.NewSocket(zmq.SUB)
    defer responder.Close()
    responder.Connect(zmqaddr)
    responder.SetSubscribe(zmqsubject)

    // start handle ws
    http.HandleFunc("/ws", handler)
    err = http.ListenAndServe(wsaddr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); ok {
        http.Error(w, "Not a websocket handshake", 400)
        return
    } else if err != nil {
        log.Println(err)
        return
    }

    log.Println("Start...")
    for {
        // gets message from ZMQ broker
        msg, _ := responder.RecvMessage(0)
        log.Println("Received ", msg)
        // websockets
        if err := conn.WriteMessage(websocket.TextMessage, []byte(msg[1])); err != nil {
            log.Println(err)
            return
        }

    }
}
