package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	    return []ent.Field{
        field.String("address").Unique(),
    }
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// Mendefinisikan relasi "one-to-many" (satu-ke-banyak)
		// Satu User bisa memiliki banyak NFTAccessories
		edge.To("accessories", NFTAccessory.Type),

		// Satu User juga bisa memiliki banyak NFTMoments
		edge.To("moments", NFTMoment.Type),
	}
}
