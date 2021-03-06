version := 0.3
packageNameNix := panner-linux-amd64-$(version).tar.gz
packageNameMac := panner-darwin-amd64-$(version).tar.gz
packageNameWindows := panner-windows-amd64-$(version).tar.gz

build_dir := output
build_dir_linux := output-linux
build_dir_mac := output-mac
build_dir_windows := output-windows

build: format configure build-linux build-mac build-windows

build_m: format configure build-mac
	cp ./$(build_dir_mac)/panner ./

format:
	go fmt ./...


configure:
		mkdir -p $(build_dir)
		mkdir -p $(build_dir_linux)
		mkdir -p $(build_dir_mac)
		mkdir -p $(build_dir_windows)


build-linux:
		env GOOS=linux GOARCH=amd64 go build -o ./$(build_dir_linux)/panner -ldflags "-X main.version=$(version)"
		@cd ./$(build_dir_linux) && tar zcf ../$(build_dir)/$(packageNameNix) .

build-mac:
		env GOOS=darwin GOARCH=amd64 go build -o ./$(build_dir_mac)/panner -ldflags "-X main.version=$(version)"
		@cd ./$(build_dir_mac) && tar zcf ../$(build_dir)/$(packageNameMac) .

build-windows:
		env GOOS=windows GOARCH=amd64 go build -o ./$(build_dir_windows)/panner.exe -ldflags "-X main.version=$(version)"
		@cd ./$(build_dir_windows) && tar zcf ../$(build_dir)/$(packageNameWindows) .

clean:
		rm -rf $(build_dir)
		rm -rf $(build_dir_linux)
		rm -rf $(build_dir_mac)
		rm -rf $(build_dir_windows)
