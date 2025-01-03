package utils

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration time.Duration

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	x, err := time.ParseDuration(value.Value)
	if err != nil {
		fmt.Println("2->", err)
		return err
	}
	*d = Duration(x)
	return nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	unQuote, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	x, err := time.ParseDuration(unQuote)
	if err != nil {
		fmt.Println("->", err)
		return err
	}
	*d = Duration(x)
	return nil
}
