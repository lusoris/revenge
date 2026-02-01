# OMDb API

> Source: https://www.omdbapi.com/
> Fetched: 2026-02-01T11:43:47.458643+00:00
> Content-Hash: 26f42876565a2122
> Type: html

---

OMDb API

- Usage
- Parameters
- Examples
- Change Log
- [API Key](apikey.aspx)

- [Become a Patron](https://www.patreon.com/join/omdb)
- [Donate](donate.htm)
- [Contact](/cdn-cgi/l/email-protection#abc9cdd9c2dfd1ebcdcacfc2c5ccd8c2ccc5cac785c8c4c6)

# OMDb API

The Open Movie Database

The OMDb API is a RESTful web service to obtain movie information, all content and images on the site are contributed and maintained by our users.
  
If you find this service useful, please consider making a [one-time donation](donate.htm) or [become a patron](https://www.patreon.com/join/omdb).

#### Poster API

The Poster API is only available to patrons.
  
Currently over 280,000 posters, updated daily with resolutions up to 2000x3000.

×

#### My Tesla Referral Link

[](https://www.tesla.com/referral/brian178422)

### Sponsors

[Emby](https://emby.media/), [Trakt](https://trakt.tv/), [FileBot](http://filebot.net/), [Reelgood](https://reelgood.com/), [Xirvik Servers](http://xirvik.com/), [Yidio](https://www.yidio.com/), [mi.tv](https://mi.tv/co/programacion), [Couchpop](https://couchpop.com/), [Edu Reviewer](https://edureviewer.com), [Flixboss](https://flixboss.com), [Scripts on Screen](https://scripts-onscreen.com/), [Topagency.webflow.io](https://topagency.webflow.io/), [Ramotion.com](https://www.ramotion.com/agency/web-app-development/), [Phone Trackers](https://celltrackingapps.com/), [Classics on DVD](https://dvdlady.com/), [Streaming App](https://streamin.app/), [How to make a FinTech app](https://www.purrweb.com/blog/how-to-create-a-fintech-app/), [Vid2 - Create movie lists with AI](https://vid2.com/), [Popflick](https://popflick.com/), [Custom Couch](https://whataroom.com/collections/custom-sofas), [TV.GURU](https://tv.guru/), [Deploy Vision - Nearshore Software Development](https://www.deployvision.com/), [Phictly](https://phictly.com/), [DoMyEssay](https://domyessay.com/do-my-homework "Need fast and reliable academic help? Let DoMyEssay do your homework and take the pressure off your shoulders."), [EssayPro](https://essaypro.com/do-my-homework "Do my homework with EssayPro"), [HomeworkGuy](https://homeworkguy.org/someone-to-take-my-online-class), [The Streaming Codex](https://github.com/LoSTxDragon/The-Streaming-Codex-Showcase "AI-powered movie recommendations that understand your taste"), [**WriteMyEssay**](https://www.writemyessay.com/ "WriteMyEssay offers a dedicated writing service for students, helping with essays and academic papers")

# Usage

Send all data requests to:

    http://www.omdbapi.com/?apikey=[yourkey]&

Poster API requests:

    http://img.omdbapi.com/?apikey=[yourkey]&

# Parameters

#### By ID or Title

Parameter | Required | Valid Options | Default Value | Description  
---|---|---|---|---  
i | Optional* |  | <empty> | A valid IMDb ID (e.g. tt1285016)  
t | Optional* |  | <empty> | Movie title to search for.  
type | No | movie, series, episode | <empty> | Type of result to return.  
y | No |  | <empty> | Year of release.  
plot | No | short, full | short | Return short or full plot.  
r | No | json, xml | json | The data type to return.  
callback | No |  | <empty> | JSONP callback name.  
v | No |  | 1 | API version (reserved for future use).  
*Please note while both "i" and "t" are optional at least one argument is required.

* * *

#### By Search

Parameter | Required | Valid options | Default Value | Description  
---|---|---|---|---  
s | Yes |  | <empty> | Movie title to search for.  
type | No | movie, series, episode | <empty> | Type of result to return.  
y | No |  | <empty> | Year of release.  
r | No | json, xml | json | The data type to return.  
page New! | No | 1-100 | 1 | Page number to return.  
callback | No |  | <empty> | JSONP callback name.  
v | No |  | 1 | API version (reserved for future use).  
  
# Examples

By Title

Title: Year: Plot: Short Full Response: JSON XML Search Reset

Request:

    [](javascript:;)

Response:

By ID

ID: Plot: Short Full Response: JSON XML Search Reset

Request:

    [](javascript:;)

Response:

# Change Log

- 04/08/19
  - Added support for eight digit IMDb IDs.
- 01/20/19
  - Supressed adult content from search results.
- 01/20/19
  - Added Swagger files ([YAML](http://www.omdbapi.com/swagger.yaml), [JSON](http://www.omdbapi.com/swagger.json)) to expose current API abilities and upcoming REST functions.
- 11/02/17
  - **FREE KEYS!** The "open" API is finally open again!
- 08/20/17
  - I created a [GitHub repository](https://github.com/omdbapi/OMDb-API/issues) for tracking bugs.
- 05/10/17
  - Due to some security concerns on how the keys were being distributed I updated the form to email them and also changed the algorithm used, which means your older keys not obtained through email will eventually stop working.
- 01/12/17
  - Removed single character restriction from title/search results.
- 06/11/16
  - "totalSeasons" count has been added to series results.
- 1/20/16
  - To accommodate search paging "totalResults" is now returned at the root level.
- 12/12/15
  - Search pagination added: [http://www.omdbapi.com/?s=Batman&**page=2**](http://www.omdbapi.com/?s=Batman&page=2)
- 11/16/15
  - Season+Episode now works with "i" parameter: [http://www.omdbapi.com/?**i=tt0944947** &Season=1](http://www.omdbapi.com/?i=tt0944947&Season=1)
  - Fixed the max pool size connection issues.
- 10/18/15
  - You can now return all episodes by using just the "Season" parameter: [http://www.omdbapi.com/?t=Game of Thrones&**Season=1**](http://www.omdbapi.com/?t=Game of Thrones&Season=1)
- 9/9/15
  - New server is up, response times should be < 500ms.
  - Setup a CDN/Caching service with [CloudFlare](http://www.cloudflare.com)
- 8/15/15
  - Created and Fixed a bad parsing error with JSON response. -Sorry about that!
  - HTTPS (with TLS) is now active: [https://www.omdbapi.com/](https://www.omdbapi.com)
- 5/10/15
  - Season+Episode search parameters added: [http://www.omdbapi.com/?t=Game of Thrones&**Season=1** &**Episode=1**](http://www.omdbapi.com/?t=Game of Thrones&Season=1&Episode=1)
