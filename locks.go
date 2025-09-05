package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

var ErrNotOwner = errors.New("not the owner of the lock")
var ErrNotFound = errors.New("could not find the lock")
var ErrLockExists = errors.New("lock already exists")

type Lock struct {
	Hostname  string    `json:"hostname"`
	Filename  string    `json:"filename"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LockRepo struct {
	locks []*Lock
	mu    sync.Mutex
}

func NewLock(hostname string, filename string, state string) *Lock {
	now := time.Now()

	return &Lock{
		Hostname:  hostname,
		Filename:  filename,
		State:     state,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (lr *LockRepo) find(filename string) *Lock {
	log.Println("Finding lock for filename:", filename)

	for _, lock := range lr.locks {
		if lock.Filename == filename {
			log.Println("Found lock for filename:", filename, "; host:", lock.Hostname)
			return lock
		}
	}

	log.Println("Could not find lock for filename:", filename)
	return nil
}

func (lr *LockRepo) add(lock *Lock) error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	existingLock := lr.find(lock.Filename)
	if existingLock != nil {
		return ErrLockExists
	}

	lr.locks = append(lr.locks, lock)
	log.Println("Created lock for filename:", lock.Filename, "; host:", lock.Hostname)
	return nil
}

func (lr *LockRepo) remove(hostname string, lock *Lock) error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	for i, l := range lr.locks {
		if l == lock {
			if l.Hostname != hostname {
				log.Println("Host attempted to remove lock for filename:", l.Filename, "; host:", hostname, "; actual host:", lock.Hostname)
				return ErrNotOwner
			}

			lr.locks = append(lr.locks[:i], lr.locks[i+1:]...)
			log.Println("Removed lock for filename:", lock.Filename, "; host:", lock.Hostname)
			return nil
		}
	}

	return ErrNotFound
}

func (lr *LockRepo) update(hostname string, state string, lock *Lock) error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lock.Hostname != hostname {
		log.Println("Host attempted to update lock for filename:", lock.Filename, "; host:", hostname, "; actual host:", lock.Hostname)
		return ErrNotOwner
	}

	lock.State = state
	lock.UpdatedAt = time.Now()
	log.Println("Updated lock for filename:", lock.Filename, "; state:", state, "; host:", lock.Hostname)

	return nil
}
