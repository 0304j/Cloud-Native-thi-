db = db.getSiblingDB('shopping');
db.createUser({
  user: "shopuser",
  pwd: "shoppass",
  roles: [{ role: "readWrite", db: "shopping" }]
});
db = db.getSiblingDB('auth_db');
db.createUser({
  user: 'authuser',
  pwd: 'authpass',
  roles: [{ role: 'readWrite', db: 'auth_db' }]
});
