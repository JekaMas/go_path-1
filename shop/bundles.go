package shop

import (
	"github.com/pkg/errors"
	"sync"
)

type Bundles struct {
	Bundles map[string]Bundle
	mu      sync.RWMutex
}

type changeBundleFunc func(bundle *Bundle)

func NewBundles() Bundles {
	return Bundles{
		Bundles: make(map[string]Bundle),
		mu:      sync.RWMutex{},
	}
}

func (b *Bundles) AddBundle(name string, main Product, discount float32, additional ...Product) error {

	if main.Type == ProductSampled {
		return ErrorBundleMainSample
	}

	sampled := getProductsWithType(additional, ProductSampled)

	bundleType := BundleNormal
	if len(sampled) == 1 {
		bundleType = BundleSample
	}

	bundle := NewBundle(main, discount, bundleType, additional...)

	if _, err := b.GetBundle(name); err == nil {
		return ErrorBundleExists
	}
	return b.SetBundle(name, bundle)
}

func (b *Bundles) RemoveBundle(name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, err := b.getBundle(name); err != nil {
		return errors.Wrap(err, "can't delete nil bundle")
	}

	delete(b.Bundles, name)
	return nil
}

func (b *Bundles) GetBundle(name string) (Bundle, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.getBundle(name)
}

func (b *Bundles) getBundle(name string) (Bundle, error) {
	bundle, ok := b.Bundles[name]

	if !ok {
		return Bundle{}, ErrorBundleNotExists
	}

	return bundle, nil
}

func (b *Bundles) SetBundle(name string, bundle Bundle) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.setBundle(name, bundle)
}

func (b *Bundles) setBundle(name string, bundle Bundle) error {

	if err := checkName(name); err != nil {
		return errors.Wrap(err, "can't set invalid bundle")
	}
	if err := checkBundle(bundle); err != nil {
		return errors.Wrap(err, "can't set invalid bundle")
	}

	b.Bundles[name] = bundle
	return nil
}

func (b *Bundles) changeBundle(name string, fns ...changeBundleFunc) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	acc, err := b.getBundle(name)
	if err != nil {
		return errors.Wrap(err, "can't change the nil bundle")
	}

	for _, fn := range fns {
		fn(&acc)
	}

	b.Bundles[name] = acc
	return nil
}

func changeDiscountFunc(discount float32) changeBundleFunc {
	return func(bundle *Bundle) {
		bundle.Discount = discount
	}
}
