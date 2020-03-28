package main

import (
	"flag"
	"fmt"
	"github.com/go-stomp/stomp"
	"strconv"
)

const minTemperatura = 20
const minLuminosidad = 400

var serverAddr = flag.String("server", "localhost:61613", "STOMP server endpoint")
var topicLecturaTemp1 = flag.String("topicTemp1", "/topic/LecturasTemperaturas1", "Topic Lectura Temperatura Oficina 1")
var topicLecturaIlum1 = flag.String("topicIlum1", "/topic/LecturasIluminacion1", "Topic Lectura Iluminacion Oficina 1")
var topicLecturaTemp2 = flag.String("topicTemp2", "/topic/LecturasTemperaturas2", "Topic Lectura Temperatura Oficina 2")
var topicLecturaIlum2 = flag.String("topicIlum2", "/topic/LecturasIluminacion2", "Topic Lectura Iluminacion Oficina 2")
var topicActuadorTemp1 = flag.String("topicActuadorTemp1", "/topic/ActuadorTemperatura1", "Topic Actuador Temperatura Oficina 1")
var topicActuadorIlum1 = flag.String("topicActuadorIlum1", "/topic/ActuadorIluminacion1", "Topic Actuador Iluminacion Oficina 1")
var topicActuadorTemp2 = flag.String("topicActuadorTemp2", "/topic/ActuadorTemperatura2", "Topic Actuador Temperatura Oficina 2")
var topicActuadorIlum2 = flag.String("topicActuadorIlum2", "/topic/ActuadorIluminacion2", "Topic Actuador Iluminacion Oficina 2")

var stop = make(chan bool)

// these are the default options that work with RabbitMQ
var options = []func(*stomp.Conn) error{
	stomp.ConnOpt.Login("guest", "guest"),
	stomp.ConnOpt.Host("/"),
}

func main(){
	flag.Parse()

	subscribedLecturaTemp1 := make(chan bool)
	subscribedLecturaIlum1 := make(chan bool)
	
	go recibirMensajesTemperatura(subscribedLecturaTemp1)
	go recibirMensajesIluminacion(subscribedLecturaIlum1)
	
	<-subscribedLecturaTemp1
	<-subscribedLecturaIlum1
	
	<-stop
	<-stop
}

func recibirMensajesIluminacion(subscribed chan bool) {
	defer func(){
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)

	if err != nil{
		println("No se puede conectar al servidor ", err.Error())
		return
	}

	sub, err := conn.Subscribe(*topicLecturaIlum1, stomp.AckAuto)
	if err != nil {
		println("No se ha podido suscribir al topic", *topicLecturaIlum1, err.Error())
		return
	}
	close(subscribed)

	for {
		msg := <-sub.C
		actualText := string(msg.Body)
		println("Iluminación Recibida de la Oficina 1 ", actualText)
		var iluminacion,err =  strconv.Atoi(actualText)

		if err != nil {
			println("Error al convertir el mensaje a entero")
			return
		}

		if iluminacion < minLuminosidad {
			go enviarMensajeActuadorIluminacion()
		}
	}
}

func enviarMensajeActuadorIluminacion() {
	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil{
		println("No se puede conectar al servidor", err.Error())
	}
	text := fmt.Sprintf("%d", minLuminosidad)
	err = conn.Send(*topicActuadorIlum1, "text/plain", []byte(text), nil)
	if err != nil {
		println("Fallo al enviar al servidor", err)
		return
	}
}

func recibirMensajesTemperatura(subscribed chan bool) {
	defer func() {
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)

	if err != nil {
		println("Fallo al conectarse con el servidor", err.Error())
		return
	}

	sub, err := conn.Subscribe(*topicLecturaTemp1, stomp.AckAuto)
	if err != nil {
		println("No se ha podido suscribir al topico", *topicLecturaTemp1, err.Error())
		return
	}
	close(subscribed)

	for {
		msg := <-sub.C
		actualText := string(msg.Body)
		println("Temperatura Recibida de la Oficina 1", actualText)
		var temperatura,err =  strconv.Atoi(actualText)

		if err != nil {
			println("Error al convertir el mensaje a entero")
			return
		}
		if minTemperatura > temperatura {
			go enviarMensajeActuadorTemperatura()
		}
	}
}

func enviarMensajeActuadorTemperatura() {
	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil{
		println("No se puede conectar al servidor", err.Error())
	}
	text := fmt.Sprintf("%d", minTemperatura)
	err = conn.Send(*topicActuadorTemp1, "text/plain", []byte(text), nil)
	if err != nil {
		println("Fallo al enviar al servidor", err)
		return
	}
}
