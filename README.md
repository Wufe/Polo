<p align="center">
    <img src="./.assets/logo.svg" width="50%" />
    <br />
    <br />
    <h1 align="center">Polo</h1>
    <p align="center">Serve your application development branches</p>
    <br /><br />
</p>

## What is Polo  

Polo is a **git-based reverse proxy**.

Allows you to create a web server which provides the ability to serve your web application in a specific time/branch/tag in history, using git.  

You just need to specify the git commit or the branch. A new session will be created for you to navigate into.

Although it provides HTTPS support, it is not intended to be used in production.  

***

## Use cases

It can be used for serving your staging / QA environment.  

Instead of having one application on a single host, you can let the user select which git branch to serve.  

Polo will then start the application and provide a reverse proxy feature for navigation.

***

## Getting started

- Download Polo from the release page
- Create one or more configuration files for your application
- Start Polo

***

## Configuration

You must provide at least one yaml configuration file describing your application remote and how to build and run it.  
The configuration file must be put next to the Polo executable file.  
You can find an example of a configuration file with all the options in the folder *examples*.  

***

## State diagram

This diagram represents the states walked since the request of a session to its destruction.

[![](https://mermaid.ink/img/eyJjb2RlIjoic3RhdGVEaWFncmFtLXYyXG4gICAgcmV0cnk6IFJldHJ5IGxvZ2ljXG4gICAgbm90ZSByaWdodCBvZiByZXRyeVxuICAgICAgICBUaGVyZSdzIGFsd2F5cyBhIGxpbWl0XG4gICAgICAgIGZvciByZXRyaWVzIGNvdW50XG4gICAgZW5kIG5vdGVcblxuICAgIFsqXSAtLT4gYnVpbGQ6IFVzZXIgcmVxdWVzdGVkIGEgc2Vzc2lvblxuICAgIHN0YXRlIGJ1aWxkIHtcbiAgICAgICAgWypdIC0tPiB2YWxpZGF0ZV9yZXF1ZXN0XG4gICAgICAgIHZhbGlkYXRlX3JlcXVlc3QgLS0-IHByb3Zpc2lvbl9mb2xkZXJzXG4gICAgICAgIHByb3Zpc2lvbl9mb2xkZXJzIC0tPiBnaXRfY2hlY2tvdXRcbiAgICAgICAgZ2l0X2NoZWNrb3V0IC0tPiBleGVjdXRlX2NvbW1hbmRzXG4gICAgICAgIGV4ZWN1dGVfY29tbWFuZHMgLS0-IFsqXVxuICAgIH1cbiAgICBidWlsZCAtLT4gaGVhbHRoY2hlY2sgOiBzdWNjZWVkZWRcbiAgICBcbiAgICBmYWlsOiBGYWlsIHJlYXNvblxuICAgIG5vdGUgcmlnaHQgb2YgZmFpbFxuICAgICAgICBDb21tYW5kIGZhaWx1cmUgL1xuICAgICAgICB0aW1lb3V0IHJlYWNoZWRcbiAgICBlbmQgbm90ZVxuXG4gICAgaGVhbHRoY2hlY2sgLS0-IHN0YXJ0IDogc3VjY2VlZGVkXG4gICAgaGVhbHRoY2hlY2sgLS0-IGJ1aWxkIDogZmFpbGVkIChyZXRyeSlcbiAgICBoZWFsdGhjaGVjayAtLT4gZGVzdHJveSA6IGZhaWxlZFxuXG5cbiAgICBidWlsZCAtLT4gYnVpbGQgOiBmYWlsZWQgKHJldHJ5KVxuICAgIGJ1aWxkIC0tPiBjbGVhbiA6IGZhaWxlZFxuXG4gICAgc3RhcnQgLS0-IHJ1bm5pbmdcbiAgICBub3RlIGxlZnQgb2Ygc3RhcnRcbiAgICAgICAgSGVyZSB0aGUgc2Vzc2lvbiBtYXggYWdlXG4gICAgICAgIHN0YXJ0cyB0byBkZWNyZW1lbnRcbiAgICBlbmQgbm90ZVxuXG4gICAgcnVubmluZyAtLT4gZGVzdHJveSA6IHJlcXVlc3RlZCBraWxsXG5cbiAgICBkZXN0cm95IC0tPiBjbGVhbiA6IGZhaWxlZFxuICAgIGRlc3Ryb3kgLS0-IGNsZWFuIDogc3VjY2VlZGVkXG5cbiAgICBjbGVhbiAtLT4gYnVpbGQgOiBzdGFydHVwIHJldHJ5XG5cbiAgICBjbGVhbiAtLT4gWypdIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQifSwidXBkYXRlRWRpdG9yIjpmYWxzZX0)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoic3RhdGVEaWFncmFtLXYyXG4gICAgcmV0cnk6IFJldHJ5IGxvZ2ljXG4gICAgbm90ZSByaWdodCBvZiByZXRyeVxuICAgICAgICBUaGVyZSdzIGFsd2F5cyBhIGxpbWl0XG4gICAgICAgIGZvciByZXRyaWVzIGNvdW50XG4gICAgZW5kIG5vdGVcblxuICAgIFsqXSAtLT4gYnVpbGQ6IFVzZXIgcmVxdWVzdGVkIGEgc2Vzc2lvblxuICAgIHN0YXRlIGJ1aWxkIHtcbiAgICAgICAgWypdIC0tPiB2YWxpZGF0ZV9yZXF1ZXN0XG4gICAgICAgIHZhbGlkYXRlX3JlcXVlc3QgLS0-IHByb3Zpc2lvbl9mb2xkZXJzXG4gICAgICAgIHByb3Zpc2lvbl9mb2xkZXJzIC0tPiBnaXRfY2hlY2tvdXRcbiAgICAgICAgZ2l0X2NoZWNrb3V0IC0tPiBleGVjdXRlX2NvbW1hbmRzXG4gICAgICAgIGV4ZWN1dGVfY29tbWFuZHMgLS0-IFsqXVxuICAgIH1cbiAgICBidWlsZCAtLT4gaGVhbHRoY2hlY2sgOiBzdWNjZWVkZWRcbiAgICBcbiAgICBmYWlsOiBGYWlsIHJlYXNvblxuICAgIG5vdGUgcmlnaHQgb2YgZmFpbFxuICAgICAgICBDb21tYW5kIGZhaWx1cmUgL1xuICAgICAgICB0aW1lb3V0IHJlYWNoZWRcbiAgICBlbmQgbm90ZVxuXG4gICAgaGVhbHRoY2hlY2sgLS0-IHN0YXJ0IDogc3VjY2VlZGVkXG4gICAgaGVhbHRoY2hlY2sgLS0-IGJ1aWxkIDogZmFpbGVkIChyZXRyeSlcbiAgICBoZWFsdGhjaGVjayAtLT4gZGVzdHJveSA6IGZhaWxlZFxuXG5cbiAgICBidWlsZCAtLT4gYnVpbGQgOiBmYWlsZWQgKHJldHJ5KVxuICAgIGJ1aWxkIC0tPiBjbGVhbiA6IGZhaWxlZFxuXG4gICAgc3RhcnQgLS0-IHJ1bm5pbmdcbiAgICBub3RlIGxlZnQgb2Ygc3RhcnRcbiAgICAgICAgSGVyZSB0aGUgc2Vzc2lvbiBtYXggYWdlXG4gICAgICAgIHN0YXJ0cyB0byBkZWNyZW1lbnRcbiAgICBlbmQgbm90ZVxuXG4gICAgcnVubmluZyAtLT4gZGVzdHJveSA6IHJlcXVlc3RlZCBraWxsXG5cbiAgICBkZXN0cm95IC0tPiBjbGVhbiA6IGZhaWxlZFxuICAgIGRlc3Ryb3kgLS0-IGNsZWFuIDogc3VjY2VlZGVkXG5cbiAgICBjbGVhbiAtLT4gYnVpbGQgOiBzdGFydHVwIHJldHJ5XG5cbiAgICBjbGVhbiAtLT4gWypdIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQifSwidXBkYXRlRWRpdG9yIjpmYWxzZX0)

***

## Known issues / missing features
- Add command line interface (e.g. polo reload)
- Add command to reload applications
- Add retries to startup of an application
- Print CLI errors on startup (cli & start commands)
- Add support to command concatenations (; and &&)
- Update session-helper design
- Update session page design (estimated time required, progress bar)
- Add "warmup" endpoints
- Ended sessions cleanup (folder structure)
- Add optional "copy mode" opposed to standard "clone mode" to initialize a application: copies the directory instead of cloning again
- Pruning branches does not work with embedded git client (prune is not supported)
- Add possibility to manually trigger a fetch in a git application folder
- Configuration reload
- Configuration reload via watching files
- Configuration CRUD via UI