'use client';

import React, { useState, useEffect } from 'react';
import { useFlowCurrentUser } from '@onflow/react-sdk';
import { Package } from 'lucide-react'; // <-- Impor ikon yang Anda gunakan

// Impor hook transaksi
import { useReqGacha } from '@/hooks/transactions/useReqGacha';
import { useRevealGacha } from '@/hooks/transactions/useRevealGacha';
import { useCheckReceipt } from '@/hooks/scripts/useCheckReceipt';

// Impor hook inventaris (hanya untuk me-refresh)
import { useGetAccessoryIDs } from '@/hooks/scripts/useGetNFTAccessory';

const GachaSection: React.FC = () => {
  const { user } = useFlowCurrentUser();

  // --- State untuk UI (Sesuai kode Anda) ---
  const [isShaking, setIsShaking] = useState(false);
  const [showGachaResult, setShowGachaResult] = useState(false);
  const [gachaResult, setGachaResult] = useState<string | null>(null);

  //check receipt
  const {
    hasReceipt,
    isLoading: isLoadingCheckReceipt,
    refetch: refetchCheckReceipt
  } = useCheckReceipt({ address: user?.addr || ""})

  // --- Hook Blockchain ---
  // 1. Hook Tahap 1 (Request)
  const { 
    request, 
    isPending: isRequestPending, 
    isSealed: isRequestSealed,
    error: requestError,
    status
  } = useReqGacha();
  
  // 2. Hook Tahap 2 (Reveal)
  const { 
    reveal, 
    isPending: isRevealPending, 
    isSealed: isRevealSealed, 
    error: revealError 
  } = useRevealGacha();

  // 3. Hook Inventaris (Hanya untuk me-refresh)
  const { refetch: refetchInventory } = useGetAccessoryIDs({
    address: user?.addr ?? ""
  });
  console.log(status, 'woi')
  // Gabungkan status 'pending'
  const isGachaPending = isRequestPending || isRevealPending;

  // --- Logika untuk Merantai (Chain) Transaksi ---

  // Efek 1: Jika Tahap 1 (Request) selesai, otomatis jalankan Tahap 2 (Reveal)
  useEffect(() => {
    if (isRequestSealed) {
      console.log("Tahap 1 (Request) Sealed. Memulai Tahap 2 (Reveal)...");
      reveal({ recipient: user?.addr || "" });
    }
  }, [isRequestSealed]);

  // Efek 2: Jika Tahap 2 (Reveal) selesai, hentikan UI & refresh data
  useEffect(() => {
    if (isRevealSealed) {
      console.log("Tahap 2 (Reveal) Sealed. Gacha selesai!");
      setIsShaking(false);
      setGachaResult("Anda mendapatkan item baru! Cek inventaris Anda.");
      setShowGachaResult(true);
      
      // Refresh data inventaris on-chain Anda
      refetchInventory();
      refetchCheckReceipt();
    }
  }, [isRevealSealed, refetchInventory]);

  // Efek 3: Tangani jika ada error di tengah jalan
  const gachaError = requestError || revealError;
  useEffect(() => {
    if (gachaError) {
      console.error("Gacha Gagal:", gachaError);
      setIsShaking(false);
      // Ubah pesan error agar lebih singkat
      let friendlyError = gachaError.message.includes("No Receipt found") 
        ? "Anda belum membeli paket." 
        : "Transaksi gacha gagal.";
      
      setGachaResult(friendlyError);
      setShowGachaResult(true);
    }
  }, [gachaError]);

  // --- Fungsi Handle Tombol (Sekarang Jauh Lebih Sederhana) ---
  const handleBuyGacha = () => {
    if (!user?.addr) {
      alert("Harap hubungkan dompet (wallet) Anda!");
      return;
    }

    if (hasReceipt) {
      setIsShaking(true);
      setShowGachaResult(false);
      setGachaResult(null);
      reveal({ recipient: user?.addr || "" });
    } else {
      // Reset UI dan mulai animasi
      setIsShaking(true);
      setShowGachaResult(false);
      setGachaResult(null);
  
      // Mulai alur dengan memanggil Tahap 1
      // Efek (useEffect) akan menangani sisanya
      request();
    }
  };

  // --- JSX (Tampilan dari Anda) ---
  return (
    <section className="container mx-auto px-4 py-16 bg-gradient-to-b from-black to-gray-900">
      <h2 className="text-2xl text-green-500 text-center mb-12 glow">Gacha Accessory Pack</h2>
      <div className="max-w-md mx-auto">
        <div className="pixel-card text-center">
          
          {/* Container Gambar (dengan 'shake') */}
          <div className={`w-full aspect-square bg-gray-900 border-2 border-green-500 mb-6 flex items-center justify-center ${isShaking ? 'shake' : ''}`}>
            <Package size={120} className="text-green-500" />
          </div>
          
          <p className="text-green-400 text-xl mb-6 pixel-text">Price: 0 FLOW (FREE)</p>
          
          {/* Tombol Beli */}
          <button
            onClick={handleBuyGacha}
            // Nonaktifkan tombol saat transaksi sedang diproses
            disabled={isGachaPending} 
            className="pixel-button w-full mb-4 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {/* Ubah teks tombol berdasarkan status */}
            {
             isLoadingCheckReceipt ? "[LOADING CHECK RECEIPT]" :
             isRequestPending ? "[ 1/2 REQUESTING... ]" :
             isRevealPending ? "[ 2/2 REVEALING... ]" :
             hasReceipt ? "[ GACHA ]" : "[ BUY GACHA PACK ]"}
          </button>

          {/* Tampilan Hasil (Sesuai kode Anda) */}
          {showGachaResult && (
            <div className="mt-6 p-4 border-2 border-green-500 bg-black">
              <p className="text-green-500 mb-2 pixel-text text-sm">
                {gachaError ? "Error:" : "You got:"}
              </p>
              <p className={`text-lg glow ${gachaError ? 'text-red-500' : 'text-green-400'}`}>
                {gachaResult}!
              </p>
            </div>
          )}
        </div>
      </div>
    </section>
  );
};

export default GachaSection;