# Slack Nats

Interact with slack over NATS


## Quick start

Run NATS:

```
docker run -p 4222:4222 -p 8222:8222 -p 6222:6222 --name gnatsd -d nats:latest
```

Run slack-nats (you must provide the token of a valid slack user):

```
dep ensure && \
go build && \
SLACK_TOKEN=<slack_token> ./slack-nats
```

Interact with slack over nats. 
Here we get the user we are running slack-nats with to join a slack channel:

```
telnet localhost 4222
PUB slack.channel.join INBOX.1 26
{"name": "hcom-nats-test"}
```

By default slack-nats will connect to nats running on `nats://localhost:4222` - to change this set the `NATS_URL`
env variable.

## Nats Subjects

The following nats subjects are currently supported.
The message bodies for requests & responses follow the corresponding slack api message bodies as closely as possible.

### Request-Reply

#### slack.channel.join

Join the specified slack channel.
Uses the slack [channels.join](https://api.slack.com/methods/channels.join) api.
For the exact requests & responses supported see [channel.go](pkg/channel/channel.go).
Provide channel id OR name.
    
e.g. (node)

```javascript
nats.requestOne('slack.channel.join', {name: 'my-slack-channel'}, {}, 3000, resp => {
    console.log(resp)
})
```

output:

```
{ channel:
   { id: 'CDNPXK2KT',
     created: 1540570962,
     is_open: false,
     is_group: false,
     is_shared: false,
     is_im: false,
     is_ext_shared: false,
     is_org_shared: false,
     is_pending_ext_shared: false,
     is_private: false,
     is_mpim: false,
     unlinked: 0,
     name_normalized: 'my-slack-channel',
     num_members: 0,
     priority: 0,
     user: '',
     name: 'hcom-nats-test',
     creator: 'U6WDH7CCC',
     is_archived: false,
     members: [ 'U6WDH7CCC', 'U7KMBRAVB' ],
     topic:
      { value: 'Testing stuff',
        creator: 'U6WDH7CCC',
        last_set: 1540916727 },
     purpose: { value: '', creator: '', last_set: 0 },
     is_channel: true,
     is_general: false,
     is_member: true,
     locale: '' },
  err: '' }
```


#### slack.channel.leave

Leave the specified slack channel. 
Uses the slack [channels.leave](https://api.slack.com/methods/channels.leave) api.
For the exact requests & responses see [channel.go](pkg/channel/channel.go).
    
e.g. (node)

```javascript
nats.requestOne('slack.channel.leave', {id: 'CDNPXK2KT'}, {}, 3000, resp => {
    console.log(resp)
})
```

output:

```
{ not_in_channel: false, err: '' }
```

#### slack.chat.postMessage

Post a message.
Uses the slack [chat.postMessage](https://api.slack.com/methods/chat.postMessage) api.
For the exact requests & responses see [chat.go](pkg/chat/chat.go).
    
e.g. (node)

```javascript
nats.requestOne('slack.chat.postMessage', { text: 'Hello there', channel: 'CDNPXK2KT' }, {}, 3000, resp => {
    console.log(resp)
})
```

output:

```
{ channel: 'CDNPXK2KT', ts: '1541506301.003000', err: '' }
```

### Publish-Subscribe 

You can subscribe to events published over the [Slack RTM api](https://api.slack.com/rtm).
Events will be published to the NATS subject `slack.event.<event_name>` where `event_name` matches the name published [here](https://api.slack.com/events).
For example all messages sent in channels where the slacks-nats user has joined will be published to `slack.event.message`.

The following slack events are currently published to nats: 

#### slack.event.message

Subscribe to this to receive all messages from all channels where the slacks-nats user has joined. Response is the slack [message event](https://api.slack.com/events/message).
    
e.g. (node)

```javascript
nats.subscribe('slack.event.message', resp => {
    console.log(resp)
})
```

output:

```
{ type: 'message',
  channel: 'CDNPXK2KT',
  user: 'U6WDH7CCC',
  text: 'hey everyone',
  ts: '1541506728.003400',
  event_ts: '1541506728.003400',
  team: 'T09D77D4P',
  replace_original: false,
  delete_original: false }
  ...
```

# Legal
This project is available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0.html).
