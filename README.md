CyberCamp 2018 Hackathon entry for team Phish 'n Chips.

## FastPhish

FastPhish aims to be a fast and reliable phishing detection framework. It's based in Go.

The following modules have been already included:

- Analysis
  - Levenshtein distance
- Ingestion Engine
- Data sources
  - Certificate Transparency Logs
  - Passive sources
    - Whoisds

**Only master branch is stable**. Please, if you use other branches, do it at your own risk.



---

### Dependencies

We only support the latest stable **Go** version (Go **1.11.2** as of now).

You need to have installed `protoc` (`libprotoc 3.6.1`), the Protocol Buffers Compiler and support for gRPC. In order to install them, please follow the official guide [here](https://google.github.io/proto-lens/installing-protoc.html) for `protoc` and make sure you `go get` the needed gRPC packages listed in the official [guide](https://grpc.io/docs/quickstart/go.html).


You also need Redis installed, for example following [this](https://www.digitalocean.com/community/tutorials/how-to-install-and-secure-redis-on-ubuntu-18-04) (only step 1 is necessary)

This project should work on any latest Linux or macOS systems, but note that it has only been actually tested under the following operating systems:

- Linux 4.15 (Ubuntu 18.04.1 LTS)
- macOS 10.14.1 (18B75)

There is no reason why it shouldn't work on Windows, but we haven't tested it.

---

### Build steps

After installing the dependencies, in order to build and run our project you have to download it and build the binaries. To do that you just have to run (outside of your `$GOPATH`) the following shell commands:

```sh
export GO11MODULE=on
# make sure you have the most recent version of proto and protoc-gen-go
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
git clone --single-branch -b master https://github.com/DarkAnHell/FastPhish
cd FastPhish
make api
make certs
go mod init github.com/DarkAnHell/FastPhish
make build
```

### Using

Edit or create config files for the binaries (examples avaliable at example_configs) before hand, and make sure to have a redis DB launched

- Launch DB (should **always** be present):
```
bin/db <redis-config.json>
```

- Launch Analyzer (should **always** be present):
```
bin/analyzer <levenshtein-config.json>
```

- Launch Certificate Transparency Logs parser: if you want to get data from CT Logs.
```
bin/ctdemo <ctlogs.json>
```

- Launch user API (should **always** be present if you use the HTTP API or the `aux_client` module):
```
bin/api
```

- Launch CLI Client:

```
bin/aux_client
```

- Launch HTTP API:

```
bin/http
```

- Use HTTP API:

```sh
curl --silent --header "Content-Type: application/json" \
  --request POST \
  --data '{"domain":"twistter.com"}' \
  http://localhost:8080/
```



If you want to use the API connection, you can write your own gRPC client to connect to it. You have an example at **aux_client**, which you can also launch to see a prepared execution