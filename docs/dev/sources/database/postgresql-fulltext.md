# PostgreSQL Full Text Search

> Source: https://www.postgresql.org/docs/current/textsearch.html
> Fetched: 2026-02-01T11:43:38.303277+00:00
> Content-Hash: cc217bdaa43ae38d
> Type: html

---

November 13, 2025: [PostgreSQL 18.1, 17.7, 16.11, 15.15, 14.20, and 13.23 Released!](/about/news/postgresql-181-177-1611-1515-1420-and-1323-released-3171/)

[Documentation](/docs/ "Documentation") → [PostgreSQL 18](/docs/18/index.html)

Supported Versions: [Current](/docs/current/textsearch.html "PostgreSQL 18 - Chapter 12. Full Text Search") ([18](/docs/18/textsearch.html "PostgreSQL 18 - Chapter 12. Full Text Search")) / [17](/docs/17/textsearch.html "PostgreSQL 17 - Chapter 12. Full Text Search") / [16](/docs/16/textsearch.html "PostgreSQL 16 - Chapter 12. Full Text Search") / [15](/docs/15/textsearch.html "PostgreSQL 15 - Chapter 12. Full Text Search") / [14](/docs/14/textsearch.html "PostgreSQL 14 - Chapter 12. Full Text Search")

Development Versions: [devel](/docs/devel/textsearch.html "PostgreSQL devel - Chapter 12. Full Text Search")

Unsupported versions: [13](/docs/13/textsearch.html "PostgreSQL 13 - Chapter 12. Full Text Search") / [12](/docs/12/textsearch.html "PostgreSQL 12 - Chapter 12. Full Text Search") / [11](/docs/11/textsearch.html "PostgreSQL 11 - Chapter 12. Full Text Search") / [10](/docs/10/textsearch.html "PostgreSQL 10 - Chapter 12. Full Text Search") / [9.6](/docs/9.6/textsearch.html "PostgreSQL 9.6 - Chapter 12. Full Text Search") / [9.5](/docs/9.5/textsearch.html "PostgreSQL 9.5 - Chapter 12. Full Text Search") / [9.4](/docs/9.4/textsearch.html "PostgreSQL 9.4 - Chapter 12. Full Text Search") / [9.3](/docs/9.3/textsearch.html "PostgreSQL 9.3 - Chapter 12. Full Text Search") / [9.2](/docs/9.2/textsearch.html "PostgreSQL 9.2 - Chapter 12. Full Text Search") / [9.1](/docs/9.1/textsearch.html "PostgreSQL 9.1 - Chapter 12. Full Text Search") / [9.0](/docs/9.0/textsearch.html "PostgreSQL 9.0 - Chapter 12. Full Text Search") / [8.4](/docs/8.4/textsearch.html "PostgreSQL 8.4 - Chapter 12. Full Text Search") / [8.3](/docs/8.3/textsearch.html "PostgreSQL 8.3 - Chapter 12. Full Text Search")

__

Chapter 12. Full Text Search  
---  
[Prev](indexes-examine.html "11.12. Examining Index Usage") | [Up](sql.html "Part II. The SQL Language") | Part II. The SQL Language | [Home](index.html "PostgreSQL 18.1 Documentation") |  [Next](textsearch-intro.html "12.1. Introduction")  
  
* * *

## Chapter 12. Full Text Search

**Table of Contents**

[12.1. Introduction](textsearch-intro.html)

[12.1.1. What Is a Document?](textsearch-intro.html#TEXTSEARCH-DOCUMENT)
[12.1.2. Basic Text Matching](textsearch-intro.html#TEXTSEARCH-MATCHING)
[12.1.3. Configurations](textsearch-intro.html#TEXTSEARCH-INTRO-CONFIGURATIONS)
[12.2. Tables and Indexes](textsearch-tables.html)

[12.2.1. Searching a Table](textsearch-tables.html#TEXTSEARCH-TABLES-SEARCH)
[12.2.2. Creating Indexes](textsearch-tables.html#TEXTSEARCH-TABLES-INDEX)
[12.3. Controlling Text Search](textsearch-controls.html)

[12.3.1. Parsing Documents](textsearch-controls.html#TEXTSEARCH-PARSING-DOCUMENTS)
[12.3.2. Parsing Queries](textsearch-controls.html#TEXTSEARCH-PARSING-QUERIES)
[12.3.3. Ranking Search Results](textsearch-controls.html#TEXTSEARCH-RANKING)
[12.3.4. Highlighting Results](textsearch-controls.html#TEXTSEARCH-HEADLINE)
[12.4. Additional Features](textsearch-features.html)

[12.4.1. Manipulating Documents](textsearch-features.html#TEXTSEARCH-MANIPULATE-TSVECTOR)
[12.4.2. Manipulating Queries](textsearch-features.html#TEXTSEARCH-MANIPULATE-TSQUERY)
[12.4.3. Triggers for Automatic Updates](textsearch-features.html#TEXTSEARCH-UPDATE-TRIGGERS)
[12.4.4. Gathering Document Statistics](textsearch-features.html#TEXTSEARCH-STATISTICS)
[12.5. Parsers](textsearch-parsers.html)
[12.6. Dictionaries](textsearch-dictionaries.html)

[12.6.1. Stop Words](textsearch-dictionaries.html#TEXTSEARCH-STOPWORDS)
[12.6.2. Simple Dictionary](textsearch-dictionaries.html#TEXTSEARCH-SIMPLE-DICTIONARY)
[12.6.3. Synonym Dictionary](textsearch-dictionaries.html#TEXTSEARCH-SYNONYM-DICTIONARY)
[12.6.4. Thesaurus Dictionary](textsearch-dictionaries.html#TEXTSEARCH-THESAURUS)
[12.6.5. Ispell Dictionary](textsearch-dictionaries.html#TEXTSEARCH-ISPELL-DICTIONARY)
[12.6.6. Snowball Dictionary](textsearch-dictionaries.html#TEXTSEARCH-SNOWBALL-DICTIONARY)
[12.7. Configuration Example](textsearch-configuration.html)
[12.8. Testing and Debugging Text Search](textsearch-debugging.html)

[12.8.1. Configuration Testing](textsearch-debugging.html#TEXTSEARCH-CONFIGURATION-TESTING)
[12.8.2. Parser Testing](textsearch-debugging.html#TEXTSEARCH-PARSER-TESTING)
[12.8.3. Dictionary Testing](textsearch-debugging.html#TEXTSEARCH-DICTIONARY-TESTING)
[12.9. Preferred Index Types for Text Search](textsearch-indexes.html)
[12.10. psql Support](textsearch-psql.html)
[12.11. Limitations](textsearch-limitations.html)

* * *

[Prev](indexes-examine.html "11.12. Examining Index Usage") | [Up](sql.html "Part II. The SQL Language") |  [Next](textsearch-intro.html "12.1. Introduction")  
---|---|---  
11.12. Examining Index Usage  | [Home](index.html "PostgreSQL 18.1 Documentation") |  12.1. Introduction  
  
## Submit correction

If you see anything in the documentation that is not correct, does not match your experience with the particular feature or requires further clarification, please use [this form](/account/comments/new/18/textsearch.html/) to report a documentation issue.
