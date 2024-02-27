# katanomi common packages :package:

All shared common packages and code across katanomi repos

 - [apis](apis): common types and functions for type definitions
 - [apis/meta](apis/meta): objects, definitions and functions shared across projects (versioned)
 - [apis/validation](apis/validation): common validation methods
 - [client](client): client related functions
 - [controllers](controllers): controller methods and objects
 - [errors](error): common error functions
 - [examples](examples): examples of how to utilize this repo methods/objects
 - [hack](hack): basic repo hacking files (not a package)
 - [logging](logging): logging related
 - [maps](maps): package to manipulate maps with sortingand other methods.
 - [manager](manager): controller-runtime manager methods
 - [multicluster](multicluster): shared multicluster interfaces and implementations for client, etc.
 - [names](names): name generation releated methods (k8s.io/apiserver inspired)
 - [namespace](namespace): namespace releated methods
 - [parallel](parallel): parallel task execution implementation
 - [plugin](plugin): plugin system files and subpackages
 - [restclient](restclient): RESTful client methods
 - [scheme](scheme): scheme related methods
 - [sharedmain](sharedmain): common main functions to init components
 - [testing](testing): automated test related methods
 - [testing/framework](testing/framework): automated test framework for e2e and integration testing
 - [user](user): user matching releated functions
 - [webhook](webhook): custom webhook methods to extend current controller-runtime webhooks

## TODO

 - [ ]: implement custom validation webhook methods
 - [ ]: add more unit tests
 - [ ]: add requirements to testing/framework to automatically enable/disable tests based on setup
 - fdsa
