module tgtsrv

go 1.20

replace asm13sam/tg => ../tg

require (
	asm13sam/tg v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.2
	github.com/gorilla/websocket v1.5.1
	github.com/mattn/go-sqlite3 v1.14.17
)

require (
	github.com/gorilla/securecookie v1.1.2 // indirect
	golang.org/x/net v0.17.0 // indirect
)
