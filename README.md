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

### Work in Progress tasks

- Analysis module selection logic
- Database integration
- Cache for recent queries
- gRPC communication between the different modules
- User facing API
- Configuration for each module (maybe)
- DNS service (maybe)
- Add more collectors, ingestors and analyzers (maybe)



---

### Dependencies

We only support the latest stable **Go** version (Go **1.11.2** as of now). You need to have installed `protoc` (`libprotoc 3.6.1`), the Protocol Buffers Compiler and support for gRPC.

This project should work on any latest Linux or macOS systems, but note that it has only been actually tested under the following operating systems:

- Linux 4.15 (Ubuntu 18.04.1 LTS)
- macOS 10.14.1 (18B75)

There is no reason why it shouldn't work on Windows, but we haven't tested it.

---

### Build steps

After installing the dependencies, in order to build and run our project you  have to download it and build the binaries. To do that you just have to run (outside of your `$GOPATH`) the following shell commands:

```sh
export GO11MODULE=on
git clone --single-branch -b master https://github.com/DarkAnHell/FastPhish
cd FastPhish
go mod init github.com/DarkAnHell/FastPhish
go build ./...
# for example, run the whoisds data collector.
bin/whoisds
# or check the results for the levenshtein analyzer
bin/analyzer
```

Note that the binaries provided are only there as a quick test, and will later on be the full modules
