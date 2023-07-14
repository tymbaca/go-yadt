FROM golang:latest as BUILDER

WORKDIR /workspaces/go-yadt
COPY go.* .
RUN go mod download
COPY . .

# Create a test binary
RUN CGO_ENABLED=0 GOOS=linux go test -c -o go-yadt.test


FROM mcr.microsoft.com/dotnet/sdk:7.0-alpine as PAGEMERGER_BUILDER
WORKDIR /pagemerger
RUN git clone https://github.com/tymbaca/pagemerger.git .
RUN dotnet add package DocumentFormat.OpenXml --version 2.20.0 
RUN dotnet add package CommandLineParser --version 2.9.1 
RUN dotnet publish -o ./result -p:PublishSingleFile=true --self-contained false


FROM alpine:latest

# Install dotnetcore
RUN apk add aspnetcore7-runtime

# Install Pagemerger
COPY --from=PAGEMERGER_BUILDER /pagemerger/result/pagemerger /usr/local/bin

WORKDIR /workspaces/go-yadt
COPY --from=BUILDER /workspaces/go-yadt/ .

ENTRYPOINT ./go-yadt.test -test.v
