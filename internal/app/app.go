package app

import (
	"log"
	"markat/internal/core"
	"markat/internal/imapclient"
	"markat/internal/workers"
	"markat/utils"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	Core   *core.Core
	Router *mux.Router
}

func Run() int {
	utils.PrintBanner("./ascii.txt")
	app := &App{}
	utils.LoadEnv(".env")

	addr := os.Getenv("APP_IMAP")
	username := os.Getenv("APP_ACCOUNT")
	password := os.Getenv("APP_PASSWORD")

	var err error
	ic, err := imapclient.StartTLS(addr)

	app.Core = core.InitCore(ic)

	if err != nil {
		log.Println(err)
		return -1
	}
	log.Printf("IMAP connection successful : %s\n", addr)

	err = ic.Login(username, password)
	if err != nil {
		log.Println(err)
		return -1
	}
	log.Printf("Login successful : %s\n", username)
	done := make(chan bool)
	workers.NoopTicker(done, ic)

	app.Router = InitRouter()

	err = serveHttp(3000, app)
	if err != nil {
		log.Println(err)
	}
	done <- true

	errchan := make(chan error, 1)
	go func() {
		errchan <- ic.Logout()
	}()

	select {
	case err := <-errchan:
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Logout")
		}
		break
	case <-time.After(3 * time.Second):
		log.Println("Force shutdown. Logout is taking too much time.")
		break
	}

	return 0
}
