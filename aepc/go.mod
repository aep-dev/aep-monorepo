module github.com/aep-dev/aep-monorepo/aepc

go 1.24.0

require (
	buf.build/gen/go/aep/api/protocolbuffers/go v1.36.10-20251109183837-26a011a354ee.1
	cloud.google.com/go/longrunning v0.7.0
	github.com/aep-dev/aep-monorepo/aep-lib-go v0.0.0-20260901132022-b69ac1bc6386
	github.com/google/cel-go v0.27.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.0
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/rs/cors v1.11.1
	github.com/spf13/cobra v1.10.2
	google.golang.org/genproto/googleapis/api v0.0.0-20251202230838-ff82c1b0f217
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.12
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.12-20260709200747-435963d16310.1 // indirect
	cel.dev/expr v0.25.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jarcoal/httpmock v1.4.0 // indirect
	github.com/jhump/protoreflect v1.17.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/exp v0.0.0-20251009144603-d2f985daa21b // indirect
	golang.org/x/net v0.48.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/ghodss/yaml v1.0.0
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)

// uncomment for local development
// replace github.com/aep-dev/aep-monorepo/aep-lib-go => ../aep-lib-go
