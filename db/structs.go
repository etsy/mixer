package db

type Page struct {
	Title string
	Body  []byte
}

type Person struct {
	Id           int64             `json:"id"`
	Name         string            `json:"name"`
	Username     string            `json:"username"`
	Assistant    string            `json:"assistant"`
	Avatar       string            `json:"avatar"`
	IsManager    bool              `json:"is_manager"`
	Mixers       *map[string]bool  `json:"mixers"`
	AssistantFor *map[string]int64 `json:"assistant_for"`
}

type Group struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Pair struct {
	Id      int64 `json:"id"`
	Week    int   `json:"week"`
	Person1 int64 `json:"person1"`
	Person2 int64 `json:"person2"`
	WeekId  int64 `json:"week_id"`
}

type Staff struct {
	Id            int    `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Title         string `json:"title"`
	IsManager     bool   `json:"is_manager"`
	Avatar        string `json:"avatar"`
	Enabled       int    `json:"enabled"`
	Auth_UserName string `json:"auth_username"`
}

type Week struct {
	Id     int64 `json:"id"`
	Number int   `json:"number"`
}
