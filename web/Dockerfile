# Stage 1: Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install build dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Copy source code and build Nuxt application
COPY . .
RUN npm run build

# Stage 2: Production run stage
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

# Copy built output from builder stage
COPY --from=builder /app/.output ./.output

EXPOSE 3000

# Run the Nitro standalone server
CMD ["node", ".output/server/index.mjs"]
