FROM golang:1.19.4-bullseye

WORKDIR /app

RUN useradd --create-home gopher \
  && chown -R gopher:gopher /app

USER gopher

COPY --chown=gopher:gopher . .
