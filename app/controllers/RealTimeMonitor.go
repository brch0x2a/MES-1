package controllers

import (
	//"database/sql"
	//"bytes"
	"log"
	"net/http"

	//"math"
	"fmt"
	//"os"
	//"text/template"
	//"image/jpeg"
	"encoding/json"
	"time"

	//"reflect"
	//"unsafe"
	//"bytes"
	"math/rand"
	"strconv"

	//"strings"
	"io/ioutil"
	"strings"

	//"./models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	// "github.com/jung-kurt/gofpdf"
)

func Monitor00(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	tmpl.ExecuteTemplate(w, "Monitor00", nil)
}

func MonitorSetup(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	tmpl.ExecuteTemplate(w, "MonitorSetup", nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	/*
		err = ws.WriteMessage(1, []byte("Hi Client"))

		if err != nil{
			log.Println(err)
		}

	*/
	clientResult, _ := readerFilter(ws)
	clientData := strings.Split(clientResult, "_")

	turn := clientData[0]
	date := clientData[1]
	line := clientData[2]

	go WriterBy(ws, turn, date, line)
}

func WriterBy(conn *websocket.Conn, turn string, date string, line string) {

	for {
		//fmt.Printf("turn:%s\tdate:%s\tline:%s\n", turn, date, line)
		fmt.Printf("flag")
		ticker := time.NewTicker(time.Second)
		//var oee_data int64

		for i := 0; i < 115; i += 10 {
			time.Sleep(time.Second / 4)
			jsonString, err := json.Marshal(i)
			if err != nil {
				fmt.Println(err)
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonString)); err != nil {
				fmt.Println(err)
				return
			}
		}

		for t := range ticker.C {

			msg, err_msg := readerFilter(conn)

			if err_msg != nil {
				conn.Close()
				break
			}

			if msg == "on" {

				fmt.Printf("Updating At: %+v\n", t)

				//data = rand.Intn(100)

				response, err := http.Get("http://localhost:3000/getOEE?line=" + line + "&turn=" + turn + "&date=" + date)

				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
					break
				}

				data, _ := ioutil.ReadAll(response.Body)

				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					fmt.Println(err)
					break
				}

			} else {
				fmt.Printf("Client socket close!")
				conn.Close()
				break
			}
		}
		defer conn.Close()
	}
}

/*----------------------------------*/
func WsEndpointEvent(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	/*
		log.Println("Client Connected")

		err = ws.WriteMessage(1, []byte("..."))

		if err != nil {
			log.Println(err)
		}
	*/
	clientResult, _ := readerFilter(ws)
	clientData := strings.Split(clientResult, "_")

	turn := clientData[0]
	date := clientData[1]
	line := clientData[2]

	go WriterByEvent(ws, turn, date, line)
}

func WriterByEvent(conn *websocket.Conn, turn string, date string, line string) {

	for {
		fmt.Printf("\n\n----------Event------------------\n")

		ticker := time.NewTicker(time.Second)

		if err := conn.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
			fmt.Println(err)
			return
		}

		for t := range ticker.C {

			msg, err_msg := readerFilter(conn)

			if err_msg != nil {
				conn.Close()
				break
			}

			if msg == "on" {

				fmt.Printf("Updating At: %+v\n", t)

				//data = rand.Intn(100)

				response, err := http.Get("http://localhost:3000/getEventFilterV00?line=" + line + "&turn=" + turn + "&date=" + date)

				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
					break
				}

				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println("data")
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					fmt.Println(err)
					break
				}

			} else {
				fmt.Printf("Client socket close!")
				conn.Close()
				break
			}
		}
		defer conn.Close()
	}
}

/*----------------------------------*/

func Writer(conn *websocket.Conn) {

	for {

		ticker := time.NewTicker(time.Second * 3 / 2)
		var oee_data int64

		for i := 0; i < 115; i += 5 {
			time.Sleep(time.Second / 4)
			jsonString, err := json.Marshal(i)
			if err != nil {
				fmt.Println(err)
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonString)); err != nil {
				fmt.Println(err)
				return
			}
		}

		for t := range ticker.C {

			fmt.Printf("Updating At: %+v\n", t)

			//data = rand.Intn(100)

			response, err := http.Get("http://localhost:3000/getOEE?line=2&turn=1&date=20190830")

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
				break
			}

			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println("oee->>>--->" + string(data))

			oee_data, _ = strconv.ParseInt(string(data[0])+string(data[1])+string(data[2]), 10, 64)
			fmt.Printf("data: %d\n", oee_data)
			oee_data = oee_data + int64(rand.Intn(5))

			fmt.Printf("data: %d\n", oee_data)

			jsonString, err := json.Marshal(oee_data)
			if err != nil {
				fmt.Println(err)
				break
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonString)); err != nil {
				fmt.Println(err)
				break
			}
		}
	}

}

func readerFilter(conn *websocket.Conn) (string, error) {

	_, p, err := conn.ReadMessage()

	fmt.Println("client says: " + string(p))
	return string(p), err
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("client says: " + string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
