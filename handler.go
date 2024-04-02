package main

import (
	"encoding/json"
	"errors"
	"jobQ/job"
	"jobQ/queue"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func NewJobsHandler() JobsHandler {
	jh := JobsHandler{
		mux:   http.NewServeMux(),
		queue: queue.New(),
	}

	jh.mux.HandleFunc("POST /jobs/enqueue", jh.enqueue)
	jh.mux.HandleFunc("POST /jobs/dequeue", jh.dequeue)
	jh.mux.HandleFunc("PUT /jobs/{jobId}/conclude", jh.conclude)
	jh.mux.HandleFunc("GET /jobs/{jobId}", jh.info)

	return jh
}

type JobsHandler struct {
	mux   *http.ServeMux
	queue *queue.Queue
}

func (jh JobsHandler) enqueue(w http.ResponseWriter, r *http.Request) {
	var j job.Job

	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := j.Type.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j.Id = (rand.Uint64() << 11) >> 11 // due to json number limitations.
	j.Status = job.Queued
	err = jh.queue.Enqueue(j)
	if err == queue.ErrIdCollision {
		w.Header().Set("Retry-After", "0")
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		log.Printf("error while writing json: %v", err)
	}
}

func (jh JobsHandler) dequeue(w http.ResponseWriter, r *http.Request) {
	j, err := jh.queue.Dequeue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		log.Printf("error while writing json: %v", err)
	}
}

func (jh JobsHandler) conclude(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("jobId"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = jh.queue.Conclude(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (jh JobsHandler) info(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("jobId"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := jh.queue.Job(id)
	if errors.Is(err, queue.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		log.Printf("error while writing json: %v", err)
	}
}
