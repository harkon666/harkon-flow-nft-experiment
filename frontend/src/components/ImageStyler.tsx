'use client'; // (Jika Anda menggunakan Next.js App Router, tambahkan ini)

import React, { useState, useEffect } from 'react';
import { Upload } from 'lucide-react'; // (Opsional, untuk ikon yang bagus)
import BingkaiKayu from '@/assets/Bingkai Kayu.png' 

const ImageStyler: React.FC = () => {
  // 1. State untuk menyimpan URL sementara dari gambar yang di-upload
  const [momentImage, setMomentImage] = useState<string | null>(null);

  // 2. Fungsi untuk menangani saat pengguna memilih file
  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    // Pastikan ada file yang dipilih
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      
      // Hapus URL lama (jika ada) untuk mencegah 'memory leak'
      if (momentImage) {
        URL.revokeObjectURL(momentImage);
      }
      
      // Buat URL sementara HANYA di browser pengguna
      // Ini super cepat dan tidak perlu upload
      setMomentImage(URL.createObjectURL(file));
    }
  };

  // 3. Efek (Effect) untuk 'cleanup' (membersihkan)
  //    Wajib untuk mencegah 'memory leak' dari URL.createObjectURL
  useEffect(() => {
    // Ini adalah 'cleanup function' yang akan berjalan
    // saat komponen di-unmount atau saat 'momentImage' berubah
    return () => {
      if (momentImage) {
        URL.revokeObjectURL(momentImage);
      }
    };
  }, [momentImage]); // Jalankan efek ini setiap kali 'momentImage' berubah

  // --- Render (Tampilan) Komponen ---
  return (
    <section className="container mx-auto px-4 py-16">
      <h2 className="text-2xl font-bold text-center mb-8">
        Demo Penumpukan (Equip) NFT
      </h2>
      
      <div className="max-w-md mx-auto">
        
        {/* --- 4. AREA PREVIEW (KANVAS) --- */}
        {/* Di sinilah keajaiban CSS terjadi */}
        <div className="relative w-full aspect-square bg-gray-100 rounded-xl overflow-hidden border-2 border-gray-200 shadow-lg">
          
          {/* Layer 1: Foto Pengguna (z-index 10) */}
          {momentImage ? (
            <img
              src={momentImage}
              alt="Preview Momen"
              className="absolute inset-0 w-full h-full object-cover z-10"
            />
          ) : (
            // Tampilan placeholder jika belum ada gambar
            <div className="absolute inset-0 z-5 flex items-center justify-center p-4">
              <p className="text-gray-500 text-center font-mono">
                Silakan upload foto Momen Anda
              </p>
            </div>
          )}

          {/* Layer 2: Bingkai Kayu (z-index 20) */}
          <img
            src={BingkaiKayu}
            alt="Bingkai Kayu"
            // 'pointer-events-none' agar bingkai tidak bisa diklik
            className="absolute inset-0 w-full h-full object-cover z-20 pointer-events-none"
          />
        </div>
        
        {/* --- 5. TOMBOL UPLOAD --- */}
        <div className="mt-6">
          {/* Kita gunakan <label> untuk men-trigger <input> yang disembunyikan */}
          <label 
            htmlFor="image-upload" 
            className="cursor-pointer w-full inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-medium rounded-lg text-white bg-green-600 hover:bg-green-700 transition-all"
          >
            <Upload size={20} className="mr-2" />
            Upload Foto Momen
          </label>
          
          {/* Input file yang sebenarnya, tapi disembunyikan */}
          <input
            id="image-upload"
            name="image-upload"
            type="file"
            // Hanya terima file gambar
            accept="image/png, image/jpeg, image/webp" 
            className="hidden"
            onChange={handleImageUpload}
          />
        </div>
        
        <p className="text-center text-sm text-gray-500 mt-4">
          Catatan: Ini adalah simulasi *frontend* (CSS stacking). Untuk rilis produksi, Anda perlu menambahkan *cropper* 1:1 terlebih dahulu.
        </p>
      </div>
    </section>
  );
};

export default ImageStyler;