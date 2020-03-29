package main

import (
	"flag"
	"fmt"
	"github.com/go-stomp/stomp"
	"strconv"
)

const minTemperaturaOficina1 = 20
const minLuminosidadOficina1 = 400
const minTemperaturaOficina2 = 30
const minLuminosidadOficina2 = 300

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
	stomp.ConnOpt.Login("user", "password"),
	stomp.ConnOpt.Host("/"),
}

func main(){
	flag.Parse()

	subscribedLecturaTemp1 := make(chan bool)
	subscribedLecturaIlum1 := make(chan bool)
	subscribedLecturaTemp2 := make(chan bool)
	subscribedLecturaIlum2 := make(chan bool)
	
	go recibirMensajesTemperatura1(subscribedLecturaTemp1)
	go recibirMensajesIluminacion1(subscribedLecturaIlum1)
	go recibirMensajesTemperatura2(subscribedLecturaTemp2)
	go recibirMensajesIluminacion2(subscribedLecturaIlum2)
	
	<-subscribedLecturaTemp1
	<-subscribedLecturaIlum1
	<-subscribedLecturaTemp2
	<-subscribedLecturaIlum2
	
	<-stop
	<-stop
}





func recibirMensajesIluminacion2(subscribed chan bool) {
	defer func(){
		stop <-true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)

	if err != nil{
		println("No se puede conectar al servidor ", err.Error())
		return
	}

	sub, err := conn.Subscribe(*topicLecturaIlum2, stomp.AckAuto)
	if err != nil {
		println("No se ha podido suscribir al topic", *topicLecturaIlum2, err.Error())
		return
	}
	close(subscribed)

	for {
		msg := <-sub.C
		actualText := string(msg.Body)
		println("Iluminación Recibida de la Oficina 2", actualText)
		var iluminacion,err =  strconv.Atoi(actualText)

		if err != nil {
			println("Error al convertir el mensaje a entero")
			return
		}

		if iluminacion < minLuminosidadOficina2 {
			go enviarMensajeActuadorIluminacion(topicActuadorIlum2, minLuminosidadOficina2)
		}
	}
}

func recibirMensajesIluminacion1(subscribed chan bool) {
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

		if iluminacion < minLuminosidadOficina1 {
			go enviarMensajeActuadorIluminacion(topicActuadorIlum1, minLuminosidadOficina1)
		}
	}
}

func enviarMensajeActuadorIluminacion(topicActuador *string, minLuminosidad int) {
	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil{
		println("No se puede conectar al servidor", err.Error())
	}
	text := fmt.Sprintf("%d", minLuminosidad)
	err = conn.Send(*topicActuador, "text/plain", []byte(text), nil)
	if err != nil {
		println("Fallo al enviar al servidor", err)
		return
	}
}

func recibirMensajesTemperatura2(subscribed chan bool) {
	defer func() {
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)

	if err != nil {
		println("Fallo al conectarse con el servidor", err.Error())
		return
	}

	sub, err := conn.Subscribe(*topicLecturaTemp2, stomp.AckAuto)
	if err != nil {
		println("No se ha podido suscribir al topico", *topicLecturaTemp2, err.Error())
		return
	}
	close(subscribed)

	for {
		msg := <-sub.C
		actualText := string(msg.Body)
		println("Temperatura Recibida de la Oficina 2", actualText)
		var temperatura,err =  strconv.Atoi(actualText)

		if err != nil {
			println("Error al convertir el mensaje a entero")
			return
		}
		if minTemperaturaOficina2 > temperatura {
			go enviarMensajeActuadorTemperatura(topicActuadorTemp2, minTemperaturaOficina2)
		}
	}
}
func recibirMensajesTemperatura1(subscribed chan bool) {
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
		if minTemperaturaOficina1 > temperatura {
			go enviarMensajeActuadorTemperatura(topicActuadorTemp1, minTemperaturaOficina1)
		}
	}
}

func enviarMensajeActuadorTemperatura(topicActuador *string,minTemperatura int) {
	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil{
		println("No se puede conectar al servidor", err.Error())
	}
	text := fmt.Sprintf("%d", minTemperatura)
	err = conn.Send(*topicActuador, "text/plain", []byte(text), nil)
	if err != nil {
		println("Fallo al enviar al servidor", err)
		return
	}
}
