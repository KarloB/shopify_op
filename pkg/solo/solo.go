package solo

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/KarloB/amount"
	qs "github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

type Ponuda struct {
	TipUsluge     int           `url:"tip_usluge"`
	PrikaziPDV    int           `url:"prikazi_porez"`
	KupacNaziv    string        `url:"kupac_naziv"`
	KupacAdresa   string        `url:"kupac_adresa"`
	NacinPlacanja NacinPlacanja `url:"-"` // nacin_placanja
	RokPlacanja   time.Time     `url:"-"` // Dozvoljen je ISO 8601 format datuma (npr. 2014-01-01).
	Usluge        []Usluga      `url:"-"`
	Napomena      string        `url:"-"`
}

type Usluga struct {
	Opis           string         `url:"opis_usluge_x"`
	JedinicnaMjera JedinicnaMjera `url:"jed_mjera_x"`
	Cijena         float64        `url:"cijena_x"`
	Kolicina       float64        `url:"kolicina_x"`
	Popust         float64        `url:"popust_x"`
	PorezStopa     int            `url:"porez_stopa_x"` //  Podržane porezne stope su 0, 5, 13 i 25.

}

type JedinicnaMjera uint8

const (
	KOM JedinicnaMjera = 2
	KG  JedinicnaMjera = 7
)

type NacinPlacanja uint8

const (
	TransakcijskiRacun NacinPlacanja = iota + 1
	Gotovina
	Kartice
	Cek
	Ostalo
)

type Solo struct {
	endpoint string
	token    string
	debugLog bool
}

func New(endpoint, token string, debugLog bool) *Solo {
	return &Solo{
		endpoint: endpoint,
		token:    token,
		debugLog: debugLog,
	}
}

func (t *Solo) CreatePonuda(in *Ponuda) error {
	vals, _ := qs.Values(in)
	encodedFirstPart := vals.Encode()

	customQueryParams := make([]string, len(in.Usluge))
	for i, u := range in.Usluge { // cant use map because of duplicate and custom keys...
		customQueryParams[i] = fmt.Sprintf(`usluga=%d&opis_usluge_%d=%s&jed_mjera_%d=%d&cijena_%d=%s&kolicina_%d=%.0f&popust_%d=%s&porez_stopa_%d=%d`,
			i+1,
			i+1, url.QueryEscape(u.Opis),
			i+1, u.JedinicnaMjera,
			i+1, amount.Amount(u.Cijena),
			i+1, u.Kolicina,
			i+1, amount.Amount(u.Popust),
			i+1, u.PorezStopa,
		)
	}

	if in.RokPlacanja.IsZero() {
		in.RokPlacanja = time.Now().AddDate(0, 0, 7)
	}

	usluge := strings.Join(customQueryParams, "&")
	encodedSecondPart := fmt.Sprintf(`nacin_placanja=%d&napomene=%s&rok_placanja=%s`, in.NacinPlacanja, url.QueryEscape(in.Napomena), in.RokPlacanja.Format("2006-01-02"))
	params := fmt.Sprintf("?token=%s&%s&%s&%s", t.token, encodedFirstPart, usluge, encodedSecondPart)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ponuda%s", t.endpoint, params), nil)
	if err != nil {
		return errors.Wrapf(err, "CreatePonuda.NewRequest")
	}

	if t.debugLog {
		reqB, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(reqB))
	}

	cl := http.DefaultClient
	res, err := cl.Do(req)
	if err != nil {
		return errors.Wrapf(err, "CreatePonuda.http.Post")
	}

	if t.debugLog {
		resB, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(resB))
	}

	if res.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "CreatePonuda.Invalid http response: %s", res.Status)
	}

	return nil
}
