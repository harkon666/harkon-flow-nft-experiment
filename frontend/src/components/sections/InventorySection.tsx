'use client';

import React, { useState } from 'react';
import { useFlowCurrentUser } from '@onflow/react-sdk';

// 1. Impor KEDUA hook (untuk Momen dan Asesori)
import { useGetAccessoryIDs } from '@/hooks/scripts/useGetNFTAccessory'; // Hook Anda yang sudah ada
import { useGetMomentIDs } from '@/hooks/scripts/useGetNFTMoment'; // Hook BARU (dibuat di bawah)

// 2. Impor KEDUA komponen kartu
import AccessoryCard from '@/components/cards/AccessoryCard'; // Kartu Anda yang sudah ada
import MomentCard from '@/components/cards/MomentCard'; // Kartu BARU (dibuat di bawah)

// Tipe untuk state tab kita
type ActiveTab = 'moments' | 'accessories';

const InventorySection: React.FC = () => {
  const { user } = useFlowCurrentUser();
  
  // 3. State untuk melacak tab yang aktif, default-nya 'moments'
  const [activeTab, setActiveTab] = useState<ActiveTab>('moments');

  // 4. Panggil KEDUA hook di top-level (Aturan React)
  // Hook untuk Asesori
  const { 
    data: accessoryIDs, 
    isLoading: isLoadingAccessories, 
    error: errorAccessories,
    refetch: refetchAccessories
  } = useGetAccessoryIDs({
    address: user?.addr ?? '' // Kirim null jika belum login
  });

  // Hook untuk Momen
  const { 
    data: momentIDs, 
    isLoading: isLoadingMoments, 
    error: errorMoments,
    refetch: refetchMoments
  } = useGetMomentIDs({ // Hook baru Anda
    address: user?.addr ?? '' // Kirim null jika belum login
  });

  const handleTransactionSuccess = () => {
    console.log("Transaksi sukses! Me-refresh data...");
    refetchMoments();
    refetchAccessories();
  };

  // --- Styling untuk Tombol Tab ---
  // (Anda bisa pindahkan ini ke file CSS Anda)
  const baseTabClass = "px-6 py-2 transition-all pixel-text text-xl";
  const activeTabClass = "bg-green-500 text-black";
  const inactiveTabClass = "bg-gray-800 text-green-400 hover:bg-gray-700";
  // ---
  
  return (
    <section className="container mx-auto px-4 py-16 text-center">
      <h2 className="text-2xl text-green-500 text-center mb-8 glow">Your Inventory</h2>

      {/* 5. Render Tombol Tab */}
      <div className="flex justify-center mb-12 pixel-card-container p-1">
        <button
          className={`${baseTabClass} ${activeTab === 'moments' ? activeTabClass : inactiveTabClass}`}
          onClick={() => setActiveTab('moments')}
        >
          Moments
        </button>
        <button
          className={`${baseTabClass} ${activeTab === 'accessories' ? activeTabClass : inactiveTabClass}`}
          onClick={() => setActiveTab('accessories')}
        >
          Accessories
        </button>
      </div>

      {/* 6. Render Konten Tab secara Kondisional */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
        
        {/* Tampilkan Konten Tab Momen */}
        {activeTab === 'moments' && (
          <>
            {isLoadingMoments && <p className="col-span-3">Loading your moments...</p>}
            {errorMoments && <p className="col-span-3 text-red-500">Error loading moments: {errorMoments.message}</p>}
            
            {!isLoadingMoments && user?.addr && momentIDs && momentIDs.length > 0 ? (
              momentIDs.map((id: number) => (
                <MomentCard 
                  key={id} 
                  id={id} 
                  ownerAddress={user.addr ?? ''} 
                  onTransactionSuccess={handleTransactionSuccess}
                />
              ))
            ) : (
              !isLoadingMoments && <p className="col-span-3">Your Moments inventory is empty.</p>
            )}
          </>
        )}

        {/* Tampilkan Konten Tab Asesori */}
        {activeTab === 'accessories' && (
          <>
            {isLoadingAccessories && <p className="col-span-3">Loading your accessories...</p>}
            {errorAccessories && <p className="col-span-3 text-red-500">Error loading accessories: {errorAccessories.message}</p>}

            {!isLoadingAccessories && user?.addr && accessoryIDs && accessoryIDs.length > 0 ? (
              accessoryIDs.map((id: number) => (
                <AccessoryCard 
                  key={id} 
                  id={id} 
                  ownerAddress={user.addr ?? ''} 
                />
              ))
            ) : (
              !isLoadingAccessories && <p className="col-span-3">Your Accessories inventory is empty.</p>
            )}
          </>
        )}

      </div>
    </section>
  );
};

export default InventorySection;