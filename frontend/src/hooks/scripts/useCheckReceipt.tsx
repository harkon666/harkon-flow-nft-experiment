'use client';

import { useFlowQuery } from '@onflow/react-sdk';

// 1. Skrip Cadence yang Anda berikan
const CHECK_RECEIPT_SCRIPT = `
  import "AccessoryPack"

  // Skrip ini akan GAGAL jika dipanggil via useFlowQuery
  // karena "storage" tidak bisa diakses
  access(all) fun main(address: Address): Bool {
    let account = getAccount(address)
    
    if account.storage.type(at: AccessoryPack.ReceiptStoragePath) != nil {
      // false = resi DITEMUKAN
      return false
    }
    // true = resi TIDAK ADA
    return true
  }
`;

interface UseCheckReceiptProps {
  address: string | "";
}

// Hook kustom Anda
export function useCheckReceipt({ address }: UseCheckReceiptProps) {

  const { data, isLoading, error, refetch } = useFlowQuery({
    cadence: CHECK_RECEIPT_SCRIPT,
    args: (arg, t) => [
      arg(address, t.Address)
    ],
    query: {
      enabled: true,
    },
  });
  console.log(data, "check")
  return {
    // Kita balik logikanya agar lebih masuk akal di UI
    // 'data' = true berarti 'tidak ada resi'
    // 'data' = false berarti 'ada resi'
    hasReceipt: data === false, 
    isLoading,
    // 'error' akan berisi "member 'storage' has restricted access"
    error, 
    refetch
  };
}