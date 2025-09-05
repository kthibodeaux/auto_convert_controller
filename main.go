package main

import (
	"net/http"
)

func main() {
	lockRepo := &LockRepo{}

	http.HandleFunc("GET /locks", locksHandler(lockRepo))
	http.HandleFunc("POST /locks", createLockHandler(lockRepo))
	http.HandleFunc("DELETE /lock", deleteLockHandler(lockRepo))
	http.HandleFunc("PUT /lock", updateLockHandler(lockRepo))

	http.ListenAndServe(":8080", nil)
}
