# Family Expenses (API)

This project is the API of the [Family Expenses](https://github.com/BLACKMIDORI/family-expenses-app) application

# Setting up

### MongoDB

Configure the MongoDB properties in `src/main/resources/application.properties`

Create a database named `familyexpenses` and also a collection named `user`, add a user as follows:
```
{
    "_id":"DevelopmentUserInstance",
    "name":"Developer",
    "roles":["admin"],
}
```
