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
	Product Product
	Err     error
}

type ImportAccountsError struct {
	Account Account
	Err     error
}

func NewMarket() Market {
	return Market{
		Accounts:    NewAccounts(),
		Products:    NewProducts(),
		Bundles:     NewBundles(),
		OrdersCache: NewOrdersCache(),
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
	//fixme: было бы здорово выделить это в отдельную функцию, которая принимает длину слайса, размер банча и возвращает [][2]int - массив индексов начала и конца подслайсов. Заодно эту логику тестами покрыть
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

		go ImportProductsCSVRecords(ctx, records[start:end], resChan, errChan)
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
	m.Products.mu.Lock()
	defer m.Products.mu.Unlock()
	// union all maps with our products map
	for key := range products {
		err := m.setProduct(key, products[key])
		if err != nil {
			return append(errs, ImportProductsError{products[key], fmt.Errorf("can't set product: %v", err)})
		}
	}

	return nil
}

func ImportProductsCSVRecords(
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
	for key := range m.Products.Products {
		export[key] = m.Products.Products[key]
	}

	return exportMapToCsv(export, reflect.ValueOf(Product{}))
}

// Accounts
func (m *Market) ImportAccountsCSV(data []byte) (errs []ImportAccountsError) { // fixme code duplicate

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll() // fixme read big data?
	if err == io.EOF {
		return append(errs, ImportAccountsError{Account{}, err})
	}
	if err != nil {
		return append(errs, ImportAccountsError{Account{}, errors.Wrap(err, "import product error")})
	}
	if len(records) < 2 {
		return append(errs, ImportAccountsError{Account{}, errors.New("empty data")})
	}

	records = records[1:]

	batchSize := 100
	length := len(records)
	rem := length % batchSize // remainder

	count := length / batchSize // goroutines count
	if rem > 0 {
		count++
	}

	resChan := make(chan map[string]Account, 1)
	errChan := make(chan ImportAccountsError, 1)

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

		go m.ImportAccountsCSVRecords(ctx, records[start:end], resChan, errChan)
	}

	// gather data
	accounts := make(map[string]Account)

	for i := 0; i < count; i++ {
		select {
		case result := <-resChan:
			// union two maps
			for key := range result {
				accounts[key] = result[key]
			}
		case err := <-errChan:
			cancel()
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return
	}

	// lock accounts to write
	m.Accounts.mu.Lock()
	defer m.Accounts.mu.Unlock()
	// union accounts maps
	for key := range accounts {
		err := m.setAccount(key, accounts[key])
		if err != nil {
			return append(errs, ImportAccountsError{Account{}, errors.Wrap(err, "accounts import error")})
		}
	}

	return nil
}

// fixme методов импорта и экспорта много. может стоит их разделить на 2 новых файла: shop_import.go, shop.export.go
func (m *Market) ImportAccountsCSVRecords(
	ctx context.Context,
	records [][]string,
	resChan chan<- map[string]Account,
	errChan chan<- ImportAccountsError) {

	accounts := make(map[string]Account)

	for _, record := range records { // read each record from csv
		select {
		case <-ctx.Done():
			errChan <- ImportAccountsError{Account{}, ErrorCancelled}
			return
		default:
		}

		if len(record) < 3 {
			errChan <- ImportAccountsError{Account{}, errors.New("not enough fields")}
			return
		}
		// first - account name
		userName := record[0]
		// second - balance
		balance, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			errChan <- ImportAccountsError{Account{userName, float32(balance), 0}, errors.Wrap(err, "parse error")}
			return
		}
		// third - account type
		typ, err := strconv.Atoi(record[2])
		if err != nil {
			errChan <- ImportAccountsError{Account{userName, float32(balance), AccountType(typ)}, errors.Wrap(err, "parse error")}
			return
		}
		// add new account to the temporary map
		account := Account{userName, float32(balance), AccountType(typ)}
		accounts[userName] = account
	}

	resChan <- accounts
}

func (m *Market) ExportAccountsCSV() ([]byte, error) {
	export := make(map[interface{}]interface{})
	for key := range m.Accounts.Accounts {
		export[key] = m.Accounts.Accounts[key]
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
