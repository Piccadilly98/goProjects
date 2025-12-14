package unit_handler

import (
	"context"
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests/init_test_server"
)

type RegistrationHandlerTest struct {
}

func TestRegistrationHanlder(t *testing.T) {
	serv, err := init_test_server.InitTestServer(false)
	if err != nil {
		t.Fatalf("error in init server: %v\n", err)
	}
	_, err = serv.Db.RegistrationBoard(context.Background(), init_test_server.GetPtrStr("esp32_1"), nil, nil, init_test_server.GetPtrStr("registred"))
	if err != nil {
		t.Fatal(err)
	}
}
