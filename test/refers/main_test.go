package refers_test

import (
	"log"
	refersrest "refers_rest/internal/refers_rest"
	"refers_rest/pkg/settings"
	"runtime"
	"testing"
)

var (
	err error
	// настройка подключения
	config *settings.Config
)

func TestMain(m *testing.M) {
	config, err = settings.InitEnv()
	if err != nil {
		log.Fatal(err)
	}
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	m.Run()
}


func TestRunRest(t *testing.T) {
	g, err := refersrest.New(config)
	if err != nil {
		t.Fatal(err)
	}

	err = g.Start()
	if err != nil {
		t.Fatal(err)
	}

}