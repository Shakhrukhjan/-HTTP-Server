package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	host := "0.0.0.0"
	port := "8080"
	err := execute(host, port)
	if err != nil {
		os.Exit(1)
	}
	http.HandleFunc("/", Response_page)
	http.ListenAndServe(port, nil)

}
func execute(host string, port string) (err error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		cerr := listener.Close()
		if cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			return err
		}
		err = handler(conn)
		if err != nil {
			log.Print(err)
			return err
		}
	}
}
func Response_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет Golang")
}
func handler(con net.Conn) (err error) {
	defer func() {

		if temp := con.Close(); temp != nil {
			if err == nil {
				err = temp
				return
			}
			log.Print(err)
		}
	}()
	buff := make([]byte, 4096)
	n, err := con.Read(buff)
	if err == io.EOF {
		log.Printf("%s", buff[:n])
		return nil
	}
	if err != nil {
		return err
	}
	data := buff[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {
	}
	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
	}
	method, path, version := parts[0], parts[1], parts[2]
	if method != "GET" {
	}
	if version != "HTTP/1.1" {
	}
	if path == "/" {
		body, err := ioutil.ReadFile("static/index.html")
		if err != nil {
			return fmt.Errorf("can't read index.html: %w", err)
		}
		marker := "{{year}}"
		year := time.Now().Year()
		body = bytes.ReplaceAll(body, []byte(marker), []byte(strconv.Itoa(year)))
		_, err = con.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Lenth: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: text/html\r\n" +
				"Connection:close\r\n" +
				"r\n" +
				string(body),
		))
		if err != nil {
			return err
		}
	}
	return nil
}
