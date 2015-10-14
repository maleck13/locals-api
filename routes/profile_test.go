package routes_test
import (
"fmt"
"github.com/maleck13/locals-api"
	"net/http/httptest"
)

var(
	server *httptest.Server
	pingUrl string
//reader *strings.Reader
)


func init(){
	server = httptest.NewServer(main.SetUpRoutes());
	pingUrl = fmt.Sprintf("%s/profile", server.URL)
	fmt.Println(pingUrl);
}
