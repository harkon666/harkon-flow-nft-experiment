package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NFTAccessory holds the schema definition for the NFTAccessory entity.
type NFTAccessory struct {
	ent.Schema
}

// Fields of the NFTAccessory.
func (NFTAccessory) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("nft_id").
			Unique(), // ID NFT on-chain
		field.String("name"),
		field.String("description"),
		field.String("thumbnail"),
		field.String("equipment_type"),
	}
}

// Edges of the NFTAccessory.
func (NFTAccessory) Edges() []ent.Edge {
	return []ent.Edge{
		// Mendefinisikan relasi "many-to-one" (banyak-ke-satu)
		edge.From("owner", User.Type).
			Ref("accessories"). // Merujuk ke edge "accessories" di skema User
			Unique().           // Aksesori hanya punya satu owner
			Required(),         // Aksesori wajib punya owner

		edge.From("equipped_on_moment", NFTMoment.Type).
			Ref("equipped_accessories"). // Merujuk ke edge di NFTMoment
			Unique(),

		edge.From("listing", Listing.Type).
			Ref("nft_accessory").
			Unique(),
	}
}
