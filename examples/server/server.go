package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"vallon.me/shortening"
)

func main() {
	goji.Get("/", home)

	goji.Get("/:id", redirect)
	goji.Post("/create", create)

	goji.Serve()
}

var (
	count uint64 = 0
	urls         = struct {
		set map[uint64]string
		sync.RWMutex
	}{make(map[uint64]string), sync.RWMutex{}}
)

func create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url := r.PostForm.Get("url")
	if len(url) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id := atomic.AddUint64(&count, 1)

	urls.Lock()
	defer urls.Unlock()
	urls.set[id] = url

	fmt.Fprintf(w, "%s", shortening.Encode(id))
}

func redirect(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := shortening.Decode([]byte(c.URLParams["id"]))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	urls.RLock()
	defer urls.RUnlock()

	if url, ok := urls.set[id]; ok {
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, page)
}

var page = `<!DOCTYPE html>
<html>
<head>
	<title>Shortening - Fat free URLs</title>
	<script src="//code.jquery.com/jquery-2.1.3.min.js"></script>
	<script type="text/javascript">
		$(function() {
			$('#shortBtn').on('click', function() {
				$.post("/create", {url: $("#url").val()})
					.done(function(data) {
						$('<a>', {
							text: data,
							href: "/" + data,
						}).appendTo('#shortText');
						$('#shortText').append("<br />");
					})
					.fail(function(data) {
						$("#shortText").append("ERROR");
						$('#shortText').append("<br />");
					});
			});
		});
	</script>
</head>
<body>
	<center>
		<input type="text" id="url" placeholder="http://google.com" />
		<button id="shortBtn" type="button" />Create Short URL</button>
		<br />
		<div id="shortText"></div>
	</center>
</body>
</html>`
