ARCH=
OS=
all: compile

armlinux:ARCH=arm
armlinux:OS=linux
armlinux:show
armlinux:compile

x86linux:ARCH=386
x86linux:OS=linux
x86linux:show
x86linux:compile

x64linux:ARCH=amd64
x64linux:OS=linux
x64linux:show
x64linux:compile

x86windows:ARCH=386
x86windows:OS=windows
x86windows:show
x86windows:compile

x64windows:ARCH=amd64
x64windows:OS=windows
x64windows:show
x64windows:compile

x86darwin:ARCH=386
x86darwin:OS=darwin
x86darwin:show
x86darwin:compile

x64darwin:ARCH=amd64
x64darwin:OS=darwin
x64darwin:show
x64darwin:compile

show:
	@echo "正在编译$(ARCH)平台,$(OS)内核的可执行文件"
help:
	@echo "请输入对应的系统(armlinux x86linux x64linux x86windows x64windows x86darwin x64darwin)"
compile:
	GOARCH=$(ARCH) GOOS=$(OS) CGO_ENABLED=1 GOPATH=$(shell pwd)  go build -o server-$(OS)-$(ARCH) dianxie/server
clean:
	rm -f server-*