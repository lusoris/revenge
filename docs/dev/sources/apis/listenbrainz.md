# ListenBrainz API

> Auto-fetched from [https://listenbrainz.readthedocs.io/en/latest/users/api/index.html](https://listenbrainz.readthedocs.io/en/latest/users/api/index.html)
> Last Updated: 2026-01-28T21:45:41.345787+00:00

---

ListenBrainz
API Documentation
ListenBrainz API
Authentication
Get the User token
Add the User token to your requests
Reference
Core
Playlists
Recordings
Statistics
Popularity
Metadata
Social
Recommendations
Art
Settings
Miscellaneous
OpenAPI specification
Rate limiting
Usage Examples
Prerequisites
Examples
Submitting Listens
Getting Listen History
Lookup MBIDs
Love/hate feedback
Latest Import
JSON Documentation
Submission JSON
Fetching listen JSON
Payload JSON details
Client Metadata examples
BrainzPlayer on the ListenBrainz website playing a video from YouTube
BrainzPlayer on the ListenBrainz website playing a video from Spotify
Using Otter for Funkwhale on Android, and submitting with Simple Scrobbler
Rhythmbox player listening to Jamendo
Listening to a recording from Bandcamp and submitting with the browser extension WebScrobbler
Client Libraries
Haskell
Go
Rust
.NET
Python
Java
Swift
Last.FM Compatible API for ListenBrainz
AudioScrobbler API v1.2
Last.FM API
For development
For users
Data Dumps
Dump mirrors
File Descriptions
listenbrainz-public-dump.tar.zst
listenbrainz-listens-dump.tar.zst
listenbrainz-listens-dump-spark.tar.zst
Structure of the listens dump
Incremental dumps
ListenBrainz Data Update Intervals
ListenBrainz Data Update Intervals
Listens and Listen Counts
User Statistics
MBID Mapper & MusicBrainz Metadata Cache
ListenBrainz data infrastructure
Developer Documentation
Server development
Set up ListenBrainz Server development environment
Clone listenbrainz-server
Install docker
Register a MusicBrainz application
Update config.py
Initialize ListenBrainz containers
Initialize ListenBrainz databases
Run the magic script
Listenbrainz containers
Test your changes with unit tests
Lint your code
Using develop.sh
Spark development
Set up the webserver
Create listenbrainz_spark/config.py
Initialize ListenBrainz Spark containers
Bring containers up
Import data into the spark environment
Working with request_consumer
Test your changes with unit tests
Architecture
Services
Listen Flow
Frontend Rendering
Spark Architecture
Developing request_consumer
Start the webserver
Start the spark containers
Start the spark reader
MBID Mapping
Database tables
Fuzzy lookups
MBID Mapper
Scripts
ListenBrainz
./develop.sh manage
Dump Manager
./develop.sh manage dump
ListenBrainz Spark
./develop.sh manage spark
Troubleshooting
Docker Installations
Windows
Maintainer Documentation
Production Deployment
Cron
Building Docker Images
Production Images
Test Images
Using Github Actions
Using docker/push.sh script
Data Dumps
Check FTP Dumps age script
Logs
Manually triggering dumps
MBID Mapping
Containers
Data sources
Debugging lookups
Debugging Spotify Reader
RabbitMQ
Maintenance
Tolerance to connectivity issues
Maintenance mode
Data importance
Data persistence
Procedures
Implementation details
Updating Production Database Schema
Pull Requests Policy
ListenBrainz
ListenBrainz API
View page source
ListenBrainz API
¶
All endpoints have this root URL for our current production site.
API Root URL
:
https://api.listenbrainz.org
Note
All ListenBrainz services are only available on
HTTPS
!
Authentication
¶
ListenBrainz makes use of private API keys called user tokens to authenticate requests and ensure the proper
access controls on user data. A user token is a unique alphanumeric string linked to a user account. To retrieve
your user token, follow this guide.
Get the User token
¶
Every account has a User token associated with it, to get the token:
Sign up or Log in your an account using this
link
.
Navigate to
settings
page to find your user Token (See image below for reference).
Copy the User Token to your clipboard.
Note
You may also reset your user token by clicking the Reset token button on the
settings
page.
Add the User token to your requests
¶
The user token must be included in the request header for its usage.
To format the header correctly, you can use the following piece of code:
# The following token must be valid, but it doesn't have to be the token of the user you're
# trying to get the listen history of.
TOKEN
=
'YOUR_TOKEN_HERE'
AUTH_HEADER
=
{
"Authorization"
:
"Token
{0}
"
.
format
(
TOKEN
)
}
Then include the formatted header in the request to use it.
response
=
requests
.
get
(
...
# Your request url and params go here.
...
headers
=
AUTH_HEADER
,
)
Note
A complete usage example for a request employing Authorization headers to make authenticated requests to ListenBrainz
can be found on the
API Usage
page.
Reference
¶
Core
Playlists
Recordings
Statistics
Popularity
Metadata
Social
Recommendations
Art
Settings
Miscellaneous
OpenAPI specification
¶
Contributor
rain0r
went through the trouble of making
an OpenAPI 3 specification for the ListenBrainz API. Many thanks! Check it out here:
https://github.com/rain0r/listenbrainz-openapi
Rate limiting
¶
The ListenBrainz API is rate limited via the use of rate limiting headers that
are sent as part of the HTTP response headers. Each call will include the
following headers:
X-RateLimit-Limit
: Number of requests allowed in given time window
X-RateLimit-Remaining
: Number of requests remaining in current time
window
X-RateLimit-Reset-In
: Number of seconds when current time window expires
(
recommended
: this header is resilient against clients with incorrect
clocks)
X-RateLimit-Reset
: UNIX epoch number of seconds (without timezone) when
current time window expires
[
1
]
Rate limiting is automatic and the client must use these headers to determine
the rate to make API calls. If the client exceeds the number of requests
allowed, the server will respond with error code
429:
Too
Many
Requests
.
Requests that provide the
Authorization
header with a valid user token may
receive higher rate limits than those without valid user tokens.
[
1
]
Provided for compatibility with other APIs, but we still recommend using
X-RateLimit-Reset-In
wherever possible
Previous
Next
© Copyright 2017-2026, MetaBrainz Foundation.
Built with
Sphinx
using a
theme
provided by
Read the Docs
.