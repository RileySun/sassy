# Sassy

Sassy is an API agnostic SaaS usage tracker, authenticator, report generator, and administrative website all rolled into one!  

## Installation

Server Requirements:

*   1 mysql/mariadb server for API Server Database (Must have an active account with permissions)
*   1 mysql/mariadb server for Authentication Server Database (Must have an active account with permissions)
*   1 Docker compatible server with docker or alternative installed

  
Configure your .env file:

*   Copy the file labeled "ENV FILE" in the "setup" directory to the main directory
*   Fill out the sections labeled "API Database Credentials" & "Auth Database Credentials" with the correct information for the two servers.
*   Change the section labeled "ADMIN\_USER" & "ADMIN\_PASS" away from the default ones (Security risk if you don't)
*   Re-Save the file ".env" (This file _MUST_ be in the main directory)

  
Import the SQL files into your database:

*   Import the file labeled "API\_SETUP.sql" in the directory "setup" into the API Server Database (Or copy paste the code into a database query)
*   Import the file labeled "AUTH\_SETUP.sql" in the directory "setup" into the Authentication Server Database (Or copy paste the code into a database query)

  
Open your command line prompt and move to the main directory  
  
For running on machine:

*   docker compose up --build

  
For building a docker image

*   docker build . -t sassy-prod --target production

  

## API

The API sassy is wrapped around in this example is a database storing information on models, images, and videos.  

## Authentication

The API and the Authentication Servers have been separated for security reasons.  
The Authentication Server is capable of Checking Authorizations of Access Tokens (API calls must use Bearer Authorization Headers) & Generating new Access Tokens

## Administrative

Included in sassy is an administrative website for tracking usage and generating reports.  
By default the address for this site is "http://localhost:9090" but can be changed in the ".env" file  
The login to the site uses the "ADMIN\_USER" & "ADMIN\_PASS" values stored in the ".env" file