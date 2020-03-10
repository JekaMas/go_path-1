package cancel

import (
	"context"
	"errors"
	shop "go_path/struct"
)

var (
	ErrorCancelled = errors.New("operation cancelled")
)

type exportResult struct {
	bts []byte
	err error
}

func cancelExportFunc(f func(ch chan exportResult), cancelCtx context.Context) ([]byte, error) {
	ch := make(chan exportResult, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res.bts, res.err
	case <-cancelCtx.Done(): // cancel
		return []byte{}, ErrorCancelled
	}
}

func cancelImportFunc(f func(ch chan error), cancelCtx context.Context) error {
	ch := make(chan error, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res
	case <-cancelCtx.Done(): // cancel
		return ErrorCancelled
	}
}

func cancelImportProductsFunc(f func(ch chan []shop.ImportProductsError), cancelCtx context.Context) []shop.ImportProductsError {
	ch := make(chan []shop.ImportProductsError, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res
	case <-cancelCtx.Done(): // cancel
		return []shop.ImportProductsError{{shop.Product{}, ErrorCancelled}}
	}
}

func cancelImportAccountsFunc(f func(ch chan []shop.ImportAccountsError), cancelCtx context.Context) []shop.ImportAccountsError {
	ch := make(chan []shop.ImportAccountsError, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res
	case <-cancelCtx.Done(): // cancel
		return []shop.ImportAccountsError{{shop.Account{}, ErrorCancelled}}
	}
}

type CancelDecorator struct {
	shop.Shop
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
}

func NewCancelDecorator(s shop.Shop) CancelDecorator {
	ctx, cancel := context.WithCancel(context.Background())
	return CancelDecorator{
		Shop:       s,
		cancelCtx:  ctx,
		cancelFunc: cancel,
	}
}

/* --- Implementation ----------------------------------------------------------------------------------------------- */

func (c *CancelDecorator) CancelImportExport() error {
	c.cancelFunc()
	return nil // fixme hmmm...                ^^^^^
}

func (c CancelDecorator) Export() ([]byte, error) {
	return cancelExportFunc(func(ch chan exportResult) {
		bts, err := c.Shop.Export()
		ch <- exportResult{bts, err}
	}, c.cancelCtx)
}

func (c CancelDecorator) ExportProductsCSV() ([]byte, error) {
	return cancelExportFunc(func(ch chan exportResult) {
		bts, err := c.Shop.ExportProductsCSV()
		ch <- exportResult{bts, err}
	}, c.cancelCtx)
}

func (c CancelDecorator) ExportAccountsCSV() ([]byte, error) {
	return cancelExportFunc(func(ch chan exportResult) {
		bts, err := c.Shop.ExportAccountsCSV()
		ch <- exportResult{bts, err}
	}, c.cancelCtx)
}

func (c CancelDecorator) Import(data []byte) error {
	return cancelImportFunc(func(ch chan error) {
		ch <- c.Shop.Import(data)
	}, c.cancelCtx)
}

func (c CancelDecorator) ImportProductsCSV(data []byte) (errs []shop.ImportProductsError) {
	return cancelImportProductsFunc(func(ch chan []shop.ImportProductsError) {
		ch <- c.Shop.ImportProductsCSV(data)
	}, c.cancelCtx)
}

func (c CancelDecorator) ImportAccountsCSV(data []byte) (errs []shop.ImportAccountsError) {
	return cancelImportAccountsFunc(func(ch chan []shop.ImportAccountsError) {
		ch <- c.Shop.ImportAccountsCSV(data)
	}, c.cancelCtx)
}
