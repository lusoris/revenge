# Wikipedia API

> Source: https://en.wikipedia.org/w/api.php?action=help&modules=query
> Fetched: 2026-02-01T11:51:15.944548+00:00
> Content-Hash: ce7bd2c677113ec4
> Type: html

---

# MediaWiki API help

This is an auto-generated MediaWiki API documentation page. 

Documentation and examples: <https://www.mediawiki.org/wiki/Special:MyLanguage/API:Main_page>

## action=query

([main](/wiki/Special:ApiHelp/main) | **query**)

  * This module requires read rights.
  * Source: MediaWiki
  * License: [GPL-2.0-or-later](/wiki/Special:Version/License/MediaWiki "Special:Version/License/MediaWiki")



Fetch data from and about MediaWiki. 

All data modifications will first have to use query to acquire a token to prevent abuse from malicious sites. 

  * <https://www.mediawiki.org/wiki/Special:MyLanguage/API:Query>
  * <https://www.mediawiki.org/wiki/Special:MyLanguage/API:Meta>
  * <https://www.mediawiki.org/wiki/Special:MyLanguage/API:Properties>
  * <https://www.mediawiki.org/wiki/Special:MyLanguage/API:Lists>



Specific parameters:

Other general parameters are available.

prop
    

Which properties to get for the queried pages. 

[categories](/w/api.php?action=help&modules=query%2Bcategories)
    List all categories the pages belong to.
[categoryinfo](/w/api.php?action=help&modules=query%2Bcategoryinfo)
    Returns information about the given categories.
[contributors](/w/api.php?action=help&modules=query%2Bcontributors)
    Get the list of registered contributors (including temporary users) and the count of anonymous contributors to a page.
[coordinates](/w/api.php?action=help&modules=query%2Bcoordinates)
    Returns coordinates of the given pages.
[deletedrevisions](/w/api.php?action=help&modules=query%2Bdeletedrevisions)
    Get deleted revision information.
[duplicatefiles](/w/api.php?action=help&modules=query%2Bduplicatefiles)
    List all files that are duplicates of the given files based on hash values.
[extlinks](/w/api.php?action=help&modules=query%2Bextlinks)
    Returns all external URLs (not interwikis) from the given pages.
[extracts](/w/api.php?action=help&modules=query%2Bextracts)
    Returns plain-text or limited HTML extracts of the given pages.
[fileusage](/w/api.php?action=help&modules=query%2Bfileusage)
    Find all pages that use the given files.
[flagged](/w/api.php?action=help&modules=query%2Bflagged)
    Get information about the flagging status of the given pages.
[globalusage](/w/api.php?action=help&modules=query%2Bglobalusage)
    Returns global image usage for a certain image.
[growthimagesuggestiondata](/w/api.php?action=help&modules=query%2Bgrowthimagesuggestiondata)
    Fetch associated [image suggestion data](https://www.mediawiki.org/wiki/wikitech:Add_Image), if available
[imageinfo](/w/api.php?action=help&modules=query%2Bimageinfo)
    Returns file information and upload history.
[images](/w/api.php?action=help&modules=query%2Bimages)
    Returns all files contained on the given pages.
[info](/w/api.php?action=help&modules=query%2Binfo)
    Get basic page information.
[isreviewed](/w/api.php?action=help&modules=query%2Bisreviewed)
    Determine if a page is marked as reviewed.
[iwlinks](/w/api.php?action=help&modules=query%2Biwlinks)
    Returns all interwiki links from the given pages.
[langlinks](/w/api.php?action=help&modules=query%2Blanglinks)
    Returns all interlanguage links from the given pages.
[langlinkscount](/w/api.php?action=help&modules=query%2Blanglinkscount)
    Get the number of other language versions.
[links](/w/api.php?action=help&modules=query%2Blinks)
    Returns all links from the given pages.
[linkshere](/w/api.php?action=help&modules=query%2Blinkshere)
    Find all pages that link to the given pages.
[mmcontent](/w/api.php?action=help&modules=query%2Bmmcontent)
    Get the description and targets of a spamlist
[pageassessments](/w/api.php?action=help&modules=query%2Bpageassessments)
    Return associated projects and assessments for the given pages.
[pageimages](/w/api.php?action=help&modules=query%2Bpageimages)
    Returns information about images on the page, such as thumbnail and presence of photos.
[pageprops](/w/api.php?action=help&modules=query%2Bpageprops)
    Get various page properties defined in the page content.
[pageterms](/w/api.php?action=help&modules=query%2Bpageterms)
    Get the Wikidata terms (typically labels, descriptions and aliases) associated with a page via a sitelink.
[pageviews](/w/api.php?action=help&modules=query%2Bpageviews)
    Shows per-page pageview data (the number of daily pageviews for each of the last pvipdays days).
[redirects](/w/api.php?action=help&modules=query%2Bredirects)
    Returns all redirects to the given pages.
[revisions](/w/api.php?action=help&modules=query%2Brevisions)
    Get revision information.
[stashimageinfo](/w/api.php?action=help&modules=query%2Bstashimageinfo)
    Returns file information for stashed files.
[templates](/w/api.php?action=help&modules=query%2Btemplates)
    Returns all pages transcluded on the given pages.
[transcludedin](/w/api.php?action=help&modules=query%2Btranscludedin)
    Find all pages that transclude the given pages.
[transcodestatus](/w/api.php?action=help&modules=query%2Btranscodestatus)
    Get transcode status for a given file page.
[videoinfo](/w/api.php?action=help&modules=query%2Bvideoinfo)
    Extends imageinfo to include video source (derivatives) information
[wbentityusage](/w/api.php?action=help&modules=query%2Bwbentityusage)
    Returns all entity IDs used in the given pages.
[cirrusbuilddoc](/w/api.php?action=help&modules=query%2Bcirrusbuilddoc)
    Internal. Dump of a CirrusSearch article document from the database servers
[cirruscompsuggestbuilddoc](/w/api.php?action=help&modules=query%2Bcirruscompsuggestbuilddoc)
    Internal. Dump of the document used by the completion suggester
[cirrusdoc](/w/api.php?action=help&modules=query%2Bcirrusdoc)
    Internal. Dump of a CirrusSearch article document from the search servers
[description](/w/api.php?action=help&modules=query%2Bdescription)
    Internal. Get a short description a.k.a. subtitle explaining what the target page is about.
[mapdata](/w/api.php?action=help&modules=query%2Bmapdata)
    Internal. Request all Kartographer map data for the given pages
    Values (separate with `|` or [alternative](/w/api.php?action=help&modules=main#main/datatypes)): [categories](/w/api.php?action=help&modules=query%2Bcategories), [categoryinfo](/w/api.php?action=help&modules=query%2Bcategoryinfo), [contributors](/w/api.php?action=help&modules=query%2Bcontributors), [coordinates](/w/api.php?action=help&modules=query%2Bcoordinates), [deletedrevisions](/w/api.php?action=help&modules=query%2Bdeletedrevisions), [duplicatefiles](/w/api.php?action=help&modules=query%2Bduplicatefiles), [extlinks](/w/api.php?action=help&modules=query%2Bextlinks), [extracts](/w/api.php?action=help&modules=query%2Bextracts), [fileusage](/w/api.php?action=help&modules=query%2Bfileusage), [flagged](/w/api.php?action=help&modules=query%2Bflagged), [globalusage](/w/api.php?action=help&modules=query%2Bglobalusage), [growthimagesuggestiondata](/w/api.php?action=help&modules=query%2Bgrowthimagesuggestiondata), [imageinfo](/w/api.php?action=help&modules=query%2Bimageinfo), [images](/w/api.php?action=help&modules=query%2Bimages), [info](/w/api.php?action=help&modules=query%2Binfo), [isreviewed](/w/api.php?action=help&modules=query%2Bisreviewed), [iwlinks](/w/api.php?action=help&modules=query%2Biwlinks), [langlinks](/w/api.php?action=help&modules=query%2Blanglinks), [langlinkscount](/w/api.php?action=help&modules=query%2Blanglinkscount), [links](/w/api.php?action=help&modules=query%2Blinks), [linkshere](/w/api.php?action=help&modules=query%2Blinkshere), [mmcontent](/w/api.php?action=help&modules=query%2Bmmcontent), [pageassessments](/w/api.php?action=help&modules=query%2Bpageassessments), [pageimages](/w/api.php?action=help&modules=query%2Bpageimages), [pageprops](/w/api.php?action=help&modules=query%2Bpageprops), [pageterms](/w/api.php?action=help&modules=query%2Bpageterms), [pageviews](/w/api.php?action=help&modules=query%2Bpageviews), [redirects](/w/api.php?action=help&modules=query%2Bredirects), [revisions](/w/api.php?action=help&modules=query%2Brevisions), [stashimageinfo](/w/api.php?action=help&modules=query%2Bstashimageinfo), [templates](/w/api.php?action=help&modules=query%2Btemplates), [transcludedin](/w/api.php?action=help&modules=query%2Btranscludedin), [transcodestatus](/w/api.php?action=help&modules=query%2Btranscodestatus), [videoinfo](/w/api.php?action=help&modules=query%2Bvideoinfo), [wbentityusage](/w/api.php?action=help&modules=query%2Bwbentityusage), [cirrusbuilddoc](/w/api.php?action=help&modules=query%2Bcirrusbuilddoc), [cirruscompsuggestbuilddoc](/w/api.php?action=help&modules=query%2Bcirruscompsuggestbuilddoc), [cirrusdoc](/w/api.php?action=help&modules=query%2Bcirrusdoc), [description](/w/api.php?action=help&modules=query%2Bdescription), [mapdata](/w/api.php?action=help&modules=query%2Bmapdata)
list
    

Which lists to get. 

[abusefilters](/w/api.php?action=help&modules=query%2Babusefilters)
    Show details of the edit filters.
[abuselog](/w/api.php?action=help&modules=query%2Babuselog)
    Show events that were caught by one of the edit filters.
[allcategories](/w/api.php?action=help&modules=query%2Ballcategories)
    Enumerate all categories.
[alldeletedrevisions](/w/api.php?action=help&modules=query%2Balldeletedrevisions)
    List all deleted revisions by a user or in a namespace.
[allfileusages](/w/api.php?action=help&modules=query%2Ballfileusages)
    List all file usages, including non-existing.
[allimages](/w/api.php?action=help&modules=query%2Ballimages)
    Enumerate all images sequentially.
[alllinks](/w/api.php?action=help&modules=query%2Balllinks)
    Enumerate all links that point to a given namespace.
[allpages](/w/api.php?action=help&modules=query%2Ballpages)
    Enumerate all pages sequentially in a given namespace.
[allredirects](/w/api.php?action=help&modules=query%2Ballredirects)
    List all redirects to a namespace.
[allrevisions](/w/api.php?action=help&modules=query%2Ballrevisions)
    List all revisions.
[alltransclusions](/w/api.php?action=help&modules=query%2Balltransclusions)
    List all transclusions (pages embedded using {{x}}), including non-existing.
[allusers](/w/api.php?action=help&modules=query%2Ballusers)
    Enumerate all registered users.
[automatictranslationdenselanguages](/w/api.php?action=help&modules=query%2Bautomatictranslationdenselanguages)
    Fetch the list of sitelinks for the article that corresponds to a given Wikidata ID, ordered by article size.
[backlinks](/w/api.php?action=help&modules=query%2Bbacklinks)
    Find all pages that link to the given page.
[betafeatures](/w/api.php?action=help&modules=query%2Bbetafeatures)
    List all BetaFeatures
[blocks](/w/api.php?action=help&modules=query%2Bblocks)
    List all blocked users and IP addresses.
[categorymembers](/w/api.php?action=help&modules=query%2Bcategorymembers)
    List all pages in a given category.
[centralnoticeactivecampaigns](/w/api.php?action=help&modules=query%2Bcentralnoticeactivecampaigns)
    Get a list of currently active campaigns with start and end dates and associated banners.
[centralnoticelogs](/w/api.php?action=help&modules=query%2Bcentralnoticelogs)
    Get a log of campaign configuration changes.
[checkuserlog](/w/api.php?action=help&modules=query%2Bcheckuserlog)
    Get entries from the CheckUser log.
[codexicons](/w/api.php?action=help&modules=query%2Bcodexicons)
    Get Codex icons
[contenttranslation](/w/api.php?action=help&modules=query%2Bcontenttranslation)
    Query Content Translation database for translations.
[contenttranslationcorpora](/w/api.php?action=help&modules=query%2Bcontenttranslationcorpora)
    Get the section-aligned parallel text for a given translation. See also `list=cxpublishedtranslations`. Dumps are provided in different formats for high volume access.
[contenttranslationfavoritesuggestions](/w/api.php?action=help&modules=query%2Bcontenttranslationfavoritesuggestions)
    Get user's favorite suggestions for Content Translation.
[cxpublishedtranslations](/w/api.php?action=help&modules=query%2Bcxpublishedtranslations)
    Fetch all published translations information.
[cxtranslatorstats](/w/api.php?action=help&modules=query%2Bcxtranslatorstats)
    Fetch the translation statistics for the given user.
[embeddedin](/w/api.php?action=help&modules=query%2Bembeddedin)
    Find all pages that embed (transclude) the given title.
[exturlusage](/w/api.php?action=help&modules=query%2Bexturlusage)
    Enumerate pages that contain a given URL.
[filearchive](/w/api.php?action=help&modules=query%2Bfilearchive)
    Enumerate all deleted files sequentially.
[gadgetcategories](/w/api.php?action=help&modules=query%2Bgadgetcategories)
    Returns a list of gadget sections.
[gadgets](/w/api.php?action=help&modules=query%2Bgadgets)
    Returns a list of gadgets used on this wiki.
[geosearch](/w/api.php?action=help&modules=query%2Bgeosearch)
    Returns pages having coordinates that are located in a certain area.
[globalallusers](/w/api.php?action=help&modules=query%2Bglobalallusers)
    Enumerate all global users.
[globalblocks](/w/api.php?action=help&modules=query%2Bglobalblocks)
    List all globally blocked IP addresses.
[globalgroups](/w/api.php?action=help&modules=query%2Bglobalgroups)
    Enumerate all global groups.
[growthmentorlist](/w/api.php?action=help&modules=query%2Bgrowthmentorlist)
    List all the mentors
[growthmentormentee](/w/api.php?action=help&modules=query%2Bgrowthmentormentee)
    Get all mentees assigned to a given mentor
[growthstarredmentees](/w/api.php?action=help&modules=query%2Bgrowthstarredmentees)
    Get list of mentees starred by the currently logged in mentor
[imageusage](/w/api.php?action=help&modules=query%2Bimageusage)
    Find all pages that use the given image title.
[iwbacklinks](/w/api.php?action=help&modules=query%2Biwbacklinks)
    Find all pages that link to the given interwiki link.
[langbacklinks](/w/api.php?action=help&modules=query%2Blangbacklinks)
    Find all pages that link to the given language link.
[linterrors](/w/api.php?action=help&modules=query%2Blinterrors)
    Get a list of lint errors
[logevents](/w/api.php?action=help&modules=query%2Blogevents)
    Get events from logs.
[mostviewed](/w/api.php?action=help&modules=query%2Bmostviewed)
    Lists the most viewed pages (based on last day's pageview count).
[mystashedfiles](/w/api.php?action=help&modules=query%2Bmystashedfiles)
    Get a list of files in the current user's upload stash.
[oldreviewedpages](/w/api.php?action=help&modules=query%2Boldreviewedpages)
    Enumerates pages that have changes pending review.
[pagecollectionsmetadata](/w/api.php?action=help&modules=query%2Bpagecollectionsmetadata)
    Fetch page collection information for the given title.
[pagepropnames](/w/api.php?action=help&modules=query%2Bpagepropnames)
    List all page property names in use on the wiki.
[pageswithprop](/w/api.php?action=help&modules=query%2Bpageswithprop)
    List all pages using a given page property.
[prefixsearch](/w/api.php?action=help&modules=query%2Bprefixsearch)
    Perform a prefix search for page titles.
[projectpages](/w/api.php?action=help&modules=query%2Bprojectpages)
    List all pages associated with one or more projects.
[projects](/w/api.php?action=help&modules=query%2Bprojects)
    List all the projects.
[protectedtitles](/w/api.php?action=help&modules=query%2Bprotectedtitles)
    List all titles protected from creation.
[querypage](/w/api.php?action=help&modules=query%2Bquerypage)
    Get a list provided by a QueryPage-based special page.
[random](/w/api.php?action=help&modules=query%2Brandom)
    Get a set of random pages.
[recentchanges](/w/api.php?action=help&modules=query%2Brecentchanges)
    Enumerate recent changes.
[search](/w/api.php?action=help&modules=query%2Bsearch)
    Perform a full text search.
[tags](/w/api.php?action=help&modules=query%2Btags)
    List change tags.
[trackingcategories](/w/api.php?action=help&modules=query%2Btrackingcategories)
    Enumerate all existing tracking categories defined in [Special:TrackingCategories](/wiki/Special:TrackingCategories "Special:TrackingCategories"). A tracking category exists if it contains pages or if its category page exists.
[usercontribs](/w/api.php?action=help&modules=query%2Busercontribs)
    Get all edits by a user.
[users](/w/api.php?action=help&modules=query%2Busers)
    Get information about a list of users.
[watchlist](/w/api.php?action=help&modules=query%2Bwatchlist)
    Get recent changes to pages in the current user's watchlist.
[watchlistraw](/w/api.php?action=help&modules=query%2Bwatchlistraw)
    Get all pages on the current user's watchlist.
[wblistentityusage](/w/api.php?action=help&modules=query%2Bwblistentityusage)
    Returns all pages that use the given entity IDs.
[wikisets](/w/api.php?action=help&modules=query%2Bwikisets)
    Enumerate all wiki sets.
[checkuser](/w/api.php?action=help&modules=query%2Bcheckuser)
    Deprecated. **This API has been disabled by the site administrators. Querying the API will return no data.** Check which IP addresses are used by a given username or which usernames are used by a given IP address.
[deletedrevs](/w/api.php?action=help&modules=query%2Bdeletedrevs)
    Deprecated. List deleted revisions.
[growthtasks](/w/api.php?action=help&modules=query%2Bgrowthtasks)
    Internal. Get task recommendations suitable for newcomers.
[readinglistentries](/w/api.php?action=help&modules=query%2Breadinglistentries)
    Internal. List the pages of a certain list.
    Values (separate with `|` or [alternative](/w/api.php?action=help&modules=main#main/datatypes)): [abusefilters](/w/api.php?action=help&modules=query%2Babusefilters), [abuselog](/w/api.php?action=help&modules=query%2Babuselog), [allcategories](/w/api.php?action=help&modules=query%2Ballcategories), [alldeletedrevisions](/w/api.php?action=help&modules=query%2Balldeletedrevisions), [allfileusages](/w/api.php?action=help&modules=query%2Ballfileusages), [allimages](/w/api.php?action=help&modules=query%2Ballimages), [alllinks](/w/api.php?action=help&modules=query%2Balllinks), [allpages](/w/api.php?action=help&modules=query%2Ballpages), [allredirects](/w/api.php?action=help&modules=query%2Ballredirects), [allrevisions](/w/api.php?action=help&modules=query%2Ballrevisions), [alltransclusions](/w/api.php?action=help&modules=query%2Balltransclusions), [allusers](/w/api.php?action=help&modules=query%2Ballusers), [automatictranslationdenselanguages](/w/api.php?action=help&modules=query%2Bautomatictranslationdenselanguages), [backlinks](/w/api.php?action=help&modules=query%2Bbacklinks), [betafeatures](/w/api.php?action=help&modules=query%2Bbetafeatures), [blocks](/w/api.php?action=help&modules=query%2Bblocks), [categorymembers](/w/api.php?action=help&modules=query%2Bcategorymembers), [centralnoticeactivecampaigns](/w/api.php?action=help&modules=query%2Bcentralnoticeactivecampaigns), [centralnoticelogs](/w/api.php?action=help&modules=query%2Bcentralnoticelogs), [checkuserlog](/w/api.php?action=help&modules=query%2Bcheckuserlog), [codexicons](/w/api.php?action=help&modules=query%2Bcodexicons), [contenttranslation](/w/api.php?action=help&modules=query%2Bcontenttranslation), [contenttranslationcorpora](/w/api.php?action=help&modules=query%2Bcontenttranslationcorpora), [contenttranslationfavoritesuggestions](/w/api.php?action=help&modules=query%2Bcontenttranslationfavoritesuggestions), [cxpublishedtranslations](/w/api.php?action=help&modules=query%2Bcxpublishedtranslations), [cxtranslatorstats](/w/api.php?action=help&modules=query%2Bcxtranslatorstats), [embeddedin](/w/api.php?action=help&modules=query%2Bembeddedin), [exturlusage](/w/api.php?action=help&modules=query%2Bexturlusage), [filearchive](/w/api.php?action=help&modules=query%2Bfilearchive), [gadgetcategories](/w/api.php?action=help&modules=query%2Bgadgetcategories), [gadgets](/w/api.php?action=help&modules=query%2Bgadgets), [geosearch](/w/api.php?action=help&modules=query%2Bgeosearch), [globalallusers](/w/api.php?action=help&modules=query%2Bglobalallusers), [globalblocks](/w/api.php?action=help&modules=query%2Bglobalblocks), [globalgroups](/w/api.php?action=help&modules=query%2Bglobalgroups), [growthmentorlist](/w/api.php?action=help&modules=query%2Bgrowthmentorlist), [growthmentormentee](/w/api.php?action=help&modules=query%2Bgrowthmentormentee), [growthstarredmentees](/w/api.php?action=help&modules=query%2Bgrowthstarredmentees), [imageusage](/w/api.php?action=help&modules=query%2Bimageusage), [iwbacklinks](/w/api.php?action=help&modules=query%2Biwbacklinks), [langbacklinks](/w/api.php?action=help&modules=query%2Blangbacklinks), [linterrors](/w/api.php?action=help&modules=query%2Blinterrors), [logevents](/w/api.php?action=help&modules=query%2Blogevents), [mostviewed](/w/api.php?action=help&modules=query%2Bmostviewed), [mystashedfiles](/w/api.php?action=help&modules=query%2Bmystashedfiles), [oldreviewedpages](/w/api.php?action=help&modules=query%2Boldreviewedpages), [pagecollectionsmetadata](/w/api.php?action=help&modules=query%2Bpagecollectionsmetadata), [pagepropnames](/w/api.php?action=help&modules=query%2Bpagepropnames), [pageswithprop](/w/api.php?action=help&modules=query%2Bpageswithprop), [prefixsearch](/w/api.php?action=help&modules=query%2Bprefixsearch), [projectpages](/w/api.php?action=help&modules=query%2Bprojectpages), [projects](/w/api.php?action=help&modules=query%2Bprojects), [protectedtitles](/w/api.php?action=help&modules=query%2Bprotectedtitles), [querypage](/w/api.php?action=help&modules=query%2Bquerypage), [random](/w/api.php?action=help&modules=query%2Brandom), [recentchanges](/w/api.php?action=help&modules=query%2Brecentchanges), [search](/w/api.php?action=help&modules=query%2Bsearch), [tags](/w/api.php?action=help&modules=query%2Btags), [trackingcategories](/w/api.php?action=help&modules=query%2Btrackingcategories), [usercontribs](/w/api.php?action=help&modules=query%2Busercontribs), [users](/w/api.php?action=help&modules=query%2Busers), [watchlist](/w/api.php?action=help&modules=query%2Bwatchlist), [watchlistraw](/w/api.php?action=help&modules=query%2Bwatchlistraw), [wblistentityusage](/w/api.php?action=help&modules=query%2Bwblistentityusage), [wikisets](/w/api.php?action=help&modules=query%2Bwikisets), [checkuser](/w/api.php?action=help&modules=query%2Bcheckuser), [deletedrevs](/w/api.php?action=help&modules=query%2Bdeletedrevs), [growthtasks](/w/api.php?action=help&modules=query%2Bgrowthtasks), [readinglistentries](/w/api.php?action=help&modules=query%2Breadinglistentries)
    Maximum number of values is 50 (500 for clients that are allowed higher limits).
meta
    

Which metadata to get. 

[allmessages](/w/api.php?action=help&modules=query%2Ballmessages)
    Return messages from this site.
[authmanagerinfo](/w/api.php?action=help&modules=query%2Bauthmanagerinfo)
    Retrieve information about the current authentication status.
[babel](/w/api.php?action=help&modules=query%2Bbabel)
    Get information about what languages the user knows
[communityconfiguration](/w/api.php?action=help&modules=query%2Bcommunityconfiguration)
    Read the community configuration
[cxconfig](/w/api.php?action=help&modules=query%2Bcxconfig)
    Get ContentTranslation local configuration settings.
[featureusage](/w/api.php?action=help&modules=query%2Bfeatureusage)
    Get a summary of logged API feature usages for a user agent.
[filerepoinfo](/w/api.php?action=help&modules=query%2Bfilerepoinfo)
    Return meta information about image repositories configured on the wiki.
[globalpreferences](/w/api.php?action=help&modules=query%2Bglobalpreferences)
    Retrieve global preferences for the current user.
[globalrenamestatus](/w/api.php?action=help&modules=query%2Bglobalrenamestatus)
    Show information about global renames that are in progress.
[globaluserinfo](/w/api.php?action=help&modules=query%2Bglobaluserinfo)
    Show information about a global user.
[growthmenteestatus](/w/api.php?action=help&modules=query%2Bgrowthmenteestatus)
    Query current user's mentee status; see documentation of action=growthsetmenteestatus for detailed information about individual statuses.
[growthmentorstatus](/w/api.php?action=help&modules=query%2Bgrowthmentorstatus)
    Query current user's mentor status
[languageinfo](/w/api.php?action=help&modules=query%2Blanguageinfo)
    Return information about available languages.
[linterstats](/w/api.php?action=help&modules=query%2Blinterstats)
    Get number of lint errors in each category
[notifications](/w/api.php?action=help&modules=query%2Bnotifications)
    Get notifications waiting for the current user.
[ores](/w/api.php?action=help&modules=query%2Bores)
    Return ORES configuration and model data for this wiki.
[siteinfo](/w/api.php?action=help&modules=query%2Bsiteinfo)
    Return general information about the site.
[siteviews](/w/api.php?action=help&modules=query%2Bsiteviews)
    Shows sitewide pageview data (daily pageview totals for each of the last pvisdays days).
[tokens](/w/api.php?action=help&modules=query%2Btokens)
    Gets tokens for data-modifying actions.
[unreadnotificationpages](/w/api.php?action=help&modules=query%2Bunreadnotificationpages)
    Get pages for which there are unread notifications for the current user.
[userinfo](/w/api.php?action=help&modules=query%2Buserinfo)
    Get information about the current user.
[wikibase](/w/api.php?action=help&modules=query%2Bwikibase)
    Get information about the Wikibase client and the associated Wikibase repository.
[checkuserformattedblockinfo](/w/api.php?action=help&modules=query%2Bcheckuserformattedblockinfo)
    Internal. Return formatted block details for sitewide blocks affecting the current user.
[cxdeletedtranslations](/w/api.php?action=help&modules=query%2Bcxdeletedtranslations)
    Internal. Get the number of your published translations that were deleted.
[growthnextsuggestedtasktype](/w/api.php?action=help&modules=query%2Bgrowthnextsuggestedtasktype)
    Internal. Get a suggested task type for a user to try next.
[oath](/w/api.php?action=help&modules=query%2Boath)
    Internal. Check to see if two-factor authentication (OATH) is enabled for a user.
[readinglists](/w/api.php?action=help&modules=query%2Breadinglists)
    Internal. List or filter the user's reading lists and show metadata about them.
    Values (separate with `|` or [alternative](/w/api.php?action=help&modules=main#main/datatypes)): [allmessages](/w/api.php?action=help&modules=query%2Ballmessages), [authmanagerinfo](/w/api.php?action=help&modules=query%2Bauthmanagerinfo), [babel](/w/api.php?action=help&modules=query%2Bbabel), [communityconfiguration](/w/api.php?action=help&modules=query%2Bcommunityconfiguration), [cxconfig](/w/api.php?action=help&modules=query%2Bcxconfig), [featureusage](/w/api.php?action=help&modules=query%2Bfeatureusage), [filerepoinfo](/w/api.php?action=help&modules=query%2Bfilerepoinfo), [globalpreferences](/w/api.php?action=help&modules=query%2Bglobalpreferences), [globalrenamestatus](/w/api.php?action=help&modules=query%2Bglobalrenamestatus), [globaluserinfo](/w/api.php?action=help&modules=query%2Bglobaluserinfo), [growthmenteestatus](/w/api.php?action=help&modules=query%2Bgrowthmenteestatus), [growthmentorstatus](/w/api.php?action=help&modules=query%2Bgrowthmentorstatus), [languageinfo](/w/api.php?action=help&modules=query%2Blanguageinfo), [linterstats](/w/api.php?action=help&modules=query%2Blinterstats), [notifications](/w/api.php?action=help&modules=query%2Bnotifications), [ores](/w/api.php?action=help&modules=query%2Bores), [siteinfo](/w/api.php?action=help&modules=query%2Bsiteinfo), [siteviews](/w/api.php?action=help&modules=query%2Bsiteviews), [tokens](/w/api.php?action=help&modules=query%2Btokens), [unreadnotificationpages](/w/api.php?action=help&modules=query%2Bunreadnotificationpages), [userinfo](/w/api.php?action=help&modules=query%2Buserinfo), [wikibase](/w/api.php?action=help&modules=query%2Bwikibase), [checkuserformattedblockinfo](/w/api.php?action=help&modules=query%2Bcheckuserformattedblockinfo), [cxdeletedtranslations](/w/api.php?action=help&modules=query%2Bcxdeletedtranslations), [growthnextsuggestedtasktype](/w/api.php?action=help&modules=query%2Bgrowthnextsuggestedtasktype), [oath](/w/api.php?action=help&modules=query%2Boath), [readinglists](/w/api.php?action=help&modules=query%2Breadinglists)
indexpageids
    

Include an additional pageids section listing all returned page IDs. 

    Type: boolean ([details](/w/api.php?action=help&modules=main#main/datatype/boolean))
export
    

Export the current revisions of all given or generated pages. 

    Type: boolean ([details](/w/api.php?action=help&modules=main#main/datatype/boolean))
exportnowrap
    

Return the export XML without wrapping it in an XML result (same format as [Special:Export](/wiki/Special:Export "Special:Export")). Can only be used with query+export. 

    Type: boolean ([details](/w/api.php?action=help&modules=main#main/datatype/boolean))
exportschema
    

Target the given version of the XML dump format when exporting. Can only be used with query+export. 

    One of the following values: 0.10, 0.11
    Default: 0.11
iwurl
    

Whether to get the full URL if the title is an interwiki link. 

    Type: boolean ([details](/w/api.php?action=help&modules=main#main/datatype/boolean))
continue
    

When more results are available, use this to continue. More detailed information on how to continue queries [can be found on mediawiki.org](https://www.mediawiki.org/wiki/Special:MyLanguage/API:Continue "mw:Special:MyLanguage/API:Continue"). 

rawcontinue
    

Return raw query-continue data for continuation. 

    Type: boolean ([details](/w/api.php?action=help&modules=main#main/datatype/boolean))
titles
    

A list of titles to work on. 

    Separate values with `|` or [alternative](/w/api.php?action=help&modules=main#main/datatypes).
    Maximum number of values is 50 (500 for clients that are allowed higher limits).
pageids
    

A list of page IDs to work on. 

    Type: list of integers
    Separate values with `|` or [alternative](/w/api.php?action=help&modules=main#main/datatypes).
    Maximum number of values is 50 (500 for clients that are allowed higher limits).
revids
    

A list of revision IDs to work on. Note that almost all query modules will convert revision IDs to the corresponding page ID and work on the latest revision instead. Only `prop=revisions` uses exact revisions for its response. 

    Type: list of integers
    Separate values with `|` or [alternative](/w/api.php?action=help&modules=main#main/datatypes).
    Maximum number of values is 50 (500 for clients that are allowed higher limits).
generator
    

Get the list of pages to work on by executing the specified query module. 

**Note:** Generator parameter names must be prefixed with a "g", see examples. 

[allcategories](/w/api.php?action=help&modules=query%2Ballcategories)
    Enumerate all categories.
[alldeletedrevisions](/w/api.php?action=help&modules=query%2Balldeletedrevisions)
    List all deleted revisions by a user or in a namespace.
[allfileusages](/w/api.php?action=help&modules=query%2Ballfileusages)
    List all file usages, including non-existing.
[allimages](/w/api.php?action=help&modules=query%2Ballimages)
    Enumerate all images sequentially.
[alllinks](/w/api.php?action=help&modules=query%2Balllinks)
    Enumerate all links that point to a given namespace.
[allpages](/w/api.php?action=help&modules=query%2Ballpages)
    Enumerate all pages sequentially in a given namespace.
[allredirects](/w/api.php?action=help&modules=query%2Ballredirects)
    List all redirects to a namespace.
[allrevisions](/w/api.php?action=help&modules=query%2Ballrevisions)
    List all revisions.
[alltransclusions](/w/api.php?action=help&modules=query%2Balltransclusions)
    List all transclusions (pages embedded using {{x}}), including non-existing.
[backlinks](/w/api.php?action=help&modules=query%2Bbacklinks)
    Find all pages that link to the given page.
[categories](/w/api.php?action=help&modules=query%2Bcategories)
    List all categories the pages belong to.
[categorymembers](/w/api.php?action=help&modules=query%2Bcategorymembers)
    List all pages in a given category.
[deletedrevisions](/w/api.php?action=help&modules=query%2Bdeletedrevisions)
    Get deleted revision information.
[duplicatefiles](/w/api.php?action=help&modules=query%2Bduplicatefiles)
    List all files that are duplicates of the given files based on hash values.
[embeddedin](/w/api.php?action=help&modules=query%2Bembeddedin)
    Find all pages that embed (transclude) the given title.
[exturlusage](/w/api.php?action=help&modules=query%2Bexturlusage)
    Enumerate pages that contain a given URL.
[fileusage](/w/api.php?action=help&modules=query%2Bfileusage)
    Find all pages that use the given files.
[geosearch](/w/api.php?action=help&modules=query%2Bgeosearch)
    Returns pages having coordinates that are located in a certain area.
[images](/w/api.php?action=help&modules=query%2Bimages)
    Returns all files contained on the given pages.
[imageusage](/w/api.php?action=help&modules=query%2Bimageusage)
    Find all pages that use the given image title.
[iwbacklinks](/w/api.php?action=help&modules=query%2Biwbacklinks)
    Find all pages that link to the given interwiki link.
[langbacklinks](/w/api.php?action=help&modules=query%2Blangbacklinks)
    Find all pages that link to the given language link.
[links](/w/api.php?action=help&modules=query%2Blinks)
    Returns all links from the given pages.
[linkshere](/w/api.php?action=help&modules=query%2Blinkshere)
    Find all pages that link to the given pages.
[mostviewed](/w/api.php?action=help&modules=query%2Bmostviewed)
    Lists the most viewed pages (based on last day's pageview count).
[oldreviewedpages](/w/api.php?action=help&modules=query%2Boldreviewedpages)
    Enumerates pages that have changes pending review.
[pageswithprop](/w/api.php?action=help&modules=query%2Bpageswithprop)
    List all pages using a given page property.
[prefixsearch](/w/api.php?action=help&modules=query%2Bprefixsearch)
    Perform a prefix search for page titles.
[projectpages](/w/api.php?action=help&modules=query%2Bprojectpages)
    List all pages associated with one or more projects.
[protectedtitles](/w/api.php?action=help&modules=query%2Bprotectedtitles)
    List all titles protected from creation.
[querypage](/w/api.php?action=help&modules=query%2Bquerypage)
    Get a list provided by a QueryPage-based special page.
[random](/w/api.php?action=help&modules=query%2Brandom)
    Get a set of random pages.
[recentchanges](/w/api.php?action=help&modules=query%2Brecentchanges)
    Enumerate recent changes.
[redirects](/w/api.php?action=help&modules=query%2Bredirects)
    Returns all redirects to the given pages.
[revisions](/w/api.php?action=help&modules=query%2Brevisions)
    Get revision information.
[search](/w/api.php?action=help&modules=query%2Bsearch)
    Perform a full text search.
[templates](/w/api.php?action=help&modules=query%2Btemplates)
    Returns all pages transcluded on the given pages.
[trackingcategories](/w/api.php?action=help&modules=query%2Btrackingcategories)
    Enumerate all existing tracking categories defined in [Special:TrackingCategories](/wiki/Special:TrackingCategories "Special:TrackingCategories"). A tracking category exists if it contains pages or if its category page exists.
[transcludedin](/w/api.php?action=help&modules=query%2Btranscludedin)
    Find all pages that transclude the given pages.
[watchlist](/w/api.php?action=help&modules=query%2Bwatchlist)
    Get recent changes to pages in the current user's watchlist.
[watchlistraw](/w/api.php?action=help&modules=query%2Bwatchlistraw)
    Get all pages on the current user's watchlist.
[wblistentityusage](/w/api.php?action=help&modules=query%2Bwblistentityusage)
    Returns all pages that use the given entity IDs.
[growthtasks](/w/api.php?action=help&modules=query%2Bgrowthtasks)
    Internal. Get task recommendations suitable for newcomers.
[readinglistentries](/w/api.php?action=help&modules=query%2Breadinglistentries)
    Internal. List the pages of a certain list.
    One of the following values: [allcategories](/w/api.php?action=help&modules=query%2Ballcategories), [alldeletedrevisions](/w/api.php?action=help&modules=query%2Balldeletedrevisions), [allfileusages](/w/api.php?action=help&modules=query%2Ballfileusages), [allimages](/w/api.php?action=help&modules=query%2Ballimages), [alllinks](/w/api.php?action=help&modules=query%2Balllinks), [allpages](/w/api.php?action=help&modules=query%2Ballpages), [allredirects](/w/api.php?action=help&modules=query%2Ballredirects), [allrevisions](/w/api.php?action=help&modules=query%2Ballrevisions), [alltransclusions](/w/api.php?action=help&modules=query%2Balltransclusions), [backlinks](/w/api.php?action=help&modules=query%2Bbacklinks), [categories](/w/api.php?action=help&modules=query%2Bcategories), [categorymembers](/w/api.php?action=help&modules=query%2Bcategorymembers), [deletedrevisions](/w/api.php?action=help&modules=query%2Bdeletedrevisions), [duplicatefiles](/w/api.php?action=help&modules=query%2Bduplicatefiles), [embeddedin](/w/api.php?action=help&modules=query%2Bembeddedin), [exturlusage](/w/api.php?action=help&modules=query%2Bexturlusage), [fileusage](/w/api.php?action=help&modules=query%2Bfileusage), [geosearch](/w/api.php?action=help&modules=query%2Bgeosearch), [images](/w/api.php?action=help&modules=query%2Bimages), [imageusage](/w/api.php?action=help&modules=query%2Bimageusage), [iwbacklinks](/w/api.php?action=help&modules=query%2Biwbacklinks), [langbacklinks](/w/api.php?action=help&modules=query%2Blangbacklinks), [links](/w/api.php?action=help&modules=query%2Blinks), [linkshere](/w/api.php?action=help&modules=query%2Blinkshere), [mostviewed](/w/api.php?action=help&modules=query%2Bmostviewed), [oldreviewedpages](/w/api.php?action=help&modules=query%2Boldreviewedpages), [pageswithprop](/w/api.php?action=help&modules=query%2Bpageswithprop), [prefixsearch](/w/api.php?action=help&modules=query%2Bprefixsearch), [projectpages](/w/api.php?action=help&modules=query%2Bprojectpages), [protectedtitles](/w/api.php?action=help&modules=query%2Bprotectedtitles), [querypage](/w/api.php?action=help&modules=query%2Bquerypage), [random](/w/api.php?action=help&modules=query%2Brandom), [recentchanges](/w/api.php?action=help&modules=query%2Brecentchanges), [redirects](/w/api.php?action=help&modules=query%2Bredirects), [revisions](/w/api.php?action=help&modules=query%2Brevisions), [search](/w/api.php?action=help&modules=query%2Bsearch), [templates](/w/api.php?action=help&modules=query%2Btemplates), [trackingcategories](/w/api.php?action=help&modules=query%2Btrackingcategories), [transcludedin](/w/api.php?action=help&modules=query%2Btranscludedin), [watchlist](/w/api.php?action=help&modules=query%2Bwatchlist), [watchlistraw](/w/api.php?action=help&modules=query%2Bwatchlistraw), [wblistentityusage](/w/api.php?action=help&modules=query%2Bwblistentityusage), [growthtasks](/w/api.php?action=help&modules=query%2Bgrowthtasks), [readinglistentries](/w/api.php?action=help&modules=query%2Breadinglistentries)
redirects
    

Automatically resolve redirects in query+titles, query+pageids, and query+revids, and in pages returned by query+generator. 

    Type: boolean ([details](/w/api.php?action=help&modules=main#main/datatype/boolean))
converttitles
    

Convert titles to other variants if necessary. Only works if the wiki's content language supports variant conversion. Languages that support variant conversion include ban, crh, en, gan, iu, ku, mni, sh, shi, sr, tg, tly, uz, wuu, zgh and zh. 

    Type: boolean ([details](/w/api.php?action=help&modules=main#main/datatype/boolean))

Examples:

Fetch [site info](/w/api.php?action=help&modules=query%2Bsiteinfo) and [revisions](/w/api.php?action=help&modules=query%2Brevisions) of [Main Page](/wiki/Main_Page "Main Page").
    [api.php?action=query&prop=revisions&meta=siteinfo&titles=Main%20Page&rvprop=user|comment&continue=](/w/api.php?action=query&prop=revisions&meta=siteinfo&titles=Main%20Page&rvprop=user|comment&continue=) [[open in sandbox]](/wiki/Special:ApiSandbox#action=query&prop=revisions&meta=siteinfo&titles=Main%20Page&rvprop=user|comment&continue=)
Fetch revisions of pages beginning with `API/`.
    [api.php?action=query&generator=allpages&gapprefix=API/&prop=revisions&continue=](/w/api.php?action=query&generator=allpages&gapprefix=API/&prop=revisions&continue=) [[open in sandbox]](/wiki/Special:ApiSandbox#action=query&generator=allpages&gapprefix=API/&prop=revisions&continue=)

Retrieved from "<https://en.wikipedia.org/wiki/Special:ApiHelp>"
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
