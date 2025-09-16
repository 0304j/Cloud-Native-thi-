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

// Switch to kitchen_db database and create kitchenuser  
db = db.getSiblingDB('kitchen_db');
db.createUser({
  user: 'kitchenuser',
  pwd: 'kitchenpass',
  roles: [{ role: 'readWrite', db: 'kitchen_db' }]
});

console.log("Created kitchenuser for kitchen_db database");

// Create a test collection to ensure databases exist
db.getSiblingDB('shopping').createCollection('products');
db.getSiblingDB('auth_db').createCollection('users');
db.getSiblingDB('kitchen_db').createCollection('kitchen_orders');

console.log("MongoDB initialization completed");
