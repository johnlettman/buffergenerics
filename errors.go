package buffergenerics

import (
	"fmt"
	"reflect"
)

type ErrUnknownKind struct {
	error
	reflect.Kind
}

func NewErrUnknownKind(kind reflect.Kind) ErrUnknownKind {
	return ErrUnknownKind{
		error: fmt.Errorf("unknown kind: %v", kind),
		Kind:  kind,
	}
}
