# Chirpy
A program that represents a simplified twiter like backend. Users are able to regester accounts, create chirps, edit chirps, and delete their account. 

## Motivation

This project is from the boot.dev course, and designed to teach about how http servers.

## Quickstart 
Download the project and run go build -o out && ./out. This should start a new server on local host 8080. At this point you'll be able to create and receive calls to the api.

## Usage
This will require you to connect the document to the frontend of your choosing. You can see the following options that are available:
	GET /api/healthz - Confirms health of server and general functioning
	GET /admin/metrics - Provides metrics regarding how many people have viewed the home page
	POST /api/reset - Resets the metrics database
	POST /api/users - Create a new user
	POST /api/chirps - Create a new chirp
	GET /api/chirps - Receives a chirp. Optional flags include author_id={id} which will only give you one users chirps and sort={ord}, either asc or desc, which will sort the results accordingly. Default is asc
	POST /admin/reset - Resets the database to default
	GET /api/chirps/{chirpID} - Gets a single chirp at {chirpID}
	POST /api/login - Allow a user to login
	POST /api/refresh - Refresh a users long term cookie
	POST /api/revoke - Revoke a users long term cookie
	PUT /api/users - Update the users long term key
	DELETE /api/chirps/{chirpID} - Deletes a chirp. Only works if the person requesting the delete is the poster of the chirp
	POST /api/polka/webhooks - A webhook for our "payment processor" that enables you to update a users verified status

## Contributing
Contribute by forking the repo opening a pull request! All pull requests should be submited to the main branch.
