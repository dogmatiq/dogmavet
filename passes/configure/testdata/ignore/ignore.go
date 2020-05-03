package ignore

import (
	"io"

	"github.com/dogmatiq/dogma"
)

func nonMethod() {}

type (
	wrongName     struct{}
	wrongArity    struct{}
	unnamedType   struct{}
	notDogmaType  struct{}
	notConfigurer struct{}
)

func (wrongName) NotNamedConfigure()                          {}
func (wrongArity) Configure(a, b dogma.ApplicationConfigurer) {}
func (unnamedType) Configure(interface{})                     {}
func (notDogmaType) Configure(io.Writer)                      {}
func (notConfigurer) Configure(dogma.Application)             {}
