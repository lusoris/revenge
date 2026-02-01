# MyAnimeList API

> Source: https://myanimelist.net/apiconfig/references/api/v2
> Fetched: 2026-02-01T11:44:12.627872+00:00
> Content-Hash: a38de9c011d9d658
> Type: html

---

- Versioning
- Common formats
- Common parameters
- Common status codes
- Authentication
- anime
  - getGet anime list
  - getGet anime details
  - getGet anime ranking
  - getGet seasonal anime
  - getGet suggested anime
- user animelist
  - patchUpdate my anime list status
  - delDelete my anime list item.
  - getGet user anime list
- forum
  - getGet forum boards
  - getGet forum topic detail
  - getGet forum topics
- manga
  - getGet manga list
  - getGet manga details
  - getGet manga ranking
- user mangalist
  - patchUpdate my manga list status
  - delDelete my manga list item.
  - getGet user manga list
- user
  - getGet my user information

[Documentation Powered by ReDoc](https://github.com/Redocly/redoc)

# MyAnimeList API (beta ver.) (2)

MyAnimeList.net is the property of MyAnimeList Co., Ltd. All Rights Reserved.

## Versioning

There are multiple versions of the API. You can specify a version by including it in the request uri:

(During closed beta, version starts with '0.')

`https://api.myanimelist.net/v2`

The version is incremented if any backwards incompatible changes are introduced to the API.

Examples of backwards incompatible changes are belows:

- Removing existing endpoints.
- Removing existing fields of API responses.
- Changing mean of the existing fields.

On the other hand, belows are examples of backwards compatible changes:

- Adding new endpoints.
- Adding new optional request parameters.
- Adding new fields to existing API responses.
- Changing the order of fields in existing API responses.
- Changing the contents of fields that suffixed by `_formatted`.

## Common formats

### List / Pagination

    {
      "data": [
        {...},
        {...},
        ...
      ],
      "paging": {
        "previous": "https://xxx",
        "next": "https://xxx"
      }
    }

### Date / Time format

Format | Type | Sample | Description
---|---|---|---  
date-time | string | "2015-03-02T06:03:11+00:00" | ISO 8601
date | string | "2017-10-23" or
"2017-10" or
"2017" |
time | string | "01:35" |
  
### Error format

    {
      "error": "invalid_token"
      "message": "token is invalid",
    }

## Common parameters

### Parameters for endpoints that return a list

Name | Description  
---|---  
limit |
offset |
  
### Choosing fields

By default, the API doesn’t return all fields.

You can choose the fields that you want returned with the `fields` parameter.

Example:

    fields=synopsis,my_list_status{priority,comments}

### Not Safe For Work

By default, some APIs don’t return nsfw content.

You can control this behavior with the `nsfw` parameter.

Name | Description  
---|---  
nsfw | `true` or `false`  
  
## Common status codes

Status code | Error code | Description  
---|---|---  
400 Bad Request | - | Invalid Parameters  
401 Unauthorized | invalid_token | Expired access tokens, Invalid access tokens, etc.  
403 Forbidden | - | DoS detected etc.  
404 Not Found | - |
  
# Authentication

## main_auth

Security Scheme Type |  OAuth2  
---|---  
implicit OAuth Flow | **Authorization URL:** https://myanimelist.net/v1/oauth2/authorize**Scopes:**

- `write:users` \- The API client can see and modify basic profile information and users' list data, post information to MyAnimelist on behalf of users.

## client_auth

When user login is not required, the X-MAL-CLIENT-ID request header can be used to authenticate the client by setting your API client ID.

Security Scheme Type |  API Key  
---|---  
Header parameter name:|  X-MAL-CLIENT-ID  
  
# anime

## Get anime list

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### query Parameters

q| string Search.  
---|---  
limit| integer Default: 100 The maximum value is 100.  
offset| integer Default: 0  
fields| string  
  
### Responses

**200**

OK

get/anime

https://api.myanimelist.net/v2/anime

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/anime?q=one&limit=4' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Get anime details

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### path Parameters

anime_id required | integer  
---|---  
  
##### query Parameters

fields| string  
---|---  
  
### Responses

**200**

OK

get/anime/{anime_id}

https://api.myanimelist.net/v2/anime/{anime_id}

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/anime/30230?fields=id,title,main_picture,alternative_titles,start_date,end_date,synopsis,mean,rank,popularity,num_list_users,num_scoring_users,nsfw,created_at,updated_at,media_type,status,genres,my_list_status,num_episodes,start_season,broadcast,source,average_episode_duration,rating,pictures,background,related_anime,related_manga,recommendations,studios,statistics' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Get anime ranking

The returned anime contains the `ranking` field.

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### query Parameters

ranking_type required | string | value |
---|---  
all | Top Anime Series  
airing | Top Airing Anime  
upcoming | Top Upcoming Anime  
tv | Top Anime TV Series  
ova | Top Anime OVA Series  
movie | Top Anime Movies  
special | Top Anime Specials  
bypopularity | Top Anime by Popularity  
favorite | Top Favorited Anime  
  
limit| integer Default: 100 The maximum value is 500.  
offset| integer Default: 0  
fields| string  
  
### Responses

**200**

OK

get/anime/ranking

https://api.myanimelist.net/v2/anime/ranking

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/anime/ranking?ranking_type=all&limit=4' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Get seasonal anime

Get seasonal anime.

Season name | Months  
---|---  
winter | January, February, March  
spring | April, May, June  
summer | July, August, September  
fall | October, November, December  
  
##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### path Parameters

year required | integer  
---|---  
season required | string  
  
##### query Parameters

sort| string Valid values: | Value | Order  
---|---  
`anime_score` | Descending  
`anime_num_list_users` | Descending  
  
limit| integer Default: 100 The maximum value is 500.  
offset| integer Default: 0  
fields| string  
  
### Responses

**200**

OK

get/anime/season/{year}/{season}

https://api.myanimelist.net/v2/anime/season/{year}/{season}

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/anime/season/2017/summer?limit=4' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Get suggested anime

Returns suggested anime for the authorized user.

If the user is new comer, this endpoint returns an empty list.

##### Authorizations

main_auth (`write:users`)

##### query Parameters

limit| integer Default: 100 The maximum value is 100.  
---|---  
offset| integer Default: 0  
fields| string  
  
### Responses

**200**

OK.

get/anime/suggestions

https://api.myanimelist.net/v2/anime/suggestions

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/anime/suggestions?limit=4' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

# user animelist

## Update my anime list status

Add specified anime to my anime list.

If specified anime already exists, update its status.

This endpoint updates only values specified by the parameter.

##### Authorizations

main_auth (`write:users`)

##### path Parameters

anime_id required | integer  
---|---  
  
##### Request Body schema: application/x-www-form-urlencoded

status| string

- watching
- completed
- on_hold
- dropped
- plan_to_watch

---|---  
is_rewatching| boolean  
score| integer 0-10  
num_watched_episodes| integer  
priority| integer 0-2  
num_times_rewatched| integer  
rewatch_value| integer 0-5  
tags| string  
comments| string  
  
### Responses

**200**

OK

patch/anime/{anime_id}/my_list_status

https://api.myanimelist.net/v2/anime/{anime_id}/my_list_status

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/anime/17074/my_list_status' \
    -X PUT \
    -d status=completed \
    -d score=8 \
    -d num_watched_episodes=3 \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Delete my anime list item

If the specified anime does not exist in user's anime list, this endpoint does nothing and returns `404 Not Found`.

So be careful when retrying.

##### Authorizations

main_auth (`write:users`)

##### path Parameters

anime_id required | integer  
---|---  
  
### Responses

**200**

OK

**404**

Not Found

delete/anime/{anime_id}/my_list_status

https://api.myanimelist.net/v2/anime/{anime_id}/my_list_status

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/anime/21/my_list_status' \
    -X DELETE \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200
- 404

Content type

application/json

Copy

Expand all  Collapse all

`null`

## Get user anime list

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### path Parameters

user_name required | string User name or `@me`.  
---|---  
  
##### query Parameters

status| string Filters returned anime list by these statuses. To return all anime, don't specify this field. Valid values:

- watching
- completed
- on_hold
- dropped
- plan_to_watch

---|---  
sort| string Valid values: | Value | Order  
---|---  
`list_score` | Descending  
`list_updated_at` | Descending  
`anime_title` | Ascending  
`anime_start_date` | Descending  
`anime_id` (Under Development) | Ascending  
  
limit| integer Default: 100 The maximum value is 1000.  
offset| integer Default: 0  
  
### Responses

**200**

OK

get/users/{user_name}/animelist

https://api.myanimelist.net/v2/users/{user_name}/animelist

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/users/@me/animelist?fields=list_status&limit=4' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

# forum

## Get forum boards

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

### Responses

**200**

OK

get/forum/boards

https://api.myanimelist.net/v2/forum/boards

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/forum/boards' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Get forum topic detail

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### path Parameters

topic_id required | integer  
---|---  
  
##### query Parameters

limit| integer <= 100 Default: 100  
---|---  
offset| integer Default: 0  
  
### Responses

**200**

OK

get/forum/topic/{topic_id}

https://api.myanimelist.net/v2/forum/topic/{topic_id}

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/forum/topic/481' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Get forum topics

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### query Parameters

board_id| integer  
---|---  
subboard_id| integer  
limit| integer <= 100 Default: 100  
offset| integer Default: 0  
sort| string Default: "recent" Currently, only "recent" can be set.  
q| string  
topic_user_name| string  
user_name| string  
  
### Responses

**200**

OK

get/forum/topics

https://api.myanimelist.net/v2/forum/topics

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/forum/topics?q=love&subboard_id=2&limit=10' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

# manga

## Get manga list

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### query Parameters

q| string Search.  
---|---  
limit| integer Default: 100 The maximum value is 100.  
offset| integer Default: 0  
fields| string  
  
### Responses

**200**

OK

get/manga

https://api.myanimelist.net/v2/manga

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/manga?q=berserk' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Get manga details

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### path Parameters

manga_id required | integer  
---|---  
  
##### query Parameters

fields| string  
---|---  
  
### Responses

**200**

OK

get/manga/{manga_id}

https://api.myanimelist.net/v2/manga/{manga_id}

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/manga/2?fields=id,title,main_picture,alternative_titles,start_date,end_date,synopsis,mean,rank,popularity,num_list_users,num_scoring_users,nsfw,created_at,updated_at,media_type,status,genres,my_list_status,num_volumes,num_chapters,authors{first_name,last_name},pictures,background,related_anime,related_manga,recommendations,serialization{name}' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Get manga ranking

The returned manga contains the `ranking` field.

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### query Parameters

ranking_type required | string | value |
---|---  
all | All  
manga | Top Manga  
novels | Top Novels  
oneshots | Top One-shots  
doujin | Top Doujinshi  
manhwa | Top Manhwa  
manhua | Top Manhua  
bypopularity | Most Popular  
favorite | Most Favorited  
  
limit| integer Default: 100 The maximum value is 500.  
offset| integer Default: 0  
fields| string  
  
### Responses

**200**

OK

get/manga/ranking

https://api.myanimelist.net/v2/manga/ranking

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/manga/ranking?ranking_type=all&limit=4' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

# user mangalist

## Update my manga list status

Add specified manga to my manga list.

If specified manga already exists, update its status.

This endpoint updates only values specified by the parameter.

##### Authorizations

main_auth (`write:users`)

##### path Parameters

manga_id required | integer  
---|---  
  
##### Request Body schema: application/x-www-form-urlencoded

status| string

- reading
- completed
- on_hold
- dropped
- plan_to_read

---|---  
is_rereading| boolean  
score| integer 0-10  
num_volumes_read| integer  
num_chapters_read| integer  
priority| integer 0-2  
num_times_reread| integer  
reread_value| integer 0-5  
tags| string  
comments| string  
  
### Responses

**200**

OK

patch/manga/{manga_id}/my_list_status

https://api.myanimelist.net/v2/manga/{manga_id}/my_list_status

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/manga/2/my_list_status' \
    -X PUT \
    -d status=completed \
    -d score=8 \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

## Delete my manga list item

If the specified manga does not exist in user's manga list, this endpoint does nothing and returns `404 Not Found`.

So be careful when retrying.

##### Authorizations

main_auth (`write:users`)

##### path Parameters

manga_id required | integer  
---|---  
  
### Responses

**200**

OK

**404**

Not Found

delete/manga/{manga_id}/my_list_status

https://api.myanimelist.net/v2/manga/{manga_id}/my_list_status

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/manga/2/my_list_status' \
    -X DELETE \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200
- 404

Content type

application/json

Copy

Expand all  Collapse all

`null`

## Get user manga list

##### Authorizations

main_auth (`write:users`) client_auth (`-`)

##### path Parameters

user_name required | string User name or `@me`.  
---|---  
  
##### query Parameters

status| string Filters returned manga list by these statuses. To return all manga, don't specify this field. Valid values:

- reading
- completed
- on_hold
- dropped
- plan_to_read

---|---  
sort| string Valid values: | Value | Order  
---|---  
`list_score` | Descending  
`list_updated_at` | Descending  
`manga_title` | Ascending  
`manga_start_date` | Descending  
`manga_id` (Under Development) | Ascending  
  
limit| integer Default: 100 The maximum value is 1000.  
offset| integer Default: 0  
  
### Responses

**200**

OK

get/users/{user_name}/mangalist

https://api.myanimelist.net/v2/users/{user_name}/mangalist

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/users/@me/mangalist?fields=list_status&limit=4' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample

# user

## Get my user information

##### Authorizations

main_auth (`write:users`)

##### path Parameters

user_id required | string You can only specify `@me`.  
---|---  
  
##### query Parameters

fields| string  
---|---  
  
### Responses

**200**

OK

get/users/{user_name}

https://api.myanimelist.net/v2/users/{user_name}

### Request samples

- curl

Copy

    curl 'https://api.myanimelist.net/v2/users/@me?fields=anime_statistics' \
    -H 'Authorization: Bearer YOUR_TOKEN'
    

### Response samples

- 200

Content type

*/*

*/*

application/json

No sample
