# Examples

This directory contains various examples demonstrating different features of Streamlit Go.

## Basic Example

A simple example showing basic usage of Streamlit Go.

```bash
cd basic
go run main.go
```

Then visit http://localhost:8502 in your browser.

## Session Widgets Example

An example demonstrating session-scoped widgets and multi-user support.

```bash
cd session-widgets
go run main.go
```

Then visit:
- User 1: http://localhost:8504?sessionId=user-1-session
- User 2: http://localhost:8504?sessionId=user-2-session
- Default user: http://localhost:8504?sessionId=default-session-id

## Widget Update Demo

An example demonstrating widget-level updates and deletions via HTTP POST requests.

```bash
cd widget-update-demo
go run main.go
```

Then visit http://localhost:8505 in your browser.