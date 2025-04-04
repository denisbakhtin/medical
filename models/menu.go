package models

const AboutPageId = 4
const ContactsPageId = 7
const SessionPageId = 10

// MenuItem represents main menu item
type MenuItem struct {
	URL      string
	Title    string
	CSSClass string
	IsActive bool
	Children []MenuItem
}
