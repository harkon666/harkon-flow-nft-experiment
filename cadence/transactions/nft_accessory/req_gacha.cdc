/// This transaction is what an account would run
/// to set itself up to receive NFTs

import "AccessoryPack"
import "NonFungibleToken"
import "FungibleToken"
import "FlowToken"
import "FungibleTokenMetadataViews"

transaction {

    prepare(signer: auth(BorrowValue, SaveValue) &Account) {
        
        // Commit my bet and get a receipt
        let vaultData = FlowToken.resolveContractView(resourceType: nil, viewType: Type<FungibleTokenMetadataViews.FTVaultData>()) as! FungibleTokenMetadataViews.FTVaultData?
            ?? panic("Could not resolve FTVaultData view. The ExampleToken"
                .concat(" contract needs to implement the FTVaultData Metadata view in order to execute this transaction."))

        // Get a reference to the signer's stored vault
        let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(from: vaultData.storagePath)
            ?? panic("The signer does not store an ExampleToken.Vault object at the path "
                    .concat(vaultData.storagePath.toString())
                    .concat(". The signer must initialize their account with this vault first!"))
        let sentVault <- vaultRef.withdraw(amount: AccessoryPack.gachaPrice)
        let receipt <- AccessoryPack.RequestGacha(payment: <-sentVault)

        // Check that I don't already have a receipt stored
        if signer.storage.type(at: AccessoryPack.ReceiptStoragePath) != nil {
            panic("Storage collision at path=".concat(AccessoryPack.ReceiptStoragePath.toString()).concat(" a Receipt is already stored!"))
        }

        // Save that receipt to my storage
        // Note: production systems would consider handling path collisions
        signer.storage.save(<-receipt, to: AccessoryPack.ReceiptStoragePath)
    }
}