package shopify

import (
	"encoding/csv"
	"os"

	"github.com/pkg/errors"
)

type Order struct {
	//
}

func ParseExcel(filePath string) ([]Order, error) {

	return nil, nil
}

func parseCsv(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "file open %s", filePath)
	}
	reader := csv.NewReader(f)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrapf(err, "read CSV failed")
	}
	return records, nil
}
