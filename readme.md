# PCS PubSubClient

[![Go Reference](https://pkg.go.dev/badge/github.com/username/repo)](https://pkg.go.dev/github.com/username/repo)
[![Go Report Card](https://goreportcard.com/badge/github.com/username/repo)](https://goreportcard.com/report/github.com/username/repo)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/username/repo/blob/main/LICENSE)

## Overview

Internal Use of PCS Pub/Sub service.  
**_Only for PCS Applications_**.

## Installation

Use `go get` to install the package:

```
go get github.com/PCS-Indonesia/pcspubsub/pubsubclient
```

## Command Message

```
{
    "command" : "insert|update|delete|notify",
    "payload" : <content>,
    "id"      : <id>,
    "detail"  : <additional data>
}
```

## Examples

### Listen to subscription

```
package main

import (
    "github.com/PCS-Indonesia/pcspubsub/pubsubclient"
)

func main() {
    pubsubclient.NewPubSubClient("project-id","your-credentials.json")
    pubsubclient.ReceiveMessages("subscription-name", func() {
        // Do something
    })
}
```

### Publish to topic

```
package main

import (
    "github.com/PCS-Indonesia/pcspubsub/pubsubclient"
)

func main() {
    pubsubclient.NewPubSubClient("project-id","your-credentials.json")
    var message pubsubclient.CommandMessage
    message.Command = "command"
    message.Payload = "your json body"
    message.ID      = "id"
    message.Detail = "detail"
    err := pubsubclient.PublishMessage("topic-name", message)
    if err != nil {
        // Error Handling
    }
}
```

### Test sample

-   Copy credential file to root directory
-   Copy `.env.example` as `.env` and fill it
-   Run pub with

```
go run sample.go pub
```

-   Run sub with

```
go run sample.go sub
```

### Message Ordering & Attribute Filter

This package already enables **message ordering** (`OrderingKey`) and **attribute filtering** (`origin`).

To use these features correctly:

-   The **subscriber** must enable **message ordering** in its subscription settings.
-   The **subscriber** should also configure an **attribute filter** to prevent receiving its own published messages.

Example filter:

```text
attributes.origin != "apps-2"
```

In this example, "apps-2" comes from the message.Detail field in the published message.
