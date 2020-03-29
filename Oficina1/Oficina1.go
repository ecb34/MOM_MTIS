package main

import (
	"flag"
	"fmt"
	"github.com/go-stomp/stomp"
	"math/rand"
	"strconv"
	"time"
)

var temperatura = rand.Intn(50)
var iluminacion = rand.Intn(1000-200) + 200

var serverAddr = flag.String("server", "localhost:61613", "STOMP server endpoint")
var topicLecturaTemp = flag.String("topicTemp", "/topic/LecturasTemperaturas1", "Topic Lectura temperatura")
var topicLecturaIlum = flag.String("topicIlum", "/topic/LecturasIluminacion1", "Topic Lectura Iluminacion")
var topicActuadorTemp = flag.String("topicActuadorTemp", "/topic/ActuadorTemperatura1", "Topic Actuador Temperatura")
var topicActuadorIlum = flag.String("topicActuadorIlum", "/topic/ActuadorIluminacion1", "Topic Actuador Iluminacion")

var stop = make(chan bool)

// these are the default options that work with RabbitMQ
var options = []func(*stomp.Conn) error{
	stomp.ConnOpt.Login("guest", "guest"),
	stomp.ConnOpt.Host("/"),
}

func main() {
	flag.Parse()

	subscribedActuadorTemp := make(chan bool)
	subscribedActuadorIlum := make(chan bool)

	go recibirMensajeActuadorTemperatura(subscribedActuadorTemp)
	go recibirMensajeActuadorIluminacion(subscribedActuadorIlum)

	// wait until we know the receiver has subscribed
	<-subscribedActuadorTemp
	<-subscribedActuadorIlum

	go enviarMensajeLecturaTemperatura()
	go enviarMensajeLecturaIluminacion()

	<-stop
	<-stop
	<-stop
	<-stop
}



func enviarMensajeLecturaIluminacion() {
	defer func() {
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil{
		println("No se puede conectar al servidor", err.Error())
	}
	for{
		text := fmt.Sprintf("%d", iluminacion)
		err = conn.Send(*topicLecturaIlum, "text/plain", []byte(text), nil)
		if err != nil {
			println("Fallo al enviar al servidor", err)
			return
		}
		iluminacion = rand.Intn(1000-200) + 200
		time.Sleep(5 * time.Second)
	}
}

func enviarMensajeLecturaTemperatura() {
	defer func() {
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil {
		println("No se puede conectar al servidor", err.Error())
		return
	}
	for {
		text := fmt.Sprintf("%d", temperatura)
		err = conn.Send(*topicLecturaTemp, "text/plain",
			[]byte(text), nil)
		if err != nil {
			println("Fallo al enviar al servidor", err)
			return
		}
		temperatura = rand.Intn(50)
		time.Sleep(5 * time.Second)
	}
}

func recibirMensajeActuadorIluminacion(subscribed chan bool) {
	defer func(){
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)

	if err != nil{
		println("No se puede conectar con el servidor", err.Error())
		return
	}

	sub, err := conn.Subscribe(*topicActuadorIlum, stomp.AckAuto)
	if err != nil {
		println("No se puede suscribir al topic", *topicLecturaIlum, err.Error())
		return
	}
	close(subscribed)

	for {
		msg := <-sub.C
		actualText := string(msg.Body)
		println("Subiendo Iluminación a", actualText)
		iluminacion, err  = strconv.Atoi(actualText)

		if err != nil{
			println("Error en el dato recibido, no se puede convertir a entero")
		}
	}
}

func recibirMensajeActuadorTemperatura(subscribed chan bool) {
	defer func() {
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)

	if err != nil {
		println("No se puede conectar al servidor", err.Error())
		return
	}

	sub, err := conn.Subscribe(*topicActuadorTemp, stomp.AckAuto)
	if err != nil {
		println("No se ha podido suscribir al topic", *topicActuadorTemp, err.Error())
		return
	}
	close(subscribed)

	for {
		msg := <-sub.C
		actualText := string(msg.Body)
		println("Subiendo Temperatura a", actualText)
		temperatura, err  = strconv.Atoi(actualText)

		if err != nil{
			println("Error en el dato recibido, no se puede convertir a entero")
		}
	}
}
