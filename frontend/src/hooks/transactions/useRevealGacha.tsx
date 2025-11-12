'use client';

import { useFlowMutate, useFlowTransactionStatus } from '@onflow/react-sdk';

// 1. Skrip Cadence Anda (Tahap 2: Reveal)
//    (Ini adalah kode Anda, sedikit dimodifikasi agar lebih aman
//     dengan mengimpor alamat kontrak dari FCL config)
const REVEAL_GACHA_TRANSACTION = `
  import "AccessoryPack"
  import "NFTAccessory"
  import "NonFungibleToken"
  import "MetadataViews"
  /// Retrieves the saved Receipt and redeems it to reveal the gacha result, depositing winnings with any luck
  ///
  transaction(recipient: Address) {

      let recipientCollectionRef: &{NonFungibleToken.Receiver}
      let minterRef: &NFTAccessory.NFTMinter
      let randomNumber: UInt8

      prepare(signer: auth(BorrowValue, LoadValue) &Account) {
          let collectionData = NFTAccessory.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
              ?? panic("Could not resolve NFTCollectionData view. The NFTAccessory contract needs to implement the NFTCollectionData Metadata view in order to execute this transaction")
          // Load my receipt from storage
          let receipt <- signer.storage.load<@AccessoryPack.Receipt>(from: AccessoryPack.ReceiptStoragePath)
              ?? panic("No Receipt found in storage at path=".concat(AccessoryPack.ReceiptStoragePath.toString()))
          self.minterRef = signer.storage.borrow<&NFTAccessory.NFTMinter>(from: NFTAccessory.MinterStoragePath)
            ?? panic("no admin")
          self.randomNumber = AccessoryPack.RevealGacha(receipt: <-receipt)
          
          self.recipientCollectionRef = getAccount(recipient).capabilities.borrow<&{NonFungibleToken.Receiver}>(collectionData.publicPath)
              ?? panic("The recipient does not have a NonFungibleToken Receiver at "
                      .concat(collectionData.publicPath.toString())
                      .concat(" that is capable of receiving an NFT.")
                      .concat("The recipient must initialize their account with this collection and receiver first!"))
      }

      execute {
          AccessoryPack.distributeAccessory(self.randomNumber, recipient: self.recipientCollectionRef)
          log("berhasil njay")
      }
  }
`;

// 2. Tipe untuk argumen
interface UseRevealGachaOptions {
  recipient: string;
}

// 3. Tipe kembalian (return type)
interface UseRevealGachaReturn {
  reveal: (options: UseRevealGachaOptions) => void;
  isPending: boolean;
  txId: string | null;
  error: Error | null;
  isSealed: boolean;
  status: number | undefined;
}

// 4. Hook kustom Anda
export function useRevealGacha(): UseRevealGachaReturn {
  
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

  // 5. Fungsi 'reveal' yang akan dipanggil oleh komponen
  const reveal = (options: UseRevealGachaOptions) => {
    const { recipient } = options;

    mutate({
      cadence: REVEAL_GACHA_TRANSACTION,
      args: (arg, t) => [
        arg(recipient, t.Address)
      ],
    });
  };

  // 6. Kembalikan (return) interface yang bersih
  return {
    reveal,
    isPending: isMutating,
    txId: txId || null,
    error: txError || txStatusError || null,
    isSealed: transactionStatus?.status === 4,
    status: transactionStatus?.status,
  };
}