package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"../apis"
	"../router"
	"../services"
	"../stores"
)

// TODO(tho) add comment
type Config struct {
	ListenSpec  string
	VirtualHost string
}

// TODO(tho) add comment
func Run(cfg *Config) error {
	log.Printf("Starting daemon\n")

	router := router.NewRouter()
	userStore, err := stores.NewDelegationStore()
	if err != nil {
		return err
	}

	apis.SetupDelegationRoutes(router, services.NewDelegationService(userStore))
	apis.Start(cfg.ListenSpec, router)

	waitForSignal()

	return nil
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Caught signal %v, exiting...\n", s)
}
