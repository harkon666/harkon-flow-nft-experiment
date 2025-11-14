// src/lib/utils.ts

/**
 * Mengubah URL IPFS (ipfs://) menjadi URL HTTPS
 * yang bisa dibaca oleh browser melalui gateway publik.
 */
export const resolveIpfsUrl = (url: string | null | undefined): string => {
  if (!url || url === '') {
    return ''; // Kembalikan string kosong jika tidak ada URL
  }

  // Cek apakah URL sudah merupakan HTTPS
  if (url.startsWith('https://') || url.startsWith('http://')) {
    return url;
  }

  // Cek jika ini adalah URL IPFS
  if (url.startsWith('ipfs://')) {
    // Ganti 'ipfs://' dengan URL gateway
    // Anda bisa mengganti 'ipfs.io' dengan gateway lain
    // seperti 'cloudflare-ipfs.com' atau 'gateway.pinata.cloud'
    return url.replace('ipfs://', 'https://white-lazy-marten-351.mypinata.cloud/ipfs/');
  }

  // Kembalikan apa adanya jika formatnya tidak dikenal
  return url;
};