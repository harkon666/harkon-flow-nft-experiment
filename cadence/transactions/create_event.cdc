import "EventManager"
import "MetadataViews"

// Transaksi ini HANYA bisa dijalankan oleh Admin/Backend
// untuk membuat event baru di 'EventManager'

transaction(
    hostAddress: Address,
    eventName: String,
    description: String,
    thumbnailURL: String,
    eventType: UInt8, // 0 untuk online, 1 untuk offline
    location: String,
    lat: Fix64,       // Kirim 0.0 jika event online
    long: Fix64,      // Kirim 0.0 jika event online
    startDate: UFix64,
    endDate: UFix64,
    quota: UInt64
) {

    // Referensi ke 'Admin' resource
    let adminRef: &EventManager.Admin

    prepare(signer: auth(BorrowValue) &Account) {
        
        // Pinjam 'kunci' Admin dari storage 'signer'
        // Ini memastikan hanya admin yang bisa menjalankan ini
        self.adminRef = signer.storage.borrow<&EventManager.Admin>(
            from: /storage/EventManagerAdmin
        ) ?? panic("Tidak bisa meminjam resource Admin EventManager")
    }

    execute {
        // Panggil fungsi 'createEvent' di dalam resource Admin
        let newID = self.adminRef.createEvent(
            hostAddress: hostAddress,
            eventName: eventName,
            description: description,
            thumbnailURL: thumbnailURL,
            eventType: eventType,
            location: location,
            lat: lat,
            long: long,
            startDate: startDate,
            endDate: endDate,
            quota: quota
        )
        
        log("Event baru berhasil dibuat dengan ID: ".concat(newID.toString()))
    }
}