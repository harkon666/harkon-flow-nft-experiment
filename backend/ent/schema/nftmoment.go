package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NFTMoment holds the schema definition for the NFTMoment entity.
type NFTMoment struct {
	ent.Schema
}

// Fields of the NFTMoment.
func (NFTMoment) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("nft_id").
			Unique(), // ID NFT on-chain
		field.String("name"),
		field.String("description"),
		field.String("thumbnail"),
	}
}

// Edges of the NFTMoment.
func (NFTMoment) Edges() []ent.Edge {
	return []ent.Edge{
		// Mendefinisikan relasi "many-to-one" (banyak-ke-satu)
		edge.From("owner", User.Type).
			Ref("moments"). // Merujuk ke edge "moments" di skema User
			Unique().       // Momen hanya punya satu owner
			Required(),     // Momen wajib punya owner

		edge.To("equipped_accessories", NFTAccessory.Type),

		edge.From("minted_with_pass", EventPass.Type).
			Ref("moment").
			Unique(),
	}
}
