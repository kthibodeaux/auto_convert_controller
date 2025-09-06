package main

import (
	"encoding/json"
	"log/slog"
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
			slog.Error(err.Error())
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
			slog.Error(err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		lock := NewLock(req.Hostname, req.Filename, req.State)
		err := lockRepo.add(lock)
		if err != nil {
			switch err {
			case ErrLockExists:
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				slog.Error(err.Error())
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
			slog.Error(err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		lock := lockRepo.find(req.Filename)
		err := lockRepo.remove(req.Hostname, lock)
		if err != nil {
			switch err {
			case ErrNotOwner:
				http.Error(w, err.Error(), http.StatusForbidden)
			case ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				slog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteAllLocksHandler(lockRepo *LockRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var locksToRemove []*Lock
		for _, lock := range lockRepo.locks {
			if lock.Hostname == req.Hostname {
				locksToRemove = append(locksToRemove, lock)
			}
		}

		for _, lock := range locksToRemove {
			err := lockRepo.remove(req.Hostname, lock)
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func updateLockHandler(lockRepo *LockRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error(err.Error())
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
			switch err {
			case ErrNotOwner:
				http.Error(w, err.Error(), http.StatusForbidden)
			default:
				slog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
