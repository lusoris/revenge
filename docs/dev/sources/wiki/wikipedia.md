# Wikipedia API

> Source: https://en.wikipedia.org/w/api.php?action=help&modules=query
> Fetched: 2026-01-30T23:56:37.865700+00:00
> Content-Hash: 3efb3fcf2b7bde5d
> Type: html

---

MediaWiki API help

This is an auto-generated MediaWiki API documentation page.

Documentation and examples:

https://www.mediawiki.org/wiki/Special:MyLanguage/API:Main_page

action=query

(

main

|

query

)

This module requires read rights.

Source:

MediaWiki

License:

GPL-2.0-or-later

Fetch data from and about MediaWiki.

All data modifications will first have to use query to acquire a token to prevent abuse from malicious sites.

https://www.mediawiki.org/wiki/Special:MyLanguage/API:Query

https://www.mediawiki.org/wiki/Special:MyLanguage/API:Meta

https://www.mediawiki.org/wiki/Special:MyLanguage/API:Properties

https://www.mediawiki.org/wiki/Special:MyLanguage/API:Lists

Specific parameters:

Other general parameters are available.

prop

Which properties to get for the queried pages.

categories

List all categories the pages belong to.

categoryinfo

Returns information about the given categories.

contributors

Get the list of registered contributors (including temporary users) and the count of anonymous contributors to a page.

coordinates

Returns coordinates of the given pages.

deletedrevisions

Get deleted revision information.

duplicatefiles

List all files that are duplicates of the given files based on hash values.

extlinks

Returns all external URLs (not interwikis) from the given pages.

extracts

Returns plain-text or limited HTML extracts of the given pages.

fileusage

Find all pages that use the given files.

flagged

Get information about the flagging status of the given pages.

globalusage

Returns global image usage for a certain image.

growthimagesuggestiondata

Fetch associated

image suggestion data

, if available

imageinfo

Returns file information and upload history.

images

Returns all files contained on the given pages.

info

Get basic page information.

isreviewed

Determine if a page is marked as reviewed.

iwlinks

Returns all interwiki links from the given pages.

langlinks

Returns all interlanguage links from the given pages.

langlinkscount

Get the number of other language versions.

links

Returns all links from the given pages.

linkshere

Find all pages that link to the given pages.

mmcontent

Get the description and targets of a spamlist

pageassessments

Return associated projects and assessments for the given pages.

pageimages

Returns information about images on the page, such as thumbnail and presence of photos.

pageprops

Get various page properties defined in the page content.

pageterms

Get the Wikidata terms (typically labels, descriptions and aliases) associated with a page via a sitelink.

pageviews

Shows per-page pageview data (the number of daily pageviews for each of the last

pvipdays

days).

redirects

Returns all redirects to the given pages.

revisions

Get revision information.

stashimageinfo

Returns file information for stashed files.

templates

Returns all pages transcluded on the given pages.

transcludedin

Find all pages that transclude the given pages.

transcodestatus

Get transcode status for a given file page.

videoinfo

Extends imageinfo to include video source (derivatives) information

wbentityusage

Returns all entity IDs used in the given pages.

cirrusbuilddoc

Internal.

Dump of a CirrusSearch article document from the database servers

cirruscompsuggestbuilddoc

Internal.

Dump of the document used by the completion suggester

cirrusdoc

Internal.

Dump of a CirrusSearch article document from the search servers

description

Internal.

Get a short description a.k.a. subtitle explaining what the target page is about.

mapdata

Internal.

Request all Kartographer map data for the given pages

Values (separate with

|

or

alternative

):

categories

,

categoryinfo

,

contributors

,

coordinates

,

deletedrevisions

,

duplicatefiles

,

extlinks

,

extracts

,

fileusage

,

flagged

,

globalusage

,

growthimagesuggestiondata

,

imageinfo

,

images

,

info

,

isreviewed

,

iwlinks

,

langlinks

,

langlinkscount

,

links

,

linkshere

,

mmcontent

,

pageassessments

,

pageimages

,

pageprops

,

pageterms

,

pageviews

,

redirects

,

revisions

,

stashimageinfo

,

templates

,

transcludedin

,

transcodestatus

,

videoinfo

,

wbentityusage

,

cirrusbuilddoc

,

cirruscompsuggestbuilddoc

,

cirrusdoc

,

description

,

mapdata

list

Which lists to get.

abusefilters

Show details of the edit filters.

abuselog

Show events that were caught by one of the edit filters.

allcategories

Enumerate all categories.

alldeletedrevisions

List all deleted revisions by a user or in a namespace.

allfileusages

List all file usages, including non-existing.

allimages

Enumerate all images sequentially.

alllinks

Enumerate all links that point to a given namespace.

allpages

Enumerate all pages sequentially in a given namespace.

allredirects

List all redirects to a namespace.

allrevisions

List all revisions.

alltransclusions

List all transclusions (pages embedded using {{x}}), including non-existing.

allusers

Enumerate all registered users.

automatictranslationdenselanguages

Fetch the list of sitelinks for the article that corresponds to a given Wikidata ID, ordered by article size.

backlinks

Find all pages that link to the given page.

betafeatures

List all BetaFeatures

blocks

List all blocked users and IP addresses.

categorymembers

List all pages in a given category.

centralnoticeactivecampaigns

Get a list of currently active campaigns with start and end dates and associated banners.

centralnoticelogs

Get a log of campaign configuration changes.

checkuserlog

Get entries from the CheckUser log.

codexicons

Get Codex icons

contenttranslation

Query Content Translation database for translations.

contenttranslationcorpora

Get the section-aligned parallel text for a given translation. See also

list=cxpublishedtranslations

. Dumps are provided in different formats for high volume access.

contenttranslationfavoritesuggestions

Get user's favorite suggestions for Content Translation.

cxpublishedtranslations

Fetch all published translations information.

cxtranslatorstats

Fetch the translation statistics for the given user.

embeddedin

Find all pages that embed (transclude) the given title.

exturlusage

Enumerate pages that contain a given URL.

filearchive

Enumerate all deleted files sequentially.

gadgetcategories

Returns a list of gadget sections.

gadgets

Returns a list of gadgets used on this wiki.

geosearch

Returns pages having coordinates that are located in a certain area.

globalallusers

Enumerate all global users.

globalblocks

List all globally blocked IP addresses.

globalgroups

Enumerate all global groups.

growthmentorlist

List all the mentors

growthmentormentee

Get all mentees assigned to a given mentor

growthstarredmentees

Get list of mentees starred by the currently logged in mentor

imageusage

Find all pages that use the given image title.

iwbacklinks

Find all pages that link to the given interwiki link.

langbacklinks

Find all pages that link to the given language link.

linterrors

Get a list of lint errors

logevents

Get events from logs.

mostviewed

Lists the most viewed pages (based on last day's pageview count).

mystashedfiles

Get a list of files in the current user's upload stash.

oldreviewedpages

Enumerates pages that have changes pending review.

pagecollectionsmetadata

Fetch page collection information for the given title.

pagepropnames

List all page property names in use on the wiki.

pageswithprop

List all pages using a given page property.

prefixsearch

Perform a prefix search for page titles.

projectpages

List all pages associated with one or more projects.

projects

List all the projects.

protectedtitles

List all titles protected from creation.

querypage

Get a list provided by a QueryPage-based special page.

random

Get a set of random pages.

recentchanges

Enumerate recent changes.

search

Perform a full text search.

tags

List change tags.

trackingcategories

Enumerate all existing tracking categories defined in

Special:TrackingCategories

. A tracking category exists if it contains pages or if its category page exists.

usercontribs

Get all edits by a user.

users

Get information about a list of users.

watchlist

Get recent changes to pages in the current user's watchlist.

watchlistraw

Get all pages on the current user's watchlist.

wblistentityusage

Returns all pages that use the given entity IDs.

wikisets

Enumerate all wiki sets.

checkuser

Deprecated.

This API has been disabled by the site administrators. Querying the API will return no data.

Check which IP addresses are used by a given username or which usernames are used by a given IP address.

deletedrevs

Deprecated.

List deleted revisions.

growthtasks

Internal.

Get task recommendations suitable for newcomers.

readinglistentries

Internal.

List the pages of a certain list.

Values (separate with

|

or

alternative

):

abusefilters

,

abuselog

,

allcategories

,

alldeletedrevisions

,

allfileusages

,

allimages

,

alllinks

,

allpages

,

allredirects

,

allrevisions

,

alltransclusions

,

allusers

,

automatictranslationdenselanguages

,

backlinks

,

betafeatures

,

blocks

,

categorymembers

,

centralnoticeactivecampaigns

,

centralnoticelogs

,

checkuserlog

,

codexicons

,

contenttranslation

,

contenttranslationcorpora

,

contenttranslationfavoritesuggestions

,

cxpublishedtranslations

,

cxtranslatorstats

,

embeddedin

,

exturlusage

,

filearchive

,

gadgetcategories

,

gadgets

,

geosearch

,

globalallusers

,

globalblocks

,

globalgroups

,

growthmentorlist

,

growthmentormentee

,

growthstarredmentees

,

imageusage

,

iwbacklinks

,

langbacklinks

,

linterrors

,

logevents

,

mostviewed

,

mystashedfiles

,

oldreviewedpages

,

pagecollectionsmetadata

,

pagepropnames

,

pageswithprop

,

prefixsearch

,

projectpages

,

projects

,

protectedtitles

,

querypage

,

random

,

recentchanges

,

search

,

tags

,

trackingcategories

,

usercontribs

,

users

,

watchlist

,

watchlistraw

,

wblistentityusage

,

wikisets

,

checkuser

,

deletedrevs

,

growthtasks

,

readinglistentries

Maximum number of values is 50 (500 for clients that are allowed higher limits).

meta

Which metadata to get.

allmessages

Return messages from this site.

authmanagerinfo

Retrieve information about the current authentication status.

babel

Get information about what languages the user knows

communityconfiguration

Read the community configuration

cxconfig

Get ContentTranslation local configuration settings.

featureusage

Get a summary of logged API feature usages for a user agent.

filerepoinfo

Return meta information about image repositories configured on the wiki.

globalpreferences

Retrieve global preferences for the current user.

globalrenamestatus

Show information about global renames that are in progress.

globaluserinfo

Show information about a global user.

growthmenteestatus

Query current user's mentee status; see documentation of action=growthsetmenteestatus for detailed information about individual statuses.

growthmentorstatus

Query current user's mentor status

languageinfo

Return information about available languages.

linterstats

Get number of lint errors in each category

notifications

Get notifications waiting for the current user.

ores

Return ORES configuration and model data for this wiki.

siteinfo

Return general information about the site.

siteviews

Shows sitewide pageview data (daily pageview totals for each of the last

pvisdays

days).

tokens

Gets tokens for data-modifying actions.

unreadnotificationpages

Get pages for which there are unread notifications for the current user.

userinfo

Get information about the current user.

wikibase

Get information about the Wikibase client and the associated Wikibase repository.

checkuserformattedblockinfo

Internal.

Return formatted block details for sitewide blocks affecting the current user.

cxdeletedtranslations

Internal.

Get the number of your published translations that were deleted.

growthnextsuggestedtasktype

Internal.

Get a suggested task type for a user to try next.

oath

Internal.

Check to see if two-factor authentication (OATH) is enabled for a user.

readinglists

Internal.

List or filter the user's reading lists and show metadata about them.

Values (separate with

|

or

alternative

):

allmessages

,

authmanagerinfo

,

babel

,

communityconfiguration

,

cxconfig

,

featureusage

,

filerepoinfo

,

globalpreferences

,

globalrenamestatus

,

globaluserinfo

,

growthmenteestatus

,

growthmentorstatus

,

languageinfo

,

linterstats

,

notifications

,

ores

,

siteinfo

,

siteviews

,

tokens

,

unreadnotificationpages

,

userinfo

,

wikibase

,

checkuserformattedblockinfo

,

cxdeletedtranslations

,

growthnextsuggestedtasktype

,

oath

,

readinglists

indexpageids

Include an additional pageids section listing all returned page IDs.

Type: boolean (

details

)

export

Export the current revisions of all given or generated pages.

Type: boolean (

details

)

exportnowrap

Return the export XML without wrapping it in an XML result (same format as

Special:Export

). Can only be used with query+export.

Type: boolean (

details

)

exportschema

Target the given version of the XML dump format when exporting. Can only be used with

query+export

.

One of the following values: 0.10, 0.11

Default: 0.11

iwurl

Whether to get the full URL if the title is an interwiki link.

Type: boolean (

details

)

continue

When more results are available, use this to continue. More detailed information on how to continue queries

can be found on mediawiki.org

.

rawcontinue

Return raw

query-continue

data for continuation.

Type: boolean (

details

)

titles

A list of titles to work on.

Separate values with

|

or

alternative

.

Maximum number of values is 50 (500 for clients that are allowed higher limits).

pageids

A list of page IDs to work on.

Type: list of integers

Separate values with

|

or

alternative

.

Maximum number of values is 50 (500 for clients that are allowed higher limits).

revids

A list of revision IDs to work on. Note that almost all query modules will convert revision IDs to the corresponding page ID and work on the latest revision instead. Only

prop=revisions

uses exact revisions for its response.

Type: list of integers

Separate values with

|

or

alternative

.

Maximum number of values is 50 (500 for clients that are allowed higher limits).

generator

Get the list of pages to work on by executing the specified query module.

Note:

Generator parameter names must be prefixed with a "g", see examples.

allcategories

Enumerate all categories.

alldeletedrevisions

List all deleted revisions by a user or in a namespace.

allfileusages

List all file usages, including non-existing.

allimages

Enumerate all images sequentially.

alllinks

Enumerate all links that point to a given namespace.

allpages

Enumerate all pages sequentially in a given namespace.

allredirects

List all redirects to a namespace.

allrevisions

List all revisions.

alltransclusions

List all transclusions (pages embedded using {{x}}), including non-existing.

backlinks

Find all pages that link to the given page.

categories

List all categories the pages belong to.

categorymembers

List all pages in a given category.

deletedrevisions

Get deleted revision information.

duplicatefiles

List all files that are duplicates of the given files based on hash values.

embeddedin

Find all pages that embed (transclude) the given title.

exturlusage

Enumerate pages that contain a given URL.

fileusage

Find all pages that use the given files.

geosearch

Returns pages having coordinates that are located in a certain area.

images

Returns all files contained on the given pages.

imageusage

Find all pages that use the given image title.

iwbacklinks

Find all pages that link to the given interwiki link.

langbacklinks

Find all pages that link to the given language link.

links

Returns all links from the given pages.

linkshere

Find all pages that link to the given pages.

mostviewed

Lists the most viewed pages (based on last day's pageview count).

oldreviewedpages

Enumerates pages that have changes pending review.

pageswithprop

List all pages using a given page property.

prefixsearch

Perform a prefix search for page titles.

projectpages

List all pages associated with one or more projects.

protectedtitles

List all titles protected from creation.

querypage

Get a list provided by a QueryPage-based special page.

random

Get a set of random pages.

recentchanges

Enumerate recent changes.

redirects

Returns all redirects to the given pages.

revisions

Get revision information.

search

Perform a full text search.

templates

Returns all pages transcluded on the given pages.

trackingcategories

Enumerate all existing tracking categories defined in

Special:TrackingCategories

. A tracking category exists if it contains pages or if its category page exists.

transcludedin

Find all pages that transclude the given pages.

watchlist

Get recent changes to pages in the current user's watchlist.

watchlistraw

Get all pages on the current user's watchlist.

wblistentityusage

Returns all pages that use the given entity IDs.

growthtasks

Internal.

Get task recommendations suitable for newcomers.

readinglistentries

Internal.

List the pages of a certain list.

One of the following values:

allcategories

,

alldeletedrevisions

,

allfileusages

,

allimages

,

alllinks

,

allpages

,

allredirects

,

allrevisions

,

alltransclusions

,

backlinks

,

categories

,

categorymembers

,

deletedrevisions

,

duplicatefiles

,

embeddedin

,

exturlusage

,

fileusage

,

geosearch

,

images

,

imageusage

,

iwbacklinks

,

langbacklinks

,

links

,

linkshere

,

mostviewed

,

oldreviewedpages

,

pageswithprop

,

prefixsearch

,

projectpages

,

protectedtitles

,

querypage

,

random

,

recentchanges

,

redirects

,

revisions

,

search

,

templates

,

trackingcategories

,

transcludedin

,

watchlist

,

watchlistraw

,

wblistentityusage

,

growthtasks

,

readinglistentries

redirects

Automatically resolve redirects in

query+titles

,

query+pageids

, and

query+revids

, and in pages returned by

query+generator

.

Type: boolean (

details

)

converttitles

Convert titles to other variants if necessary. Only works if the wiki's content language supports variant conversion. Languages that support variant conversion include ban, crh, en, gan, iu, ku, mni, sh, shi, sr, tg, tly, uz, wuu, zgh and zh.

Type: boolean (

details

)

Examples:

Fetch

site info

and

revisions

of

Main Page

.

api.php?action=query&prop=revisions&meta=siteinfo&titles=Main%20Page&rvprop=user|comment&continue=

[open in sandbox]

Fetch revisions of pages beginning with

API/

.

api.php?action=query&generator=allpages&gapprefix=API/&prop=revisions&continue=

[open in sandbox]

Retrieved from "

https://en.wikipedia.org/wiki/Special:ApiHelp

"