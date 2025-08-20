// Create users for different databases
// This script runs when MongoDB container starts for the first time

// Switch to shopping database and create shopuser
db = db.getSiblingDB('shopping');
db.createUser({
  user: "shopuser",
  pwd: "shoppass",
  roles: [{ role: "readWrite", db: "shopping" }]
});

console.log("Created shopuser for shopping database");

// Switch to auth_db database and create authuser  
db = db.getSiblingDB('auth_db');
db.createUser({
  user: 'authuser',
  pwd: 'authpass',
  roles: [{ role: 'readWrite', db: 'auth_db' }]
});

console.log("Created authuser for auth_db database");

// Create a test collection to ensure databases exist
db.getSiblingDB('shopping').createCollection('products');
db.getSiblingDB('auth_db').createCollection('users');

console.log("MongoDB initialization completed");
