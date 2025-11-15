import "EventManager"

// Transaksi ini dijalankan oleh ADMIN/BACKEND
// untuk secara manual melakukan check-in atas nama pengguna

transaction(
    eventID: UInt64,
    userAddress: Address // Alamat pengguna yang di-scan
) {

    // Referensi ke 'Admin' resource
    let adminRef: &EventManager.Admin
    let eventRef: &EventManager.Event
    prepare(signer: auth(BorrowValue) &Account) {
        
        // Pinjam 'kunci' Admin dari storage 'signer' (backend)
        self.eventRef = EventManager.events[eventID] as! &EventManager.Event
        self.adminRef = signer.storage.borrow<&EventManager.Admin>(
            from: /storage/EventManagerAdmin
        ) ?? panic("Tidak bisa meminjam resource Admin EventManager")
    }

    execute {
        // Panggil fungsi 'checkInUserToEvent' di dalam resource Admin
        self.adminRef.checkInUserToEvent(eventID: eventID, userAddress: userAddress)
        
        log("Admin berhasil melakukan check-in untuk ".concat(userAddress.toString()))
    }
}