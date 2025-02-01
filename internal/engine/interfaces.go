package engine

import "github.com/Blocktunium/gonyx/pkg/http"

type RestfulApp interface {
	Routes() []http.HttpRoute
	GetName() string
}

type GrpcApp interface {
	GetName() string
	GetServerNames() []string
}
