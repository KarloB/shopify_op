package main

import (
	"flag"
	"log"
	"time"

	"github.com/KarloB/shopify_op/internal/model"
	shopify "github.com/KarloB/shopify_op/pkg/shopify"
	"github.com/KarloB/shopify_op/pkg/solo"
)

func main() {
	var csvFilePath, soloEndpoint, soloToken string

	flag.StringVar(&csvFilePath, "csvFilePath", "", "Shopify orders CSV export")
	flag.StringVar(&soloEndpoint, "soloEndpoint", "", "Solo API endpoint")
	flag.StringVar(&soloToken, "soloToken", "", "Solo token")
	flag.Parse()

	orders, err := shopify.ParseCsv(csvFilePath)
	if err != nil {
		panic(err)
	}

	if err := model.ValidateOrders(orders); err != nil {
		panic(err)
	}

	s := solo.New(soloEndpoint, soloToken, true)
	for i, o := range orders {
		ponuda := model.ShopifyToPonuda(&orders[i])
		if err := s.CreatePonuda(ponuda); err != nil {
			log.Printf("Order %v create ponuda failed: %v\n", o.ID, err)
		} else {
			log.Printf("Order %v ponuda success\n", o.ID)
		}

		time.Sleep(10 * time.Second)
	}
}
