# But why?

At a high level, I'd like to generate direct dependencies (using `--include-imports`), while also
being able to exclude certain paths or modules entirely.

For example, A -> B -> C(protovalidate). If I consume B, then C is "forced" on my project, even
though it'll never be used.

## Overview

The current repository has a single .proto file that depends on
[bufbuild/registry](https://buf.build/bufbuild/registry).

```shell
proto
└── api
    └── v1
        └── api.proto
```

which imports

```
import "buf/registry/owner/v1/user.proto";
```

Now, the [bufbuild/registry](https://buf.build/bufbuild/registry) module itself depends on
[bufbuild/protovalidate](https://buf.build/bufbuild/protovalidate).

However, there's no way for the current project to generate code that **does not** include the
protovalidate-related bits.

1. There's the blank import in the generated file: gen/go/buf/registry/owner/v1/user.pb.go

This is tricky because it's the code generator (`protoc-gen-go`) itself adding a blank import.

```diff
package ownerv1

import (
    _ "github.com/mfridman/try-protovalidate/gen/go/buf/registry/priv/extension/v1beta1"
-	_ "github.com/mfridman/try-protovalidate/gen/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)
```

2. There's generated code from the protovalidate module.

```
gen/go/buf/validate
├── expression.pb.go
├── priv
│   └── private.pb.go
└── validate.pb.go
```

## Possible solutions

I thought modifying the `buf.gen.yaml` might work:

```yaml
version: v2
inputs:
  - module: buf.build/bufbuild/protovalidate
    exclude_paths:
      - buf/validate
  - directory: proto/
```

But this results in an error:

> Failure: no .proto files were targeted. This can occur if no .proto files are found in your input,
> --path points to files that do not exist, **or --exclude-path excludes all files.**

Also this does not work:

```yaml
version: v2
inputs:
  - directory: proto/
    exclude_paths:
      - buf/validate
```

TL;DR - it's not possible to avoid generating the transitive dependency when consuming a BSR module
that depends on, say, [bufbuild/protovalidate](https://buf.build/bufbuild/protovalidate).

Following the example above, as a consumer of B I'd like to say A depends on _only_ B, but not on C.
