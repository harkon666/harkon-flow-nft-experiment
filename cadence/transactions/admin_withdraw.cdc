// transactions/admin_withdraw.cdc
// Transaksi ini HANYA bisa dijalankan oleh Akun Admin

import "FungibleToken"
import "FlowToken"
import "AccessoryPack"

// Tentukan jumlah yang ingin Anda tarik
transaction(amount: UFix64) {

    let adminRef: &AccessoryPack.Admin
    let vaultRef: &FlowToken.Vault

    prepare(signer: auth(BorrowValue) &Account) {
        
        // 1. Pinjam "kunci" Admin dari storage pribadi Anda
        self.adminRef = signer.storage.borrow<&AccessoryPack.Admin>(
            from: AccessoryPack.AdminStoragePath
        ) ?? panic("Tidak bisa meminjam resource Admin")

        // 2. Pinjam "brankas" pribadi Anda (tempat menyimpan hasil)
        self.vaultRef = signer.storage.borrow<&FlowToken.Vault>(
            from: /storage/flowTokenVault // Path standar FlowToken
        ) ?? panic("Tidak bisa meminjam brankas pribadi FlowToken")
    }

    execute {
        // 3. Panggil fungsi withdraw, yang mengembalikan 'Vault'
        let withdrawnVault <- self.adminRef.adminWithdraw(amount: amount)
        
        // 4. Setor (deposit) dana ke brankas pribadi Anda
        self.vaultRef.deposit(from: <-withdrawnVault)

        log("Berhasil menarik ".concat(amount.toString()).concat(" FLOW"))
    }
}