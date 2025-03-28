package mongokit

import (
	"reflect"
	"testing"
)

func TestMongoWrapper_Initialize(t *testing.T) {
	makeReadyConfigManager()

	wrapper := &MongoWrapper{name: "db/server4", config: &Mongo{
		DatabaseName: "m",
		Username:     "u1",
		Password:     "p1",
		Host:         "127.0.0.1",
		Port:         "27017",
		Options: map[string]string{
			"maxpoolsize":      "100",
			"w":                "majority",
			"connecttimeoutms": "30000",
		},
	}}
	newWrapper, err := NewMongoWrapper("db/server4")

	if err != nil {
		t.Errorf("Creating Mongo Wrapper --> Expected: %v, but got %v", nil, err)
		return
	}

	if !reflect.DeepEqual(wrapper.config, newWrapper.config) {
		t.Errorf("Creating Mongo Wrapper --> Expected: %v, but got %v", wrapper.config, newWrapper.config)
		return
	}
}
