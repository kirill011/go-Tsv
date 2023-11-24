package filereaderwriter

import (
	"fmt"
	"go-Tsv/internal/config"
	"go-Tsv/internal/database"
	"go-Tsv/internal/database/databaseStruct"
	"os"
	"strconv"
	"strings"

	"github.com/dogenzaka/tsv"
	"github.com/fumiama/go-docx"
)

type FileReaderWriter struct {
}

func New() *FileReaderWriter {
	return &FileReaderWriter{}
}

type TsvRow struct {
	N          string `tsv:n`
	Mqtt       string `tsv: mqtt`
	Invid      string `tsv: invid`
	Unit_guid  string `tsv: unit_guid`
	Msg_id     string `tsv: msg_id`
	Text       string `tsv: text`
	Context    string `tsv: context`
	Class      string `tsv: class`
	Level      string `tsv: level`
	Area       string `tsv: area`
	Addr       string `tsv: addr`
	Block      string `tsv: block`
	Type       string `tsv: type`
	Bit        string `tsv: bit`
	Invert_bit string `tsv: invert_bit`
}

func (f *FileReaderWriter) ReadFile(fileName string, db *database.Database, cfg *config.Config) error {

	file, err := os.Open(cfg.DirInput + "/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	row := TsvRow{}

	table := make([]databaseStruct.Messages, 0, 21)

	parser, err := tsv.NewParser(file, &row)
	if err != nil {
		return err
	}

	for {
		eof, err := parser.Next()
		if err != nil {
			return err
		}
		n, err := strconv.Atoi(strings.ReplaceAll(row.N, " ", ""))
		if err != nil {
			return err
		}
		level, err := strconv.Atoi(strings.ReplaceAll(row.Level, " ", ""))
		if err != nil {
			return err
		}
		bit := 0
		if row.Bit != "" {
			bit, err = strconv.Atoi(strings.ReplaceAll(row.Bit, " ", ""))
			if err != nil {
				return err
			}
		}
		invertBit := 0
		if row.Bit != "" {
			invertBit, err = strconv.Atoi(strings.ReplaceAll(row.Invert_bit, " ", ""))
			if err != nil {
				return err
			}
		}
		block := false
		if row.Block != "" {
			block, err = strconv.ParseBool(strings.ReplaceAll(row.Block, " ", ""))
			if err != nil {
				return err
			}
		}
		table = append(table, databaseStruct.Messages{N: n,
			Mqtt:       strings.ReplaceAll(row.Mqtt, " ", ""),
			Invid:      strings.ReplaceAll(row.Invid, " ", ""),
			Unit_guid:  strings.ReplaceAll(row.Unit_guid, " ", ""),
			Msg_id:     strings.ReplaceAll(row.Msg_id, " ", ""),
			Text:       strings.ReplaceAll(row.Text, " ", ""),
			Context:    strings.ReplaceAll(row.Context, " ", ""),
			Class:      strings.ReplaceAll(row.Class, " ", ""),
			Level:      level,
			Area:       strings.ReplaceAll(row.Area, " ", ""),
			Addr:       strings.ReplaceAll(row.Addr, " ", ""),
			Block:      block,
			Type:       strings.ReplaceAll(row.Type, " ", ""),
			Bit:        bit,
			Invert_bit: invertBit,
		})
		if eof {
			break
		}

	}

	fileStruct := databaseStruct.LoadFiles{Filename: fileName, FileSaved: true, Message: table}

	err = db.InsertMessages(&fileStruct)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileReaderWriter) CreateFile(fileName string, db *database.Database, cfg *config.Config) error {
	guids, err := db.SelectGuid(fileName)
	if err != nil {
		return err
	}

	for _, guid := range guids {
		newFile := docx.NewA4()
		res, err := db.SelectMessages(guid.LoadFilesId, guid.Unit_guid)
		if err != nil {
			return err
		}

		table := newFile.AddTable(len(res)+1, 15)

		for i, row := range table.TableRows {
			if i == 0 {
				insertTextToCell(row, 0, "n")
				insertTextToCell(row, 1, "mqtt")
				insertTextToCell(row, 2, "invid")
				insertTextToCell(row, 3, "unit_guid")
				insertTextToCell(row, 4, "msg_id")
				insertTextToCell(row, 5, "text")
				insertTextToCell(row, 6, "context")
				insertTextToCell(row, 7, "class")
				insertTextToCell(row, 8, "level")
				insertTextToCell(row, 9, "area")
				insertTextToCell(row, 10, "addr")
				insertTextToCell(row, 11, "block")
				insertTextToCell(row, 12, "type")
				insertTextToCell(row, 13, "bit")
				insertTextToCell(row, 14, "invert_bit")
			} else {
				insertTextToCell(row, 0, strconv.Itoa(res[i-1].N))
				insertTextToCell(row, 1, res[i-1].Mqtt)
				insertTextToCell(row, 2, res[i-1].Invid)
				insertTextToCell(row, 3, res[i-1].Unit_guid)
				insertTextToCell(row, 4, res[i-1].Msg_id)
				insertTextToCell(row, 5, res[i-1].Text)
				insertTextToCell(row, 6, res[i-1].Context)
				insertTextToCell(row, 7, res[i-1].Class)
				insertTextToCell(row, 8, strconv.Itoa(res[i-1].Level))
				insertTextToCell(row, 9, res[i-1].Area)
				insertTextToCell(row, 10, res[i-1].Addr)
				insertTextToCell(row, 11, strconv.FormatBool(res[i-1].Block))
				insertTextToCell(row, 12, res[i-1].Type)
				insertTextToCell(row, 13, strconv.Itoa(res[i-1].Bit))
				insertTextToCell(row, 14, strconv.Itoa(res[i-1].Invert_bit))
			}

		}

		cfg.InfoLog.Println(fmt.Sprintf("Create file %s/%s.docx", cfg.DirOutput, guid.Unit_guid))
		f, err := os.Create(fmt.Sprintf("%s/%s.docx", cfg.DirOutput, guid.Unit_guid))
		if err != nil {
			return err
		}

		_, err = newFile.WriteTo(f)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func insertTextToCell(row *docx.WTableRow, cellNumber int, text string) {
	par := row.TableCells[cellNumber].AddParagraph()
	par.AddText(text)
}
