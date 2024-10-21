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

func NewFileService(sheetName string) (IFileService, error) {
	ctx := context.Background()
	ds, err := NewDriveService(ctx)
	if err != nil {
		return nil, err
	}
	ss, err := NewSheetService(ctx)
	if err != nil {
		return nil, err
	}
	sheet, err := NewSheet(sheetName)
	if err != nil {
		return nil, err
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
	writeRange := fmt.Sprintf("%s!A1", fs.Sheet.Properties.Title)
	_, err := fs.SheetsService.Spreadsheets.Values.Update(fs.Sheet.SpreadsheetId, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	folderID := os.Getenv("MSD_FOLDER_ID")
	if folderID == "" {
		log.Fatalf("MSD_FOLDER_ID is not set")
	}

	file := &drive.File{}

	_, err = fs.DriveService.Files.Update(fs.Sheet.SpreadsheetId, file).
		AddParents(folderID).
		Fields("id, parents").
		Do()

	if err != nil {
		log.Fatalf("Unable to move file to folder after retries: %v", err)
	}

	fmt.Printf("スプレッドシートURL: https://docs.google.com/spreadsheets/d/%s\n", fs.Sheet.SpreadsheetId)

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

func NewSheet(title string) (*sheets.Spreadsheet, error) {
	ctx := context.Background()

	sheetsService, err := NewSheetService(ctx)
	if err != nil {
		return nil, err
	}

	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
	}

	spreadsheet, err = sheetsService.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		return nil, err
	}

	return spreadsheet, nil
}
