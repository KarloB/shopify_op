package shopify

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/pkg/errors"
)

type Order struct {
	ID       string // "Name"
	Customer Customer
	Items    []OrderItem
	Shipping float64
}

type Customer struct {
	Email    string // "Email"
	Name     string // "Billing Name" "Shipping Name"
	Street   string // "Billing Street" "Shipping Street"
	Postcode string // "Billing Zip" "Shipping Zip"
	City     string // "Billing City" "Shipping City"
	Phone    string // "Billing Phone" "Shipping Phone"
}

type OrderItem struct {
	Name     string  // "Lineitem name"
	Price    float64 // "Lineitem price"
	Quantity float64
}

func ParseCsv(filePath string) ([]Order, error) {
	csvRows, err := parseCsv(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "ParseCsv")
	}

	rowGroup := make(map[string][][]string)

	headerIndex := make(map[string]int)
	for i := range csvRows {
		if i == 0 {
			for j := range csvRows[i] {
				headerIndex[csvRows[i][j]] = j
			}
			continue
		}
		curr := rowGroup[csvRows[i][headerIndex["Name"]]]
		curr = append(curr, csvRows[i])
		rowGroup[csvRows[i][headerIndex["Name"]]] = curr
	}

	var orders []Order
	for key, rows := range rowGroup {
		customer := getCustomer(rows, headerIndex)
		items := getOrderItems(rows, headerIndex)

		o := Order{
			ID:       key,
			Customer: customer,
			Items:    items,
			Shipping: getShipping(rows, headerIndex),
		}
		orders = append(orders, o)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})

	for i, o := range orders {
		fmt.Printf("%d %v\n", i, o.ID)
	}

	return orders, nil
}

func getShipping(rows [][]string, headerIndex map[string]int) float64 {
	for _, row := range rows {
		price, _ := strconv.ParseFloat(row[headerIndex["Shipping"]], 64)
		if price != 0 {
			return price
		}
	}

	return 0
}

func getOrderItems(rows [][]string, headerIndex map[string]int) []OrderItem {
	var res []OrderItem
	for _, row := range rows {
		price, _ := strconv.ParseFloat(row[headerIndex["Lineitem price"]], 64)
		qty, _ := strconv.ParseFloat(row[headerIndex["Lineitem quantity"]], 64)
		item := OrderItem{
			Name:     row[headerIndex["Lineitem name"]],
			Price:    price,
			Quantity: qty,
		}
		res = append(res, item)
	}

	return res
}

func getCustomer(rows [][]string, headerIndex map[string]int) Customer {
	customer := Customer{}
	for _, row := range rows {
		if email := row[headerIndex["Email"]]; email != "" {
			customer.Email = email
		}
		if name := row[headerIndex["Billing Name"]]; name != "" {
			customer.Name = name
		}
		if s := row[headerIndex["Billing Street"]]; s != "" {
			customer.Street = s
		}
		if z := row[headerIndex["Billing Zip"]]; z != "" {
			customer.Postcode = z
		}
		if c := row[headerIndex["Billing City"]]; c != "" {
			customer.City = c
		}
		if c := row[headerIndex["Billing Phone"]]; c != "" {
			customer.Phone = c
		}
		if name := row[headerIndex["Shipping Name"]]; name != "" {
			customer.Name = name
		}
		if s := row[headerIndex["Shipping Street"]]; s != "" {
			customer.Street = s
		}
		if z := row[headerIndex["Shipping Zip"]]; z != "" {
			customer.Postcode = z
		}
		if c := row[headerIndex["Shipping City"]]; c != "" {
			customer.City = c
		}
		if c := row[headerIndex["Shipping Phone"]]; c != "" {
			customer.Phone = c
		}
	}
	return customer
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
