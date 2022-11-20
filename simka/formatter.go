package simka

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"simka/helper"
)

type EventFormatter struct {
	Summary     string
	Description string
	StartHour   int
	StartMin    int
	EndHour     int
	EndMin      int
	Recurrence  string
	HtoDay      int
}

func CalendarEventFormatter(events helper.ListMataKuliah) ([]EventFormatter, error) {
	var ruangan, matkul, dosen, kelas, hari, jam string
	listEvent := make([]EventFormatter, 0)

	for _, event := range events.Data {
		matkul = event.MataKuliah
		dosen = event.KodeDosen
		ruangan = event.Ruang
		kelas = event.Kelas
		hari = event.Hari
		jam = event.Jam

		summary := fmt.Sprintf("%s, Ruang %s", matkul, ruangan)
		description := fmt.Sprintf("Matkul %s\nRuang %s\nDosen %s\nKelas %s", matkul, ruangan, dosen, kelas)
		day, htoDay := convertDay(hari)
		clockEvent, err := parseTime(jam)
		if err != nil {
			return listEvent, err
		}

		evnt := EventFormatter{
			Summary:     summary,
			Description: description,
			StartHour:   clockEvent["startHour"],
			StartMin:    clockEvent["startEnd"],
			EndHour:     clockEvent["endHour"],
			EndMin:      clockEvent["endMin"],
			Recurrence:  fmt.Sprintf("RRULE:FREQ=WEEKLY;BYDAY=%s;BYMONTH=", day),
			HtoDay:      htoDay,
		}

		listEvent = append(listEvent, evnt)
	}

	return listEvent, nil
}

func parseTime(jam string) (map[string]int, error) {
	var startHour, startMin, endHour, endMin int
	var errorParseClock bool
	var err error
	errMsg := "Can't parse time"

	defaultStart := 7
	defaultMinStart := 0

	defaultEnd := 8
	defaultMinEnd := 0

	clockEvent := map[string]int{}

	// jam = 10.20 - 10.40
	// split - untuk memishakan waktu masuk kelas dan akhir kelas
	rawParse := strings.Split(jam, "-")

	if len(strings.ReplaceAll(rawParse[0], " ", "")) == 0 {
		startHour = defaultStart
		startMin = defaultMinStart
		endHour = defaultEnd
		endMin = defaultMinEnd
		errorParseClock = true
	}

	if !errorParseClock {

		// hasil split 10.20 untuk masuk kelas 10.40 untuk keluar kelas
		// split . untuk mengambil jam dan menit
		start := strings.Split(rawParse[0], ".")
		end := strings.Split(rawParse[1], ".")

		startHour, err = strconv.Atoi(strings.ReplaceAll(start[0], " ", ""))
		if err != nil {
			return clockEvent, errors.New(errMsg)
		}
		startMin, err = strconv.Atoi(strings.ReplaceAll(start[1], " ", ""))
		if err != nil {
			return clockEvent, errors.New(errMsg)
		}

		endHour, err = strconv.Atoi(strings.ReplaceAll(end[0], " ", ""))
		if err != nil {
			return clockEvent, errors.New(errMsg)
		}
		endMin, err = strconv.Atoi(strings.ReplaceAll(end[1], " ", ""))
		if err != nil {
			startMin = defaultMinStart
		}

	}

	clockEvent["startHour"] = startHour
	clockEvent["startMin"] = startMin
	clockEvent["endHour"] = endHour
	clockEvent["endMin"] = endMin

	return clockEvent, nil
}

func getDay(hari string) (string, int) {
	var currDay string
	var intCurrDay int

	switch hari {
	case "Senin":
		currDay = "MO"
		intCurrDay = 1
	case "Selasa":
		currDay = "TU"
		intCurrDay = 2
	case "Rabu":
		currDay = "WE"
		intCurrDay = 3
	case "Kamis":
		currDay = "TH"
		intCurrDay = 4
	case "Jumat":
		currDay = "FR"
		intCurrDay = 5
	case "Sabtu":
		currDay = "SA"
		intCurrDay = 6
	case "Minggu":
		currDay = "SU"
		intCurrDay = 7
	default:
		currDay = "SA"
		intCurrDay = 6
	}

	return currDay, intCurrDay
}

func convertDay(hari string) (string, int) {
	rawDay := time.Now().Weekday()
	intDay := int(rawDay)

	index := 0
	overFlow := false
	statusOverFlow := false

	day, targetIntday := getDay(hari)
	hTodayH := 0

	for true {
		// cek apakah index sekarang melebihi batas hari yg ada
		if index == 7 {
			overFlow = true
		}

		// jika overflow reset index ke 0
		// ubah overFlow ke false
		// dan buat statusOverFlow ke true
		if overFlow {
			index = 0
			overFlow = false
			statusOverFlow = true
		}

		// ini akan memastikan index harus sama dengan int weekday hari ini
		// jika hari ini jumat dan index di hari senin, maka lakukan iterasi,
		// sampai hari index sama dengan hari jumat
		if !overFlow && !statusOverFlow && index >= intDay {
			hTodayH += 1
			if index == targetIntday {
				break
			}
		}

		if statusOverFlow {
			hTodayH += 1
			if index == targetIntday {
				break
			}
		}

		index += 1
	}

	return day, hTodayH - 1
}
