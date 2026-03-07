/** @type {import('next').NextConfig} */
const nextConfig = {
  // In local dev without `vercel dev`, proxy /api to a local Go server.
  // On Vercel, routing is handled by vercel.json.
  async rewrites() {
    if (process.env.NODE_ENV !== 'production') {
      return [
        {
          source: '/api/:path*',
          destination: `${process.env.API_URL ?? 'http://localhost:8080'}/api/:path*`,
        },
      ]
    }
    return []
  },
}

export default nextConfig
