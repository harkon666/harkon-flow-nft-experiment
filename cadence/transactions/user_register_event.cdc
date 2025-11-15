import "EventManager"

// Transaksi ini dijalankan oleh PENGGUNA
// untuk mendaftarkan diri mereka ke sebuah event

transaction(eventID: UInt64) {

    prepare(signer: &Account) {
        
        // Panggil fungsi publik di kontrak EventManager
        // Kita meneruskan 'signer' (akun pengguna)
        EventManager.userRegisterForEvent(
            eventID: eventID,
            userAccount: signer
        )
    }

    execute {
        log("Pengguna berhasil mendaftar untuk event ID: ".concat(eventID.toString()))
    }
}