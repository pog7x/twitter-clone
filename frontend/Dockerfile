FROM node:latest AS builder
ENV VITE_BASE_API_URL=http://0.0.0.0:8080
WORKDIR /frontend
COPY . .
RUN [ -d "node_modules" ] && npm run build || npm install && npm run build

FROM nginx:latest
COPY --from=builder /frontend/dist/ /usr/share/nginx/html/
COPY ./nginx.conf /etc/nginx/conf.d/nginx.conf
CMD ["nginx", "-g", "daemon off;" ]