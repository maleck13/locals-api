package routes
import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/maleck13/locals-api/data"
	"strconv"
	"encoding/json"
	"errors"
)

func Authenticate(wr http.ResponseWriter, req *http.Request)(error){
	auth := req.URL.Query().Get("auth")
	if "ssV09fGpX" != auth{
		return errors.New("auth failed");
	}
	return nil;
}

func CreateProfile(wr http.ResponseWriter, req *http.Request){
	var(
		decoder * json.Decoder
		profile * data.Profile
		err error
		encoder *json.Encoder
	)

	decoder = json.NewDecoder(req.Body)
	profile = data.NewProfile();
	encoder = json.NewEncoder(wr);

	err = decoder.Decode(profile)

	if nil != err{
		wr.WriteHeader(400)
		encoder.Encode(err)
		return;
	}
	if profile.Exists(profile.Email){
		wr.WriteHeader(409)
		encoder.Encode(errors.New("user already exists"))
		return;
	}
	_,err = profile.Save()

	if nil != err{
		wr.WriteHeader(500)
		encoder.Encode(err)
		return;
	}
	encoder.Encode(profile)
	//get handle to email channel and send email to user

}

func UpdateProfile(wr http.ResponseWriter, req *http.Request){

}

func DeleteProfile(wr http.ResponseWriter, req *http.Request){

}

func ListProfiles(wr http.ResponseWriter, req *http.Request){
	err:=Authenticate(wr,req);
	if nil != err{
		wr.WriteHeader(401)
		return;
	}
	profile := data.NewProfile()
	params:= mux.Vars(req);
	profiles,err := profile.FindByCounty(params["county"])
	if nil !=err{
		log.Print("error getting by id " + err.Error());
		wr.WriteHeader(500);
		return;
	}
	log.Printf("encoding profiles %i", len(profiles))
	enc :=json.NewEncoder(wr);
	enc.Encode(profiles);
}

func GetProfile(wr http.ResponseWriter, req *http.Request){
	err:=Authenticate(wr,req);
	if nil != err{
		wr.WriteHeader(401)
		return;
	}
	profile := data.NewProfile()
	params:= mux.Vars(req);
	id,err:= strconv.ParseInt(params["id"],10,64);
	p,err:=profile.FindById(id);
	if ae, ok := err.(*data.NoResult); ok { // if it is a NoResult
		if(404 == ae.Code()){
			wr.WriteHeader(404);
			return;
		}

	}
	if nil !=err{
		log.Print("error getting by id " + err.Error());
		wr.WriteHeader(500);
		return;
	}
	enc :=json.NewEncoder(wr);
	enc.Encode(p);
}
