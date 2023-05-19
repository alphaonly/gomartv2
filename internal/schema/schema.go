package schema

import (
	"encoding/json"
	"fmt"
	"time"
)

type PreviousBytes []byte
type CtxUName string
type ContextKey int

const PKey1 ContextKey = 123455
const CtxKeyUName ContextKey = 1343456


type CreatedTime time.Time

func (t CreatedTime) MarshalJSON() ([]byte, error) {
	value := time.Time(t)
	//created, err := time.Parse(time.RFC3339, d.created_at.String)

	bytes, err := json.Marshal(value.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error marshal  CreatedTime %v", value)
	}
	return bytes, nil
}
func (t *CreatedTime) UnmarshalJSON(b []byte) error {
	var createdTimeString string
	err := json.Unmarshal(b, &createdTimeString)
	if err != nil {
		return fmt.Errorf("error unmarshal  CreatedTime %v", b)
	}
	createdTime, err := time.Parse(time.RFC3339, createdTimeString)
	if err != nil {
		return fmt.Errorf("error parse  CreatedTime to RFC3339 %v", b)
	}
	*t = CreatedTime(createdTime)
	return nil
}

// type Duration time.Duration

// func (d Duration) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(time.Duration(d).String())
// }

// func (d *Duration) UnmarshalJSON(b []byte) error {
// 	var v interface{}
// 	if err := json.Unmarshal(b, &v); err != nil {
// 		return err
// 	}
// 	switch value := v.(type) {
// 	case float64:
// 		*d = Duration(time.Duration(value))
// 		return nil
// 	case string:
// 		tmp, err := time.ParseDuration(value)
// 		if err != nil {
// 			return err
// 		}
// 		*d = Duration(tmp)
// 		return nil
// 	default:
// 		return errors.New("invalid duration")
// 	}
// }
