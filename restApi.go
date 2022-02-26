package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

type cars struct {
	ID           string `json:"id"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	Power        int    `json:"power"`
	Color        string `json:"color"`
}

type carHandler struct {
	sync.Mutex
	store map[string]cars
}

func (c *carHandler) req(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("id") != "" {
			c.getByID(w, r)
			return
		}
		c.get(w, r)
		return

	case "POST":
		c.post(w, r)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func (c *carHandler) get(w http.ResponseWriter, r *http.Request) {
	car := make([]cars, len(c.store))
	c.Lock()
	i := 0
	for _, cars := range c.store {
		car[i] = cars
		i++
	}
	c.Unlock()

	jsonBytes, err := json.Marshal(car)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func (c *carHandler) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var lcars cars
	err = json.Unmarshal(bodyBytes, &lcars)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	c.Lock()
	c.store[lcars.ID] = lcars
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The item was added"))
	defer c.Unlock()
}

func newCarHandler() *carHandler {
	return &carHandler{
		store: map[string]cars{
			"1": {
				ID:           "1",
				Manufacturer: "Tesla",
				Model:        "Model X",
				Year:         2022,
				Power:        1020,
				Color:        "Matte Black",
			},
			"2": {
				ID:           "2",
				Manufacturer: "Ford",
				Model:        "Focus RS",
				Year:         2016,
				Power:        345,
				Color:        "Blue",
			},
			"3": {
				ID:           "3",
				Manufacturer: "Subaru",
				Model:        "Impreza",
				Year:         2010,
				Power:        305,
				Color:        "Metalic Silver",
			},
		},
	}
}

func (c *carHandler) getByID(w http.ResponseWriter, r *http.Request) {
	c.Lock()
	id := r.URL.Query().Get("id")
	c.Unlock()

	if _, ok := c.store[id]; ok {

		c.Lock()
		car := c.store[id]
		jsonBytes, err := json.Marshal(car)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", "application/json")
		w.Write(jsonBytes)
		c.Unlock()

	} else {

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(id))
	}
}

func main() {
	carHandler := newCarHandler()
	http.HandleFunc("/cars", carHandler.req)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}
