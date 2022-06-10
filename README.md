# For start this API must run mongo database and redis

`Create a directory with public named in root of project`

```
mkdir public
```

`Install and start Mongo with this command (in macOs):`

```
brew install mongodb-community
brew services start mongodb-community
```

`Install and start Redis with this command (in macOs):`

```
brew install redis
brew services start redis-server
```

`Install and start minio client (in macOs):`

```
brew install minio/stable/minio
minio server {{dataStoreDirectory}} --console-address ":9001"
dataStoreDirectory : Desktop/Documents/data/
```
