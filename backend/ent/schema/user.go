package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

type HighlightedEventPassIds struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("address").Unique(),
		field.String("nickname").Optional(),
		field.String("bio").Optional(),
		field.String("pfp").Optional(),
		field.String("short_description").Optional(),
		field.String("bg_image").Optional(),
		field.JSON("highlighted_eventPass_ids", []uint64{}).Optional(),
		field.Uint64("highlighted_moment_id").Optional(),
		field.JSON("socials", map[string]string{}).Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("event_passes", EventPass.Type),
		edge.To("hosted_events", Event.Type),
		edge.To("moments", NFTMoment.Type),
		edge.To("accessories", NFTAccessory.Type),
		edge.To("attendances", Attendance.Type),
		edge.To("listings", Listing.Type),
	}
}
