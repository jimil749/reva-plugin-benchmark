go build -o hashicorp-plugin ./pkg/plugins/hashicorp/netrpc
go build -o hashicorp-plugin-grpc ./pkg/plugins/hashicorp/grpc/
go build -o buildmode=plugin -o go-plugin.so ./pkg/plugins/go-native/
go build -o pieplugin ./pkg/plugins/pie/
go test -bench=.