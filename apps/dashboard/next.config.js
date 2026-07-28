/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: '/app',
  allowedDevOrigins: ['*.monkeycode-ai.online'],
  transpilePackages: ['@aige/ui', '@aige/types', '@aige/utils'],
  async rewrites() {
    return [
      {
        source: '/admin',
        destination: 'http://localhost:3001/admin',
        basePath: false,
      },
      {
        source: '/admin/:path*',
        destination: 'http://localhost:3001/admin/:path*',
        basePath: false,
      },
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ];
  },
};

module.exports = nextConfig;
