# mmcdole/gofeed

> Source: https://pkg.go.dev/github.com/mmcdole/gofeed
> Fetched: 2026-01-30T23:49:40.442797+00:00
> Content-Hash: 2bcbad7bb28ba649
> Type: html

---

Index

¶

Variables

type Auth

type DefaultAtomTranslator

func (t *DefaultAtomTranslator) Translate(feed interface{}) (*Feed, error)

type DefaultJSONTranslator

func (t *DefaultJSONTranslator) Translate(feed interface{}) (*Feed, error)

type DefaultRSSTranslator

func (t *DefaultRSSTranslator) Translate(feed interface{}) (*Feed, error)

type Enclosure

type Feed

func (f Feed) Len() int

func (f Feed) Less(i, k int) bool

func (f Feed) String() string

func (f Feed) Swap(i, k int)

type FeedType

func DetectFeedType(feed io.Reader) FeedType

type HTTPError

func (err HTTPError) Error() string

type Image

type Item

type Parser

func NewParser() *Parser

func (f *Parser) Parse(feed io.Reader) (*Feed, error)

func (f *Parser) ParseString(feed string) (*Feed, error)

func (f *Parser) ParseURL(feedURL string) (feed *Feed, err error)

func (f *Parser) ParseURLWithContext(feedURL string, ctx context.Context) (feed *Feed, err error)

type Person

type Translator

Examples

¶

DetectFeedType

Parser.Parse

Parser.ParseString

Parser.ParseURL

Constants

¶

This section is empty.

Variables

¶

View Source

var ErrFeedTypeNotDetected =

errors

.

New

("Failed to detect feed type")

ErrFeedTypeNotDetected is returned when the detection system can not figure
out the Feed format

Functions

¶

This section is empty.

Types

¶

type

Auth

¶

added in

v1.2.0

type Auth struct {

Username

string

Password

string

}

Auth is a structure allowing to
use the BasicAuth during the HTTP request
It must be instantiated with your new Parser

type

DefaultAtomTranslator

¶

type DefaultAtomTranslator struct{}

DefaultAtomTranslator converts an atom.Feed struct
into the generic Feed struct.

This default implementation defines a set of
mapping rules between atom.Feed -> Feed
for each of the fields in Feed.

func (*DefaultAtomTranslator)

Translate

¶

func (t *

DefaultAtomTranslator

) Translate(feed interface{}) (*

Feed

,

error

)

Translate converts an Atom feed into the universal
feed type.

type

DefaultJSONTranslator

¶

added in

v1.1.0

type DefaultJSONTranslator struct{}

DefaultJSONTranslator converts an json.Feed struct
into the generic Feed struct.

This default implementation defines a set of
mapping rules between json.Feed -> Feed
for each of the fields in Feed.

func (*DefaultJSONTranslator)

Translate

¶

added in

v1.1.0

func (t *

DefaultJSONTranslator

) Translate(feed interface{}) (*

Feed

,

error

)

Translate converts an JSON feed into the universal
feed type.

type

DefaultRSSTranslator

¶

type DefaultRSSTranslator struct{}

DefaultRSSTranslator converts an rss.Feed struct
into the generic Feed struct.

This default implementation defines a set of
mapping rules between rss.Feed -> Feed
for each of the fields in Feed.

func (*DefaultRSSTranslator)

Translate

¶

func (t *

DefaultRSSTranslator

) Translate(feed interface{}) (*

Feed

,

error

)

Translate converts an RSS feed into the universal
feed type.

type

Enclosure

¶

type Enclosure struct {

URL

string

`json:"url,omitempty"`

Length

string

`json:"length,omitempty"`

Type

string

`json:"type,omitempty"`

}

Enclosure is a file associated with a given Item.

type

Feed

¶

type Feed struct {

Title

string

`json:"title,omitempty"`

Description

string

`json:"description,omitempty"`

Link

string

`json:"link,omitempty"`

FeedLink

string

`json:"feedLink,omitempty"`

Links           []

string

`json:"links,omitempty"`

Updated

string

`json:"updated,omitempty"`

UpdatedParsed   *

time

.

Time

`json:"updatedParsed,omitempty"`

Published

string

`json:"published,omitempty"`

PublishedParsed *

time

.

Time

`json:"publishedParsed,omitempty"`

Author          *

Person

`json:"author,omitempty"`

// Deprecated: Use feed.Authors instead

Authors         []*

Person

`json:"authors,omitempty"`

Language

string

`json:"language,omitempty"`

Image           *

Image

`json:"image,omitempty"`

Copyright

string

`json:"copyright,omitempty"`

Generator

string

`json:"generator,omitempty"`

Categories      []

string

`json:"categories,omitempty"`

DublinCoreExt   *

ext

.

DublinCoreExtension

`json:"dcExt,omitempty"`

ITunesExt       *

ext

.

ITunesFeedExtension

`json:"itunesExt,omitempty"`

Extensions

ext

.

Extensions

`json:"extensions,omitempty"`

Custom          map[

string

]

string

`json:"custom,omitempty"`

Items           []*

Item

`json:"items"`

FeedType

string

`json:"feedType"`

FeedVersion

string

`json:"feedVersion"`

}

Feed is the universal Feed type that atom.Feed
and rss.Feed gets translated to. It represents
a web feed.
Sorting with sort.Sort will order the Items by
oldest to newest publish time.

func (Feed)

Len

¶

func (f

Feed

) Len()

int

Len returns the length of Items.

func (Feed)

Less

¶

func (f

Feed

) Less(i, k

int

)

bool

Less compares PublishedParsed of Items[i], Items[k]
and returns true if Items[i] is less than Items[k].

func (Feed)

String

¶

func (f

Feed

) String()

string

func (Feed)

Swap

¶

func (f

Feed

) Swap(i, k

int

)

Swap swaps Items[i] and Items[k].

type

FeedType

¶

type FeedType

int

FeedType represents one of the possible feed
types that we can detect.

const (

// FeedTypeUnknown represents a feed that could not have its

// type determiend.

FeedTypeUnknown

FeedType

=

iota

// FeedTypeAtom repesents an Atom feed

FeedTypeAtom

// FeedTypeRSS represents an RSS feed

FeedTypeRSS

// FeedTypeJSON represents a JSON feed

FeedTypeJSON
)

func

DetectFeedType

¶

func DetectFeedType(feed

io

.

Reader

)

FeedType

DetectFeedType attempts to determine the type of feed
by looking for specific xml elements unique to the
various feed types.

Example

¶

package main

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

func main() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	feedType := gofeed.DetectFeedType(strings.NewReader(feedData))
	if feedType == gofeed.FeedTypeRSS {
		fmt.Println("Wow! This is an RSS feed!")
	}
}

Share

Format

Run

type

HTTPError

¶

type HTTPError struct {

StatusCode

int

Status

string

}

HTTPError represents an HTTP error returned by a server.

func (HTTPError)

Error

¶

func (err

HTTPError

) Error()

string

type

Image

¶

type Image struct {

URL

string

`json:"url,omitempty"`

Title

string

`json:"title,omitempty"`

}

Image is an image that is the artwork for a given
feed or item.

type

Item

¶

type Item struct {

Title

string

`json:"title,omitempty"`

Description

string

`json:"description,omitempty"`

Content

string

`json:"content,omitempty"`

Link

string

`json:"link,omitempty"`

Links           []

string

`json:"links,omitempty"`

Updated

string

`json:"updated,omitempty"`

UpdatedParsed   *

time

.

Time

`json:"updatedParsed,omitempty"`

Published

string

`json:"published,omitempty"`

PublishedParsed *

time

.

Time

`json:"publishedParsed,omitempty"`

Author          *

Person

`json:"author,omitempty"`

// Deprecated: Use item.Authors instead

Authors         []*

Person

`json:"authors,omitempty"`

GUID

string

`json:"guid,omitempty"`

Image           *

Image

`json:"image,omitempty"`

Categories      []

string

`json:"categories,omitempty"`

Enclosures      []*

Enclosure

`json:"enclosures,omitempty"`

DublinCoreExt   *

ext

.

DublinCoreExtension

`json:"dcExt,omitempty"`

ITunesExt       *

ext

.

ITunesItemExtension

`json:"itunesExt,omitempty"`

Extensions

ext

.

Extensions

`json:"extensions,omitempty"`

Custom          map[

string

]

string

`json:"custom,omitempty"`

}

Item is the universal Item type that atom.Entry
and rss.Item gets translated to.  It represents
a single entry in a given feed.

type

Parser

¶

type Parser struct {

AtomTranslator

Translator

RSSTranslator

Translator

JSONTranslator

Translator

UserAgent

string

AuthConfig     *

Auth

Client         *

http

.

Client

// contains filtered or unexported fields

}

Parser is a universal feed parser that detects
a given feed type, parsers it, and translates it
to the universal feed type.

func

NewParser

¶

func NewParser() *

Parser

NewParser creates a universal feed parser.

func (*Parser)

Parse

¶

func (f *

Parser

) Parse(feed

io

.

Reader

) (*

Feed

,

error

)

Parse parses a RSS or Atom or JSON feed into
the universal gofeed.Feed.  It takes an
io.Reader which should return the xml/json content.

Example

¶

package main

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

func main() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	fp := gofeed.NewParser()
	feed, err := fp.Parse(strings.NewReader(feedData))
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}

Share

Format

Run

func (*Parser)

ParseString

¶

func (f *

Parser

) ParseString(feed

string

) (*

Feed

,

error

)

ParseString parses a feed XML string and into the
universal feed type.

Example

¶

package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(feedData)
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}

Share

Format

Run

func (*Parser)

ParseURL

¶

func (f *

Parser

) ParseURL(feedURL

string

) (feed *

Feed

, err

error

)

ParseURL fetches the contents of a given url and
attempts to parse the response into the universal feed type.

Example

¶

package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("http://feeds.twit.tv/twit.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}

Share

Format

Run

func (*Parser)

ParseURLWithContext

¶

func (f *

Parser

) ParseURLWithContext(feedURL

string

, ctx

context

.

Context

) (feed *

Feed

, err

error

)

ParseURLWithContext fetches contents of a given url and
attempts to parse the response into the universal feed type.
You can instantiate the Auth structure with your Username and Password
to use the BasicAuth during the HTTP call.
It will be automatically added to the header of the request
Request could be canceled or timeout via given context

type

Person

¶

type Person struct {

Name

string

`json:"name,omitempty"`

Email

string

`json:"email,omitempty"`

}

Person is an individual specified in a feed
(e.g. an author)

type

Translator

¶

type Translator interface {

Translate(feed interface{}) (*

Feed

,

error

)

}

Translator converts a particular feed (atom.Feed or rss.Feed of json.Feed)
into the generic Feed struct