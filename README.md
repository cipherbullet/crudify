# Crudify
Crudify is a simple Golang project that implements basic CRUD function using Golang as it's my first experience with this language.
I tried to use Storage interface but making some general database functions that allows you to use them for all models using context and reflection in Go!

## How to use it?
keep calm and run ``make run``! then the application starts to listen on port ``8080`` by default.
if you like to edit listen port, just edit it in the ``main.go`` file
```go
s := NewServer("0.0.0.0", "8080", pgStorage)
```
> Don't forget to put your own database credentials in the ``main.go`` file

## How to add new Model?
To add new model:
1. Define your model ``struct`` in the ``models.go`` file
2. Add your model into ``migrate.go`` file for ``AutoMigration``
3. Implement your ``Handlers`` in the ``handlers.go`` file

## How to use different database?
Here, I am using ``postgresql`` as it's my favorite relational database :) but if you like to use different databse:
1. Create separate ``.go`` file and call it with your db name in ``main`` package
2. Implement all ``StorageInterface`` methods in that file
3. Use it in the ``main.go`` to pass it to ``NewServer()``

> Just be sure you are implementing general functions to use reflection and context to get and return all models using same methods

## Contribution
Just fork the repo and make your changes in new branch starting with ``feat/your-feature-name``. create PR when you are done!