'use client';

import React, { useState, useEffect } from 'react';
import { useFlowCurrentUser } from '@onflow/react-sdk';
import { Upload, CheckCircle, AlertTriangle } from 'lucide-react'; // Impor ikon

// Tentukan props yang diterima komponen ini
interface MintMomentSectionProps {
  // Kita perlu fungsi 'refetch' dari hook 'useGetMomentIDs'
  // agar inventaris bisa otomatis update setelah minting
  refetchMoments: () => void;
}

const MintMomentSection: React.FC<MintMomentSectionProps> = ({ refetchMoments }) => {
  const { user } = useFlowCurrentUser();
  
  // State untuk data formulir
  const [name, setName] = useState<string>('');
  const [description, setDescription] = useState<string>('');
  
  // State untuk file gambar
  const [thumbnailFile, setThumbnailFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  
  // State untuk UI (loading & pesan)
  const [isPending, setIsPending] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Fungsi untuk menangani 'upload' file dan membuat 'preview'
  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    // Reset pesan
    setError(null);
    setSuccess(null);

    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      setThumbnailFile(file);

      // (Di sini Anda akan menjalankan 'cropper' 1:1 Anda)
      // Untuk demo ini, kita langsung pakai file-nya

      // Buat URL preview di browser
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl); // Hapus URL lama
      }
      setPreviewUrl(URL.createObjectURL(file));
    }
  };

  // Cleanup memory leak dari URL.createObjectURL
  useEffect(() => {
    return () => {
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl);
      }
    };
  }, [previewUrl]);

  // Fungsi untuk 'submit' ke API backend Go
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    // 1. Validasi Input
    if (!user?.addr) {
      setError("Harap hubungkan dompet (wallet) Anda terlebih dahulu.");
      return;
    }
    if (!thumbnailFile) {
      setError("Harap pilih gambar untuk Momen Anda.");
      return;
    }
    if (name.trim() === '') {
      setError("Nama Momen tidak boleh kosong.");
      return;
    }

    setIsPending(true);

    // 2. Buat FormData
    // Ini WAJIB untuk mengirim file ke backend Echo Anda
    const formData = new FormData();
    formData.append("recipient", user.addr);
    formData.append("name", name);
    formData.append("description", description);
    formData.append("thumbnail", thumbnailFile); // "thumbnail" harus cocok dengan c.FormFile("thumbnail")

    try {
      // 3. Panggil API Backend Anda
      const response = await fetch("http://localhost:8000/moment", {
        method: "POST",
        body: formData,
        // JANGAN set 'Content-Type', browser akan otomatis
      });

      const result = await response.json();

      if (!response.ok) {
        // Tangkap error dari backend (misal: "Konten dilarang")
        throw new Error(result.error || "Gagal minting. Terjadi kesalahan.");
      }

      // 4. Sukses!
      setSuccess(`Minting sukses! URL IPFS: ${result.thumbnail}`);
      
      // Kosongkan form
      setName('');
      setDescription('');
      setThumbnailFile(null);
      setPreviewUrl(null);

      // Panggil refetch untuk memperbarui inventaris!
      refetchMoments();

    } catch (err: any) {
      setError(err.message || "Gagal menghubungi server.");
    } finally {
      setIsPending(false);
    }
  };

  return (
    <section className="container mx-auto px-4 py-16">
      <h2 className="text-2xl text-green-500 text-center mb-12 glow">Mint Your Moment</h2>
      
      {/* Ubah div luar menjadi <form> */}
      <form onSubmit={handleSubmit} className="max-w-md mx-auto">
        <div className="pixel-card text-center">

          {/* --- AREA PREVIEW GAMBAR (Diganti) --- */}
          <div className="w-full aspect-square bg-gray-900 border-2 border-green-500 mb-6 flex items-center justify-center overflow-hidden">
            {previewUrl ? (
              // 1. Tampilkan gambar preview jika ada
              <img
                src={previewUrl}
                alt="Preview Momen"
                className="w-full h-full object-cover"
                style={{ imageRendering: 'pixelated' }}
              />
            ) : (
              // 2. Tampilkan tombol upload jika kosong
              <label 
                htmlFor="moment-upload" 
                className="cursor-pointer flex flex-col items-center text-green-400 opacity-70 hover:opacity-100 transition-opacity"
              >
                <Upload size={64} />
                <span className="pixel-text text-sm mt-4">Pilih Gambar (1:1)</span>
              </label>
            )}
          </div>
          {/* Input file yang sebenarnya, tapi disembunyikan */}
          <input
            id="moment-upload"
            type="file"
            accept="image/png, image/jpeg, image/webp"
            className="hidden"
            onChange={handleImageChange}
            disabled={isPending}
          />
          {/* --- AKHIR AREA PREVIEW --- */}

          {/* --- Input Teks Baru --- */}
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Nama Momen (Wajib)"
            className="w-full p-3 bg-black border-2 border-green-500 text-green-400 pixel-text mb-4"
            disabled={isPending}
          />
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Deskripsi Momen (Opsional)"
            rows={3}
            className="w-full p-3 bg-black border-2 border-green-500 text-green-400 pixel-text mb-6"
            disabled={isPending}
          />
          {/* --- Akhir Input Teks --- */}

          <p className="text-green-400 text-xl mb-6 pixel-text">Price: FREE</p>
          
          <button 
            type="submit" 
            className="pixel-button w-full disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={isPending || !thumbnailFile || !user?.addr} // Nonaktifkan jika sedang 'pending' atau 'file' kosong
          >
            {isPending ? "MINTING..." : "[ MINT MOMENT ]"}
          </button>

          {/* --- Pesan Status --- */}
          {error && (
            <div className="mt-4 p-3 border-2 border-red-500 bg-black text-red-500 pixel-text text-sm flex items-center gap-2">
              <AlertTriangle size={16} />
              <span>{error}</span>
            </div>
          )}
          {success && (
            <div className="mt-4 p-3 border-2 border-green-500 bg-black text-green-500 pixel-text text-sm flex items-center gap-2">
              <CheckCircle size={16} />
              <span>{success}</span>
            </div>
          )}

        </div>
      </form>
    </section>
  );
};

export default MintMomentSection;