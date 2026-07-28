/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: '/admin',
  allowedDevOrigins: ['*.monkeycode-ai.online'],
  transpilePackages: ['@aige/ui', '@aige/types', '@aige/utils'],
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ];
  },
};

module.exports = nextConfig;
