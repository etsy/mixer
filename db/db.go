package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/etsy/mixer/config"

	_ "github.com/etsy/mixer/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
)

var (
	MixerDb *sql.DB
)

// Connect() needs to be called before anything else.
func ConnectMixer() {
	var err error

	if nil == MixerDb {
		err = config.Config.Load()
		if err != nil {
			log.Println("error reading or parsing config config:", err)
		}

		MixerDb, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.Config.Database["mixer"].User,
			config.Config.Database["mixer"].Password,
			config.Config.Database["mixer"].Hostname,
			config.Config.Database["mixer"].Port,
			config.Config.Database["mixer"].Name))

		// https://github.com/go-sql-driver/mysql/issues/257#issuecomment-53886663
		MixerDb.SetMaxIdleConns(0)   // this is the root problem! set it to 0 to remove all idle connections
		MixerDb.SetMaxOpenConns(500) // or whatever is appropriate for your setup.

		if err != nil {
			log.Println(err)
		}
	}
}

///////////////////////////////////////////////////////////////////////////////////////////

// modifyPersonGroups will insert a new group for a person or delete
// their old groups based on their preferences
func modifyPersonGroups(p Person) {
	ConnectMixer()

	for k, v := range *p.Mixers {
		if v == true {
			stmt, err := MixerDb.Prepare("insert ignore into people_groups (people_id, groups_id) values (?, (select id from groups where name = ?))")
			if err != nil {
				log.Println(err)
				return
			}
			_, err = stmt.Exec(p.Id, k)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			stmt, err := MixerDb.Prepare("delete from people_groups where people_id = ? and groups_id = (select id from groups where name = ?)")
			if err != nil {
				log.Println(err)
				return
			}
			_, err = stmt.Exec(p.Id, k)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func UpdatePerson(p Person) {
	ConnectMixer()

	stmt, err := MixerDb.Prepare("update people set name = ?, assistant = ? where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec(p.Name, p.Assistant, p.Id)
	if err != nil {
		log.Println(err)
		return
	}

	// rectify the groups this person is in
	modifyPersonGroups(p)
}

func InsertPerson(p Person) Person {
	ConnectMixer()
	stmt, err := MixerDb.Prepare("insert ignore into people (name, username, assistant) values (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return p
	}
	res, err := stmt.Exec(p.Name, p.Username, p.Assistant)
	if err != nil {
		log.Println(err)
		return p
	}

	person_id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error:", err.Error())
	} else {
		p.Id = person_id
		modifyPersonGroups(p)
	}

	return p
}

func InsertWeek(week int, group string) Week {
	ConnectMixer()
	w := Week{}
	stmt, err := MixerDb.Prepare("insert ignore into week (number, date, group_id) values (?, ?, (select id from groups where name = ?))")
	if err != nil {
		log.Println(err)
		return w
	}
	res, err := stmt.Exec(week, time.Now().Unix(), group)
	if err != nil {
		log.Println(err)
		return w
	}

	week_id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error:", err.Error())
	} else {
		w.Id = week_id
		w.Number = week
	}

	return w
}

func InsertPair(w Week, p1 Person, p2 Person) {
	ConnectMixer()
	stmt, err := MixerDb.Prepare("insert ignore into pairs (person1, person2, week_id) values (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec(p1.Id, p2.Id, w.Id)
	if err != nil {
		log.Println(err)
	}
}

func InsertOddPerson(w Week, p Person) {
	ConnectMixer()
	stmt, err := MixerDb.Prepare("insert into odd_person_out (person, week_id) values (?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec(p.Id, w.Id)
	if err != nil {
		log.Println(err)
	}
}

func GetLastWeek(groupname string) int {
	ConnectMixer()
	var week int = 0
	err := MixerDb.QueryRow("SELECT IFNULL(MAX(number), 0) FROM week where group_id = (select id from groups where name = ?)", groupname).Scan(&week)
	if err != nil {
		log.Print(err)
	}

	return week
}

func GetPersonDataFromUsername(username string) Person {
	ConnectMixer()

	p := Person{}
	err := MixerDb.QueryRow("SELECT id FROM people WHERE username = ?", username).Scan(&p.Id)
	if err != nil {
		// new user
		p.Id = 0
		p.Username = username
		// so as not to deref a nil pointer
		temp := make(map[string]bool)
		p.Mixers = &temp
		s := GetStaffData(p.Username)
		p.Avatar = s.Avatar
		p.Name = s.FirstName + " " + s.LastName
		p.IsManager = s.IsManager
		return p
	} else {
		p, err = GetPersonData(p.Id)
		if err != nil {
			log.Print(err)
		}
		return p
	}

	return p
}

func GetRandomPeopleParticipating(groupname string) []Person {
	ConnectMixer()

	people := make([]Person, 0)
	rows, err := MixerDb.Query("select p.id, p.name, p.username, p.assistant from people p join people_groups pg on p.id = pg.people_id and pg.groups_id = (select id from groups where name = ?) order by rand()", groupname)
	if err != nil {
		log.Println(err)
		return people
	}
	defer rows.Close()

	for rows.Next() {
		p := Person{}
		if err := rows.Scan(&p.Id, &p.Name, &p.Username, &p.Assistant); err != nil {
			log.Println(err)
		}

		people = append(people, p)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	return people
}

func CleanupAlumni() {
	ConnectMixer()

	rows, err := MixerDb.Query("select id, name, username, assistant from people")
	if err != nil {
		log.Print(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := Person{}
		if err := rows.Scan(&p.Id, &p.Name, &p.Username, &p.Assistant); err != nil {
			log.Println(err)
		}

		s := GetStaffData(p.Username)

		// delete their mixers if they are no longer enabled/employed
		if s.Enabled == 0 {
			fmt.Printf("disabled: %#v\n\n", p)
			stmt2, err2 := MixerDb.Prepare("delete from people_groups where people_id = ?")
			if err2 != nil {
				log.Println(err)
			}
			_, err2 = stmt2.Exec(p.Id)
			if err2 != nil {
				log.Println(err2)
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

}

func GetPeopleData(groupname string) []Person {
	ConnectMixer()

	people := make([]Person, 0)
	rows, err := MixerDb.Query("select p.id, p.name, p.username, p.assistant from people p join people_groups pg on p.id = pg.people_id and pg.groups_id = (select id from groups where name = ?)", groupname)
	if err != nil {
		log.Println(err)
		return people
	}
	defer rows.Close()

	// now get the staff data for these people
	for rows.Next() {
		p := Person{}
		if err := rows.Scan(&p.Id, &p.Name, &p.Username, &p.Assistant); err != nil {
			log.Println(err)
		}

		s := GetStaffData(p.Username)
		if len(s.Avatar) > 0 {
			p.Avatar = s.Avatar
		} else {
			p.Avatar = config.Config.Staff.DefaultAvatarUrl
		}

		people = append(people, p)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	return people

}

func InsertStaffData(s Staff) {
	ConnectMixer()

	stmt, err := MixerDb.Prepare(`insert into staff (auth_username, staff_id, first_name, last_name, title, is_manager, avatar, enabled)
                                  values (?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update enabled = ?, title = ?, is_manager = ?, avatar = ?`)
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(s.Auth_UserName, s.Id, s.FirstName, s.LastName, s.Title, s.IsManager, s.Avatar, s.Enabled,
		s.Enabled, s.Title, s.IsManager, s.Avatar)
	if err != nil {
		log.Println(err)
	}
}

func GetStaffData(username string) (p Staff) {
	ConnectMixer()
	staff := Staff{}
	err := MixerDb.QueryRow(`select staff_id, first_name, last_name, title, IFNULL(avatar, ""), is_manager, enabled
	                         from staff where auth_username = ?`, username).Scan(
		&staff.Id, &staff.FirstName, &staff.LastName, &staff.Title,
		&staff.Avatar, &staff.IsManager, &staff.Enabled)
	if err != nil {
		// just return empty staff object if this fails
		/*log.Printf("couldnt find staff data for %s", username)*/
		return staff
	}

	// also accept these titles as a manager
	if strings.Contains(staff.Title, "Manager") ||
		strings.Contains(staff.Title, "Staff") ||
		strings.Contains(staff.Title, "Principal") ||
		strings.Contains(staff.Title, "Distinguished") ||
		strings.Contains(staff.Title, "Director") {
		staff.IsManager = true
	}

	return staff
}

func GetPersonData(id int64) (Person, error) {
	ConnectMixer()

	p := Person{}
	err := MixerDb.QueryRow(`select p.id, p.name, p.username, p.assistant
		                     from people p
		                     left join people_groups pg on p.id = pg.people_id
		                     where p.id = ?`, id).Scan(&p.Id, &p.Name, &p.Username, &p.Assistant)
	if err != nil {
		p.Id = 0
		return p, fmt.Errorf("user not found for %v", id)
	}

	// get the mixer data //////////////////////////////
	p, err = getMixersForPerson(p)
	if err != nil {
		log.Println(err)
		return p, err
	}

	// get the staff data //////////////////////////////
	s := GetStaffData(p.Username)
	p.Avatar = s.Avatar
	p.IsManager = s.IsManager
	if p.Name == "" {
		p.Name = s.FirstName + " " + s.LastName
	}

	// get "assistant for" data (if they are an assistant)
	p, err = getAssistantForData(p)
	if err != nil {
		return p, err
	}

	return p, nil
}

func getAssistantForData(p Person) (Person, error) {
	ConnectMixer()

	rows, err := MixerDb.Query("select id, username from people where assistant = ?", p.Username)
	if err != nil {
		return p, err
	}
	defer rows.Close()

	for rows.Next() {

		boss := Person{}
		if err := rows.Scan(&boss.Id, &boss.Username); err != nil {
			log.Println(err)
		}

		if p.AssistantFor == nil {
			// so as not to deref a nil pointer
			temp := make(map[string]int64)
			p.AssistantFor = &temp
		}
		(*p.AssistantFor)[boss.Username] = boss.Id
	}

	return p, nil
}

func GetMixers() ([]Group, error) {
	ConnectMixer()

	groups := make([]Group, 0)
	rows, err := MixerDb.Query("select id, name, description from groups")
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {

		g := Group{}
		if err := rows.Scan(&g.Id, &g.Name, &g.Description); err != nil {
			log.Println(err)
		}

		groups = append(groups, g)
	}

	return groups, nil
}

func GetMixerData(groupname string) (Group, error) {
	ConnectMixer()

	g := Group{}
	err := MixerDb.QueryRow("select id, name, description from groups where name = ?", groupname).Scan(&g.Id, &g.Name, &g.Description)

	return g, err
}

func getMixersForPerson(p Person) (Person, error) {
	ConnectMixer()

	rows, err := MixerDb.Query("select id, name, description from groups")
	if err != nil {
		return p, err
	}
	defer rows.Close()

	for rows.Next() {

		g := Group{}
		if err := rows.Scan(&g.Id, &g.Name, &g.Description); err != nil {
			log.Println(err)
		}

		var mixergroup bool = false
		err = MixerDb.QueryRow(`select count(pg.id)
		                        from people_groups pg
		                        where pg.people_id = ? and pg.groups_id = ?`,
			p.Id, g.Id).Scan(&mixergroup)

		if err != nil {
			return p, err
		}
		if p.Mixers == nil {
			// so as not to deref a nil pointer
			temp := make(map[string]bool)
			p.Mixers = &temp
		}
		(*p.Mixers)[g.Name] = mixergroup
		if err != nil {
			return p, fmt.Errorf("error finding groups for person id %v", p.Id)
		}
	}

	return p, nil
}

func GetLastPairing(cur_id int64, other_id int64) Pair {
	ConnectMixer()

	p := Pair{}
	err := MixerDb.QueryRow(`select id, person1, person2, week_id
	                         from pairs where person1 IN (?, ?) and person2 IN (?, ?)
	                         order by week_id desc limit 1`,
		cur_id, other_id, cur_id, other_id).Scan(&p.Id, &p.Person1, &p.Person2, &p.WeekId)
	if err != nil {
		// return some 404 status to the app
		p.Id = 0
		/*log.Print(err)*/
	}

	return p
}
