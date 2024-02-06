// db.auth(
//     process.env.MONGO_INITDB_ROOT_USERNAME,
//     process.env.MONGO_INITDB_ROOT_PASSWORD
var adminDb = db.getSiblingDB('admin');

// Создание пользователя с правами доступа к новой базе данных
adminDb.createUser({
    user: process.env.MONGO_DB_DEV_USERNAME,
    pwd: process.env.MONGO_DB_DEV_PASSWORD,
    roles: [
        {
            role: 'readWrite',
            db: process.env.MONGO_DB_DATABASE
        }
    ]
});

// Подключение к созданной базе данных
var appDB = db.getSiblingDB(process.env.MONGO_DB_DATABASE);

// Вставка данных в коллекцию
appDB[process.env.MONGO_DB_COLLECTION].insertOne({
    "name": "John",
    "file": 30,
    "city": "New York"
});