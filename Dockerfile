# Base image for Node.js
FROM m.daocloud.io/docker.io/library/node:20.14

# Set the working directory
WORKDIR /app

# Copy package.json and pnpm-lock.yaml for dependency installation
COPY package.json pnpm-lock.yaml ./

# Install pnpm
RUN npm install -g pnpm

# Set Taobao registry for pnpm
RUN pnpm config set registry https://registry.npmmirror.com

# Install dependencies
RUN pnpm install

# Copy the rest of the project files
COPY . .

# Expose the application port
EXPOSE 18968

# Run the application
CMD ["npm", "run", "dev"]