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
