FROM --platform=linux/amd64 node:21.2.0-alpine as build

ARG REACT_APP_BACKEND_URL
ENV REACT_APP_BACKEND_URL=$REACT_APP_BACKEND_URL

WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH

COPY package.json ./

RUN yarn install --immutable --immutable-cache --check-cache --silent
RUN yarn add react-scripts@5.0.1 -g --silent
COPY . ./
RUN yarn build

# production environment
FROM --platform=linux/amd64 nginx:stable-alpine

COPY --from=build /app/build /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]