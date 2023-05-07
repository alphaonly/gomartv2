// Package configuration - describes the configuration of application
package configuration

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// ServerDefaultJSON - default values of configuration parameters in JSON format
const ServerDefaultJSON = `{
"RUN_ADDRESS":"localhost:8080",
"DATABASE_URI": "postgres://postgres:mypassword@localhost:5432/yandex",
"ACCRUAL_SYSTEM_ADDRESS":"http://localhost:8080",
"RESTORE":true,"KEY":"",
"ACCRUAL_TIME":200
}`

// ServerConfiguration  - a structure that describes conf parameters.
type ServerConfiguration struct {
	RunAddress           string          `json:"RUN_ADDRESS,omitempty"`            // Host and port to receive HTTP requests
	Port                 string          `json:"PORT,omitempty"`                   // a port derives from RunAddress
	DatabaseURI          string          `json:"DATABASE_URI,omitempty"`           // an URI to database
	AccrualSystemAddress string          `json:"ACCRUAL_SYSTEM_ADDRESS,omitempty"` // An address to a remote accrual system
	AccrualTime          int64           `json:"ACCRUAL_TIME,omitempty"`           // a time to update data from accrual system
	EnvChanged           map[string]bool // technical map to check which parameters were given in app start
}

// ServerConfigurationOption  - a type to implement a Options pattern for giving conf parameters in app start
type ServerConfigurationOption func(*ServerConfiguration)

// UnMarshalServerDefaults - unmarshal default parameter values
func UnMarshalServerDefaults(s string) ServerConfiguration {
	sc := ServerConfiguration{}
	err := json.Unmarshal([]byte(s), &sc)
	if err != nil {
		log.Fatal("cannot unmarshal server configuration")
	}
	return sc

}

// NewServerConfiguration - it is a factory that returns an instance of server configuration.
func NewServerConfiguration() *ServerConfiguration {
	c := UnMarshalServerDefaults(ServerDefaultJSON)
	c.Port = ":" + strings.Split(c.RunAddress, ":")[1]
	c.EnvChanged = make(map[string]bool)
	return &c

}

// NewServerConf - it is a factory that returns an instance of server configuration using Options parttern.
func NewServerConf(options ...ServerConfigurationOption) *ServerConfiguration {
	c := UnMarshalServerDefaults(ServerDefaultJSON)
	c.EnvChanged = make(map[string]bool)
	for _, option := range options {
		option(&c)
	}
	return &c
}

// UpdateSCFromEnvironment  - updates conf parameters values from os environment
func UpdateSCFromEnvironment(c *ServerConfiguration) {
	c.RunAddress = getEnv("RUN_ADDRESS", &StrValue{c.RunAddress}, c.EnvChanged).(string)
	c.AccrualSystemAddress = getEnv("ACCRUAL_SYSTEM_ADDRESS", &StrValue{c.AccrualSystemAddress}, c.EnvChanged).(string)
	//PORT is derived from ADDRESS
	c.Port = ":" + strings.Split(c.RunAddress, ":")[1]
	c.DatabaseURI = getEnv("DATABASE_URI", &StrValue{c.DatabaseURI}, c.EnvChanged).(string)
}

// UpdateSCFromFlags - updates conf parameters values from given flags in terminal string
func UpdateSCFromFlags(c *ServerConfiguration) {

	dc := NewServerConfiguration()

	var (
		a = flag.String("a", dc.RunAddress, "Domain name and :port")
		r = flag.String("r", dc.AccrualSystemAddress, "Restore from external storage:true/false")
		d = flag.String("d", dc.DatabaseURI, "database destination string")
	)
	flag.Parse()

	message := "variable %v  updated from flags, value %v"
	//Если значение из переменных равно значению по умолчанию, тогда берем из flagS
	if !c.EnvChanged["RUN_ADDRESS"] {
		c.RunAddress = *a
		c.Port = ":" + strings.Split(c.RunAddress, ":")[1]
		log.Printf(message, "RUN_ADDRESS", c.RunAddress)
		log.Printf(message, "PORT", c.Port)
	}
	if !c.EnvChanged["ACCRUAL_SYSTEM_ADDRESS"] {
		c.AccrualSystemAddress = *r
		log.Printf(message, "ACCRUAL_SYSTEM_ADDRESS", c.AccrualSystemAddress)
	}
	if !c.EnvChanged["DATABASE_URI"] {
		c.DatabaseURI = *d
		log.Printf(message, "DATABASE_URI", c.DatabaseURI)
	}
}

// VariableValue - an interface to communicate parameters values that variates by types
type VariableValue interface {
	Get() interface{} // returns parameter's value
	Set(string)       // sets parameter's value
}

// StrValue - a string implementation of VariableValue
type StrValue struct {
	value string
}

// NewStrValue - a factory that returns impl. of VariableValue for string type
func NewStrValue(s string) VariableValue {
	return &StrValue{value: s}
}

// Get - gets a string value of parameter
func (v *StrValue) Get() interface{} {
	return v.value
}

// Set - sets a string value of parameter
func (v *StrValue) Set(s string) {
	v.value = s
}

type IntValue struct {
	value int
}

// Get - gets an int value of parameter
func (v IntValue) Get() interface{} {
	return v.value
}

// Set - sets an int value of parameter
func (v *IntValue) Set(s string) {
	var err error
	v.value, err = strconv.Atoi(s)
	if err != nil {
		log.Fatal("Int Parse error")
	}
}

// NewIntValue - a factory that returns impl. of VariableValue for int type
func NewIntValue(s string) VariableValue {
	changedValue, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Int64 Parse error")
	}
	return &IntValue{value: changedValue}
}

type BoolValue struct {
	value bool
}

// Get - gets a bool value of parameter
func (v BoolValue) Get() interface{} {
	return v.value
}

// Set - sets a bool value of parameter
func (v *BoolValue) Set(s string) {
	var err error
	v.value, err = strconv.ParseBool(s)
	if err != nil {
		log.Fatal("Bool Parse error")
	}
}

// NewBoolValue - a factory that returns impl. of VariableValue for bool type
func NewBoolValue(s string) VariableValue {
	changedValue, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal("Bool Parse error")
	}
	return &BoolValue{value: changedValue}
}

func getEnv(variableName string, variableValue VariableValue, changed map[string]bool) (changedValue interface{}) {
	var stringVal string

	if variableValue == nil {
		log.Fatal("nil pointer in getEnv")
	}
	var exists bool
	stringVal, exists = os.LookupEnv(variableName)
	if !exists {
		log.Printf("variable "+variableName+" not presented in environment, remains default:%v", variableValue.Get())
		changed[variableName] = false
		return variableValue.Get()
	}
	variableValue.Set(stringVal)
	changed[variableName] = true
	log.Println("variable " + variableName + " presented in environment, value: " + stringVal)

	return variableValue.Get()
}
