package graph
//go:generate go run github.com/99designs/gqlgen
import "github.com/brandonforster/resolver/graph/model"

// TODO: have this backed by a SQLite DB and not an array
type Resolver struct{
	IPs []*model.IPAddress
}
