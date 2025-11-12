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
    nftMomentId: UInt64,
  ) {
      let momentCollectionRef: auth(NFTMoment.Equip) &NFTMoment.Collection
      let accessoryCollectionRef: &NFTAccessory.Collection
      let frameNFT: @NFTAccessory.NFT
      prepare(signer: auth(BorrowValue) &Account) {

        let accessoryCollectionData = NFTAccessory.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The NFTAccessory contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")
        let momentCollectionData = NFTMoment.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view. The NFTMoment contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")
        self.momentCollectionRef = signer.storage.borrow<auth(NFTMoment.Equip) &NFTMoment.Collection>(from: momentCollectionData.storagePath)
            ?? panic("No Moment Collection Ressource in Storage")
        self.accessoryCollectionRef = signer.storage.borrow<&NFTAccessory.Collection>(from: accessoryCollectionData.storagePath)
            ?? panic("No Accessory Collection Ressource in Storage")
        // borrow a reference to the signer's NFT collection
        let withdrawRef = signer.storage.borrow<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(
          from: accessoryCollectionData.storagePath
        ) ?? panic("The signer does not store a NFT Collection object at the path \(accessoryCollectionData.storagePath)"
                        .concat("The signer must initialize their account with this collection first!"))
        self.frameNFT <- withdrawRef.withdraw(withdrawID: nftAccessoryId) as! @NFTAccessory.NFT

        assert(
          self.frameNFT.getType().identifier == "A.f8d6e0586b0a20c7.NFTAccessory.NFT",
          message: "The NFT that was withdrawn to transfer is not the type that was requested <A.f8d6e0586b0a20c7.NFTAccessory.NFT>."
        )
      }

      execute {
        let prevEquipFrame <- self.momentCollectionRef.equipFrame(momentNFTID: nftMomentId, frameNFT: <-self.frameNFT)
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