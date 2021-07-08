package models

type MetaAM_instance struct {
	Id           int
	Planned_init string
	Planned_end  string
	Line         string
	LineId       int
	JobId        int
	Lila         string
	LilaColor    string
	Component    string
	Phase        int
	State        string
}

type MetaLogAM_instance struct {
	Id           int
	Planned_init string
	Planned_end  string
	Line         string
	LineId       int

	JobId       int
	Description string

	Machine        string
	Component      string
	ComponentPhoto string

	EPP      string
	EPPPhoto string

	Lila      string
	LilaColor string

	OperatorProfile  string
	OperatorNickName string
	OperatorFname    string
	OperatorLname    string

	JobInit     string
	JobEmd      string
	MinutesStop int
	MinutesRun  int

	ApproverProfile  string
	ApproverNickName string
	ApproverFname    string
	ApproverLname    string
	State            string
	StateColor       string

	Phase int

	Note string
}

type AMStat struct {
	Id             int
	Line           string
	TotalJobs      int
	OpenJobs       int
	InProgressJobs int
	ClosedJobs     int
}
