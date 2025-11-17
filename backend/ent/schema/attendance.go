package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Attendance memegang skema untuk 'join table' (tabel gabungan)
// yang menghubungkan Users dan Events.
type Attendance struct {
	ent.Schema
}

// Fields dari Attendance.
func (Attendance) Fields() []ent.Field {
	return []ent.Field{
		// Kita bisa simpan status dari kontrak Anda
		// false = registered, true = checked_in
		field.Bool("checked_in").
			Default(false),

		// Kita simpan kapan mereka mendaftar
		field.Time("registration_time").
			Default(time.Now),
	}
}

// Edges (relasi) dari Attendance.
// Ini adalah inti dari Many-to-Many
func (Attendance) Edges() []ent.Edge {
	return []ent.Edge{
		// Relasi ke User
		// Satu 'Attendance' dimiliki oleh satu User
		edge.From("user", User.Type).
			Ref("attendances"). // Merujuk ke edge di skema User
			Unique().
			Required(),

		// Relasi ke Event
		// Satu 'Attendance' milik satu Event
		edge.From("event", Event.Type).
			Ref("attendances"). // Merujuk ke edge di skema Event
			Unique().
			Required(),
	}
}
