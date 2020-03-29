FROM golang:latest 
RUN echo deb http://deb.debian.org/debian testing main contrib non-free >> /etc/apt/sources.list
RUN echo deb-src http://deb.debian.org/debian testing main contrib non-free >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y thrift-compiler
RUN mkdir /app 
ADD go.mod /app/
ADD go.sum /app/
WORKDIR /app/ 
RUN go mod download
ADD . /app/
CMD ["make", "test"]