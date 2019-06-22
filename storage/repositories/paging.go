package repositories

type Paging struct {
	Limit    int
	NextPage string
	LastPage string
	Sort     string
}

func DefaultPager() Paging {
	p := Paging{}
	p.Limit = 20
	p.LastPage = "0"
	p.NextPage = "20"
	p.Sort = "asc"
	return p
}
