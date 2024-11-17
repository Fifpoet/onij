# Base image for Go
FROM m.daocloud.io/docker.io/library/golang:1.22

ENV GOPROXY=https://goproxy.cn,direct

# Set the working directory
WORKDIR /app

# Copy Go project files
COPY . .

# Install dependencies and build the project
RUN go mod tidy
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]