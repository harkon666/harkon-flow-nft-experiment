'use client';

import React from 'react';
import { useFlowCurrentUser } from '@onflow/react-sdk';
// Hook dari File 1
import { useGetListingIDs } from '@/hooks/scripts/useGetListingIds';
// Komponen dari File 3
import ListingCard from '@/components/cards/ListingCard';

const YourSalesSection: React.FC = () => {
  const { user } = useFlowCurrentUser();
  
  // Ambil daftar ID dari listingan (penjualan) Anda
  const { 
    listingIDs, 
    isLoading: isLoadingIDs, 
    error: errorIDs 
  } = useGetListingIDs({
    address: user?.addr ?? ""
  });
  console.log(listingIDs)
  // Tampilkan hanya jika pengguna sudah login
  if (!user?.loggedIn) {
    return null; // atau tampilkan pesan 'Silakan login'
  }

  return (
    <section className="container mx-auto px-4 py-16 text-center">
      <h2 className="text-2xl text-green-500 text-center mb-12 glow">Your Items for Sale</h2>

      {isLoadingIDs && <p>Loading your sale items...</p>}
      {errorIDs && <p className="text-red-500">Error loading sales: {errorIDs.message}</p>}

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
        {console.log(!isLoadingIDs && listingIDs && listingIDs.length > 0)}
        {!isLoadingIDs && listingIDs && listingIDs.length > 0 ? (
          // Map setiap ID ke komponen ListingCard
          listingIDs.map((id: number) => (
            <ListingCard 
              key={id} 
              listingResourceID={id} 
              ownerAddress={user.addr!} // '!' aman karena kita sudah cek login
            />
          ))
        ) : (
          !isLoadingIDs && <p className="col-span-3">You have no accessories listed for sale.</p>
        )}
      </div>
    </section>
  );
};

export default YourSalesSection;