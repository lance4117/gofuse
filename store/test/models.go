package test

type Blog struct {
	Id      int    `xorm:"not null pk autoincr INT"`
	Title   string `xorm:"VARCHAR(100)"`
	Content string `xorm:"TEXT"`
	Ctm     int64  `xorm:"BIGINT"`
	Uid     int    `xorm:"INT"`
}

type User struct {
	Id   int64  `xorm:"pk autoincr BIGINT"`
	Name string `xorm:"VARCHAR(100)"`
	Ctm  int64  `xorm:"BIGINT"`
}
