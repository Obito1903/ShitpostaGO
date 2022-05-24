package api_test

import (
	"testing"

	"github.com/Obito1903/shitpostaGo/pkg/api"
)

func TestInitDB(t *testing.T) {
	api.Start("../../appTest")

}
