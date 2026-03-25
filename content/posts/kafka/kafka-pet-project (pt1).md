title: Building a Kafka Pet Project (Part 1)
date: 2026-03-17
description: A series about building a simple Kafka-based event processing pipeline in Go, with producer, consumer, and schema registry integration.
tags: Kafka
series: Kafka Pet Project
series_title: Introduction
---

## Introduction

I'm taking an online course from [EngineerPro](https://engineerprogurus.com/) about building a distributed message queue project (Kafka). This is the first part of a series documenting the process of building a simple Kafka-based event processing pipeline, including a producer, consumer, and schema registry integration.

## Language Choice

First, we need to choose the programming language for our project. Several languages came to mind: Go, Rust, and Zig. Let's determine which is most suitable:

#### 1. Go

- Statically typed, compiled programming language.
- Widely supported with big community.
- Strong concurrent programming capability.

→ Easy to learn

#### 2. Rust

- Statically typed, compiled programming language (Like Go)
- Not beginner-friendly, very hard to use (ownership + lifetime is weird)
- Hard to setup for concurrent programming (especially async)

→ Hard but interesting — performance is great.

#### 3. Zig

- Statically typed, compiled programming language (Like Go and Rust).
- No garbage collector, manual memory management like Rust
- Supports both sync and async. Give you raw power, more control and better use of OS's API

→ Interesting but too niche for this project.

> **Note:** This is the first time I have heard of this programming language. When my mentor mentioned it, I was ??? (Wtf is this). Zig is still an evolving language (still stable 0.15) and its community is so small. This is something to consider.

### Conclusion

These are three programming languages that our mentor (Mr. Bao Do) recommended us to choose, two of them are offcially supported (which means we have the offcial code sample provided), these are Zig and Go. When i looked deeper into Distributed Message Queue like Kafka, RabbitMQ, Pulsar... there are more languages like Java, Elixir.... but now, i think i will do this project both in Zig and Go, cause it seems more fun to work with both a low level language and a high level language. 

## Requirements

- The outcome of out class is a small project about Distributed MQ (but i will try to improve and add more features to it)
