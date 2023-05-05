package typeutils

import "fmt"

// TypeUtilsPackagePath should be set to the known package path for this package.
const TypeUtilsPackagePath = "github.com/madkins23/go-utils/typeutils"

//////////////////////////////////////////////////////////////////////////

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

//////////////////////////////////////////////////////////////////////////

// YAML works with a single copyFn that will copy from map to item and vice versa.

var copyFn func(dest, src interface{}) error

// JSON won't unmarshal to a map hidden within an interface so it requires two functions.

var copyMapFromItemFn func(toMap map[string]interface{}, fromItem interface{}) error
var copyItemFromMapFn func(toItem interface{}, fromMap map[string]interface{}) error

var copyFnMissing = fmt.Errorf("no copy function")

// Normally these would be more specific, but for testing purposes they are generic and weird.

func (a *alpha) PushToMap(toMap map[string]interface{}) error {
	if copyMapFromItemFn != nil {
		return copyMapFromItemFn(toMap, a)
	} else if copyFn != nil {
		return copyFn(toMap, a)
	}

	return copyFnMissing
}

func (a *alpha) PullFromMap(fromMap map[string]interface{}) error {
	if copyItemFromMapFn != nil {
		return copyItemFromMapFn(a, fromMap)
	} else if copyFn != nil {
		return copyFn(a, fromMap)
	}

	return copyFnMissing
}

func (b *bravo) PushToMap(toMap map[string]interface{}) error {
	if copyMapFromItemFn != nil {
		return copyMapFromItemFn(toMap, b)
	} else if copyFn != nil {
		return copyFn(toMap, b)
	}

	return copyFnMissing
}

func (b *bravo) PullFromMap(fromMap map[string]interface{}) error {
	if copyItemFromMapFn != nil {
		return copyItemFromMapFn(b, fromMap)
	} else if copyFn != nil {
		return copyFn(b, fromMap)
	}

	return copyFnMissing
}
