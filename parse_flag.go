package configo

import (
	"fmt"

	"github.com/han0110/configo/util"
)

// ParseFlag parses arguments into a map.
func ParseFlag(args []string) (*util.FlattenMap, error) {
	f := flagSet{
		args:        args,
		values:      make(map[string]string),
		sliceValues: make(map[string][]string),
	}

	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		return nil, err
	}

	flattenMap := util.NewFlattenMap()
	for _, key := range f.keys {
		if values, ok := f.sliceValues[key]; ok {
			delete(f.values, key)
			for index, value := range values {
				flattenMap.Set(fmt.Sprintf("%s%s%d", key, util.CharDot, index), value)
			}
			continue
		}
		flattenMap.Set(key, f.values[key])
	}

	return flattenMap, nil
}

type flagSet struct {
	args        []string
	keys        []string
	values      map[string]string
	sliceValues map[string][]string
}

func (f *flagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}

	arg := f.args[0]
	if len(arg) < 2 || arg[0] != '-' {
		return false, nil
	}
	numMinuses := 1
	if arg[1] == '-' {
		numMinuses = 2
		if arg == "--" { // "--" terminates the flags
			f.args = f.args[1:]
			return false, nil
		}
	}

	arg = arg[numMinuses:]
	if arg[0] == '-' || arg[0] == '=' {
		return false, fmt.Errorf("bad flag syntax: %s", f.args[0])
	}

	// It's a flag. Does it have an argument?
	f.args = f.args[1:]
	for i := 1; i < len(arg); i++ { // '=' will not be first
		if arg[i] == '=' {
			f.setValue(arg[0:i], arg[i+1:])
			return true, nil
		}
	}

	// End of the arguments, consider it to be a boolean flag
	if len(f.args) == 0 {
		f.setValue(arg, "true")
		return true, nil
	}

	// Looking for next argument
	nextArg, value := f.args[0], ""
	if nextArg[0] == '-' {
		// Start with '-', consider it to be another flag's arg,
		// and current flag to be a boolean flag
		value = "true"
	} else {
		// Take the next argument as value, then shift.
		value, f.args = nextArg, f.args[1:]
	}
	f.setValue(arg, value)
	return true, nil
}

func (f *flagSet) setValue(key, value string) {
	if val, exist := f.values[key]; exist {
		if len(f.sliceValues[key]) == 0 { // First element has not appended into slice
			f.sliceValues[key] = []string{val, value}
		} else {
			f.sliceValues[key] = append(f.sliceValues[key], value)
		}
		return
	}
	f.values[key] = value
	f.keys = append(f.keys, key)
}
