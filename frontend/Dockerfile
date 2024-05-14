FROM node:latest as build-stage
ENV VITE_BASE_API_URL=http://0.0.0.0:8080
WORKDIR /frontend
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:latest as production-stage
COPY --from=build-stage /frontend/dist/ /usr/share/nginx/html/
COPY ./nginx.conf /etc/nginx/conf.d/nginx.conf
CMD ["nginx", "-g", "daemon off;" ]