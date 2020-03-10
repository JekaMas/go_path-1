package shop

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"runtime"
	"strconv"
	"testing"
)

func TestMarket_ExportProductsCSV(t *testing.T) {
	m := testExportMarket(10, 0)
	spew.Dump(m.ExportProductsCSV())
}

func TestMarket_ExportAccountsCSV(t *testing.T) {
	m := testExportMarket(0, 10)
	spew.Dump(m.ExportAccountsCSV())
}

func TestMarket_ImportProductsCSV(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU()) // max procs

	m := testExportMarket(10_000, 0)
	m2 := NewMarket()

	data, _ := m.ExportProductsCSV()
	errs := m2.ImportProductsCSV(data)

	for _, err := range errs {
		t.Errorf("Error: '%v' Product: %v", err.Err, err.Product)
	}

	if !reflect.DeepEqual(m, m2) {
		t.Errorf("Imported and exported markets are not equals!")
	}

	cases := [][]byte{
		{}, {0}, nil,
	}

	for _, cas := range cases {
		errs = m2.ImportProductsCSV(cas)

		if len(errs) == 0 {
			t.Errorf("Corrupted data success import: %v", cas)
		}
	}

	data, _ = m.ExportProductsCSV()
	data[100] = '\000'
	errs = m2.ImportProductsCSV(data)
	if len(errs) == 0 {
		t.Errorf("Corrupted data success import")
	}
	//else {
	//	spew.Dump(errs)
	//}

	data, _ = m.ExportProductsCSV()
	data = data[100:]
	errs = m2.ImportProductsCSV(data)

	if len(errs) == 0 {
		t.Errorf("Corrupted data success import")
	}
}

func TestMarket_ImportAccountsCSV(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU()) // max procs

	m := testExportMarket(0, 300)
	m2 := NewMarket()

	data, _ := m.ExportAccountsCSV()
	errs := m2.ImportAccountsCSV(data)

	for _, err := range errs {
		t.Errorf("Error: '%v' Account: %v", err.Err, err.Account)
	}

	if !reflect.DeepEqual(m, m2) {
		t.Errorf("Imported and exported markets are not equals!")
	}

	cases := [][]byte{
		{}, {0}, nil,
	}

	for _, cas := range cases {
		errs = m2.ImportAccountsCSV(cas)

		if len(errs) == 0 {
			t.Errorf("Corrupted data success import: %v", cas)
		}
	}

	data, _ = m.ExportAccountsCSV()
	data[25] = '\000'
	errs = m2.ImportAccountsCSV(data)

	if len(errs) == 0 {
		t.Errorf("Corrupted data success import")
	}

	data, _ = m.ExportAccountsCSV()
	data = data[100:]
	errs = m2.ImportAccountsCSV(data)

	if len(errs) == 0 {
		t.Errorf("Corrupted data success import")
	}
}

func testExportMarket(productsN int, accountsN int) Market {

	m := NewMarket()
	baseName := "Bimbo"

	for i := 1; i <= productsN; i++ {
		p := NewProduct(baseName+strconv.Itoa(i), float32(i), ProductNormal)
		err := m.AddProduct(p)

		if err != nil {
			fmt.Errorf("error while genereting test market: %v", err)
		}
	}

	baseName = "Jc"
	for i := 1; i <= accountsN; i++ {
		userName := baseName + strconv.Itoa(i)
		err := m.Register(userName)
		if err != nil {
			fmt.Errorf("error while genereting test market: %v", err)
		}

		err = m.AddBalance(userName, float32(i))
		if err != nil {
			fmt.Errorf("error while genereting test market: %v", err)
		}
	}

	return m
}
