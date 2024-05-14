FROM golang:1.19-alpine AS build

RUN apk add --update --no-cache upx

WORKDIR /go/src/app

COPY . /go/src/app/
RUN go mod download
# Add options for building go with static library
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static'" -o ./backend
# Compress the binary file using Ultimate Packer for eXecutables
RUN upx ./backend

# Argument DB
ARG dbUrl
ARG dbUser
ARG dbPass
ARG dbName
ARG dbSecret

# Argument S3
ARG aws_s3_host
ARG aws_s3_access
ARG aws_s3_secret
ARG aws_s3_bucket
ARG aws_s3_zone

# Argument Redis
ARG redis_host
ARG redis_pass

# ENV DB
ENV MONGO_HOST=${dbUrl}
ENV MONGO_DBNAME=${dbName}
ENV MONGO_USER=${dbUser}
ENV MONGO_PASS=${dbPass}
ENV JWT_SECRET=${dbSecret}

# ENV S3
ENV AWS_S3_URL=${aws_s3_host}
ENV AWS_S3_ACCESS=${aws_s3_access}
ENV AWS_S3_SECRET=${aws_s3_secret}
ENV AWS_S3_BUCKET=${aws_s3_bucket}
ENV AWS_S3_ZONE=${aws_s3_zone}

# ENV REDIS
ENV REDIS_HOST=${redis_host}
ENV REDIS_PASS=${redis_pass}


# Stage 2 final
FROM scratch

# Argument DB
ARG dbUrl
ARG dbUser
ARG dbPass
ARG dbName
ARG dbSecret

# Argument S3
ARG aws_s3_host
ARG aws_s3_access
ARG aws_s3_secret
ARG aws_s3_bucket
ARG aws_s3_zone

# Argument Redis
ARG redis_host
ARG redis_pass

# ENV DB
ENV MONGO_HOST=${dbUrl}
ENV MONGO_DBNAME=${dbName}
ENV MONGO_USER=${dbUser}
ENV MONGO_PASS=${dbPass}
ENV JWT_SECRET=${dbSecret}

# ENV S3
ENV AWS_S3_URL=${aws_s3_host}
ENV AWS_S3_ACCESS=${aws_s3_access}
ENV AWS_S3_SECRET=${aws_s3_secret}
ENV AWS_S3_BUCKET=${aws_s3_bucket}
ENV AWS_S3_ZONE=${aws_s3_zone}

# ENV REDIS
ENV REDIS_HOST=${redis_host}
ENV REDIS_PASS=${redis_pass}

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/backend /backend

ENTRYPOINT ["/backend"]

EXPOSE 8080