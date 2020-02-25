package shop

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"reflect"
	"strconv"
	"unicode"
)

var (
	ErrorInvalidDiscount = errors.New("invalid discount")
	ErrorEmptyField      = errors.New("empty field")
	ErrorCancelled       = errors.New("cancelled")
)

type ImportProductsError struct {
	product Product
	err     error
}

type ImportAccountsError struct {
	product Product
	err     error
}

func NewMarket() Market {
	return Market{
		Accounts:    make(map[string]Account),
		Products:    make(map[string]Product),
		Bundles:     make(map[string]Bundle),
		OrdersCache: make(map[string]float32),
	}
}

/* -- Importer, Exporter -------------------------------------------------------------------------------------------- */

func (m *Market) Import(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Market) Export() ([]byte, error) {
	return json.MarshalIndent(m, "", "")
}

// Products
func (m *Market) ImportProductsCSV(data []byte) (errs []ImportProductsError) {
	// create new reader
	reader := csv.NewReader(bytes.NewReader(data))
	// read all data at once
	records, err := reader.ReadAll() // fixme read big data?
	// end of file at the beginning
	if err == io.EOF {
		return append(errs, ImportProductsError{Product{}, err})
	}
	// io error
	if err != nil {
		return append(errs, ImportProductsError{Product{}, errors.Wrap(err, "import product error")})
	}
	// with headers
	if len(records) < 2 {
		return append(errs, ImportProductsError{Product{}, errors.New("empty data")})
	}
	// skip headers
	records = records[1:]
	// steps per goroutine
	batchSize := 1000
	length := len(records)
	rem := length % batchSize // remainder
	// goroutines count
	count := length / batchSize
	if rem > 0 {
		count++
	}
	// result channels
	resChan := make(chan map[string]Product, 1)
	errChan := make(chan ImportProductsError, 1)
	// context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < length; i += batchSize {

		start := i
		end := start

		if length < batchSize {
			end = rem
		} else if i+batchSize > length {
			end += rem
		} else {
			end += batchSize
		}

		go m.ImportProductsCSVRecords(ctx, records[start:end], resChan, errChan)
	}

	// gather data
	products := make(map[string]Product)

	for i := 0; i < count; i++ {
		select {
		case result := <-resChan:
			// union two maps
			for key := range result {
				products[key] = result[key]
			}
		case err := <-errChan:
			cancel()
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return
	}

	// finish
	m.productsMutex.Lock()
	defer m.productsMutex.Unlock()
	// union all maps with our products map
	for key := range products {
		m.Products[key] = products[key]
	}

	return nil
}

func (m *Market) ImportProductsCSVRecords(
	ctx context.Context,
	records [][]string,
	resChan chan<- map[string]Product,
	errChan chan<- ImportProductsError) {

	products := make(map[string]Product)

	for _, record := range records { // read each record from csv
		select {
		case <-ctx.Done():
			errChan <- ImportProductsError{Product{}, ErrorCancelled}
			return
		default:
		}

		if len(record) < 3 {
			errChan <- ImportProductsError{Product{}, errors.New("not enough fields")}
			return
		}
		// first - name
		name := record[0]
		// second - price
		price, err := strconv.ParseFloat(record[1], 32)
		if err != nil { // fixme no full information?
			errChan <- ImportProductsError{Product{name, float32(price), 0}, errors.Wrap(err, "parse error")}
			return
		}
		// third - type
		typ, err := strconv.Atoi(record[2])
		if err != nil { // fixme no full information?
			errChan <- ImportProductsError{Product{name, float32(price), ProductType(typ)}, errors.Wrap(err, "parse error")}
			return
		}
		// add new product to temporary map
		product := NewProduct(name, float32(price), ProductType(typ))
		products[product.Name] = product
	}

	resChan <- products
}

func (m *Market) ExportProductsCSV() ([]byte, error) {
	export := make(map[interface{}]interface{}) // fixme no generics, nice
	for key := range m.Products {
		export[key] = m.Products[key]
	}

	return exportMapToCsv(export, reflect.ValueOf(Product{}))
}

// Accounts
func (m *Market) ImportAccountsCSV(data []byte) error {

	accounts := make(map[string]Account)

	reader := csv.NewReader(bytes.NewReader(data))
	// headers
	_, err := reader.Read()
	// end of file at beginning
	if err == io.EOF {
		return nil
	}
	// io error
	if err != nil {
		return errors.Wrap(err, "import product error")
	}

	for { // read each record from csv
		record, err := reader.Read()
		// end of file
		if err == io.EOF {
			break
		}
		// io error
		if err != nil {
			return errors.Wrap(err, "import product error")
		}

		// first - account name
		name := record[0]
		// second - balance
		balance, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			return errors.Wrap(err, "parse error")
		}
		// third - account type
		typ, err := strconv.Atoi(record[2])
		if err != nil {
			return errors.Wrap(err, "parse error")
		}
		// add new product to temporary map
		account := Account{name, float32(balance), AccountType(typ)}
		accounts[account.Name] = account
	}

	// lock accounts to write
	m.accountsMutex.Lock()
	defer m.accountsMutex.Unlock()
	// union accounts maps
	for key := range accounts {
		m.Accounts[key] = accounts[key]
	}

	return nil
}

func (m *Market) ExportAccountsCSV() ([]byte, error) {
	export := make(map[interface{}]interface{})
	for key := range m.Accounts {
		export[key] = m.Accounts[key]
	}

	return exportMapToCsv(export, reflect.ValueOf(Product{}))
}

/* --- Utils -------------------------------------------------------------------------------------------------------- */

func exportMapToCsv(m map[interface{}]interface{}, typeValue reflect.Value) ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := csv.NewWriter(&buffer)

	// create slice of headers
	var headers []string
	for i := 0; i < typeValue.NumField(); i++ {
		headers = append(headers, typeValue.Type().Field(i).Name)
	}

	err := writer.Write(headers)
	if err != nil {
		return nil, errors.Wrap(err, "can't write headers")
	}

	for key := range m {
		var record []string
		// go through struct values
		typeValue = reflect.ValueOf(m[key])
		for i := 0; i < typeValue.NumField(); i++ {
			field := typeValue.Field(i).Interface()
			var value string
			// print full float representation
			switch field.(type) {
			case float64, float32:
				value = fmt.Sprintf("%f", field)
			default:
				value = fmt.Sprintf("%v", field)
			}

			record = append(record, value)
		}

		err = writer.Write(record)
		if err != nil {
			return nil, errors.Wrap(err, "can't writer record")
		}
	}

	writer.Flush()
	return buffer.Bytes(), writer.Error()
}

/* --- Checks ------------------------------------------------------------------------------------------------------- */

func checkName(name string) error {
	if len(name) == 0 {
		return ErrorEmptyField
	}

	// TODO max chars count
	//if len(userName) > MAX_NAME_LENGTH {
	//	return ErrorAccountInvalidName
	//}

	// TODO spaces at beginning
	for _, r := range name { // for each rune
		if !(unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			unicode.IsSpace(r)) {
			return ErrorAccountInvalidName
		}
	}

	return nil
}
