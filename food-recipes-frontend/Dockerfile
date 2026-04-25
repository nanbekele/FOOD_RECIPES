# Use stable Node.js version
FROM node:20.18-alpine

# Set working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Copy frontend source code
COPY . .

# Expose frontend port
EXPOSE 3000

# Start Nuxt dev server
CMD ["npm", "run", "dev"]
