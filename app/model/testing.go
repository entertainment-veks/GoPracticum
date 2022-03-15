package model

func TestLink() *Link {
	return &Link{
		Link:   "https://www.google.com",
		Code:   "CODE",
		UserID: "userid",
	}
}
