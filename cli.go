package main

import (
	"io/ioutil"
	"os"
)

func main() {
	session := NewSession(os.Getenv("POLAR_USERNAME"), os.Getenv("POLAR_PASSWORD"))
	trainings := session.GetTrainings("2018-01-01", "2018-12-24", []int{Running})
	for _, t := range trainings {
		println(t.StartDate, t.Distance.String(), t.Duration.String())
		data := session.GetTrainingGpx(t.ID)
		ioutil.WriteFile(t.StartDate+".gpx", data, 0644)
	}
}
