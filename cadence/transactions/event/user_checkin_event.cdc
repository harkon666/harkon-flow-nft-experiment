import "EventManager"
import "EventPass"

// Transaksi ini dijalankan oleh ADMIN/BACKEND
// untuk secara manual melakukan check-in atas nama pengguna

transaction(
    eventID: UInt64,
    userAddress: Address // Alamat pengguna yang di-scan
) {

    // Referensi ke 'Admin' resource
    let adminRef: &EventManager.Admin
    let recipient: &EventPass.Collection
    let eventRef: &EventManager.Event
    let eventPassMinterRef: &EventPass.NFTMinter
    prepare(signer: auth(BorrowValue) &Account) {
        
        // Pinjam 'kunci' Admin dari storage 'signer' (backend)
        self.eventRef = EventManager.events[eventID] as! &EventManager.Event
        self.recipient = getAccount(userAddress).capabilities.borrow<&EventPass.Collection>(
          EventPass.CollectionPublicPath
        ) ?? panic("cant borrow ressource recipient collection EventPass")
        self.adminRef = signer.storage.borrow<&EventManager.Admin>(
            from: EventManager.eventManagerStoragePath
        ) ?? panic("cant borrow resource Admin EventManager")
        self.eventPassMinterRef = signer.storage.borrow<&EventPass.NFTMinter>(
          from: EventPass.MinterStoragePath
        ) ?? panic("cant borrow ressource minter EventPass")
    }

    execute {
        let thumbnail = self.eventRef.eventPassImg != nil ?
          self.eventRef.eventPassImg?.uri() :
          "https://white-lazy-marten-351.mypinata.cloud/ipfs/bafybeibv7mz4yvpuw5ejbovka3h2zhrzyf7jptikz7fzsuprlgw3h6qtnq"

        self.adminRef.checkInUserToEvent(eventID: eventID, userAddress: userAddress)
        self.eventPassMinterRef.mintNFT(
          recipient: self.recipient,
          name: self.eventRef.eventName,
          description: self.eventRef.description,
          thumbnail: thumbnail!,
          eventType: self.eventRef.eventType.rawValue,
          eventID: self.eventRef.eventID
        )
        
        log("checkIn success and event pass minted to ".concat(userAddress.toString()))
    }
}