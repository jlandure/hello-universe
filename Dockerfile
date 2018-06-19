FROM zenika/alpine-golang as BUILD
COPY handler.go main.go version.go /go/
RUN go get github.com/kelseyhightower/kargo
RUN apk add --no-cache gcc g++
RUN CGO_ENABLED=0 GOOS=linux go build -o hello-universe -a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo .


FROM scratch
COPY --from=BUILD /go/hello-universe .
ENTRYPOINT ["./hello-universe"]
