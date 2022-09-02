package model

import (
	"fmt"

	shopify "github.com/KarloB/shopify_op/pkg/shopify"
	"github.com/KarloB/shopify_op/pkg/solo"
)

// ShopifyToPonuda convert
func ShopifyToPonuda(in *shopify.Order) *solo.Ponuda {
	usluge := make([]solo.Usluga, len(in.Items))

	for i := range in.Items {
		usluge[i] = solo.Usluga{
			Opis:           in.Items[i].Name,
			Cijena:         in.Items[i].Price,
			JedinicnaMjera: solo.KOM,
			Kolicina:       in.Items[i].Quantity,
		}
	}

	if in.Shipping != 0 {
		usluge = append(usluge, solo.Usluga{
			Opis:           "Trošak Dostave",
			Cijena:         in.Shipping,
			JedinicnaMjera: solo.KOM,
			Kolicina:       1,
		})
	}

	return &solo.Ponuda{
		TipUsluge:     1,
		PrikaziPDV:    0,
		KupacNaziv:    in.Customer.Name,
		KupacAdresa:   fmt.Sprintf("%s %s %s", in.Customer.Street, in.Customer.Postcode, in.Customer.City),
		Napomena:      fmt.Sprintf("sphynx.hr narudzba: %s", in.ID),
		NacinPlacanja: solo.TransakcijskiRacun,
		Usluge:        usluge,
	}
}

func ValidateOrders(in []shopify.Order) error {
	for i, o := range in {
		if o.ID == "" {
			return fmt.Errorf("Order %d id is empty", i)
		}
		if len(o.Items) == 0 {
			return fmt.Errorf("Order %s has no items", o.ID)
		}
		if o.Customer.Name == "" {
			return fmt.Errorf("Order %s customer name is empty", o.ID)
		}
		for _, item := range o.Items {
			if item.Name == "" {
				return fmt.Errorf("Order %v Item %d name is empty", o.ID, i)
			}
			if item.Quantity == 0 {
				return fmt.Errorf("Order %v Item %d quantity is empty", o.ID, i)
			}
			if item.Price <= 0 {
				return fmt.Errorf("Order %v Item %d price is empty", o.ID, i)
			}
		}
	}

	return nil
}
