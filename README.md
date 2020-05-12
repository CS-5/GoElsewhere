# GoElsewhere ↪️
![](demo.gif)
## A self-hosted URL redirector (think [bit.ly](https://bit.ly))

_Note: There is still a lot of work to be done here. In it's current state, everything mostly works (although untested), but is lacking key functionality._

This project started as me wanting to find a self hosted alternative to the PHP-backed [YOURLS](https://yourls.org/), which is very good... but doesn't feel overly modern. Origionally I was just going to stick to an API backend and call it a day, but this seemed like a nice introductory project to building a React-based front end. Start to finish, in it's base state, this project took about a week, meaning there is a lot to change, and the React is far from perfect (I'm definitely still learning).

### Feature List

(In no specific order)

- [x] Basic API (Creating, deleting, and updating redirects) - Backend
- [ ] Improved API (Changing codes, access control?) - Backend
- [ ] Code refactoring and cleanup - Backend and Frontend
- [x] Basic Web Interface - Frontend
- [ ] Web Interface User Authentication - Backend and Frontend
- [ ] API Authentication - Backend and Frontent
- [ ] Per-Redirect Stats - Backend and Frontend
- [ ] Docker Image
- [ ] Different DB Backend(?)
- [ ] Searching - Frontend
- [ ] Better Go logging - Backend
- [ ] Better customization (Name, theme, login methods) - Frontend
- [ ] Redirect Sorting - Frontend
- [ ] Condensed and Comfortable UI Layouts - Frontend
- [ ] Better UI Look and Feel - Frontend
- [ ] Make the API more intuitive - Backend

### Enviorment Variables _YOU_ need to know

| Variable     | Default            | Description                                                                                                                                                        |
| ------------ | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| HTTP_NAME    | http://0.0.0.0     | The FQDN/IP you will use to access the web interface and where the generated codes will be valid. Must begin with HTTP/HTTPS and must not end with a trailing "/". |
| HTTP_IP      |                    | The IP to bind to. Not required if you plan on binding to all interfaces.                                                                                          |
| HTTP_PORT    | 80                 | The Port to bind the webserver to. NOTE: You will need to expose this port when creating your docker container                                                     |
| DB_DIRECTORY | data               | The location the SQLite database will be stored.                                                                                                                   |
| DEFAULT_URL  | https://google.com | The default redirect that will occur when not specifying a code or "admin" to get to the admin interface.                                                          |

### Docker

~~`docker pull reg.carsonseese.com/external/GoElsewhere`~~

~~`docker run --name GoElsewhere -e HTTP_NAME='https://my.url' -e DEFAULT_URL='https://google.com' -p 80:80 reg.carsonseese.com/external/GoElsewhere`~~
