# Build stage
FROM node:18-alpine as build

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Create the directory and copy built files
RUN mkdir -p /usr/share/nginx/html
COPY --from=build /app/dist/. /usr/share/nginx/html/

# Expose port 5000
EXPOSE 5000

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
