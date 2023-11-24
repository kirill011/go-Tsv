package databaseStruct

import (
	"gorm.io/gorm"
)

type Messages struct {
	gorm.Model
	N           int    `tsv:n`
	Mqtt        string `tsv: mqtt`
	Invid       string `tsv: invid`
	Unit_guid   string `tsv: unit_guid`
	Msg_id      string `tsv: msg_id`
	Text        string `tsv: text`
	Context     string `tsv: context`
	Class       string `tsv: class`
	Level       int    `tsv: level`
	Area        string `tsv: area`
	Addr        string `tsv: addr`
	Block       bool   `tsv: block`
	Type        string `tsv: type`
	Bit         int    `tsv: bit`
	Invert_bit  int    `tsv: invert_bit`
	LoadFilesID uint
}

type LoadFiles struct {
	gorm.Model
	Filename  string
	FileSaved bool
	Message   []Messages
}

type Errors struct {
	gorm.Model
	ErrorMessage string
	Filename     string
}
