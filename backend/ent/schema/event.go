package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Event memegang skema untuk tipe Event.
type Event struct {
	ent.Schema
}

// Fields dari Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		// Ini adalah ID dari kontrak EventManager (nextEventID)
		field.Uint64("event_id").
			Unique(),
		field.String("name"),
		field.String("description"),
		field.String("thumbnail"),
		field.Uint8("event_type"),
		field.String("location"),
		field.Float("lat"),
		field.Float("long"),
		field.Time("start_date"),
		field.Time("end_date"),
		field.Uint64("quota"),
	}
}

// Edges (relasi) dari Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		// Relasi Many-to-One (Banyak-ke-Satu)
		// Banyak Event dimiliki oleh satu Host (User)
		edge.From("host", User.Type).
			Ref("hosted_events"). // Merujuk ke edge di skema User
			Unique().
			Required(), // Setiap event harus punya host

		// Relasi One-to-Many (Satu-ke-Banyak)
		// Satu Event akan menghasilkan banyak EventPass
		edge.To("passes_issued", EventPass.Type),
		edge.To("attendances", Attendance.Type),
	}
}
