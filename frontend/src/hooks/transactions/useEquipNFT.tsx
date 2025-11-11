'use client';

import { useFlowMutate, useFlowTransactionStatus } from '@onflow/react-sdk';

// 1. Definisikan string Cadence Anda
//    Ini adalah kode transaksi yang Anda berikan
const EQUIP_NFT_TRANSACTION = `
  import "NonFungibleToken"
  import "NFTMoment"
  import "NFTAccessory"
  import "MetadataViews"

  transaction(
    nftAccessoryId: UInt64,
    nftMomentId: UInt64
  ) {
    let momentNFT: &NFTMoment.NFT
    let accessoryCollectionRef: &NFTAccessory.Collection
    let frameNFT: @NFTAccessory.NFT
    
    prepare(signer: auth(BorrowValue) &Account) {

      let accessoryCollectionData = NFTAccessory.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
          ?? panic("Could not resolve NFTCollectionData view...")
      let momentCollectionData = NFTMoment.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
          ?? panic("Could not resolve NFTMoment CollectionData view...")
      
      let momentCollectionRef = signer.storage.borrow<&NFTMoment.Collection>(from: momentCollectionData.storagePath)
          ?? panic("No Moment Collection Ressource in Storage")
      
      self.accessoryCollectionRef = signer.storage.borrow<&NFTAccessory.Collection>(from: accessoryCollectionData.storagePath)
          ?? panic("No Accessory Collection Ressource in Storage")
      
      self.momentNFT = momentCollectionRef.borrowNFT(nftMomentId) as! &NFTMoment.NFT
      
      // Pinjam referensi untuk 'withdraw' (sesuai kode Anda)
      let withdrawRef = signer.storage.borrow<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(
          from: accessoryCollectionData.storagePath
      ) ?? panic("Could not borrow withdraw reference from signer")

      self.frameNFT <- withdrawRef.withdraw(withdrawID: nftAccessoryId) as! @NFTAccessory.NFT

      // Asumsi 'A.f8d6e0586b0a20c7' adalah alamat Anda di emulator
      assert(
        self.frameNFT.getType().identifier == "A.f8d6e0586b0a20c7.NFTAccessory.NFT",
        message: "Withdrawn NFT is not the correct Accessory type."
      )
    }

    execute {
      let prevEquipFrame <- self.momentNFT.equipFrame(frameNFT: <-self.frameNFT)
      if prevEquipFrame != nil {
        self.accessoryCollectionRef.deposit(token: <-prevEquipFrame!)
      } else {
        destroy prevEquipFrame
      }
    }
  }
`;

// 2. Definisikan tipe untuk argumen yang diterima fungsi 'equip'
interface UseEquipNFTOptions {
  nftAccessoryId: number; // JavaScript 'number'
  nftMomentId: number;    // JavaScript 'number'
}

// 3. Definisikan tipe untuk nilai yang dikembalikan oleh hook
interface UseEquipNFTReturn {
  equip: (options: UseEquipNFTOptions) => void; // Fungsi untuk memanggil transaksi
  isPending: boolean;
  txId: string | null;
  error: Error | null;
  isSealed: Boolean;
  status: number | undefined;
}

// 4. Buat hook kustom Anda
export function useEquipNFT(): UseEquipNFTReturn {
  
  const { 
    mutate, 
    isPending: isMutating, 
    data: txId, 
    error: txError 
  } = useFlowMutate();
  const { 
    transactionStatus,  
    error: txStatusError 
  } = useFlowTransactionStatus({
    id: txId || '',
  });
  
  // 5. Buat fungsi 'equip' yang akan dipanggil oleh komponen
  const equip = (options: UseEquipNFTOptions) => {
    const { nftAccessoryId, nftMomentId } = options;
    console.log(nftAccessoryId, nftMomentId)
    mutate({
      cadence: EQUIP_NFT_TRANSACTION,
      
      args: (arg, t) => [
        // Konversi JS 'number' ke Cadence 'UInt64' (sebagai string)
        arg(Number(nftAccessoryId), t.UInt64),
        arg(Number(nftMomentId), t.UInt64),
      ],
    });
  };

  // 6. Kembalikan (return) interface yang bersih
  return {
    equip,
    isPending: isMutating,
    txId: txId || null,
    error: txError || txStatusError || null,
    isSealed: transactionStatus?.status === 4,
    status: transactionStatus?.status,
  };
}