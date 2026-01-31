# appleboy/go-fcm

> Source: https://pkg.go.dev/github.com/appleboy/go-fcm
> Fetched: 2026-01-31T10:57:14.608151+00:00
> Content-Hash: 424be5f1e1c1f8b9
> Type: html

---

### Overview ¶

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
    

### Index ¶

  * type Client
  *     * func NewClient(ctx context.Context, opts ...Option) (*Client, error)
  *     * func (c *Client) Send(ctx context.Context, message ...*messaging.Message) (*messaging.BatchResponse, error)
    * func (c *Client) SendDryRun(ctx context.Context, message ...*messaging.Message) (*messaging.BatchResponse, error)
    * func (c *Client) SendMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
    * func (c *Client) SendMulticastDryRun(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
    * func (c *Client) SubscribeTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)
    * func (c *Client) UnsubscribeTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)
  * type Option
  *     * func WithCredentialsFile(filename string) Option
    * func WithCredentialsJSON(data []byte) Option
    * func WithCustomClientOption(opts ...option.ClientOption) Option
    * func WithDebug(debug bool) Option
    * func WithEndpoint(endpoint string) Option
    * func WithHTTPClient(httpClient *http.Client) Option
    * func WithHTTPProxy(proxyURL string) Option
    * func WithProjectID(projectID string) Option
    * func WithServiceAccount(serviceAccount string) Option
    * func WithTokenSource(s oauth2.TokenSource) Option



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

This section is empty.

### Types ¶

####  type [Client](https://github.com/appleboy/go-fcm/blob/v1.2.7/client.go#L31) ¶
    
    
    type Client struct {
    	// contains filtered or unexported fields
    }

Client abstracts the interaction between the application server and the FCM server via HTTP protocol. The developer must obtain an API key from the Google APIs Console page and pass it to the `Client` so that it can perform authorized requests on the application server's behalf. To send a message to one or more devices use the Client's Send. 

If the `HTTP` field is nil, a zeroed http.Client will be allocated and used to send messages. 

Authorization Scopes Requires one of the following OAuth scopes: \- <https://www.googleapis.com/auth/firebase.messaging>

####  func [NewClient](https://github.com/appleboy/go-fcm/blob/v1.2.7/client.go#L43) ¶
    
    
    func NewClient(ctx [context](/context).[Context](/context#Context), opts ...Option) (*Client, [error](/builtin#error))

NewClient creates new Firebase Cloud Messaging Client based on API key and with default endpoint and http client. 

####  func (*Client) [Send](https://github.com/appleboy/go-fcm/blob/v1.2.7/client.go#L111) ¶
    
    
    func (c *Client) Send(
    	ctx [context](/context).[Context](/context#Context),
    	message ...*[messaging](/firebase.google.com/go/v4/messaging).[Message](/firebase.google.com/go/v4/messaging#Message),
    ) (*[messaging](/firebase.google.com/go/v4/messaging).[BatchResponse](/firebase.google.com/go/v4/messaging#BatchResponse), [error](/builtin#error))

SendWithContext sends a message to the FCM server without retrying in case of service unavailability. A non-nil error is returned if a non-recoverable error occurs (i.e. if the response status is not "200 OK"). Behaves just like regular send, but uses external context. 

####  func (*Client) [SendDryRun](https://github.com/appleboy/go-fcm/blob/v1.2.7/client.go#L125) ¶ added in v1.0.0
    
    
    func (c *Client) SendDryRun(
    	ctx [context](/context).[Context](/context#Context),
    	message ...*[messaging](/firebase.google.com/go/v4/messaging).[Message](/firebase.google.com/go/v4/messaging#Message),
    ) (*[messaging](/firebase.google.com/go/v4/messaging).[BatchResponse](/firebase.google.com/go/v4/messaging#BatchResponse), [error](/builtin#error))

SendDryRun sends the messages in the given array via Firebase Cloud Messaging in the dry run (validation only) mode. 

####  func (*Client) [SendMulticast](https://github.com/appleboy/go-fcm/blob/v1.2.7/client.go#L138) ¶ added in v1.0.0
    
    
    func (c *Client) SendMulticast(
    	ctx [context](/context).[Context](/context#Context),
    	message *[messaging](/firebase.google.com/go/v4/messaging).[MulticastMessage](/firebase.google.com/go/v4/messaging#MulticastMessage),
    ) (*[messaging](/firebase.google.com/go/v4/messaging).[BatchResponse](/firebase.google.com/go/v4/messaging#BatchResponse), [error](/builtin#error))

SendEachForMulticast sends the given multicast message to all the FCM registration tokens specified. 

####  func (*Client) [SendMulticastDryRun](https://github.com/appleboy/go-fcm/blob/v1.2.7/client.go#L152) ¶ added in v1.0.0
    
    
    func (c *Client) SendMulticastDryRun(
    	ctx [context](/context).[Context](/context#Context),
    	message *[messaging](/firebase.google.com/go/v4/messaging).[MulticastMessage](/firebase.google.com/go/v4/messaging#MulticastMessage),
    ) (*[messaging](/firebase.google.com/go/v4/messaging).[BatchResponse](/firebase.google.com/go/v4/messaging#BatchResponse), [error](/builtin#error))

SendEachForMulticastDryRun sends the given multicast message to all the specified FCM registration tokens in the dry run (validation only) mode. 

####  func (*Client) [SubscribeTopic](https://github.com/appleboy/go-fcm/blob/v1.2.7/client.go#L167) ¶ added in v1.0.0
    
    
    func (c *Client) SubscribeTopic(
    	ctx [context](/context).[Context](/context#Context),
    	tokens [][string](/builtin#string),
    	topic [string](/builtin#string),
    ) (*[messaging](/firebase.google.com/go/v4/messaging).[TopicManagementResponse](/firebase.google.com/go/v4/messaging#TopicManagementResponse), [error](/builtin#error))

SubscribeToTopic subscribes a list of registration tokens to a topic. 

The tokens list must not be empty, and have at most 1000 tokens. 

####  func (*Client) [UnsubscribeTopic](https://github.com/appleboy/go-fcm/blob/v1.2.7/client.go#L183) ¶ added in v1.0.0
    
    
    func (c *Client) UnsubscribeTopic(
    	ctx [context](/context).[Context](/context#Context),
    	tokens [][string](/builtin#string),
    	topic [string](/builtin#string),
    ) (*[messaging](/firebase.google.com/go/v4/messaging).[TopicManagementResponse](/firebase.google.com/go/v4/messaging#TopicManagementResponse), [error](/builtin#error))

UnsubscribeFromTopic unsubscribes a list of registration tokens from a topic. 

The tokens list must not be empty, and have at most 1000 tokens. 

####  type [Option](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L14) ¶
    
    
    type Option func(*Client) [error](/builtin#error)

Option configurates Client with defined option. 

####  func [WithCredentialsFile](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L42) ¶ added in v1.0.0
    
    
    func WithCredentialsFile(filename [string](/builtin#string)) Option

WithCredentialsFile returns a ClientOption that authenticates API calls with the given service account or refresh token JSON credentials file. 

####  func [WithCredentialsJSON](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L57) ¶ added in v1.0.0
    
    
    func WithCredentialsJSON(data [][byte](/builtin#byte)) Option

WithCredentialsJSON returns a ClientOption that authenticates API calls with the given service account or refresh token JSON credentials. 

####  func [WithCustomClientOption](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L116) ¶ added in v1.2.0
    
    
    func WithCustomClientOption(opts ...[option](/google.golang.org/api/option).[ClientOption](/google.golang.org/api/option#ClientOption)) Option

WithCustomClientOption is an option function that allows you to provide custom client options. It appends the provided custom options to the client's options list. The custom options are applied when sending requests to the FCM server. If no custom options are provided, this function does nothing. 

Parameters: 

  * opts: The custom client options to be appended to the client's options list.



Returns: 

  * An error if there was an issue appending the custom options to the client's options list, or nil otherwise.



####  func [WithDebug](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L99) ¶ added in v1.1.0
    
    
    func WithDebug(debug [bool](/builtin#bool)) Option

WithDebug returns Option to configure debug mode. 

####  func [WithEndpoint](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L66) ¶
    
    
    func WithEndpoint(endpoint [string](/builtin#string)) Option

WithEndpoint returns Option to configure endpoint. 

####  func [WithHTTPClient](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L17) ¶
    
    
    func WithHTTPClient(httpClient *[http](/net/http).[Client](/net/http#Client)) Option

WithHTTPClient returns Option to configure HTTP Client. 

####  func [WithHTTPProxy](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L25) ¶ added in v0.1.6
    
    
    func WithHTTPProxy(proxyURL [string](/builtin#string)) Option

WithHTTPProxy returns Option to configure HTTP Client proxy. 

####  func [WithProjectID](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L82) ¶ added in v1.0.0
    
    
    func WithProjectID(projectID [string](/builtin#string)) Option

WithProjectID returns Option to configure project ID. 

####  func [WithServiceAccount](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L74) ¶ added in v1.0.0
    
    
    func WithServiceAccount(serviceAccount [string](/builtin#string)) Option

WithServiceAccount returns Option to configure service account. 

####  func [WithTokenSource](https://github.com/appleboy/go-fcm/blob/v1.2.7/option.go#L91) ¶ added in v1.0.0
    
    
    func WithTokenSource(s [oauth2](/golang.org/x/oauth2).[TokenSource](/golang.org/x/oauth2#TokenSource)) Option

WithTokenSource returns a ClientOption that specifies an OAuth2 token source to be used as the basis for authentication. 
