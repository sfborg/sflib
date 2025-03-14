package coldp

import "strings"

type ReferenceType int

const (
	UnknownRT ReferenceType = iota
	Article
	ArticleJournal
	ArticleMagazine
	ArticleNewspaper
	Bill
	Book
	Broadcast
	Chapter
	Dataset
	Entry
	EntryDictionary
	EntryEncyclopedia
	Figure
	Graphic
	Interview
	LegalCase
	Legislation
	ManuscriptRT
	Map
	MotionPicture
	MusicalScore
	Pamphlet
	PaperConference
	Patent
	PersonalCommunication
	Post
	PostWeblog
	Report
	Review
	ReviewBook
	Song
	Speech
	Thesis
	Treaty
	Webpage
)

var refTypeToString = map[ReferenceType]string{
	Article:               "ARTICLE",
	ArticleJournal:        "ARTICLE_JOURNAL",
	ArticleMagazine:       "ARTICLE_MAGAZINE",
	ArticleNewspaper:      "ARTICLE_NEWSPAPER",
	Bill:                  "BILL",
	Book:                  "BOOK",
	Broadcast:             "BROADCAST",
	Chapter:               "CHAPTER",
	Dataset:               "DATASET",
	Entry:                 "ENTRY",
	EntryDictionary:       "ENTRY_DICTIONARY",
	EntryEncyclopedia:     "ENTRY_ENCYCLOPEDIA",
	Figure:                "FIGURE",
	Graphic:               "GRAPHIC",
	Interview:             "INTERVIEW",
	LegalCase:             "LEGAL_CASE",
	Legislation:           "LEGISLATION",
	ManuscriptRT:          "MANUSCRIPT",
	Map:                   "MAP",
	MotionPicture:         "MOTION_PICTURE",
	MusicalScore:          "MUSICAL_SCORE",
	Pamphlet:              "PAMPHLET",
	PaperConference:       "PAPER_CONFERENCE",
	Patent:                "PATENT",
	PersonalCommunication: "PERSONAL_COMMUNICATION",
	Post:                  "POST",
	PostWeblog:            "POST_WEBLOG",
	Report:                "REPORT",
	Review:                "REVIEW",
	ReviewBook:            "REVIEW_BOOK",
	Song:                  "SONG",
	Speech:                "SPEECH",
	Thesis:                "THESIS",
	Treaty:                "TREATY",
	Webpage:               "WEBPAGE",
}

var stringToReferenceType = func() map[string]ReferenceType {
	res := make(map[string]ReferenceType)
	for k, v := range refTypeToString {
		res[v] = k
	}
	return res
}()

// ID returns the string ID of ReferenceType.
func (rg ReferenceType) ID() string {
	if res, ok := refTypeToString[rg]; ok {
		return res
	}
	return ""
}

// String returns the string representation of ReferenceType.
func (rg ReferenceType) String() string {
	return ToStr(rg.ID())
}

func NewReferenceType(s string) ReferenceType {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, "", "_")
	if res, ok := stringToReferenceType[s]; ok {
		return res
	}
	return UnknownRT
}
