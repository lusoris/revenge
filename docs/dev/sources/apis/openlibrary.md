# Open Library API

> Source: https://openlibrary.org/developers/api
> Fetched: 2026-01-30T23:50:40.394173+00:00
> Content-Hash: 627177c6a1a4e763
> Type: html

---

Last edited by

The Poison Path Collection

3 days ago |

History

Edit

APIs

Developer Center (Home)

|

Web API

|

Client Library

|

Data Dumps

|

Source Code

|

Report an Issue

|

Licensing

Developer Center

→ APIs

Open Library offers a suite of APIs to help developers get up and running with our data. This includes

RESTful APIs

, which make Open Library data availabile in JSON, YAML and RDF/XML formats. There's also an earlier, now deprecated

JSON API

which is preserved for backward compatibility.

Try out our interactive

OpenAPI sandbox

!

Identifying your Application

If your application will make regular, frequent use of Open Library's APIs (e.g. multiple calls per minute), please add a HEADER that specifies a User-Agent string with (a) the name of your application and (b) your contact email or phone number, so we may contact you when we notice high request volume. Failing to include these headers may result in your application being blocked.

For help adding user-agent headers to your API calls, please refer to this

javascript

and

python

example.

Additional rate limit information can be accessed

here

.

Index of APIs

Book Search

API - Search results for books, authors, and more

Work & Edition

APIs - Retrieve a specific work or edition by identifier

My Books

API - Retrieve books on a patron's public reading log

Authors

API - Retrieve an author and their works by author identifier

Subjects

API - Fetch books by subject name

Search inside

API - Search for matching text within millions of books

Partner

API -- Formerly the "Read" API, fetch one or more books by library identifiers (ISBNs, OCLC, LCCNs)

Covers

API - Fetch book covers by ISBN or Open Library identifier

Recent Changes

API - Programatic access to changes across Open Library

Lists

API - Reading, modifying, or creating user lists

History of an item can be accessed by appending

?m=history

to the page like

this

.

Bulk Access

Please do not use our APIs for bulk download

of Open Library data because this affects our ability to serve patrons. We make our data

publicly available

each month for partners. If you want a dump of complete data, please read about our

Bulk Download

options, or email us at

openlibrary@archive.org

.

More APIs

Did you know, nearly every page on Open Library is or has an API. You can return structured bibliographic data for any page by adding a .rdf/.json/.yml extension to the end of any Open Library identifier. For instance:

https://openlibrary.org/works/OL15626917W.json

or

https://openlibrary.org/authors/OL33421A.json

. Many pages, such as the Books, Authors, and Lists, will include links to their RDF and JSON formats.

Questions

We encourage developers to ask questions by opening issues on

GitHub

and on our

gitter

chat channel.

Friends using Open Library APIs

Several developers are creating amazing things with the Open Library APIs:

Trove

by

the National Library of Australia

Trove is a new discovery experience focused on Australia and Australians. It supplements what search engines provide with reliable information from Australia's memory institutions. The system hits Open Library when public domain books turn up in searches, and displays links to Open Library.

Koha

Koha is an open source library system for public libraries that includes catalog searches and member organizing. It uses Open Library covers, displays OL related subjects, and lendable eBooks using the Read API.

Evergreen

Evergreen is highly-scalable software for libraries that helps library patrons find library materials, and helps libraries manage, catalog, and circulate those materials. It uses Open Library for covers, tables of contents, with plans to expand into other areas.

read.gov

by

the Library of Congress

OK, this isn't exactly Open Library, but it's still awesome! The Library of Congress have modified the Internet Archive's Book Reader to sit perfectly within their Rare Books Collection site.

OpenBook WordPress Plug-in

by

John Miedema

OpenBook is useful for anyone who wants to add book covers and other book data on a WordPress website. OpenBook links to detailed book information in Open Library, the main data source, as well as other book sites. Users have complete control over the display through templates. OpenBook can link to library records by configuring an OpenURL resolver or through a

WorldCat

link. OpenBook inserts

COinS

so that other applications like

Zotero

can pick up the book data.

Umlaut

by

Jason Ronallo

Umlaut is a middle-tier OpenURL link resolver that adds functions and services to commercial link resolving software.

Virtual Shelf

by

Jonathan Breitbart and Devin Blong

(UC Berkeley School of Information)

The Virtual Shelf is a visualization created by two students at the UC Berkeley School of Information. The project includes the student's master thesis, with research into the searching and browsing patterns of library patrons. The Open Library RESTful API was utilized during the project as a source of metadata for the user interface.

RDC UI Toolkit

by

Rural Design Collective

This group created a suite of tools that facilitates the creation of localized user interfaces for public domain books. The RDC used the Open Library Covers API and the Internet Archive Book Reader in their online demonstration customized for the OLPC XO.

Dreambooks.club

by

Bernat Fortet

Dreambooks is a portal and community where parents and children can discover new books to read together. Think of it as the online equivalent of your library's children's corner. All the book data is powered by OpenLibary's API.

MyBooks.Life

by

Mark Webster

MyBooks.Life is an android app and website designed primarily to manage TBR (to-be-read) lists. You can keep track of your reading progress, make notes, manage your wishlist, and rate your books. MyBooks.Life uses Open Library data to power its search.

Bookmind

Bookmind is now available at

https://apps.apple.com/app/bookmind/id6593662584

. It uses open library’s api exclusively for book data. You can even see the rough prototype source at

https://github.com/dave-ruest/Bookmind

.

Hobbyverse

let's you track all your hobbies in one place. Users can add their books to their digital library and track their progress reading books, view what books their friends are reading, earn achievements, etc.

ReadOtter

ReadOtter, a classroom library management app designed to help teachers organize their classroom libraries.

Chapter

Chapter is an online reading library and reading organizer app.

Land of Readers

is a free, easy-to-use book discovery tool designed to help readers find books that match their interests, age group, and reading level.

Austen

is a web app that uses the Open Library API to generate visual character relationship diagrams for books using AI.

Here

is a live demo of their work.

mcp-open-library

by

Ben Smith

A Model Context Protocol (MCP) server for the Open Library API that enables AI assistants to search for book and author information. The source code can be found

here

. So far on a MCP server website its been called 5.2k times:

https://smithery.ai/server/@8enSmith/mcp-open-library

. It has also been published to npm:

https://www.npmjs.com/package/mcp-open-library

.

Are you using the Open Library APIs? We'd love to hear about it! Please email us at

openlibrary@archive.org

.

History

Created November 12, 2009

86 revisions

3 days ago

Edited by

The Poison Path Collection

Added mcp-open-library as one of friends using OL APIs.

July 18, 2025

Edited by

The Poison Path Collection

Corrected a formatting error with Austen.

July 18, 2025

Edited by

The Poison Path Collection

Added Austen to the list of developers with appropriate formatting.

July 16, 2025

Edited by

The Poison Path Collection

Added three friends of OL at the end and cleaned up one link

November 12, 2009

Created by

George

Building out the sitemap