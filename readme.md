# PCS PubSubClient

[![Go Reference](https://pkg.go.dev/badge/github.com/username/repo)](https://pkg.go.dev/github.com/username/repo)
[![Go Report Card](https://goreportcard.com/badge/github.com/username/repo)](https://goreportcard.com/report/github.com/username/repo)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/username/repo/blob/main/LICENSE)

## Overview

Internal Use of PCS Pub/Sub service.  
***Only for PCS Applications***.

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
    message.Deatail = "detail"
    err := pubsubclient.PublishMessage("topic-name", message)
    if err != nil {
        // Error Handling
    }
}
```
