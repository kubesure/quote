FROM node:8.15.1-alpine
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 8030
CMD [ "npm", "start" ]