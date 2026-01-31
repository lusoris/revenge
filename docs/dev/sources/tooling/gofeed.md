# mmcdole/gofeed

> Source: https://pkg.go.dev/github.com/mmcdole/gofeed
> Fetched: 2026-01-31T10:57:00.223324+00:00
> Content-Hash: fec402a0a93e1299
> Type: html

---

### Index ¶

  * Variables
  * type Auth
  * type DefaultAtomTranslator
  *     * func (t *DefaultAtomTranslator) Translate(feed interface{}) (*Feed, error)
  * type DefaultJSONTranslator
  *     * func (t *DefaultJSONTranslator) Translate(feed interface{}) (*Feed, error)
  * type DefaultRSSTranslator
  *     * func (t *DefaultRSSTranslator) Translate(feed interface{}) (*Feed, error)
  * type Enclosure
  * type Feed
  *     * func (f Feed) Len() int
    * func (f Feed) Less(i, k int) bool
    * func (f Feed) String() string
    * func (f Feed) Swap(i, k int)
  * type FeedType
  *     * func DetectFeedType(feed io.Reader) FeedType
  * type HTTPError
  *     * func (err HTTPError) Error() string
  * type Image
  * type Item
  * type Parser
  *     * func NewParser() *Parser
  *     * func (f *Parser) Parse(feed io.Reader) (*Feed, error)
    * func (f *Parser) ParseString(feed string) (*Feed, error)
    * func (f *Parser) ParseURL(feedURL string) (feed *Feed, err error)
    * func (f *Parser) ParseURLWithContext(feedURL string, ctx context.Context) (feed *Feed, err error)
  * type Person
  * type Translator



### Examples ¶

  * DetectFeedType
  * Parser.Parse
  * Parser.ParseString
  * Parser.ParseURL



### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L19)
    
    
    var ErrFeedTypeNotDetected = [errors](/errors).[New](/errors#New)("Failed to detect feed type")

ErrFeedTypeNotDetected is returned when the detection system can not figure out the Feed format 

### Functions ¶

This section is empty.

### Types ¶

####  type [Auth](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L49) ¶ added in v1.2.0
    
    
    type Auth struct {
    	Username [string](/builtin#string)
    	Password [string](/builtin#string)
    }

Auth is a structure allowing to use the BasicAuth during the HTTP request It must be instantiated with your new Parser 

####  type [DefaultAtomTranslator](https://github.com/mmcdole/gofeed/blob/v1.3.0/translator.go#L535) ¶
    
    
    type DefaultAtomTranslator struct{}

DefaultAtomTranslator converts an atom.Feed struct into the generic Feed struct. 

This default implementation defines a set of mapping rules between atom.Feed -> Feed for each of the fields in Feed. 

####  func (*DefaultAtomTranslator) [Translate](https://github.com/mmcdole/gofeed/blob/v1.3.0/translator.go#L539) ¶
    
    
    func (t *DefaultAtomTranslator) Translate(feed interface{}) (*Feed, [error](/builtin#error))

Translate converts an Atom feed into the universal feed type. 

####  type [DefaultJSONTranslator](https://github.com/mmcdole/gofeed/blob/v1.3.0/translator.go#L863) ¶ added in v1.1.0
    
    
    type DefaultJSONTranslator struct{}

DefaultJSONTranslator converts an json.Feed struct into the generic Feed struct. 

This default implementation defines a set of mapping rules between json.Feed -> Feed for each of the fields in Feed. 

####  func (*DefaultJSONTranslator) [Translate](https://github.com/mmcdole/gofeed/blob/v1.3.0/translator.go#L867) ¶ added in v1.1.0
    
    
    func (t *DefaultJSONTranslator) Translate(feed interface{}) (*Feed, [error](/builtin#error))

Translate converts an JSON feed into the universal feed type. 

####  type [DefaultRSSTranslator](https://github.com/mmcdole/gofeed/blob/v1.3.0/translator.go#L30) ¶
    
    
    type DefaultRSSTranslator struct{}

DefaultRSSTranslator converts an rss.Feed struct into the generic Feed struct. 

This default implementation defines a set of mapping rules between rss.Feed -> Feed for each of the fields in Feed. 

####  func (*DefaultRSSTranslator) [Translate](https://github.com/mmcdole/gofeed/blob/v1.3.0/translator.go#L34) ¶
    
    
    func (t *DefaultRSSTranslator) Translate(feed interface{}) (*Feed, [error](/builtin#error))

Translate converts an RSS feed into the universal feed type. 

####  type [Enclosure](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L86) ¶
    
    
    type Enclosure struct {
    	URL    [string](/builtin#string) `json:"url,omitempty"`
    	Length [string](/builtin#string) `json:"length,omitempty"`
    	Type   [string](/builtin#string) `json:"type,omitempty"`
    }

Enclosure is a file associated with a given Item. 

####  type [Feed](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L15) ¶
    
    
    type Feed struct {
    	Title           [string](/builtin#string)                   `json:"title,omitempty"`
    	Description     [string](/builtin#string)                   `json:"description,omitempty"`
    	Link            [string](/builtin#string)                   `json:"link,omitempty"`
    	FeedLink        [string](/builtin#string)                   `json:"feedLink,omitempty"`
    	Links           [][string](/builtin#string)                 `json:"links,omitempty"`
    	Updated         [string](/builtin#string)                   `json:"updated,omitempty"`
    	UpdatedParsed   *[time](/time).[Time](/time#Time)               `json:"updatedParsed,omitempty"`
    	Published       [string](/builtin#string)                   `json:"published,omitempty"`
    	PublishedParsed *[time](/time).[Time](/time#Time)               `json:"publishedParsed,omitempty"`
    	Author          *Person                  `json:"author,omitempty"` // Deprecated: Use feed.Authors instead
    	Authors         []*Person                `json:"authors,omitempty"`
    	Language        [string](/builtin#string)                   `json:"language,omitempty"`
    	Image           *Image                   `json:"image,omitempty"`
    	Copyright       [string](/builtin#string)                   `json:"copyright,omitempty"`
    	Generator       [string](/builtin#string)                   `json:"generator,omitempty"`
    	Categories      [][string](/builtin#string)                 `json:"categories,omitempty"`
    	DublinCoreExt   *[ext](/github.com/mmcdole/gofeed@v1.3.0/extensions).[DublinCoreExtension](/github.com/mmcdole/gofeed@v1.3.0/extensions#DublinCoreExtension) `json:"dcExt,omitempty"`
    	ITunesExt       *[ext](/github.com/mmcdole/gofeed@v1.3.0/extensions).[ITunesFeedExtension](/github.com/mmcdole/gofeed@v1.3.0/extensions#ITunesFeedExtension) `json:"itunesExt,omitempty"`
    	Extensions      [ext](/github.com/mmcdole/gofeed@v1.3.0/extensions).[Extensions](/github.com/mmcdole/gofeed@v1.3.0/extensions#Extensions)           `json:"extensions,omitempty"`
    	Custom          map[[string](/builtin#string)][string](/builtin#string)        `json:"custom,omitempty"`
    	Items           []*Item                  `json:"items"`
    	FeedType        [string](/builtin#string)                   `json:"feedType"`
    	FeedVersion     [string](/builtin#string)                   `json:"feedVersion"`
    }

Feed is the universal Feed type that atom.Feed and rss.Feed gets translated to. It represents a web feed. Sorting with sort.Sort will order the Items by oldest to newest publish time. 

####  func (Feed) [Len](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L93) ¶
    
    
    func (f Feed) Len() [int](/builtin#int)

Len returns the length of Items. 

####  func (Feed) [Less](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L99) ¶
    
    
    func (f Feed) Less(i, k [int](/builtin#int)) [bool](/builtin#bool)

Less compares PublishedParsed of Items[i], Items[k] and returns true if Items[i] is less than Items[k]. 

####  func (Feed) [String](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L41) ¶
    
    
    func (f Feed) String() [string](/builtin#string)

####  func (Feed) [Swap](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L106) ¶
    
    
    func (f Feed) Swap(i, k [int](/builtin#int))

Swap swaps Items[i] and Items[k]. 

####  type [FeedType](https://github.com/mmcdole/gofeed/blob/v1.3.0/detector.go#L15) ¶
    
    
    type FeedType [int](/builtin#int)

FeedType represents one of the possible feed types that we can detect. 
    
    
    const (
    	// FeedTypeUnknown represents a feed that could not have its
    	// type determiend.
    	FeedTypeUnknown FeedType = [iota](/builtin#iota)
    	// FeedTypeAtom repesents an Atom feed
    	FeedTypeAtom
    	// FeedTypeRSS represents an RSS feed
    	FeedTypeRSS
    	// FeedTypeJSON represents a JSON feed
    	FeedTypeJSON
    )

####  func [DetectFeedType](https://github.com/mmcdole/gofeed/blob/v1.3.0/detector.go#L32) ¶
    
    
    func DetectFeedType(feed [io](/io).[Reader](/io#Reader)) FeedType

DetectFeedType attempts to determine the type of feed by looking for specific xml elements unique to the various feed types. 

Example ¶
    
    
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
    

Share Format Run

####  type [HTTPError](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L22) ¶
    
    
    type HTTPError struct {
    	StatusCode [int](/builtin#int)
    	Status     [string](/builtin#string)
    }

HTTPError represents an HTTP error returned by a server. 

####  func (HTTPError) [Error](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L27) ¶
    
    
    func (err HTTPError) Error() [string](/builtin#string)

####  type [Image](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L80) ¶
    
    
    type Image struct {
    	URL   [string](/builtin#string) `json:"url,omitempty"`
    	Title [string](/builtin#string) `json:"title,omitempty"`
    }

Image is an image that is the artwork for a given feed or item. 

####  type [Item](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L49) ¶
    
    
    type Item struct {
    	Title           [string](/builtin#string)                   `json:"title,omitempty"`
    	Description     [string](/builtin#string)                   `json:"description,omitempty"`
    	Content         [string](/builtin#string)                   `json:"content,omitempty"`
    	Link            [string](/builtin#string)                   `json:"link,omitempty"`
    	Links           [][string](/builtin#string)                 `json:"links,omitempty"`
    	Updated         [string](/builtin#string)                   `json:"updated,omitempty"`
    	UpdatedParsed   *[time](/time).[Time](/time#Time)               `json:"updatedParsed,omitempty"`
    	Published       [string](/builtin#string)                   `json:"published,omitempty"`
    	PublishedParsed *[time](/time).[Time](/time#Time)               `json:"publishedParsed,omitempty"`
    	Author          *Person                  `json:"author,omitempty"` // Deprecated: Use item.Authors instead
    	Authors         []*Person                `json:"authors,omitempty"`
    	GUID            [string](/builtin#string)                   `json:"guid,omitempty"`
    	Image           *Image                   `json:"image,omitempty"`
    	Categories      [][string](/builtin#string)                 `json:"categories,omitempty"`
    	Enclosures      []*Enclosure             `json:"enclosures,omitempty"`
    	DublinCoreExt   *[ext](/github.com/mmcdole/gofeed@v1.3.0/extensions).[DublinCoreExtension](/github.com/mmcdole/gofeed@v1.3.0/extensions#DublinCoreExtension) `json:"dcExt,omitempty"`
    	ITunesExt       *[ext](/github.com/mmcdole/gofeed@v1.3.0/extensions).[ITunesItemExtension](/github.com/mmcdole/gofeed@v1.3.0/extensions#ITunesItemExtension) `json:"itunesExt,omitempty"`
    	Extensions      [ext](/github.com/mmcdole/gofeed@v1.3.0/extensions).[Extensions](/github.com/mmcdole/gofeed@v1.3.0/extensions#Extensions)           `json:"extensions,omitempty"`
    	Custom          map[[string](/builtin#string)][string](/builtin#string)        `json:"custom,omitempty"`
    }

Item is the universal Item type that atom.Entry and rss.Item gets translated to. It represents a single entry in a given feed. 

####  type [Parser](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L34) ¶
    
    
    type Parser struct {
    	AtomTranslator Translator
    	RSSTranslator  Translator
    	JSONTranslator Translator
    	UserAgent      [string](/builtin#string)
    	AuthConfig     *Auth
    	Client         *[http](/net/http).[Client](/net/http#Client)
    	// contains filtered or unexported fields
    }

Parser is a universal feed parser that detects a given feed type, parsers it, and translates it to the universal feed type. 

####  func [NewParser](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L55) ¶
    
    
    func NewParser() *Parser

NewParser creates a universal feed parser. 

####  func (*Parser) [Parse](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L68) ¶
    
    
    func (f *Parser) Parse(feed [io](/io).[Reader](/io#Reader)) (*Feed, [error](/builtin#error))

Parse parses a RSS or Atom or JSON feed into the universal gofeed.Feed. It takes an io.Reader which should return the xml/json content. 

Example ¶
    
    
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
    

Share Format Run

####  func (*Parser) [ParseString](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L146) ¶
    
    
    func (f *Parser) ParseString(feed [string](/builtin#string)) (*Feed, [error](/builtin#error))

ParseString parses a feed XML string and into the universal feed type. 

Example ¶
    
    
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
    

Share Format Run

####  func (*Parser) [ParseURL](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L96) ¶
    
    
    func (f *Parser) ParseURL(feedURL [string](/builtin#string)) (feed *Feed, err [error](/builtin#error))

ParseURL fetches the contents of a given url and attempts to parse the response into the universal feed type. 

Example ¶
    
    
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
    

Share Format Run

####  func (*Parser) [ParseURLWithContext](https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L106) ¶
    
    
    func (f *Parser) ParseURLWithContext(feedURL [string](/builtin#string), ctx [context](/context).[Context](/context#Context)) (feed *Feed, err [error](/builtin#error))

ParseURLWithContext fetches contents of a given url and attempts to parse the response into the universal feed type. You can instantiate the Auth structure with your Username and Password to use the BasicAuth during the HTTP call. It will be automatically added to the header of the request Request could be canceled or timeout via given context 

####  type [Person](https://github.com/mmcdole/gofeed/blob/v1.3.0/feed.go#L73) ¶
    
    
    type Person struct {
    	Name  [string](/builtin#string) `json:"name,omitempty"`
    	Email [string](/builtin#string) `json:"email,omitempty"`
    }

Person is an individual specified in a feed (e.g. an author) 

####  type [Translator](https://github.com/mmcdole/gofeed/blob/v1.3.0/translator.go#L20) ¶
    
    
    type Translator interface {
    	Translate(feed interface{}) (*Feed, [error](/builtin#error))
    }

Translator converts a particular feed (atom.Feed or rss.Feed of json.Feed) into the generic Feed struct 
