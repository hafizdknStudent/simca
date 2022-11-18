package helper

import (
	"strings"

	"github.com/antchfx/htmlquery"
)

type DataMataKuliah struct {
	KodeMk      string `json:"kode_mk"`
	MataKuliah  string `json:"mata_kuliah"`
	SksTeori    string `json:"sks_teori"`
	SksPrak     string `json:"sks_prak"`
	KodeDosen   string `json:"kode_dosen"`
	Ruang       string `json:"ruang"`
	Hari        string `json:"hari"`
	Jam         string `json:"jam"`
	Kelas       string `json:"kelas"`
	StatusKelas string `json:"status_kelas"`
	BU          string `json:"bu"`
}

type ListMataKuliah struct {
	Data []DataMataKuliah
}

func ParseAnswerQuestion(url string) (string, error) {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return "", err
	}

	xpathAnswerValue := "//*[@placeholder='jawaban']//following-sibling::input"

	respNode := htmlquery.FindOne(doc, xpathAnswerValue)
	trueAnswerValue := htmlquery.SelectAttr(respNode, "value")

	return trueAnswerValue, err
}

func ParseDataTableMatkul(datDoc string) map[string][]string {
	doc, err := htmlquery.Parse(strings.NewReader(datDoc))
	if err != nil {
		panic(err)
	}

	xpathTableTrMatkul := "//h3/../..//tbody/tr"
	xpathTableTdMatkul := "//td"

	tr := htmlquery.Find(doc, xpathTableTrMatkul)
	dataMatkul := make(map[string][]string, 0)

	for index, item := range tr {
		if index == 0 {
			continue
		}

		tmp := make([]string, 0)
		td := htmlquery.Find(item, xpathTableTdMatkul)

		for _, itemTd := range td {
			tmp = append(tmp, htmlquery.InnerText(itemTd))
		}

		dataMatkul[tmp[0]] = tmp
	}

	return dataMatkul
}

func StoreDataMatkul(dataMatkul map[string][]string) ListMataKuliah {
	var listMatkul ListMataKuliah
	for _, item := range dataMatkul {
		t := DataMataKuliah{
			KodeMk:      item[0],
			MataKuliah:  item[1],
			SksTeori:    item[2],
			SksPrak:     item[3],
			KodeDosen:   item[4],
			Ruang:       item[5],
			Hari:        item[6],
			Jam:         item[7],
			Kelas:       item[8],
			StatusKelas: item[9],
			BU:          item[10],
		}
		listMatkul.Data = append(listMatkul.Data, t)
	}
	return listMatkul
}
