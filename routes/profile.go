package routes

import (
	"encoding/json"
	"errors"
	"github.com/maleck13/locals-api/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/maleck13/locals-api/data"
	"github.com/maleck13/locals-api/service"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

func Authenticate(wr http.ResponseWriter, req *http.Request) error {
	auth := req.URL.Query().Get("auth")
	if "ssV09fGpX" != auth {
		return errors.New("auth failed")
	}
	return nil
}

func CreateProfile(wr http.ResponseWriter, req *http.Request) {
	var (
		decoder *json.Decoder
		profile *data.Profile
		err     error
		encoder *json.Encoder
	)

	decoder = json.NewDecoder(req.Body)
	profile = data.NewProfile()
	encoder = json.NewEncoder(wr)

	if err = decoder.Decode(profile); err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
	}

	if profile.Exists(profile.Email) {
		wr.WriteHeader(http.StatusConflict)
		encoder.Encode(NewErrorJSONExists())
		return
	}
	if _, err = profile.Save(); err != nil {
		log.Printf("error occurred saving profile %s ", err.Error())
		wr.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(NewErrorJSONUnexpectedError(err.Error()))
		return
	}

	encoder.Encode(profile)
	//get handle to email channel and send email to user
	go service.SendMailTemplate(service.MAIL_TEMPLATE_INTEREST, service.NewMailSender(), service.MAIL_FROM, profile.Email)

}

func UpdateProfile(wr http.ResponseWriter, req *http.Request) {
	var (
		decoder *json.Decoder
		profile *data.Profile
		err     error
		encoder *json.Encoder
	)

	decoder = json.NewDecoder(req.Body)
	profile = data.NewProfile()
	encoder = json.NewEncoder(wr)
	if err = decoder.Decode(profile); nil != err {
		wr.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewErrorJSONBadRequest())
		return
	}

	if err = profile.Update(); err != nil {
		log.Printf("error occurred saving profile %s ", err.Error())
		wr.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(NewErrorJSONUnexpectedError(err.Error()))
		return
	}
	encoder.Encode(profile)
}

func DeleteProfile(wr http.ResponseWriter, req *http.Request) {

}

func ListProfiles(wr http.ResponseWriter, req *http.Request) {
	err := Authenticate(wr, req)
	if nil != err {
		wr.WriteHeader(401)
		return
	}
	params := mux.Vars(req)
	profile := data.NewProfile()
	log.Printf("listing by county %s ", params)
	profiles, err := profile.FindByCounty(params["county"])

	if nil != err {
		log.Print("error getting list " + err.Error())
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("encoding profiles %i", len(profiles))
	enc := json.NewEncoder(wr)
	enc.Encode(profiles)
}

func GetProfile(wr http.ResponseWriter, req *http.Request) {
	if err := Authenticate(wr, req); nil != err {
		wr.WriteHeader(http.StatusUnauthorized)
		return
	}
	enc := json.NewEncoder(wr)
	profile := data.NewProfile()
	params := mux.Vars(req)
	log.Println(params)
	log.Println(req.URL.String())
	id, err := strconv.ParseInt(params["id"], 10, 64)
	p, err := profile.FindById(id)
	if _, ok := err.(*data.NoResult); ok { // if it is a NoResult
		wr.WriteHeader(http.StatusNotFound)
		enc.Encode(NewErrorJSONNotFound())
		return
	}
	if nil != err {
		wr.WriteHeader(http.StatusInternalServerError)
		log.Println("faild to get profile " + err.Error())
		return
	}

	enc.Encode(p)
}

func UploadProfilePic(wr http.ResponseWriter, req *http.Request) {
	var (
		profileImgLoc string
		id            int64
		err           error
		p             *data.Profile
		file          multipart.File
		header        *multipart.FileHeader
		enc           *json.Encoder
	)

	enc = json.NewEncoder(wr)
	params := mux.Vars(req)
	id, err = strconv.ParseInt(params["id"], 10, 64)
	p, err = data.FindProfileById(id)
	req.ParseMultipartForm(10 << 20) //approx 10MB
	file, header, err = req.FormFile("file")

	handleUploadErr := func(err error, status int) {
		if nil != err {
			wr.WriteHeader(status)
			enc.Encode(NewErrorJSON(err.Error(), status))
		}
	}

	if err != nil {
		log.Println("error upload pic " + err.Error())
		handleUploadErr(err, http.StatusBadRequest)
		return
	}
	defer file.Close()
	uploadedFilePath, err := service.SaveUploadedFile(file, header.Filename)
	if err != nil {
		log.Println("failed to create thumbnail file  " + err.Error())
		handleUploadErr(err, http.StatusInternalServerError)
		return
	}

	uploadedFilePath, err = service.ThumbnailMultipart(file, header.Filename)
	if err != nil {
		log.Println("failed to create thumbnail file  " + err.Error())
		handleUploadErr(err, http.StatusInternalServerError)
		return
	}

	profileImgLoc, err = data.PutInBucket(uploadedFilePath, header.Filename)

	if err != nil {
		log.Println("failed up upload to s3 " + err.Error())
		handleUploadErr(err, http.StatusInternalServerError)
		return
	}

	p.UpdateProfilePic(profileImgLoc)
	enc.Encode(p)

}
