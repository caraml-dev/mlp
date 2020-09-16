# Build optimized production build
FROM node:13-alpine as builder

ARG REACT_APP_API_URL
ARG REACT_APP_ENVIRONMENT
ARG REACT_APP_OAUTH_CLIENT_ID
ARG REACT_APP_SENTRY_DSN

RUN mkdir -p /opt/mlp_ui
COPY . /opt/mlp_ui
WORKDIR /opt/mlp_ui

RUN yarn
RUN yarn lib build
RUN yarn app build

# Build serving image
FROM node:13-alpine
ENV PORT "8080"

EXPOSE ${PORT}

RUN yarn global add serve

RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /app

COPY --from=builder /opt/mlp_ui/packages/app/build build/

ENTRYPOINT ["sh", "-c"]
CMD ["echo 'Running UI...' && serve -s build -l ${PORT}"]