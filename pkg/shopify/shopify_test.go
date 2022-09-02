package shopify

import "testing"

func TestParseCsv(t *testing.T) {
	sh, err := ParseCsv("/home/karlo/Downloads/orders_export_1.csv")
	if err != nil {
		t.Fatal(err)
	}

	if len(sh) != 48 {
		t.Fatalf("Expected 48 orders")
	}

}
