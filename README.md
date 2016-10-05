# Hawkulark
## Hawkular Kubernetes Agent

## Building
```
mkdir GOPATH
export GOPATH=`pwd`/GOPATH
mkdir -p ~GOPATH/src/github.com/hawkular
cd ~GOPATH/src/github.com/hawkular
git clone https://github.com/${USERID}/hawkulark
cd hawkulark
git remote add upstream https://github.com/hawkular/hawkulark
export PATH=$GOPATH/bin:$PATH
make
```
