package db_test

import (
	"fmt"
	"testing"

	. "github.com/etsy/mixer/db"
)

func TestParticipants(*testing.T) {
	participants := GetRandomPeopleParticipating("Managers")
	fmt.Printf("# people: %d\n", len(participants))

	for _, v := range participants {
		fmt.Printf("\ncurrent person: %#v\n", v)
	}

}

func TestInsertWeek(*testing.T) {
	w := InsertWeek(1, "Managers")
	fmt.Printf("# people: %#v\n", w)
}

func TestGetLastWeek(*testing.T) {
	w := GetLastWeek("Managers")
	fmt.Printf("week: %d\n", w)
}

func TestGetStaffData(*testing.T) {
	s := GetStaffData("genericuser")
	fmt.Printf("staff: %#v\n", s)
}

func TestCleanupAlumni(*testing.T) {
	CleanupAlumni()
}
