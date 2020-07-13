FROM golang:1
COPY . /work
WORKDIR /work
RUN make build

FROM gcr.io/distroless/base
COPY --from=0 /work/server /server
ENTRYPOINT ["/server"]
