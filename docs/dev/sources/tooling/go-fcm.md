# appleboy/go-fcm

> Source: https://pkg.go.dev/github.com/appleboy/go-fcm
> Fetched: 2026-01-30T23:49:54.222180+00:00
> Content-Hash: f078de80a93c5c73
> Type: html

---

Overview

¶

Package fcm provides Firebase Cloud Messaging functionality for Golang

Here is a simple example illustrating how to use FCM library:

func main() {
	ctx := context.Background()
	client, err := fcm.NewClient(
		ctx,
		fcm.WithCredentialsFile("path/to/serviceAccountKey.json"),
	)
	if err != nil {
		log.Fatal(err)
	}

// Send to a single device
token := "test"
resp, err := client.Send(

ctx,
&messaging.Message{
	Token: token,
	Data: map[string]string{
		"foo": "bar",
	},
},

)

if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success count:", resp.SuccessCount)
	fmt.Println("failure count:", resp.FailureCount)
	fmt.Println("message id:", resp.Responses[0].MessageID)
	fmt.Println("error msg:", resp.Responses[0].Error)
}

Index

¶

type Client

func NewClient(ctx context.Context, opts ...Option) (*Client, error)

func (c *Client) Send(ctx context.Context, message ...*messaging.Message) (*messaging.BatchResponse, error)

func (c *Client) SendDryRun(ctx context.Context, message ...*messaging.Message) (*messaging.BatchResponse, error)

func (c *Client) SendMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)

func (c *Client) SendMulticastDryRun(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)

func (c *Client) SubscribeTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)

func (c *Client) UnsubscribeTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)

type Option

func WithCredentialsFile(filename string) Option

func WithCredentialsJSON(data []byte) Option

func WithCustomClientOption(opts ...option.ClientOption) Option

func WithDebug(debug bool) Option

func WithEndpoint(endpoint string) Option

func WithHTTPClient(httpClient *http.Client) Option

func WithHTTPProxy(proxyURL string) Option

func WithProjectID(projectID string) Option

func WithServiceAccount(serviceAccount string) Option

func WithTokenSource(s oauth2.TokenSource) Option

Constants

¶

This section is empty.

Variables

¶

This section is empty.

Functions

¶

This section is empty.

Types

¶

type

Client

¶

type Client struct {

// contains filtered or unexported fields

}

Client abstracts the interaction between the application server and the
FCM server via HTTP protocol. The developer must obtain an API key from the
Google APIs Console page and pass it to the `Client` so that it can
perform authorized requests on the application server's behalf.
To send a message to one or more devices use the Client's Send.

If the `HTTP` field is nil, a zeroed http.Client will be allocated and used
to send messages.

Authorization Scopes
Requires one of the following OAuth scopes:
-

https://www.googleapis.com/auth/firebase.messaging

func

NewClient

¶

func NewClient(ctx

context

.

Context

, opts ...

Option

) (*

Client

,

error

)

NewClient creates new Firebase Cloud Messaging Client based on API key and
with default endpoint and http client.

func (*Client)

Send

¶

func (c *

Client

) Send(
	ctx

context

.

Context

,
	message ...*

messaging

.

Message

,
) (*

messaging

.

BatchResponse

,

error

)

SendWithContext sends a message to the FCM server without retrying in case of service
unavailability. A non-nil error is returned if a non-recoverable error
occurs (i.e. if the response status is not "200 OK").
Behaves just like regular send, but uses external context.

func (*Client)

SendDryRun

¶

added in

v1.0.0

func (c *

Client

) SendDryRun(
	ctx

context

.

Context

,
	message ...*

messaging

.

Message

,
) (*

messaging

.

BatchResponse

,

error

)

SendDryRun sends the messages in the given array via Firebase Cloud Messaging in the
dry run (validation only) mode.

func (*Client)

SendMulticast

¶

added in

v1.0.0

func (c *

Client

) SendMulticast(
	ctx

context

.

Context

,
	message *

messaging

.

MulticastMessage

,
) (*

messaging

.

BatchResponse

,

error

)

SendEachForMulticast sends the given multicast message to all the FCM registration tokens specified.

func (*Client)

SendMulticastDryRun

¶

added in

v1.0.0

func (c *

Client

) SendMulticastDryRun(
	ctx

context

.

Context

,
	message *

messaging

.

MulticastMessage

,
) (*

messaging

.

BatchResponse

,

error

)

SendEachForMulticastDryRun sends the given multicast message to all the specified FCM registration
tokens in the dry run (validation only) mode.

func (*Client)

SubscribeTopic

¶

added in

v1.0.0

func (c *

Client

) SubscribeTopic(
	ctx

context

.

Context

,
	tokens []

string

,
	topic

string

,
) (*

messaging

.

TopicManagementResponse

,

error

)

SubscribeToTopic subscribes a list of registration tokens to a topic.

The tokens list must not be empty, and have at most 1000 tokens.

func (*Client)

UnsubscribeTopic

¶

added in

v1.0.0

func (c *

Client

) UnsubscribeTopic(
	ctx

context

.

Context

,
	tokens []

string

,
	topic

string

,
) (*

messaging

.

TopicManagementResponse

,

error

)

UnsubscribeFromTopic unsubscribes a list of registration tokens from a topic.

The tokens list must not be empty, and have at most 1000 tokens.

type

Option

¶

type Option func(*

Client

)

error

Option configurates Client with defined option.

func

WithCredentialsFile

¶

added in

v1.0.0

func WithCredentialsFile(filename

string

)

Option

WithCredentialsFile returns a ClientOption that authenticates
API calls with the given service account or refresh token JSON
credentials file.

func

WithCredentialsJSON

¶

added in

v1.0.0

func WithCredentialsJSON(data []

byte

)

Option

WithCredentialsJSON returns a ClientOption that authenticates
API calls with the given service account or refresh token JSON
credentials.

func

WithCustomClientOption

¶

added in

v1.2.0

func WithCustomClientOption(opts ...

option

.

ClientOption

)

Option

WithCustomClientOption is an option function that allows you to provide custom client options.
It appends the provided custom options to the client's options list.
The custom options are applied when sending requests to the FCM server.
If no custom options are provided, this function does nothing.

Parameters:

opts: The custom client options to be appended to the client's options list.

Returns:

An error if there was an issue appending the custom options to the client's options list, or nil otherwise.

func

WithDebug

¶

added in

v1.1.0

func WithDebug(debug

bool

)

Option

WithDebug returns Option to configure debug mode.

func

WithEndpoint

¶

func WithEndpoint(endpoint

string

)

Option

WithEndpoint returns Option to configure endpoint.

func

WithHTTPClient

¶

func WithHTTPClient(httpClient *

http

.

Client

)

Option

WithHTTPClient returns Option to configure HTTP Client.

func

WithHTTPProxy

¶

added in

v0.1.6

func WithHTTPProxy(proxyURL

string

)

Option

WithHTTPProxy returns Option to configure HTTP Client proxy.

func

WithProjectID

¶

added in

v1.0.0

func WithProjectID(projectID

string

)

Option

WithProjectID returns Option to configure project ID.

func

WithServiceAccount

¶

added in

v1.0.0

func WithServiceAccount(serviceAccount

string

)

Option

WithServiceAccount returns Option to configure service account.

func

WithTokenSource

¶

added in

v1.0.0

func WithTokenSource(s

oauth2

.

TokenSource

)

Option

WithTokenSource returns a ClientOption that specifies an OAuth2 token
source to be used as the basis for authentication.