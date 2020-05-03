package application

import (
	"strings"

	"github.com/dogmatiq/dogma"
)

type (
	nonConstant      struct{}
	invalidValues    struct{}
	nonUUIDKey       struct{}
	alternateUUIDKey struct{}
)

// Configure calls c.Identity() with non-constant values in order to verify that
// the check does not fail if the values are not known at compile time.
func (nonConstant) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(
		strings.ToUpper("name"),
		strings.ToUpper("key"),
	)
}

// Configure calls c.Identity() with invalid values. It uses validation from
// configkit.NewIdentity().
func (invalidValues) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(
		"", // want `invalid name "", names must be non-empty, printable UTF-8 strings with no whitespace`
		"", // want `invalid key "", keys must be non-empty, printable UTF-8 strings with no whitespace`
	)
}

// Configure calls c.Identity() with a non-UUID key.
func (nonUUIDKey) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(
		"name",
		"key", // want `identity keys should be UUIDs \(invalid UUID length: 3\)`
	)
}

// Configure calls c.Identity() with a UUID key that is given in a non-standard
// format.
func (alternateUUIDKey) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(
		"name",
		"{bc5c4138-9ead-4d17-a425-88d0e4cb3059}", // want `identity key UUIDs should use RFC-4122 hex notation \(xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\)`
	)
}
