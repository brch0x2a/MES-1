package models

// User holds a users account information
type User struct {
	Username        string
	Authenticated   bool
	Usertype        int
	UserId          int
	Date            string
	Turn            int64
	Line            int64
	Profile_picture string
	Presentation    int64
	Coordinator     int64
	Config          bool
	WeightSubId     int
	CRQSId          int
}

func (u User) HasPermission(feature string) bool {
	if u.Usertype > 7 && feature == "Menu" {
		return true
	} else if u.Usertype > 2 && feature == "Warehouse" {
		return true
	} else {
		return false
	}
}

func (u User) HasPermissionTo(feature string) bool {
	if u.Usertype > 1 && feature == "Tag" {
		return true
	} else if u.Usertype > 1 && feature == "woOption" {
		return true
	} else {
		return false
	}
}

type User_holder struct {
	Id        int
	Nick_name string
	//Password  string
	Fname           string
	Lname           string
	Profile_picture string
	Privilege       string
	Factory         string
}

type MetaUser struct {
	Line            string
	LineID          string
	Presentation    string
	Profile_picture string
	Nick_name       string
	Fname           string
	Lname           string
	ProductPhoto    string
}
