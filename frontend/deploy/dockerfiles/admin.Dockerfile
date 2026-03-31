# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

# Copy workspace root and lockfile
COPY package.json package-lock.json ./

# Copy workspace packages (needed for workspace resolution)
COPY packages ./packages
COPY admin ./admin

# Install all dependencies and build
RUN npm ci
RUN npm run build -w admin

# Production stage
FROM node:22-alpine AS runner

WORKDIR /app

# Copy workspace root and lockfile
COPY package.json package-lock.json ./

# Copy workspace packages
COPY packages ./packages
COPY admin/package.json ./admin/

# Install production dependencies only
RUN npm ci --omit=dev

# Copy built admin app from builder
COPY --from=builder /app/admin/build ./admin/build

WORKDIR /app/admin

ARG COMMIT_HASH=unknown
ARG COMMIT_TIME=unknown
ARG BUILD_TIME=unknown
ARG VERSION=unknown

ENV NODE_ENV=production
ENV PORT=3000
ENV COMMIT_HASH=$COMMIT_HASH
ENV COMMIT_TIME=$COMMIT_TIME
ENV BUILD_TIME=$BUILD_TIME
ENV VERSION=$VERSION
EXPOSE 3000

CMD ["node", "build"]
