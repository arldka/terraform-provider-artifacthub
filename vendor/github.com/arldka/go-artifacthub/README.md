# go-artifacthub

Auto generated using https://github.com/deepmap/oapi-codegen from https://artifacthub.io/docs/api/

## Usage

```
client, err := artifacthub.NewClientWithResponses("https://artifacthub.io/api/v1", artifacthub.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
  req.Header.Set("X-API-KEY-ID", "KEY_ID")
  req.Header.Set("X-API-KEY-SECRET", "KEY_SECRET")
  return nil
}))
if err != nil {
  return nil, err
}

resp, err := client.GetHelmPackageDetailsWithResponse(ctx, "repo", "package")
if err != nil {
  return nil, err
}
```
