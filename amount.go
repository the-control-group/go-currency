package currency

// stdlib
import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type Amount struct {
	Nil     bool
	Dollars int
	Cents   int
}

func (a *Amount) Float64() float64 {
	return float64(a.Dollars) + float64(a.Cents) / 100
}

func (a *Amount) ToString() string {
	return fmt.Sprintf("%d.%02d", a.Dollars, a.Cents)
}

func (a *Amount) String() string {
	return fmt.Sprintf("%d.%02d", a.Dollars, a.Cents)
}

func (a *Amount) MarshalText() ([]byte, error) {
	if a.Nil {
		return nil, nil
	}
	return []byte(fmt.Sprintf("%d.%02d", a.Dollars, a.Cents)), nil
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	if a.Nil {
		return []byte(`null`), nil
	}
	return []byte(fmt.Sprintf(`%d.%02d`, a.Dollars, a.Cents)), nil
}

func (a *Amount) MarshalRQL() (interface{}, error) {
	return a.MarshalJSON()
}

func (a *Amount) UnmarshalJSON(data []byte) error {
	if bytes.Compare([]byte(`null`), bytes.ToLower(data)) == 0 {
		a.Nil = true
		return nil
	}
	_, err := fmt.Sscanf(string(data), `%d.%d`, &a.Dollars, &a.Cents)
	if err != nil {
		_, err = fmt.Sscanf(string(data), `%d`, &a.Dollars)
		if err == nil {
			a.Cents = 0
		}
	}
	return err
}

func (a *Amount) UnmarshalText(data []byte) error {
	if bytes.Compare([]byte(`null`), bytes.ToLower(data)) == 0 {
		a.Nil = true
		return nil
	}
	_, err := fmt.Sscanf(string(data), `%d.%d`, &a.Dollars, &a.Cents)
	return err
}

func (a *Amount) UnmarshalRQL(data interface{}) error {
	return a.Scan(data)
}

func (a *Amount) Scan(src interface{}) error {
	switch src.(type) {
	case float64:
		str := strconv.FormatFloat(src.(float64), 'f', 2, 64)
		return a.UnmarshalText([]byte(str))
	case string:
		return a.UnmarshalText([]byte(src.(string)))
	case []byte:
		return a.UnmarshalText(src.([]byte))
	case nil:
		a.Nil = true
		return nil
	default:
		return errors.New("Unexpected type")
	}
}
