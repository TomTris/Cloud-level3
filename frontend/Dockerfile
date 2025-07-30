FROM node:20-alpine as build
WORKDIR /app

# Copy package files and install all dependencies including plugin-vue
COPY package*.json ./
COPY vite.config.js ./
RUN npm install

# Copy the rest of the code (including src/, index.html, etc)
COPY . .

# Build your app
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
# Serve under /, and proxy /api to backend
COPY nginx.conf /etc/nginx/conf.d/default.conf
