import "NFTMoment"
import "NFTAccessory"
import "MetadataViews"
import "ViewResolver"

// Skrip ini hanya mengambil view kustom Anda
access(all) fun main(address: Address, id: UInt64): &NFTAccessory.NFT? {
    
    let account = getAccount(address)
    let collection = account.capabilities.borrow<&{ViewResolver.ResolverCollection}>(
      NFTMoment.CollectionPublicPath
    )
        ?? panic("Tidak bisa meminjam koleksi")

    let resolver = collection.borrowViewResolver(id: id) 
        ?? panic("Tidak bisa meminjam resolver")
    
    // 1. Minta 'view' kustom Anda secara spesifik
    let view = resolver.resolveView(Type<NFTMoment.NFTMomentEquipmentMetadataView>())

    // 2. 'Cast' (ubah) hasilnya ke tipe struct Anda
    if let equipmentView = view as! NFTMoment.NFTMomentEquipmentMetadataView? {
        
        // 3. Kembalikan datanya!
        return equipmentView.equippedFrame
    }
    
    return nil
}