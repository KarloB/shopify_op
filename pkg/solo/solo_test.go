package solo

import (
	"os"
	"testing"
	"time"
)

func TestSolo(t *testing.T) {
	apiKey := os.Getenv("SOLO_KEY")
	if apiKey == "" {
		t.Fatalf("Solo api key not defined")
	}
	soloEndpoint := os.Getenv("SOLO_ENDPOINT")
	if soloEndpoint == "" {
		t.Fatalf("Solo endpoint not defined")
	}

	s := New(soloEndpoint, apiKey, true)

	testData := &Ponuda{
		TipUsluge:     1,
		PrikaziPDV:    0,
		KupacNaziv:    "John Matrix Terminatorović",
		KupacAdresa:   "Predatorova 58",
		NacinPlacanja: 1,
		RokPlacanja:   time.Now().AddDate(0, 0, 7),
		Napomena:      "Shopify order ID 12345",
		Usluge: []Usluga{
			{
				Opis:           "Voskić lijepi crni",
				Cijena:         35,
				Kolicina:       2,
				JedinicnaMjera: KOM,
			},
			{
				Opis:           "Svijeća",
				Cijena:         1234.56,
				Kolicina:       1,
				JedinicnaMjera: KOM,
			},
			{
				Opis:           "Šargija",
				Cijena:         67.54,
				Kolicina:       2,
				JedinicnaMjera: KOM,
			},
		},
	}

	if err := s.CreatePonuda(testData); err != nil {
		t.Fatal(err)
	}
}
