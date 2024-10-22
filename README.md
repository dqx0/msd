## MSD計算用CLI
[core/msd.go](pkg/core/msd.go)
### 前提条件
- GCP Google Drive API, Google Sheets API, ADC
- 環境変数 MSD_FOLDER_ID(ファイル出力先)
- Go 1.23~
- データファイル
### 実行方法
```go
go run cmd/tracers/main.go -path=/path/to/data/folder
```
