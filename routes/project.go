package routes

import (
	"encoding/json"
	"github.com/maleck13/locals-api/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/maleck13/locals-api/data"
	"log"
	"net/http"
	"strconv"
)

func CreateProject(wr http.ResponseWriter, req *http.Request) {
	var (
		decoder *json.Decoder
		encoder *json.Encoder
		proj    *data.Project
		err     error
	)

	decoder = json.NewDecoder(req.Body)
	encoder = json.NewEncoder(wr)
	proj = data.NewProject()

	if err = decoder.Decode(proj); err != nil {
		log.Println(err.Error())
		wr.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewErrorJSON(err.Error(), 400))
		return
	}

	if _, err = proj.Save(); err != nil {
		log.Println(err.Error())
		wr.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(NewErrorJSONUnexpectedError(err.Error()))
		return
	}

	wr.WriteHeader(http.StatusCreated)
	encoder.Encode(proj)
}

func ListProjects(wr http.ResponseWriter, req *http.Request)  {

	var(
		from int
		to int
		err error
		encoder *json.Encoder
		proj *data.Project
		listProj []*data.Project
	)
	proj = data.NewProject()
	encoder = json.NewEncoder(wr)
	f := req.URL.Query().Get("from")
	t := req.URL.Query().Get("to")
	if from,err = strconv.Atoi(f); err != nil{
		log.Println(err.Error())
		wr.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewErrorJSONBadRequest())
		return;
	}
	if to,err = strconv.Atoi(t); err != nil{
		wr.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewErrorJSONBadRequest())
		return;
	}
	listProj,err = proj.ListProjects(from,to)
	encoder.Encode(listProj)
}

func GetProject(wr http.ResponseWriter, req *http.Request) {
	var (
		encoder *json.Encoder
		proj    *data.Project
		err     error
	)

	encoder = json.NewEncoder(wr)
	proj = data.NewProject()
	params := mux.Vars(req)
	log.Println(params)
	log.Println(req.URL.String())
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewErrorJSONBadRequest())
		return
	}

	proj, err = proj.FindById(id)
	if _, ok := err.(*data.NoResult); ok {
		wr.WriteHeader(http.StatusNotFound)
		encoder.Encode(NewErrorJSONNotFound())
		return
	}
	if nil != err {
		wr.WriteHeader(http.StatusInternalServerError)
		log.Println("faild to get profile " + err.Error())
		return
	}

	encoder.Encode(proj)
}
