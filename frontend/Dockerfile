FROM node:18-alpine AS builder
WORKDIR /frontend
COPY . .
CMD [ -d "node_modules" ] && npm run build || npm install && npm run build

FROM nginx:stable-alpine AS nginx
COPY --from=builder /frontend/dist/ /usr/share/nginx/html/
COPY ./nginx.conf /etc/nginx/conf.d/nginx.conf
CMD ["nginx", "-g", "daemon off;" ]