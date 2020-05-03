package application

import "github.com/dogmatiq/dogma"

type (
	unnamed  struct{}
	misnamed struct{}
)

func (unnamed) Configure(dogma.ApplicationConfigurer) {} // want `configurer parameter should be named 'c'`

func (misnamed) Configure(x dogma.ApplicationConfigurer) { // want `configurer parameter should be named 'c'`
	x.Identity("name", "bf63e337-eaae-4459-8a6a-deca84e9c672")
}
