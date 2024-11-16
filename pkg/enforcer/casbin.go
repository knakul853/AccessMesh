package enforcer

import (
	"log"

	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/knakul853/accessmesh/internal/store"
)

type Enforcer struct {
	*casbin.Enforcer
}

func NewCasbinEnforcer(store *store.MongoStore) (*Enforcer, error) {
	log.Println("Creating Casbin enforcer...")
	adapter, err := mongodbadapter.NewAdapter(store.GetURI())
	if err != nil {
		log.Printf("Error creating MongoDB adapter: %v", err)
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("model.conf", adapter)
	if err != nil {
		log.Printf("Error creating Casbin enforcer: %v", err)
		return nil, err
	}

	log.Println("Casbin enforcer created successfully.")
	return &Enforcer{enforcer}, nil
}
