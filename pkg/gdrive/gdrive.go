package gdrive

import (
	"context"
	b64 "encoding/base64"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GDrive struct{}

func New() *GDrive {
	return &GDrive{}
}

func (g *GDrive) AppendToSpreadsheet(ctx context.Context, spreadsheetID string, sheetID int, records [][]string) error {
	credBytes, err := b64.StdEncoding.DecodeString(os.Getenv("GDRIVE_KEY"))
	if err != nil {
		return err
	}

	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return err
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	batch := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{},
	}

	for _, record := range records {
		cells := []*sheets.CellData{}

		for idx := range record {
			data := &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{},
			}
			val := record[idx]

			if len(val) > 0 && val[0] == '=' {
				data.UserEnteredValue.FormulaValue = &val
			} else {
				data.UserEnteredValue.StringValue = &val
			}

			cells = append(cells, data)
		}

		req := &sheets.Request{
			AppendCells: &sheets.AppendCellsRequest{
				SheetId: int64(sheetID),
				Rows: []*sheets.RowData{
					{Values: cells},
				},
				Fields: "*",
			},
		}

		batch.Requests = append(batch.Requests, req)
	}

	res, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, &batch).Context(ctx).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return err // todo handle !200
	}

	return nil
}
