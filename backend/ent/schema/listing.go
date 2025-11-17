package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Listing memegang skema untuk 'Listing' (Penjualan) di marketplace.
type Listing struct {
	ent.Schema
}

// Fields dari Listing.
func (Listing) Fields() []ent.Field {
	return []ent.Field{
		// Ini adalah 'listingResourceID' dari event. HARUS unik.
		field.Uint64("listing_id").
			Unique(),

		// Harga jual dalam UFix64 (disimpan sebagai float)
		field.Float("price"),

		// Tipe vault pembayaran (misal: "A.0ae...FlowToken.Vault")
		field.String("payment_vault_type"),

		// ID kustom (opsional)
		field.String("custom_id").
			Optional().
			Nillable(), // Izinkan 'nil' di database

		// Waktu kedaluwarsa (disimpan sebagai time.Time)
		field.Time("expiry"),

		// (Kita tidak perlu menyimpan 'nftID' atau 'sellerAddress' di sini
		// karena itu akan ditangani oleh 'Edges' (Relasi) di bawah)
	}
}

// Edges (relasi) dari Listing.
func (Listing) Edges() []ent.Edge {
	return []ent.Edge{
		// Relasi Many-to-One (Banyak Listing dimiliki oleh 1 User)
		edge.From("seller", User.Type).
			Ref("listings"). // Merujuk ke edge 'listings' di User
			Unique().        // Satu listing hanya punya satu penjual
			Required(),      // Wajib ada penjual

		// Relasi One-to-One (Satu Listing untuk 1 NFT)
		// (Saya asumsikan Anda ingin menjual 'NFTAccessory' di sini)
		edge.To("nft_accessory", NFTAccessory.Type).
			Unique().   // Satu listing hanya untuk satu aksesori
			Required(), // Wajib ada NFT yang dijual

		// (Jika Anda juga ingin menjual 'NFTMoment', tambahkan edge
		//  opsional lain ke 'NFTMoment' di sini)
	}
}
