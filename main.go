package main

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strconv"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) []student {
	var students []student
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, row := range data {
		if i > 0 { // omit header line
			st, err := createStudentFromRow(row)
			if err != nil {
				log.Fatal(err)
			}
			students = append(students, st)
		}
	}
	return students
}

func createStudentFromRow(row []string) (student, error) {
	var st student
	if len(row) != 7 {
		return student{}, errors.New("invalid csv row")
	}
	st.firstName = row[0]
	st.lastName = row[1]
	st.university = row[2]

	t1Score, errT1 := strconv.Atoi(row[3])
	t2Score, errT2 := strconv.Atoi(row[4])
	t3Score, errT3 := strconv.Atoi(row[5])
	t4Score, errT4 := strconv.Atoi(row[6])
	if errT1 != nil || errT2 != nil || errT3 != nil || errT4 != nil {
		return student{}, errors.New("invalid csv row")
	}

	st.test1Score = t1Score
	st.test2Score = t2Score
	st.test3Score = t3Score
	st.test4Score = t4Score
	return st, nil
}

func calculateGrade(students []student) []studentStat {
	var studentStats []studentStat
	for _, st := range students {
		var ss studentStat
		avg := float32(st.test1Score+st.test2Score+st.test3Score+st.test4Score) / 4
		ss.student = st
		ss.finalScore = avg
		ss.grade = getGrade(avg)
		studentStats = append(studentStats, ss)
	}
	return studentStats
}

func getGrade(score float32) Grade {
	if score >= 70 {
		return A
	} else if score < 70 && score >= 50 {
		return B
	} else if score < 50 && score >= 35 {
		return C
	} else {
		return F
	}
}

// EdgeCase - There could be multiple students with same finalScore
func findOverallTopper(gradedStudents []studentStat) studentStat {
	maxScore := float32(0)
	var topper studentStat
	for _, s := range gradedStudents {
		if s.finalScore > maxScore {
			maxScore = s.finalScore
			topper = s
		}
	}
	return topper
}

// EdgeCase - There could be multiple students with same finalScore per university
func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	topperPerUni := make(map[string]studentStat)
	studentsUniMap := createStudentUniMap(gs)
	for uni, students := range studentsUniMap {
		t := findOverallTopper(students)
		topperPerUni[uni] = t
	}
	return topperPerUni
}

func createStudentUniMap(gs []studentStat) map[string][]studentStat {
	studentsUniMap := make(map[string][]studentStat)
	for _, stStat := range gs {
		studentsUniMap[stStat.university] = append(studentsUniMap[stStat.university], stStat)
	}
	return studentsUniMap
}
