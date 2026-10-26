import type { NextConfig } from "next";

const backend = process.env.API_INTERNAL_URL ?? "http://127.0.0.1:8080";

const config: NextConfig = {
  poweredByHeader: false,
  async rewrites() {
    // The browser uses one origin. Cookies pass through the Next.js proxy.
    return [{ source: "/api/:path*", destination: `${backend}/api/:path*` }];
  },
};

export default config;
