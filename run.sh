go build -o hashicorp-plugin ./pkg/plugins/hashicorp/netrpc
go build -o buildmode=plugin -o go-plugin.so ./pkg/plugins/go-native/
go test -bench=.