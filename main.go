// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 1 * time.Second
)

var (
	addr        = flag.String("addr", ":8080", "http service address")
	homeTempl   = template.Must(template.New("").Parse(homeHTML))
	clientTempl = template.Must(template.New("").Funcs(template.FuncMap{"printTempl": PrintTempl}).Parse(clientHTML))
	filename    string
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	buttonList = make([]Button, 0)
)

type Button struct {
	Name    string
	Color   string
	Updated time.Time
}

func readFileIfModified(lastMod time.Time) ([]byte, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return nil, lastMod, nil
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fi.ModTime(), err
	}
	return p, fi.ModTime(), nil
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, lastMod time.Time) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-fileTicker.C:
			var p []byte
			var err error

			p, lastMod, err = readFileIfModified(lastMod)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
					return
				}
			}
		case pin := <-buttonPush:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			ba, err := json.Marshal(pin)
			if err != nil {
				return
			}
			if err := ws.WriteMessage(websocket.TextMessage, ba); err != nil {
				return
			}

		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n)
	}

	go writer(ws, lastMod)
	reader(ws)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p, lastMod, err := readFileIfModified(time.Time{})
	if err != nil {
		p = []byte(err.Error())
		lastMod = time.Unix(0, 0)
	}
	var v = struct {
		Host    string
		Data    string
		LastMod string
	}{
		r.Host,
		string(p),
		strconv.FormatInt(lastMod.UnixNano(), 16),
	}
	homeTempl.Execute(w, &v)
}

func serveClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var v = struct {
		Host    string
		LastMod string
		Circles []Button
	}{
		r.Host,
		strconv.FormatInt(time.Now().UnixNano(), 16),
		buttonList,
	}
	clientTempl.Execute(w, &v)
}

func serveSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
	}

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	color := r.PostFormValue("color")
	name := r.PostFormValue("name")

	pin := Pin{Name: name, Value: "1", Color: color, Updated: time.Now()}
	buttonPush <- pin

	var v = struct {
		Host    string
		LastMod string
	}{
		r.Host,
		strconv.FormatInt(time.Now().UnixNano(), 16),
	}
	clientTempl.Execute(w, &v)
}

func serveJquery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	fmt.Fprint(w, AJAX_JS)
}

func serveCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	fmt.Fprint(w, CSS)
}

func initButtons() {
	buttonList = append(buttonList, Button{Name: "Blue", Color: "blue"})
	buttonList = append(buttonList, Button{Name: "Red", Color: "red"})
	buttonList = append(buttonList, Button{Name: "Yellow", Color: "yellow"})
	buttonList = append(buttonList, Button{Name: "Green", Color: "Green"})
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Println("filename not specified")
	} else {
		filename = flag.Args()[0]
	}
	initButtons()
	InitGpioPoll()

	http.HandleFunc("/host", serveHome)
	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/", serveClient)
	http.HandleFunc("/submit", serveSubmit)
	http.HandleFunc("/jquery-3.2.1.min.js", serveJquery)
	http.HandleFunc("/styles.css", serveCss)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
<script src="jquery-3.2.1.min.js"></script>
<link rel="stylesheet" type="text/css" href="styles.css" media="all"></link>
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
	<div id="winner" class="center" style="width: 100%; min-height: 350px; text-align: left">
	<span id="winner_name" class="left"></span>
	</div>
        <pre id="fileData" class="left">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('file updated');
					append = data.textContent;
					j = JSON.parse(data.textContent);
					$("#winner_name").textContent = j["Name"];
					$("#winner").css('background-color', j["Color"]);
                    data.textContent = evt.data + "\n" +append;
                }
            })();
        </script>
    </body>
</html>
`

func PrintTempl(text string) string {
	return text
}

const clientHTML = `<!DOCTYPE html>
<html lang="en">
<script src="jquery-3.2.1.min.js"></script>
<link rel="stylesheet" type="text/css" href="styles.css" media="all"></link>
    <head>
        <title>Responder</title>
    </head>
    <body>
		<form id="respond" action="/submit" method="POST">
		  Color:<br>
		<div class="center">

		{{range .Circles}}<div id="{{ .Name }}" class="circle not-selected" color="{{.Color}}" style="background-color: {{ .Color}};"></div>{{else}}<div><strong>no rows</strong></div>{{end}}
		
		</div>
		
			<br>
		  Name:<br>
		    <input id="name" type="text" name="name" value="">
		    <br><br>
		  <a href="#" id="submit" class="button" type="button" value="Respond!">
		    <span>Respond!</span>
		  </a>
	  	</form> 
<script type="text/javascript">
var selected = $("#Blue");
$(".button").css("background-color", selected.attr("color"));
$(".circle").on('click', function(e) {
		selected = $(this)
		$(".selected").toggleClass('selected not-selected');
		$(".button").css("background-color", selected.attr("color"));
		$(this).toggleClass('not-selected selected');
});
selected.toggleClass('not-selected selected');

$("#submit").on('click', function(e) {
  $.post( "/submit", {"color": selected.attr("color"), "name": $('#name').val()}, function( data ) {
    $( ".result" ).html( data );
  });
});
</script>
    </body>
</html>
`
