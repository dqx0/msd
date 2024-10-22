package spreadsheet_test

import (
	"context"
	"os"
	"testing"

	sheet "github.com/dqx0/msd/pkg/spreadsheet"
	"github.com/stretchr/testify/assert"
)

func TestFileService_Write(t *testing.T) {
	ctx := context.Background()
	os.Setenv("MSD_FOLDER_ID", "test_folder_id")

	fileService, err := sheet.NewFileService(ctx, "シート1", "シート1")
	if err != nil {
		t.Fatal(err)
	}

	// テストデータ
	values := [][]interface{}{
		{"Item", "Cost", "Stocked", "Ship Date"},
		{"Wheel", "$20.50", "4", "3/1/2016"},
		{"Door", "$15", "2", "3/15/2016"},
		{"Engine", "$100", "1", "30/20/2016"},
	}

	// Write メソッドの呼び出し
	err = fileService.Write(values)

	// 結果の検証
	assert.NoError(t, err)
}
