// Code generated by "stringer -type=Type"; DO NOT EDIT

package scan

import "fmt"

const _Type_name = "ILLEGALEOFATOMLPARENRPAREN"

var _Type_index = [...]uint8{0, 7, 10, 14, 20, 26}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return fmt.Sprintf("Type(%d)", i)
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
