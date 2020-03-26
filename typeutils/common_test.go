package typeutils

import "fmt"

type actor interface {
	declaim() string
}

//////////////////////////////////////////////////////////////////////////

type alpha struct {
	Name    string
	Percent float32 `yaml:"percentDone"`
	extra   string
}

func (a *alpha) declaim() string {
	return fmt.Sprintf("%s is %6.2f%%  complete", a.Name, a.Percent)
}

type bravo struct {
	Finished   bool
	Iterations int
	extra      string
}

func (b *bravo) declaim() string {
	var finished string
	if !b.Finished {
		finished = "not "
	}
	return fmt.Sprintf("%sfinished after %d iterations", finished, b.Iterations)
}

// YAML works with a single copyFn that will copy from map to item and vice versa.

var copyFn func(dest, src interface{}) error

// JSON won't unmarshal to a map hidden within an interface so it requires two functions.

var copyMapFromItemFn func(toMap map[string]interface{}, fromItem interface{}) error
var copyItemFromMapFn func(toItem interface{}, fromMap map[string]interface{}) error

var copyFnMissing = fmt.Errorf("no copy function")

// Normally these would be more specific, but for testing purposes they are generic and weird.

func (alf *alpha) PushToMap(toMap map[string]interface{}) error {
	if copyMapFromItemFn != nil {
		return copyMapFromItemFn(toMap, alf)
	} else if copyFn != nil {
		return copyFn(toMap, alf)
	}

	return copyFnMissing
}

func (alf *alpha) PullFromMap(fromMap map[string]interface{}) error {
	if copyItemFromMapFn != nil {
		return copyItemFromMapFn(alf, fromMap)
	} else if copyFn != nil {
		return copyFn(alf, fromMap)
	}

	return copyFnMissing
}

func (bra *bravo) PushToMap(toMap map[string]interface{}) error {
	if copyMapFromItemFn != nil {
		return copyMapFromItemFn(toMap, bra)
	} else if copyFn != nil {
		return copyFn(toMap, bra)
	}

	return copyFnMissing
}

func (bra *bravo) PullFromMap(fromMap map[string]interface{}) error {
	if copyItemFromMapFn != nil {
		return copyItemFromMapFn(bra, fromMap)
	} else if copyFn != nil {
		return copyFn(bra, fromMap)
	}

	return copyFnMissing
}
