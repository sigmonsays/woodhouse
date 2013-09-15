woodhouse
=========

Woodhouse is an IRC bot

Features
- SSL support
- Simple HTTP API to interact with the bot
- Shell script plugins
- Quote database (eggs)
- template based quote database

Other features
- Private messages make the bot speak
- Bot greets anyone who joins
- Bot speaks egg when its nickname is mentioned
- simple YAML configuration file

HTTP API
===============

- /speak?channel=ROOM&message=MESSAGE
- /eggs
- /help

Commands
===============

- !egg add [message]
- !egg 
- !ping
- !<script> from 

Install
=========
go get github.com/sigmonsays/woodhouse/woodhouse-ircbot


