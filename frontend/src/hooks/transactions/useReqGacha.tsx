'use client';

import { useFlowMutate, useFlowTransactionStatus } from '@onflow/react-sdk';

// 1. Skrip Cadence Anda (Tahap 1: Request)
const REQUEST_GACHA_TRANSACTION = `
  import "AccessoryPack"

  transaction {
    prepare(signer: auth(BorrowValue, SaveValue) &Account) {
        
        // Commit my bet and get a receipt
        let receipt <- AccessoryPack.RequestGacha()
        
        // Check that I don't already have a receipt stored
        if signer.storage.type(at: AccessoryPack.ReceiptStoragePath) != nil {
            panic("Storage collision at path=".concat(AccessoryPack.ReceiptStoragePath.toString()).concat(" a Receipt is already stored!"))
        }

        // Save that receipt to my storage
        // Note: production systems would consider handling path collisions
        signer.storage.save(<-receipt, to: AccessoryPack.ReceiptStoragePath)
    }
  }
`;

// 2. Tipe kembalian (return type) dari hook
interface UseReqGachaReturn {
  request: () => void; // Fungsi untuk memanggil transaksi
  isPending: boolean;
  txId: string | null;
  error: Error | null;
  isSealed: boolean;
  status: number | undefined;
}

// 3. Hook kustom Anda
export function useReqGacha(): UseReqGachaReturn {
  
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

  // 4. Fungsi 'request' yang akan dipanggil oleh komponen
  const request = () => {
    mutate({
      cadence: REQUEST_GACHA_TRANSACTION,
      args: () => [], // Tidak ada argumen
    });
  };

  // 5. Kembalikan (return) interface yang bersih
  return {
    request,
    isPending: isMutating,
    txId: txId || null,
    error: txError || txStatusError || null,
    isSealed: transactionStatus?.status === 4,
    status: transactionStatus?.status,
  };
}