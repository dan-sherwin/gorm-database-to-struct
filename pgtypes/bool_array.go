package pgtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

type BoolArray []bool

func (a *BoolArray) Scan(src interface{}) error {
	if src == nil {
		*a = nil
		return nil
	}
	var input string
	switch t := src.(type) {
	case []byte:
		input = string(t)
	case string:
		input = t
	default:
		return fmt.Errorf("cannot scan type %T into BoolArray", src)
	}
	input = strings.Trim(input, "{}")
	if input == "" {
		*a = BoolArray{}
		return nil
	}
	parts := strings.Split(input, ",")
	result := make(BoolArray, len(parts))
	for i, p := range parts {
		switch strings.ToLower(strings.TrimSpace(p)) {
		case "t", "true":
			result[i] = true
		case "f", "false":
			result[i] = false
		default:
			return fmt.Errorf("invalid boolean value: %s", p)
		}
	}
	*a = result
	return nil
}

func (a BoolArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	strs := make([]string, len(a))
	for i, v := range a {
		if v {
			strs[i] = "t"
		} else {
			strs[i] = "f"
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(strs, ",")), nil
}

func (a BoolArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]bool(a))
}

func (a *BoolArray) UnmarshalJSON(data []byte) error {
	var tmp []bool
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*a = BoolArray(tmp)
	return nil
}

func (a BoolArray) MarshalText() ([]byte, error) {
	strs := make([]string, len(a))
	for i, v := range a {
		if v {
			strs[i] = "true"
		} else {
			strs[i] = "false"
		}
	}
	return []byte(strings.Join(strs, ",")), nil
}

func (a *BoolArray) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*a = BoolArray{}
		return nil
	}
	parts := strings.Split(string(data), ",")
	out := make(BoolArray, len(parts))
	for i, s := range parts {
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "true", "t":
			out[i] = true
		case "false", "f":
			out[i] = false
		default:
			return fmt.Errorf("invalid boolean value: %s", s)
		}
	}
	*a = out
	return nil
}

func (BoolArray) GormDataType() string {
	return "boolean[]"
}

func (BoolArray) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	if db.Dialector.Name() == "postgres" {
		return "boolean[]"
	}
	return ""
}

func (BoolArray) FromSlice(s []bool) BoolArray {
	return BoolArray(s)
}

func (a BoolArray) AsSlice() []bool {
	return []bool(a)
}

func (a BoolArray) String() string {
	strs := make([]string, len(a))
	for i, v := range a {
		if v {
			strs[i] = "true"
		} else {
			strs[i] = "false"
		}
	}
	return strings.Join(strs, ",")
}

func (a BoolArray) Len() int           { return len(a) }
func (a BoolArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BoolArray) Less(i, j int) bool { return !a[i] && a[j] }

func (a BoolArray) Contains(val bool) bool {
	for _, x := range a {
		if x == val {
			return true
		}
	}
	return false
}

func (a BoolArray) IndexOf(val bool) int {
	for i, x := range a {
		if x == val {
			return i
		}
	}
	return -1
}

func (a BoolArray) IsEmpty() bool {
	return len(a) == 0
}

func (a BoolArray) Unique() BoolArray {
	var out BoolArray
	seen := map[bool]bool{}
	for _, v := range a {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func (a BoolArray) Filter(f func(bool) bool) BoolArray {
	var out BoolArray
	for _, v := range a {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}

func (a BoolArray) Append(vals ...bool) BoolArray {
	return append(a, vals...)
}

func (a BoolArray) Equals(b BoolArray) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
