# README

This repository is modeled loosely after a takehome task I had for an interview I did with Unity Technologies. As a Node.js developer, I did the task in Node, but I would like to learn Golang. This seemed like a perfect opportunity to leverage my REST experience so that I could learn Golang.

The task is to do the following:
1. [x] An endpoint where users can leave feedback for a specific game session
   1. [x] A user can only submit one review per game session
   2. [x] User MUST leave a rating of 1-5 if providing feedback
   3. [x] User MAY add a comment when providing feedback
   4. [x] Multiple players can rate the same session
2. [x] Folloing RESTful principles, create HTTP endpoints to allow:
   1. [x] Players to add feedback for a session
   2. [x] Ops team members to see recent feedback left by players
   3. [x] Allow filtering by rating
3. [x] Include a README that includes at least the following:
   1. [x] API Documentation
   2. [x] Instructions for launching and testing your API locally (if not the built-in scripts)
4. [ ] Bonus items
   1.  [ ] A simple front-end
   2.  [ ] Tests
   3.  [ ] Deployment scripts/tools
   4.  [ ] Authentication
   5.  [ ] User permissions

## Testing

### Starting the service locally
1. Run `make all`
2. Run `make run`
3. Begin testing

> #### Schema changes
> Whenever you have a schema change (e.g., whenever you add/remove, or otherwise update gorm models), you will likely run into issues when starting the service. **This is most-likely a result of the schema no longer matching that of the one that exists within the `test.db` file**.
> 
> The fastest way around this is to simply delete the `test.db` file and re-run `make run`, which will generate a new `test.db` file with the new schema. At this point, the project should build successfully.

### API Documentation
#### Creating test resources
* **User**
  * Send `POST` to `/users/create`
  * Does not require a post-body
* **Session**
  * Send `POST` to `/sessions/create`
* **SessionFeedback**
  * Send `POST` to `/sessions/feedback`
  * Pass the following parameters in the POST body:
    * `sessionId`: the UUID of the Session being reviewed
    * `userId`: the UUID of the user posting the feedback
    * `rating`: the rating for the Session (1-5)
    * (optional) `comment`: Optional comment for the feedback

#### Querying resources
* Get all users
  * Send `GET` to `/users`
* Get Sessions
  * Send `GET` to `/sessions`
* Get Session feedback
  * Send `GET` to `/sessions/feedback`
  * To get all feedback for a given session, send `GET` to `/sessions/feedback?sessionId=<SESSION_ID>`
  * To get all feedback with a given rating, send `GET` to `/sessions/feedback?rating=<RATING>`
  * To get all feedback for a given session with a given rating, send `GET` to `/sessions/feedback?sessionId=<SESSION_ID>&rating=<RATING>`