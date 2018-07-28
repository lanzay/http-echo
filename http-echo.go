package main

import (
	"net/http"
	"log"
	"fmt"
	"net/http/httputil"
	"time"
	"net"
	"flag"
	"strings"
)

var ctx map[string]struct {
	count int
	text  string
	time time.Time
}

func main() {
	
	//port := flag.String("port","80","Port for echo listening...")
	port := flag.String("port","80","Port for echo listening...")
	flag.PrintDefaults()
	flag.Parse()
	
	ctx = make(map[string]struct {
		count int
		text  string
		time time.Time
	})
	
	http.HandleFunc("/", handleEcho)
	http.HandleFunc("/app/", handleApp)
	http.HandleFunc("/log/", handleLog)
	http.HandleFunc("/cls/", handleCls)
	go func() {
		addr := ":" + *port
		log.Println("Start listen on ", addr)
		log.Println(http.ListenAndServe(addr, nil))
	}()
	
	ch := make(chan struct{})
	<-ch
}
func handleEcho(w http.ResponseWriter, r *http.Request) {
	
	c := ctx["root"]
	c.count++
	c.time = time.Now()
	dumpR, _ := httputil.DumpRequest(r, true)
	hostR, _, _ := net.SplitHostPort(r.RemoteAddr)
	
	ss := fmt.Sprintf("----- Request #%d ", c.count)
	ss += fmt.Sprint(time.Now().String() + "\n")
	ss += "Remote host: " + hostR + "\n"
	log.Println(ss + string(dumpR))
	ss += string(dumpR) + "\n"
	c.text = ss + c.text
	//fmt.Fprint(w, c.text)
	ctx["root"] = c
	
	if len(ctx) == 1 {
		c := ctx["test-for-yor-01"]
		c.count = 0
		c.text = ""
		ctx["test-for-yor-01"] = c
	}
	
	fmt.Fprint(w,fmt.Sprintf("<p>ECHO FOR HTTP GET/POST/.... Try to http:\\\\%s\\app\\[Your Any App Name or ID]</p>",r.Host))
	fmt.Fprint(w,"<table style=\"border: 1px solid grey; border-spacing: 0.6em;\">")
	fmt.Fprint(w,"<caption>Table of Endpoint echo</caption>")
	fmt.Fprint(w,"<tr>" +
		"<th>EndPoint fo Loging Request</th>" +
		"<th>Log</th>" +
		"<th>Count REQ</th>" +
		"<th>Clear log</th>" +
		"</tr>")
	for k, v := range ctx {
		fmt.Fprint(w,"<tr>")
		appUrl := fmt.Sprintf("http:\\\\%s\\app\\%s", r.Host, k)
		logUrl := fmt.Sprintf("http:\\\\%s\\log\\%s", r.Host, k)
		clsUrl := fmt.Sprintf("http:\\\\%s\\cls\\%s", r.Host, k)
		fmt.Fprint(w, fmt.Sprintf("<td><a href=\"%s\" target=\"_blank\">%s</a></td>", appUrl, appUrl))
		fmt.Fprint(w, fmt.Sprintf("<td><a href=\"%s\" target=\"_blank\">%s</a></td>", logUrl, logUrl))
		fmt.Fprint(w, fmt.Sprintf("<td><a href=\"%s\" target=\"_blank\">%d</a></td>", logUrl, v.count))
		fmt.Fprint(w, fmt.Sprintf("<td><a href=\"%s\" target=\"_blank\">%s</a></td>", clsUrl, clsUrl))
		fmt.Fprint(w,"</tr>")
	}
	fmt.Fprint(w,"</table>")

}
func handleApp(w http.ResponseWriter, r *http.Request) {
	
	//Слушаем
	uri := strings.Split(r.RequestURI,"/")
	app := uri[2]
	
	c := ctx[app]
	c.count++
	c.time = time.Now()
	dumpR, _ := httputil.DumpRequest(r, true)
	hostR, _, _ := net.SplitHostPort(r.RemoteAddr)
	
	ss := fmt.Sprintf("----- Request #%d ", c.count)
	ss += fmt.Sprint(time.Now().String() + "\n")
	ss += "Remote host: " + hostR + "\n"
	log.Println(ss + string(dumpR))
	ss += string(dumpR) + "\n"
	c.text = ss + c.text
	fmt.Fprint(w, c.text)
	ctx[app] = c
	
}
func handleLog(w http.ResponseWriter, r *http.Request) {
	
	//Показываем
	uri := strings.Split(r.RequestURI,"/")
	app := uri[2]
	c := ctx[app]
	
	fmt.Fprint(w, c.text)
	
	ctx[app] = c
}
func handleCls(w http.ResponseWriter, r *http.Request) {
	
	//Показываем
	uri := strings.Split(r.RequestURI,"/")
	app := uri[2]
	c := ctx[app]
	
	c.text = ""
	c.count = 0
	c.time = time.Time{}
	
	ctx[app] = c
}
