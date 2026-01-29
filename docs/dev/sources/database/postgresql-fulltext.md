# PostgreSQL Full Text Search

> Auto-fetched from [https://www.postgresql.org/docs/current/textsearch.html](https://www.postgresql.org/docs/current/textsearch.html)
> Last Updated: 2026-01-29T20:12:31.133993+00:00

---

Home
About
Download
Documentation
Community
Developers
Support
Donate
Your account
November 13, 2025:
PostgreSQL 18.1, 17.7, 16.11, 15.15, 14.20, and 13.23 Released!
Documentation
→
PostgreSQL 18
Supported Versions:
Current
(
18
)

/
17
/
16
/
15
/
14
Development Versions:
devel
Unsupported versions:
13
/
12
/
11
/
10
/
9.6
/
9.5
/
9.4
/
9.3
/
9.2
/
9.1
/
9.0
/
8.4
/
8.3
Chapter 12. Full Text Search
Prev
Up
Part II. The SQL Language
Home
Next
Chapter 12. Full Text Search
Table of Contents
12.1. Introduction
12.1.1. What Is a Document?
12.1.2. Basic Text Matching
12.1.3. Configurations
12.2. Tables and Indexes
12.2.1. Searching a Table
12.2.2. Creating Indexes
12.3. Controlling Text Search
12.3.1. Parsing Documents
12.3.2. Parsing Queries
12.3.3. Ranking Search Results
12.3.4. Highlighting Results
12.4. Additional Features
12.4.1. Manipulating Documents
12.4.2. Manipulating Queries
12.4.3. Triggers for Automatic Updates
12.4.4. Gathering Document Statistics
12.5. Parsers
12.6. Dictionaries
12.6.1. Stop Words
12.6.2. Simple Dictionary
12.6.3. Synonym Dictionary
12.6.4. Thesaurus Dictionary
12.6.5.
Ispell
Dictionary
12.6.6.
Snowball
Dictionary
12.7. Configuration Example
12.8. Testing and Debugging Text Search
12.8.1. Configuration Testing
12.8.2. Parser Testing
12.8.3. Dictionary Testing
12.9. Preferred Index Types for Text Search
12.10.
psql
Support
12.11. Limitations
Prev
Up
Next
11.12. Examining Index Usage
Home
12.1. Introduction
Submit correction
If you see anything in the documentation that is not correct, does not match
your experience with the particular feature or requires further clarification,
please use
this form
to report a documentation issue.
Policies
|
Code of Conduct
|
About PostgreSQL
|
Contact
Copyright © 1996-2026 The PostgreSQL Global Development Group