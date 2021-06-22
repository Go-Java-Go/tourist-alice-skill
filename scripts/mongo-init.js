const username = _getEnv('MONGO_APP_USERNAME') === null || _getEnv('MONGO_APP_USERNAME') === "" ? _getEnv('MONGO_INITDB_ROOT_USERNAME') : _getEnv('MONGO_APP_USERNAME')
const password = _getEnv('MONGO_APP_PASSWORD') === null || _getEnv('MONGO_APP_PASSWORD') === "" ? _getEnv('MONGO_INITDB_ROOT_PASSWORD') : _getEnv('MONGO_APP_PASSWORD')
const database = _getEnv('MONGO_INITDB_DATABASE')


db.createUser({
    user: username,
    pwd: password,
    roles: [
        {
            role: 'dbOwner',
            db: database,
        },
    ],
});

db.createCollection("chat_state")
db.chat_state.createIndex({"user_id": 1}, {unique: true, expireAfterSeconds: 43200}) // ttl index to 12 hours

db.createCollection("users")
db.users.createIndex({"_id": 1}, {unique: true, expireAfterSeconds: 43200}) // ttl index to 12 hours
