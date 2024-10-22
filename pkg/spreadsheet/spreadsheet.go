package spreadsheet

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type FileService struct {
	Sheet         *sheets.Spreadsheet
	DriveService  *drive.Service
	SheetsService *sheets.Service
}

type IFileService interface {
	Write(values [][]interface{}) error
}

func NewFileService(ctx context.Context, spreadsheetName, sheetName string) (IFileService, error) {
	// ドライブサービスの初期化
	ds, err := NewDriveService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create drive service: %v", err)
	}

	// シートサービスの初期化
	ss, err := NewSheetService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet service: %v", err)
	}

	// スプレッドシートの作成
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: spreadsheetName,
		},
		Sheets: []*sheets.Sheet{
			{
				Properties: &sheets.SheetProperties{
					Title: sheetName, // シート（タブ）名を明示的に設定
				},
			},
		},
	}

	sheet, err := ss.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create spreadsheet: %v", err)
	}

	return &FileService{
		Sheet:         sheet,
		DriveService:  ds,
		SheetsService: ss,
	}, nil
}

func (fs *FileService) Write(values [][]interface{}) error {
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	writeRange := fmt.Sprintf("%s!A1", fs.Sheet.Sheets[0].Properties.Title)

	_, err := fs.SheetsService.Spreadsheets.Values.Update(
		fs.Sheet.SpreadsheetId,
		writeRange,
		valueRange,
	).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("failed to update values: %v", err)
	}

	folderID := os.Getenv("MSD_FOLDER_ID")
	if folderID == "" {
		return fmt.Errorf("MSD_FOLDER_ID is not set")
	}

	file := &drive.File{}
	_, err = fs.DriveService.Files.Update(
		fs.Sheet.SpreadsheetId,
		file,
	).AddParents(folderID).Fields("id, parents").Do()

	if err != nil {
		return fmt.Errorf("failed to move file to folder: %v", err)
	}

	fmt.Printf("スプレッドシート '%s' を作成しました\n", fs.Sheet.Properties.Title)
	fmt.Printf("URL: https://docs.google.com/spreadsheets/d/%s\n", fs.Sheet.SpreadsheetId)

	return nil
}

func NewDriveService(ctx context.Context) (*drive.Service, error) {
	log.Printf("Creating drive service with application default credentials...")

	scopes := []string{
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/drive.file",
	}

	credentials, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		log.Printf("Failed to get default credentials: %v", err)
		return nil, err
	}

	driveService, err := drive.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Printf("Failed to create drive service: %v", err)
		return nil, err
	}

	log.Printf("Successfully created drive service")
	return driveService, nil
}

func NewSheetService(ctx context.Context) (*sheets.Service, error) {
	log.Printf("Creating sheets service with application default credentials...")

	scopes := []string{
		"https://www.googleapis.com/auth/spreadsheets",
	}

	credentials, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		log.Printf("Failed to get default credentials: %v", err)
		return nil, err
	}

	sheetsService, err := sheets.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Printf("Failed to create sheets service: %v", err)
		return nil, err
	}

	log.Printf("Successfully created sheets service")
	return sheetsService, nil
}
