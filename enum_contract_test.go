package protobuf

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/prothegee/erposfnb-protobuf/common"
	"github.com/prothegee/erposfnb-protobuf/pos"
)

// enumFiles lists every proto file that declares an enum.
func enumFiles() map[string]protoreflect.FileDescriptor {
	return map[string]protoreflect.FileDescriptor{
		"common/common.proto": common.File_common_common_proto,
		"pos/pos.proto":       pos.File_pos_pos_proto,
	}
}

// upperSnake turns an enum name into its value name prefix.
// For example MeasureType becomes MEASURE_TYPE.
func upperSnake(name string) string {
	var builder strings.Builder

	for index, character := range name {
		if index > 0 && character >= 'A' && character <= 'Z' {
			builder.WriteByte('_')
		}

		builder.WriteRune(character)
	}

	return strings.ToUpper(builder.String())
}

// TestEnumTagZeroIsUnspecified checks that tag 0 of every enum is the reserved
// value, so a real value can never be mistaken for an unset field.
func TestEnumTagZeroIsUnspecified(t *testing.T) {
	for path, file := range enumFiles() {
		enums := file.Enums()

		for index := 0; index < enums.Len(); index++ {
			enum := enums.Get(index)
			values := enum.Values()

			for valueIndex := 0; valueIndex < values.Len(); valueIndex++ {
				value := values.Get(valueIndex)
				if value.Number() != 0 {
					continue
				}

				if !strings.HasSuffix(string(value.Name()), "_UNSPECIFIED") {
					t.Fatalf("%s: %s tag 0 is %s, want a name ending in _UNSPECIFIED", path, enum.Name(), value.Name())
				}
			}
		}
	}
}

// TestEnumValueNameUsesEnumPrefix checks that every value carries the enum name
// as a prefix, which keeps one enum value from reading like another enum value.
func TestEnumValueNameUsesEnumPrefix(t *testing.T) {
	for path, file := range enumFiles() {
		enums := file.Enums()

		for index := 0; index < enums.Len(); index++ {
			enum := enums.Get(index)
			values := enum.Values()
			prefix := upperSnake(string(enum.Name())) + "_"

			for valueIndex := 0; valueIndex < values.Len(); valueIndex++ {
				value := values.Get(valueIndex)
				if !strings.HasPrefix(string(value.Name()), prefix) {
					t.Fatalf("%s: %s value %s does not start with %s", path, enum.Name(), value.Name(), prefix)
				}
			}
		}
	}
}

// TestEnumTagNumberIsAscending checks that non zero tags climb in declaration
// order, so a new value is appended with a tag higher than every tag before it.
// A reserved gap is allowed, a repeat or a smaller tag is not.
func TestEnumTagNumberIsAscending(t *testing.T) {
	for path, file := range enumFiles() {
		enums := file.Enums()

		for index := 0; index < enums.Len(); index++ {
			enum := enums.Get(index)
			values := enum.Values()
			previous := protoreflect.EnumNumber(0)

			for valueIndex := 0; valueIndex < values.Len(); valueIndex++ {
				value := values.Get(valueIndex)
				if value.Number() == 0 {
					continue
				}

				if value.Number() <= previous {
					t.Fatalf("%s: %s tag %d does not climb past %d", path, enum.Name(), value.Number(), previous)
				}

				previous = value.Number()
			}
		}
	}
}
