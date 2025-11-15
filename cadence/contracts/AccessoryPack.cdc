import "Burner"
import "NonFungibleToken"
import "FungibleToken"
import "FlowToken"
import "RandomConsumer"
import "NFTAccessory"

/// CoinToss is a simple game contract showcasing the safe use of onchain randomness by way of a commit-reveal sheme.
///
/// See FLIP 123 for more details: https://github.com/onflow/flips/blob/main/protocol/20230728-commit-reveal.md
/// And the onflow/random-coin-toss repo for implementation context: https://github.com/onflow/random-coin-toss
///
/// NOTE: This contract is for demonstration purposes only and is not intended to be used in a production environment.
///
access(all) contract AccessoryPack {
    /// The RandomConsumer.Consumer resource used to request & fulfill randomness
    access(self) let consumer: @RandomConsumer.Consumer
    access(all) let gachaPrice: UFix64
    access(self) let vault: @FlowToken.Vault

    /// The canonical path for common Receipt storage
    /// Note: production systems would consider handling path collisions
    access(all) let ReceiptStoragePath: StoragePath
    access(all) let AdminStoragePath: StoragePath
    /* --- Events --- */
    //
    access(all) event AccessoryPackOpened(commitBlock: UInt64, receiptID: UInt64)
    access(all) event AccessoryPackRevealed(rarity: UInt8, commitBlock: UInt64, receiptID: UInt64)
    access(all) event AccessoryDistributed(recipient: Address, id: UInt64, name: String, description: String, thumbnail: String, equipmentType: String)

    /// The Receipt resource is used to store the bet amount and the associated randomness request. By listing the
    /// RandomConsumer.RequestWrapper conformance, this resource inherits all the default implementations of the
    /// interface. This is why the Receipt resource has access to the getRequestBlock() and popRequest() functions
    /// without explicitly defining them.
    ///
    access(all) resource Receipt : RandomConsumer.RequestWrapper {
        /// The amount bet by the user
        /// The associated randomness request which contains the block height at which the request was made
        /// and whether the request has been fulfilled.
        access(all) var request: @RandomConsumer.Request?

        init(request: @RandomConsumer.Request) {
            self.request <- request
        }
    }

    /* --- Commit --- */
    //
    /// In this method, the caller commits a bet. The contract takes note of the block height and bet amount, returning a
    /// Receipt resource which is used by the better to reveal the coin toss result and determine their winnings.
    ///
    access(all) fun RequestGacha(payment: @{FungibleToken.Vault}): @Receipt {
        pre {
          payment.balance == self.gachaPrice: "imbalance payment Amount"
        }
        let request: @RandomConsumer.Request <- self.consumer.requestRandomness()
        let receipt <- create Receipt(request: <-request)
        self.vault.deposit(from: <-payment)
        emit AccessoryPackOpened(commitBlock: receipt.getRequestBlock()!, receiptID: receipt.uuid)

        return <- receipt
    }

    /* --- Reveal --- */
    access(all) fun RevealGacha(receipt: @Receipt): UInt8 {
        pre {
            receipt.request != nil: 
            "CoinToss.revealCoin: Cannot reveal the coin! The provided receipt has already been revealed."
            receipt.getRequestBlock()! <= getCurrentBlock().height:
            "CoinToss.revealCoin: Cannot reveal the coin! The provided receipt was committed for block height ".concat(receipt.getRequestBlock()!.toString())
            .concat(" which is greater than the current block height of ")
            .concat(getCurrentBlock().height.toString())
            .concat(". The reveal can only happen after the committed block has passed.")
        }
        let commitBlock = receipt.getRequestBlock()!
        let receiptID = receipt.uuid

        let randomNumber = self._randomNumber(request: <-receipt.popRequest())
        Burner.burn(<-receipt)
        
        emit AccessoryPackRevealed(rarity: randomNumber, commitBlock: commitBlock, receiptID: receiptID)
        return randomNumber
    }

    //gacha accessory send to receiver
    access(all) fun distributeAccessory(_ randomNumber:UInt8, recipient: &{NonFungibleToken.Receiver}) {
      let minterRef = self.account.storage.borrow<&NFTAccessory.NFTMinter>(from: NFTAccessory.MinterStoragePath)
        ?? panic("No Minter resource in storage")
      var name: String = ""
      var description: String = ""
      var descriptionRarity: String = ""
      var thumbnail: String = ""
      var equipmentType: String = ""
      if randomNumber == 1 {
          name = "Bingkai Emas"
          description = "emas banget"
          descriptionRarity = "Super Rare"
          thumbnail = "ipfs://bafybeibknydg67jcaybwkpxvstlzkd5i4itzgxfqmdn2fs62dldkdt4d3u"
          equipmentType = "Bingkai"
      } else if randomNumber > 1 && randomNumber < 11 {
          name = "Bingkai Perak"
          description = "Perak banget"
          descriptionRarity = "Rare"
          thumbnail = "ipfs://bafybeiajk2fbvg2ycxbbqugwbu52luav523rv43lwmkjr2mchbnlpbvkdy"
          equipmentType = "Bingkai"
      } else {
          name = "Bingkai Kayu"
          description = "Kayu banget"
          descriptionRarity = "Common"
          thumbnail = "ipfs://bafybeifxyn7qqz72io7eegzarsv7omsxwxhdtwprzck4vd5ktfnwyf3oa4"
          equipmentType = "Bingkai"
      }

      let mintedNFT <- minterRef.mintNFT(
          name: name,
          description: description,
          thumbnail: thumbnail,
          equipmentType: equipmentType,
          score: UFix64(randomNumber),
          max: 100.0,
          descriptionRarity: descriptionRarity
      )

      emit AccessoryPack.AccessoryDistributed(
          recipient: recipient.owner!.address,
          id: mintedNFT.id,
          name: mintedNFT.name,
          description: mintedNFT.description,
          thumbnail: mintedNFT.thumbnail,
          equipmentType: mintedNFT.equipmentType
      )
      let textResult = mintedNFT.name.concat(" ").concat(descriptionRarity)
      //what should i do? do i need to make it the same as how NFTMoment minted nft?
      recipient.deposit(token: <-mintedNFT)
    }

    access(all) resource Admin  {
      access(all) fun adminWithdraw(amount: UFix64): @{FungibleToken.Vault} {     
            // Tarik (withdraw) dana dan kembalikan
            return <- AccessoryPack.vault.withdraw(amount: amount)
      }
    }

    /// Returns a random number between 0 and 1 using the RandomConsumer.Consumer resource contained in the contract.
    /// For the purposes of this contract, a simple modulo operation could have been used though this is not the case
    /// for all ranges. Using the Consumer.fulfillRandomInRange function ensures that we can get a random number
    /// within any range without a risk of bias.
    ///
    access(self) fun _randomNumber(request: @RandomConsumer.Request): UInt8 {
        return UInt8(self.consumer.fulfillRandomInRange(request: <-request, min: 1, max: 100))
    }

    init() {
        // Create a RandomConsumer.Consumer resource
        self.consumer <-RandomConsumer.createConsumer()
        self.gachaPrice = 1.0
        self.vault <- FlowToken.createEmptyVault(vaultType: Type<@FlowToken.Vault>())

        self.AdminStoragePath = /storage/AccessoryPackAdmin

        self.account.storage.save(
            <- create Admin(),
            to: self.AdminStoragePath
        )

        // Set the ReceiptStoragePath to a unique path for this contract - appending the address to the identifier
        // prevents storage collisions with other objects in user's storage
        self.ReceiptStoragePath = StoragePath(identifier: "AccessoryPack_".concat(self.account.address.toString()))!
    }
}