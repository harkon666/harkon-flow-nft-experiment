/// This transaction is what an account would run
/// to set itself up to receive NFTs

import "AccessoryPack"

transaction {

    prepare(signer: auth(BorrowValue, SaveValue) &Account) {
        
        // Commit my bet and get a receipt
        let receipt <- AccessoryPack.RequestGacha()
        
        // Check that I don't already have a receipt stored
        if signer.storage.type(at: AccessoryPack.ReceiptStoragePath) != nil {
            panic("Storage collision at path=".concat(AccessoryPack.ReceiptStoragePath.toString()).concat(" a Receipt is already stored!"))
        }

        // Save that receipt to my storage
        // Note: production systems would consider handling path collisions
        signer.storage.save(<-receipt, to: AccessoryPack.ReceiptStoragePath)
    }
}