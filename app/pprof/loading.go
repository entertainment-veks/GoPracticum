package pprof

import (
	"fmt"
	"go_practicum/app/config"
	"net/http"
)

func Loading(cfg config.Config) {
	for {
		client := http.DefaultClient

		_, err := CallCreateLinkHandler(client, cfg)
		if err != nil {
			fmt.Println(err)
		}

		_, err = CallCreateLinkJSONHandler(client, cfg)
		if err != nil {
			fmt.Println(err)
		}

		_, err = CallCreateAllLinkHandler(client, cfg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
