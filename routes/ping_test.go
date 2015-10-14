package routes_test
import (
"net/http/httptest"
"github.com/maleck13/locals-api"
	"fmt"
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
)


var(
	server *httptest.Server
	pingUrl string
	//reader *strings.Reader
)


func init(){
	server = httptest.NewServer(main.SetUpRoutes());
	pingUrl = fmt.Sprintf("%s/health", server.URL)
	fmt.Println(pingUrl);
}

func TestHealth(t *testing.T){
//	profileJson := `{"email":"testprofile@test.com","county":"waterford"}`
//	reader = strings.NewReader(profileJson);
	request, err := http.NewRequest("GET", pingUrl,nil)

	res, err := http.DefaultClient.Do(request)

	handleFail(t,err)

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

	con,err :=ioutil.ReadAll(res.Body)
	handleFail(t,err)
	var jMap map[string]string
	json.Unmarshal(con,&jMap)

	i,ok:= jMap["health"]
	if ! ok{
		handleFail(t,errors.New("no health key returned"))
	}
	if "ok" != i{
		handleFail(t,errors.New("value for health should be ok"))
	}

}

func handleFail(t * testing.T , err error){
	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}
}