# goacl Documentation

Maintainer: [jiharal](https://github.com/jiharal)

![Go CI](https://github.com/tlabdotcom/goacl/actions/workflows/go.yml/badge.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/tlabdotcom/goacl)](https://goreportcard.com/report/github.com/tlabdotcom/goacl)

## Overview

The goacl library is designed to manage Access Control Lists (ACL) using Casbin for authorization, Redis for event streaming, and Bun as the ORM for database interactions. This library provides a convenient way to create and manage roles, features, sub-features, and policies in your Go applications.

## Features

- Role Management: Create, update, and delete roles.
- Feature Management: Manage features and their sub-features.
- Policy Management: Define and manage access policies.
- Event Streaming: Automatically synchronize changes using Redis.
- Casbin Integration: Use Casbin for flexible access control.
