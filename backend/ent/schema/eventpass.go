package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EventPass memegang skema untuk tipe EventPass (SBT).
type EventPass struct {
	ent.Schema
}

// Fields dari EventPass.
func (EventPass) Fields() []ent.Field {
	return []ent.Field{
		// Ini adalah ID on-chain dari SBT (SBT.uuid)
		field.Uint64("pass_id").
			Unique(),
		field.String("name"),
		field.String("description"),
		field.String("thumbnail"),
		field.Uint8("event_type"),
		field.Bool("is_used").
			Default(false),
	}
}

// Edges (relasi) dari EventPass.
func (EventPass) Edges() []ent.Edge {
	return []ent.Edge{
		// Relasi Many-to-One
		// Banyak Pass dimiliki oleh satu User
		edge.From("owner", User.Type).
			Ref("event_passes").
			Unique().
			Required(),

		// Relasi Many-to-One
		// Banyak Pass berasal dari satu Event
		edge.From("event", Event.Type).
			Ref("passes_issued").
			Unique().
			Required(),

		// Relasi One-to-One
		// Satu Pass digunakan untuk me-mint satu Moment
		edge.To("moment", NFTMoment.Type).
			Unique(),
	}
}
