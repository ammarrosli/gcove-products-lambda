build:
	echo "Building lambda binaries"
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/products/bootstrap cmd/lambda/products/main.go
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/authorizer/bootstrap cmd/lambda/authorizer/main.go

zip:
	zip -j build/lambda/products.zip build/lambda/products/bootstrap
	zip -j build/lambda/authorizer.zip build/lambda/authorizer/bootstrap
