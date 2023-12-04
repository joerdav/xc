# My readme

## Tasks

### generate-templ

```bash
go run -mod=mod github.com/a-h/templ/cmd/templ generate
go mod tidy
```

### generate-translations

```bash
go run ./i18n/generate
```

### generate-all

Requires: generate-templ, generate-translations
RunDeps: async
