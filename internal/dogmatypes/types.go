package dogmatypes

import (
	"fmt"
	"go/types"
)

// PkgPath is the fully-qualified path to the Dogma package.
const PkgPath = "github.com/dogmatiq/dogma"

// Types is a container for types.Type representations of the types from the
// dogma package.
type Types struct {
	Application               *types.Interface
	AggregateMessageHandler   *types.Interface
	ProcessMessageHandler     *types.Interface
	IntegrationMessageHandler *types.Interface
	ProjectionMessageHandler  *types.Interface
}

// FromPackageImportedBy returns a Type structs containing the types the Dogma
// package imported by pkg.
func FromPackageImportedBy(pkg *types.Package) (*Types, bool, error) {
	for _, imp := range pkg.Imports() {
		if imp.Path() == PkgPath {
			t, err := FromPackage(imp)
			return t, true, err
		}
	}

	return nil, false, nil
}

// FromPackage returns a Types struct containing the types from the given
// package, which is assumed to be the Dogma package itself.
func FromPackage(pkg *types.Package) (t *Types, err error) {
	t = &Types{}

	t.Application, err = interfaceType(pkg, "Application")
	if err != nil {
		return nil, err
	}

	t.AggregateMessageHandler, err = interfaceType(pkg, "AggregateMessageHandler")
	if err != nil {
		return nil, err
	}

	t.ProcessMessageHandler, err = interfaceType(pkg, "ProcessMessageHandler")
	if err != nil {
		return nil, err
	}

	t.IntegrationMessageHandler, err = interfaceType(pkg, "IntegrationMessageHandler")
	if err != nil {
		return nil, err
	}

	t.ProjectionMessageHandler, err = interfaceType(pkg, "ProjectionMessageHandler")
	if err != nil {
		return nil, err
	}

	return t, nil
}

// interfaceType returns the Dogma interface with the given name.
func interfaceType(pkg *types.Package, n string) (*types.Interface, error) {
	obj := pkg.Scope().Lookup(n)
	if obj == nil {
		return nil, fmt.Errorf("type information for dogma.%s is not available", n)
	}

	iface, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, fmt.Errorf("dogma.%s is not an interface", n)
	}

	return iface, nil
}
