# CNCF Landscape API

Unofficial [CNCF Landscape](https://landscape.cncf.io/) REST API.

Official Web App -> [Link](https://landscape.cncf.io/)

## Motivation

A programmatic way to navigate the CNCF landscape. There is no way right now to query landscape app using an standard way.

Just created this little API for querying:
- a particular project
- projects in a specific category / subcategory

## How?

- Scrapping the web page
- Extracting the json data using [goquery](https://github.com/PuerkitoBio/goquery)
- Storing the data in a database -> [turso](https://turso.tech/) - libsql
- Creating a API for querying the registry of projects in CNCF

## TODO

- [x] Add simple api to get all projects 
- [x] Add query params for getting projects with name / category / subcategor 
- [x] Add update endpoint to update the db with latest projects 
- [ ] Add serverless group 
- [ ] Add wasm and other groups
