package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Hostname string `json:"hostname"`
	Filename string `json:"filename"`
	State    string `json:"state"`
}

func locksHandler(lockRepo *LockRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonBytes, err := json.Marshal(lockRepo.locks)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonBytes)
	}
}

func createLockHandler(lockRepo *LockRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		lock := NewLock(req.Hostname, req.Filename, req.State)
		err := lockRepo.add(lock)
		if err != nil {
			log.Println(err.Error())
			switch err {
			case ErrLockExists:
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func deleteLockHandler(lockRepo *LockRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		lock := lockRepo.find(req.Filename)
		err := lockRepo.remove(req.Hostname, lock)
		if err != nil {
			log.Println(err.Error())
			switch err {
			case ErrNotOwner:
				http.Error(w, err.Error(), http.StatusForbidden)
			case ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func updateLockHandler(lockRepo *LockRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		lock := lockRepo.find(req.Filename)
		if lock == nil {
			http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		err := lockRepo.update(req.Hostname, req.State, lock)
		if err != nil {
			log.Println(err.Error())
			switch err {
			case ErrNotOwner:
				http.Error(w, err.Error(), http.StatusForbidden)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
