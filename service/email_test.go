package service_test
import (
	"testing"
	"errors"
	"github.com/maleck13/locals-api/service"
)


type Test_Sender struct {

}

func (Test_Sender) Send(from,to,content string)error{
	var(
		err error
	)
	if "" == from || "" == to || "" == content{
		err = errors.New("missing params");
	}

	return err;
}

func Test_Email_Send (t *testing.T){
	err := service.SendMailTemplate(service.MAIL_TEMPLATE_INTEREST,Test_Sender{},"test@test.com","test@to.from")
	if nil != err{
		t.Fatal("did not expect err during send " + err.Error())
	}
}

func Test_Email_Send_Bad_Template (t *testing.T){
	err := service.SendMailTemplate("notemplate",Test_Sender{},"test@test.com","test@to.from")
	if nil == err{
		t.Fatal("expected err during send but there was none ")
	}
}