import { useFlowMutate, useFlowTransactionStatus } from "@onflow/react-sdk"; // <-- Ganti ini
import * as t from "@onflow/types";
import { arg } from "@onflow/fcl";
import { useEffect } from "react";

// 1. Definisikan tipe untuk argumen yang diterima fungsi 'mint'
interface UseMintNFTOptions {
  recipient: string;
  name: string;
  description: string;
  thumbnail: string;
}

// 2. Definisikan tipe untuk nilai yang dikembalikan oleh hook
interface UseMintNFTReturn {
  mint: (options: UseMintNFTOptions) => void;
  txId: string | null;
  isPending: boolean;
  error: Error | null;
  status: number | undefined; // (e.g., 4 = sealed)
  isSealed: boolean;
}

// (Cadence string tetap sama)
const MINT_NFT_TRANSACTION = `
  import "NonFungibleToken"
  import "NFTMoment"
  import "MetadataViews"

  transaction(
      recipient: Address,
      name: String,
      description: String,
      thumbnail: String
  ) { /* ... isi transaksi Anda ... */ }
`;

// 3. Tentukan tipe kembalian (return type) dari hook
export function useMintNFT(): UseMintNFTReturn {
  
  // Asumsikan hook ini mengembalikan tipe yang sesuai
  const { 
    mutate, 
    isPending: isMutating, 
    data: txId, 
    error: txError 
  } = useFlowMutate();

  const { transactionStatus, error: txStatusError } = useFlowTransactionStatus({
    id: txId || '',
  });

  // 4. Beri tipe pada parameter 'options'
  const mint = (options: UseMintNFTOptions) => {
    const { recipient, name, description, thumbnail } = options;

    if (!recipient || !name || !description || !thumbnail) {
      console.error("Argumen untuk minting tidak lengkap");
      return;
    }

    mutate({
      cadence: MINT_NFT_TRANSACTION,
      args: (arg, t) => [
        arg(recipient, t.Address),
        arg(name, t.String),
        arg(description, t.String),
        arg(thumbnail, t.String),
      ],
    });
  };

  // 5. Kembalikan objek yang sesuai dengan interface UseMintNFTReturn
  return {
    mint,
    txId: txId || null,
    isPending: isMutating || isStatusPending,
    error: txError || txStatusError || null,
    status: transactionStatus?.status,
    isSealed: transactionStatus?.status === 4,
  };
}