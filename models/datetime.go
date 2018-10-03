package models

import "time"

type DateTime time.Time

func (d DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte(t.UTC().Format(time.RFC3339)), nil
}

type ReadonlyString string

func (s *ReadonlyString) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}
