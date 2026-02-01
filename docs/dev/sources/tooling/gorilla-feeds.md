# gorilla/feeds

> Source: https://pkg.go.dev/github.com/gorilla/feeds
> Fetched: 2026-02-01T11:43:07.528820+00:00
> Content-Hash: e20ee4e159d4d619
> Type: html

---

### Overview ¶

  * Examples



Syndication (feed) generator library for golang. 

Installing 
    
    
    go get github.com/gorilla/feeds
    

Feeds provides a simple, generic Feed interface with a generic Item object as well as RSS, Atom and JSON Feed specific RssFeed, AtomFeed and JSONFeed objects which allow access to all of each spec's defined elements. 

#### Examples ¶

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
    

### Index ¶

  * func ToXML(feed XmlFeed) (string, error)
  * func WriteXML(feed XmlFeed, w io.Writer) error
  * type Atom
  *     * func (a *Atom) AtomFeed() *AtomFeed
    * func (a *Atom) FeedXml() interface{}
  * type AtomAuthor
  * type AtomContent
  * type AtomContributor
  * type AtomEntry
  * type AtomFeed
  *     * func (a *AtomFeed) FeedXml() interface{}
  * type AtomLink
  * type AtomPerson
  * type AtomSummary
  * type Author
  * type Enclosure
  * type Feed
  *     * func (f *Feed) Add(item *Item)
    * func (f *Feed) Sort(less func(a, b *Item) bool)
    * func (f *Feed) ToAtom() (string, error)
    * func (f *Feed) ToJSON() (string, error)
    * func (f *Feed) ToRss() (string, error)
    * func (f *Feed) WriteAtom(w io.Writer) error
    * func (f *Feed) WriteJSON(w io.Writer) error
    * func (f *Feed) WriteRss(w io.Writer) error
  * type Image
  * type Item
  * type JSON
  *     * func (f *JSON) JSONFeed() *JSONFeed
    * func (f *JSON) ToJSON() (string, error)
  * type JSONAttachment
  *     * func (a *JSONAttachment) MarshalJSON() ([]byte, error)
    * func (a *JSONAttachment) UnmarshalJSON(data []byte) error
  * type JSONAuthor
  * type JSONFeed
  *     * func (f *JSONFeed) ToJSON() (string, error)
  * type JSONHub
  * type JSONItem
  * type Link
  * type Rss
  *     * func (r *Rss) FeedXml() interface{}
    * func (r *Rss) RssFeed() *RssFeed
  * type RssContent
  * type RssEnclosure
  * type RssFeed
  *     * func (r *RssFeed) FeedXml() interface{}
  * type RssFeedXml
  * type RssGuid
  * type RssImage
  * type RssItem
  * type RssTextInput
  * type UUID
  *     * func NewUUID() *UUID
  *     * func (u *UUID) String() string
  * type XmlFeed



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

####  func [ToXML](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L78) ¶
    
    
    func ToXML(feed XmlFeed) ([string](/builtin#string), [error](/builtin#error))

turn a feed object (either a Feed, AtomFeed, or RssFeed) into xml returns an error if xml marshaling fails 

####  func [WriteXML](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L91) ¶
    
    
    func WriteXML(feed XmlFeed, w [io](/io).[Writer](/io#Writer)) [error](/builtin#error)

WriteXML writes a feed object (either a Feed, AtomFeed, or RssFeed) as XML into the writer. Returns an error if XML marshaling fails. 

### Types ¶

####  type [Atom](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L86) ¶
    
    
    type Atom struct {
    	*Feed
    }

####  func (*Atom) [AtomFeed](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L146) ¶
    
    
    func (a *Atom) AtomFeed() *AtomFeed

create a new AtomFeed with a generic Feed struct's data 

####  func (*Atom) [FeedXml](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L171) ¶
    
    
    func (a *Atom) FeedXml() interface{}

FeedXml returns an XML-Ready object for an Atom object 

####  type [AtomAuthor](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L32) ¶
    
    
    type AtomAuthor struct {
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"author"`
    	AtomPerson
    }

####  type [AtomContent](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L26) ¶
    
    
    type AtomContent struct {
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"content"`
    	Content [string](/builtin#string)   `xml:",chardata"`
    	Type    [string](/builtin#string)   `xml:"type,attr"`
    }

####  type [AtomContributor](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L37) ¶
    
    
    type AtomContributor struct {
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"contributor"`
    	AtomPerson
    }

####  type [AtomEntry](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L42) ¶
    
    
    type AtomEntry struct {
    	XMLName     [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"entry"`
    	Xmlns       [string](/builtin#string)   `xml:"xmlns,attr,omitempty"`
    	Title       [string](/builtin#string)   `xml:"title"`   // required
    	Updated     [string](/builtin#string)   `xml:"updated"` // required
    	Id          [string](/builtin#string)   `xml:"id"`      // required
    	Category    [string](/builtin#string)   `xml:"category,omitempty"`
    	Content     *AtomContent
    	Rights      [string](/builtin#string) `xml:"rights,omitempty"`
    	Source      [string](/builtin#string) `xml:"source,omitempty"`
    	Published   [string](/builtin#string) `xml:"published,omitempty"`
    	Contributor *AtomContributor
    	Links       []AtomLink   // required if no child 'content' elements
    	Summary     *AtomSummary // required if content has src or content is base64
    	Author      *AtomAuthor  // required if feed lacks an author
    }

####  type [AtomFeed](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L69) ¶
    
    
    type AtomFeed struct {
    	XMLName     [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"feed"`
    	Xmlns       [string](/builtin#string)   `xml:"xmlns,attr"`
    	Title       [string](/builtin#string)   `xml:"title"`   // required
    	Id          [string](/builtin#string)   `xml:"id"`      // required
    	Updated     [string](/builtin#string)   `xml:"updated"` // required
    	Category    [string](/builtin#string)   `xml:"category,omitempty"`
    	Icon        [string](/builtin#string)   `xml:"icon,omitempty"`
    	Logo        [string](/builtin#string)   `xml:"logo,omitempty"`
    	Rights      [string](/builtin#string)   `xml:"rights,omitempty"` // copyright used
    	Subtitle    [string](/builtin#string)   `xml:"subtitle,omitempty"`
    	Link        *AtomLink
    	Author      *AtomAuthor `xml:"author,omitempty"`
    	Contributor *AtomContributor
    	Entries     []*AtomEntry `xml:"entry"`
    }

####  func (*AtomFeed) [FeedXml](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L176) ¶
    
    
    func (a *AtomFeed) FeedXml() interface{}

FeedXml returns an XML-ready object for an AtomFeed object 

####  type [AtomLink](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L60) ¶
    
    
    type AtomLink struct {
    	//Atom 1.0 <link rel="enclosure" type="audio/mpeg" title="MP3" href="<http://www.example.org/myaudiofile.mp3>" length="1234" />
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"link"`
    	Href    [string](/builtin#string)   `xml:"href,attr"`
    	Rel     [string](/builtin#string)   `xml:"rel,attr,omitempty"`
    	Type    [string](/builtin#string)   `xml:"type,attr,omitempty"`
    	Length  [string](/builtin#string)   `xml:"length,attr,omitempty"`
    }

Multiple links with different rel can coexist 

####  type [AtomPerson](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L14) ¶
    
    
    type AtomPerson struct {
    	Name  [string](/builtin#string) `xml:"name,omitempty"`
    	Uri   [string](/builtin#string) `xml:"uri,omitempty"`
    	Email [string](/builtin#string) `xml:"email,omitempty"`
    }

####  type [AtomSummary](https://github.com/gorilla/feeds/blob/v1.2.0/atom.go#L20) ¶
    
    
    type AtomSummary struct {
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"summary"`
    	Content [string](/builtin#string)   `xml:",chardata"`
    	Type    [string](/builtin#string)   `xml:"type,attr"`
    }

####  type [Author](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L15) ¶
    
    
    type Author struct {
    	Name, Email [string](/builtin#string)
    }

####  type [Enclosure](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L24) ¶
    
    
    type Enclosure struct {
    	Url, Length, Type [string](/builtin#string)
    }

####  type [Feed](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L42) ¶
    
    
    type Feed struct {
    	Title       [string](/builtin#string)
    	Link        *Link
    	Description [string](/builtin#string)
    	Author      *Author
    	Updated     [time](/time).[Time](/time#Time)
    	Created     [time](/time).[Time](/time#Time)
    	Id          [string](/builtin#string)
    	Subtitle    [string](/builtin#string)
    	Items       []*Item
    	Copyright   [string](/builtin#string)
    	Image       *Image
    }

####  func (*Feed) [Add](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L57) ¶
    
    
    func (f *Feed) Add(item *Item)

add a new Item to a Feed 

####  func (*Feed) [Sort](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L141) ¶ added in v1.1.1
    
    
    func (f *Feed) Sort(less func(a, b *Item) [bool](/builtin#bool))

Sort sorts the Items in the feed with the given less function. 

####  func (*Feed) [ToAtom](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L103) ¶
    
    
    func (f *Feed) ToAtom() ([string](/builtin#string), [error](/builtin#error))

creates an Atom representation of this feed 

####  func (*Feed) [ToJSON](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L125) ¶
    
    
    func (f *Feed) ToJSON() ([string](/builtin#string), [error](/builtin#error))

ToJSON creates a JSON Feed representation of this feed 

####  func (*Feed) [ToRss](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L114) ¶
    
    
    func (f *Feed) ToRss() ([string](/builtin#string), [error](/builtin#error))

creates an Rss representation of this feed 

####  func (*Feed) [WriteAtom](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L109) ¶
    
    
    func (f *Feed) WriteAtom(w [io](/io).[Writer](/io#Writer)) [error](/builtin#error)

WriteAtom writes an Atom representation of this feed to the writer. 

####  func (*Feed) [WriteJSON](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L131) ¶
    
    
    func (f *Feed) WriteJSON(w [io](/io).[Writer](/io#Writer)) [error](/builtin#error)

WriteJSON writes an JSON representation of this feed to the writer. 

####  func (*Feed) [WriteRss](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L120) ¶
    
    
    func (f *Feed) WriteRss(w [io](/io).[Writer](/io#Writer)) [error](/builtin#error)

WriteRss writes an RSS representation of this feed to the writer. 

####  type [Image](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L19) ¶
    
    
    type Image struct {
    	Url, Title, Link [string](/builtin#string)
    	Width, Height    [int](/builtin#int)
    }

####  type [Item](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L28) ¶
    
    
    type Item struct {
    	Title       [string](/builtin#string)
    	Link        *Link
    	Source      *Link
    	Author      *Author
    	Description [string](/builtin#string) // used as description in rss, summary in atom
    	Id          [string](/builtin#string) // used as guid in rss, id in atom
    	IsPermaLink [string](/builtin#string) // an optional parameter for guid in rss
    	Updated     [time](/time).[Time](/time#Time)
    	Created     [time](/time).[Time](/time#Time)
    	Enclosure   *Enclosure
    	Content     [string](/builtin#string)
    }

####  type [JSON](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L114) ¶
    
    
    type JSON struct {
    	*Feed
    }

JSON is used to convert a generic Feed to a JSONFeed. 

####  func (*JSON) [JSONFeed](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L134) ¶
    
    
    func (f *JSON) JSONFeed() *JSONFeed

JSONFeed creates a new JSONFeed with a generic Feed struct's data. 

####  func (*JSON) [ToJSON](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L119) ¶
    
    
    func (f *JSON) ToJSON() ([string](/builtin#string), [error](/builtin#error))

ToJSON encodes f into a JSON string. Returns an error if marshalling fails. 

####  type [JSONAttachment](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L21) ¶
    
    
    type JSONAttachment struct {
    	Url      [string](/builtin#string)        `json:"url,omitempty"`
    	MIMEType [string](/builtin#string)        `json:"mime_type,omitempty"`
    	Title    [string](/builtin#string)        `json:"title,omitempty"`
    	Size     [int32](/builtin#int32)         `json:"size,omitempty"`
    	Duration [time](/time).[Duration](/time#Duration) `json:"duration_in_seconds,omitempty"`
    }

JSONAttachment represents a related resource. Podcasts, for instance, would include an attachment that’s an audio or video file. 

####  func (*JSONAttachment) [MarshalJSON](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L32) ¶
    
    
    func (a *JSONAttachment) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON implements the json.Marshaler interface. The Duration field is marshaled in seconds, all other fields are marshaled based upon the definitions in struct tags. 

####  func (*JSONAttachment) [UnmarshalJSON](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L46) ¶
    
    
    func (a *JSONAttachment) UnmarshalJSON(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON implements the json.Unmarshaler interface. The Duration field is expected to be in seconds, all other field types match the struct definition. 

####  type [JSONAuthor](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L13) ¶
    
    
    type JSONAuthor struct {
    	Name   [string](/builtin#string) `json:"name,omitempty"`
    	Url    [string](/builtin#string) `json:"url,omitempty"`
    	Avatar [string](/builtin#string) `json:"avatar,omitempty"`
    }

JSONAuthor represents the author of the feed or of an individual item in the feed 

####  type [JSONFeed](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L95) ¶
    
    
    type JSONFeed struct {
    	Version     [string](/builtin#string)        `json:"version"`
    	Title       [string](/builtin#string)        `json:"title"`
    	Language    [string](/builtin#string)        `json:"language,omitempty"`
    	HomePageUrl [string](/builtin#string)        `json:"home_page_url,omitempty"`
    	FeedUrl     [string](/builtin#string)        `json:"feed_url,omitempty"`
    	Description [string](/builtin#string)        `json:"description,omitempty"`
    	UserComment [string](/builtin#string)        `json:"user_comment,omitempty"`
    	NextUrl     [string](/builtin#string)        `json:"next_url,omitempty"`
    	Icon        [string](/builtin#string)        `json:"icon,omitempty"`
    	Favicon     [string](/builtin#string)        `json:"favicon,omitempty"`
    	Author      *JSONAuthor   `json:"author,omitempty"` // deprecated in JSON Feed v1.1, keeping for backwards compatibility
    	Authors     []*JSONAuthor `json:"authors,omitempty"`
    	Expired     *[bool](/builtin#bool)         `json:"expired,omitempty"`
    	Hubs        []*JSONHub    `json:"hubs,omitempty"`
    	Items       []*JSONItem   `json:"items,omitempty"`
    }

JSONFeed represents a syndication feed in the JSON Feed Version 1 format. Matching the specification found here: <https://jsonfeed.org/version/1>. 

####  func (*JSONFeed) [ToJSON](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L124) ¶
    
    
    func (f *JSONFeed) ToJSON() ([string](/builtin#string), [error](/builtin#error))

ToJSON encodes f into a JSON string. Returns an error if marshalling fails. 

####  type [JSONHub](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L88) ¶
    
    
    type JSONHub struct {
    	Type [string](/builtin#string) `json:"type"`
    	Url  [string](/builtin#string) `json:"url"`
    }

JSONHub describes an endpoint that can be used to subscribe to real-time notifications from the publisher of this feed. 

####  type [JSONItem](https://github.com/gorilla/feeds/blob/v1.2.0/json.go#L68) ¶
    
    
    type JSONItem struct {
    	Id            [string](/builtin#string)           `json:"id"`
    	Url           [string](/builtin#string)           `json:"url,omitempty"`
    	ExternalUrl   [string](/builtin#string)           `json:"external_url,omitempty"`
    	Title         [string](/builtin#string)           `json:"title,omitempty"`
    	ContentHTML   [string](/builtin#string)           `json:"content_html,omitempty"`
    	ContentText   [string](/builtin#string)           `json:"content_text,omitempty"`
    	Summary       [string](/builtin#string)           `json:"summary,omitempty"`
    	Image         [string](/builtin#string)           `json:"image,omitempty"`
    	BannerImage   [string](/builtin#string)           `json:"banner_,omitempty"`
    	PublishedDate *[time](/time).[Time](/time#Time)       `json:"date_published,omitempty"`
    	ModifiedDate  *[time](/time).[Time](/time#Time)       `json:"date_modified,omitempty"`
    	Author        *JSONAuthor      `json:"author,omitempty"` // deprecated in JSON Feed v1.1, keeping for backwards compatibility
    	Authors       []*JSONAuthor    `json:"authors,omitempty"`
    	Tags          [][string](/builtin#string)         `json:"tags,omitempty"`
    	Attachments   []JSONAttachment `json:"attachments,omitempty"`
    }

JSONItem represents a single entry/post for the feed. 

####  type [Link](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L11) ¶
    
    
    type Link struct {
    	Href, Rel, Type, Length [string](/builtin#string)
    }

####  type [Rss](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L97) ¶
    
    
    type Rss struct {
    	*Feed
    }

####  func (*Rss) [FeedXml](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L170) ¶
    
    
    func (r *Rss) FeedXml() interface{}

FeedXml returns an XML-Ready object for an Rss object 

####  func (*Rss) [RssFeed](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L133) ¶
    
    
    func (r *Rss) RssFeed() *RssFeed

create a new RssFeed with a generic Feed struct's data 

####  type [RssContent](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L21) ¶ added in v1.1.0
    
    
    type RssContent struct {
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"content:encoded"`
    	Content [string](/builtin#string)   `xml:",cdata"`
    }

####  type [RssEnclosure](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L82) ¶
    
    
    type RssEnclosure struct {
    	//RSS 2.0 <enclosure url="<http://example.com/file.mp3>" length="123456789" type="audio/mpeg" />
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"enclosure"`
    	Url     [string](/builtin#string)   `xml:"url,attr"`
    	Length  [string](/builtin#string)   `xml:"length,attr"`
    	Type    [string](/builtin#string)   `xml:"type,attr"`
    }

####  type [RssFeed](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L43) ¶
    
    
    type RssFeed struct {
    	XMLName        [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"channel"`
    	Title          [string](/builtin#string)   `xml:"title"`       // required
    	Link           [string](/builtin#string)   `xml:"link"`        // required
    	Description    [string](/builtin#string)   `xml:"description"` // required
    	Language       [string](/builtin#string)   `xml:"language,omitempty"`
    	Copyright      [string](/builtin#string)   `xml:"copyright,omitempty"`
    	ManagingEditor [string](/builtin#string)   `xml:"managingEditor,omitempty"` // Author used
    	WebMaster      [string](/builtin#string)   `xml:"webMaster,omitempty"`
    	PubDate        [string](/builtin#string)   `xml:"pubDate,omitempty"`       // created or updated
    	LastBuildDate  [string](/builtin#string)   `xml:"lastBuildDate,omitempty"` // updated used
    	Category       [string](/builtin#string)   `xml:"category,omitempty"`
    	Generator      [string](/builtin#string)   `xml:"generator,omitempty"`
    	Docs           [string](/builtin#string)   `xml:"docs,omitempty"`
    	Cloud          [string](/builtin#string)   `xml:"cloud,omitempty"`
    	Ttl            [int](/builtin#int)      `xml:"ttl,omitempty"`
    	Rating         [string](/builtin#string)   `xml:"rating,omitempty"`
    	SkipHours      [string](/builtin#string)   `xml:"skipHours,omitempty"`
    	SkipDays       [string](/builtin#string)   `xml:"skipDays,omitempty"`
    	Image          *RssImage
    	TextInput      *RssTextInput
    	Items          []*RssItem `xml:"item"`
    }

####  func (*RssFeed) [FeedXml](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L177) ¶
    
    
    func (r *RssFeed) FeedXml() interface{}

FeedXml returns an XML-ready object for an RssFeed object 

####  type [RssFeedXml](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L14) ¶ added in v1.1.1
    
    
    type RssFeedXml struct {
    	XMLName          [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"rss"`
    	Version          [string](/builtin#string)   `xml:"version,attr"`
    	ContentNamespace [string](/builtin#string)   `xml:"xmlns:content,attr"`
    	Channel          *RssFeed
    }

private wrapper around the RssFeed which gives us the <rss>..</rss> xml 

####  type [RssGuid](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L90) ¶ added in v1.2.0
    
    
    type RssGuid struct {
    	//RSS 2.0 <guid isPermaLink="true"><http://inessential.com/2002/09/01.php#a2></guid>
    	XMLName     [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"guid"`
    	Id          [string](/builtin#string)   `xml:",chardata"`
    	IsPermaLink [string](/builtin#string)   `xml:"isPermaLink,attr,omitempty"` // "true", "false", or an empty string
    }

####  type [RssImage](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L26) ¶
    
    
    type RssImage struct {
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"image"`
    	Url     [string](/builtin#string)   `xml:"url"`
    	Title   [string](/builtin#string)   `xml:"title"`
    	Link    [string](/builtin#string)   `xml:"link"`
    	Width   [int](/builtin#int)      `xml:"width,omitempty"`
    	Height  [int](/builtin#int)      `xml:"height,omitempty"`
    }

####  type [RssItem](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L67) ¶
    
    
    type RssItem struct {
    	XMLName     [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"item"`
    	Title       [string](/builtin#string)   `xml:"title"`       // required
    	Link        [string](/builtin#string)   `xml:"link"`        // required
    	Description [string](/builtin#string)   `xml:"description"` // required
    	Content     *RssContent
    	Author      [string](/builtin#string) `xml:"author,omitempty"`
    	Category    [string](/builtin#string) `xml:"category,omitempty"`
    	Comments    [string](/builtin#string) `xml:"comments,omitempty"`
    	Enclosure   *RssEnclosure
    	Guid        *RssGuid // Id used
    	PubDate     [string](/builtin#string)   `xml:"pubDate,omitempty"` // created or updated
    	Source      [string](/builtin#string)   `xml:"source,omitempty"`
    }

####  type [RssTextInput](https://github.com/gorilla/feeds/blob/v1.2.0/rss.go#L35) ¶
    
    
    type RssTextInput struct {
    	XMLName     [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"textInput"`
    	Title       [string](/builtin#string)   `xml:"title"`
    	Description [string](/builtin#string)   `xml:"description"`
    	Name        [string](/builtin#string)   `xml:"name"`
    	Link        [string](/builtin#string)   `xml:"link"`
    }

####  type [UUID](https://github.com/gorilla/feeds/blob/v1.2.0/uuid.go#L10) ¶
    
    
    type UUID [16][byte](/builtin#byte)

####  func [NewUUID](https://github.com/gorilla/feeds/blob/v1.2.0/uuid.go#L13) ¶
    
    
    func NewUUID() *UUID

create a new uuid v4 

####  func (*UUID) [String](https://github.com/gorilla/feeds/blob/v1.2.0/uuid.go#L25) ¶
    
    
    func (u *UUID) String() [string](/builtin#string)

####  type [XmlFeed](https://github.com/gorilla/feeds/blob/v1.2.0/feed.go#L72) ¶
    
    
    type XmlFeed interface {
    	FeedXml() interface{}
    }

interface used by ToXML to get a object suitable for exporting XML. 
