package cligo

import (
	"fmt"
	"reflect"
	"strings"
)

type flagDescriptor struct {
	Flag        string
	Short       string
	Description string
	Required    bool
	Default     string
	Env         string
}

func parseFlagTag(structReflect reflect.Value, index int) (flagDescriptor, bool) {

	field := structReflect.Type().Field(index)
	name := field.Name
	tag := field.Tag.Get("cligo")
	if tag == "" {
		return flagDescriptor{}, false
	}

	fd := flagDescriptor{
		Flag: fmt.Sprintf("--%s", toHyphenatedString(name)),
	}
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "required" {
			fd.Required = true
		} else if strings.HasPrefix(part, "flag=") {
			fd.Flag = fmt.Sprintf("--%s", toHyphenatedString(strings.TrimPrefix(part, "flag=")))
		} else if strings.HasPrefix(part, "short=") {
			fd.Short = "-" + strings.TrimPrefix(part, "short=")
		} else if strings.HasPrefix(part, "description=") {
			fd.Description = strings.TrimPrefix(part, "description=")
		} else if strings.HasPrefix(part, "env=") {
			fd.Env = strings.TrimPrefix(part, "env=")
			fd.Default = fmt.Sprintf("$%s", fd.Env)
		}
	}

	value := structReflect.Field(index)
	if !value.IsZero() {
		fd.Default = fmt.Sprintf("%v", value.Interface())
	}
	return fd, true
}
