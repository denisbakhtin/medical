package models

// AboutPageID is a fixed ID for about page db record
const AboutPageID = 4

// ContactsPageID is a fixed ID for contacts page db record
const ContactsPageID = 7

// SessionPageID is a fixed ID for session page db record
const SessionPageID = 10

// MenuItem represents main menu item
type MenuItem struct {
	URL      string
	Title    string
	CSSClass string
	IsActive bool
	Children []MenuItem
}
