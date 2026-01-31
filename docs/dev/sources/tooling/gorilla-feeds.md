# gorilla/feeds

> Source: https://pkg.go.dev/github.com/gorilla/feeds
> Fetched: 2026-01-30T23:49:44.989344+00:00
> Content-Hash: 6c7ad9b1147478ba
> Type: html

---

Overview

¶

Examples

Syndication (feed) generator library for golang.

Installing

go get github.com/gorilla/feeds

Feeds provides a simple, generic Feed interface with a generic Item object as well as RSS, Atom and JSON Feed specific RssFeed, AtomFeed and JSONFeed objects which allow access to all of each spec's defined elements.

Examples

¶

Create a Feed and some Items in that feed using the generic interfaces:

import (
	"time"
	. "github.com/gorilla/feeds"
)

now = time.Now()

feed := &Feed{
	Title:       "jmoiron.net blog",
	Link:        &Link{Href: "http://jmoiron.net/blog"},
	Description: "discussion about tech, footie, photos",
	Author:      &Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
	Created:     now,
	Copyright:   "This work is copyright © Benjamin Button",
}

feed.Items = []*Item{
	&Item{
		Title:       "Limiting Concurrency in Go",
		Link:        &Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
		Description: "A discussion on controlled parallelism in golang",
		Author:      &Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
		Created:     now,
	},
	&Item{
		Title:       "Logic-less Template Redux",
		Link:        &Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
		Description: "More thoughts on logicless templates",
		Created:     now,
	},
	&Item{
		Title:       "Idiomatic Code Reuse in Go",
		Link:        &Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
		Description: "How to use interfaces <em>effectively</em>",
		Created:     now,
	},
}

From here, you can output Atom, RSS, or JSON Feed versions of this feed easily

atom, err := feed.ToAtom()
rss, err := feed.ToRss()
json, err := feed.ToJSON()

You can also get access to the underlying objects that feeds uses to export its XML

atomFeed := (&Atom{Feed: feed}).AtomFeed()
rssFeed := (&Rss{Feed: feed}).RssFeed()
jsonFeed := (&JSON{Feed: feed}).JSONFeed()

From here, you can modify or add each syndication's specific fields before outputting

atomFeed.Subtitle = "plays the blues"
atom, err := ToXML(atomFeed)
rssFeed.Generator = "gorilla/feeds v1.0 (github.com/gorilla/feeds)"
rss, err := ToXML(rssFeed)
jsonFeed.NextUrl = "https://www.example.com/feed.json?page=2"
json, err := jsonFeed.ToJSON()

Index

¶

func ToXML(feed XmlFeed) (string, error)

func WriteXML(feed XmlFeed, w io.Writer) error

type Atom

func (a *Atom) AtomFeed() *AtomFeed

func (a *Atom) FeedXml() interface{}

type AtomAuthor

type AtomContent

type AtomContributor

type AtomEntry

type AtomFeed

func (a *AtomFeed) FeedXml() interface{}

type AtomLink

type AtomPerson

type AtomSummary

type Author

type Enclosure

type Feed

func (f *Feed) Add(item *Item)

func (f *Feed) Sort(less func(a, b *Item) bool)

func (f *Feed) ToAtom() (string, error)

func (f *Feed) ToJSON() (string, error)

func (f *Feed) ToRss() (string, error)

func (f *Feed) WriteAtom(w io.Writer) error

func (f *Feed) WriteJSON(w io.Writer) error

func (f *Feed) WriteRss(w io.Writer) error

type Image

type Item

type JSON

func (f *JSON) JSONFeed() *JSONFeed

func (f *JSON) ToJSON() (string, error)

type JSONAttachment

func (a *JSONAttachment) MarshalJSON() ([]byte, error)

func (a *JSONAttachment) UnmarshalJSON(data []byte) error

type JSONAuthor

type JSONFeed

func (f *JSONFeed) ToJSON() (string, error)

type JSONHub

type JSONItem

type Link

type Rss

func (r *Rss) FeedXml() interface{}

func (r *Rss) RssFeed() *RssFeed

type RssContent

type RssEnclosure

type RssFeed

func (r *RssFeed) FeedXml() interface{}

type RssFeedXml

type RssGuid

type RssImage

type RssItem

type RssTextInput

type UUID

func NewUUID() *UUID

func (u *UUID) String() string

type XmlFeed

Constants

¶

This section is empty.

Variables

¶

This section is empty.

Functions

¶

func

ToXML

¶

func ToXML(feed

XmlFeed

) (

string

,

error

)

turn a feed object (either a Feed, AtomFeed, or RssFeed) into xml
returns an error if xml marshaling fails

func

WriteXML

¶

func WriteXML(feed

XmlFeed

, w

io

.

Writer

)

error

WriteXML writes a feed object (either a Feed, AtomFeed, or RssFeed) as XML into
the writer. Returns an error if XML marshaling fails.

Types

¶

type

Atom

¶

type Atom struct {

*

Feed

}

func (*Atom)

AtomFeed

¶

func (a *

Atom

) AtomFeed() *

AtomFeed

create a new AtomFeed with a generic Feed struct's data

func (*Atom)

FeedXml

¶

func (a *

Atom

) FeedXml() interface{}

FeedXml returns an XML-Ready object for an Atom object

type

AtomAuthor

¶

type AtomAuthor struct {

XMLName

xml

.

Name

`xml:"author"`

AtomPerson

}

type

AtomContent

¶

type AtomContent struct {

XMLName

xml

.

Name

`xml:"content"`

Content

string

`xml:",chardata"`

Type

string

`xml:"type,attr"`

}

type

AtomContributor

¶

type AtomContributor struct {

XMLName

xml

.

Name

`xml:"contributor"`

AtomPerson

}

type

AtomEntry

¶

type AtomEntry struct {

XMLName

xml

.

Name

`xml:"entry"`

Xmlns

string

`xml:"xmlns,attr,omitempty"`

Title

string

`xml:"title"`

// required

Updated

string

`xml:"updated"`

// required

Id

string

`xml:"id"`

// required

Category

string

`xml:"category,omitempty"`

Content     *

AtomContent

Rights

string

`xml:"rights,omitempty"`

Source

string

`xml:"source,omitempty"`

Published

string

`xml:"published,omitempty"`

Contributor *

AtomContributor

Links       []

AtomLink

// required if no child 'content' elements

Summary     *

AtomSummary

// required if content has src or content is base64

Author      *

AtomAuthor

// required if feed lacks an author

}

type

AtomFeed

¶

type AtomFeed struct {

XMLName

xml

.

Name

`xml:"feed"`

Xmlns

string

`xml:"xmlns,attr"`

Title

string

`xml:"title"`

// required

Id

string

`xml:"id"`

// required

Updated

string

`xml:"updated"`

// required

Category

string

`xml:"category,omitempty"`

Icon

string

`xml:"icon,omitempty"`

Logo

string

`xml:"logo,omitempty"`

Rights

string

`xml:"rights,omitempty"`

// copyright used

Subtitle

string

`xml:"subtitle,omitempty"`

Link        *

AtomLink

Author      *

AtomAuthor

`xml:"author,omitempty"`

Contributor *

AtomContributor

Entries     []*

AtomEntry

`xml:"entry"`

}

func (*AtomFeed)

FeedXml

¶

func (a *

AtomFeed

) FeedXml() interface{}

FeedXml returns an XML-ready object for an AtomFeed object

type

AtomLink

¶

type AtomLink struct {

//Atom 1.0 <link rel="enclosure" type="audio/mpeg" title="MP3" href="

http://www.example.org/myaudiofile.mp3

" length="1234" />

XMLName

xml

.

Name

`xml:"link"`

Href

string

`xml:"href,attr"`

Rel

string

`xml:"rel,attr,omitempty"`

Type

string

`xml:"type,attr,omitempty"`

Length

string

`xml:"length,attr,omitempty"`

}

Multiple links with different rel can coexist

type

AtomPerson

¶

type AtomPerson struct {

Name

string

`xml:"name,omitempty"`

Uri

string

`xml:"uri,omitempty"`

Email

string

`xml:"email,omitempty"`

}

type

AtomSummary

¶

type AtomSummary struct {

XMLName

xml

.

Name

`xml:"summary"`

Content

string

`xml:",chardata"`

Type

string

`xml:"type,attr"`

}

type

Author

¶

type Author struct {

Name, Email

string

}

type

Enclosure

¶

type Enclosure struct {

Url, Length, Type

string

}

type

Feed

¶

type Feed struct {

Title

string

Link        *

Link

Description

string

Author      *

Author

Updated

time

.

Time

Created

time

.

Time

Id

string

Subtitle

string

Items       []*

Item

Copyright

string

Image       *

Image

}

func (*Feed)

Add

¶

func (f *

Feed

) Add(item *

Item

)

add a new Item to a Feed

func (*Feed)

Sort

¶

added in

v1.1.1

func (f *

Feed

) Sort(less func(a, b *

Item

)

bool

)

Sort sorts the Items in the feed with the given less function.

func (*Feed)

ToAtom

¶

func (f *

Feed

) ToAtom() (

string

,

error

)

creates an Atom representation of this feed

func (*Feed)

ToJSON

¶

func (f *

Feed

) ToJSON() (

string

,

error

)

ToJSON creates a JSON Feed representation of this feed

func (*Feed)

ToRss

¶

func (f *

Feed

) ToRss() (

string

,

error

)

creates an Rss representation of this feed

func (*Feed)

WriteAtom

¶

func (f *

Feed

) WriteAtom(w

io

.

Writer

)

error

WriteAtom writes an Atom representation of this feed to the writer.

func (*Feed)

WriteJSON

¶

func (f *

Feed

) WriteJSON(w

io

.

Writer

)

error

WriteJSON writes an JSON representation of this feed to the writer.

func (*Feed)

WriteRss

¶

func (f *

Feed

) WriteRss(w

io

.

Writer

)

error

WriteRss writes an RSS representation of this feed to the writer.

type

Image

¶

type Image struct {

Url, Title, Link

string

Width, Height

int

}

type

Item

¶

type Item struct {

Title

string

Link        *

Link

Source      *

Link

Author      *

Author

Description

string

// used as description in rss, summary in atom

Id

string

// used as guid in rss, id in atom

IsPermaLink

string

// an optional parameter for guid in rss

Updated

time

.

Time

Created

time

.

Time

Enclosure   *

Enclosure

Content

string

}

type

JSON

¶

type JSON struct {

*

Feed

}

JSON is used to convert a generic Feed to a JSONFeed.

func (*JSON)

JSONFeed

¶

func (f *

JSON

) JSONFeed() *

JSONFeed

JSONFeed creates a new JSONFeed with a generic Feed struct's data.

func (*JSON)

ToJSON

¶

func (f *

JSON

) ToJSON() (

string

,

error

)

ToJSON encodes f into a JSON string. Returns an error if marshalling fails.

type

JSONAttachment

¶

type JSONAttachment struct {

Url

string

`json:"url,omitempty"`

MIMEType

string

`json:"mime_type,omitempty"`

Title

string

`json:"title,omitempty"`

Size

int32

`json:"size,omitempty"`

Duration

time

.

Duration

`json:"duration_in_seconds,omitempty"`

}

JSONAttachment represents a related resource. Podcasts, for instance, would
include an attachment that’s an audio or video file.

func (*JSONAttachment)

MarshalJSON

¶

func (a *

JSONAttachment

) MarshalJSON() ([]

byte

,

error

)

MarshalJSON implements the json.Marshaler interface.
The Duration field is marshaled in seconds, all other fields are marshaled
based upon the definitions in struct tags.

func (*JSONAttachment)

UnmarshalJSON

¶

func (a *

JSONAttachment

) UnmarshalJSON(data []

byte

)

error

UnmarshalJSON implements the json.Unmarshaler interface.
The Duration field is expected to be in seconds, all other field types
match the struct definition.

type

JSONAuthor

¶

type JSONAuthor struct {

Name

string

`json:"name,omitempty"`

Url

string

`json:"url,omitempty"`

Avatar

string

`json:"avatar,omitempty"`

}

JSONAuthor represents the author of the feed or of an individual item
in the feed

type

JSONFeed

¶

type JSONFeed struct {

Version

string

`json:"version"`

Title

string

`json:"title"`

Language

string

`json:"language,omitempty"`

HomePageUrl

string

`json:"home_page_url,omitempty"`

FeedUrl

string

`json:"feed_url,omitempty"`

Description

string

`json:"description,omitempty"`

UserComment

string

`json:"user_comment,omitempty"`

NextUrl

string

`json:"next_url,omitempty"`

Icon

string

`json:"icon,omitempty"`

Favicon

string

`json:"favicon,omitempty"`

Author      *

JSONAuthor

`json:"author,omitempty"`

// deprecated in JSON Feed v1.1, keeping for backwards compatibility

Authors     []*

JSONAuthor

`json:"authors,omitempty"`

Expired     *

bool

`json:"expired,omitempty"`

Hubs        []*

JSONHub

`json:"hubs,omitempty"`

Items       []*

JSONItem

`json:"items,omitempty"`

}

JSONFeed represents a syndication feed in the JSON Feed Version 1 format.
Matching the specification found here:

https://jsonfeed.org/version/1

.

func (*JSONFeed)

ToJSON

¶

func (f *

JSONFeed

) ToJSON() (

string

,

error

)

ToJSON encodes f into a JSON string. Returns an error if marshalling fails.

type

JSONHub

¶

type JSONHub struct {

Type

string

`json:"type"`

Url

string

`json:"url"`

}

JSONHub describes an endpoint that can be used to subscribe to real-time
notifications from the publisher of this feed.

type

JSONItem

¶

type JSONItem struct {

Id

string

`json:"id"`

Url

string

`json:"url,omitempty"`

ExternalUrl

string

`json:"external_url,omitempty"`

Title

string

`json:"title,omitempty"`

ContentHTML

string

`json:"content_html,omitempty"`

ContentText

string

`json:"content_text,omitempty"`

Summary

string

`json:"summary,omitempty"`

Image

string

`json:"image,omitempty"`

BannerImage

string

`json:"banner_,omitempty"`

PublishedDate *

time

.

Time

`json:"date_published,omitempty"`

ModifiedDate  *

time

.

Time

`json:"date_modified,omitempty"`

Author        *

JSONAuthor

`json:"author,omitempty"`

// deprecated in JSON Feed v1.1, keeping for backwards compatibility

Authors       []*

JSONAuthor

`json:"authors,omitempty"`

Tags          []

string

`json:"tags,omitempty"`

Attachments   []

JSONAttachment

`json:"attachments,omitempty"`

}

JSONItem represents a single entry/post for the feed.

type

Link

¶

type Link struct {

Href, Rel, Type, Length

string

}

type

Rss

¶

type Rss struct {

*

Feed

}

func (*Rss)

FeedXml

¶

func (r *

Rss

) FeedXml() interface{}

FeedXml returns an XML-Ready object for an Rss object

func (*Rss)

RssFeed

¶

func (r *

Rss

) RssFeed() *

RssFeed

create a new RssFeed with a generic Feed struct's data

type

RssContent

¶

added in

v1.1.0

type RssContent struct {

XMLName

xml

.

Name

`xml:"content:encoded"`

Content

string

`xml:",cdata"`

}

type

RssEnclosure

¶

type RssEnclosure struct {

//RSS 2.0 <enclosure url="

http://example.com/file.mp3

" length="123456789" type="audio/mpeg" />

XMLName

xml

.

Name

`xml:"enclosure"`

Url

string

`xml:"url,attr"`

Length

string

`xml:"length,attr"`

Type

string

`xml:"type,attr"`

}

type

RssFeed

¶

type RssFeed struct {

XMLName

xml

.

Name

`xml:"channel"`

Title

string

`xml:"title"`

// required

Link

string

`xml:"link"`

// required

Description

string

`xml:"description"`

// required

Language

string

`xml:"language,omitempty"`

Copyright

string

`xml:"copyright,omitempty"`

ManagingEditor

string

`xml:"managingEditor,omitempty"`

// Author used

WebMaster

string

`xml:"webMaster,omitempty"`

PubDate

string

`xml:"pubDate,omitempty"`

// created or updated

LastBuildDate

string

`xml:"lastBuildDate,omitempty"`

// updated used

Category

string

`xml:"category,omitempty"`

Generator

string

`xml:"generator,omitempty"`

Docs

string

`xml:"docs,omitempty"`

Cloud

string

`xml:"cloud,omitempty"`

Ttl

int

`xml:"ttl,omitempty"`

Rating

string

`xml:"rating,omitempty"`

SkipHours

string

`xml:"skipHours,omitempty"`

SkipDays

string

`xml:"skipDays,omitempty"`

Image          *

RssImage

TextInput      *

RssTextInput

Items          []*

RssItem

`xml:"item"`

}

func (*RssFeed)

FeedXml

¶

func (r *

RssFeed

) FeedXml() interface{}

FeedXml returns an XML-ready object for an RssFeed object

type

RssFeedXml

¶

added in

v1.1.1

type RssFeedXml struct {

XMLName

xml

.

Name

`xml:"rss"`

Version

string

`xml:"version,attr"`

ContentNamespace

string

`xml:"xmlns:content,attr"`

Channel          *

RssFeed

}

private wrapper around the RssFeed which gives us the <rss>..</rss> xml

type

RssGuid

¶

added in

v1.2.0

type RssGuid struct {

//RSS 2.0 <guid isPermaLink="true">

http://inessential.com/2002/09/01.php#a2

</guid>

XMLName

xml

.

Name

`xml:"guid"`

Id

string

`xml:",chardata"`

IsPermaLink

string

`xml:"isPermaLink,attr,omitempty"`

// "true", "false", or an empty string

}

type

RssImage

¶

type RssImage struct {

XMLName

xml

.

Name

`xml:"image"`

Url

string

`xml:"url"`

Title

string

`xml:"title"`

Link

string

`xml:"link"`

Width

int

`xml:"width,omitempty"`

Height

int

`xml:"height,omitempty"`

}

type

RssItem

¶

type RssItem struct {

XMLName

xml

.

Name

`xml:"item"`

Title

string

`xml:"title"`

// required

Link

string

`xml:"link"`

// required

Description

string

`xml:"description"`

// required

Content     *

RssContent

Author

string

`xml:"author,omitempty"`

Category

string

`xml:"category,omitempty"`

Comments

string

`xml:"comments,omitempty"`

Enclosure   *

RssEnclosure

Guid        *

RssGuid

// Id used

PubDate

string

`xml:"pubDate,omitempty"`

// created or updated

Source

string

`xml:"source,omitempty"`

}

type

RssTextInput

¶

type RssTextInput struct {

XMLName

xml

.

Name

`xml:"textInput"`

Title

string

`xml:"title"`

Description

string

`xml:"description"`

Name

string

`xml:"name"`

Link

string

`xml:"link"`

}

type

UUID

¶

type UUID [16]

byte

func

NewUUID

¶

func NewUUID() *

UUID

create a new uuid v4

func (*UUID)

String

¶

func (u *

UUID

) String()

string

type

XmlFeed

¶

type XmlFeed interface {

FeedXml() interface{}

}

interface used by ToXML to get a object suitable for exporting XML.