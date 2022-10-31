ARG MLP_API_IMAGE
FROM ${MLP_API_IMAGE}

COPY ui/build ./ui/build

CMD ["mlp"]
