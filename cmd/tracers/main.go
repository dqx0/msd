package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	formatter "github.com/dqx0/msd/internal/formatter"
	reader "github.com/dqx0/msd/internal/reader"
	sheet "github.com/dqx0/msd/pkg/spreadsheet"
)

func main() {
	// フラグの定義
	dataPath := flag.String("path", "", "path to the data directory")
	flag.Parse()

	// 環境変数の取得（データのパス）
	if *dataPath == "" {
		dataPath := os.Getenv("MSD_DATA_PATH")
		if dataPath == "" {
			fmt.Println("MSD_DATA_PATH is not set")
			os.Exit(1)
		}
	}
	fmt.Println("Data path: ", *dataPath)
	// 粒子の読み込み
	particles, nStep, err := reader.GetTracers(*dataPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// スプレッドシートの準備
	ctx := context.Background()
	sheetName := "nTrac=" + strconv.Itoa(len(particles)) + ", nStep=" + strconv.Itoa(nStep)
	fileService, err := sheet.NewFileService(ctx, sheetName+" "+time.Now().Format("2006-01-02-15-04-05"), sheetName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// データの作成
	values := formatter.Create(particles)

	// スプレッドシートに書き込み
	err = fileService.Write(values)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
